// Package devicelifecycle coordinates destructive device lifecycle operations.
package devicelifecycle

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"sapasora/internal/modules/device"
	"sapasora/internal/modules/devicetoken"
)

var ErrDeviceDeletion = errors.New("device deletion failed")

// Error intentionally exposes only the cleanup stage. The underlying provider
// or repository error remains available to errors.Is/errors.As callers without
// leaking credentials or provider internals through an HTTP error message.
type Error struct {
	stage string
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s during %s", ErrDeviceDeletion, e.stage)
}

func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Is(target error) bool {
	return target == ErrDeviceDeletion || errors.Is(e.cause, target)
}

func deletionError(stage string, err error) error {
	if err == nil {
		return nil
	}
	return &Error{stage: stage, cause: err}
}

// DeviceDeletionService performs cleanup in a deliberately strict order:
// disconnect provider, revoke/delete every device token, then soft-delete the
// device. A failed step leaves the device present so a later call can retry.
type DeviceDeletionService interface {
	DeleteDevice(ctx context.Context, d *device.Device) error
}

// ProviderDisconnector is deliberately narrower than GatewayService so this
// coordinator does not pull provider implementations into unit tests.
type ProviderDisconnector interface {
	Disconnect(ctx context.Context, d *device.Device) error
}

type Service struct {
	devices device.DeviceService
	tokens  devicetoken.DeviceTokenService
	gateway ProviderDisconnector
}

// NewDeviceDeletionService creates the T0.8 device deletion coordinator.
// @wired:provide
func NewDeviceDeletionService(
	devices device.DeviceService,
	tokens devicetoken.DeviceTokenService,
	gatewayService ProviderDisconnector,
) DeviceDeletionService {
	return &Service{
		devices: devices,
		tokens:  tokens,
		gateway: gatewayService,
	}
}

func (s *Service) DeleteDevice(ctx context.Context, d *device.Device) error {
	if d == nil {
		return deletionError("validation", errors.New("nil device"))
	}
	if d.DeletedAt.Valid {
		return nil
	}

	if err := s.gateway.Disconnect(ctx, d); err != nil {
		return deletionError("provider disconnect", err)
	}

	tokens, err := s.tokens.ListDeviceTokensByDeviceID(ctx, d.ID)
	if err != nil {
		return deletionError("credential lookup", err)
	}
	for _, token := range tokens {
		if token == nil {
			return deletionError("credential cleanup", errors.New("nil device credential"))
		}
		if err := s.tokens.DeleteDeviceToken(ctx, token); err != nil {
			return deletionError("credential cleanup", err)
		}
	}

	if err := ctx.Err(); err != nil {
		return deletionError("device delete", err)
	}
	if err := s.devices.DeleteDevice(ctx, d); err != nil {
		return deletionError("device delete", err)
	}

	// GORM normally sets DeletedAt during Delete. Set it here as well so a
	// repeated call using the same object is a no-op even with a test/fake repo.
	if !d.DeletedAt.Valid {
		d.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	}
	return nil
}
