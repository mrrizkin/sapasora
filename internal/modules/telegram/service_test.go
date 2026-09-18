package telegram

import (
	"context"
	"errors"
	"testing"

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
	if err := (&Client{}).Disconnect(context.Background()); err != nil {
		t.Fatalf("Disconnect() error = %v, want nil", err)
	}
}
