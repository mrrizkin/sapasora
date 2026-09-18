package telegram

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrDeviceNotFound           = errors.New("device not found")
	ErrNoAvatarFound            = errors.New("no avatar found")
	ErrNotConnected             = errors.New("not connected")
	ErrAlreadyLoggedIn          = errors.New("already logged in")
	ErrNotLoggedIn              = errors.New("not logged in")
	ErrNoSession                = errors.New("no session")
	ErrClientAlreadyConnected   = errors.New("client already connected")
	ErrFailedToConnect          = errors.New("failed to connect")
	ErrInvalidPhoneNumber       = errors.New("invalid phone number")
	ErrEmptyBody                = errors.New("body cannot be empty")
	ErrMissingStanzaID          = errors.New("missing stanza id in contextinfo")
	ErrMissingParticipant       = errors.New("missing participant in contextinfo")
	ErrCapabilityNotImplemented = errors.New("telegram_capability_not_implemented")
)

// CapabilityNotImplemented returns a stable 501 error for Telegram features
// that are intentionally outside the current release scope.
func CapabilityNotImplemented(capability string) error {
	return fiber.NewError(fiber.StatusNotImplemented, fmt.Sprintf("%s:%s", ErrCapabilityNotImplemented, capability))
}
