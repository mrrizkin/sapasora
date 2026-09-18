package message

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// MessageDirection describes whether a message entered or left the platform.
type MessageDirection string

const (
	MessageDirectionUnknown  MessageDirection = "unknown"
	MessageDirectionInbound  MessageDirection = "inbound"
	MessageDirectionOutbound MessageDirection = "outbound"

	DirectionInbound  = MessageDirectionInbound
	DirectionOutbound = MessageDirectionOutbound
	Inbound           = MessageDirectionInbound
	Outbound          = MessageDirectionOutbound
)

func (d MessageDirection) Valid() bool {
	return d == MessageDirectionInbound || d == MessageDirectionOutbound
}

func (d MessageDirection) String() string { return string(d) }

// MessageStatus is a stored observation, not a state-machine transition API.
// Transition validation belongs to Track 7.2.
type MessageStatus string

const (
	MessageStatusUnknown    MessageStatus = "unknown"
	MessageStatusAccepted   MessageStatus = "accepted"
	MessageStatusQueued     MessageStatus = "queued"
	MessageStatusSending    MessageStatus = "sending"
	MessageStatusSent       MessageStatus = "sent"
	MessageStatusDelivered  MessageStatus = "delivered"
	MessageStatusRead       MessageStatus = "read"
	MessageStatusFailed     MessageStatus = "failed"
	MessageStatusRetrying   MessageStatus = "retrying"
	MessageStatusDeadLetter MessageStatus = "dead_letter"
	MessageStatusCanceled   MessageStatus = "canceled"

	StatusAccepted   = MessageStatusAccepted
	StatusQueued     = MessageStatusQueued
	StatusSending    = MessageStatusSending
	StatusSent       = MessageStatusSent
	StatusDelivered  = MessageStatusDelivered
	StatusRead       = MessageStatusRead
	StatusFailed     = MessageStatusFailed
	StatusRetrying   = MessageStatusRetrying
	StatusDeadLetter = MessageStatusDeadLetter
	StatusCanceled   = MessageStatusCanceled
)

func (s MessageStatus) Valid() bool {
	switch s {
	case MessageStatusAccepted, MessageStatusQueued, MessageStatusSending, MessageStatusSent,
		MessageStatusDelivered, MessageStatusRead, MessageStatusFailed, MessageStatusRetrying,
		MessageStatusDeadLetter, MessageStatusCanceled:
		return true
	default:
		return false
	}
}

func (s MessageStatus) String() string { return string(s) }

// DeliveryStatus intentionally uses the same vocabulary as MessageStatus. A
// delivery record is a provider observation; this package does not validate
// transitions between observations.
type DeliveryStatus = MessageStatus

const (
	DeliveryStatusAccepted  = MessageStatusAccepted
	DeliveryStatusSent      = MessageStatusSent
	DeliveryStatusDelivered = MessageStatusDelivered
	DeliveryStatusRead      = MessageStatusRead
	DeliveryStatusFailed    = MessageStatusFailed
)

// ContentKind identifies normalized content without importing a provider SDK.
type ContentKind string

const (
	ContentKindUnknown  ContentKind = "unknown"
	ContentKindText     ContentKind = "text"
	ContentKindTemplate ContentKind = "template"
	ContentKindImage    ContentKind = "image"
	ContentKindAudio    ContentKind = "audio"
	ContentKindVideo    ContentKind = "video"
	ContentKindDocument ContentKind = "document"
	ContentKindSticker  ContentKind = "sticker"
	ContentKindLocation ContentKind = "location"
	ContentKindContact  ContentKind = "contact"
	ContentKindPoll     ContentKind = "poll"
)

func (k ContentKind) Valid() bool {
	switch k {
	case ContentKindText, ContentKindTemplate, ContentKindImage, ContentKindAudio, ContentKindVideo,
		ContentKindDocument, ContentKindSticker, ContentKindLocation, ContentKindContact, ContentKindPoll:
		return true
	default:
		return false
	}
}

// NormalizedContentMetadata contains non-secret facts about normalized content.
type NormalizedContentMetadata struct {
	Kind            ContentKind `json:"kind"`
	CharacterCount  int         `json:"character_count"`
	ByteCount       int         `json:"byte_count"`
	AttachmentCount int         `json:"attachment_count"`
	HasText         bool        `json:"has_text"`
	NormalizedHash  string      `json:"normalized_hash,omitempty"`
	Language        string      `json:"language,omitempty"`
}

// ContentMetadata is a concise compatibility name for normalized content
// metadata.
type ContentMetadata = NormalizedContentMetadata

func (m NormalizedContentMetadata) Valid() error {
	if m.Kind != "" && m.Kind != ContentKindUnknown && !m.Kind.Valid() {
		return fmt.Errorf("invalid content kind %q", m.Kind)
	}
	if m.CharacterCount < 0 || m.ByteCount < 0 || m.AttachmentCount < 0 {
		return fmt.Errorf("content counts cannot be negative")
	}
	if m.CharacterCount > 0 && m.ByteCount == 0 {
		return fmt.Errorf("content byte count is required")
	}
	return nil
}

// AttachmentKind describes metadata about media without carrying media bytes.
type AttachmentKind string

const (
	AttachmentKindUnknown  AttachmentKind = "unknown"
	AttachmentKindImage    AttachmentKind = "image"
	AttachmentKindAudio    AttachmentKind = "audio"
	AttachmentKindVideo    AttachmentKind = "video"
	AttachmentKindDocument AttachmentKind = "document"
	AttachmentKindSticker  AttachmentKind = "sticker"
)

func (k AttachmentKind) Valid() bool {
	switch k {
	case AttachmentKindImage, AttachmentKindAudio, AttachmentKindVideo, AttachmentKindDocument, AttachmentKindSticker:
		return true
	default:
		return false
	}
}

// Message is the provider-neutral logical message record. TenantID and
// WorkspaceID are aliases for the same isolation boundary; repositories
// canonicalize both to the same value.
type Message struct {
	ID          uint64 `json:"-"`
	PublicID    string `json:"id"`
	TenantID    string `json:"-"`
	WorkspaceID string `json:"-"`

	Direction MessageDirection `json:"direction"`
	Status    MessageStatus    `json:"status"`

	Provider             string `json:"-"`
	ProviderMessageID    string `json:"-"`
	IdempotencyKey       string `json:"-"`
	RequestID            string `json:"-"`
	ChannelPublicID      string `json:"channel_id,omitempty"`
	ContactPublicID      string `json:"contact_id,omitempty"`
	ConversationPublicID string `json:"conversation_id,omitempty"`

	// NormalizedContent is retained for domain use and deliberately excluded
	// from JSON/String diagnostics.
	NormalizedContent string `json:"-"`
	// Content is a compatibility input alias. Repositories canonicalize it into
	// NormalizedContent and never expose either value in diagnostics.
	Content         string                    `json:"-"`
	ContentMetadata NormalizedContentMetadata `json:"content_metadata"`
	// NormalizedContentMetadata is a descriptive compatibility alias. New code
	// should use ContentMetadata.
	NormalizedContentMetadata NormalizedContentMetadata `json:"-"`
	// ProviderMetadata is accepted as an input alias and canonicalized into the
	// redacted map. Neither map is emitted by safe diagnostics.
	ProviderMetadata         map[string]string `json:"-"`
	RedactedProviderMetadata map[string]string `json:"-"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

func (m Message) ScopeID() string {
	if strings.TrimSpace(m.TenantID) != "" {
		return strings.TrimSpace(m.TenantID)
	}
	return strings.TrimSpace(m.WorkspaceID)
}

func (m Message) Valid() error {
	if _, err := validateScope(m.TenantID, m.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	if strings.TrimSpace(m.PublicID) == "" {
		return fmt.Errorf("%w: public id is required", ErrInvalidMessage)
	}
	if !m.Direction.Valid() {
		return fmt.Errorf("%w: invalid direction %q", ErrInvalidMessage, m.Direction)
	}
	if !m.Status.Valid() {
		return fmt.Errorf("%w: invalid status %q", ErrInvalidMessage, m.Status)
	}
	if err := validateReference(m.Provider, "provider"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	if err := validateBoundedReference(m.ProviderMessageID, "provider message id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	if err := validateBoundedReference(m.IdempotencyKey, "idempotency key", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	if err := validateBoundedReference(m.RequestID, "request id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	for name, value := range map[string]string{
		"channel public id":      m.ChannelPublicID,
		"contact public id":      m.ContactPublicID,
		"conversation public id": m.ConversationPublicID,
	} {
		if err := validateBoundedReference(value, name, 256); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
		}
	}
	metadata := m.ContentMetadata
	if metadata == (NormalizedContentMetadata{}) {
		metadata = m.NormalizedContentMetadata
	}
	if err := metadata.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	if m.DeletedAt != nil && m.DeletedAt.IsZero() {
		return fmt.Errorf("%w: deleted timestamp is invalid", ErrInvalidMessage)
	}
	return nil
}

func (m Message) String() string {
	return fmt.Sprintf("Message{id=%s direction=%s status=%s content_present=%t attachments=%d}", m.PublicID, m.Direction, m.Status, m.NormalizedContent != "" || m.Content != "", m.ContentMetadata.AttachmentCount)
}

func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID             string                    `json:"id"`
		Direction      MessageDirection          `json:"direction"`
		Status         MessageStatus             `json:"status"`
		ChannelID      string                    `json:"channel_id,omitempty"`
		ContactID      string                    `json:"contact_id,omitempty"`
		ConversationID string                    `json:"conversation_id,omitempty"`
		Content        NormalizedContentMetadata `json:"content_metadata"`
		CreatedAt      time.Time                 `json:"created_at"`
		UpdatedAt      time.Time                 `json:"updated_at"`
	}{m.PublicID, m.Direction, m.Status, m.ChannelPublicID, m.ContactPublicID, m.ConversationPublicID, effectiveContentMetadata(m), m.CreatedAt, m.UpdatedAt})
}

// MessageAttachment stores only metadata and opaque references to media.
type MessageAttachment struct {
	ID          uint64 `json:"-"`
	PublicID    string `json:"id"`
	TenantID    string `json:"-"`
	WorkspaceID string `json:"-"`

	MessageID       uint64            `json:"-"`
	MessagePublicID string            `json:"message_id"`
	Kind            AttachmentKind    `json:"kind"`
	FileName        string            `json:"-"`
	MIMEType        string            `json:"mime_type,omitempty"`
	SizeBytes       int64             `json:"size_bytes"`
	StorageRef      string            `json:"-"`
	Checksum        string            `json:"-"`
	ProviderMediaID string            `json:"-"`
	Metadata        map[string]string `json:"-"`
	CreatedAt       time.Time         `json:"created_at"`
}

func (a MessageAttachment) ScopeID() string {
	if strings.TrimSpace(a.TenantID) != "" {
		return strings.TrimSpace(a.TenantID)
	}
	return strings.TrimSpace(a.WorkspaceID)
}

func (a MessageAttachment) Valid() error {
	if _, err := validateScope(a.TenantID, a.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidAttachment, err)
	}
	if strings.TrimSpace(a.PublicID) == "" || strings.TrimSpace(a.MessagePublicID) == "" {
		return fmt.Errorf("%w: public and message ids are required", ErrInvalidAttachment)
	}
	if err := validateBoundedReference(a.MessagePublicID, "message public id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidAttachment, err)
	}
	if !a.Kind.Valid() {
		return fmt.Errorf("%w: invalid attachment kind %q", ErrInvalidAttachment, a.Kind)
	}
	if a.SizeBytes < 0 {
		return fmt.Errorf("%w: size cannot be negative", ErrInvalidAttachment)
	}
	if err := validateBoundedReference(a.MIMEType, "mime type", 128); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidAttachment, err)
	}
	if a.CreatedAt.IsZero() {
		return fmt.Errorf("%w: created at is required", ErrInvalidAttachment)
	}
	return nil
}

func (a MessageAttachment) String() string {
	return fmt.Sprintf("MessageAttachment{id=%s message_id=%s kind=%s size_bytes=%d}", a.PublicID, a.MessagePublicID, a.Kind, a.SizeBytes)
}

func (a MessageAttachment) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string         `json:"id"`
		MessageID string         `json:"message_id"`
		Kind      AttachmentKind `json:"kind"`
		MIMEType  string         `json:"mime_type,omitempty"`
		SizeBytes int64          `json:"size_bytes"`
		CreatedAt time.Time      `json:"created_at"`
	}{a.PublicID, a.MessagePublicID, a.Kind, a.MIMEType, a.SizeBytes, a.CreatedAt})
}

// AttachmentMetadata is the aggregate-qualified name for attachment metadata.
type AttachmentMetadata = MessageAttachment

// MessageAttachmentMetadata is a descriptive alias used by integrations.
type MessageAttachmentMetadata = MessageAttachment

// DeliveryRecord records one provider-neutral delivery observation.
type DeliveryRecord struct {
	ID          uint64 `json:"-"`
	PublicID    string `json:"id"`
	TenantID    string `json:"-"`
	WorkspaceID string `json:"-"`

	MessageID                uint64            `json:"-"`
	MessagePublicID          string            `json:"message_id"`
	ChannelPublicID          string            `json:"channel_id,omitempty"`
	Provider                 string            `json:"-"`
	ProviderMessageID        string            `json:"-"`
	Status                   DeliveryStatus    `json:"status"`
	ErrorCode                string            `json:"error_code,omitempty"`
	RedactedProviderMetadata map[string]string `json:"-"`
	OccurredAt               time.Time         `json:"occurred_at"`
	RecordedAt               time.Time         `json:"recorded_at"`
}

func (d DeliveryRecord) ScopeID() string {
	if strings.TrimSpace(d.TenantID) != "" {
		return strings.TrimSpace(d.TenantID)
	}
	return strings.TrimSpace(d.WorkspaceID)
}

func (d DeliveryRecord) Valid() error {
	if _, err := validateScope(d.TenantID, d.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDeliveryRecord, err)
	}
	if strings.TrimSpace(d.PublicID) == "" || strings.TrimSpace(d.MessagePublicID) == "" {
		return fmt.Errorf("%w: public and message ids are required", ErrInvalidDeliveryRecord)
	}
	if !d.Status.Valid() {
		return fmt.Errorf("%w: invalid delivery status %q", ErrInvalidDeliveryRecord, d.Status)
	}
	if err := validateReference(d.Provider, "provider"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDeliveryRecord, err)
	}
	if err := validateBoundedReference(d.MessagePublicID, "message public id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDeliveryRecord, err)
	}
	if err := validateBoundedReference(d.ChannelPublicID, "channel public id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDeliveryRecord, err)
	}
	if err := validateBoundedReference(d.ProviderMessageID, "provider message id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDeliveryRecord, err)
	}
	if d.OccurredAt.IsZero() || d.RecordedAt.IsZero() {
		return fmt.Errorf("%w: delivery timestamps are required", ErrInvalidDeliveryRecord)
	}
	return nil
}

func (d DeliveryRecord) String() string {
	return fmt.Sprintf("DeliveryRecord{id=%s message_id=%s status=%s provider_present=%t}", d.PublicID, d.MessagePublicID, d.Status, d.Provider != "")
}

func (d DeliveryRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID         string         `json:"id"`
		MessageID  string         `json:"message_id"`
		ChannelID  string         `json:"channel_id,omitempty"`
		Status     DeliveryStatus `json:"status"`
		ErrorCode  string         `json:"error_code,omitempty"`
		OccurredAt time.Time      `json:"occurred_at"`
		RecordedAt time.Time      `json:"recorded_at"`
	}{d.PublicID, d.MessagePublicID, d.ChannelPublicID, d.Status, d.ErrorCode, d.OccurredAt, d.RecordedAt})
}

// MessageDelivery is the aggregate-qualified name for a delivery record.
type MessageDelivery = DeliveryRecord

// MessageEvent is immutable history for a logical message. Append-only
// storage is enforced by the repository API, which has no update/delete path.
type MessageEvent struct {
	ID          uint64 `json:"-"`
	PublicID    string `json:"id"`
	TenantID    string `json:"-"`
	WorkspaceID string `json:"-"`

	MessageID                uint64            `json:"-"`
	MessagePublicID          string            `json:"message_id"`
	Sequence                 uint64            `json:"sequence"`
	Type                     MessageEventType  `json:"type"`
	Status                   MessageStatus     `json:"status,omitempty"`
	ProviderMessageID        string            `json:"-"`
	RequestID                string            `json:"-"`
	RedactedProviderMetadata map[string]string `json:"-"`
	OccurredAt               time.Time         `json:"occurred_at"`
	RecordedAt               time.Time         `json:"recorded_at"`
}

// MessageEventType names an immutable observation. It does not define valid
// state transitions; that is intentionally deferred to Track 7.2.
type MessageEventType string

const (
	MessageEventCreated    MessageEventType = "created"
	MessageEventAccepted   MessageEventType = "accepted"
	MessageEventQueued     MessageEventType = "queued"
	MessageEventSending    MessageEventType = "sending"
	MessageEventSent       MessageEventType = "sent"
	MessageEventDelivered  MessageEventType = "delivered"
	MessageEventRead       MessageEventType = "read"
	MessageEventFailed     MessageEventType = "failed"
	MessageEventRetrying   MessageEventType = "retrying"
	MessageEventDeadLetter MessageEventType = "dead_letter"
	MessageEventCanceled   MessageEventType = "canceled"

	MessageEventTypeCreated    = MessageEventCreated
	MessageEventTypeAccepted   = MessageEventAccepted
	MessageEventTypeQueued     = MessageEventQueued
	MessageEventTypeSending    = MessageEventSending
	MessageEventTypeSent       = MessageEventSent
	MessageEventTypeDelivered  = MessageEventDelivered
	MessageEventTypeRead       = MessageEventRead
	MessageEventTypeFailed     = MessageEventFailed
	MessageEventTypeRetrying   = MessageEventRetrying
	MessageEventTypeDeadLetter = MessageEventDeadLetter
	MessageEventTypeCanceled   = MessageEventCanceled
)

func (t MessageEventType) Valid() bool {
	switch t {
	case MessageEventCreated, MessageEventAccepted, MessageEventQueued, MessageEventSending, MessageEventSent,
		MessageEventDelivered, MessageEventRead, MessageEventFailed, MessageEventRetrying, MessageEventDeadLetter, MessageEventCanceled:
		return true
	default:
		return false
	}
}

func (e MessageEvent) ScopeID() string {
	if strings.TrimSpace(e.TenantID) != "" {
		return strings.TrimSpace(e.TenantID)
	}
	return strings.TrimSpace(e.WorkspaceID)
}

func (e MessageEvent) Valid() error {
	if _, err := validateScope(e.TenantID, e.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessageEvent, err)
	}
	if strings.TrimSpace(e.PublicID) == "" || strings.TrimSpace(e.MessagePublicID) == "" {
		return fmt.Errorf("%w: public and message ids are required", ErrInvalidMessageEvent)
	}
	if e.Sequence == 0 {
		return fmt.Errorf("%w: sequence is required", ErrInvalidMessageEvent)
	}
	if !e.Type.Valid() {
		return fmt.Errorf("%w: invalid event type %q", ErrInvalidMessageEvent, e.Type)
	}
	if e.Status != "" && e.Status != MessageStatusUnknown && !e.Status.Valid() {
		return fmt.Errorf("%w: invalid event status %q", ErrInvalidMessageEvent, e.Status)
	}
	if err := validateBoundedReference(e.MessagePublicID, "message public id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessageEvent, err)
	}
	if err := validateBoundedReference(e.ProviderMessageID, "provider message id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessageEvent, err)
	}
	if err := validateBoundedReference(e.RequestID, "request id", 256); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessageEvent, err)
	}
	if e.OccurredAt.IsZero() || e.RecordedAt.IsZero() {
		return fmt.Errorf("%w: event timestamps are required", ErrInvalidMessageEvent)
	}
	return nil
}

func (e MessageEvent) String() string {
	return fmt.Sprintf("MessageEvent{id=%s message_id=%s sequence=%d type=%s status=%s}", e.PublicID, e.MessagePublicID, e.Sequence, e.Type, e.Status)
}

func (e MessageEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID         string           `json:"id"`
		MessageID  string           `json:"message_id"`
		Sequence   uint64           `json:"sequence"`
		Type       MessageEventType `json:"type"`
		Status     MessageStatus    `json:"status,omitempty"`
		OccurredAt time.Time        `json:"occurred_at"`
		RecordedAt time.Time        `json:"recorded_at"`
	}{e.PublicID, e.MessagePublicID, e.Sequence, e.Type, e.Status, e.OccurredAt, e.RecordedAt})
}

func validateScope(tenantID, workspaceID string) (string, error) {
	tenantID, workspaceID = strings.TrimSpace(tenantID), strings.TrimSpace(workspaceID)
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

func validateReference(value, name string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return validateBoundedReference(value, name, 256)
}

func validateBoundedReference(value, name string, max int) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if len(value) > max {
		return fmt.Errorf("%s is too long", name)
	}
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("%s contains a newline", name)
	}
	return nil
}

func effectiveContentMetadata(m Message) NormalizedContentMetadata {
	metadata := m.ContentMetadata
	if metadata == (NormalizedContentMetadata{}) {
		metadata = m.NormalizedContentMetadata
	}
	return metadata
}

func normalizeContent(value string) string { return strings.TrimSpace(value) }

func contentMetadata(value string, current NormalizedContentMetadata) NormalizedContentMetadata {
	if current == (NormalizedContentMetadata{}) && value != "" {
		current = NormalizedContentMetadata{Kind: ContentKindText, HasText: true}
	}
	if value != "" {
		current.CharacterCount = utf8.RuneCountInString(value)
		current.ByteCount = len(value)
		current.HasText = true
	}
	return current
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

func cloneMessage(value *Message) *Message {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.ContentMetadata = effectiveContentMetadata(*value)
	copyValue.NormalizedContentMetadata = copyValue.ContentMetadata
	copyValue.ProviderMetadata = cloneStringMap(value.ProviderMetadata)
	copyValue.RedactedProviderMetadata = cloneStringMap(value.RedactedProviderMetadata)
	copyValue.DeletedAt = cloneTime(value.DeletedAt)
	return &copyValue
}

func cloneAttachment(value *MessageAttachment) *MessageAttachment {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.Metadata = cloneStringMap(value.Metadata)
	return &copyValue
}

func cloneDelivery(value *DeliveryRecord) *DeliveryRecord {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.RedactedProviderMetadata = cloneStringMap(value.RedactedProviderMetadata)
	return &copyValue
}

func cloneEvent(value *MessageEvent) *MessageEvent {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.RedactedProviderMetadata = cloneStringMap(value.RedactedProviderMetadata)
	return &copyValue
}
