package contact

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ContactService is the use-case boundary for tenant-scoped contact CRUD.
// Callers provide the tenant explicitly; a model's TenantID is treated as an
// ownership assertion and may not select a different tenant.
type ContactService interface {
	CreateContact(context.Context, string, *Contact) error
	GetContact(context.Context, string, string) (*Contact, error)
	GetContactByPublicID(context.Context, string, string) (*Contact, error)
	GetContactByID(context.Context, string, uint64) (*Contact, error)
	ListContacts(context.Context, string, ContactFilter) ([]*Contact, error)
	UpdateContact(context.Context, string, *Contact) error
	DeleteContact(context.Context, string, string) error

	CreateContactAddress(context.Context, string, *ContactAddress) error
	GetContactAddress(context.Context, string, string) (*ContactAddress, error)
	GetContactAddressByPublicID(context.Context, string, string) (*ContactAddress, error)
	GetContactAddressByID(context.Context, string, uint64) (*ContactAddress, error)
	ListContactAddresses(context.Context, string, string) ([]*ContactAddress, error)
	ListContactAddressesByContactID(context.Context, string, uint64) ([]*ContactAddress, error)
	UpdateContactAddress(context.Context, string, *ContactAddress) error
	DeleteContactAddress(context.Context, string, string) error

	TransitionConsent(context.Context, string, string, ConsentTransition) error
	OptIn(context.Context, string, string, ConsentMetadata) error
	OptOut(context.Context, string, string, ConsentMetadata) error
	ListConsentEvents(context.Context, string, string) ([]*ConsentEvent, error)
	SuppressAddress(context.Context, string, string, *SuppressionRecord) error
	ListSuppressions(context.Context, string, string) ([]*SuppressionRecord, error)
	BlockAddress(context.Context, string, string, ConsentMetadata) error
	UnsuppressAddress(context.Context, string, string, SuppressionReason) error
	IsAddressSendable(context.Context, string, AddressIdentity) (bool, error)
	IsContactAddressSendable(context.Context, string, string) (bool, error)
}

// ContactUseCase is the explicit use-case name for integrations that prefer
// not to call application services "services".
type ContactUseCase = ContactService

// Service coordinates validation and tenant ownership checks before calling a
// ContactRepository. It contains no transport or authorization policy; the
// authenticated tenant must be selected by the integration boundary.
type Service struct {
	repository ContactRepository
}

// ContactServiceImpl is retained as a descriptive implementation alias.
type ContactServiceImpl = Service

var _ ContactService = (*Service)(nil)

// NewContactService constructs the contact CRUD use-case boundary.
// @wired:provide
func NewContactService(repository ContactRepository) ContactService {
	return &Service{repository: repository}
}

// NewService is a concise constructor for package-local composition.
func NewService(repository ContactRepository) ContactService {
	return NewContactService(repository)
}

// NewContactUseCase is an explicit constructor alias for application wiring.
func NewContactUseCase(repository ContactRepository) ContactUseCase {
	return NewContactService(repository)
}

func (s *Service) ready() error {
	if s == nil || s.repository == nil {
		return fmt.Errorf("%w: contact repository is unavailable", ErrInvalid)
	}
	return nil
}

func validateTenant(tenantID string) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant id is required", ErrInvalid)
	}
	return nil
}

func validateOwnership(tenantID, modelTenantID string) error {
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if modelTenantID != "" && modelTenantID != tenantID {
		return fmt.Errorf("%w: resource belongs to another tenant", ErrInvalid)
	}
	return nil
}

func validatePublicID(publicID string) error {
	if strings.TrimSpace(publicID) == "" {
		return fmt.Errorf("%w: public id is required", ErrInvalid)
	}
	return nil
}

func (s *Service) CreateContact(ctx context.Context, tenantID string, value *Contact) error {
	if err := s.ready(); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: contact is nil", ErrInvalid)
	}
	if err := validateOwnership(tenantID, value.TenantID); err != nil {
		return err
	}
	candidate := cloneContact(value)
	candidate.TenantID = tenantID
	if candidate.Status == "" || candidate.Status == ContactStatusUnknown {
		candidate.Status = ContactStatusActive
	}
	if candidate.Source == "" || candidate.Source == ContactSourceUnknown {
		candidate.Source = ContactSourceManual
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if err := s.repository.CreateContact(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*value = *cloneContact(candidate)
	return nil
}

func (s *Service) GetContact(ctx context.Context, tenantID, publicID string) (*Contact, error) {
	return s.GetContactByPublicID(ctx, tenantID, publicID)
}

func (s *Service) GetContactByPublicID(ctx context.Context, tenantID, publicID string) (*Contact, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(publicID); err != nil {
		return nil, err
	}
	value, err := s.repository.GetContactByPublicID(ctx, tenantID, publicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

func (s *Service) GetContactByID(ctx context.Context, tenantID string, id uint64) (*Contact, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, fmt.Errorf("%w: contact id is required", ErrInvalid)
	}
	value, err := s.repository.GetContactByID(ctx, tenantID, id)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

func (s *Service) ListContacts(ctx context.Context, tenantID string, filter ContactFilter) ([]*Contact, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	values, err := s.repository.ListContacts(ctx, tenantID, filter)
	if err != nil {
		return nil, safeContactError(err)
	}
	return values, nil
}

func (s *Service) UpdateContact(ctx context.Context, tenantID string, value *Contact) error {
	if err := s.ready(); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: contact is nil", ErrInvalid)
	}
	if err := validateOwnership(tenantID, value.TenantID); err != nil {
		return err
	}
	candidate := cloneContact(value)
	candidate.TenantID = tenantID
	if err := validatePublicID(candidate.PublicID); err != nil {
		return err
	}
	if candidate.ID == 0 {
		stored, err := s.repository.GetContactByPublicID(ctx, tenantID, candidate.PublicID)
		if err != nil {
			return safeContactError(err)
		}
		candidate.ID = stored.ID
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if err := s.repository.UpdateContact(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*value = *cloneContact(candidate)
	return nil
}

func (s *Service) DeleteContact(ctx context.Context, tenantID, publicID string) error {
	if err := s.ready(); err != nil {
		return err
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if err := validatePublicID(publicID); err != nil {
		return err
	}
	return safeContactError(s.repository.DeleteContact(ctx, tenantID, publicID))
}

func (s *Service) CreateContactAddress(ctx context.Context, tenantID string, value *ContactAddress) error {
	if err := s.ready(); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: address is nil", ErrInvalid)
	}
	if err := validateOwnership(tenantID, value.TenantID); err != nil {
		return err
	}
	candidate, err := s.prepareAddressForWrite(tenantID, value)
	if err != nil {
		return err
	}
	if candidate.ContactID == 0 {
		return fmt.Errorf("%w: contact id is required", ErrInvalid)
	}
	if _, err := s.repository.GetContactByID(ctx, tenantID, candidate.ContactID); err != nil {
		return safeContactError(err)
	}
	if err := s.repository.CreateContactAddress(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*value = *cloneContactAddress(candidate)
	return nil
}

func (s *Service) GetContactAddress(ctx context.Context, tenantID, publicID string) (*ContactAddress, error) {
	return s.GetContactAddressByPublicID(ctx, tenantID, publicID)
}

func (s *Service) GetContactAddressByPublicID(ctx context.Context, tenantID, publicID string) (*ContactAddress, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(publicID); err != nil {
		return nil, err
	}
	value, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, publicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

func (s *Service) GetContactAddressByID(ctx context.Context, tenantID string, id uint64) (*ContactAddress, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, fmt.Errorf("%w: address id is required", ErrInvalid)
	}
	value, err := s.repository.GetContactAddressByID(ctx, tenantID, id)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

func (s *Service) ListContactAddresses(ctx context.Context, tenantID, contactPublicID string) ([]*ContactAddress, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(contactPublicID); err != nil {
		return nil, err
	}
	contact, err := s.repository.GetContactByPublicID(ctx, tenantID, contactPublicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	values, err := s.repository.ListContactAddresses(ctx, tenantID, contact.ID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return values, nil
}

func (s *Service) ListContactAddressesByContactID(ctx context.Context, tenantID string, contactID uint64) ([]*ContactAddress, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if contactID == 0 {
		return nil, fmt.Errorf("%w: contact id is required", ErrInvalid)
	}
	values, err := s.repository.ListContactAddresses(ctx, tenantID, contactID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return values, nil
}

func (s *Service) UpdateContactAddress(ctx context.Context, tenantID string, value *ContactAddress) error {
	if err := s.ready(); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: address is nil", ErrInvalid)
	}
	if err := validateOwnership(tenantID, value.TenantID); err != nil {
		return err
	}
	if err := validatePublicID(value.PublicID); err != nil {
		return err
	}
	stored, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, value.PublicID)
	if err != nil {
		return safeContactError(err)
	}
	input := cloneContactAddress(value)
	if input.ContactID == 0 {
		input.ContactID = stored.ContactID
	}
	if input.Consent.State == ConsentStateUnknown && input.Consent.Source == ConsentSourceUnknown && input.Consent.OccurredAt == nil && input.Consent.EvidenceRef == "" && input.Consent.ActorID == "" {
		input.Consent = cloneConsentMetadata(stored.Consent)
	}
	candidate, err := s.prepareAddressForWrite(tenantID, input)
	if err != nil {
		return err
	}
	if candidate.ID == 0 {
		candidate.ID = stored.ID
	}
	if candidate.ContactID == 0 {
		candidate.ContactID = stored.ContactID
	}
	if candidate.ID != stored.ID || candidate.ContactID != stored.ContactID {
		return ErrConflict
	}
	if _, err := s.repository.GetContactByID(ctx, tenantID, candidate.ContactID); err != nil {
		return safeContactError(err)
	}
	if err := s.repository.UpdateContactAddress(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*value = *cloneContactAddress(candidate)
	return nil
}

func (s *Service) DeleteContactAddress(ctx context.Context, tenantID, publicID string) error {
	if err := s.ready(); err != nil {
		return err
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if err := validatePublicID(publicID); err != nil {
		return err
	}
	return safeContactError(s.repository.DeleteContactAddress(ctx, tenantID, publicID))
}

func (s *Service) TransitionConsent(ctx context.Context, tenantID, addressPublicID string, transition ConsentTransition) error {
	if err := s.ready(); err != nil {
		return err
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if err := validatePublicID(addressPublicID); err != nil {
		return err
	}
	address, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, addressPublicID)
	if err != nil {
		return safeContactError(err)
	}
	if transition.Source == ConsentSourceUnknown {
		transition.Source = ConsentSourceManual
	}
	if err := transition.Valid(); err != nil || transition.State == ConsentStateUnknown {
		return fmt.Errorf("%w: invalid consent transition", ErrInvalid)
	}
	_, err = s.repository.TransitionContactAddressConsent(ctx, tenantID, address.ID, transition)
	return safeContactError(err)
}

func (s *Service) OptIn(ctx context.Context, tenantID, addressPublicID string, metadata ConsentMetadata) error {
	metadata.State = ConsentStateOptedIn
	return s.TransitionConsent(ctx, tenantID, addressPublicID, metadata)
}

func (s *Service) OptOut(ctx context.Context, tenantID, addressPublicID string, metadata ConsentMetadata) error {
	metadata.State = ConsentStateOptedOut
	return s.TransitionConsent(ctx, tenantID, addressPublicID, metadata)
}

func (s *Service) ListConsentEvents(ctx context.Context, tenantID, addressPublicID string) ([]*ConsentEvent, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(addressPublicID); err != nil {
		return nil, err
	}
	address, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, addressPublicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	events, err := s.repository.ListConsentEvents(ctx, tenantID, address.ID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return events, nil
}

func (s *Service) SuppressAddress(ctx context.Context, tenantID, addressPublicID string, value *SuppressionRecord) error {
	if err := s.ready(); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: suppression is nil", ErrInvalid)
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if err := validatePublicID(addressPublicID); err != nil {
		return err
	}
	address, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, addressPublicID)
	if err != nil {
		return safeContactError(err)
	}
	candidate := cloneSuppression(value)
	candidate.TenantID = tenantID
	candidate.Identity = address.Identity
	if candidate.Source == ConsentSourceUnknown {
		candidate.Source = ConsentSourceManual
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if err := s.repository.CreateSuppression(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*value = *cloneSuppression(candidate)
	return nil
}

func (s *Service) ListSuppressions(ctx context.Context, tenantID, addressPublicID string) ([]*SuppressionRecord, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(addressPublicID); err != nil {
		return nil, err
	}
	address, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, addressPublicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	values, err := s.repository.ListSuppressions(ctx, tenantID, address.Identity)
	if err != nil {
		return nil, safeContactError(err)
	}
	return values, nil
}

func (s *Service) BlockAddress(ctx context.Context, tenantID, addressPublicID string, metadata ConsentMetadata) error {
	if metadata.Source == ConsentSourceUnknown {
		metadata.Source = ConsentSourceManual
	}
	record := &SuppressionRecord{Reason: SuppressionReasonBlocked, Source: metadata.Source, OccurredAt: cloneTime(metadata.OccurredAt), EvidenceRef: metadata.EvidenceRef, ActorID: metadata.ActorID}
	return s.SuppressAddress(ctx, tenantID, addressPublicID, record)
}

func (s *Service) UnsuppressAddress(ctx context.Context, tenantID, addressPublicID string, reason SuppressionReason) error {
	if err := s.ready(); err != nil {
		return err
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if err := validatePublicID(addressPublicID); err != nil {
		return err
	}
	address, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, addressPublicID)
	if err != nil {
		return safeContactError(err)
	}
	return safeContactError(s.repository.ResolveSuppression(ctx, tenantID, address.Identity, reason))
}

func (s *Service) IsAddressSendable(ctx context.Context, tenantID string, identity AddressIdentity) (bool, error) {
	if err := s.ready(); err != nil {
		return false, err
	}
	if err := validateTenant(tenantID); err != nil {
		return false, err
	}
	canonical, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value)
	if err != nil {
		return false, fmt.Errorf("%w: invalid address identity", ErrInvalid)
	}
	value, err := s.repository.IsAddressSendable(ctx, tenantID, canonical)
	if err != nil {
		return false, safeContactError(err)
	}
	return value, nil
}

func (s *Service) IsContactAddressSendable(ctx context.Context, tenantID, addressPublicID string) (bool, error) {
	if err := s.ready(); err != nil {
		return false, err
	}
	if err := validateTenant(tenantID); err != nil {
		return false, err
	}
	if err := validatePublicID(addressPublicID); err != nil {
		return false, err
	}
	address, err := s.repository.GetContactAddressByPublicID(ctx, tenantID, addressPublicID)
	if err != nil {
		return false, safeContactError(err)
	}
	value, err := s.repository.IsAddressSendable(ctx, tenantID, address.Identity)
	if err != nil {
		return false, safeContactError(err)
	}
	return value, nil
}

func (s *Service) prepareAddressForWrite(tenantID string, value *ContactAddress) (*ContactAddress, error) {
	candidate, err := prepareAddress(value)
	if err != nil {
		return nil, err
	}
	candidate.TenantID = tenantID
	if candidate.Source == "" || candidate.Source == ContactSourceUnknown {
		candidate.Source = ContactSourceManual
	}
	if candidate.Consent.State == "" || candidate.Consent.State == ConsentStateUnknown {
		candidate.Consent.State = ConsentStateUnknown
	}
	if candidate.Consent.Source == "" || candidate.Consent.Source == ConsentSourceUnknown {
		candidate.Consent.Source = ConsentSourceUnknown
	}
	if err := candidate.Valid(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	return candidate, nil
}

func safeContactError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound
	case errors.Is(err, ErrConflict):
		return ErrConflict
	case errors.Is(err, ErrInvalid):
		return ErrInvalid
	case errors.Is(err, context.Canceled):
		return context.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return context.DeadlineExceeded
	default:
		return fmt.Errorf("%w: contact repository operation failed", ErrInvalid)
	}
}
