package conversation

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ConversationStatus is the provider-neutral inbox lifecycle state.
type ConversationStatus string

const (
	ConversationStatusUnknown  ConversationStatus = "unknown"
	ConversationStatusOpen     ConversationStatus = "open"
	ConversationStatusPending  ConversationStatus = "pending"
	ConversationStatusSnoozed  ConversationStatus = "snoozed"
	ConversationStatusResolved ConversationStatus = "resolved"
	ConversationStatusClosed   ConversationStatus = "closed"

	// Short names make domain filters convenient without losing the explicit
	// conversation prefix in the canonical constants above.
	ConversationOpen     = ConversationStatusOpen
	ConversationPending  = ConversationStatusPending
	ConversationSnoozed  = ConversationStatusSnoozed
	ConversationResolved = ConversationStatusResolved
	ConversationClosed   = ConversationStatusClosed

	StatusOpen     = ConversationStatusOpen
	StatusPending  = ConversationStatusPending
	StatusSnoozed  = ConversationStatusSnoozed
	StatusResolved = ConversationStatusResolved
	StatusClosed   = ConversationStatusClosed
)

func (s ConversationStatus) Valid() bool {
	switch s {
	case ConversationStatusOpen, ConversationStatusPending, ConversationStatusSnoozed, ConversationStatusResolved, ConversationStatusClosed:
		return true
	default:
		return false
	}
}

func (s ConversationStatus) String() string { return string(s) }

// ConversationPriority controls inbox ordering without coupling to a UI.
type ConversationPriority string

const (
	ConversationPriorityUnknown ConversationPriority = "unknown"
	ConversationPriorityLow     ConversationPriority = "low"
	ConversationPriorityNormal  ConversationPriority = "normal"
	ConversationPriorityHigh    ConversationPriority = "high"
	ConversationPriorityUrgent  ConversationPriority = "urgent"

	PriorityLow    = ConversationPriorityLow
	PriorityNormal = ConversationPriorityNormal
	PriorityHigh   = ConversationPriorityHigh
	PriorityUrgent = ConversationPriorityUrgent
)

func (p ConversationPriority) Valid() bool {
	switch p {
	case ConversationPriorityLow, ConversationPriorityNormal, ConversationPriorityHigh, ConversationPriorityUrgent:
		return true
	default:
		return false
	}
}

func (p ConversationPriority) String() string { return string(p) }

// ParticipantRole describes a participant without importing an agent or
// provider identity type.
type ParticipantRole string

const (
	ParticipantRoleUnknown  ParticipantRole = "unknown"
	ParticipantRoleCustomer ParticipantRole = "customer"
	ParticipantRoleAgent    ParticipantRole = "agent"
	ParticipantRoleTeam     ParticipantRole = "team"
	ParticipantRoleBot      ParticipantRole = "bot"
	ParticipantRoleSystem   ParticipantRole = "system"

	RoleCustomer = ParticipantRoleCustomer
	RoleAgent    = ParticipantRoleAgent
	RoleTeam     = ParticipantRoleTeam
	RoleBot      = ParticipantRoleBot
	RoleSystem   = ParticipantRoleSystem
)

func (r ParticipantRole) Valid() bool {
	switch r {
	case ParticipantRoleCustomer, ParticipantRoleAgent, ParticipantRoleTeam, ParticipantRoleBot, ParticipantRoleSystem:
		return true
	default:
		return false
	}
}

// Conversation is a tenant/workspace-owned communication thread. IDs are
// split so the internal ID cannot escape a backend boundary; PublicID is the
// opaque external reference used by repositories and future transports.
type Conversation struct {
	ID       uint64 `json:"-"`
	PublicID string `json:"id"`
	TenantID string `json:"-"`
	// WorkspaceID is an explicit synonym for TenantID for callers that model
	// the tenant boundary as a workspace. Repositories canonicalize both.
	WorkspaceID string `json:"-"`

	ChannelPublicID string `json:"channel_id,omitempty"`
	ContactPublicID string `json:"contact_id,omitempty"`

	Status           ConversationStatus   `json:"status"`
	Priority         ConversationPriority `json:"priority"`
	Tags             []string             `json:"tags,omitempty"`
	AssigneePublicID string               `json:"assignee_id,omitempty"`
	TeamPublicID     string               `json:"team_id,omitempty"`

	SLAFirstResponseDueAt *time.Time `json:"sla_first_response_due_at,omitempty"`
	SLAResolutionDueAt    *time.Time `json:"sla_resolution_due_at,omitempty"`
	SLAFirstRespondedAt   *time.Time `json:"sla_first_responded_at,omitempty"`
	SLAResolvedAt         *time.Time `json:"sla_resolved_at,omitempty"`
	SLABreachedAt         *time.Time `json:"sla_breached_at,omitempty"`

	LastActivityAt time.Time  `json:"last_activity_at"`
	UnreadCount    int        `json:"unread_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"-"`
}

// ScopeID returns the canonical tenant/workspace boundary.
func (c Conversation) ScopeID() string {
	if strings.TrimSpace(c.TenantID) != "" {
		return strings.TrimSpace(c.TenantID)
	}
	return strings.TrimSpace(c.WorkspaceID)
}

func validateScope(tenantID, workspaceID string) (string, error) {
	tenantID = strings.TrimSpace(tenantID)
	workspaceID = strings.TrimSpace(workspaceID)
	if tenantID != "" && workspaceID != "" && tenantID != workspaceID {
		return "", fmt.Errorf("tenant and workspace ids differ")
	}
	if tenantID == "" {
		tenantID = workspaceID
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant/workspace id is required")
	}
	return tenantID, nil
}

func (c Conversation) Valid() error {
	if _, err := validateScope(c.TenantID, c.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversation, err)
	}
	if strings.TrimSpace(c.PublicID) == "" {
		return fmt.Errorf("%w: public id is required", ErrInvalidConversation)
	}
	if !c.Status.Valid() {
		return fmt.Errorf("%w: invalid status %q", ErrInvalidConversation, c.Status)
	}
	if !c.Priority.Valid() {
		return fmt.Errorf("%w: invalid priority %q", ErrInvalidConversation, c.Priority)
	}
	if err := validateTags(c.Tags); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversation, err)
	}
	if c.UnreadCount < 0 {
		return fmt.Errorf("%w: unread count cannot be negative", ErrInvalidConversation)
	}
	if err := validateOptionalTimes(c.SLAFirstResponseDueAt, c.SLAResolutionDueAt, c.SLAFirstRespondedAt, c.SLAResolvedAt, c.SLABreachedAt, c.DeletedAt); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversation, err)
	}
	return nil
}

func (c Conversation) String() string {
	return fmt.Sprintf("Conversation{id=%s status=%s priority=%s unread=%d tags=%d}", c.PublicID, c.Status, c.Priority, c.UnreadCount, len(c.Tags))
}

// MarshalJSON is a safe diagnostic projection. Scope internals and deleted
// state are intentionally omitted; participant PII is never part of this
// aggregate projection.
func (c Conversation) MarshalJSON() ([]byte, error) {
	type projection struct {
		ID                    string               `json:"id"`
		ChannelID             string               `json:"channel_id,omitempty"`
		ContactID             string               `json:"contact_id,omitempty"`
		Status                ConversationStatus   `json:"status"`
		Priority              ConversationPriority `json:"priority"`
		Tags                  []string             `json:"tags,omitempty"`
		AssigneeID            string               `json:"assignee_id,omitempty"`
		TeamID                string               `json:"team_id,omitempty"`
		SLAFirstResponseDueAt *time.Time           `json:"sla_first_response_due_at,omitempty"`
		SLAResolutionDueAt    *time.Time           `json:"sla_resolution_due_at,omitempty"`
		SLAFirstRespondedAt   *time.Time           `json:"sla_first_responded_at,omitempty"`
		SLAResolvedAt         *time.Time           `json:"sla_resolved_at,omitempty"`
		SLABreachedAt         *time.Time           `json:"sla_breached_at,omitempty"`
		LastActivityAt        time.Time            `json:"last_activity_at"`
		UnreadCount           int                  `json:"unread_count"`
		CreatedAt             time.Time            `json:"created_at"`
		UpdatedAt             time.Time            `json:"updated_at"`
	}
	return json.Marshal(projection{c.PublicID, c.ChannelPublicID, c.ContactPublicID, c.Status, c.Priority, append([]string(nil), c.Tags...), c.AssigneePublicID, c.TeamPublicID, cloneTime(c.SLAFirstResponseDueAt), cloneTime(c.SLAResolutionDueAt), cloneTime(c.SLAFirstRespondedAt), cloneTime(c.SLAResolvedAt), cloneTime(c.SLABreachedAt), c.LastActivityAt, c.UnreadCount, c.CreatedAt, c.UpdatedAt})
}

// Participant associates a conversation with a contact or another domain
// actor. DisplayName and Address are useful to domain callers but are PII and
// are excluded from diagnostics and JSON.
type Participant struct {
	ID          uint64 `json:"-"`
	PublicID    string `json:"id"`
	TenantID    string `json:"-"`
	WorkspaceID string `json:"-"`

	ConversationID       uint64            `json:"-"`
	ConversationPublicID string            `json:"conversation_id"`
	ContactPublicID      string            `json:"contact_id,omitempty"`
	Role                 ParticipantRole   `json:"role"`
	DisplayName          string            `json:"-"`
	Address              string            `json:"-"`
	Metadata             map[string]string `json:"-"`
	JoinedAt             time.Time         `json:"joined_at"`
	LeftAt               *time.Time        `json:"left_at,omitempty"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
	DeletedAt            *time.Time        `json:"-"`
}

func (p Participant) ScopeID() string {
	if strings.TrimSpace(p.TenantID) != "" {
		return strings.TrimSpace(p.TenantID)
	}
	return strings.TrimSpace(p.WorkspaceID)
}

func (p Participant) Valid() error {
	if _, err := validateScope(p.TenantID, p.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidParticipant, err)
	}
	if strings.TrimSpace(p.PublicID) == "" || strings.TrimSpace(p.ConversationPublicID) == "" {
		return fmt.Errorf("%w: public and conversation ids are required", ErrInvalidParticipant)
	}
	if !p.Role.Valid() {
		return fmt.Errorf("%w: invalid participant role %q", ErrInvalidParticipant, p.Role)
	}
	if p.JoinedAt.IsZero() {
		return fmt.Errorf("%w: joined at is required", ErrInvalidParticipant)
	}
	if err := validateOptionalTimes(p.LeftAt, p.DeletedAt); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidParticipant, err)
	}
	return nil
}

func (p Participant) String() string {
	return fmt.Sprintf("Participant{id=%s conversation_id=%s role=%s contact_present=%t}", p.PublicID, p.ConversationPublicID, p.Role, p.ContactPublicID != "")
}

func (p Participant) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID             string          `json:"id"`
		ConversationID string          `json:"conversation_id"`
		ContactID      string          `json:"contact_id,omitempty"`
		Role           ParticipantRole `json:"role"`
		JoinedAt       time.Time       `json:"joined_at"`
		LeftAt         *time.Time      `json:"left_at,omitempty"`
		CreatedAt      time.Time       `json:"created_at"`
		UpdatedAt      time.Time       `json:"updated_at"`
	}{p.PublicID, p.ConversationPublicID, p.ContactPublicID, p.Role, p.JoinedAt, cloneTime(p.LeftAt), p.CreatedAt, p.UpdatedAt})
}

// ConversationAssignment is an appendable assignment record. Agent and team
// IDs are opaque references; their existence and authorization are integration
// concerns, not provider concerns of this package.
type ConversationAssignment struct {
	ID          uint64 `json:"-"`
	PublicID    string `json:"id"`
	TenantID    string `json:"-"`
	WorkspaceID string `json:"-"`

	ConversationID       uint64     `json:"-"`
	ConversationPublicID string     `json:"conversation_id"`
	AssigneePublicID     string     `json:"assignee_id,omitempty"`
	TeamPublicID         string     `json:"team_id,omitempty"`
	AssignedByPublicID   string     `json:"assigned_by_id,omitempty"`
	AssignedAt           time.Time  `json:"assigned_at"`
	UnassignedAt         *time.Time `json:"unassigned_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"-"`
}

// Assignment is a concise domain alias.
type Assignment = ConversationAssignment

// ConversationParticipant is the aggregate-qualified participant alias.
type ConversationParticipant = Participant

func (a ConversationAssignment) ScopeID() string {
	if strings.TrimSpace(a.TenantID) != "" {
		return strings.TrimSpace(a.TenantID)
	}
	return strings.TrimSpace(a.WorkspaceID)
}

func (a ConversationAssignment) Valid() error {
	if _, err := validateScope(a.TenantID, a.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversationAssignment, err)
	}
	if strings.TrimSpace(a.PublicID) == "" || strings.TrimSpace(a.ConversationPublicID) == "" {
		return fmt.Errorf("%w: public and conversation ids are required", ErrInvalidConversationAssignment)
	}
	if strings.TrimSpace(a.AssigneePublicID) == "" && strings.TrimSpace(a.TeamPublicID) == "" {
		return fmt.Errorf("%w: assignee or team is required", ErrInvalidConversationAssignment)
	}
	if a.AssignedAt.IsZero() {
		return fmt.Errorf("%w: assigned at is required", ErrInvalidConversationAssignment)
	}
	if err := validateOptionalTimes(a.UnassignedAt, a.DeletedAt); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversationAssignment, err)
	}
	return nil
}

func (a ConversationAssignment) String() string {
	return fmt.Sprintf("ConversationAssignment{id=%s conversation_id=%s assignee_present=%t team_present=%t active=%t}", a.PublicID, a.ConversationPublicID, a.AssigneePublicID != "", a.TeamPublicID != "", a.UnassignedAt == nil)
}

func validateTags(tags []string) error {
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return fmt.Errorf("tag cannot be empty")
		}
		if len(tag) > 64 {
			return fmt.Errorf("tag is too long")
		}
		key := strings.ToLower(tag)
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate tag")
		}
		seen[key] = struct{}{}
	}
	return nil
}

func normalizeTags(tags []string) ([]string, error) {
	result := make([]string, 0, len(tags))
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" || len(tag) > 64 {
			return nil, fmt.Errorf("invalid tag")
		}
		if _, exists := seen[tag]; exists {
			return nil, fmt.Errorf("duplicate tag")
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}
	return result, nil
}

func validateOptionalTimes(values ...*time.Time) error {
	for _, value := range values {
		if value != nil && value.IsZero() {
			return fmt.Errorf("timestamp is invalid")
		}
	}
	return nil
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	copyValue := make(map[string]string, len(values))
	for key, value := range values {
		copyValue[key] = value
	}
	return copyValue
}

func cloneConversation(value *Conversation) *Conversation {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.Tags = append([]string(nil), value.Tags...)
	copyValue.SLAFirstResponseDueAt = cloneTime(value.SLAFirstResponseDueAt)
	copyValue.SLAResolutionDueAt = cloneTime(value.SLAResolutionDueAt)
	copyValue.SLAFirstRespondedAt = cloneTime(value.SLAFirstRespondedAt)
	copyValue.SLAResolvedAt = cloneTime(value.SLAResolvedAt)
	copyValue.SLABreachedAt = cloneTime(value.SLABreachedAt)
	copyValue.DeletedAt = cloneTime(value.DeletedAt)
	return &copyValue
}

func cloneParticipant(value *Participant) *Participant {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.Metadata = cloneStringMap(value.Metadata)
	copyValue.LeftAt = cloneTime(value.LeftAt)
	copyValue.DeletedAt = cloneTime(value.DeletedAt)
	return &copyValue
}

func cloneAssignment(value *ConversationAssignment) *ConversationAssignment {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.UnassignedAt = cloneTime(value.UnassignedAt)
	copyValue.DeletedAt = cloneTime(value.DeletedAt)
	return &copyValue
}

func canonicalizeScope(tenantID, workspaceID *string) error {
	scope, err := validateScope(*tenantID, *workspaceID)
	if err != nil {
		return err
	}
	*tenantID = scope
	*workspaceID = scope
	return nil
}
