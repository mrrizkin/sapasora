package contact

import (
	"context"
	"errors"
	"testing"
)

func TestContactServiceCRUDIsTenantScoped(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)

	contact := Contact{PublicID: "contact-a", DisplayName: "Alice"}
	if err := service.CreateContact(ctx, "tenant-a", &contact); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}
	if contact.TenantID != "tenant-a" || contact.Status != ContactStatusActive || contact.Source != ContactSourceManual {
		t.Fatalf("CreateContact() defaults = %+v", contact)
	}
	other := Contact{PublicID: "contact-b"}
	if err := service.CreateContact(ctx, "tenant-b", &other); err != nil {
		t.Fatalf("CreateContact(other) error = %v", err)
	}

	got, err := service.GetContact(ctx, "tenant-a", contact.PublicID)
	if err != nil || got.ID != contact.ID {
		t.Fatalf("GetContact() = %+v, %v", got, err)
	}
	if got, err := service.GetContactByID(ctx, "tenant-a", contact.ID); err != nil || got.PublicID != contact.PublicID {
		t.Fatalf("GetContactByID() = %+v, %v", got, err)
	}
	duplicateContact := Contact{PublicID: contact.PublicID}
	if err := service.CreateContact(ctx, "tenant-a", &duplicateContact); !errors.Is(err, ErrConflict) {
		t.Fatalf("duplicate contact error = %v, want ErrConflict", err)
	}
	if _, err := service.GetContact(ctx, "tenant-b", contact.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant GetContact() error = %v, want ErrNotFound", err)
	}

	contact.DisplayName = "Alice Updated"
	if err := service.UpdateContact(ctx, "tenant-a", &contact); err != nil {
		t.Fatalf("UpdateContact() error = %v", err)
	}
	if got, err := service.GetContact(ctx, "tenant-a", contact.PublicID); err != nil || got.DisplayName != "Alice Updated" {
		t.Fatalf("updated contact = %+v, %v", got, err)
	}
	contact.TenantID = "tenant-b"
	if err := service.UpdateContact(ctx, "tenant-a", &contact); !errors.Is(err, ErrInvalid) {
		t.Fatalf("cross-tenant UpdateContact() error = %v, want ErrInvalid", err)
	}

	status := ContactStatusActive
	contacts, err := service.ListContacts(ctx, "tenant-a", ContactFilter{Status: &status})
	if err != nil || len(contacts) != 1 || contacts[0].PublicID != "contact-a" {
		t.Fatalf("ListContacts() = %+v, %v", contacts, err)
	}
	if contacts, err := service.ListContacts(ctx, "tenant-b", ContactFilter{}); err != nil || len(contacts) != 1 {
		t.Fatalf("other tenant ListContacts() = %+v, %v", contacts, err)
	}
}

func TestContactServiceAddressCRUDValidatesOwnershipAndIdentity(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)
	contact := Contact{PublicID: "contact-a"}
	if err := service.CreateContact(ctx, "tenant-a", &contact); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}

	address, err := NewContactAddress("tenant-a", contact.ID, AddressKindEmail, "", "Alice@Example.com")
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	address.PublicID = "address-a"
	if err := service.CreateContactAddress(ctx, "tenant-a", &address); err != nil {
		t.Fatalf("CreateContactAddress() error = %v", err)
	}
	if address.NormalizedValue != "alice@example.com" || address.Identity.Value != address.NormalizedValue {
		t.Fatalf("address was not normalized = %+v", address)
	}
	if got, err := service.GetContactAddressByID(ctx, "tenant-a", address.ID); err != nil || got.PublicID != address.PublicID {
		t.Fatalf("GetContactAddressByID() = %+v, %v", got, err)
	}
	if addresses, err := service.ListContactAddressesByContactID(ctx, "tenant-a", contact.ID); err != nil || len(addresses) != 1 {
		t.Fatalf("ListContactAddressesByContactID() = %+v, %v", addresses, err)
	}

	duplicate, err := NewContactAddress("tenant-a", contact.ID, AddressKindEmail, "", " ALICE@example.com ")
	if err != nil {
		t.Fatalf("NewContactAddress(duplicate) error = %v", err)
	}
	if err := service.CreateContactAddress(ctx, "tenant-a", &duplicate); !errors.Is(err, ErrConflict) {
		t.Fatalf("duplicate identity error = %v, want ErrConflict", err)
	}
	if _, err := service.GetContactAddress(ctx, "tenant-b", address.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant GetContactAddress() error = %v, want ErrNotFound", err)
	}

	addresses, err := service.ListContactAddresses(ctx, "tenant-a", contact.PublicID)
	if err != nil || len(addresses) != 1 || addresses[0].PublicID != address.PublicID {
		t.Fatalf("ListContactAddresses() = %+v, %v", addresses, err)
	}
	address.Value = "alice.work"
	address.Identity = AddressIdentity{Kind: AddressKindUsername, Namespace: "chat", Value: "alice.work"}
	address.Kind = AddressKindUsername
	address.Namespace = "chat"
	address.ContactID = 0 // the service resolves immutable ownership from the public ID
	if err := service.UpdateContactAddress(ctx, "tenant-a", &address); err != nil {
		t.Fatalf("UpdateContactAddress() error = %v", err)
	}
	if address.NormalizedValue != "alice.work" || address.Identity.Value != "alice.work" {
		t.Fatalf("updated identity was not normalized = %+v", address)
	}
	if _, err := repository.FindContactAddressByIdentity(ctx, "tenant-a", AddressIdentity{Kind: AddressKindEmail, Namespace: "global", Value: "alice@example.com"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("old identity lookup error = %v, want ErrNotFound", err)
	}

	if err := service.DeleteContactAddress(ctx, "tenant-a", address.PublicID); err != nil {
		t.Fatalf("DeleteContactAddress() error = %v", err)
	}
	if _, err := service.GetContactAddress(ctx, "tenant-a", address.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted address lookup error = %v, want ErrNotFound", err)
	}
	if err := service.DeleteContactAddress(ctx, "tenant-a", address.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("repeated DeleteContactAddress() error = %v, want ErrNotFound", err)
	}
}

func TestContactServiceSoftDeleteCascadesAndReleasesIdentity(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)
	contact := Contact{PublicID: "contact-a"}
	if err := service.CreateContact(ctx, "tenant-a", &contact); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}
	address, err := NewContactAddress("tenant-a", contact.ID, AddressKindPhone, "", "+1 555 123 4567")
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	address.PublicID = "address-a"
	if err := service.CreateContactAddress(ctx, "tenant-a", &address); err != nil {
		t.Fatalf("CreateContactAddress() error = %v", err)
	}

	if err := service.DeleteContact(ctx, "tenant-a", contact.PublicID); err != nil {
		t.Fatalf("DeleteContact() error = %v", err)
	}
	if _, err := service.GetContact(ctx, "tenant-a", contact.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted contact lookup error = %v, want ErrNotFound", err)
	}
	if addresses, err := service.ListContactAddresses(ctx, "tenant-a", contact.PublicID); !errors.Is(err, ErrNotFound) || addresses != nil {
		t.Fatalf("deleted contact addresses = %+v, %v", addresses, err)
	}
	if err := service.DeleteContact(ctx, "tenant-a", contact.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("repeated DeleteContact() error = %v, want ErrNotFound", err)
	}

	replacement := Contact{PublicID: "contact-replacement"}
	if err := service.CreateContact(ctx, "tenant-a", &replacement); err != nil {
		t.Fatalf("CreateContact(replacement) error = %v", err)
	}
	reused, err := NewContactAddress("tenant-a", replacement.ID, AddressKindPhone, "", "+15551234567")
	if err != nil {
		t.Fatalf("NewContactAddress(reused) error = %v", err)
	}
	if err := service.CreateContactAddress(ctx, "tenant-a", &reused); err != nil {
		t.Fatalf("CreateContactAddress(reused identity) error = %v", err)
	}
}

func TestContactServiceRejectsInvalidAndCanceledOperations(t *testing.T) {
	service := NewContactService(NewInMemoryContactRepository())
	if err := service.CreateContact(context.Background(), "", &Contact{}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("blank tenant error = %v, want ErrInvalid", err)
	}
	if err := service.CreateContact(context.Background(), "tenant-a", nil); !errors.Is(err, ErrInvalid) {
		t.Fatalf("nil contact error = %v, want ErrInvalid", err)
	}
	if _, err := service.GetContact(context.Background(), "tenant-a", ""); !errors.Is(err, ErrInvalid) {
		t.Fatalf("blank public ID error = %v, want ErrInvalid", err)
	}
	invalidAddress := &ContactAddress{TenantID: "tenant-a", ContactID: 1, Kind: AddressKindEmail, Namespace: "global", Value: "not-an-email"}
	if err := service.CreateContactAddress(context.Background(), "tenant-a", invalidAddress); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid address error = %v, want ErrInvalid", err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	contact := Contact{PublicID: "contact-a"}
	if err := service.CreateContact(canceled, "tenant-a", &contact); err == nil || errors.Is(err, ErrInvalid) {
		t.Fatalf("canceled CreateContact() error = %v, want context cancellation", err)
	}
}
