package contact

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"sapasora/platform/support/hash"
)

// AudienceHardLimit is the maximum number of contacts or recipients that one
// evaluation may inspect. It is deliberately finite so a preview cannot turn
// into an unbounded query. Event, delivery-status, and interaction-time
// predicates are intentionally not part of this foundation.
const AudienceHardLimit = 10000

var (
	ErrAudienceEvaluationLimit = errors.New("audience evaluation limit exceeded")
	ErrAudienceImmutable       = errors.New("audience snapshot is immutable")
)

// StaticList is a tenant-owned, manually maintained contact list.
type StaticList struct {
	ID         uint64    `json:"id"`
	PublicID   string    `json:"public_id"`
	TenantID   string    `json:"tenant_id"`
	Name       string    `json:"-"`
	ContactIDs []uint64  `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (l StaticList) Valid() error {
	if strings.TrimSpace(l.TenantID) == "" || strings.TrimSpace(l.Name) == "" {
		return fmt.Errorf("tenant and list name are required")
	}
	if len(l.Name) > 200 {
		return fmt.Errorf("list name is too long")
	}
	return nil
}

func (l StaticList) String() string {
	return fmt.Sprintf("StaticList{public_id=%s contacts=%d}", l.PublicID, len(l.ContactIDs))
}

func (l StaticList) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           uint64    `json:"id"`
		PublicID     string    `json:"public_id"`
		TenantID     string    `json:"tenant_id"`
		ContactCount int       `json:"contact_count"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}{l.ID, l.PublicID, l.TenantID, len(l.ContactIDs), l.CreatedAt, l.UpdatedAt})
}

// SegmentField names are an allowlist, not an expression language. The
// contact.custom.* form is allowed only for a named custom field. Address
// fields are evaluated against a contact's linked addresses.
const (
	SegmentFieldStatus           = "contact.status"
	SegmentFieldSource           = "contact.source"
	SegmentFieldLocale           = "contact.locale"
	SegmentFieldTimezone         = "contact.timezone"
	SegmentFieldOwnerID          = "contact.owner_id"
	SegmentFieldLifecycleStage   = "contact.lifecycle_stage"
	SegmentFieldAddressKind      = "address.kind"
	SegmentFieldAddressNamespace = "address.namespace"
	SegmentFieldAddressProvider  = "address.provider"
	SegmentFieldAddressConsent   = "address.consent_state"
)

type PredicateOperator string

const (
	PredicateEquals     PredicateOperator = "equals"
	PredicateNotEquals  PredicateOperator = "not_equals"
	PredicateContains   PredicateOperator = "contains"
	PredicateStartsWith PredicateOperator = "starts_with"
	PredicateExists     PredicateOperator = "exists"
	PredicateIn         PredicateOperator = "in"
)

type SegmentMatch string

const (
	SegmentMatchAll SegmentMatch = "all"
	SegmentMatchAny SegmentMatch = "any"
)

// FieldPredicate is a safe comparison against one allowlisted contact or
// address field. Value and Values are never included in diagnostic output.
type FieldPredicate struct {
	Field    string            `json:"field"`
	Operator PredicateOperator `json:"operator"`
	Value    string            `json:"-"`
	Values   []string          `json:"-"`
}

func (p FieldPredicate) Valid() error {
	if !allowedSegmentField(p.Field) {
		return fmt.Errorf("unsupported segment field")
	}
	if !validPredicateOperator(p.Operator) {
		return fmt.Errorf("unsupported segment operator")
	}
	if len(p.Value) > 256 || len(p.Values) > 32 {
		return fmt.Errorf("segment predicate value is too large")
	}
	for _, value := range p.Values {
		if len(value) > 256 {
			return fmt.Errorf("segment predicate value is too large")
		}
	}
	if p.Operator == PredicateIn && len(p.Values) == 0 {
		return fmt.Errorf("in predicate requires values")
	}
	if p.Operator != PredicateExists && p.Operator != PredicateIn && strings.TrimSpace(p.Value) == "" {
		return fmt.Errorf("predicate value is required")
	}
	return nil
}

func (p FieldPredicate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Field      string            `json:"field"`
		Operator   PredicateOperator `json:"operator"`
		ValueSet   bool              `json:"value_present"`
		ValueCount int               `json:"value_count"`
	}{p.Field, p.Operator, p.Value != "" || len(p.Values) > 0, len(p.Values)})
}

type TagOperator string

const (
	TagHas    TagOperator = "has"
	TagNotHas TagOperator = "not_has"
)

type TagPredicate struct {
	Tag      string      `json:"-"`
	Operator TagOperator `json:"operator"`
}

func (p TagPredicate) Valid() error {
	if strings.TrimSpace(p.Tag) == "" || len(p.Tag) > 128 {
		return fmt.Errorf("tag is required and must be short")
	}
	if p.Operator != TagHas && p.Operator != TagNotHas {
		return fmt.Errorf("unsupported tag operator")
	}
	return nil
}

func (p TagPredicate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Operator   TagOperator `json:"operator"`
		TagPresent bool        `json:"tag_present"`
	}{p.Operator, strings.TrimSpace(p.Tag) != ""})
}

// SegmentDefinition contains only bounded field and tag predicates. It does
// not accept event, delivery-status, or interaction-time filters.
type SegmentDefinition struct {
	Match  SegmentMatch     `json:"match"`
	Fields []FieldPredicate `json:"fields"`
	Tags   []TagPredicate   `json:"tags"`
	// FieldPredicates and TagPredicates are descriptive aliases accepted by the
	// evaluator for callers that prefer the longer names.
	FieldPredicates []FieldPredicate `json:"-"`
	TagPredicates   []TagPredicate   `json:"-"`
	MaxCandidates   int              `json:"max_candidates"`
}

func (d SegmentDefinition) predicates() ([]FieldPredicate, []TagPredicate) {
	fields := append([]FieldPredicate(nil), d.Fields...)
	fields = append(fields, d.FieldPredicates...)
	tags := append([]TagPredicate(nil), d.Tags...)
	tags = append(tags, d.TagPredicates...)
	return fields, tags
}

func (d SegmentDefinition) Valid() error {
	if d.Match == "" {
		d.Match = SegmentMatchAll
	}
	if d.Match != SegmentMatchAll && d.Match != SegmentMatchAny {
		return fmt.Errorf("invalid segment match mode")
	}
	if d.MaxCandidates < 0 || d.MaxCandidates > AudienceHardLimit {
		return fmt.Errorf("invalid segment evaluation bound")
	}
	fields, tags := d.predicates()
	if len(fields) == 0 && len(tags) == 0 {
		return fmt.Errorf("segment requires at least one predicate")
	}
	if len(fields)+len(tags) > 32 {
		return fmt.Errorf("too many segment predicates")
	}
	for _, predicate := range fields {
		if err := predicate.Valid(); err != nil {
			return err
		}
	}
	for _, predicate := range tags {
		if err := predicate.Valid(); err != nil {
			return err
		}
	}
	return nil
}

func (d SegmentDefinition) MarshalJSON() ([]byte, error) {
	fields, tags := d.predicates()
	fieldNames := make([]string, 0, len(fields))
	for _, field := range fields {
		fieldNames = append(fieldNames, field.Field)
	}
	return json.Marshal(struct {
		Match          SegmentMatch `json:"match"`
		Fields         []string     `json:"fields"`
		TagCount       int          `json:"tag_count"`
		PredicateCount int          `json:"predicate_count"`
		MaxCandidates  int          `json:"max_candidates"`
	}{d.Match, fieldNames, len(tags), len(fields) + len(tags), d.MaxCandidates})
}

// DynamicSegment is a tenant-owned, reusable segment definition.
type DynamicSegment struct {
	ID         uint64            `json:"id"`
	PublicID   string            `json:"public_id"`
	TenantID   string            `json:"tenant_id"`
	Name       string            `json:"-"`
	Definition SegmentDefinition `json:"definition"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type Segment = DynamicSegment

func (s DynamicSegment) Valid() error {
	if strings.TrimSpace(s.TenantID) == "" || strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("tenant and segment name are required")
	}
	if len(s.Name) > 200 {
		return fmt.Errorf("segment name is too long")
	}
	return s.Definition.Valid()
}

func (s DynamicSegment) String() string {
	return fmt.Sprintf("DynamicSegment{public_id=%s}", s.PublicID)
}

// ExclusionList holds tenant-scoped contact or address exclusions.
type ExclusionList struct {
	ID                uint64
	PublicID          string
	TenantID          string
	ContactIDs        []uint64
	AddressIdentities []AddressIdentity
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (e ExclusionList) Valid() error {
	if strings.TrimSpace(e.TenantID) == "" {
		return fmt.Errorf("tenant id is required")
	}
	if len(e.ContactIDs)+len(e.AddressIdentities) > AudienceHardLimit {
		return fmt.Errorf("exclusion list is too large")
	}
	for _, identity := range e.AddressIdentities {
		if _, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value); err != nil {
			return fmt.Errorf("invalid exclusion identity")
		}
	}
	return nil
}

func (e ExclusionList) String() string {
	return fmt.Sprintf("ExclusionList{public_id=%s contacts=%d addresses=%d}", e.PublicID, len(e.ContactIDs), len(e.AddressIdentities))
}

func (e ExclusionList) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           uint64    `json:"id"`
		PublicID     string    `json:"public_id"`
		TenantID     string    `json:"tenant_id"`
		ContactCount int       `json:"contact_count"`
		AddressCount int       `json:"address_count"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}{e.ID, e.PublicID, e.TenantID, len(e.ContactIDs), len(e.AddressIdentities), e.CreatedAt, e.UpdatedAt})
}

type AudienceSelection struct {
	StaticListPublicID    string
	SegmentPublicID       string
	ExclusionListPublicID string
	MaxContacts           int
	MaxRecipients         int
}

func (s AudienceSelection) Valid() error {
	if (strings.TrimSpace(s.StaticListPublicID) == "") == (strings.TrimSpace(s.SegmentPublicID) == "") {
		return fmt.Errorf("exactly one static list or segment is required")
	}
	if s.MaxContacts < 0 || s.MaxRecipients < 0 || s.MaxContacts > AudienceHardLimit || s.MaxRecipients > AudienceHardLimit {
		return fmt.Errorf("invalid audience evaluation bound")
	}
	return nil
}

// AudienceRecipient contains only tenant-local IDs plus an identity used for
// suppression checks. AddressIdentity's JSON projection never emits its value.
type AudienceRecipient struct {
	ContactID uint64          `json:"contact_id"`
	AddressID uint64          `json:"address_id"`
	Identity  AddressIdentity `json:"identity"`
}

type AudiencePreview struct {
	TenantID            string `json:"tenant_id"`
	CandidateContacts   int    `json:"candidate_contacts"`
	IncludedContacts    int    `json:"included_contacts"`
	CandidateRecipients int    `json:"candidate_recipients"`
	IncludedRecipients  int    `json:"included_recipients"`
	ExcludedContacts    int    `json:"excluded_contacts"`
	ExcludedRecipients  int    `json:"excluded_recipients"`
	Count               int    `json:"count"`
}

func (p AudiencePreview) String() string {
	return fmt.Sprintf("AudiencePreview{tenant=%s contacts=%d recipients=%d excluded=%d}", p.TenantID, p.IncludedContacts, p.IncludedRecipients, p.ExcludedRecipients)
}

type AudienceSnapshot struct {
	ID         uint64              `json:"id"`
	PublicID   string              `json:"public_id"`
	TenantID   string              `json:"tenant_id"`
	Selection  AudienceSelection   `json:"-"`
	Recipients []AudienceRecipient `json:"recipients"`
	CreatedAt  time.Time           `json:"created_at"`
}

func (s AudienceSnapshot) Valid() error {
	if strings.TrimSpace(s.TenantID) == "" || strings.TrimSpace(s.PublicID) == "" {
		return fmt.Errorf("snapshot ownership is required")
	}
	if err := s.Selection.Valid(); err != nil {
		return err
	}
	if len(s.Recipients) > AudienceHardLimit {
		return fmt.Errorf("snapshot is too large")
	}
	return nil
}

func (s AudienceSnapshot) String() string {
	return fmt.Sprintf("AudienceSnapshot{public_id=%s recipients=%d immutable=true}", s.PublicID, len(s.Recipients))
}

type SuppressionRevalidation struct {
	SnapshotID string              `json:"snapshot_id"`
	Sendable   []AudienceRecipient `json:"sendable"`
	Suppressed []AudienceRecipient `json:"suppressed"`
}

func (r SuppressionRevalidation) String() string {
	return fmt.Sprintf("SuppressionRevalidation{snapshot=%s sendable=%d suppressed=%d}", r.SnapshotID, len(r.Sendable), len(r.Suppressed))
}

// AudienceRepository is the persistence boundary for audience definitions and
// immutable snapshots. It is separate from ContactRepository so existing
// contact adapters remain source-compatible while adding this track.
type AudienceRepository interface {
	CreateStaticList(context.Context, *StaticList) error
	GetStaticList(context.Context, string, string) (*StaticList, error)
	UpdateStaticList(context.Context, *StaticList) error
	CreateDynamicSegment(context.Context, *DynamicSegment) error
	GetDynamicSegment(context.Context, string, string) (*DynamicSegment, error)
	UpdateDynamicSegment(context.Context, *DynamicSegment) error
	CreateExclusionList(context.Context, *ExclusionList) error
	GetExclusionList(context.Context, string, string) (*ExclusionList, error)
	UpdateExclusionList(context.Context, *ExclusionList) error
	CreateAudienceSnapshot(context.Context, *AudienceSnapshot) error
	GetAudienceSnapshot(context.Context, string, string) (*AudienceSnapshot, error)
}

// AudienceService evaluates audiences against the contact repository and
// stores definitions/snapshots through the audience repository. It has no HTTP,
// UI, provider, or send-worker integration.
type AudienceService struct {
	contacts  ContactRepository
	audiences AudienceRepository
}

func NewAudienceService(contacts ContactRepository, stores ...AudienceRepository) *AudienceService {
	var audiences AudienceRepository
	if len(stores) > 0 {
		audiences = stores[0]
	} else if candidate, ok := contacts.(AudienceRepository); ok {
		audiences = candidate
	}
	return &AudienceService{contacts: contacts, audiences: audiences}
}

func NewContactAudienceService(contacts ContactRepository, stores ...AudienceRepository) *AudienceService {
	return NewAudienceService(contacts, stores...)
}

func (s *AudienceService) ready() error {
	if s == nil || s.contacts == nil || s.audiences == nil {
		return fmt.Errorf("%w: audience repositories are unavailable", ErrInvalid)
	}
	return nil
}

func (s *AudienceService) CreateStaticList(ctx context.Context, tenantID string, list *StaticList) error {
	if err := s.ready(); err != nil {
		return err
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if list == nil {
		return fmt.Errorf("%w: static list is nil", ErrInvalid)
	}
	candidate := cloneStaticList(list)
	if err := validateOwnership(tenantID, candidate.TenantID); err != nil {
		return err
	}
	candidate.TenantID = tenantID
	candidate.ContactIDs, _ = uniqueIDs(candidate.ContactIDs)
	for _, id := range candidate.ContactIDs {
		if _, err := s.contacts.GetContactByID(ctx, tenantID, id); err != nil {
			return safeContactError(err)
		}
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if err := s.audiences.CreateStaticList(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*list = *cloneStaticList(candidate)
	return nil
}

func (s *AudienceService) GetStaticList(ctx context.Context, tenantID, publicID string) (*StaticList, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(publicID); err != nil {
		return nil, err
	}
	value, err := s.audiences.GetStaticList(ctx, tenantID, publicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

func (s *AudienceService) SetStaticListMembers(ctx context.Context, tenantID, publicID string, contactIDs []uint64) error {
	list, err := s.GetStaticList(ctx, tenantID, publicID)
	if err != nil {
		return err
	}
	list.ContactIDs, _ = uniqueIDs(contactIDs)
	for _, id := range list.ContactIDs {
		if _, err := s.contacts.GetContactByID(ctx, tenantID, id); err != nil {
			return safeContactError(err)
		}
	}
	return safeContactError(s.audiences.UpdateStaticList(ctx, list))
}

func (s *AudienceService) CreateDynamicSegment(ctx context.Context, tenantID string, segment *DynamicSegment) error {
	if err := s.ready(); err != nil {
		return err
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if segment == nil {
		return fmt.Errorf("%w: segment is nil", ErrInvalid)
	}
	candidate := cloneDynamicSegment(segment)
	if err := validateOwnership(tenantID, candidate.TenantID); err != nil {
		return err
	}
	candidate.TenantID = tenantID
	if candidate.Definition.Match == "" {
		candidate.Definition.Match = SegmentMatchAll
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if err := s.audiences.CreateDynamicSegment(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*segment = *cloneDynamicSegment(candidate)
	return nil
}

func (s *AudienceService) GetDynamicSegment(ctx context.Context, tenantID, publicID string) (*DynamicSegment, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(publicID); err != nil {
		return nil, err
	}
	value, err := s.audiences.GetDynamicSegment(ctx, tenantID, publicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

func (s *AudienceService) PreviewSegmentCount(ctx context.Context, tenantID, segmentPublicID string) (int, error) {
	preview, err := s.PreviewAudience(ctx, tenantID, AudienceSelection{SegmentPublicID: segmentPublicID})
	if err != nil {
		return 0, err
	}
	return preview.IncludedContacts, nil
}

func (s *AudienceService) CreateExclusionList(ctx context.Context, tenantID string, exclusion *ExclusionList) error {
	if err := s.ready(); err != nil {
		return err
	}
	if err := validateTenant(tenantID); err != nil {
		return err
	}
	if exclusion == nil {
		return fmt.Errorf("%w: exclusion list is nil", ErrInvalid)
	}
	candidate := cloneExclusionList(exclusion)
	if err := validateOwnership(tenantID, candidate.TenantID); err != nil {
		return err
	}
	candidate.TenantID = tenantID
	candidate.ContactIDs, _ = uniqueIDs(candidate.ContactIDs)
	for _, id := range candidate.ContactIDs {
		if _, err := s.contacts.GetContactByID(ctx, tenantID, id); err != nil {
			return safeContactError(err)
		}
	}
	for index, identity := range candidate.AddressIdentities {
		canonical, err := NormalizeAddress(identity.Kind, identity.Namespace, identity.Value)
		if err != nil {
			return fmt.Errorf("%w: invalid exclusion identity", ErrInvalid)
		}
		candidate.AddressIdentities[index] = canonical
	}
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if err := s.audiences.CreateExclusionList(ctx, candidate); err != nil {
		return safeContactError(err)
	}
	*exclusion = *cloneExclusionList(candidate)
	return nil
}

func (s *AudienceService) GetExclusionList(ctx context.Context, tenantID, publicID string) (*ExclusionList, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(publicID); err != nil {
		return nil, err
	}
	value, err := s.audiences.GetExclusionList(ctx, tenantID, publicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

func (s *AudienceService) PreviewAudience(ctx context.Context, tenantID string, selection AudienceSelection) (AudiencePreview, error) {
	if err := s.ready(); err != nil {
		return AudiencePreview{}, err
	}
	if err := validateTenant(tenantID); err != nil {
		return AudiencePreview{}, err
	}
	if err := selection.Valid(); err != nil {
		return AudiencePreview{}, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	contactIDs, err := s.resolveContactIDs(ctx, tenantID, selection)
	if err != nil {
		return AudiencePreview{}, err
	}
	recipients, excludedContacts, excludedRecipients, err := s.buildRecipients(ctx, tenantID, contactIDs, selection)
	if err != nil {
		return AudiencePreview{}, err
	}
	includedContacts := len(contactIDs) - excludedContacts
	return AudiencePreview{TenantID: tenantID, CandidateContacts: len(contactIDs), IncludedContacts: includedContacts, CandidateRecipients: len(recipients) + excludedRecipients, IncludedRecipients: len(recipients), ExcludedContacts: excludedContacts, ExcludedRecipients: excludedRecipients, Count: len(recipients)}, nil
}

func (s *AudienceService) CreateAudienceSnapshot(ctx context.Context, tenantID string, selection AudienceSelection) (*AudienceSnapshot, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := selection.Valid(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	contactIDs, err := s.resolveContactIDs(ctx, tenantID, selection)
	if err != nil {
		return nil, err
	}
	recipients, _, _, err := s.buildRecipients(ctx, tenantID, contactIDs, selection)
	if err != nil {
		return nil, err
	}
	limit := selection.MaxRecipients
	if limit == 0 {
		limit = AudienceHardLimit
	}
	if len(recipients) > limit {
		return nil, ErrAudienceEvaluationLimit
	}
	snapshot := &AudienceSnapshot{TenantID: tenantID, Selection: selection, Recipients: recipients, CreatedAt: time.Now().UTC()}
	if err := s.audiences.CreateAudienceSnapshot(ctx, snapshot); err != nil {
		return nil, safeContactError(err)
	}
	return cloneAudienceSnapshot(snapshot), nil
}

func (s *AudienceService) GetAudienceSnapshot(ctx context.Context, tenantID, publicID string) (*AudienceSnapshot, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if err := validateTenant(tenantID); err != nil {
		return nil, err
	}
	if err := validatePublicID(publicID); err != nil {
		return nil, err
	}
	value, err := s.audiences.GetAudienceSnapshot(ctx, tenantID, publicID)
	if err != nil {
		return nil, safeContactError(err)
	}
	return value, nil
}

// RevalidateSnapshotForSend is the last-mile suppression/consent check. It
// does not mutate the immutable snapshot and deliberately leaves actual send
// execution to a later send-worker track.
func (s *AudienceService) RevalidateSnapshotForSend(ctx context.Context, tenantID, snapshotPublicID string) (SuppressionRevalidation, error) {
	snapshot, err := s.GetAudienceSnapshot(ctx, tenantID, snapshotPublicID)
	if err != nil {
		return SuppressionRevalidation{}, err
	}
	result := SuppressionRevalidation{SnapshotID: snapshot.PublicID, Sendable: make([]AudienceRecipient, 0, len(snapshot.Recipients)), Suppressed: make([]AudienceRecipient, 0)}
	for _, recipient := range snapshot.Recipients {
		if err := ctx.Err(); err != nil {
			return SuppressionRevalidation{}, err
		}
		sendable, err := s.contacts.IsAddressSendable(ctx, tenantID, recipient.Identity)
		if errors.Is(err, ErrNotFound) {
			result.Suppressed = append(result.Suppressed, recipient)
			continue
		}
		if err != nil {
			return SuppressionRevalidation{}, safeContactError(err)
		}
		if sendable {
			result.Sendable = append(result.Sendable, recipient)
		} else {
			result.Suppressed = append(result.Suppressed, recipient)
		}
	}
	return result, nil
}

func (s *AudienceService) Preview(ctx context.Context, tenantID string, selection AudienceSelection) (AudiencePreview, error) {
	return s.PreviewAudience(ctx, tenantID, selection)
}

func (s *AudienceService) SnapshotAudience(ctx context.Context, tenantID string, selection AudienceSelection) (*AudienceSnapshot, error) {
	return s.CreateAudienceSnapshot(ctx, tenantID, selection)
}

func (s *AudienceService) resolveContactIDs(ctx context.Context, tenantID string, selection AudienceSelection) ([]uint64, error) {
	limit := selection.MaxContacts
	if limit == 0 {
		limit = AudienceHardLimit
	}
	if selection.StaticListPublicID != "" {
		list, err := s.GetStaticList(ctx, tenantID, selection.StaticListPublicID)
		if err != nil {
			return nil, err
		}
		if len(list.ContactIDs) > limit {
			return nil, ErrAudienceEvaluationLimit
		}
		ids, _ := uniqueIDs(list.ContactIDs)
		return ids, nil
	}
	segment, err := s.GetDynamicSegment(ctx, tenantID, selection.SegmentPublicID)
	if err != nil {
		return nil, err
	}
	if segment.Definition.MaxCandidates > 0 && segment.Definition.MaxCandidates < limit {
		limit = segment.Definition.MaxCandidates
	}
	contacts, err := s.contacts.ListContacts(ctx, tenantID, ContactFilter{})
	if err != nil {
		return nil, safeContactError(err)
	}
	if len(contacts) > limit {
		return nil, ErrAudienceEvaluationLimit
	}
	ids := make([]uint64, 0)
	addressEvaluations := 0
	for _, contact := range contacts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		addresses, err := s.contacts.ListContactAddresses(ctx, tenantID, contact.ID)
		if err != nil {
			return nil, safeContactError(err)
		}
		addressEvaluations += len(addresses)
		if addressEvaluations > AudienceHardLimit {
			return nil, ErrAudienceEvaluationLimit
		}
		if segmentMatches(contact, addresses, segment.Definition) {
			ids = append(ids, contact.ID)
		}
	}
	return ids, nil
}

func (s *AudienceService) buildRecipients(ctx context.Context, tenantID string, contactIDs []uint64, selection AudienceSelection) ([]AudienceRecipient, int, int, error) {
	var exclusion *ExclusionList
	var err error
	if selection.ExclusionListPublicID != "" {
		exclusion, err = s.GetExclusionList(ctx, tenantID, selection.ExclusionListPublicID)
		if err != nil {
			return nil, 0, 0, err
		}
	}
	excludedContacts := 0
	excludedRecipients := 0
	contactExcluded := make(map[uint64]struct{})
	addressExcluded := make(map[string]struct{})
	if exclusion != nil {
		for _, id := range exclusion.ContactIDs {
			contactExcluded[id] = struct{}{}
		}
		for _, identity := range exclusion.AddressIdentities {
			addressExcluded[identity.Key()] = struct{}{}
		}
	}
	recipients := make([]AudienceRecipient, 0)
	addressEvaluations := 0
	for _, contactID := range contactIDs {
		if err := ctx.Err(); err != nil {
			return nil, 0, 0, err
		}
		if _, excluded := contactExcluded[contactID]; excluded {
			excludedContacts++
			addresses, addressErr := s.contacts.ListContactAddresses(ctx, tenantID, contactID)
			if addressErr != nil && !errors.Is(addressErr, ErrNotFound) {
				return nil, 0, 0, safeContactError(addressErr)
			}
			excludedRecipients += len(addresses)
			if excludedRecipients > AudienceHardLimit {
				return nil, 0, 0, ErrAudienceEvaluationLimit
			}
			continue
		}
		addresses, err := s.contacts.ListContactAddresses(ctx, tenantID, contactID)
		if err != nil {
			return nil, 0, 0, safeContactError(err)
		}
		addressEvaluations += len(addresses)
		if addressEvaluations > AudienceHardLimit {
			return nil, 0, 0, ErrAudienceEvaluationLimit
		}
		for _, address := range addresses {
			if _, excluded := addressExcluded[address.Identity.Key()]; excluded {
				excludedRecipients++
				continue
			}
			recipients = append(recipients, AudienceRecipient{ContactID: contactID, AddressID: address.ID, Identity: address.Identity})
			if len(recipients) > effectiveRecipientLimit(selection.MaxRecipients) {
				return nil, 0, 0, ErrAudienceEvaluationLimit
			}
		}
	}
	return recipients, excludedContacts, excludedRecipients, nil
}

func effectiveRecipientLimit(limit int) int {
	if limit <= 0 || limit > AudienceHardLimit {
		return AudienceHardLimit
	}
	return limit
}

func allowedSegmentField(field string) bool {
	if field == SegmentFieldStatus || field == SegmentFieldSource || field == SegmentFieldLocale || field == SegmentFieldTimezone || field == SegmentFieldOwnerID || field == SegmentFieldLifecycleStage || field == SegmentFieldAddressKind || field == SegmentFieldAddressNamespace || field == SegmentFieldAddressProvider || field == SegmentFieldAddressConsent {
		return true
	}
	return strings.HasPrefix(field, "contact.custom.") && len(field) > len("contact.custom.") && len(field) <= 128
}

func validPredicateOperator(operator PredicateOperator) bool {
	return operator == PredicateEquals || operator == PredicateNotEquals || operator == PredicateContains || operator == PredicateStartsWith || operator == PredicateExists || operator == PredicateIn
}

func segmentMatches(contact *Contact, addresses []*ContactAddress, definition SegmentDefinition) bool {
	fields, tags := definition.predicates()
	matches := make([]bool, 0, len(fields)+len(tags))
	for _, predicate := range fields {
		matches = append(matches, fieldPredicateMatches(contact, addresses, predicate))
	}
	for _, predicate := range tags {
		matches = append(matches, tagPredicateMatches(contact, predicate))
	}
	if len(matches) == 0 {
		return false
	}
	if definition.Match == SegmentMatchAny {
		for _, match := range matches {
			if match {
				return true
			}
		}
		return false
	}
	for _, match := range matches {
		if !match {
			return false
		}
	}
	return true
}

func fieldPredicateMatches(contact *Contact, addresses []*ContactAddress, predicate FieldPredicate) bool {
	if strings.HasPrefix(predicate.Field, "address.") {
		for _, address := range addresses {
			if value, present := addressFieldValue(address, predicate.Field); predicateMatches(value, present, predicate) {
				return true
			}
		}
		return false
	}
	value, present := contactFieldValue(contact, predicate.Field)
	return predicateMatches(value, present, predicate)
}

func predicateMatches(value string, present bool, predicate FieldPredicate) bool {
	if predicate.Operator == PredicateExists {
		return present == (strings.EqualFold(strings.TrimSpace(predicate.Value), "true") || predicate.Value == "")
	}
	if !present {
		return predicate.Operator == PredicateNotEquals
	}
	value = strings.ToLower(value)
	if predicate.Operator == PredicateIn {
		for _, candidate := range predicate.Values {
			if value == strings.ToLower(candidate) {
				return true
			}
		}
		return false
	}
	want := strings.ToLower(predicate.Value)
	switch predicate.Operator {
	case PredicateEquals:
		return value == want
	case PredicateNotEquals:
		return value != want
	case PredicateContains:
		return strings.Contains(value, want)
	case PredicateStartsWith:
		return strings.HasPrefix(value, want)
	default:
		return false
	}
}

func contactFieldValue(contact *Contact, field string) (string, bool) {
	switch field {
	case SegmentFieldStatus:
		return contact.Status.String(), true
	case SegmentFieldSource:
		return contact.Source.String(), true
	case SegmentFieldLocale:
		return contact.Locale, contact.Locale != ""
	case SegmentFieldTimezone:
		return contact.Timezone, contact.Timezone != ""
	case SegmentFieldOwnerID:
		return contact.OwnerID, contact.OwnerID != ""
	case SegmentFieldLifecycleStage:
		return contact.LifecycleStage, contact.LifecycleStage != ""
	}
	if strings.HasPrefix(field, "contact.custom.") {
		key := strings.TrimPrefix(field, "contact.custom.")
		value, ok := contact.CustomFields[key]
		return value, ok
	}
	return "", false
}

func addressFieldValue(address *ContactAddress, field string) (string, bool) {
	switch field {
	case SegmentFieldAddressKind:
		return address.Kind.String(), true
	case SegmentFieldAddressNamespace:
		return address.Namespace, address.Namespace != ""
	case SegmentFieldAddressProvider:
		return address.Provider, address.Provider != ""
	case SegmentFieldAddressConsent:
		return string(address.Consent.State), address.Consent.State != ConsentStateUnknown
	default:
		return "", false
	}
}

func tagPredicateMatches(contact *Contact, predicate TagPredicate) bool {
	want := strings.ToLower(strings.TrimSpace(predicate.Tag))
	found := false
	for _, tag := range contact.Tags {
		if strings.ToLower(strings.TrimSpace(tag)) == want {
			found = true
			break
		}
	}
	if predicate.Operator == TagNotHas {
		return !found
	}
	return found
}

func uniqueIDs(values []uint64) ([]uint64, bool) {
	result := make([]uint64, 0, len(values))
	seen := make(map[uint64]struct{}, len(values))
	changed := false
	for _, value := range values {
		if value == 0 {
			changed = true
			continue
		}
		if _, ok := seen[value]; ok {
			changed = true
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result, changed
}

func cloneStaticList(value *StaticList) *StaticList {
	if value == nil {
		return nil
	}
	copy := *value
	copy.ContactIDs = append([]uint64(nil), value.ContactIDs...)
	return &copy
}

func cloneDynamicSegment(value *DynamicSegment) *DynamicSegment {
	if value == nil {
		return nil
	}
	copy := *value
	copy.Definition.Fields = append([]FieldPredicate(nil), value.Definition.Fields...)
	copy.Definition.Tags = append([]TagPredicate(nil), value.Definition.Tags...)
	copy.Definition.FieldPredicates = append([]FieldPredicate(nil), value.Definition.FieldPredicates...)
	copy.Definition.TagPredicates = append([]TagPredicate(nil), value.Definition.TagPredicates...)
	for index := range copy.Definition.Fields {
		copy.Definition.Fields[index].Values = append([]string(nil), value.Definition.Fields[index].Values...)
	}
	for index := range copy.Definition.FieldPredicates {
		copy.Definition.FieldPredicates[index].Values = append([]string(nil), value.Definition.FieldPredicates[index].Values...)
	}
	return &copy
}

func cloneExclusionList(value *ExclusionList) *ExclusionList {
	if value == nil {
		return nil
	}
	copy := *value
	copy.ContactIDs = append([]uint64(nil), value.ContactIDs...)
	copy.AddressIdentities = append([]AddressIdentity(nil), value.AddressIdentities...)
	return &copy
}

func cloneAudienceSnapshot(value *AudienceSnapshot) *AudienceSnapshot {
	if value == nil {
		return nil
	}
	copy := *value
	copy.Selection = value.Selection
	copy.Recipients = append([]AudienceRecipient(nil), value.Recipients...)
	return &copy
}

// The following methods are the concurrency-safe in-memory audience adapter.
func (r *InMemoryContactRepository) CreateStaticList(ctx context.Context, value *StaticList) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("%w: static list is nil", ErrInvalid)
	}
	candidate := cloneStaticList(value)
	if candidate.PublicID == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.UpdatedAt = time.Now().UTC()
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	key := scopedKey(candidate.TenantID, candidate.PublicID)
	if _, exists := r.staticListByKey[key]; exists {
		return ErrConflict
	}
	candidate.ID = r.nextStaticListID
	r.nextStaticListID++
	r.staticLists[candidate.ID] = candidate
	r.staticListByKey[key] = candidate.ID
	*value = *cloneStaticList(candidate)
	return nil
}

func (r *InMemoryContactRepository) GetStaticList(ctx context.Context, tenantID, publicID string) (*StaticList, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.staticListByKey[scopedKey(tenantID, publicID)]
	if !ok || r.staticLists[id] == nil {
		return nil, ErrNotFound
	}
	return cloneStaticList(r.staticLists[id]), nil
}

func (r *InMemoryContactRepository) UpdateStaticList(ctx context.Context, value *StaticList) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return ErrInvalid
	}
	candidate := cloneStaticList(value)
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.staticLists[candidate.ID]
	if !ok || stored.TenantID != candidate.TenantID || stored.PublicID != candidate.PublicID {
		return ErrNotFound
	}
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	r.staticLists[candidate.ID] = candidate
	*value = *cloneStaticList(candidate)
	return nil
}

func (r *InMemoryContactRepository) CreateDynamicSegment(ctx context.Context, value *DynamicSegment) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return ErrInvalid
	}
	candidate := cloneDynamicSegment(value)
	if candidate.PublicID == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.UpdatedAt = time.Now().UTC()
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	key := scopedKey(candidate.TenantID, candidate.PublicID)
	if _, exists := r.segmentByKey[key]; exists {
		return ErrConflict
	}
	candidate.ID = r.nextSegmentID
	r.nextSegmentID++
	r.segments[candidate.ID] = candidate
	r.segmentByKey[key] = candidate.ID
	*value = *cloneDynamicSegment(candidate)
	return nil
}

func (r *InMemoryContactRepository) GetDynamicSegment(ctx context.Context, tenantID, publicID string) (*DynamicSegment, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.segmentByKey[scopedKey(tenantID, publicID)]
	if !ok || r.segments[id] == nil {
		return nil, ErrNotFound
	}
	return cloneDynamicSegment(r.segments[id]), nil
}

func (r *InMemoryContactRepository) UpdateDynamicSegment(ctx context.Context, value *DynamicSegment) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return ErrInvalid
	}
	candidate := cloneDynamicSegment(value)
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.segments[candidate.ID]
	if !ok || stored.TenantID != candidate.TenantID || stored.PublicID != candidate.PublicID {
		return ErrNotFound
	}
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	r.segments[candidate.ID] = candidate
	*value = *cloneDynamicSegment(candidate)
	return nil
}

func (r *InMemoryContactRepository) CreateExclusionList(ctx context.Context, value *ExclusionList) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return ErrInvalid
	}
	candidate := cloneExclusionList(value)
	if candidate.PublicID == "" {
		candidate.PublicID = hash.NanoID(21)
	}
	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now().UTC()
	}
	candidate.UpdatedAt = time.Now().UTC()
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensureInitializedLocked()
	key := scopedKey(candidate.TenantID, candidate.PublicID)
	if _, exists := r.exclusionByKey[key]; exists {
		return ErrConflict
	}
	candidate.ID = r.nextExclusionID
	r.nextExclusionID++
	r.exclusionLists[candidate.ID] = candidate
	r.exclusionByKey[key] = candidate.ID
	*value = *cloneExclusionList(candidate)
	return nil
}

func (r *InMemoryContactRepository) GetExclusionList(ctx context.Context, tenantID, publicID string) (*ExclusionList, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.exclusionByKey[scopedKey(tenantID, publicID)]
	if !ok || r.exclusionLists[id] == nil {
		return nil, ErrNotFound
	}
	return cloneExclusionList(r.exclusionLists[id]), nil
}

func (r *InMemoryContactRepository) UpdateExclusionList(ctx context.Context, value *ExclusionList) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return ErrInvalid
	}
	candidate := cloneExclusionList(value)
	if err := candidate.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.exclusionLists[candidate.ID]
	if !ok || stored.TenantID != candidate.TenantID || stored.PublicID != candidate.PublicID {
		return ErrNotFound
	}
	candidate.CreatedAt = stored.CreatedAt
	candidate.UpdatedAt = time.Now().UTC()
	r.exclusionLists[candidate.ID] = candidate
	*value = *cloneExclusionList(candidate)
	return nil
}

func (r *InMemoryContactRepository) CreateAudienceSnapshot(ctx context.Context, value *AudienceSnapshot) error {
	if err := repositoryContextError(ctx); err != nil {
		return err
	}
	if value == nil {
		return ErrInvalid
	}
	candidate := cloneAudienceSnapshot(value)
	if candidate.PublicID == "" {
		candidate.PublicID = hash.NanoID(21)
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
	key := scopedKey(candidate.TenantID, candidate.PublicID)
	if _, exists := r.snapshotByKey[key]; exists {
		return ErrConflict
	}
	candidate.ID = r.nextSnapshotID
	r.nextSnapshotID++
	r.snapshots[candidate.ID] = candidate
	r.snapshotByKey[key] = candidate.ID
	*value = *cloneAudienceSnapshot(candidate)
	return nil
}

func (r *InMemoryContactRepository) GetAudienceSnapshot(ctx context.Context, tenantID, publicID string) (*AudienceSnapshot, error) {
	if err := repositoryContextError(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.snapshotByKey[scopedKey(tenantID, publicID)]
	if !ok || r.snapshots[id] == nil {
		return nil, ErrNotFound
	}
	return cloneAudienceSnapshot(r.snapshots[id]), nil
}

var _ AudienceRepository = (*InMemoryContactRepository)(nil)
