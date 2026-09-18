package message

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"sapasora/platform/support/hash"
)

// InMemoryRepository is a concurrency-safe domain/test adapter. It does not
// emulate provider behavior or durable database transactions.
type InMemoryRepository struct {
	mu sync.RWMutex

	nextMessageID    uint64
	nextAttachmentID uint64
	nextDeliveryID   uint64
	nextEventID      uint64
	messages         map[uint64]*Message
	messageByKey     map[string]uint64
	idempotencyByKey map[string]uint64
	attachments      map[uint64]*MessageAttachment
	attachmentByKey  map[string]uint64
	deliveries       map[uint64]*DeliveryRecord
	deliveryByKey    map[string]uint64
	events           map[string][]*MessageEvent
}

var _ Repository = (*InMemoryRepository)(nil)

// InMemoryMessageRepository is the descriptive implementation name.
type InMemoryMessageRepository = InMemoryRepository

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		nextMessageID:    1,
		nextAttachmentID: 1,
		nextDeliveryID:   1,
		nextEventID:      1,
		messages:         make(map[uint64]*Message),
		messageByKey:     make(map[string]uint64),
		idempotencyByKey: make(map[string]uint64),
		attachments:      make(map[uint64]*MessageAttachment),
		attachmentByKey:  make(map[string]uint64),
		deliveries:       make(map[uint64]*DeliveryRecord),
		deliveryByKey:    make(map[string]uint64),
		events:           make(map[string][]*MessageEvent),
	}
}

func NewInMemoryMessageRepository() *InMemoryRepository { return NewInMemoryRepository() }

func (r *InMemoryRepository) ensureInitializedLocked() {
	if r.nextMessageID == 0 {
		r.nextMessageID = 1
	}
	if r.nextAttachmentID == 0 {
		r.nextAttachmentID = 1
	}
	if r.nextDeliveryID == 0 {
		r.nextDeliveryID = 1
	}
	if r.nextEventID == 0 {
		r.nextEventID = 1
	}
	if r.messages == nil {
		r.messages = make(map[uint64]*Message)
	}
	if r.messageByKey == nil {
		r.messageByKey = make(map[string]uint64)
	}
	if r.idempotencyByKey == nil {
		r.idempotencyByKey = make(map[string]uint64)
	}
	if r.attachments == nil {
		r.attachments = make(map[uint64]*MessageAttachment)
	}
	if r.attachmentByKey == nil {
		r.attachmentByKey = make(map[string]uint64)
	}
	if r.deliveries == nil {
		r.deliveries = make(map[uint64]*DeliveryRecord)
	}
	if r.deliveryByKey == nil {
		r.deliveryByKey = make(map[string]uint64)
	}
	if r.events == nil {
		r.events = make(map[string][]*MessageEvent)
	}
}

func repositoryContextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("message repository context is nil")
	}
	return ctx.Err()
}

func scopedKey(scope, publicID string) string {
	return fmt.Sprintf("%d:%s%d:%s", len(scope), scope, len(publicID), publicID)
}

func idempotencyKey(scope, key string) string { return scopedKey(scope, key) }

func (r *InMemoryRepository) CreateMessage(ctx context.Context, value *Message) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: message is nil", ErrInvalidMessage)
	}
	candidate := cloneMessage(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.Status == MessageStatusUnknown || candidate.Status == "" {
		candidate.Status = MessageStatusAccepted
	}
	candidate.Provider = strings.TrimSpace(candidate.Provider)
	candidate.NormalizedContent = normalizeContent(candidate.NormalizedContent)
	if candidate.NormalizedContent == "" {
		candidate.NormalizedContent = normalizeContent(candidate.Content)
	}
	candidate.Content = candidate.NormalizedContent
	candidate.ContentMetadata = contentMetadata(candidate.NormalizedContent, effectiveContentMetadata(*candidate))
	candidate.NormalizedContentMetadata = candidate.ContentMetadata
	candidate.RedactedProviderMetadata = redactMetadata(firstMetadata(candidate.RedactedProviderMetadata, candidate.ProviderMetadata))
	candidate.ProviderMetadata = cloneStringMap(candidate.RedactedProviderMetadata)
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.UpdatedAt = time.Now().UTC()
	candidate.DeletedAt = nil
	if err := candidate.Valid(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	key := scopedKey(candidate.ScopeID(), candidate.PublicID)
	if _, exists := r.messageByKey[key]; exists {
		return ErrMessageConflict
	}
	if candidate.IdempotencyKey != "" {
		if _, exists := r.idempotencyByKey[idempotencyKey(candidate.ScopeID(), candidate.IdempotencyKey)]; exists {
			return ErrMessageConflict
		}
	}
	candidate.ID = r.nextMessageID
	r.nextMessageID++
	r.messages[candidate.ID] = candidate
	r.messageByKey[key] = candidate.ID
	if candidate.IdempotencyKey != "" {
		r.idempotencyByKey[idempotencyKey(candidate.ScopeID(), candidate.IdempotencyKey)] = candidate.ID
	}
	*value = *cloneMessage(candidate)
	return nil
}

func (r *InMemoryRepository) GetMessageByPublicID(ctx context.Context, tenantID, publicID string) (*Message, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID, publicID = strings.TrimSpace(tenantID), strings.TrimSpace(publicID)
	if tenantID == "" || publicID == "" {
		return nil, fmt.Errorf("%w: tenant and public ids are required", ErrInvalidMessage)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.messageByKey[scopedKey(tenantID, publicID)]
	if !ok {
		return nil, ErrMessageNotFound
	}
	value, ok := r.messages[id]
	if !ok || value.DeletedAt != nil {
		return nil, ErrMessageNotFound
	}
	return cloneMessage(value), nil
}

func (r *InMemoryRepository) GetMessageByID(ctx context.Context, tenantID string, id uint64) (*Message, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" || id == 0 {
		return nil, fmt.Errorf("%w: tenant and message id are required", ErrInvalidMessage)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.messages[id]
	if !ok || value.ScopeID() != tenantID || value.DeletedAt != nil {
		return nil, ErrMessageNotFound
	}
	return cloneMessage(value), nil
}

func (r *InMemoryRepository) ListMessages(ctx context.Context, tenantID string, filter MessageFilter) ([]*Message, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, fmt.Errorf("%w: tenant id is required", ErrInvalidMessage)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Message, 0)
	for id := uint64(1); id < r.nextMessageID; id++ {
		value, ok := r.messages[id]
		if !ok || value.ScopeID() != tenantID || value.DeletedAt != nil || !matchesMessage(value, filter) {
			continue
		}
		result = append(result, cloneMessage(value))
	}
	return result, nil
}

func matchesMessage(value *Message, filter MessageFilter) bool {
	if filter.Direction != nil && value.Direction != *filter.Direction {
		return false
	}
	if filter.Status != nil && value.Status != *filter.Status {
		return false
	}
	if filter.ChannelPublicID != nil && value.ChannelPublicID != *filter.ChannelPublicID {
		return false
	}
	if filter.ContactPublicID != nil && value.ContactPublicID != *filter.ContactPublicID {
		return false
	}
	if filter.ConversationPublicID != nil && value.ConversationPublicID != *filter.ConversationPublicID {
		return false
	}
	if filter.IdempotencyKey != nil && value.IdempotencyKey != *filter.IdempotencyKey {
		return false
	}
	if filter.ProviderMessageID != nil && value.ProviderMessageID != *filter.ProviderMessageID {
		return false
	}
	return true
}

func (r *InMemoryRepository) UpdateMessage(ctx context.Context, value *Message) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: message is nil", ErrInvalidMessage)
	}
	candidate := cloneMessage(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	candidate.NormalizedContent = normalizeContent(candidate.NormalizedContent)
	if candidate.NormalizedContent == "" {
		candidate.NormalizedContent = normalizeContent(candidate.Content)
	}
	candidate.Content = candidate.NormalizedContent
	candidate.ContentMetadata = contentMetadata(candidate.NormalizedContent, effectiveContentMetadata(*candidate))
	candidate.NormalizedContentMetadata = candidate.ContentMetadata
	candidate.RedactedProviderMetadata = redactMetadata(firstMetadata(candidate.RedactedProviderMetadata, candidate.ProviderMetadata))
	candidate.ProviderMetadata = cloneStringMap(candidate.RedactedProviderMetadata)
	if err := candidate.Valid(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.messages[candidate.ID]
	if !ok || stored.ScopeID() != candidate.ScopeID() || stored.DeletedAt != nil {
		return ErrMessageNotFound
	}
	if stored.PublicID != candidate.PublicID {
		return ErrMessageConflict
	}
	if stored.IdempotencyKey != candidate.IdempotencyKey {
		if candidate.IdempotencyKey != "" {
			if existing, exists := r.idempotencyByKey[idempotencyKey(candidate.ScopeID(), candidate.IdempotencyKey)]; exists && existing != candidate.ID {
				return ErrMessageConflict
			}
		}
		if stored.IdempotencyKey != "" {
			delete(r.idempotencyByKey, idempotencyKey(candidate.ScopeID(), stored.IdempotencyKey))
		}
		if candidate.IdempotencyKey != "" {
			r.idempotencyByKey[idempotencyKey(candidate.ScopeID(), candidate.IdempotencyKey)] = candidate.ID
		}
	}
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	candidate.DeletedAt = nil
	r.messages[candidate.ID] = candidate
	*value = *cloneMessage(candidate)
	return nil
}

func (r *InMemoryRepository) CreateAttachment(ctx context.Context, value *MessageAttachment) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: attachment is nil", ErrInvalidAttachment)
	}
	candidate := cloneAttachment(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidAttachment, err)
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.Metadata = redactMetadata(candidate.Metadata)
	if err := candidate.Valid(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	messageID, ok := r.activeMessageIDLocked(candidate.ScopeID(), candidate.MessagePublicID)
	if !ok {
		return ErrMessageNotFound
	}
	if candidate.MessageID != 0 && candidate.MessageID != messageID {
		return ErrAttachmentConflict
	}
	candidate.MessageID = messageID
	key := scopedKey(candidate.ScopeID(), candidate.PublicID)
	if _, exists := r.attachmentByKey[key]; exists {
		return ErrAttachmentConflict
	}
	candidate.ID = r.nextAttachmentID
	r.nextAttachmentID++
	r.attachments[candidate.ID] = candidate
	r.attachmentByKey[key] = candidate.ID
	message := r.messages[messageID]
	message.ContentMetadata = effectiveContentMetadata(*message)
	message.ContentMetadata.AttachmentCount++
	message.NormalizedContentMetadata = message.ContentMetadata
	message.UpdatedAt = time.Now().UTC()
	*value = *cloneAttachment(candidate)
	return nil
}

func (r *InMemoryRepository) GetAttachmentByPublicID(ctx context.Context, tenantID, publicID string) (*MessageAttachment, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.attachmentByKey[scopedKey(strings.TrimSpace(tenantID), strings.TrimSpace(publicID))]
	if !ok {
		return nil, ErrAttachmentNotFound
	}
	value, ok := r.attachments[id]
	if !ok {
		return nil, ErrAttachmentNotFound
	}
	return cloneAttachment(value), nil
}

func (r *InMemoryRepository) GetAttachmentByID(ctx context.Context, tenantID string, id uint64) (*MessageAttachment, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" || id == 0 {
		return nil, fmt.Errorf("%w: tenant and attachment id are required", ErrInvalidAttachment)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.attachments[id]
	if !ok || value.ScopeID() != tenantID {
		return nil, ErrAttachmentNotFound
	}
	return cloneAttachment(value), nil
}

func (r *InMemoryRepository) ListAttachments(ctx context.Context, tenantID, messagePublicID string) ([]*MessageAttachment, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID, messagePublicID = strings.TrimSpace(tenantID), strings.TrimSpace(messagePublicID)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.activeMessageIDLocked(tenantID, messagePublicID); !ok {
		return nil, ErrMessageNotFound
	}
	result := make([]*MessageAttachment, 0)
	for id := uint64(1); id < r.nextAttachmentID; id++ {
		value, ok := r.attachments[id]
		if ok && value.ScopeID() == tenantID && value.MessagePublicID == messagePublicID {
			result = append(result, cloneAttachment(value))
		}
	}
	return result, nil
}

func (r *InMemoryRepository) CreateDeliveryRecord(ctx context.Context, value *DeliveryRecord) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: delivery record is nil", ErrInvalidDeliveryRecord)
	}
	candidate := cloneDelivery(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDeliveryRecord, err)
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.Status == MessageStatusUnknown || candidate.Status == "" {
		candidate.Status = MessageStatusAccepted
	}
	if candidate.OccurredAt.IsZero() {
		candidate.OccurredAt = time.Now().UTC()
	}
	if candidate.RecordedAt.IsZero() {
		candidate.RecordedAt = time.Now().UTC()
	}
	candidate.Provider = strings.TrimSpace(candidate.Provider)
	candidate.RedactedProviderMetadata = redactMetadata(candidate.RedactedProviderMetadata)
	if err := candidate.Valid(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	messageID, ok := r.activeMessageIDLocked(candidate.ScopeID(), candidate.MessagePublicID)
	if !ok {
		return ErrMessageNotFound
	}
	if candidate.MessageID != 0 && candidate.MessageID != messageID {
		return ErrDeliveryRecordConflict
	}
	candidate.MessageID = messageID
	key := scopedKey(candidate.ScopeID(), candidate.PublicID)
	if _, exists := r.deliveryByKey[key]; exists {
		return ErrDeliveryRecordConflict
	}
	candidate.ID = r.nextDeliveryID
	r.nextDeliveryID++
	r.deliveries[candidate.ID] = candidate
	r.deliveryByKey[key] = candidate.ID
	*value = *cloneDelivery(candidate)
	return nil
}

func (r *InMemoryRepository) GetDeliveryRecordByPublicID(ctx context.Context, tenantID, publicID string) (*DeliveryRecord, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.deliveryByKey[scopedKey(strings.TrimSpace(tenantID), strings.TrimSpace(publicID))]
	if !ok {
		return nil, ErrDeliveryRecordNotFound
	}
	value, ok := r.deliveries[id]
	if !ok {
		return nil, ErrDeliveryRecordNotFound
	}
	return cloneDelivery(value), nil
}

func (r *InMemoryRepository) GetDeliveryRecordByID(ctx context.Context, tenantID string, id uint64) (*DeliveryRecord, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" || id == 0 {
		return nil, fmt.Errorf("%w: tenant and delivery id are required", ErrInvalidDeliveryRecord)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.deliveries[id]
	if !ok || value.ScopeID() != tenantID {
		return nil, ErrDeliveryRecordNotFound
	}
	return cloneDelivery(value), nil
}

func (r *InMemoryRepository) ListDeliveryRecords(ctx context.Context, tenantID, messagePublicID string) ([]*DeliveryRecord, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID, messagePublicID = strings.TrimSpace(tenantID), strings.TrimSpace(messagePublicID)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.activeMessageIDLocked(tenantID, messagePublicID); !ok {
		return nil, ErrMessageNotFound
	}
	result := make([]*DeliveryRecord, 0)
	for id := uint64(1); id < r.nextDeliveryID; id++ {
		value, ok := r.deliveries[id]
		if ok && value.ScopeID() == tenantID && value.MessagePublicID == messagePublicID {
			result = append(result, cloneDelivery(value))
		}
	}
	return result, nil
}

func (r *InMemoryRepository) AppendMessageEvent(ctx context.Context, value *MessageEvent) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: event is nil", ErrInvalidMessageEvent)
	}
	candidate := cloneEvent(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMessageEvent, err)
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.OccurredAt.IsZero() {
		candidate.OccurredAt = time.Now().UTC()
	}
	if candidate.RecordedAt.IsZero() {
		candidate.RecordedAt = time.Now().UTC()
	}
	candidate.RedactedProviderMetadata = redactMetadata(candidate.RedactedProviderMetadata)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	messageID, ok := r.activeMessageIDLocked(candidate.ScopeID(), candidate.MessagePublicID)
	if !ok {
		return ErrMessageNotFound
	}
	if candidate.MessageID != 0 && candidate.MessageID != messageID {
		return ErrMessageEventConflict
	}
	candidate.MessageID = messageID
	key := scopedKey(candidate.ScopeID(), candidate.MessagePublicID)
	candidate.Sequence = uint64(len(r.events[key]) + 1)
	if err := candidate.Valid(); err != nil {
		return err
	}
	candidate.ID = r.nextEventID
	r.nextEventID++
	r.events[key] = append(r.events[key], candidate)
	*value = *cloneEvent(candidate)
	return nil
}

func (r *InMemoryRepository) ListMessageEvents(ctx context.Context, tenantID, messagePublicID string) ([]*MessageEvent, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID, messagePublicID = strings.TrimSpace(tenantID), strings.TrimSpace(messagePublicID)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.activeMessageIDLocked(tenantID, messagePublicID); !ok {
		return nil, ErrMessageNotFound
	}
	stored := r.events[scopedKey(tenantID, messagePublicID)]
	result := make([]*MessageEvent, 0, len(stored))
	for _, event := range stored {
		result = append(result, cloneEvent(event))
	}
	return result, nil
}

func (r *InMemoryRepository) activeMessageIDLocked(scope, publicID string) (uint64, bool) {
	id, ok := r.messageByKey[scopedKey(scope, publicID)]
	if !ok {
		return 0, false
	}
	value, ok := r.messages[id]
	return id, ok && value.DeletedAt == nil
}

func canonicalizeScope(tenantID, workspaceID *string) error {
	scope, err := validateScope(*tenantID, *workspaceID)
	if err != nil {
		return err
	}
	*tenantID, *workspaceID = scope, scope
	return nil
}

func firstMetadata(primary, fallback map[string]string) map[string]string {
	if primary != nil {
		return primary
	}
	return fallback
}

func redactMetadata(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	result := make(map[string]string, len(values))
	for key, value := range values {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		lower := strings.ToLower(key)
		if containsSensitiveMetadataName(lower) {
			value = "[REDACTED]"
		}
		result[key] = strings.TrimSpace(value)
	}
	return result
}

func containsSensitiveMetadataName(key string) bool {
	for _, part := range []string{"token", "secret", "password", "credential", "authorization", "cookie", "payload", "body", "content"} {
		if strings.Contains(key, part) {
			return true
		}
	}
	return false
}

// Concise aliases mirror other domain repositories.
func (r *InMemoryRepository) Create(ctx context.Context, value *Message) error {
	return r.CreateMessage(ctx, value)
}
func (r *InMemoryRepository) Get(ctx context.Context, tenantID, publicID string) (*Message, error) {
	return r.GetMessageByPublicID(ctx, tenantID, publicID)
}
func (r *InMemoryRepository) List(ctx context.Context, tenantID string, filter MessageFilter) ([]*Message, error) {
	return r.ListMessages(ctx, tenantID, filter)
}
func (r *InMemoryRepository) Update(ctx context.Context, value *Message) error {
	return r.UpdateMessage(ctx, value)
}
func (r *InMemoryRepository) CreateMessageAttachment(ctx context.Context, value *MessageAttachment) error {
	return r.CreateAttachment(ctx, value)
}
func (r *InMemoryRepository) CreateDelivery(ctx context.Context, value *DeliveryRecord) error {
	return r.CreateDeliveryRecord(ctx, value)
}
func (r *InMemoryRepository) CreateMessageDelivery(ctx context.Context, value *MessageDelivery) error {
	return r.CreateDeliveryRecord(ctx, value)
}
func (r *InMemoryRepository) CreateMessageEvent(ctx context.Context, value *MessageEvent) error {
	return r.AppendMessageEvent(ctx, value)
}
func (r *InMemoryRepository) RecordMessageEvent(ctx context.Context, value *MessageEvent) error {
	return r.AppendMessageEvent(ctx, value)
}
func (r *InMemoryRepository) AppendEvent(ctx context.Context, value *MessageEvent) error {
	return r.AppendMessageEvent(ctx, value)
}
func (r *InMemoryRepository) ListEvents(ctx context.Context, tenantID, messagePublicID string) ([]*MessageEvent, error) {
	return r.ListMessageEvents(ctx, tenantID, messagePublicID)
}
func (r *InMemoryRepository) ListMessageDeliveries(ctx context.Context, tenantID, messagePublicID string) ([]*MessageDelivery, error) {
	return r.ListDeliveryRecords(ctx, tenantID, messagePublicID)
}
