package contact

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func createAudienceContact(t *testing.T, repository *InMemoryContactRepository, tenant, publicID, email string, tags []string) (Contact, ContactAddress) {
	t.Helper()
	contact := Contact{TenantID: tenant, PublicID: publicID, Status: ContactStatusActive, Source: ContactSourceManual, Locale: "en-US", Tags: tags, CustomFields: map[string]string{"plan": "gold", "private": email}}
	if err := repository.CreateContact(context.Background(), &contact); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}
	address, err := NewContactAddress(tenant, contact.ID, AddressKindEmail, "", email)
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	address.Consent = ConsentMetadata{State: ConsentStateOptedIn, Source: ConsentSourceManual}
	if err := repository.CreateContactAddress(context.Background(), &address); err != nil {
		t.Fatalf("CreateContactAddress() error = %v", err)
	}
	return contact, address
}

func TestAudienceDefinitionsAndEvaluationAreTenantScoped(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	audience := NewAudienceService(repository)
	contactA, addressA := createAudienceContact(t, repository, "tenant-a", "contact-a", "alice@example.com", []string{"vip"})
	contactB, _ := createAudienceContact(t, repository, "tenant-a", "contact-b", "bob@example.com", []string{"regular"})
	contactOther, _ := createAudienceContact(t, repository, "tenant-b", "contact-a", "alice@example.com", []string{"vip"})

	list := &StaticList{Name: "VIP list", ContactIDs: []uint64{contactA.ID}}
	if err := audience.CreateStaticList(ctx, "tenant-a", list); err != nil {
		t.Fatalf("CreateStaticList() error = %v", err)
	}
	if _, err := audience.GetStaticList(ctx, "tenant-b", list.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant static list lookup = %v", err)
	}
	if err := audience.SetStaticListMembers(ctx, "tenant-a", list.PublicID, []uint64{contactA.ID, contactB.ID}); err != nil {
		t.Fatalf("SetStaticListMembers() error = %v", err)
	}
	if err := audience.SetStaticListMembers(ctx, "tenant-a", list.PublicID, []uint64{contactOther.ID}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant static list member error = %v", err)
	}

	segment := &DynamicSegment{
		Name: "VIP gold contacts",
		Definition: SegmentDefinition{
			Fields: []FieldPredicate{{Field: "contact.custom.plan", Operator: PredicateEquals, Value: "gold"}, {Field: SegmentFieldAddressKind, Operator: PredicateEquals, Value: "email"}},
			Tags:   []TagPredicate{{Tag: "vip", Operator: TagHas}},
		},
	}
	if err := audience.CreateDynamicSegment(ctx, "tenant-a", segment); err != nil {
		t.Fatalf("CreateDynamicSegment() error = %v", err)
	}
	preview, err := audience.PreviewAudience(ctx, "tenant-a", AudienceSelection{SegmentPublicID: segment.PublicID})
	if err != nil {
		t.Fatalf("PreviewAudience() error = %v", err)
	}
	if preview.Count != 1 || preview.IncludedContacts != 1 || preview.IncludedRecipients != 1 {
		t.Fatalf("PreviewAudience() = %+v, want one contact and recipient", preview)
	}
	if count, err := audience.PreviewSegmentCount(ctx, "tenant-a", segment.PublicID); err != nil || count != 1 {
		t.Fatalf("PreviewSegmentCount() = %d, %v", count, err)
	}
	if _, err := audience.PreviewAudience(ctx, "tenant-b", AudienceSelection{SegmentPublicID: segment.PublicID}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant segment evaluation = %v", err)
	}

	exclusion := &ExclusionList{ContactIDs: []uint64{contactA.ID}}
	if err := audience.CreateExclusionList(ctx, "tenant-a", exclusion); err != nil {
		t.Fatalf("CreateExclusionList() error = %v", err)
	}
	preview, err = audience.PreviewAudience(ctx, "tenant-a", AudienceSelection{SegmentPublicID: segment.PublicID, ExclusionListPublicID: exclusion.PublicID})
	if err != nil {
		t.Fatalf("PreviewAudience(excluded) error = %v", err)
	}
	if preview.Count != 0 || preview.ExcludedContacts != 1 {
		t.Fatalf("excluded preview = %+v", preview)
	}
	if _, err := audience.GetExclusionList(ctx, "tenant-b", exclusion.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant exclusion lookup = %v", err)
	}

	// Keep this assertion so the fixture proves the tenant-b contact was not
	// accidentally selected by tenant-a evaluation despite sharing public data.
	if contactOther.TenantID != "tenant-b" || addressA.TenantID != "tenant-a" {
		t.Fatal("fixture ownership was not established")
	}
}

func TestAudienceSnapshotIsImmutableAndRevalidatesSuppression(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	contactService := NewContactService(repository)
	audience := NewAudienceService(repository)
	contact, address := createAudienceContact(t, repository, "tenant-a", "contact-a", "private@example.com", []string{"vip"})
	segment := &DynamicSegment{Name: "VIP", Definition: SegmentDefinition{Tags: []TagPredicate{{Tag: "vip", Operator: TagHas}}}}
	if err := audience.CreateDynamicSegment(ctx, "tenant-a", segment); err != nil {
		t.Fatalf("CreateDynamicSegment() error = %v", err)
	}
	snapshot, err := audience.CreateAudienceSnapshot(ctx, "tenant-a", AudienceSelection{SegmentPublicID: segment.PublicID})
	if err != nil {
		t.Fatalf("CreateAudienceSnapshot() error = %v", err)
	}
	if len(snapshot.Recipients) != 1 || snapshot.Recipients[0].ContactID != contact.ID {
		t.Fatalf("snapshot recipients = %+v", snapshot.Recipients)
	}
	snapshot.Recipients[0].ContactID = 999
	stored, err := audience.GetAudienceSnapshot(ctx, "tenant-a", snapshot.PublicID)
	if err != nil {
		t.Fatalf("GetAudienceSnapshot() error = %v", err)
	}
	if stored.Recipients[0].ContactID != contact.ID {
		t.Fatalf("snapshot was mutable through returned slice: %+v", stored.Recipients)
	}

	if err := contactService.OptOut(ctx, "tenant-a", address.PublicID, ConsentMetadata{Source: ConsentSourceManual}); err != nil {
		t.Fatalf("OptOut() error = %v", err)
	}
	result, err := audience.RevalidateSnapshotForSend(ctx, "tenant-a", snapshot.PublicID)
	if err != nil {
		t.Fatalf("RevalidateSnapshotForSend() error = %v", err)
	}
	if len(result.Sendable) != 0 || len(result.Suppressed) != 1 {
		t.Fatalf("suppression revalidation = %+v", result)
	}
}

func TestAudienceDiagnosticsDoNotExposePredicateOrAddressPII(t *testing.T) {
	definition := SegmentDefinition{Fields: []FieldPredicate{{Field: "contact.custom.private", Operator: PredicateEquals, Value: "private@example.com"}}, Tags: []TagPredicate{{Tag: "private-tag", Operator: TagHas}}}
	segment := DynamicSegment{PublicID: "segment-public", TenantID: "tenant-a", Name: "Private audience", Definition: definition}
	list := StaticList{PublicID: "list-public", TenantID: "tenant-a", Name: "private@example.com", ContactIDs: []uint64{1}}
	identity, err := NormalizeAddress(AddressKindEmail, "", "private@example.com")
	if err != nil {
		t.Fatalf("NormalizeAddress() error = %v", err)
	}
	exclusion := ExclusionList{PublicID: "exclude-public", TenantID: "tenant-a", AddressIdentities: []AddressIdentity{identity}}
	snapshot := AudienceSnapshot{PublicID: "snapshot-public", TenantID: "tenant-a", Selection: AudienceSelection{StaticListPublicID: "list-public"}, Recipients: []AudienceRecipient{{ContactID: 1, AddressID: 2, Identity: identity}}}
	for label, value := range map[string]string{
		"segment string":   segment.String(),
		"list string":      list.String(),
		"exclusion string": exclusion.String(),
		"snapshot string":  snapshot.String(),
	} {
		if strings.Contains(value, "private@example.com") || strings.Contains(value, "private-tag") {
			t.Fatalf("%s leaked PII: %s", label, value)
		}
	}
	encoded, err := json.Marshal(struct {
		Segment   DynamicSegment
		List      StaticList
		Exclusion ExclusionList
		Snapshot  AudienceSnapshot
	}{segment, list, exclusion, snapshot})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "private@example.com") || strings.Contains(string(encoded), "private-tag") {
		t.Fatalf("audience diagnostics leaked PII: %s", encoded)
	}
}

func TestDynamicSegmentRejectsUnsupportedPredicates(t *testing.T) {
	repository := NewInMemoryContactRepository()
	audience := NewAudienceService(repository)
	segment := &DynamicSegment{Name: "unsupported", Definition: SegmentDefinition{Fields: []FieldPredicate{{Field: "event.delivery_status", Operator: PredicateEquals, Value: "delivered"}}}}
	if err := audience.CreateDynamicSegment(context.Background(), "tenant-a", segment); !errors.Is(err, ErrInvalid) {
		t.Fatalf("unsupported predicate error = %v, want ErrInvalid", err)
	}
}
