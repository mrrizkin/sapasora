package telegram

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zelenin/go-tdlib/client"
)

func TestActiveUsernameHandlesMissingTelegramUsername(t *testing.T) {
	if got := activeUsername(nil); got != "" {
		t.Fatalf("activeUsername(nil) = %q, want empty", got)
	}

	if got := activeUsername(&client.User{}); got != "" {
		t.Fatalf("activeUsername(empty) = %q, want empty", got)
	}
}

func TestCapabilityNotImplementedIsStable501(t *testing.T) {
	err := CapabilityNotImplemented("send_audio")
	var fiberErr *fiber.Error
	if !errors.As(err, &fiberErr) {
		t.Fatalf("error type = %T, want *fiber.Error", err)
	}
	if fiberErr.Code != fiber.StatusNotImplemented {
		t.Fatalf("status = %d, want %d", fiberErr.Code, fiber.StatusNotImplemented)
	}
	if fiberErr.Message != "telegram_capability_not_implemented:send_audio" {
		t.Fatalf("message = %q, want stable capability code", fiberErr.Message)
	}
}

func TestDisconnectNilTDLibClientIsSafe(t *testing.T) {
	c := &Client{}
	if err := c.Disconnect(context.Background()); err != nil {
		t.Fatalf("Disconnect() error = %v, want nil", err)
	}
	if got := c.State(); got != ClientStateDisconnected {
		t.Fatalf("state = %q, want %q", got, ClientStateDisconnected)
	}
}

func TestClientLifecycleStates(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	c := &Client{
		params:     &client.SetTdlibParametersRequest{},
		deviceInfo: &TDLibDeviceInfo{ID: "device"},
		newClient: func(client.AuthorizationStateHandler, client.Option) (*client.Client, error) {
			close(started)
			<-release
			return &client.Client{}, nil
		},
	}

	if got := c.State(); got != ClientStateDisconnected {
		t.Fatalf("initial state = %q, want %q", got, ClientStateDisconnected)
	}

	connectDone := make(chan error, 1)
	go func() {
		connectDone <- c.Connect(func(string, string) error { return nil }, time.Second)
	}()
	<-started
	if got := c.State(); got != ClientStateConnecting {
		t.Fatalf("connecting state = %q, want %q", got, ClientStateConnecting)
	}
	close(release)
	if err := <-connectDone; err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	if got := c.State(); got != ClientStateConnected || !c.IsConnected(context.Background()) {
		t.Fatalf("connected state = %q, connected = %v", got, c.IsConnected(context.Background()))
	}

	closed := make(chan struct{})
	c.SetClosedHandler(func() { close(closed) })
	c.EventHandler(&client.UpdateAuthorizationState{
		AuthorizationState: &client.AuthorizationStateClosed{},
	})
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("closed handler was not called")
	}
	if got := c.State(); got != ClientStateDisconnected || c.IsConnected(context.Background()) {
		t.Fatalf("closed state = %q, connected = %v", got, c.IsConnected(context.Background()))
	}
}

func TestClientConnectFailureIsObservableAndRetryable(t *testing.T) {
	attempts := 0
	c := &Client{
		params:     &client.SetTdlibParametersRequest{},
		deviceInfo: &TDLibDeviceInfo{ID: "device"},
		newClient: func(client.AuthorizationStateHandler, client.Option) (*client.Client, error) {
			attempts++
			if attempts == 1 {
				return nil, errors.New("tdlib unavailable")
			}
			return &client.Client{}, nil
		},
	}

	if err := c.Connect(func(string, string) error { return nil }, time.Second); err == nil {
		t.Fatal("Connect() error = nil, want provider error")
	}
	if got := c.State(); got != ClientStateError || c.IsConnected(context.Background()) {
		t.Fatalf("failed state = %q, connected = %v", got, c.IsConnected(context.Background()))
	}
	if err := c.Connect(func(string, string) error { return nil }, time.Second); err != nil {
		t.Fatalf("retry Connect() error = %v", err)
	}
	if got := c.State(); got != ClientStateConnected {
		t.Fatalf("retry state = %q, want %q", got, ClientStateConnected)
	}
}
