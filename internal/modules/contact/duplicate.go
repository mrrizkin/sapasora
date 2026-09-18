package contact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// DuplicateMatchRule is the deliberately narrow matching policy used by this
// foundation. Only an exact canonical address identity is a match; fuzzy name
// or provider-specific heuristics are intentionally deferred.
type DuplicateMatchRule string

const DuplicateMatchRuleExactAddressIdentity DuplicateMatchRule = "exact_address_identity"

// DuplicateIdentityMatch is a safe diagnostic projection. IdentityKey is a
// non-PII fingerprint; the normalized value is never returned in diagnostics.
type DuplicateIdentityMatch struct {
	Kind         AddressKind        `json:"kind"`
	Namespace    string             `json:"namespace"`
	IdentityKey  string             `json:"identity_key"`
	MatchingRule DuplicateMatchRule `json:"matching_rule"`
}

// DuplicateCandidate identifies a tenant-local contact that shares one or
// more exact address identities with a proposed identity set.
type DuplicateCandidate struct {
	TenantID             string                   `json:"tenant_id"`
	ContactID            uint64                   `json:"contact_id"`
	ContactPublicID      string                   `json:"contact_public_id"`
	ContactStatus        ContactStatus            `json:"contact_status"`
	MatchedIdentityCount int                      `json:"matched_identity_count"`
	Matches              []DuplicateIdentityMatch `json:"matches"`
}

// MergeRequest requires an explicit confirmation and a token from a current
// merge preview. The token prevents a stale preview from authorizing a merge.
type MergeRequest struct {
	SourceContactPublicID string `json:"source_contact_public_id"`
	TargetContactPublicID string `json:"target_contact_public_id"`
	ConfirmationToken     string `json:"confirmation_token"`
	Confirmed             bool   `json:"confirmed"`
}

// MergeConflict is safe to log and serialize. It deliberately contains no
// display names, raw values, normalized values, or provider identifiers.
type MergeConflict struct {
	Code            string      `json:"code"`
	Kind            AddressKind `json:"kind,omitempty"`
	Namespace       string      `json:"namespace,omitempty"`
	IdentityKey     string      `json:"identity_key,omitempty"`
	SourceContactID uint64      `json:"source_contact_id"`
	TargetContactID uint64      `json:"target_contact_id"`
	SourceAddressID uint64      `json:"source_address_id,omitempty"`
	TargetAddressID uint64      `json:"target_address_id,omitempty"`
}

// MergePreview is a transport-neutral, safe preview. It is not an approval;
// MergeContacts still requires Confirmed=true and the current token.
type MergePreview struct {
	TenantID              string          `json:"tenant_id"`
	SourceContactID       uint64          `json:"source_contact_id"`
	SourceContactPublicID string          `json:"source_contact_public_id"`
	TargetContactID       uint64          `json:"target_contact_id"`
	TargetContactPublicID string          `json:"target_contact_public_id"`
	SourceAddressCount    int             `json:"source_address_count"`
	TargetAddressCount    int             `json:"target_address_count"`
	SourceAddressIDs      []uint64        `json:"source_address_ids"`
	TargetAddressIDs      []uint64        `json:"target_address_ids"`
	Conflicts             []MergeConflict `json:"conflicts,omitempty"`
	CanMerge              bool            `json:"can_merge"`
	ConfirmationToken     string          `json:"confirmation_token,omitempty"`
}

// MergeAddressMetadata is immutable undo metadata. It contains only stable
// identifiers and identity fingerprints, never address values.
type MergeAddressMetadata struct {
	AddressID   uint64 `json:"address_id"`
	ContactID   uint64 `json:"contact_id"`
	IdentityKey string `json:"identity_key"`
}

// MergeAudit is an immutable record of an accepted merge. Undo state is kept
// in a separate append-only MergeUndoAudit rather than mutating this record.
type MergeAudit struct {
	ID                    string                 `json:"id"`
	TenantID              string                 `json:"tenant_id"`
	SourceContactID       uint64                 `json:"source_contact_id"`
	SourceContactPublicID string                 `json:"source_contact_public_id"`
	TargetContactID       uint64                 `json:"target_contact_id"`
	TargetContactPublicID string                 `json:"target_contact_public_id"`
	Addresses             []MergeAddressMetadata `json:"addresses"`
	ConfirmationToken     string                 `json:"confirmation_token"`
	CreatedAt             time.Time              `json:"created_at"`
	Undoable              bool                   `json:"undoable"`
	Undone                bool                   `json:"undone"`
}

// MergeUndoAudit is append-only metadata for a successful safe undo.
type MergeUndoAudit struct {
	ID                 string    `json:"id"`
	MergeID            string    `json:"merge_id"`
	TenantID           string    `json:"tenant_id"`
	SourceContactID    uint64    `json:"source_contact_id"`
	TargetContactID    uint64    `json:"target_contact_id"`
	RestoredAddressIDs []uint64  `json:"restored_address_ids"`
	CreatedAt          time.Time `json:"created_at"`
}

func (m MergeAudit) String() string {
	return fmt.Sprintf("MergeAudit{id=%s tenant_id=%s source_contact_id=%d target_contact_id=%d address_count=%d}", m.ID, m.TenantID, m.SourceContactID, m.TargetContactID, len(m.Addresses))
}

func (m MergeUndoAudit) String() string {
	return fmt.Sprintf("MergeUndoAudit{id=%s merge_id=%s tenant_id=%s restored_address_count=%d}", m.ID, m.MergeID, m.TenantID, len(m.RestoredAddressIDs))
}

type storedMergeAudit struct {
	audit     *MergeAudit
	addresses []MergeAddressMetadata
}

// MergeRepository is an optional atomic capability of a contact repository.
// Keeping it separate preserves compatibility with CRUD-only adapters while
// requiring durable adapters to implement merge and undo transactionally.
type MergeRepository interface {
	MergeContact(context.Context, string, MergeRequest) (*MergeAudit, error)
	UndoContactMerge(context.Context, string, string) (*MergeUndoAudit, error)
	ListMergeAudits(context.Context, string) ([]*MergeAudit, error)
}

// FindDuplicateCandidates matches proposed identities against active contacts
// in exactly one tenant. Inputs are canonicalized before lookup and results
// are deterministic by contact ID and identity key.
func (s *Service) FindDuplicateCandidates(ctx context.Context, tenantID string, identities []AddressIdentity) ([]*DuplicateCandidate, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is nil", ErrInvalid)
	}
	canonical, err := canonicalIdentities(identities)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid duplicate identity", ErrInvalid)
	}
	if len(canonical) == 0 {
		return nil, fmt.Errorf("%w: at least one identity is required", ErrInvalid)
	}
	return s.findDuplicateCandidates(ctx, tenantID, canonical, 0)
}

// ListDuplicateCandidates is a descriptive alias for integrations that use
// list terminology.
func (s *Service) ListDuplicateCandidates(ctx context.Context, tenantID string, identities []AddressIdentity) ([]*DuplicateCandidate, error) {
	return s.FindDuplicateCandidates(ctx, tenantID, identities)
}

// FindContactDuplicateCandidates checks the identities already attached to a
// contact, excluding that contact itself. It is useful for provider sync and
// import previews; the current uniqueness invariant normally makes this empty
// for an already persisted contact.
func (s *Service) FindContactDuplicateCandidates(ctx context.Context, tenantID string, contactPublicID string) ([]*DuplicateCandidate, error) {
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
	addresses, err := s.repository.ListContactAddresses(ctx, tenantID, contact.ID)
	if err != nil {
		return nil, safeContactError(err)
	}
	identities := make([]AddressIdentity, 0, len(addresses))
	for _, address := range addresses {
		if address != nil {
			identities = append(identities, address.Identity)
		}
	}
	if len(identities) == 0 {
		return []*DuplicateCandidate{}, nil
	}
	canonical, err := canonicalIdentities(identities)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid stored duplicate identity", ErrInvalid)
	}
	return s.findDuplicateCandidates(ctx, tenantID, canonical, contact.ID)
}

func (s *Service) findDuplicateCandidates(ctx context.Context, tenantID string, identities []AddressIdentity, excludedContactID uint64) ([]*DuplicateCandidate, error) {
	byContact := make(map[uint64]*DuplicateCandidate)
	for _, identity := range identities {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		address, err := s.repository.FindContactAddressByIdentity(ctx, tenantID, identity)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			}
			return nil, safeContactError(err)
		}
		if address == nil || address.ContactID == 0 || address.ContactID == excludedContactID {
			continue
		}
		candidate, ok := byContact[address.ContactID]
		if !ok {
			contact, getErr := s.repository.GetContactByID(ctx, tenantID, address.ContactID)
			if getErr != nil {
				return nil, safeContactError(getErr)
			}
			candidate = &DuplicateCandidate{TenantID: tenantID, ContactID: contact.ID, ContactPublicID: contact.PublicID, ContactStatus: contact.Status}
			byContact[address.ContactID] = candidate
		}
		candidate.Matches = append(candidate.Matches, DuplicateIdentityMatch{Kind: identity.Kind, Namespace: identity.Namespace, IdentityKey: identity.Key(), MatchingRule: DuplicateMatchRuleExactAddressIdentity})
	}
	result := make([]*DuplicateCandidate, 0, len(byContact))
	for _, candidate := range byContact {
		sort.Slice(candidate.Matches, func(i, j int) bool { return candidate.Matches[i].IdentityKey < candidate.Matches[j].IdentityKey })
		candidate.MatchedIdentityCount = len(candidate.Matches)
		result = append(result, candidate)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ContactID != result[j].ContactID {
			return result[i].ContactID < result[j].ContactID
		}
		return result[i].ContactPublicID < result[j].ContactPublicID
	})
	return result, nil
}

func (s *Service) PreviewMerge(ctx context.Context, tenantID string, request MergeRequest) (*MergePreview, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validateMergeRequestIDs(request); err != nil {
		return nil, err
	}
	source, err := s.repository.GetContactByPublicID(ctx, tenantID, request.SourceContactPublicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	target, err := s.repository.GetContactByPublicID(ctx, tenantID, request.TargetContactPublicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	sourceAddresses, err := s.repository.ListContactAddresses(ctx, tenantID, source.ID)
	if err != nil {
		return nil, safeContactError(err)
	}
	targetAddresses, err := s.repository.ListContactAddresses(ctx, tenantID, target.ID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return newMergePreview(tenantID, source, target, sourceAddresses, targetAddresses), nil
}

// PreviewContactMerge is a descriptive alias for transport integrations.
func (s *Service) PreviewContactMerge(ctx context.Context, tenantID string, request MergeRequest) (*MergePreview, error) {
	return s.PreviewMerge(ctx, tenantID, request)
}

func (s *Service) MergeContacts(ctx context.Context, tenantID string, request MergeRequest) (*MergeAudit, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if !request.Confirmed {
		return nil, ErrMergeConfirmationRequired
	}
	if strings.TrimSpace(request.ConfirmationToken) == "" {
		return nil, fmt.Errorf("%w: merge preview token is required", ErrInvalid)
	}
	preview, err := s.PreviewMerge(ctx, tenantID, request)
	if err != nil {
		return nil, err
	}
	if !preview.CanMerge || preview.ConfirmationToken != request.ConfirmationToken {
		return nil, ErrConflict
	}
	mergeRepository, ok := s.repository.(MergeRepository)
	if !ok {
		return nil, fmt.Errorf("%w: merge repository capability is unavailable", ErrInvalid)
	}
	audit, err := mergeRepository.MergeContact(ctx, tenantID, request)
	if err != nil {
		return nil, safeContactError(err)
	}
	return audit, nil
}

func (s *Service) UndoMerge(ctx context.Context, tenantID, mergeID string) (*MergeUndoAudit, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(mergeID) == "" {
		return nil, fmt.Errorf("%w: merge id is required", ErrInvalid)
	}
	mergeRepository, ok := s.repository.(MergeRepository)
	if !ok {
		return nil, fmt.Errorf("%w: merge repository capability is unavailable", ErrInvalid)
	}
	undo, err := mergeRepository.UndoContactMerge(ctx, tenantID, mergeID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return undo, nil
}

// UndoContactMerge is a descriptive alias for transport integrations.
func (s *Service) UndoContactMerge(ctx context.Context, tenantID, mergeID string) (*MergeUndoAudit, error) {
	return s.UndoMerge(ctx, tenantID, mergeID)
}

func (s *Service) ListMergeAudits(ctx context.Context, tenantID string) ([]*MergeAudit, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	mergeRepository, ok := s.repository.(MergeRepository)
	if !ok {
		return nil, fmt.Errorf("%w: merge repository capability is unavailable", ErrInvalid)
	}
	audits, err := mergeRepository.ListMergeAudits(ctx, tenantID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return audits, nil
}

func validateMergeRequestIDs(request MergeRequest) error {
	if err := validatePublicID(request.SourceContactPublicID); err != nil {
		return err
	}
	if err := validatePublicID(request.TargetContactPublicID); err != nil {
		return err
	}
	return nil
}

func canonicalIdentities(identities []AddressIdentity) ([]AddressIdentity, error) {
	seen := make(map[string]AddressIdentity, len(identities))
	for _, identity := range identities {
		canonical, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value)
		if err != nil {
			return nil, err
		}
		seen[canonical.Key()] = canonical
	}
	result := make([]AddressIdentity, 0, len(seen))
	for _, identity := range seen {
		result = append(result, identity)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key() < result[j].Key() })
	return result, nil
}

func newMergePreview(tenantID string, source, target *Contact, sourceAddresses, targetAddresses []*ContactAddress) *MergePreview {
	preview := &MergePreview{TenantID: tenantID, SourceContactID: source.ID, SourceContactPublicID: source.PublicID, TargetContactID: target.ID, TargetContactPublicID: target.PublicID, SourceAddressCount: len(sourceAddresses), TargetAddressCount: len(targetAddresses), CanMerge: true}
	for _, address := range sourceAddresses {
		if address != nil {
			preview.SourceAddressIDs = append(preview.SourceAddressIDs, address.ID)
		}
	}
	for _, address := range targetAddresses {
		if address != nil {
			preview.TargetAddressIDs = append(preview.TargetAddressIDs, address.ID)
		}
	}
	sort.Slice(preview.SourceAddressIDs, func(i, j int) bool { return preview.SourceAddressIDs[i] < preview.SourceAddressIDs[j] })
	sort.Slice(preview.TargetAddressIDs, func(i, j int) bool { return preview.TargetAddressIDs[i] < preview.TargetAddressIDs[j] })
	if source.ID == target.ID || source.PublicID == target.PublicID {
		preview.Conflicts = append(preview.Conflicts, MergeConflict{Code: "same_contact", SourceContactID: source.ID, TargetContactID: target.ID})
	}
	targetByIdentity := make(map[string]*ContactAddress, len(targetAddresses))
	for _, address := range targetAddresses {
		if address != nil {
			targetByIdentity[address.Identity.Key()] = address
		}
	}
	if source.ID != target.ID {
		for _, address := range sourceAddresses {
			if address == nil {
				continue
			}
			if targetAddress, ok := targetByIdentity[address.Identity.Key()]; ok {
				preview.Conflicts = append(preview.Conflicts, MergeConflict{Code: "identity_conflict", Kind: address.Identity.Kind, Namespace: address.Identity.Namespace, IdentityKey: address.Identity.Key(), SourceContactID: source.ID, TargetContactID: target.ID, SourceAddressID: address.ID, TargetAddressID: targetAddress.ID})
			}
		}
	}
	if len(preview.Conflicts) != 0 {
		preview.CanMerge = false
	}
	preview.ConfirmationToken = mergeConfirmationToken(tenantID, source, target, sourceAddresses, targetAddresses)
	return preview
}

func mergeConfirmationToken(tenantID string, source, target *Contact, sourceAddresses, targetAddresses []*ContactAddress) string {
	parts := []string{tenantID, fmt.Sprintf("source:%d:%d", source.ID, source.UpdatedAt.UnixNano()), fmt.Sprintf("target:%d:%d", target.ID, target.UpdatedAt.UnixNano())}
	appendAddressParts := func(prefix string, addresses []*ContactAddress) {
		values := make([]string, 0, len(addresses))
		for _, address := range addresses {
			if address != nil {
				values = append(values, fmt.Sprintf("%s:%d:%s:%d", prefix, address.ID, address.Identity.Key(), address.ContactID))
			}
		}
		sort.Strings(values)
		parts = append(parts, values...)
	}
	appendAddressParts("source", sourceAddresses)
	appendAddressParts("target", targetAddresses)
	digest := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(digest[:])
}

func cloneMergeAudit(value *MergeAudit) *MergeAudit {
	if value == nil {
		return nil
	}
	copy := *value
	copy.Addresses = append([]MergeAddressMetadata(nil), value.Addresses...)
	return &copy
}

func cloneMergeUndoAudit(value *MergeUndoAudit) *MergeUndoAudit {
	if value == nil {
		return nil
	}
	copy := *value
	copy.RestoredAddressIDs = append([]uint64(nil), value.RestoredAddressIDs...)
	return &copy
}
