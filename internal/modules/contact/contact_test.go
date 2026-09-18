package contact

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestNormalizeAddresses(t *testing.T) {
	tests := []struct {
		name string
		kind AddressKind
		in   string
		want string
	}{
		{name: "phone formatting", kind: AddressKindPhone, in: "+1 (555) 123-4567", want: "+15551234567"},
		{name: "phone international prefix", kind: AddressKindPhone, in: "0044 20 7946 0958", want: "+442079460958"},
		{name: "email case and whitespace", kind: AddressKindEmail, in: "  Alice@Example.COM ", want: "alice@example.com"},
		{name: "username at prefix", kind: AddressKindUsername, in: " @Support.Team ", want: "support.team"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NormalizeAddress(test.kind, "", test.in)
			if err != nil {
				t.Fatalf("NormalizeAddress() error = %v", err)
			}
			if got.Value != test.want || got.Namespace != "global" {
				t.Fatalf("NormalizeAddress() = %+v, want value %q in global namespace", got, test.want)
			}
		})
	}

	for _, test := range []struct {
		name string
		kind AddressKind
		in   string
	}{
		{name: "local phone cannot be interpreted", kind: AddressKindPhone, in: "555-1234"},
		{name: "invalid email", kind: AddressKindEmail, in: "not-an-email"},
		{name: "invalid username punctuation", kind: AddressKindUsername, in: "alice/example"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NormalizeAddress(test.kind, "", test.in); err == nil {
				t.Fatal("NormalizeAddress() error = nil")
			} else if strings.Contains(err.Error(), test.in) {
				t.Fatalf("normalization error leaked input: %v", err)
			}
		})
	}
}

func TestContactDiagnosticsDoNotExposePII(t *testing.T) {
	const (
		name       = "Private Example"
		phone      = "+1 (555) 123-4567"
		normalized = "+15551234567"
		evidence   = "proof-private-contact"
	)
	contact := Contact{
		PublicID:    "contact-public",
		TenantID:    "tenant-a",
		DisplayName: name,
		Status:      ContactStatusActive,
		Source:      ContactSourceManual,
		Notes:       "private note",
	}
	address, err := NewContactAddress("tenant-a", 1, AddressKindPhone, "", phone)
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	address.PublicID = "address-public"
	address.ProviderAddressID = "provider-private-id"
	address.Consent.EvidenceRef = evidence

	for label, value := range map[string]any{
		"contact string":  contact.String(),
		"address string":  address.String(),
		"identity string": address.Identity.String(),
	} {
		text := fmt.Sprint(value)
		for _, secret := range []string{name, phone, normalized, evidence, "provider-private-id"} {
			if strings.Contains(text, secret) {
				t.Fatalf("%s leaked %q: %s", label, secret, text)
			}
		}
	}
	consentDiagnostics := fmt.Sprint(address.Consent)
	if strings.Contains(consentDiagnostics, evidence) {
		t.Fatalf("consent String() leaked evidence: %s", consentDiagnostics)
	}
	event := ConsentEvent{
		TenantID:  "tenant-a",
		AddressID: 1,
		Sequence:  1,
		State:     ConsentStateOptedIn,
		Metadata:  ConsentMetadata{State: ConsentStateOptedIn, Source: ConsentSourceAPI, EvidenceRef: evidence, ActorID: "actor-private"},
	}
	suppression := SuppressionRecord{
		TenantID:    "tenant-a",
		Identity:    address.Identity,
		Reason:      SuppressionReasonBlocked,
		Source:      ConsentSourceManual,
		EvidenceRef: evidence,
		ActorID:     "actor-private",
	}
	for label, value := range map[string]string{
		"consent event string": event.String(),
		"suppression string":   suppression.String(),
	} {
		if strings.Contains(value, evidence) || strings.Contains(value, normalized) {
			t.Fatalf("%s leaked private data: %s", label, value)
		}
	}
	encodedEvent, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("json.Marshal(event) error = %v", err)
	}
	encodedSuppression, err := json.Marshal(suppression)
	if err != nil {
		t.Fatalf("json.Marshal(suppression) error = %v", err)
	}
	if strings.Contains(string(encodedEvent)+string(encodedSuppression), evidence) || strings.Contains(string(encodedSuppression), normalized) {
		t.Fatalf("consent/suppression JSON leaked private data: %s %s", encodedEvent, encodedSuppression)
	}
	encodedConsent, err := json.Marshal(address.Consent)
	if err != nil {
		t.Fatalf("json.Marshal(consent) error = %v", err)
	}
	if strings.Contains(string(encodedConsent), evidence) {
		t.Fatalf("consent JSON leaked evidence: %s", encodedConsent)
	}

	encodedContact, err := json.Marshal(contact)
	if err != nil {
		t.Fatalf("json.Marshal(contact) error = %v", err)
	}
	encodedAddress, err := json.Marshal(address)
	if err != nil {
		t.Fatalf("json.Marshal(address) error = %v", err)
	}
	jsonText := string(encodedContact) + string(encodedAddress)
	for _, secret := range []string{name, phone, normalized, evidence, "provider-private-id"} {
		if strings.Contains(jsonText, secret) {
			t.Fatalf("JSON diagnostics leaked %q: %s", secret, jsonText)
		}
	}
}

func newTestContact(t *testing.T, repository ContactRepository, tenant, publicID string) Contact {
	t.Helper()
	contact := Contact{TenantID: tenant, PublicID: publicID, Status: ContactStatusActive, Source: ContactSourceManual}
	if err := repository.CreateContact(context.Background(), &contact); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}
	return contact
}

func TestInMemoryContactRepositoryIsTenantScopedAndUniqueByNormalizedIdentity(t *testing.T) {
	repository := NewInMemoryContactRepository()
	ctx := context.Background()
	contactA := newTestContact(t, repository, "tenant-a", "contact-a")
	contactB := newTestContact(t, repository, "tenant-b", "contact-b")

	address, err := NewContactAddress("tenant-a", contactA.ID, AddressKindEmail, "", "Alice@Example.com")
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	address.Consent = ConsentMetadata{State: ConsentStateOptedIn, Source: ConsentSourceAPI}
	if err := repository.CreateContactAddress(ctx, &address); err != nil {
		t.Fatalf("CreateContactAddress() error = %v", err)
	}
	if address.ID == 0 || address.NormalizedValue != "alice@example.com" {
		t.Fatalf("CreateContactAddress() did not assign normalized persistence fields: %+v", address)
	}

	lookup, err := repository.FindContactAddressByIdentity(ctx, "tenant-a", AddressIdentity{Kind: AddressKindEmail, Namespace: "global", Value: "ALICE@EXAMPLE.COM"})
	if err != nil || lookup.ID != address.ID {
		t.Fatalf("FindContactAddressByIdentity() = %+v, %v", lookup, err)
	}
	if _, err := repository.FindContactAddressByIdentity(ctx, "tenant-b", address.Identity); !errors.Is(err, ErrContactAddressNotFound) {
		t.Fatalf("cross-tenant identity lookup error = %v", err)
	}

	duplicate := address
	duplicate.ID = 0
	duplicate.PublicID = "address-duplicate"
	if err := repository.CreateContactAddress(ctx, &duplicate); !errors.Is(err, ErrContactAddressConflict) {
		t.Fatalf("duplicate identity error = %v, want ErrContactAddressConflict", err)
	}

	otherTenant, err := NewContactAddress("tenant-b", contactB.ID, AddressKindEmail, "", "alice@example.com")
	if err != nil {
		t.Fatalf("NewContactAddress(other tenant) error = %v", err)
	}
	if err := repository.CreateContactAddress(ctx, &otherTenant); err != nil {
		t.Fatalf("same identity in another tenant error = %v", err)
	}
	if _, err := repository.GetContactByPublicID(ctx, "tenant-b", contactA.PublicID); !errors.Is(err, ErrContactNotFound) {
		t.Fatalf("cross-tenant contact lookup error = %v", err)
	}
}

func TestInMemoryContactRepositoryDefensiveCopiesAndIdentityUpdate(t *testing.T) {
	repository := NewInMemoryContactRepository()
	contact := newTestContact(t, repository, "tenant-a", "contact-a")
	contact.DisplayName = "Alice"
	contact.Tags = []string{"vip"}
	contact.CustomFields = map[string]string{"note": "private"}
	if err := repository.UpdateContact(context.Background(), &contact); err != nil {
		t.Fatalf("UpdateContact() error = %v", err)
	}
	contact.Tags[0] = "mutated"
	contact.CustomFields["note"] = "mutated"
	storedContact, err := repository.GetContactByPublicID(context.Background(), "tenant-a", "contact-a")
	if err != nil || storedContact.Tags[0] != "vip" || storedContact.CustomFields["note"] != "private" {
		t.Fatalf("contact defensive copy failed: %+v, %v", storedContact, err)
	}

	address, err := NewContactAddress("tenant-a", contact.ID, AddressKindUsername, "chat", "Alice")
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	if err := repository.CreateContactAddress(context.Background(), &address); err != nil {
		t.Fatalf("CreateContactAddress() error = %v", err)
	}
	address.Value = "changed caller value"
	storedAddress, err := repository.GetContactAddressByPublicID(context.Background(), "tenant-a", address.PublicID)
	if err != nil || storedAddress.Value != "Alice" {
		t.Fatalf("address defensive copy failed: %+v, %v", storedAddress, err)
	}

	storedAddress.Value = "Bob"
	if err := repository.UpdateContactAddress(context.Background(), storedAddress); err != nil {
		t.Fatalf("UpdateContactAddress() error = %v", err)
	}
	storedAddress.ContactID = contact.ID + 1
	if err := repository.UpdateContactAddress(context.Background(), storedAddress); !errors.Is(err, ErrContactAddressConflict) {
		t.Fatalf("cross-contact address update error = %v, want ErrContactAddressConflict", err)
	}
	storedAddress.ContactID = contact.ID
	if _, err := repository.FindContactAddressByIdentity(context.Background(), "tenant-a", address.Identity); !errors.Is(err, ErrContactAddressNotFound) {
		t.Fatalf("old identity remained indexed: %v", err)
	}
	newIdentity, err := NormalizeAddress(AddressKindUsername, "chat", "Bob")
	if err != nil {
		t.Fatalf("NormalizeAddress() error = %v", err)
	}
	if _, err := repository.FindContactAddressByIdentity(context.Background(), "tenant-a", newIdentity); err != nil {
		t.Fatalf("new identity lookup error = %v", err)
	}
}

func TestInMemoryContactRepositoryConcurrentUniqueCreate(t *testing.T) {
	repository := NewInMemoryContactRepository()
	contact := newTestContact(t, repository, "tenant-a", "contact-a")
	const workers = 32
	results := make(chan error, workers)
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			address, err := NewContactAddress("tenant-a", contact.ID, AddressKindPhone, "", "+1 555 123 4567")
			if err == nil {
				err = repository.CreateContactAddress(context.Background(), &address)
			}
			results <- err
		}()
	}
	wait.Wait()
	close(results)

	created := 0
	conflicts := 0
	for err := range results {
		switch {
		case err == nil:
			created++
		case errors.Is(err, ErrContactAddressConflict):
			conflicts++
		default:
			t.Fatalf("concurrent create error = %v", err)
		}
	}
	if created != 1 || conflicts != workers-1 {
		t.Fatalf("concurrent creates: created=%d conflicts=%d", created, conflicts)
	}
	addresses, err := repository.ListContactAddresses(context.Background(), "tenant-a", contact.ID)
	if err != nil || len(addresses) != 1 {
		t.Fatalf("ListContactAddresses() = %d, %v", len(addresses), err)
	}
}
