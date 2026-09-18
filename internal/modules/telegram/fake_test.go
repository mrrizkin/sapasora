package telegram

import (
	"context"
	"errors"
	"testing"

	"sapasora/internal/modules/device"
)

func TestFakeAdapterRecordsTelegramOperationAndPropagatesError(t *testing.T) {
	wantErr := errors.New("provider unavailable")
	adapter := &FakeAdapter{Err: wantErr}

	_, err := adapter.GetContacts(context.Background(), &device.Device{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("GetContacts() error = %v, want %v", err, wantErr)
	}
	if got := adapter.LastCall(); got != "get_contacts" {
		t.Fatalf("LastCall() = %q, want get_contacts", got)
	}
	if got := adapter.Calls(); len(got) != 1 || got[0] != "get_contacts" {
		t.Fatalf("Calls() = %v, want [get_contacts]", got)
	}
}

func TestFakeAdapterRecordsCallsWithoutSharingBackingStorage(t *testing.T) {
	adapter := &FakeAdapter{}
	_, _ = adapter.GetStatus(context.Background(), &device.Device{})
	_, _ = adapter.SendText(context.Background(), &device.Device{}, &SendTextRequest{})

	calls := adapter.Calls()
	calls[0] = "mutated"
	if got := adapter.Calls()[0]; got != "get_status" {
		t.Fatalf("Calls() exposed backing storage; got %q", got)
	}
	if got := adapter.LastCall(); got != "send_text" {
		t.Fatalf("LastCall() = %q, want send_text", got)
	}
}
