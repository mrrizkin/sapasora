// Package outbox contains the provider-neutral message outbox and worker
// foundation. It deliberately has no database, HTTP, or provider dependency.
package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"
)

// Scope identifies the tenant and workspace that own an outbox record.
type Scope struct {
	TenantID    string `json:"tenant_id"`
	WorkspaceID string `json:"workspace_id"`
}

func (s Scope) valid() bool { return s.TenantID != "" && s.WorkspaceID != "" }
func (s Scope) matches(other Scope) bool {
	return s.TenantID == other.TenantID && s.WorkspaceID == other.WorkspaceID
}

// EventStatus is the publication state of an outbox event.
type EventStatus string

const (
	EventPending   EventStatus = "pending"
	EventPublished EventStatus = "published"
)

// JobStatus is the lifecycle state of a message send job.
type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobLeased    JobStatus = "leased"
	JobAccepted  JobStatus = "accepted"
	JobSucceeded JobStatus = "succeeded"
	JobRetrying  JobStatus = "retrying"
	JobDead      JobStatus = "dead"
	JobCancelled JobStatus = "cancelled"
)

var (
	ErrNoWork          = errors.New("outbox: no work available")
	ErrNotFound        = errors.New("outbox: record not found")
	ErrConflict        = errors.New("outbox: record already exists")
	ErrInvalidScope    = errors.New("outbox: tenant and workspace are required")
	ErrNotPublished    = errors.New("outbox: event is not published")
	ErrNotOwner        = errors.New("outbox: lease is owned by another worker")
	ErrLeaseExpired    = errors.New("outbox: lease expired")
	ErrAlreadyAccepted = errors.New("outbox: provider already accepted the job")
	ErrNotReplayable   = errors.New("outbox: job is not dead-lettered")
	ErrNilContext      = errors.New("outbox: context is nil")
)

// OutboxEvent is an immutable, tenant-scoped event. Payload is never included
// in JSON diagnostics or String output; adapters can decode it at the boundary.
type OutboxEvent struct {
	ID string `json:"id"`
	// TenantID and WorkspaceID are kept as direct fields to make the
	// isolation boundary obvious to repositories and operators.
	TenantID    string      `json:"tenant_id"`
	WorkspaceID string      `json:"workspace_id"`
	Scope       Scope       `json:"-"`
	Type        string      `json:"type"`
	AggregateID string      `json:"aggregate_id,omitempty"`
	Payload     []byte      `json:"-"`
	Status      EventStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	PublishedAt *time.Time  `json:"published_at,omitempty"`
}

// SendJob is a provider-neutral delivery request. Payload is intentionally
// excluded from JSON and diagnostics to prevent message content leakage.
type SendJob struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	WorkspaceID     string     `json:"workspace_id"`
	Scope           Scope      `json:"-"`
	EventID         string     `json:"event_id"`
	Provider        string     `json:"provider"`
	Channel         string     `json:"channel"`
	Payload         []byte     `json:"-"`
	Status          JobStatus  `json:"status"`
	Attempt         int        `json:"attempt"`
	MaxAttempts     int        `json:"max_attempts"`
	NextAttemptAt   time.Time  `json:"next_attempt_at"`
	LeaseOwner      string     `json:"lease_owner,omitempty"`
	LeasedUntil     *time.Time `json:"leased_until,omitempty"`
	CancelRequested bool       `json:"cancel_requested,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	AcceptedAt      *time.Time `json:"accepted_at,omitempty"`
	FinishedAt      *time.Time `json:"finished_at,omitempty"`
	LastError       string     `json:"last_error,omitempty"`
	ReplayCount     int        `json:"replay_count"`
}

// SendResult describes the provider acceptance boundary. Accepted means the
// provider has accepted responsibility; after that point cancellation is not
// possible through this outbox.
type SendResult struct {
	Accepted bool
}

// Sender is the only worker dependency required from a future provider
// adapter. It must not log or return message payloads in errors.
type Sender interface {
	Send(context.Context, SendJob) (SendResult, error)
}

// EventStore persists events and jobs. The in-memory implementation below is
// intended for tests and local foundation work; a database adapter can satisfy
// this interface later without changing worker behavior.
type EventStore interface {
	AppendEvent(context.Context, OutboxEvent) error
	PublishEvent(context.Context, Scope, string) error
	EnqueueJob(context.Context, SendJob) error
	Claim(context.Context, Scope, string, time.Time, time.Duration) (SendJob, error)
	RenewLease(context.Context, Scope, string, string, time.Time, time.Duration) error
	AcceptJob(context.Context, Scope, string, string, time.Time) error
	RetryJob(context.Context, Scope, string, string, time.Time, time.Time, error) error
	DeadLetterJob(context.Context, Scope, string, string, time.Time, error) error
	CancelJob(context.Context, Scope, string, time.Time) error
	ReplayDeadLetter(context.Context, Scope, string, time.Time) error
	GetEvent(context.Context, Scope, string) (OutboxEvent, error)
	GetJob(context.Context, Scope, string) (SendJob, error)
}

var idCounter atomic.Uint64

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), idCounter.Add(1))
}

// NewEvent creates a pending event with a copied payload.
func NewEvent(scope Scope, eventType, aggregateID string, payload []byte, now time.Time) OutboxEvent {
	return OutboxEvent{ID: newID("evt"), TenantID: scope.TenantID, WorkspaceID: scope.WorkspaceID, Scope: scope, Type: eventType, AggregateID: aggregateID,
		Payload: append([]byte(nil), payload...), Status: EventPending, CreatedAt: now}
}

// NewSendJob creates a pending job linked to an already published event.
func NewSendJob(scope Scope, eventID, provider, channel string, payload []byte, maxAttempts int, now time.Time) SendJob {
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	return SendJob{ID: newID("job"), TenantID: scope.TenantID, WorkspaceID: scope.WorkspaceID, Scope: scope, EventID: eventID, Provider: provider, Channel: channel,
		Payload: append([]byte(nil), payload...), Status: JobPending, MaxAttempts: maxAttempts,
		NextAttemptAt: now, CreatedAt: now}
}

func validateEvent(event OutboxEvent) error {
	if !normalizeScope(event.Scope, event.TenantID, event.WorkspaceID) || event.Type == "" || event.ID == "" || event.CreatedAt.IsZero() {
		return ErrInvalidScope
	}
	if event.Status != "" && event.Status != EventPending {
		return fmt.Errorf("outbox: event must be pending when appended")
	}
	return nil
}

func validateJob(job SendJob) error {
	if !normalizeScope(job.Scope, job.TenantID, job.WorkspaceID) || job.ID == "" || job.EventID == "" || job.CreatedAt.IsZero() {
		return ErrInvalidScope
	}
	if job.MaxAttempts <= 0 {
		return fmt.Errorf("outbox: max attempts must be positive")
	}
	if job.Status != "" && job.Status != JobPending {
		return fmt.Errorf("outbox: new job must be pending")
	}
	return nil
}

func normalizeScope(scope Scope, tenantID, workspaceID string) bool {
	if scope.valid() {
		return (tenantID == "" || tenantID == scope.TenantID) &&
			(workspaceID == "" || workspaceID == scope.WorkspaceID)
	}
	return scope == (Scope{}) && tenantID != "" && workspaceID != ""
}

func eventScope(event OutboxEvent) Scope {
	if event.Scope.valid() {
		return event.Scope
	}
	return Scope{TenantID: event.TenantID, WorkspaceID: event.WorkspaceID}
}

func cloneEvent(event OutboxEvent) OutboxEvent {
	if !event.Scope.valid() {
		event.Scope = Scope{TenantID: event.TenantID, WorkspaceID: event.WorkspaceID}
	}
	event.TenantID = event.Scope.TenantID
	event.WorkspaceID = event.Scope.WorkspaceID
	event.Payload = append([]byte(nil), event.Payload...)
	if event.PublishedAt != nil {
		t := *event.PublishedAt
		event.PublishedAt = &t
	}
	return event
}

func cloneJob(job SendJob) SendJob {
	if !job.Scope.valid() {
		job.Scope = Scope{TenantID: job.TenantID, WorkspaceID: job.WorkspaceID}
	}
	job.TenantID = job.Scope.TenantID
	job.WorkspaceID = job.Scope.WorkspaceID
	job.Payload = append([]byte(nil), job.Payload...)
	if job.LeasedUntil != nil {
		t := *job.LeasedUntil
		job.LeasedUntil = &t
	}
	if job.AcceptedAt != nil {
		t := *job.AcceptedAt
		job.AcceptedAt = &t
	}
	if job.FinishedAt != nil {
		t := *job.FinishedAt
		job.FinishedAt = &t
	}
	return job
}

// MarshalSafe is useful for logs and diagnostics and guarantees that payload
// bytes are not serialized. It is intentionally small and provider-neutral.
func MarshalSafe(value any) ([]byte, error) { return json.Marshal(value) }

// RedactedError returns a diagnostic that contains no provider or message
// payload. Error type is enough for retry/DLQ inspection; detailed errors stay
// in provider-specific secure telemetry when that boundary is added.
func RedactedError(err error) string {
	if err == nil {
		return ""
	}
	return reflect.TypeOf(err).String()
}

// String returns safe diagnostic metadata without exposing payload bytes.
func (e OutboxEvent) String() string {
	encoded, err := MarshalSafe(e)
	if err != nil {
		return "outbox event"
	}
	return string(encoded)
}

// String returns safe diagnostic metadata without exposing payload bytes.
func (j SendJob) String() string {
	encoded, err := MarshalSafe(j)
	if err != nil {
		return "outbox send job"
	}
	return string(encoded)
}
