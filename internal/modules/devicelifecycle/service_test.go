package devicelifecycle

import (
	"context"
	"errors"
	"testing"

	"sapasora/internal/modules/device"
	"sapasora/internal/modules/devicetoken"

	"gorm.io/gorm"
)

type fakeDisconnector struct {
	events *[]string
	err    error
}

func (f *fakeDisconnector) Disconnect(context.Context, *device.Device) error {
	*f.events = append(*f.events, "disconnect")
	return f.err
}

type fakeTokens struct {
	devicetoken.DeviceTokenService
	events    *[]string
	tokens    []*devicetoken.DeviceToken
	failToken uint
	err       error
}

func (f *fakeTokens) ListDeviceTokensByDeviceID(context.Context, uint) ([]*devicetoken.DeviceToken, error) {
	result := append([]*devicetoken.DeviceToken(nil), f.tokens...)
	return result, nil
}

func (f *fakeTokens) DeleteDeviceToken(_ context.Context, token *devicetoken.DeviceToken) error {
	*f.events = append(*f.events, "token:"+token.PublicID)
	if token.ID == f.failToken {
		return f.err
	}
	for i, candidate := range f.tokens {
		if candidate == token {
			f.tokens = append(f.tokens[:i], f.tokens[i+1:]...)
			return nil
		}
	}
	return nil
}

type fakeDevices struct {
	device.DeviceService
	events *[]string
	calls  int
}

func (f *fakeDevices) DeleteDevice(_ context.Context, d *device.Device) error {
	*f.events = append(*f.events, "device")
	f.calls++
	d.DeletedAt = gorm.DeletedAt{Valid: true}
	return nil
}

func newLifecycleService(events *[]string, tokens []*devicetoken.DeviceToken) (*Service, *fakeTokens, *fakeDevices) {
	fakeTokenService := &fakeTokens{events: events, tokens: tokens}
	fakeDeviceService := &fakeDevices{events: events}
	return &Service{
		devices: fakeDeviceService,
		tokens:  fakeTokenService,
		gateway: &fakeDisconnector{events: events},
	}, fakeTokenService, fakeDeviceService
}

func TestDeleteDeviceDisconnectsAndRemovesAllTokensBeforeDevice(t *testing.T) {
	events := []string{}
	service, _, devices := newLifecycleService(&events, []*devicetoken.DeviceToken{
		{ID: 1, PublicID: "token-1"},
		{ID: 2, PublicID: "token-2"},
		{ID: 3, PublicID: "token-3"},
	})

	err := service.DeleteDevice(context.Background(), &device.Device{ID: 42, PublicID: "device-1"})

	if err != nil {
		t.Fatalf("DeleteDevice() error = %v", err)
	}
	if got, want := events, []string{"disconnect", "token:token-1", "token:token-2", "token:token-3", "device"}; !equalStrings(got, want) {
		t.Fatalf("cleanup order = %#v, want %#v", got, want)
	}
	if devices.calls != 1 {
		t.Fatalf("device delete calls = %d, want 1", devices.calls)
	}
}

func TestDeleteDeviceDoesNotContinueAfterDisconnectFailure(t *testing.T) {
	events := []string{}
	service, tokens, devices := newLifecycleService(&events, []*devicetoken.DeviceToken{{ID: 1, PublicID: "secret-token"}})
	service.gateway = &fakeDisconnector{events: &events, err: errors.New("provider secret")}

	err := service.DeleteDevice(context.Background(), &device.Device{ID: 42})

	if !errors.Is(err, ErrDeviceDeletion) {
		t.Fatalf("error = %v, want ErrDeviceDeletion", err)
	}
	if err.Error() == "" || contains(err.Error(), "provider secret") || contains(err.Error(), "secret-token") {
		t.Fatalf("deletion error leaked sensitive detail: %q", err)
	}
	if len(tokens.tokens) != 1 || devices.calls != 0 {
		t.Fatalf("cleanup continued after disconnect failure: tokens=%d device_calls=%d", len(tokens.tokens), devices.calls)
	}
}

func TestDeleteDeviceRetriesRemainingTokensAfterFailure(t *testing.T) {
	events := []string{}
	service, tokens, devices := newLifecycleService(&events, []*devicetoken.DeviceToken{
		{ID: 1, PublicID: "token-1"},
		{ID: 2, PublicID: "token-2"},
		{ID: 3, PublicID: "token-3"},
	})
	tokens.failToken = 2
	tokens.err = errors.New("token store unavailable")
	d := &device.Device{ID: 42}

	if err := service.DeleteDevice(context.Background(), d); err == nil {
		t.Fatal("first deletion succeeded despite token cleanup failure")
	}
	if devices.calls != 0 || len(tokens.tokens) != 2 {
		t.Fatalf("failed deletion removed the device or wrong tokens: device_calls=%d tokens=%d", devices.calls, len(tokens.tokens))
	}

	tokens.failToken = 0
	if err := service.DeleteDevice(context.Background(), d); err != nil {
		t.Fatalf("retry deletion error = %v", err)
	}
	if devices.calls != 1 || len(tokens.tokens) != 0 {
		t.Fatalf("retry did not finish cleanup: device_calls=%d tokens=%d", devices.calls, len(tokens.tokens))
	}
}

func TestDeleteDeviceIsIdempotentAfterSuccessfulDeletion(t *testing.T) {
	events := []string{}
	service, _, devices := newLifecycleService(&events, []*devicetoken.DeviceToken{{ID: 1, PublicID: "token-1"}})
	d := &device.Device{ID: 42}

	if err := service.DeleteDevice(context.Background(), d); err != nil {
		t.Fatal(err)
	}
	firstCallCount := len(events)
	if err := service.DeleteDevice(context.Background(), d); err != nil {
		t.Fatalf("repeated deletion error = %v", err)
	}
	if len(events) != firstCallCount || devices.calls != 1 {
		t.Fatalf("repeated deletion repeated cleanup: events=%#v device_calls=%d", events, devices.calls)
	}
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func contains(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
