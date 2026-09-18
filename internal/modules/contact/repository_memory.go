package contact

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"sapasora/platform/support/hash"
)

// InMemoryContactRepository is a concurrency-safe repository for domain tests
// and local composition. It enforces the same tenant and normalized-identity
// invariants a durable adapter must enforce with database constraints.
type InMemoryContactRepository struct {
	mu sync.RWMutex

	nextContactID uint64
	nextAddressID uint64
	contacts      map[uint64]*Contact
	contactByKey  map[string]uint64
	addresses     map[uint64]*ContactAddress
	addressByKey  map[string]uint64
	identityIndex map[string]uint64
}

var _ ContactRepository = (*InMemoryContactRepository)(nil)

// InMemoryRepository is the concise in-memory implementation name.
type InMemoryRepository = InMemoryContactRepository

func NewInMemoryContactRepository() *InMemoryContactRepository {
	return &InMemoryContactRepository{
		nextContactID: 1,
		nextAddressID: 1,
		contacts:      make(map[uint64]*Contact),
		contactByKey:  make(map[string]uint64),
		addresses:     make(map[uint64]*ContactAddress),
		addressByKey:  make(map[string]uint64),
		identityIndex: make(map[string]uint64),
	}
}

func NewInMemoryRepository() *InMemoryContactRepository { return NewInMemoryContactRepository() }

func (r *InMemoryContactRepository) ensureInitializedLocked() {
	if r.nextContactID == 0 {
		r.nextContactID = 1
	}
	if r.nextAddressID == 0 {
		r.nextAddressID = 1
	}
	if r.contacts == nil {
		r.contacts = make(map[uint64]*Contact)
	}
	if r.contactByKey == nil {
		r.contactByKey = make(map[string]uint64)
	}
	if r.addresses == nil {
		r.addresses = make(map[uint64]*ContactAddress)
	}
	if r.addressByKey == nil {
		r.addressByKey = make(map[string]uint64)
	}
	if r.identityIndex == nil {
		r.identityIndex = make(map[string]uint64)
	}
}

func repositoryContextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("contact repository context is nil")
	}
	return ctx.Err()
}

func scopedKey(tenantID, publicID string) string {
	return fmt.Sprintf("%d:%s%d:%s", len(tenantID), tenantID, len(publicID), publicID)
}

func tenantIdentityKey(tenantID string, identity AddressIdentity) string {
	return fmt.Sprintf("%d:%s%s", len(tenantID), tenantID, identity.Key())
}

func (r *InMemoryContactRepository) CreateContact(ctx context.Context, contact *Contact) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if contact == nil {
		return fmt.Errorf("%w: contact is nil", ErrInvalidContact)
	}
	candidate := cloneContact(contact)
	if strings.TrimSpace(candidate.PublicID) == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.Status == "" || candidate.Status == ContactStatusUnknown {
		candidate.Status = ContactStatusActive
	}
	if candidate.Source == "" || candidate.Source == ContactSourceUnknown {
		candidate.Source = ContactSourceManual
	}
	if candidate.DeletedAt != nil {
		return fmt.Errorf("%w: deleted contact cannot be created", ErrInvalidContact)
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.UpdatedAt = time.Now().UTC()
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidContact, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	key := scopedKey(candidate.TenantID, candidate.PublicID)
	if _, exists := r.contactByKey[key]; exists {
		return ErrContactConflict
	}
	candidate.ID = r.nextContactID
	r.nextContactID++
	r.contacts[candidate.ID] = candidate
	r.contactByKey[key] = candidate.ID
	*contact = *cloneContact(candidate)
	return nil
}

func (r *InMemoryContactRepository) GetContactByPublicID(ctx context.Context, tenantID, publicID string) (*Contact, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.contactByKey[scopedKey(tenantID, publicID)]
	if !ok {
		return nil, ErrContactNotFound
	}
	contact, ok := r.contacts[id]
	if !ok || contact.DeletedAt != nil {
		return nil, ErrContactNotFound
	}
	return cloneContact(contact), nil
}

func (r *InMemoryContactRepository) GetContactByID(ctx context.Context, tenantID string, id uint64) (*Contact, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	contact, ok := r.contacts[id]
	if !ok || contact.TenantID != tenantID || contact.DeletedAt != nil {
		return nil, ErrContactNotFound
	}
	return cloneContact(contact), nil
}

func (r *InMemoryContactRepository) ListContacts(ctx context.Context, tenantID string, filter ContactFilter) ([]*Contact, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant id is required", ErrInvalidContact)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Contact, 0)
	for id := uint64(1); id < r.nextContactID; id++ {
		contact, ok := r.contacts[id]
		if !ok || contact.TenantID != tenantID || contact.DeletedAt != nil {
			continue
		}
		if filter.Status != nil && contact.Status != *filter.Status {
			continue
		}
		if filter.Source != nil && contact.Source != *filter.Source {
			continue
		}
		result = append(result, cloneContact(contact))
	}
	return result, nil
}

func (r *InMemoryContactRepository) UpdateContact(ctx context.Context, contact *Contact) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if contact == nil {
		return fmt.Errorf("%w: contact is nil", ErrInvalidContact)
	}
	candidate := cloneContact(contact)
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidContact, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.contacts[candidate.ID]
	if !ok || stored.TenantID != candidate.TenantID || stored.DeletedAt != nil {
		return ErrContactNotFound
	}
	if stored.PublicID != candidate.PublicID {
		return ErrContactConflict
	}
	candidate.DeletedAt = nil
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	r.contacts[candidate.ID] = candidate
	*contact = *cloneContact(candidate)
	return nil
}

func (r *InMemoryContactRepository) DeleteContact(ctx context.Context, tenantID, publicID string) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.contactByKey[scopedKey(tenantID, publicID)]
	if !ok {
		return ErrContactNotFound
	}
	stored, ok := r.contacts[id]
	if !ok || stored.DeletedAt != nil {
		return ErrContactNotFound
	}
	now := time.Now().UTC()
	stored.DeletedAt = &now
	stored.UpdatedAt = now
	for _, address := range r.addresses {
		if address.TenantID != tenantID || address.ContactID != id || address.DeletedAt != nil {
			continue
		}
		address.DeletedAt = &now
		address.UpdatedAt = now
		delete(r.identityIndex, tenantIdentityKey(tenantID, address.Identity))
	}
	return nil
}

func (r *InMemoryContactRepository) CreateContactAddress(ctx context.Context, address *ContactAddress) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if address == nil {
		return fmt.Errorf("%w: address is nil", ErrInvalidContactAddress)
	}
	candidate, err := prepareAddress(address)
	if err != nil {
		return err
	}
	if candidate.PublicID == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.Source == "" || candidate.Source == ContactSourceUnknown {
		candidate.Source = ContactSourceManual
	}
	if candidate.Consent.State == "" || candidate.Consent.State == ConsentStateUnknown {
		candidate.Consent.State = ConsentStateUnknown
	}
	if candidate.Consent.Source == "" || candidate.Consent.Source == ConsentSourceUnknown {
		candidate.Consent.Source = ConsentSourceUnknown
	}
	if candidate.DeletedAt != nil {
		return fmt.Errorf("%w: deleted address cannot be created", ErrInvalidContactAddress)
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.UpdatedAt = time.Now().UTC()
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidContactAddress, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	contact, ok := r.contacts[candidate.ContactID]
	if !ok || contact.TenantID != candidate.TenantID || contact.DeletedAt != nil {
		return ErrContactNotFound
	}
	publicKey := scopedKey(candidate.TenantID, candidate.PublicID)
	identityKey := tenantIdentityKey(candidate.TenantID, candidate.Identity)
	if _, exists := r.addressByKey[publicKey]; exists {
		return ErrContactAddressConflict
	}
	if _, exists := r.identityIndex[identityKey]; exists {
		return ErrContactAddressConflict
	}
	candidate.ID = r.nextAddressID
	r.nextAddressID++
	r.addresses[candidate.ID] = candidate
	r.addressByKey[publicKey] = candidate.ID
	r.identityIndex[identityKey] = candidate.ID
	*address = *cloneContactAddress(candidate)
	return nil
}

func (r *InMemoryContactRepository) GetContactAddressByPublicID(ctx context.Context, tenantID, publicID string) (*ContactAddress, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.addressByKey[scopedKey(tenantID, publicID)]
	if !ok {
		return nil, ErrContactAddressNotFound
	}
	address, ok := r.addresses[id]
	if !ok || address.DeletedAt != nil {
		return nil, ErrContactAddressNotFound
	}
	return cloneContactAddress(address), nil
}

func (r *InMemoryContactRepository) GetContactAddressByID(ctx context.Context, tenantID string, id uint64) (*ContactAddress, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	address, ok := r.addresses[id]
	if !ok || address.TenantID != tenantID || address.DeletedAt != nil {
		return nil, ErrContactAddressNotFound
	}
	return cloneContactAddress(address), nil
}

func (r *InMemoryContactRepository) FindContactAddressByIdentity(ctx context.Context, tenantID string, identity AddressIdentity) (*ContactAddress, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	canonical, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidContactAddress, err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.identityIndex[tenantIdentityKey(tenantID, canonical)]
	if !ok {
		return nil, ErrContactAddressNotFound
	}
	address, ok := r.addresses[id]
	if !ok || address.DeletedAt != nil {
		return nil, ErrContactAddressNotFound
	}
	return cloneContactAddress(address), nil
}

func (r *InMemoryContactRepository) ListContactAddresses(ctx context.Context, tenantID string, contactID uint64) ([]*ContactAddress, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" || contactID == 0 {
		return nil, fmt.Errorf("%w: tenant and contact are required", ErrInvalidContactAddress)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	contact, ok := r.contacts[contactID]
	if !ok || contact.TenantID != tenantID || contact.DeletedAt != nil {
		return nil, ErrContactNotFound
	}
	result := make([]*ContactAddress, 0)
	for id := uint64(1); id < r.nextAddressID; id++ {
		address, ok := r.addresses[id]
		if ok && address.TenantID == tenantID && address.ContactID == contactID && address.DeletedAt == nil {
			result = append(result, cloneContactAddress(address))
		}
	}
	return result, nil
}

func (r *InMemoryContactRepository) UpdateContactAddress(ctx context.Context, address *ContactAddress) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if address == nil {
		return fmt.Errorf("%w: address is nil", ErrInvalidContactAddress)
	}
	candidate, err := prepareAddress(address)
	if err != nil {
		return err
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidContactAddress, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.addresses[candidate.ID]
	if !ok || stored.TenantID != candidate.TenantID || stored.DeletedAt != nil {
		return ErrContactAddressNotFound
	}
	if stored.PublicID != candidate.PublicID || stored.ContactID != candidate.ContactID {
		return ErrContactAddressConflict
	}
	contact, ok := r.contacts[candidate.ContactID]
	if !ok || contact.TenantID != candidate.TenantID || contact.DeletedAt != nil {
		return ErrContactNotFound
	}
	newKey := tenantIdentityKey(candidate.TenantID, candidate.Identity)
	if existingID, exists := r.identityIndex[newKey]; exists && existingID != candidate.ID {
		return ErrContactAddressConflict
	}
	oldKey := tenantIdentityKey(stored.TenantID, stored.Identity)
	delete(r.identityIndex, oldKey)
	candidate.DeletedAt = nil
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	r.addresses[candidate.ID] = candidate
	r.identityIndex[newKey] = candidate.ID
	*address = *cloneContactAddress(candidate)
	return nil
}

func (r *InMemoryContactRepository) DeleteContactAddress(ctx context.Context, tenantID, publicID string) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.addressByKey[scopedKey(tenantID, publicID)]
	if !ok {
		return ErrContactAddressNotFound
	}
	address, ok := r.addresses[id]
	if !ok || address.TenantID != tenantID || address.DeletedAt != nil {
		return ErrContactAddressNotFound
	}
	now := time.Now().UTC()
	address.DeletedAt = &now
	address.UpdatedAt = now
	delete(r.identityIndex, tenantIdentityKey(tenantID, address.Identity))
	return nil
}

// Short aliases make the in-memory implementation convenient in small tests.
func (r *InMemoryContactRepository) Create(ctx context.Context, contact *Contact) error {
	return r.CreateContact(ctx, contact)
}
func (r *InMemoryContactRepository) Get(ctx context.Context, tenantID, publicID string) (*Contact, error) {
	return r.GetContactByPublicID(ctx, tenantID, publicID)
}

func prepareAddress(address *ContactAddress) (*ContactAddress, error) {
	candidate := cloneContactAddress(address)
	kind := candidate.Kind
	namespace := candidate.Namespace
	value := candidate.Value
	identityWasProvided := candidate.Identity.Value != ""
	if kind == AddressKindUnknown {
		kind = candidate.Identity.Kind
	}
	if namespace == "" {
		namespace = candidate.Identity.Namespace
	}
	if value == "" {
		value = candidate.NormalizedValue
	}
	if value == "" {
		value = candidate.Identity.Value
	}
	identity, err := NormalizeAddress(kind, namespace, value)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidContactAddress, err)
	}
	if identityWasProvided && candidate.Value == "" && candidate.NormalizedValue == "" && candidate.Identity.Value != identity.Value {
		return nil, fmt.Errorf("%w: inconsistent address identity", ErrInvalidContactAddress)
	}
	candidate.Identity = identity
	candidate.Kind = identity.Kind
	candidate.Namespace = identity.Namespace
	candidate.NormalizedValue = identity.Value
	return candidate, nil
}

func cloneContact(contact *Contact) *Contact {
	if contact == nil {
		return nil
	}
	copy := *contact
	copy.SourceMetadata = cloneStringMap(contact.SourceMetadata)
	copy.CustomFields = cloneStringMap(contact.CustomFields)
	copy.DeletedAt = cloneTime(contact.DeletedAt)
	copy.Tags = append([]string(nil), contact.Tags...)
	return &copy
}

func cloneContactAddress(address *ContactAddress) *ContactAddress {
	if address == nil {
		return nil
	}
	copy := *address
	copy.Consent = address.Consent
	if address.Consent.OccurredAt != nil {
		occurredAt := *address.Consent.OccurredAt
		copy.Consent.OccurredAt = &occurredAt
	}
	copy.DeletedAt = cloneTime(address.DeletedAt)
	return &copy
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	copy := make(map[string]string, len(values))
	for key, value := range values {
		copy[key] = value
	}
	return copy
}
