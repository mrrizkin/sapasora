package message

import "context"

// MessageFilter is transport-neutral. Every list operation remains scoped by
// the explicit tenant/workspace argument.
type MessageFilter struct {
	Direction            *MessageDirection
	Status               *MessageStatus
	ChannelPublicID      *string
	ContactPublicID      *string
	ConversationPublicID *string
	IdempotencyKey       *string
	ProviderMessageID    *string
}

// Repository is the persistence boundary for Track 7.1. It intentionally has
// no provider send methods, API types, or state transition operations.
type Repository interface {
	CreateMessage(context.Context, *Message) error
	GetMessageByPublicID(context.Context, string, string) (*Message, error)
	GetMessageByID(context.Context, string, uint64) (*Message, error)
	ListMessages(context.Context, string, MessageFilter) ([]*Message, error)
	UpdateMessage(context.Context, *Message) error

	CreateAttachment(context.Context, *MessageAttachment) error
	GetAttachmentByPublicID(context.Context, string, string) (*MessageAttachment, error)
	GetAttachmentByID(context.Context, string, uint64) (*MessageAttachment, error)
	ListAttachments(context.Context, string, string) ([]*MessageAttachment, error)

	CreateDeliveryRecord(context.Context, *DeliveryRecord) error
	GetDeliveryRecordByPublicID(context.Context, string, string) (*DeliveryRecord, error)
	GetDeliveryRecordByID(context.Context, string, uint64) (*DeliveryRecord, error)
	ListDeliveryRecords(context.Context, string, string) ([]*DeliveryRecord, error)

	AppendMessageEvent(context.Context, *MessageEvent) error
	ListMessageEvents(context.Context, string, string) ([]*MessageEvent, error)
}

// MessageRepository is the explicit aggregate name used by callers.
type MessageRepository = Repository

type MessageFilterByStatus = MessageFilter
