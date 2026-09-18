package channel

import (
	"context"
	"errors"
	"testing"
	"time"
)

func testChannelAccount(tenantID, publicID string) ChannelAccount {
	return ChannelAccount{
		PublicID: publicID,
		TenantID: tenantID,
		Name:     "support",
		Type:     ChannelTypeWhatsApp,
		Provider: ProviderWhatsApp,
		Status:   StatePending,
		Capabilities: CapabilitySet{
			Provider: string(ProviderWhatsApp),
			Items:    []Capability{CapabilitySendText, CapabilityReceiveEvent},
		},
	}
}

func TestInMemoryChannelAccountRepositoryScopesPublicIDsToTenant(t *testing.T) {
	repository := NewInMemoryChannelAccountRepository()
	ctx := context.Background()
	account := testChannelAccount("tenant-a", "channel-a")

	if err := repository.CreateChannelAccount(ctx, &account); err != nil {
		t.Fatalf("CreateChannelAccount() error = %v", err)
	}
	if account.ID == 0 || account.CreatedAt.IsZero() || account.UpdatedAt.IsZero() {
		t.Fatalf("CreateChannelAccount() did not assign persistence fields: %+v", account)
	}

	got, err := repository.GetChannelAccountByPublicID(ctx, "tenant-a", "channel-a")
	if err != nil {
		t.Fatalf("GetChannelAccountByPublicID() error = %v", err)
	}
	if got.ID != account.ID || got.TenantID != "tenant-a" {
		t.Fatalf("GetChannelAccountByPublicID() = %+v", got)
	}
	if _, err := repository.GetChannelAccountByPublicID(ctx, "tenant-b", "channel-a"); !errors.Is(err, ErrChannelAccountNotFound) {
		t.Fatalf("cross-tenant lookup error = %v, want ErrChannelAccountNotFound", err)
	}

	otherTenant := testChannelAccount("tenant-b", "channel-b")
	if err := repository.CreateChannelAccount(ctx, &otherTenant); err != nil {
		t.Fatalf("second tenant create returned error = %v", err)
	}
	if _, err := repository.GetChannelAccountByPublicID(ctx, "tenant-b", "channel-a"); err == nil {
		t.Fatal("tenant-scoped lookup returned another tenant's account")
	}
	if _, err := repository.GetChannelAccountByPublicID(ctx, "tenant-b", "channel-b"); err != nil {
		t.Fatalf("second tenant lookup error = %v", err)
	}
	duplicate := testChannelAccount("tenant-b", "channel-a")
	if err := repository.CreateChannelAccount(ctx, &duplicate); !errors.Is(err, ErrChannelAccountConflict) {
		t.Fatalf("duplicate public ID error = %v, want ErrChannelAccountConflict", err)
	}
}

func TestInMemoryChannelAccountRepositoryListUpdateAndDefensiveCopies(t *testing.T) {
	repository := NewInMemoryChannelAccountRepository()
	ctx := context.Background()
	account := testChannelAccount("tenant-a", "channel-a")
	account.Credential = &CredentialReference{ID: "credential-a", Kind: CredentialKindOAuth}
	if err := repository.CreateChannelAccount(ctx, &account); err != nil {
		t.Fatalf("CreateChannelAccount() error = %v", err)
	}

	account.Capabilities.Items[0] = CapabilityGetUser
	account.Name = "mutated caller"
	got, err := repository.GetChannelAccountByPublicID(ctx, "tenant-a", "channel-a")
	if err != nil {
		t.Fatalf("GetChannelAccountByPublicID() error = %v", err)
	}
	if got.Name != "support" || !got.Capabilities.Has(CapabilitySendText) {
		t.Fatalf("repository retained caller mutation: %+v", got)
	}

	got.Name = "updated support"
	got.Status = StateConnected
	if err := repository.UpdateChannelAccount(ctx, got); err != nil {
		t.Fatalf("UpdateChannelAccount() error = %v", err)
	}
	status := StateConnected
	accounts, err := repository.ListChannelAccounts(ctx, "tenant-a", ChannelAccountFilter{Status: &status})
	if err != nil {
		t.Fatalf("ListChannelAccounts() error = %v", err)
	}
	if len(accounts) != 1 || accounts[0].Name != "updated support" {
		t.Fatalf("ListChannelAccounts() = %+v", accounts)
	}

	if err := repository.DeleteChannelAccount(ctx, "tenant-a", "channel-a"); err != nil {
		t.Fatalf("DeleteChannelAccount() error = %v", err)
	}
	if _, err := repository.GetChannelAccountByPublicID(ctx, "tenant-a", "channel-a"); !errors.Is(err, ErrChannelAccountNotFound) {
		t.Fatalf("deleted lookup error = %v", err)
	}
	if accounts, err := repository.ListChannelAccounts(ctx, "tenant-a", ChannelAccountFilter{}); err != nil || len(accounts) != 0 {
		t.Fatalf("deleted account list = %v, %v", accounts, err)
	}
}

func TestInMemoryChannelAccountRepositoryConnectionEventsAreTenantScopedAndAppendOnly(t *testing.T) {
	repository := NewInMemoryChannelAccountRepository()
	ctx := context.Background()
	account := testChannelAccount("tenant-a", "channel-a")
	if err := repository.CreateChannelAccount(ctx, &account); err != nil {
		t.Fatalf("CreateChannelAccount() error = %v", err)
	}

	event := ConnectionEvent{
		TenantID:               "tenant-a",
		ChannelAccountPublicID: "channel-a",
		Type:                   ConnectionEventConnected,
		FromState:              StateConnecting,
		ToState:                StateConnected,
		Metadata:               map[string]string{"source": "test"},
	}
	if err := repository.RecordConnectionEvent(ctx, &event); err != nil {
		t.Fatalf("RecordConnectionEvent() error = %v", err)
	}
	if event.ID == 0 || event.PublicID == "" || event.OccurredAt.IsZero() {
		t.Fatalf("RecordConnectionEvent() did not assign event fields: %+v", event)
	}

	events, err := repository.ListConnectionEvents(ctx, "tenant-a", "channel-a")
	if err != nil {
		t.Fatalf("ListConnectionEvents() error = %v", err)
	}
	if len(events) != 1 || events[0].Type != ConnectionEventConnected {
		t.Fatalf("ListConnectionEvents() = %+v", events)
	}
	events[0].Metadata["source"] = "mutated"
	events, err = repository.ListConnectionEvents(ctx, "tenant-a", "channel-a")
	if err != nil || events[0].Metadata["source"] != "test" {
		t.Fatalf("event snapshot was not defensive: %+v, %v", events, err)
	}

	updated, err := repository.GetChannelAccountByPublicID(ctx, "tenant-a", "channel-a")
	if err != nil || updated.LastConnectedAt == nil || updated.LastEventAt == nil {
		t.Fatalf("event did not update account references: %+v, %v", updated, err)
	}
	if _, err := repository.ListConnectionEvents(ctx, "tenant-b", "channel-a"); !errors.Is(err, ErrChannelAccountNotFound) {
		t.Fatalf("cross-tenant event lookup error = %v", err)
	}

	if err := repository.RecordConnectionEvent(ctx, &ConnectionEvent{
		TenantID:               "tenant-b",
		ChannelAccountPublicID: "channel-a",
		Type:                   ConnectionEventError,
		ToState:                StateError,
		OccurredAt:             time.Now(),
	}); !errors.Is(err, ErrChannelAccountNotFound) {
		t.Fatalf("cross-tenant event write error = %v", err)
	}
}

func TestChannelAccountReferencesNeverContainSecretValues(t *testing.T) {
	account := testChannelAccount("tenant-a", "channel-a")
	account.Credential = &CredentialReference{ID: "vault-ref", Kind: CredentialKindAccessToken}
	account.Session = &SessionReference{ID: "session-ref", Kind: SessionKindRuntime}
	if err := account.Valid(); err != nil {
		t.Fatalf("ChannelAccount.Valid() error = %v", err)
	}
	if account.Credential.ID == "" || account.Session.ID == "" {
		t.Fatal("references must retain opaque IDs")
	}
}
