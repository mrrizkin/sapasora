package conversation

import (
	"context"
	"time"
)

// ConversationFilter is deliberately transport-neutral. Nil fields mean no
// filter; list ordering is stable by internal creation ID in the in-memory
// implementation.
type ConversationFilter struct {
	Status           *ConversationStatus
	Priority         *ConversationPriority
	ChannelPublicID  *string
	ContactPublicID  *string
	AssigneePublicID *string
	TeamPublicID     *string
	Tag              *string
	UnreadOnly       bool
	SLAOverdueAt     *time.Time
}

// Repository is the tenant/workspace-scoped persistence boundary for the
// conversation foundation. Every lookup carries the scope explicitly.
type Repository interface {
	CreateConversation(context.Context, *Conversation) error
	GetConversationByPublicID(context.Context, string, string) (*Conversation, error)
	GetConversationByID(context.Context, string, uint64) (*Conversation, error)
	ListConversations(context.Context, string, ConversationFilter) ([]*Conversation, error)
	UpdateConversation(context.Context, *Conversation) error
	DeleteConversation(context.Context, string, string) error

	CreateParticipant(context.Context, *Participant) error
	GetParticipantByPublicID(context.Context, string, string) (*Participant, error)
	GetParticipantByID(context.Context, string, uint64) (*Participant, error)
	ListParticipants(context.Context, string, string) ([]*Participant, error)
	UpdateParticipant(context.Context, *Participant) error
	DeleteParticipant(context.Context, string, string) error

	CreateAssignment(context.Context, *ConversationAssignment) error
	GetAssignmentByPublicID(context.Context, string, string) (*ConversationAssignment, error)
	GetAssignmentByID(context.Context, string, uint64) (*ConversationAssignment, error)
	ListAssignments(context.Context, string, string) ([]*ConversationAssignment, error)
	UpdateAssignment(context.Context, *ConversationAssignment) error
	DeleteAssignment(context.Context, string, string) error
}

// ConversationRepository is the descriptive interface name used by callers
// that want the aggregate name in their dependency declarations.
type ConversationRepository = Repository

// ConversationFilterByStatus is retained as a readable alias for integrations.
type ConversationFilterByStatus = ConversationFilter
