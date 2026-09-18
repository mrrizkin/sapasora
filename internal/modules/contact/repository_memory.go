package contact

import (
	"context"
	"errors"
	"fmt"
	"sort"
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

	nextContactID      uint64
	nextAddressID      uint64
	nextConsentEventID uint64
	nextSuppressionID  uint64
	nextStaticListID   uint64
	nextSegmentID      uint64
	nextExclusionID    uint64
	nextSnapshotID     uint64
	nextMergeID        uint64
	contacts           map[uint64]*Contact
	contactByKey       map[string]uint64
	addresses          map[uint64]*ContactAddress
	addressByKey       map[string]uint64
	identityIndex      map[string]uint64
	consentEvents      map[uint64][]*ConsentEvent
	suppressions       map[uint64]*SuppressionRecord
	staticLists        map[uint64]*StaticList
	staticListByKey    map[string]uint64
	segments           map[uint64]*DynamicSegment
	segmentByKey       map[string]uint64
	exclusionLists     map[uint64]*ExclusionList
	exclusionByKey     map[string]uint64
	snapshots          map[uint64]*AudienceSnapshot
	snapshotByKey      map[string]uint64
	mergeAudits        map[uint64]*storedMergeAudit
	mergeAuditByKey    map[string]uint64
	mergeUndos         map[string]*MergeUndoAudit
}

var _ ContactRepository = (*InMemoryContactRepository)(nil)

// InMemoryRepository is the concise in-memory implementation name.
type InMemoryRepository = InMemoryContactRepository

func NewInMemoryContactRepository() *InMemoryContactRepository {
	return &InMemoryContactRepository{
		nextContactID:      1,
		nextAddressID:      1,
		nextConsentEventID: 1,
		nextSuppressionID:  1,
		nextStaticListID:   1,
		nextSegmentID:      1,
		nextExclusionID:    1,
		nextSnapshotID:     1,
		nextMergeID:        1,
		contacts:           make(map[uint64]*Contact),
		contactByKey:       make(map[string]uint64),
		addresses:          make(map[uint64]*ContactAddress),
		addressByKey:       make(map[string]uint64),
		identityIndex:      make(map[string]uint64),
		consentEvents:      make(map[uint64][]*ConsentEvent),
		suppressions:       make(map[uint64]*SuppressionRecord),
		staticLists:        make(map[uint64]*StaticList),
		staticListByKey:    make(map[string]uint64),
		segments:           make(map[uint64]*DynamicSegment),
		segmentByKey:       make(map[string]uint64),
		exclusionLists:     make(map[uint64]*ExclusionList),
		exclusionByKey:     make(map[string]uint64),
		snapshots:          make(map[uint64]*AudienceSnapshot),
		snapshotByKey:      make(map[string]uint64),
		mergeAudits:        make(map[uint64]*storedMergeAudit),
		mergeAuditByKey:    make(map[string]uint64),
		mergeUndos:         make(map[string]*MergeUndoAudit),
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
	if r.nextConsentEventID == 0 {
		r.nextConsentEventID = 1
	}
	if r.nextSuppressionID == 0 {
		r.nextSuppressionID = 1
	}
	if r.nextStaticListID == 0 {
		r.nextStaticListID = 1
	}
	if r.nextSegmentID == 0 {
		r.nextSegmentID = 1
	}
	if r.nextExclusionID == 0 {
		r.nextExclusionID = 1
	}
	if r.nextSnapshotID == 0 {
		r.nextSnapshotID = 1
	}
	if r.nextMergeID == 0 {
		r.nextMergeID = 1
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
	if r.consentEvents == nil {
		r.consentEvents = make(map[uint64][]*ConsentEvent)
	}
	if r.suppressions == nil {
		r.suppressions = make(map[uint64]*SuppressionRecord)
	}
	if r.staticLists == nil {
		r.staticLists = make(map[uint64]*StaticList)
	}
	if r.staticListByKey == nil {
		r.staticListByKey = make(map[string]uint64)
	}
	if r.segments == nil {
		r.segments = make(map[uint64]*DynamicSegment)
	}
	if r.segmentByKey == nil {
		r.segmentByKey = make(map[string]uint64)
	}
	if r.exclusionLists == nil {
		r.exclusionLists = make(map[uint64]*ExclusionList)
	}
	if r.exclusionByKey == nil {
		r.exclusionByKey = make(map[string]uint64)
	}
	if r.snapshots == nil {
		r.snapshots = make(map[uint64]*AudienceSnapshot)
	}
	if r.snapshotByKey == nil {
		r.snapshotByKey = make(map[string]uint64)
	}
	if r.mergeAudits == nil {
		r.mergeAudits = make(map[uint64]*storedMergeAudit)
	}
	if r.mergeAuditByKey == nil {
		r.mergeAuditByKey = make(map[string]uint64)
	}
	if r.mergeUndos == nil {
		r.mergeUndos = make(map[string]*MergeUndoAudit)
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
	if candidate.Consent.State != ConsentStateUnknown {
		if candidate.Consent.Source == ConsentSourceUnknown {
			candidate.Consent.Source = ConsentSourceManual
		}
		if candidate.Consent.OccurredAt == nil {
			occurredAt := candidate.CreatedAt
			candidate.Consent.OccurredAt = &occurredAt
		}
		r.appendConsentEventLocked(candidate, candidate.Consent, candidate.UpdatedAt)
		if candidate.Consent.State == ConsentStateOptedOut {
			r.upsertSuppressionLocked(candidate, SuppressionReasonOptOut, candidate.Consent.Source, candidate.Consent.OccurredAt, candidate.Consent.EvidenceRef, candidate.Consent.ActorID, candidate.UpdatedAt)
		}
	}
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
	if !consentMetadataEqual(candidate.Consent, stored.Consent) {
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

func (r *InMemoryContactRepository) TransitionContactAddressConsent(ctx context.Context, tenantID string, addressID uint64, transition ConsentTransition) (*ConsentEvent, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" || addressID == 0 {
		return nil, fmt.Errorf("%w: tenant and address are required", ErrInvalid)
	}
	if transition.State != ConsentStateOptedIn && transition.State != ConsentStateOptedOut {
		return nil, fmt.Errorf("%w: consent transition state is required", ErrInvalid)
	}
	if err := transition.Valid(); err != nil || transition.Source == ConsentSourceUnknown {
		return nil, fmt.Errorf("%w: invalid consent transition", ErrInvalid)
	}
	occurredAt := time.Now().UTC()
	if transition.OccurredAt != nil {
		if transition.OccurredAt.IsZero() {
			return nil, fmt.Errorf("%w: consent timestamp is invalid", ErrInvalid)
		}
		occurredAt = transition.OccurredAt.UTC()
	}
	transition.OccurredAt = &occurredAt

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	address, ok := r.addresses[addressID]
	if !ok || address.TenantID != tenantID || address.DeletedAt != nil {
		return nil, ErrContactAddressNotFound
	}
	if address.Consent.OccurredAt != nil && occurredAt.Before(address.Consent.OccurredAt.UTC()) {
		return nil, fmt.Errorf("%w: consent transition is older than current state", ErrConflict)
	}
	recordedAt := time.Now().UTC()
	event := &ConsentEvent{
		ID:         r.nextConsentEventID,
		TenantID:   tenantID,
		AddressID:  addressID,
		Sequence:   uint64(len(r.consentEvents[addressID]) + 1),
		State:      transition.State,
		Metadata:   cloneConsentMetadata(transition),
		RecordedAt: recordedAt,
	}
	r.nextConsentEventID++
	r.consentEvents[addressID] = append(r.consentEvents[addressID], event)
	address.Consent = cloneConsentMetadata(transition)
	address.UpdatedAt = recordedAt
	if transition.State == ConsentStateOptedOut {
		r.upsertSuppressionLocked(address, SuppressionReasonOptOut, transition.Source, transition.OccurredAt, transition.EvidenceRef, transition.ActorID, recordedAt)
	} else {
		r.resolveSuppressionLocked(tenantID, address.Identity, SuppressionReasonOptOut, recordedAt)
	}
	return cloneConsentEvent(event), nil
}

func (r *InMemoryContactRepository) ListConsentEvents(ctx context.Context, tenantID string, addressID uint64) ([]*ConsentEvent, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" || addressID == 0 {
		return nil, fmt.Errorf("%w: tenant and address are required", ErrInvalid)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	address, ok := r.addresses[addressID]
	if !ok || address.TenantID != tenantID || address.DeletedAt != nil {
		return nil, ErrContactAddressNotFound
	}
	events := r.consentEvents[addressID]
	result := make([]*ConsentEvent, 0, len(events))
	for _, event := range events {
		result = append(result, cloneConsentEvent(event))
	}
	return result, nil
}

func (r *InMemoryContactRepository) CreateSuppression(ctx context.Context, suppression *SuppressionRecord) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if suppression == nil {
		return fmt.Errorf("%w: suppression is nil", ErrInvalid)
	}
	candidate := cloneSuppression(suppression)
	identity, err := NormalizeAddress(candidate.Identity.Kind, candidate.Identity.Namespace, candidate.Identity.Value)
	if err != nil {
		return fmt.Errorf("%w: invalid suppression identity", ErrInvalid)
	}
	candidate.Identity = identity
	if candidate.Source == ConsentSourceUnknown {
		candidate.Source = ConsentSourceManual
	}
	if candidate.OccurredAt == nil {
		now := time.Now().UTC()
		candidate.OccurredAt = &now
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	addressID, ok := r.identityIndex[tenantIdentityKey(candidate.TenantID, candidate.Identity)]
	address, addressExists := r.addresses[addressID]
	if !ok || !addressExists || address.DeletedAt != nil {
		return ErrContactAddressNotFound
	}
	for id := uint64(1); id < r.nextSuppressionID; id++ {
		existing := r.suppressions[id]
		if existing != nil && existing.TenantID == candidate.TenantID && existing.Identity == candidate.Identity && existing.Reason == candidate.Reason && existing.Active() {
			*suppression = *cloneSuppression(existing)
			return nil
		}
	}
	candidate.ID = r.nextSuppressionID
	r.nextSuppressionID++
	r.suppressions[candidate.ID] = candidate
	*suppression = *cloneSuppression(candidate)
	return nil
}

func (r *InMemoryContactRepository) ListSuppressions(ctx context.Context, tenantID string, identity AddressIdentity) ([]*SuppressionRecord, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	canonical, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid suppression identity", ErrInvalid)
	}
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant id is required", ErrInvalid)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*SuppressionRecord, 0)
	for id := uint64(1); id < r.nextSuppressionID; id++ {
		suppression := r.suppressions[id]
		if suppression != nil && suppression.TenantID == tenantID && suppression.Identity == canonical {
			result = append(result, cloneSuppression(suppression))
		}
	}
	return result, nil
}

func (r *InMemoryContactRepository) ResolveSuppression(ctx context.Context, tenantID string, identity AddressIdentity, reason SuppressionReason) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	canonical, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value)
	if err != nil || !reason.Valid() {
		return fmt.Errorf("%w: invalid suppression", ErrInvalid)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.resolveSuppressionLocked(tenantID, canonical, reason, time.Now().UTC()) {
		return ErrContactAddressNotFound
	}
	return nil
}

func (r *InMemoryContactRepository) IsAddressSendable(ctx context.Context, tenantID string, identity AddressIdentity) (bool, error) {
	if err := repositoryContextError(ctx); err != nil {
		return false, err
	}
	canonical, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value)
	if err != nil {
		return false, fmt.Errorf("%w: invalid address identity", ErrInvalid)
	}
	if strings.TrimSpace(tenantID) == "" {
		return false, fmt.Errorf("%w: tenant id is required", ErrInvalid)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	addressID, ok := r.identityIndex[tenantIdentityKey(tenantID, canonical)]
	address, addressExists := r.addresses[addressID]
	if !ok || !addressExists || address.DeletedAt != nil {
		return false, ErrContactAddressNotFound
	}
	for _, suppression := range r.suppressions {
		if suppression.TenantID == tenantID && suppression.Identity == canonical && suppression.Active() {
			return false, nil
		}
	}
	return address.Consent.State == ConsentStateOptedIn, nil
}

// MergeContact atomically applies an explicitly confirmed merge. The
// repository recomputes the preview token while holding its write lock, so a
// concurrent address/contact change can only fail safely with ErrConflict.
func (r *InMemoryContactRepository) MergeContact(ctx context.Context, tenantID string, request MergeRequest) (*MergeAudit, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(request.ConfirmationToken) == "" {
		return nil, ErrInvalid
	}
	if !request.Confirmed {
		return nil, ErrMergeConfirmationRequired
	}
	if err := validateMergeRequestIDs(request); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	sourceID, sourceOK := r.contactByKey[scopedKey(tenantID, request.SourceContactPublicID)]
	targetID, targetOK := r.contactByKey[scopedKey(tenantID, request.TargetContactPublicID)]
	source, sourceExists := r.contacts[sourceID]
	target, targetExists := r.contacts[targetID]
	if !sourceOK || !targetOK || !sourceExists || !targetExists || source.DeletedAt != nil || target.DeletedAt != nil {
		return nil, ErrContactNotFound
	}
	sourceAddresses := r.activeAddressesLocked(tenantID, source.ID)
	targetAddresses := r.activeAddressesLocked(tenantID, target.ID)
	preview := newMergePreview(tenantID, source, target, sourceAddresses, targetAddresses)
	if request.ConfirmationToken != preview.ConfirmationToken {
		return nil, ErrConflict
	}
	if !preview.CanMerge {
		return nil, ErrConflict
	}
	now := time.Now().UTC()
	audit := &MergeAudit{ID: fmt.Sprintf("merge-%d", r.nextMergeID), TenantID: tenantID, SourceContactID: source.ID, SourceContactPublicID: source.PublicID, TargetContactID: target.ID, TargetContactPublicID: target.PublicID, ConfirmationToken: request.ConfirmationToken, CreatedAt: now, Undoable: true}
	r.nextMergeID++
	for _, address := range sourceAddresses {
		if address == nil {
			continue
		}
		audit.Addresses = append(audit.Addresses, MergeAddressMetadata{AddressID: address.ID, ContactID: source.ID, IdentityKey: address.Identity.Key()})
		stored := r.addresses[address.ID]
		stored.ContactID = target.ID
		stored.UpdatedAt = now
	}
	source.DeletedAt = &now
	source.UpdatedAt = now
	r.mergeAudits[r.nextMergeID-1] = &storedMergeAudit{audit: cloneMergeAudit(audit), addresses: append([]MergeAddressMetadata(nil), audit.Addresses...)}
	r.mergeAuditByKey[audit.ID] = r.nextMergeID - 1
	return cloneMergeAudit(audit), nil
}

func (r *InMemoryContactRepository) UndoContactMerge(ctx context.Context, tenantID, mergeID string) (*MergeUndoAudit, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(mergeID) == "" {
		return nil, ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	mergeNumber, ok := r.mergeAuditByKey[mergeID]
	if !ok {
		return nil, ErrNotFound
	}
	storedAudit := r.mergeAudits[mergeNumber]
	if storedAudit == nil || storedAudit.audit.TenantID != tenantID {
		return nil, ErrNotFound
	}
	if _, alreadyUndone := r.mergeUndos[mergeID]; alreadyUndone {
		return nil, ErrConflict
	}
	source := r.contacts[storedAudit.audit.SourceContactID]
	target := r.contacts[storedAudit.audit.TargetContactID]
	if source == nil || target == nil || source.TenantID != tenantID || target.TenantID != tenantID || source.DeletedAt == nil || target.DeletedAt != nil {
		return nil, ErrConflict
	}
	for _, metadata := range storedAudit.addresses {
		address := r.addresses[metadata.AddressID]
		if address == nil || address.TenantID != tenantID || address.DeletedAt != nil || address.ContactID != target.ID || address.Identity.Key() != metadata.IdentityKey {
			return nil, ErrConflict
		}
	}
	now := time.Now().UTC()
	undo := &MergeUndoAudit{ID: fmt.Sprintf("%s-undo", mergeID), MergeID: mergeID, TenantID: tenantID, SourceContactID: source.ID, TargetContactID: target.ID, CreatedAt: now}
	for _, metadata := range storedAudit.addresses {
		address := r.addresses[metadata.AddressID]
		address.ContactID = source.ID
		address.UpdatedAt = now
		undo.RestoredAddressIDs = append(undo.RestoredAddressIDs, address.ID)
	}
	source.DeletedAt = nil
	source.UpdatedAt = now
	r.mergeUndos[mergeID] = cloneMergeUndoAudit(undo)
	return cloneMergeUndoAudit(undo), nil
}

func (r *InMemoryContactRepository) ListMergeAudits(ctx context.Context, tenantID string) ([]*MergeAudit, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) == "" {
		return nil, ErrInvalid
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*MergeAudit, 0)
	for id := uint64(1); id < r.nextMergeID; id++ {
		stored := r.mergeAudits[id]
		if stored == nil || stored.audit.TenantID != tenantID {
			continue
		}
		audit := cloneMergeAudit(stored.audit)
		_, audit.Undone = r.mergeUndos[audit.ID]
		audit.Undoable = !audit.Undone
		result = append(result, audit)
	}
	return result, nil
}

func (r *InMemoryContactRepository) activeAddressesLocked(tenantID string, contactID uint64) []*ContactAddress {
	addresses := make([]*ContactAddress, 0)
	for id := uint64(1); id < r.nextAddressID; id++ {
		address := r.addresses[id]
		if address != nil && address.TenantID == tenantID && address.ContactID == contactID && address.DeletedAt == nil {
			addresses = append(addresses, cloneContactAddress(address))
		}
	}
	sort.Slice(addresses, func(i, j int) bool { return addresses[i].ID < addresses[j].ID })
	return addresses
}

func (r *InMemoryContactRepository) appendConsentEventLocked(address *ContactAddress, metadata ConsentMetadata, recordedAt time.Time) {
	event := &ConsentEvent{
		ID:         r.nextConsentEventID,
		TenantID:   address.TenantID,
		AddressID:  address.ID,
		Sequence:   uint64(len(r.consentEvents[address.ID]) + 1),
		State:      metadata.State,
		Metadata:   cloneConsentMetadata(metadata),
		RecordedAt: recordedAt,
	}
	r.nextConsentEventID++
	r.consentEvents[address.ID] = append(r.consentEvents[address.ID], event)
}

func (r *InMemoryContactRepository) upsertSuppressionLocked(address *ContactAddress, reason SuppressionReason, source ConsentSource, occurredAt *time.Time, evidenceRef, actorID string, createdAt time.Time) {
	for id := uint64(1); id < r.nextSuppressionID; id++ {
		existing := r.suppressions[id]
		if existing != nil && existing.TenantID == address.TenantID && existing.Identity == address.Identity && existing.Reason == reason && existing.Active() {
			return
		}
	}
	record := &SuppressionRecord{ID: r.nextSuppressionID, TenantID: address.TenantID, Identity: address.Identity, Reason: reason, Source: source, OccurredAt: cloneTime(occurredAt), EvidenceRef: evidenceRef, ActorID: actorID, CreatedAt: createdAt}
	r.nextSuppressionID++
	r.suppressions[record.ID] = record
}

func (r *InMemoryContactRepository) resolveSuppressionLocked(tenantID string, identity AddressIdentity, reason SuppressionReason, resolvedAt time.Time) bool {
	resolved := false
	for _, suppression := range r.suppressions {
		if suppression.TenantID == tenantID && suppression.Identity == identity && suppression.Reason == reason && suppression.Active() {
			at := resolvedAt
			suppression.ResolvedAt = &at
			resolved = true
		}
	}
	return resolved
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
	copy.Consent = cloneConsentMetadata(address.Consent)
	copy.DeletedAt = cloneTime(address.DeletedAt)
	return &copy
}

func cloneConsentMetadata(metadata ConsentMetadata) ConsentMetadata {
	metadata.OccurredAt = cloneTime(metadata.OccurredAt)
	return metadata
}

func consentMetadataEqual(left, right ConsentMetadata) bool {
	if left.State != right.State || left.Source != right.Source || left.EvidenceRef != right.EvidenceRef || left.ActorID != right.ActorID {
		return false
	}
	if left.OccurredAt == nil || right.OccurredAt == nil {
		return left.OccurredAt == nil && right.OccurredAt == nil
	}
	return left.OccurredAt.Equal(*right.OccurredAt)
}

func cloneConsentEvent(event *ConsentEvent) *ConsentEvent {
	if event == nil {
		return nil
	}
	copy := *event
	copy.Metadata = cloneConsentMetadata(event.Metadata)
	return &copy
}

func cloneSuppression(suppression *SuppressionRecord) *SuppressionRecord {
	if suppression == nil {
		return nil
	}
	copy := *suppression
	copy.Identity = suppression.Identity
	copy.OccurredAt = cloneTime(suppression.OccurredAt)
	copy.ResolvedAt = cloneTime(suppression.ResolvedAt)
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
