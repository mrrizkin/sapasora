package contact

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
)

func createContactAddressForTest(t *testing.T, service ContactService, tenant string, contact Contact, publicID string, kind AddressKind, value string) ContactAddress {
	t.Helper()
	address, err := NewContactAddress(tenant, contact.ID, kind, "", value)
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	address.PublicID = publicID
	if err := service.CreateContactAddress(context.Background(), tenant, &address); err != nil {
		t.Fatalf("CreateContactAddress() error = %v", err)
	}
	return address
}

func TestDuplicateCandidatesAreTenantScopedDeterministicAndSafe(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)
	first := newTestContact(t, repository, "tenant-a", "first")
	second := newTestContact(t, repository, "tenant-a", "second")
	otherTenant := newTestContact(t, repository, "tenant-b", "other")
	createContactAddressForTest(t, service, "tenant-a", first, "first-email", AddressKindEmail, "first@example.com")
	createContactAddressForTest(t, service, "tenant-a", second, "second-email", AddressKindEmail, "second@example.com")
	createContactAddressForTest(t, service, "tenant-b", otherTenant, "other-email", AddressKindEmail, "first@example.com")

	firstIdentity, _ := NormalizeAddress(AddressKindEmail, "", "FIRST@example.com")
	secondIdentity, _ := NormalizeAddress(AddressKindEmail, "", "second@example.com")
	candidates, err := service.FindDuplicateCandidates(ctx, "tenant-a", []AddressIdentity{secondIdentity, firstIdentity, firstIdentity})
	if err != nil {
		t.Fatalf("FindDuplicateCandidates() error = %v", err)
	}
	if len(candidates) != 2 || candidates[0].ContactPublicID != "first" || candidates[1].ContactPublicID != "second" {
		t.Fatalf("candidate order = %+v", candidates)
	}
	if candidates[0].MatchedIdentityCount != 1 || candidates[0].Matches[0].IdentityKey != firstIdentity.Key() {
		t.Fatalf("candidate match = %+v", candidates[0])
	}
	encoded, err := json.Marshal(candidates)
	if err != nil {
		t.Fatalf("json.Marshal(candidates) error = %v", err)
	}
	if strings.Contains(string(encoded), "first@example.com") || strings.Contains(string(encoded), "FIRST@example.com") {
		t.Fatalf("candidate diagnostics leaked PII: %s", encoded)
	}
	if candidates, err := service.FindDuplicateCandidates(ctx, "tenant-b", []AddressIdentity{firstIdentity}); err != nil || len(candidates) != 1 || candidates[0].ContactPublicID != "other" {
		t.Fatalf("tenant-b candidates = %+v, %v", candidates, err)
	}
	if candidates, err := service.FindDuplicateCandidates(ctx, "tenant-a", []AddressIdentity{firstIdentity}); err != nil || len(candidates) != 1 || candidates[0].ContactPublicID != "first" {
		t.Fatalf("tenant-a candidates = %+v, %v", candidates, err)
	}
}

func TestMergeRequiresFreshExplicitConfirmationAndSupportsSafeUndo(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)
	source := newTestContact(t, repository, "tenant-a", "source")
	target := newTestContact(t, repository, "tenant-a", "target")
	sourceAddress := createContactAddressForTest(t, service, "tenant-a", source, "source-email", AddressKindEmail, "source@example.com")
	createContactAddressForTest(t, service, "tenant-a", target, "target-phone", AddressKindPhone, "+1 555 123 4567")

	request := MergeRequest{SourceContactPublicID: source.PublicID, TargetContactPublicID: target.PublicID}
	preview, err := service.PreviewMerge(ctx, "tenant-a", request)
	if err != nil {
		t.Fatalf("PreviewMerge() error = %v", err)
	}
	if !preview.CanMerge || preview.SourceAddressCount != 1 || preview.TargetAddressCount != 1 || preview.ConfirmationToken == "" {
		t.Fatalf("merge preview = %+v", preview)
	}
	if _, err := service.MergeContacts(ctx, "tenant-a", MergeRequest{SourceContactPublicID: source.PublicID, TargetContactPublicID: target.PublicID, ConfirmationToken: preview.ConfirmationToken}); !errors.Is(err, ErrMergeConfirmationRequired) {
		t.Fatalf("unconfirmed merge error = %v, want ErrMergeConfirmationRequired", err)
	}
	mergeRequest := request
	mergeRequest.ConfirmationToken = preview.ConfirmationToken
	mergeRequest.Confirmed = true
	audit, err := service.MergeContacts(ctx, "tenant-a", mergeRequest)
	if err != nil {
		t.Fatalf("MergeContacts() error = %v", err)
	}
	if audit.ID == "" || len(audit.Addresses) != 1 || audit.Addresses[0].AddressID != sourceAddress.ID || !audit.Undoable || audit.Undone {
		t.Fatalf("merge audit = %+v", audit)
	}
	if _, err := service.GetContact(ctx, "tenant-a", source.PublicID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("merged source lookup error = %v, want ErrNotFound", err)
	}
	addresses, err := service.ListContactAddressesByContactID(ctx, "tenant-a", target.ID)
	if err != nil || len(addresses) != 2 {
		t.Fatalf("target addresses after merge = %+v, %v", addresses, err)
	}
	undo, err := service.UndoMerge(ctx, "tenant-a", audit.ID)
	if err != nil {
		t.Fatalf("UndoMerge() error = %v", err)
	}
	if len(undo.RestoredAddressIDs) != 1 || undo.RestoredAddressIDs[0] != sourceAddress.ID {
		t.Fatalf("undo audit = %+v", undo)
	}
	if restored, err := service.GetContact(ctx, "tenant-a", source.PublicID); err != nil || restored.ID != source.ID {
		t.Fatalf("restored source = %+v, %v", restored, err)
	}
	audits, err := service.ListMergeAudits(ctx, "tenant-a")
	if err != nil || len(audits) != 1 || !audits[0].Undone || audits[0].Undoable {
		t.Fatalf("merge audits after undo = %+v, %v", audits, err)
	}
	if _, err := service.UndoMerge(ctx, "tenant-a", audit.ID); !errors.Is(err, ErrConflict) {
		t.Fatalf("repeated UndoMerge() error = %v, want ErrConflict", err)
	}
}

func TestMergeRejectsCrossTenantSameContactAndStalePreview(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)
	source := newTestContact(t, repository, "tenant-a", "source")
	target := newTestContact(t, repository, "tenant-a", "target")
	other := newTestContact(t, repository, "tenant-b", "target")
	createContactAddressForTest(t, service, "tenant-a", source, "source-email", AddressKindEmail, "source@example.com")
	request := MergeRequest{SourceContactPublicID: source.PublicID, TargetContactPublicID: target.PublicID}
	preview, err := service.PreviewMerge(ctx, "tenant-a", request)
	if err != nil {
		t.Fatalf("PreviewMerge() error = %v", err)
	}
	if _, err := service.PreviewMerge(ctx, "tenant-b", MergeRequest{SourceContactPublicID: source.PublicID, TargetContactPublicID: other.PublicID}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant preview error = %v, want ErrNotFound", err)
	}
	samePreview, err := service.PreviewMerge(ctx, "tenant-a", MergeRequest{SourceContactPublicID: source.PublicID, TargetContactPublicID: source.PublicID})
	if err != nil || samePreview.CanMerge || len(samePreview.Conflicts) != 1 || samePreview.Conflicts[0].Code != "same_contact" {
		t.Fatalf("same-contact preview = %+v, %v", samePreview, err)
	}
	address := createContactAddressForTest(t, service, "tenant-a", target, "target-email", AddressKindEmail, "target@example.com")
	address.Value = "target2@example.com"
	if err := service.UpdateContactAddress(ctx, "tenant-a", &address); err != nil {
		t.Fatalf("UpdateContactAddress() error = %v", err)
	}
	mergeRequest := request
	mergeRequest.ConfirmationToken = preview.ConfirmationToken
	mergeRequest.Confirmed = true
	if _, err := service.MergeContacts(ctx, "tenant-a", mergeRequest); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale merge error = %v, want ErrConflict", err)
	}
}

func TestMergeConcurrentConfirmationIsSingleWinner(t *testing.T) {
	ctx := context.Background()
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)
	source := newTestContact(t, repository, "tenant-a", "source")
	target := newTestContact(t, repository, "tenant-a", "target")
	createContactAddressForTest(t, service, "tenant-a", source, "source-email", AddressKindEmail, "source@example.com")
	preview, err := service.PreviewMerge(ctx, "tenant-a", MergeRequest{SourceContactPublicID: source.PublicID, TargetContactPublicID: target.PublicID})
	if err != nil {
		t.Fatalf("PreviewMerge() error = %v", err)
	}
	request := MergeRequest{SourceContactPublicID: source.PublicID, TargetContactPublicID: target.PublicID, ConfirmationToken: preview.ConfirmationToken, Confirmed: true}
	const workers = 16
	results := make(chan error, workers)
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			_, err := service.MergeContacts(ctx, "tenant-a", request)
			results <- err
		}()
	}
	wait.Wait()
	close(results)
	wins := 0
	for err := range results {
		if err == nil {
			wins++
		} else if !errors.Is(err, ErrConflict) && !errors.Is(err, ErrNotFound) {
			t.Fatalf("concurrent merge error = %v", err)
		}
	}
	if wins != 1 {
		t.Fatalf("concurrent merge wins = %d, want 1", wins)
	}
}
