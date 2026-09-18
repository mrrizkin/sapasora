package channel

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"sapasora/platform/support/hash"
)

// InMemoryChannelAccountRepository is a concurrency-safe repository useful for
// domain tests and local composition. It deliberately has no database schema
// or provider behavior.
type InMemoryChannelAccountRepository struct {
	mu sync.RWMutex

	nextAccountID uint64
	nextEventID   uint64
	accounts      map[uint64]*ChannelAccount
	byPublicID    map[string]uint64
	publicIDs     map[string]uint64
	events        map[string][]*ConnectionEvent
}

var _ ChannelAccountRepository = (*InMemoryChannelAccountRepository)(nil)

// InMemoryRepository is a concise name for the in-memory domain repository.
type InMemoryRepository = InMemoryChannelAccountRepository

// NewInMemoryChannelAccountRepository creates an empty repository.
func NewInMemoryChannelAccountRepository() *InMemoryChannelAccountRepository {
	return &InMemoryChannelAccountRepository{
		nextAccountID: 1,
		nextEventID:   1,
		accounts:      make(map[uint64]*ChannelAccount),
		byPublicID:    make(map[string]uint64),
		publicIDs:     make(map[string]uint64),
		events:        make(map[string][]*ConnectionEvent),
	}
}

// NewInMemoryRepository creates an empty repository using the concise name.
func NewInMemoryRepository() *InMemoryChannelAccountRepository {
	return NewInMemoryChannelAccountRepository()
}

func (r *InMemoryChannelAccountRepository) ensureInitializedLocked() {
	if r.nextAccountID == 0 {
		r.nextAccountID = 1
	}
	if r.nextEventID == 0 {
		r.nextEventID = 1
	}
	if r.accounts == nil {
		r.accounts = make(map[uint64]*ChannelAccount)
	}
	if r.byPublicID == nil {
		r.byPublicID = make(map[string]uint64)
	}
	if r.publicIDs == nil {
		r.publicIDs = make(map[string]uint64)
	}
	if r.events == nil {
		r.events = make(map[string][]*ConnectionEvent)
	}
}

func contextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("channel repository context is nil")
	}
	return ctx.Err()
}

func tenantKey(tenantID, publicID string) string {
	return tenantID + "\x00" + publicID
}

func (r *InMemoryChannelAccountRepository) CreateChannelAccount(ctx context.Context, account *ChannelAccount) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("%w: account is nil", ErrInvalidChannelAccount)
	}

	candidate := cloneChannelAccount(account)
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.Status == StateUnknown {
		candidate.Status = StatePending
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.UpdatedAt = time.Now().UTC()
	if err := candidate.Valid(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	key := tenantKey(candidate.TenantID, candidate.PublicID)
	if _, exists := r.publicIDs[candidate.PublicID]; exists {
		return fmt.Errorf("%w: public id %q", ErrChannelAccountConflict, candidate.PublicID)
	}
	candidate.ID = r.nextAccountID
	r.nextAccountID++
	r.accounts[candidate.ID] = candidate
	r.byPublicID[key] = candidate.ID
	r.publicIDs[candidate.PublicID] = candidate.ID
	*account = *cloneChannelAccount(candidate)
	return nil
}

func (r *InMemoryChannelAccountRepository) GetChannelAccountByPublicID(ctx context.Context, tenantID, publicID string) (*ChannelAccount, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byPublicID[tenantKey(tenantID, publicID)]
	if !ok {
		return nil, ErrChannelAccountNotFound
	}
	account, ok := r.accounts[id]
	if !ok || account.DeletedAt != nil {
		return nil, ErrChannelAccountNotFound
	}
	return cloneChannelAccount(account), nil
}

// GetChannelAccountByID is retained for internal persistence adapters. Callers
// still need the tenant ID so an internal ID cannot bypass tenant isolation.
func (r *InMemoryChannelAccountRepository) GetChannelAccountByID(ctx context.Context, tenantID string, id uint64) (*ChannelAccount, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	account, ok := r.accounts[id]
	if !ok || account.TenantID != tenantID || account.DeletedAt != nil {
		return nil, ErrChannelAccountNotFound
	}
	return cloneChannelAccount(account), nil
}

func (r *InMemoryChannelAccountRepository) ListChannelAccounts(ctx context.Context, tenantID string, filter ChannelAccountFilter) ([]*ChannelAccount, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant id is required", ErrInvalidChannelAccount)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*ChannelAccount, 0)
	for id := uint64(1); id < r.nextAccountID; id++ {
		account, ok := r.accounts[id]
		if !ok || account.TenantID != tenantID || account.DeletedAt != nil {
			continue
		}
		if filter.Type != nil && account.Type != *filter.Type {
			continue
		}
		if filter.Provider != nil && account.Provider != *filter.Provider {
			continue
		}
		if filter.Status != nil && account.Status != *filter.Status {
			continue
		}
		result = append(result, cloneChannelAccount(account))
	}
	return result, nil
}

func (r *InMemoryChannelAccountRepository) UpdateChannelAccount(ctx context.Context, account *ChannelAccount) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("%w: account is nil", ErrInvalidChannelAccount)
	}
	candidate := cloneChannelAccount(account)
	if err := candidate.Valid(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.accounts[candidate.ID]
	if !ok || stored.TenantID != candidate.TenantID || stored.DeletedAt != nil {
		return ErrChannelAccountNotFound
	}
	if stored.PublicID != candidate.PublicID {
		return fmt.Errorf("%w: public id is immutable", ErrChannelAccountConflict)
	}
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	r.accounts[candidate.ID] = candidate
	*account = *cloneChannelAccount(candidate)
	return nil
}

func (r *InMemoryChannelAccountRepository) DeleteChannelAccount(ctx context.Context, tenantID, publicID string) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.byPublicID[tenantKey(tenantID, publicID)]
	if !ok {
		return ErrChannelAccountNotFound
	}
	account, ok := r.accounts[id]
	if !ok || account.DeletedAt != nil {
		return ErrChannelAccountNotFound
	}
	now := time.Now().UTC()
	account.DeletedAt = &now
	account.UpdatedAt = now
	return nil
}

func (r *InMemoryChannelAccountRepository) RecordConnectionEvent(ctx context.Context, event *ConnectionEvent) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("%w: event is nil", ErrInvalidConnectionEvent)
	}
	candidate := cloneConnectionEvent(event)
	if candidate.OccurredAt.IsZero() {
		candidate.OccurredAt = time.Now().UTC()
	}
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if err := candidate.Valid(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	id, ok := r.byPublicID[tenantKey(candidate.TenantID, candidate.ChannelAccountPublicID)]
	account, exists := r.accounts[id]
	if !ok || !exists || account.DeletedAt != nil {
		return ErrChannelAccountNotFound
	}
	candidate.ID = r.nextEventID
	r.nextEventID++
	key := tenantKey(candidate.TenantID, candidate.ChannelAccountPublicID)
	r.events[key] = append(r.events[key], candidate)
	account.Status = candidate.ToState
	account.LastEventAt = cloneTime(&candidate.OccurredAt)
	if candidate.ToState == StateConnected {
		account.LastConnectedAt = cloneTime(&candidate.OccurredAt)
	}
	if candidate.ErrorCode != "" {
		account.LastError = candidate.ErrorCode
	}
	account.UpdatedAt = time.Now().UTC()
	*event = *cloneConnectionEvent(candidate)
	return nil
}

func (r *InMemoryChannelAccountRepository) ListConnectionEvents(ctx context.Context, tenantID, channelPublicID string) ([]*ConnectionEvent, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byPublicID[tenantKey(tenantID, channelPublicID)]
	if !ok {
		return nil, ErrChannelAccountNotFound
	}
	if account, exists := r.accounts[id]; !exists || account.DeletedAt != nil {
		return nil, ErrChannelAccountNotFound
	}
	stored := r.events[tenantKey(tenantID, channelPublicID)]
	result := make([]*ConnectionEvent, 0, len(stored))
	for _, event := range stored {
		result = append(result, cloneConnectionEvent(event))
	}
	return result, nil
}

// Short method aliases make the in-memory implementation convenient in small
// domain tests without changing the tenant-scoped repository interface.
func (r *InMemoryChannelAccountRepository) Create(ctx context.Context, account *ChannelAccount) error {
	return r.CreateChannelAccount(ctx, account)
}

func (r *InMemoryChannelAccountRepository) Get(ctx context.Context, tenantID, publicID string) (*ChannelAccount, error) {
	return r.GetChannelAccountByPublicID(ctx, tenantID, publicID)
}

func (r *InMemoryChannelAccountRepository) List(ctx context.Context, tenantID string, filter ChannelAccountFilter) ([]*ChannelAccount, error) {
	return r.ListChannelAccounts(ctx, tenantID, filter)
}

func (r *InMemoryChannelAccountRepository) Update(ctx context.Context, account *ChannelAccount) error {
	return r.UpdateChannelAccount(ctx, account)
}

func (r *InMemoryChannelAccountRepository) Delete(ctx context.Context, tenantID, publicID string) error {
	return r.DeleteChannelAccount(ctx, tenantID, publicID)
}
