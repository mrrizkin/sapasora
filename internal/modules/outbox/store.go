package outbox

import (
	"context"
	"sort"
	"sync"
	"time"
)

// MemoryStore is a concurrency-safe, provider-neutral outbox store. It makes
// the publication and lease boundaries explicit so a durable adapter can copy
// the same state transitions later.
type MemoryStore struct {
	mu     sync.RWMutex
	events map[string]OutboxEvent
	jobs   map[string]SendJob
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{events: make(map[string]OutboxEvent), jobs: make(map[string]SendJob)}
}

func contextErr(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *MemoryStore) AppendEvent(ctx context.Context, event OutboxEvent) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if err := validateEvent(event); err != nil {
		return err
	}
	if !event.Scope.valid() {
		event.Scope = Scope{TenantID: event.TenantID, WorkspaceID: event.WorkspaceID}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[event.ID]; exists {
		return ErrConflict
	}
	if event.Status == "" {
		event.Status = EventPending
	}
	s.events[event.ID] = cloneEvent(event)
	return nil
}

// PublishEvent is the safe publication boundary. It is idempotent: a relay
// may retry it after a crash without producing a second publication state.
func (s *MemoryStore) PublishEvent(ctx context.Context, scope Scope, eventID string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if !scope.valid() {
		return ErrInvalidScope
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	event, ok := s.events[eventID]
	if !ok || !event.Scope.matches(scope) {
		return ErrNotFound
	}
	if event.Status == EventPublished {
		return nil
	}
	now := time.Now()
	event.Status = EventPublished
	event.PublishedAt = &now
	s.events[eventID] = cloneEvent(event)
	return nil
}

func (s *MemoryStore) EnqueueJob(ctx context.Context, job SendJob) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if err := validateJob(job); err != nil {
		return err
	}
	if !job.Scope.valid() {
		job.Scope = Scope{TenantID: job.TenantID, WorkspaceID: job.WorkspaceID}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.jobs[job.ID]; exists {
		return ErrConflict
	}
	event, ok := s.events[job.EventID]
	if !ok {
		return ErrNotFound
	}
	if !event.Scope.matches(job.Scope) {
		return ErrInvalidScope
	}
	if event.Status != EventPublished {
		return ErrNotPublished
	}
	if job.Status == "" {
		job.Status = JobPending
	}
	s.jobs[job.ID] = cloneJob(job)
	return nil
}

func scopeAllowed(requested, actual Scope) bool {
	return requested.valid() && requested.matches(actual)
}

func (s *MemoryStore) Claim(ctx context.Context, scope Scope, owner string, now time.Time, lease time.Duration) (SendJob, error) {
	if err := contextErr(ctx); err != nil {
		return SendJob{}, err
	}
	if !scope.valid() {
		return SendJob{}, ErrInvalidScope
	}
	if owner == "" || lease <= 0 {
		return SendJob{}, ErrNotOwner
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	ids := make([]string, 0, len(s.jobs))
	for id := range s.jobs {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		a, b := s.jobs[ids[i]], s.jobs[ids[j]]
		if a.NextAttemptAt.Equal(b.NextAttemptAt) {
			return a.ID < b.ID
		}
		return a.NextAttemptAt.Before(b.NextAttemptAt)
	})
	for _, id := range ids {
		job := s.jobs[id]
		if !scopeAllowed(scope, job.Scope) {
			continue
		}
		if job.Status == JobLeased {
			if job.LeasedUntil == nil || job.LeasedUntil.After(now) {
				continue
			}
			// An expired lease returns the job to the pending state before
			// this claimant owns it. A cancellation request that outlives a
			// crashed worker must not leave the job stuck forever.
			job.LeaseOwner = ""
			job.LeasedUntil = nil
			if job.CancelRequested {
				job.Status = JobCancelled
				job.FinishedAt = &now
				s.jobs[id] = cloneJob(job)
				continue
			}
			job.Status = JobPending
		}
		if job.CancelRequested || (job.Status != JobPending && job.Status != JobRetrying) {
			continue
		}
		if job.NextAttemptAt.After(now) {
			continue
		}
		job.Attempt++
		job.Status = JobLeased
		job.LeaseOwner = owner
		until := now.Add(lease)
		job.LeasedUntil = &until
		s.jobs[id] = cloneJob(job)
		return cloneJob(job), nil
	}
	return SendJob{}, ErrNoWork
}

func (s *MemoryStore) getOwnedJobLocked(scope Scope, jobID, owner string, now time.Time) (SendJob, error) {
	if !scope.valid() {
		return SendJob{}, ErrInvalidScope
	}
	job, ok := s.jobs[jobID]
	if !ok || !scopeAllowed(scope, job.Scope) {
		return SendJob{}, ErrNotFound
	}
	if job.Status == JobAccepted || job.Status == JobSucceeded {
		return SendJob{}, ErrAlreadyAccepted
	}
	if job.Status != JobLeased || job.LeaseOwner != owner {
		return SendJob{}, ErrNotOwner
	}
	if job.LeasedUntil == nil || !job.LeasedUntil.After(now) {
		return SendJob{}, ErrLeaseExpired
	}
	return job, nil
}

func (s *MemoryStore) RenewLease(ctx context.Context, scope Scope, jobID, owner string, now time.Time, lease time.Duration) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if lease <= 0 {
		return ErrLeaseExpired
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	job, err := s.getOwnedJobLocked(scope, jobID, owner, now)
	if err != nil {
		return err
	}
	until := now.Add(lease)
	job.LeasedUntil = &until
	s.jobs[jobID] = cloneJob(job)
	return nil
}

func (s *MemoryStore) AcceptJob(ctx context.Context, scope Scope, jobID, owner string, now time.Time) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	job, err := s.getOwnedJobLocked(scope, jobID, owner, now)
	if err != nil {
		return err
	}
	job.Status = JobAccepted
	job.LeaseOwner = ""
	job.LeasedUntil = nil
	job.CancelRequested = false
	job.AcceptedAt = &now
	job.FinishedAt = &now
	s.jobs[jobID] = cloneJob(job)
	return nil
}

func (s *MemoryStore) RetryJob(ctx context.Context, scope Scope, jobID, owner string, now, next time.Time, cause error) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	job, err := s.getOwnedJobLocked(scope, jobID, owner, now)
	if err != nil {
		return err
	}
	if job.CancelRequested {
		job.Status = JobCancelled
		job.FinishedAt = &now
	} else {
		job.Status = JobRetrying
		job.NextAttemptAt = next
	}
	job.LastError = RedactedError(cause)
	job.LeaseOwner = ""
	job.LeasedUntil = nil
	s.jobs[jobID] = cloneJob(job)
	return nil
}

func (s *MemoryStore) DeadLetterJob(ctx context.Context, scope Scope, jobID, owner string, now time.Time, cause error) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	job, err := s.getOwnedJobLocked(scope, jobID, owner, now)
	if err != nil {
		return err
	}
	if job.CancelRequested {
		job.Status = JobCancelled
	} else {
		job.Status = JobDead
	}
	job.FinishedAt = &now
	job.LastError = RedactedError(cause)
	job.LeaseOwner = ""
	job.LeasedUntil = nil
	s.jobs[jobID] = cloneJob(job)
	return nil
}

func (s *MemoryStore) CancelJob(ctx context.Context, scope Scope, jobID string, now time.Time) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if !scope.valid() {
		return ErrInvalidScope
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[jobID]
	if !ok || !scopeAllowed(scope, job.Scope) {
		return ErrNotFound
	}
	switch job.Status {
	case JobPending, JobRetrying:
		job.Status = JobCancelled
		job.FinishedAt = &now
	case JobLeased:
		// The worker resolves this request after the provider call. This
		// closes the cancellation race at the acceptance boundary.
		job.CancelRequested = true
	case JobAccepted, JobSucceeded:
		return ErrAlreadyAccepted
	case JobDead, JobCancelled:
		return nil
	}
	s.jobs[jobID] = cloneJob(job)
	return nil
}

func (s *MemoryStore) ReplayDeadLetter(ctx context.Context, scope Scope, jobID string, now time.Time) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if !scope.valid() {
		return ErrInvalidScope
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[jobID]
	if !ok || !scopeAllowed(scope, job.Scope) {
		return ErrNotFound
	}
	if job.Status != JobDead {
		return ErrNotReplayable
	}
	job.Status = JobPending
	job.Attempt = 0
	job.NextAttemptAt = now
	job.FinishedAt = nil
	job.LastError = ""
	job.CancelRequested = false
	job.ReplayCount++
	s.jobs[jobID] = cloneJob(job)
	return nil
}

func (s *MemoryStore) GetEvent(ctx context.Context, scope Scope, eventID string) (OutboxEvent, error) {
	if err := contextErr(ctx); err != nil {
		return OutboxEvent{}, err
	}
	if !scope.valid() {
		return OutboxEvent{}, ErrInvalidScope
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[eventID]
	if !ok || !scopeAllowed(scope, event.Scope) {
		return OutboxEvent{}, ErrNotFound
	}
	return cloneEvent(event), nil
}

func (s *MemoryStore) GetJob(ctx context.Context, scope Scope, jobID string) (SendJob, error) {
	if err := contextErr(ctx); err != nil {
		return SendJob{}, err
	}
	if !scope.valid() {
		return SendJob{}, ErrInvalidScope
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[jobID]
	if !ok || !scopeAllowed(scope, job.Scope) {
		return SendJob{}, ErrNotFound
	}
	return cloneJob(job), nil
}

// DeadLetters returns safe metadata for operator tooling; payload bytes remain
// in the returned value for the adapter but are excluded from JSON diagnostics.
func (s *MemoryStore) DeadLetters(ctx context.Context, scope Scope) ([]SendJob, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if !scope.valid() {
		return nil, ErrInvalidScope
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobs := make([]SendJob, 0)
	for _, job := range s.jobs {
		if job.Status == JobDead && scopeAllowed(scope, job.Scope) {
			jobs = append(jobs, cloneJob(job))
		}
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].ID < jobs[j].ID })
	return jobs, nil
}
