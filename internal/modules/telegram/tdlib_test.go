package telegram

import (
	"context"
	"errors"
	"testing"
	"time"

	"sapasora/internal/modules/device"

	"github.com/zelenin/go-tdlib/client"
)

type lifecycleDeviceServiceStub struct {
	jid           string
	status        device.DeviceStatus
	jidCalls      int
	connectedCall int
}

func (s *lifecycleDeviceServiceStub) SetDeviceJID(_ context.Context, d *device.Device, jid string) error {
	s.jid = jid
	d.Jid.String = jid
	d.Jid.Valid = true
	s.jidCalls++
	return nil
}

func (s *lifecycleDeviceServiceStub) SetDeviceStatusConnected(_ context.Context, d *device.Device) error {
	s.status = device.DeviceStatusConnected
	d.Status = device.DeviceStatusConnected
	s.connectedCall++
	return nil
}

func TestSyncConnectedDevicePersistsTelegramIdentityAndStatus(t *testing.T) {
	service := &lifecycleDeviceServiceStub{}
	telegramClient := &Client{
		getMe: func(context.Context) (*client.User, error) {
			return &client.User{Id: 987654321}, nil
		},
	}
	deviceInfo := &device.Device{PublicID: "device"}

	if err := syncConnectedDevice(context.Background(), service, deviceInfo, telegramClient); err != nil {
		t.Fatalf("syncConnectedDevice() error = %v", err)
	}
	if service.jid != "987654321" || deviceInfo.Jid.String != "987654321" {
		t.Fatalf("stored JID = %q, want Telegram account ID", service.jid)
	}
	if service.status != device.DeviceStatusConnected || deviceInfo.Status != device.DeviceStatusConnected {
		t.Fatalf("stored status = %q, want connected", service.status)
	}
	if service.jidCalls != 1 || service.connectedCall != 1 {
		t.Fatalf("calls = jid:%d connected:%d, want one each", service.jidCalls, service.connectedCall)
	}
}

func TestSyncConnectedDeviceRejectsEmptyIdentity(t *testing.T) {
	service := &lifecycleDeviceServiceStub{}
	telegramClient := &Client{
		getMe: func(context.Context) (*client.User, error) {
			return &client.User{}, nil
		},
	}

	if err := syncConnectedDevice(context.Background(), service, &device.Device{}, telegramClient); err == nil {
		t.Fatal("syncConnectedDevice() error = nil, want empty identity error")
	}
	if service.jidCalls != 0 || service.connectedCall != 0 {
		t.Fatal("device service was called for an empty identity")
	}
}

func TestTDLibAsyncConnectionReportsErrorToSupervisor(t *testing.T) {
	tdlib := &TDLib{}
	wantErr := errors.New("provider unavailable")
	deviceInfo := &device.Device{Name: "test"}
	if err := tdlib.beginConnect(); err != nil {
		t.Fatalf("beginConnect() error = %v", err)
	}

	done := make(chan struct{})
	go func() {
		tdlib.superviseAsyncConnect(deviceInfo, func() error { return wantErr })
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("async connection supervisor did not finish")
	}
	if !errors.Is(tdlib.LastStartupError(), wantErr) {
		t.Fatalf("last startup error = %v, want %v", tdlib.LastStartupError(), wantErr)
	}
}

func TestTDLibAsyncConnectionConvertsPanicToError(t *testing.T) {
	tdlib := &TDLib{}
	if err := tdlib.beginConnect(); err != nil {
		t.Fatalf("beginConnect() error = %v", err)
	}

	done := make(chan struct{})
	go func() {
		tdlib.superviseAsyncConnect(&device.Device{}, func() error {
			panic("provider secret should not escape")
		})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("async connection supervisor did not recover panic")
	}
	if !errors.Is(tdlib.LastStartupError(), ErrAsyncConnectionPanic) {
		t.Fatalf("last startup error = %v, want %v", tdlib.LastStartupError(), ErrAsyncConnectionPanic)
	}
}

func TestTDLibStopWaitsForAsyncConnections(t *testing.T) {
	tdlib := &TDLib{clientStore: NewStore[*Client]()}
	if err := tdlib.beginConnect(); err != nil {
		t.Fatalf("beginConnect() error = %v", err)
	}

	stopped := make(chan struct{})
	go func() {
		_ = tdlib.Stop(context.Background())
		close(stopped)
	}()

	select {
	case <-stopped:
		t.Fatal("Stop returned before the connection goroutine finished")
	case <-time.After(10 * time.Millisecond):
	}

	tdlib.endConnect()
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("Stop did not wait for the connection goroutine")
	}

	if err := tdlib.beginConnect(); !errors.Is(err, ErrProviderStopping) {
		t.Fatalf("beginConnect() after Stop = %v, want %v", err, ErrProviderStopping)
	}
}
