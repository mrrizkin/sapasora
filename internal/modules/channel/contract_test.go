package channel

import (
	"context"
	"errors"
	"testing"
	"time"
)

var (
	_ Connect               = (*FakeAdapter)(nil)
	_ Disconnect            = (*FakeAdapter)(nil)
	_ Health                = (*FakeAdapter)(nil)
	_ Capabilities          = (*FakeAdapter)(nil)
	_ Send                  = (*FakeAdapter)(nil)
	_ GetContact            = (*FakeAdapter)(nil)
	_ GetUser               = (*FakeAdapter)(nil)
	_ ReceiveEvent          = (*FakeAdapter)(nil)
	_ ValidateAddress       = (*FakeAdapter)(nil)
	_ NormalizeError        = (*FakeAdapter)(nil)
	_ NormalizeInboundEvent = (*FakeAdapter)(nil)
)

func TestFakeAdapterContract(t *testing.T) {
	ctx := context.Background()
	adapter := NewFakeAdapter()

	state, err := adapter.Connect(ctx, ConnectRequest{AccountID: "account-1", CredentialRef: "secret-ref"})
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	if state != StateConnected {
		t.Fatalf("Connect() state = %q, want %q", state, StateConnected)
	}

	health, err := adapter.Health(ctx)
	if err != nil {
		t.Fatalf("Health() error = %v", err)
	}
	if !health.Healthy || health.State != StateConnected {
		t.Fatalf("Health() = %+v, want connected and healthy", health)
	}
	if !adapter.Capabilities().Has(CapabilitySendText) {
		t.Fatal("Capabilities() does not advertise text sending")
	}

	address, err := adapter.ValidateAddress(ctx, Address{Value: "  user-1  "})
	if err != nil {
		t.Fatalf("ValidateAddress() error = %v", err)
	}
	if address != (Address{Kind: AddressExternal, Value: "user-1"}) {
		t.Fatalf("ValidateAddress() = %+v", address)
	}

	contact := Contact{Address: address, DisplayName: "User One"}
	adapter.SetContact(contact)
	gotContact, err := adapter.GetContact(ctx, GetContactRequest{Address: address})
	if err != nil {
		t.Fatalf("GetContact() error = %v", err)
	}
	if gotContact != contact {
		t.Fatalf("GetContact() = %+v, want %+v", gotContact, contact)
	}

	user := User{Address: address, DisplayName: "User One", Username: "user1"}
	adapter.SetUser(user)
	gotUser, err := adapter.GetUser(ctx, GetUserRequest{Address: address})
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if gotUser != user {
		t.Fatalf("GetUser() = %+v, want %+v", gotUser, user)
	}

	result, err := adapter.Send(ctx, SendRequest{To: address, Message: Message{Kind: MessageText, Body: "hello"}})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if result.MessageID != "fake-message-1" || result.Status != DeliveryAccepted {
		t.Fatalf("Send() = %+v, want deterministic accepted result", result)
	}

	adapter.QueueEvent(RawEvent{ID: "event-1", Type: string(InboundMessageEvent), Payload: []byte(`{"kind":"message"}`)})
	raw, err := adapter.ReceiveEvent(ctx)
	if err != nil {
		t.Fatalf("ReceiveEvent() error = %v", err)
	}
	event, err := adapter.NormalizeInboundEvent(ctx, raw)
	if err != nil {
		t.Fatalf("NormalizeInboundEvent() error = %v", err)
	}
	if event.ID != "event-1" || event.Kind != InboundMessageEvent {
		t.Fatalf("normalized event = %+v", event)
	}
	if _, err := adapter.ReceiveEvent(ctx); !errors.Is(err, ErrNoEvent) {
		t.Fatalf("empty ReceiveEvent() error = %v, want ErrNoEvent", err)
	}

	if err := adapter.Disconnect(ctx, DisconnectRequest{Reason: "test"}); err != nil {
		t.Fatalf("Disconnect() error = %v", err)
	}
	if got := adapter.LastCall(); got != "disconnect" {
		t.Fatalf("LastCall() = %q, want disconnect", got)
	}
}

func TestFakeAdapterContractHonorsConfiguredErrors(t *testing.T) {
	adapter := NewFakeAdapter()
	adapter.SetError("send", errors.New("send failed"))

	_, err := adapter.Send(context.Background(), SendRequest{})
	if err == nil || err.Error() != "send failed" {
		t.Fatalf("Send() error = %v, want configured error", err)
	}

	adapter.SetError("validate_address", nil)
	_, err = adapter.ValidateAddress(context.Background(), Address{})
	if !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("ValidateAddress() error = %v, want ErrInvalidAddress", err)
	}
}

func TestErrorMappingCatalog(t *testing.T) {
	catalog, err := NewErrorMappingCatalog([]ErrorMapping{
		{ProviderCode: "too_many_requests", Code: "rate_limited", Category: ErrorCategoryRateLimited, Retryable: true, RetryAfter: time.Second},
	})
	if err != nil {
		t.Fatalf("NewErrorMappingCatalog() error = %v", err)
	}

	mapping, ok := catalog.Lookup("too_many_requests")
	if !ok || mapping.Category != ErrorCategoryRateLimited || !mapping.Retryable {
		t.Fatalf("Lookup() = %+v, found = %v", mapping, ok)
	}

	normalized := catalog.Normalize("provider-x", "too_many_requests", errors.New("provider detail"))
	if normalized.Provider != "provider-x" || normalized.Code != "rate_limited" || !normalized.Retryable {
		t.Fatalf("Normalize() = %+v", normalized)
	}
	if !IsCategory(normalized, ErrorCategoryRateLimited) {
		t.Fatalf("IsCategory() = false for %+v", normalized)
	}

	if _, err := NewErrorMappingCatalog([]ErrorMapping{
		{ProviderCode: "duplicate", Code: "one", Category: ErrorCategoryInternal},
		{ProviderCode: "duplicate", Code: "two", Category: ErrorCategoryInternal},
	}); err == nil {
		t.Fatal("NewErrorMappingCatalog() accepted duplicate provider code")
	}
}
