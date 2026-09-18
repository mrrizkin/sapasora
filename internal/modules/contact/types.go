package contact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ContactStatus is the lifecycle state of a contact. Consent is deliberately
// modeled separately on each address and is not inferred from this status.
type ContactStatus string

const (
	ContactStatusUnknown  ContactStatus = "unknown"
	ContactStatusActive   ContactStatus = "active"
	ContactStatusInactive ContactStatus = "inactive"
	ContactStatusBlocked  ContactStatus = "blocked"
	ContactStatusArchived ContactStatus = "archived"
)

func (s ContactStatus) Valid() bool {
	switch s {
	case ContactStatusActive, ContactStatusInactive, ContactStatusBlocked, ContactStatusArchived:
		return true
	default:
		return false
	}
}

func (s ContactStatus) String() string { return string(s) }

// ContactSource records the origin of a contact or address without coupling
// the domain to a provider implementation.
type ContactSource string

const (
	ContactSourceUnknown      ContactSource = "unknown"
	ContactSourceManual       ContactSource = "manual"
	ContactSourceAPI          ContactSource = "api"
	ContactSourceImport       ContactSource = "import"
	ContactSourceProviderSync ContactSource = "provider_sync"
	ContactSourceInbound      ContactSource = "inbound"
	ContactSourceAutomation   ContactSource = "automation"
)

func (s ContactSource) Valid() bool {
	switch s {
	case ContactSourceManual, ContactSourceAPI, ContactSourceImport, ContactSourceProviderSync, ContactSourceInbound, ContactSourceAutomation:
		return true
	default:
		return false
	}
}

func (s ContactSource) String() string { return string(s) }

// AddressKind identifies the canonical identity algorithm used by an address.
type AddressKind string

const (
	AddressKindUnknown  AddressKind = "unknown"
	AddressKindPhone    AddressKind = "phone"
	AddressKindEmail    AddressKind = "email"
	AddressKindUsername AddressKind = "username"
)

func (k AddressKind) Valid() bool {
	switch k {
	case AddressKindPhone, AddressKindEmail, AddressKindUsername:
		return true
	default:
		return false
	}
}

func (k AddressKind) String() string { return string(k) }

// ConsentState is the current consent state for one contact address.
type ConsentState string

const (
	ConsentStateUnknown  ConsentState = "unknown"
	ConsentStateOptedIn  ConsentState = "opted_in"
	ConsentStateOptedOut ConsentState = "opted_out"
)

func (s ConsentState) Valid() bool {
	switch s {
	case ConsentStateUnknown, ConsentStateOptedIn, ConsentStateOptedOut:
		return true
	default:
		return false
	}
}

// ConsentSource describes how consent evidence was obtained.
type ConsentSource string

const (
	ConsentSourceUnknown      ConsentSource = "unknown"
	ConsentSourceManual       ConsentSource = "manual"
	ConsentSourceAPI          ConsentSource = "api"
	ConsentSourceImport       ConsentSource = "import"
	ConsentSourceInbound      ConsentSource = "inbound"
	ConsentSourceProviderSync ConsentSource = "provider_sync"
)

func (s ConsentSource) Valid() bool {
	switch s {
	case ConsentSourceUnknown, ConsentSourceManual, ConsentSourceAPI, ConsentSourceImport, ConsentSourceInbound, ConsentSourceProviderSync:
		return true
	default:
		return false
	}
}

// ConsentMetadata is address-level metadata. EvidenceRef and ActorID are
// opaque references; raw proof or message content does not belong here.
type ConsentMetadata struct {
	State       ConsentState  `json:"state"`
	Source      ConsentSource `json:"source"`
	OccurredAt  *time.Time    `json:"occurred_at,omitempty"`
	EvidenceRef string        `json:"-"`
	ActorID     string        `json:"-"`
}

func (m ConsentMetadata) Valid() error {
	if !m.State.Valid() {
		return fmt.Errorf("invalid consent state")
	}
	if !m.Source.Valid() {
		return fmt.Errorf("invalid consent source")
	}
	return nil
}

func (m ConsentMetadata) String() string {
	return fmt.Sprintf("ConsentMetadata{state=%s source=%s occurred_at=%t}", m.State, m.Source, m.OccurredAt != nil)
}

func (m ConsentMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		State      ConsentState  `json:"state"`
		Source     ConsentSource `json:"source"`
		OccurredAt *time.Time    `json:"occurred_at,omitempty"`
	}{m.State, m.Source, m.OccurredAt})
}

// AddressIdentity is the normalized, unique identity of an address. Value is
// intentionally excluded from JSON and String diagnostics because it is PII.
// Namespace is usually "global"; provider-specific usernames can use an
// explicit namespace to avoid accidental cross-provider collisions.
type AddressIdentity struct {
	Kind      AddressKind `json:"kind"`
	Namespace string      `json:"namespace"`
	Value     string      `json:"-"`
}

func (i AddressIdentity) Valid() error {
	if !i.Kind.Valid() {
		return errors.New("invalid address kind")
	}
	if strings.TrimSpace(i.Namespace) == "" {
		return errors.New("address namespace is required")
	}
	if strings.TrimSpace(i.Value) == "" {
		return errors.New("normalized address is required")
	}
	canonical, err := NormalizeAddress(i.Kind, i.Namespace, i.Value)
	if err != nil || canonical != i {
		return errors.New("address identity is not canonical")
	}
	return nil
}

// Key is a stable, non-PII fingerprint suitable for a repository uniqueness
// index. Length-prefixing avoids delimiter collisions before hashing.
func (i AddressIdentity) Key() string {
	payload := fmt.Sprintf("%d:%s%d:%s%d:%s", len(i.Kind), i.Kind, len(i.Namespace), i.Namespace, len(i.Value), i.Value)
	digest := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(digest[:])
}

func (i AddressIdentity) String() string {
	return fmt.Sprintf("AddressIdentity{kind=%s namespace=%s value=[REDACTED]}", i.Kind, i.Namespace)
}

// MarshalJSON emits only a safe diagnostic projection. Callers that need the
// identity value must use an explicit domain operation, never diagnostics.
// ContactAddressIdentity is the descriptive name used by integrations that
// prefer the aggregate-qualified type name.
type ContactAddressIdentity = AddressIdentity

func (i AddressIdentity) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind      AddressKind `json:"kind"`
		Namespace string      `json:"namespace"`
		Present   bool        `json:"value_present"`
	}{i.Kind, i.Namespace, i.Value != ""})
}

// Contact is the tenant-owned profile. Fields such as DisplayName, Notes, and
// CustomFields are retained for domain use but omitted from diagnostic JSON.
type Contact struct {
	ID             uint64            `json:"id"`
	PublicID       string            `json:"public_id"`
	TenantID       string            `json:"tenant_id"`
	DisplayName    string            `json:"-"`
	Status         ContactStatus     `json:"status"`
	Source         ContactSource     `json:"source"`
	SourceMetadata map[string]string `json:"-"`
	Locale         string            `json:"-"`
	Timezone       string            `json:"-"`
	OwnerID        string            `json:"-"`
	LifecycleStage string            `json:"-"`
	CustomFields   map[string]string `json:"-"`
	Tags           []string          `json:"-"`
	Notes          string            `json:"-"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      *time.Time        `json:"-"`
}

func (c Contact) Valid() error {
	if strings.TrimSpace(c.TenantID) == "" {
		return fmt.Errorf("tenant id is required")
	}
	if !c.Status.Valid() {
		return fmt.Errorf("invalid contact status")
	}
	if !c.Source.Valid() {
		return fmt.Errorf("invalid contact source")
	}
	if c.DeletedAt != nil && c.DeletedAt.IsZero() {
		return fmt.Errorf("deleted timestamp is invalid")
	}
	return nil
}

func (c Contact) String() string {
	return fmt.Sprintf("Contact{public_id=%s status=%s source=%s}", c.PublicID, c.Status, c.Source)
}

func (c Contact) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        uint64        `json:"id"`
		PublicID  string        `json:"public_id"`
		TenantID  string        `json:"tenant_id"`
		Status    ContactStatus `json:"status"`
		Source    ContactSource `json:"source"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
	}{c.ID, c.PublicID, c.TenantID, c.Status, c.Source, c.CreatedAt, c.UpdatedAt})
}

// ContactAddress links a contact to one normalized identity. Value is the
// caller-supplied form and NormalizedValue is canonical; both are PII and are
// excluded from String/JSON diagnostics. ProviderAddressID is opaque provider
// metadata and is also excluded because it may itself be an address.
type ContactAddress struct {
	ID                uint64          `json:"id"`
	PublicID          string          `json:"public_id"`
	TenantID          string          `json:"tenant_id"`
	ContactID         uint64          `json:"contact_id"`
	Identity          AddressIdentity `json:"-"`
	Kind              AddressKind     `json:"kind"`
	Namespace         string          `json:"namespace"`
	Value             string          `json:"-"`
	NormalizedValue   string          `json:"-"`
	Provider          string          `json:"provider,omitempty"`
	ProviderAddressID string          `json:"-"`
	Source            ContactSource   `json:"source"`
	Consent           ConsentMetadata `json:"consent"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *time.Time      `json:"-"`
}

func (a ContactAddress) Valid() error {
	if strings.TrimSpace(a.TenantID) == "" {
		return errors.New("tenant id is required")
	}
	if a.ContactID == 0 {
		return errors.New("contact id is required")
	}
	if err := a.Identity.Valid(); err != nil {
		return err
	}
	if a.Kind != a.Identity.Kind || a.Namespace != a.Identity.Namespace || a.NormalizedValue != a.Identity.Value {
		return errors.New("address identity fields are inconsistent")
	}
	if !a.Source.Valid() {
		return errors.New("invalid contact source")
	}
	if a.DeletedAt != nil && a.DeletedAt.IsZero() {
		return errors.New("deleted timestamp is invalid")
	}
	return a.Consent.Valid()
}

func (a ContactAddress) String() string {
	return fmt.Sprintf("ContactAddress{public_id=%s kind=%s namespace=%s source=%s consent=%s}", a.PublicID, a.Kind, a.Namespace, a.Source, a.Consent.State)
}

func (a ContactAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID              uint64          `json:"id"`
		PublicID        string          `json:"public_id"`
		TenantID        string          `json:"tenant_id"`
		ContactID       uint64          `json:"contact_id"`
		Kind            AddressKind     `json:"kind"`
		Namespace       string          `json:"namespace"`
		Provider        string          `json:"provider,omitempty"`
		Source          ContactSource   `json:"source"`
		Consent         ConsentMetadata `json:"consent"`
		IdentityPresent bool            `json:"identity_present"`
		CreatedAt       time.Time       `json:"created_at"`
		UpdatedAt       time.Time       `json:"updated_at"`
	}{a.ID, a.PublicID, a.TenantID, a.ContactID, a.Kind, a.Namespace, a.Provider, a.Source, a.Consent, a.Identity.Value != "", a.CreatedAt, a.UpdatedAt})
}

// NewContactAddress normalizes the supplied identity before it reaches a
// repository. Repository writes repeat the normalization defensively.
func NewContactAddress(tenantID string, contactID uint64, kind AddressKind, namespace, value string) (ContactAddress, error) {
	identity, err := NormalizeAddress(kind, namespace, value)
	if err != nil {
		return ContactAddress{}, err
	}
	address := ContactAddress{
		TenantID:        tenantID,
		ContactID:       contactID,
		Identity:        identity,
		Kind:            identity.Kind,
		Namespace:       identity.Namespace,
		Value:           value,
		NormalizedValue: identity.Value,
		Source:          ContactSourceManual,
		Consent: ConsentMetadata{
			State:  ConsentStateUnknown,
			Source: ConsentSourceUnknown,
		},
	}
	return address, nil
}
