package conversation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"sapasora/platform/support/hash"
)

// InMemoryRepository is a concurrency-safe repository for domain tests and
// local composition. It has no provider, transport, or database behavior.
type InMemoryRepository struct {
	mu sync.RWMutex

	nextConversationID uint64
	nextParticipantID  uint64
	nextAssignmentID   uint64
	conversations      map[uint64]*Conversation
	conversationByKey  map[string]uint64
	participants       map[uint64]*Participant
	participantByKey   map[string]uint64
	assignments        map[uint64]*ConversationAssignment
	assignmentByKey    map[string]uint64
}

var _ Repository = (*InMemoryRepository)(nil)

// InMemoryConversationRepository is the explicit implementation name.
type InMemoryConversationRepository = InMemoryRepository

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		nextConversationID: 1,
		nextParticipantID:  1,
		nextAssignmentID:   1,
		conversations:      make(map[uint64]*Conversation),
		conversationByKey:  make(map[string]uint64),
		participants:       make(map[uint64]*Participant),
		participantByKey:   make(map[string]uint64),
		assignments:        make(map[uint64]*ConversationAssignment),
		assignmentByKey:    make(map[string]uint64),
	}
}

func NewInMemoryConversationRepository() *InMemoryRepository { return NewInMemoryRepository() }

func (r *InMemoryRepository) ensureInitializedLocked() {
	if r.nextConversationID == 0 {
		r.nextConversationID = 1
	}
	if r.nextParticipantID == 0 {
		r.nextParticipantID = 1
	}
	if r.nextAssignmentID == 0 {
		r.nextAssignmentID = 1
	}
	if r.conversations == nil {
		r.conversations = make(map[uint64]*Conversation)
	}
	if r.conversationByKey == nil {
		r.conversationByKey = make(map[string]uint64)
	}
	if r.participants == nil {
		r.participants = make(map[uint64]*Participant)
	}
	if r.participantByKey == nil {
		r.participantByKey = make(map[string]uint64)
	}
	if r.assignments == nil {
		r.assignments = make(map[uint64]*ConversationAssignment)
	}
	if r.assignmentByKey == nil {
		r.assignmentByKey = make(map[string]uint64)
	}
}

func repositoryContextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("conversation repository context is nil")
	}
	return ctx.Err()
}

func scopedKey(scope, publicID string) string {
	return fmt.Sprintf("%d:%s%d:%s", len(scope), scope, len(publicID), publicID)
}

func (r *InMemoryRepository) CreateConversation(ctx context.Context, value *Conversation) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: conversation is nil", ErrInvalidConversation)
	}
	candidate := cloneConversation(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversation, err)
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.Status == "" || candidate.Status == ConversationStatusUnknown {
		candidate.Status = ConversationStatusOpen
	}
	if candidate.Priority == "" || candidate.Priority == ConversationPriorityUnknown {
		candidate.Priority = ConversationPriorityNormal
	}
	var err error
	candidate.Tags, err = normalizeTags(candidate.Tags)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversation, err)
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	if candidate.LastActivityAt.IsZero() {
		candidate.LastActivityAt = candidate.CreatedAt
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
	if _, exists := r.conversationByKey[key]; exists {
		return ErrConversationConflict
	}
	candidate.ID = r.nextConversationID
	r.nextConversationID++
	r.conversations[candidate.ID] = candidate
	r.conversationByKey[key] = candidate.ID
	*value = *cloneConversation(candidate)
	return nil
}

func (r *InMemoryRepository) GetConversationByPublicID(ctx context.Context, tenantID, publicID string) (*Conversation, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(publicID) == "" {
		return nil, fmt.Errorf("%w: tenant and public ids are required", ErrInvalidConversation)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.conversationByKey[scopedKey(strings.TrimSpace(tenantID), publicID)]
	if !ok {
		return nil, ErrConversationNotFound
	}
	value, ok := r.conversations[id]
	if !ok || value.DeletedAt != nil {
		return nil, ErrConversationNotFound
	}
	return cloneConversation(value), nil
}

func (r *InMemoryRepository) GetConversationByID(ctx context.Context, tenantID string, id uint64) (*Conversation, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" || id == 0 {
		return nil, fmt.Errorf("%w: tenant and conversation id are required", ErrInvalidConversation)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.conversations[id]
	if !ok || value.ScopeID() != strings.TrimSpace(tenantID) || value.DeletedAt != nil {
		return nil, ErrConversationNotFound
	}
	return cloneConversation(value), nil
}

func (r *InMemoryRepository) ListConversations(ctx context.Context, tenantID string, filter ConversationFilter) ([]*Conversation, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, fmt.Errorf("%w: tenant id is required", ErrInvalidConversation)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Conversation, 0)
	for id := uint64(1); id < r.nextConversationID; id++ {
		value, ok := r.conversations[id]
		if !ok || value.ScopeID() != tenantID || value.DeletedAt != nil || !matchesConversation(value, filter) {
			continue
		}
		result = append(result, cloneConversation(value))
	}
	return result, nil
}

func matchesConversation(value *Conversation, filter ConversationFilter) bool {
	if filter.Status != nil && value.Status != *filter.Status {
		return false
	}
	if filter.Priority != nil && value.Priority != *filter.Priority {
		return false
	}
	if filter.ChannelPublicID != nil && value.ChannelPublicID != *filter.ChannelPublicID {
		return false
	}
	if filter.ContactPublicID != nil && value.ContactPublicID != *filter.ContactPublicID {
		return false
	}
	if filter.AssigneePublicID != nil && value.AssigneePublicID != *filter.AssigneePublicID {
		return false
	}
	if filter.TeamPublicID != nil && value.TeamPublicID != *filter.TeamPublicID {
		return false
	}
	if filter.Tag != nil {
		tag := strings.ToLower(strings.TrimSpace(*filter.Tag))
		found := false
		for _, candidate := range value.Tags {
			if candidate == tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if filter.UnreadOnly && value.UnreadCount == 0 {
		return false
	}
	if filter.SLAOverdueAt != nil && !conversationSLAOverdue(value, *filter.SLAOverdueAt) {
		return false
	}
	return true
}

func conversationSLAOverdue(value *Conversation, at time.Time) bool {
	if value.SLABreachedAt != nil {
		return true
	}
	if value.SLAFirstRespondedAt == nil && value.SLAFirstResponseDueAt != nil && value.SLAFirstResponseDueAt.Before(at) {
		return true
	}
	return value.SLAResolvedAt == nil && value.SLAResolutionDueAt != nil && value.SLAResolutionDueAt.Before(at)
}

func (r *InMemoryRepository) UpdateConversation(ctx context.Context, value *Conversation) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: conversation is nil", ErrInvalidConversation)
	}
	candidate := cloneConversation(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversation, err)
	}
	var err error
	candidate.Tags, err = normalizeTags(candidate.Tags)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversation, err)
	}
	if err := candidate.Valid(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.conversations[candidate.ID]
	if !ok || stored.ScopeID() != candidate.ScopeID() || stored.DeletedAt != nil {
		return ErrConversationNotFound
	}
	if stored.PublicID != candidate.PublicID {
		return ErrConversationConflict
	}
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	candidate.DeletedAt = nil
	r.conversations[candidate.ID] = candidate
	*value = *cloneConversation(candidate)
	return nil
}

func (r *InMemoryRepository) DeleteConversation(ctx context.Context, tenantID, publicID string) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.conversationByKey[scopedKey(strings.TrimSpace(tenantID), publicID)]
	if !ok {
		return ErrConversationNotFound
	}
	value, ok := r.conversations[id]
	if !ok || value.DeletedAt != nil {
		return ErrConversationNotFound
	}
	now := time.Now().UTC()
	value.DeletedAt = &now
	value.UpdatedAt = now
	return nil
}

func (r *InMemoryRepository) CreateParticipant(ctx context.Context, value *Participant) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: participant is nil", ErrInvalidParticipant)
	}
	candidate := cloneParticipant(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidParticipant, err)
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.JoinedAt.IsZero() {
		candidate.JoinedAt = time.Now().UTC()
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = candidate.JoinedAt
	}
	candidate.UpdatedAt = time.Now().UTC()
	candidate.DeletedAt = nil
	if err := candidate.Valid(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	conversationID, exists := r.activeConversationIDByKeyLocked(candidate.ScopeID(), candidate.ConversationPublicID)
	if !exists {
		return ErrConversationNotFound
	}
	if candidate.ConversationID != 0 && candidate.ConversationID != conversationID {
		return ErrParticipantConflict
	}
	candidate.ConversationID = conversationID
	key := scopedKey(candidate.ScopeID(), candidate.PublicID)
	if _, exists := r.participantByKey[key]; exists {
		return ErrParticipantConflict
	}
	candidate.ID = r.nextParticipantID
	r.nextParticipantID++
	r.participants[candidate.ID] = candidate
	r.participantByKey[key] = candidate.ID
	*value = *cloneParticipant(candidate)
	return nil
}

func (r *InMemoryRepository) GetParticipantByPublicID(ctx context.Context, tenantID, publicID string) (*Participant, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.participantByKey[scopedKey(strings.TrimSpace(tenantID), publicID)]
	if !ok {
		return nil, ErrParticipantNotFound
	}
	value, ok := r.participants[id]
	if !ok || value.DeletedAt != nil {
		return nil, ErrParticipantNotFound
	}
	return cloneParticipant(value), nil
}

func (r *InMemoryRepository) GetParticipantByID(ctx context.Context, tenantID string, id uint64) (*Participant, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.participants[id]
	if !ok || value.ScopeID() != strings.TrimSpace(tenantID) || value.DeletedAt != nil {
		return nil, ErrParticipantNotFound
	}
	return cloneParticipant(value), nil
}

func (r *InMemoryRepository) ListParticipants(ctx context.Context, tenantID, conversationPublicID string) ([]*Participant, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID = strings.TrimSpace(tenantID)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.activeConversationByKeyLocked(tenantID, conversationPublicID) {
		return nil, ErrConversationNotFound
	}
	result := make([]*Participant, 0)
	for id := uint64(1); id < r.nextParticipantID; id++ {
		value, ok := r.participants[id]
		if ok && value.ScopeID() == tenantID && value.ConversationPublicID == conversationPublicID && value.DeletedAt == nil {
			result = append(result, cloneParticipant(value))
		}
	}
	return result, nil
}

func (r *InMemoryRepository) UpdateParticipant(ctx context.Context, value *Participant) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: participant is nil", ErrInvalidParticipant)
	}
	candidate := cloneParticipant(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidParticipant, err)
	}
	if err := candidate.Valid(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.participants[candidate.ID]
	if !ok || stored.ScopeID() != candidate.ScopeID() || stored.DeletedAt != nil {
		return ErrParticipantNotFound
	}
	if stored.PublicID != candidate.PublicID || stored.ConversationPublicID != candidate.ConversationPublicID {
		return ErrParticipantConflict
	}
	conversationID, exists := r.activeConversationIDByKeyLocked(candidate.ScopeID(), candidate.ConversationPublicID)
	if !exists {
		return ErrConversationNotFound
	}
	if candidate.ConversationID != 0 && candidate.ConversationID != conversationID {
		return ErrParticipantConflict
	}
	candidate.ConversationID = conversationID
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	candidate.DeletedAt = nil
	r.participants[candidate.ID] = candidate
	*value = *cloneParticipant(candidate)
	return nil
}

func (r *InMemoryRepository) DeleteParticipant(ctx context.Context, tenantID, publicID string) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.participantByKey[scopedKey(strings.TrimSpace(tenantID), publicID)]
	if !ok {
		return ErrParticipantNotFound
	}
	value, ok := r.participants[id]
	if !ok || value.DeletedAt != nil {
		return ErrParticipantNotFound
	}
	now := time.Now().UTC()
	value.DeletedAt = &now
	value.UpdatedAt = now
	return nil
}

func (r *InMemoryRepository) CreateAssignment(ctx context.Context, value *ConversationAssignment) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: assignment is nil", ErrInvalidConversationAssignment)
	}
	candidate := cloneAssignment(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversationAssignment, err)
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.AssignedAt.IsZero() {
		candidate.AssignedAt = time.Now().UTC()
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = candidate.AssignedAt
	}
	candidate.UpdatedAt = time.Now().UTC()
	candidate.DeletedAt = nil
	if err := candidate.Valid(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	conversationID, exists := r.activeConversationIDByKeyLocked(candidate.ScopeID(), candidate.ConversationPublicID)
	if !exists {
		return ErrConversationNotFound
	}
	if candidate.ConversationID != 0 && candidate.ConversationID != conversationID {
		return ErrConversationAssignmentConflict
	}
	candidate.ConversationID = conversationID
	key := scopedKey(candidate.ScopeID(), candidate.PublicID)
	if _, exists := r.assignmentByKey[key]; exists {
		return ErrConversationAssignmentConflict
	}
	candidate.ID = r.nextAssignmentID
	r.nextAssignmentID++
	r.assignments[candidate.ID] = candidate
	r.assignmentByKey[key] = candidate.ID
	*value = *cloneAssignment(candidate)
	return nil
}

func (r *InMemoryRepository) GetAssignmentByPublicID(ctx context.Context, tenantID, publicID string) (*ConversationAssignment, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.assignmentByKey[scopedKey(strings.TrimSpace(tenantID), publicID)]
	if !ok {
		return nil, ErrConversationAssignmentNotFound
	}
	value, ok := r.assignments[id]
	if !ok || value.DeletedAt != nil {
		return nil, ErrConversationAssignmentNotFound
	}
	return cloneAssignment(value), nil
}

func (r *InMemoryRepository) GetAssignmentByID(ctx context.Context, tenantID string, id uint64) (*ConversationAssignment, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.assignments[id]
	if !ok || value.ScopeID() != strings.TrimSpace(tenantID) || value.DeletedAt != nil {
		return nil, ErrConversationAssignmentNotFound
	}
	return cloneAssignment(value), nil
}

func (r *InMemoryRepository) ListAssignments(ctx context.Context, tenantID, conversationPublicID string) ([]*ConversationAssignment, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	tenantID = strings.TrimSpace(tenantID)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.activeConversationByKeyLocked(tenantID, conversationPublicID) {
		return nil, ErrConversationNotFound
	}
	result := make([]*ConversationAssignment, 0)
	for id := uint64(1); id < r.nextAssignmentID; id++ {
		value, ok := r.assignments[id]
		if ok && value.ScopeID() == tenantID && value.ConversationPublicID == conversationPublicID && value.DeletedAt == nil {
			result = append(result, cloneAssignment(value))
		}
	}
	return result, nil
}

func (r *InMemoryRepository) UpdateAssignment(ctx context.Context, value *ConversationAssignment) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: assignment is nil", ErrInvalidConversationAssignment)
	}
	candidate := cloneAssignment(value)
	if err := canonicalizeScope(&candidate.TenantID, &candidate.WorkspaceID); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConversationAssignment, err)
	}
	if err := candidate.Valid(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.assignments[candidate.ID]
	if !ok || stored.ScopeID() != candidate.ScopeID() || stored.DeletedAt != nil {
		return ErrConversationAssignmentNotFound
	}
	if stored.PublicID != candidate.PublicID || stored.ConversationPublicID != candidate.ConversationPublicID {
		return ErrConversationAssignmentConflict
	}
	conversationID, exists := r.activeConversationIDByKeyLocked(candidate.ScopeID(), candidate.ConversationPublicID)
	if !exists {
		return ErrConversationNotFound
	}
	if candidate.ConversationID != 0 && candidate.ConversationID != conversationID {
		return ErrConversationAssignmentConflict
	}
	candidate.ConversationID = conversationID
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	candidate.DeletedAt = nil
	r.assignments[candidate.ID] = candidate
	*value = *cloneAssignment(candidate)
	return nil
}

func (r *InMemoryRepository) DeleteAssignment(ctx context.Context, tenantID, publicID string) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.assignmentByKey[scopedKey(strings.TrimSpace(tenantID), publicID)]
	if !ok {
		return ErrConversationAssignmentNotFound
	}
	value, ok := r.assignments[id]
	if !ok || value.DeletedAt != nil {
		return ErrConversationAssignmentNotFound
	}
	now := time.Now().UTC()
	value.DeletedAt = &now
	value.UpdatedAt = now
	return nil
}

func (r *InMemoryRepository) activeConversationByKeyLocked(scope, publicID string) bool {
	_, ok := r.activeConversationIDByKeyLocked(scope, publicID)
	return ok
}

func (r *InMemoryRepository) activeConversationIDByKeyLocked(scope, publicID string) (uint64, bool) {
	id, ok := r.conversationByKey[scopedKey(scope, publicID)]
	if !ok {
		return 0, false
	}
	value, ok := r.conversations[id]
	return id, ok && value.DeletedAt == nil
}

// Concise aliases mirror the other domain repositories for conversation CRUD.
func (r *InMemoryRepository) Create(ctx context.Context, value *Conversation) error {
	return r.CreateConversation(ctx, value)
}

func (r *InMemoryRepository) Get(ctx context.Context, tenantID, publicID string) (*Conversation, error) {
	return r.GetConversationByPublicID(ctx, tenantID, publicID)
}

func (r *InMemoryRepository) List(ctx context.Context, tenantID string, filter ConversationFilter) ([]*Conversation, error) {
	return r.ListConversations(ctx, tenantID, filter)
}

func (r *InMemoryRepository) Update(ctx context.Context, value *Conversation) error {
	return r.UpdateConversation(ctx, value)
}

func (r *InMemoryRepository) Delete(ctx context.Context, tenantID, publicID string) error {
	return r.DeleteConversation(ctx, tenantID, publicID)
}
