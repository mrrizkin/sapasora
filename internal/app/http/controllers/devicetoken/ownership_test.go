package devicetoken

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"sapasora/internal/app/http/controllers"
	"sapasora/internal/modules/apikey"
	devicemodule "sapasora/internal/modules/device"
	tokenmodule "sapasora/internal/modules/devicetoken"
	"sapasora/internal/modules/permission"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type ownershipDeviceLookupStub struct {
	devicemodule.DeviceService
	gotOwnerID uint
	err        error
}

func (s *ownershipDeviceLookupStub) GetDeviceByPublicIDForUser(_ context.Context, _ string, ownerID uint) (*devicemodule.Device, error) {
	s.gotOwnerID = ownerID
	if s.err != nil {
		return nil, s.err
	}
	return &devicemodule.Device{ID: 11, PublicID: "device-7", UserID: ownerID}, nil
}

type ownershipTokenServiceStub struct {
	tokenmodule.DeviceTokenService
	globalLookupCalled bool
	scopedLookupCalled bool
	gotOwnerID         uint
	created            *tokenmodule.DeviceToken
}

func (s *ownershipTokenServiceStub) GetDeviceTokenByPublicID(_ context.Context, _ string) (*tokenmodule.DeviceToken, error) {
	s.globalLookupCalled = true
	return nil, nil
}

func (s *ownershipTokenServiceStub) GetDeviceTokenByPublicIDForUser(_ context.Context, _ string, ownerID uint) (*tokenmodule.DeviceToken, error) {
	s.scopedLookupCalled = true
	s.gotOwnerID = ownerID
	return nil, gorm.ErrRecordNotFound
}

func (s *ownershipTokenServiceStub) CreateDeviceToken(_ context.Context, created *tokenmodule.DeviceToken) error {
	s.created = created
	return nil
}

func newDeviceTokenOwnershipTestApp(deviceService devicemodule.DeviceService, tokenService tokenmodule.DeviceTokenService) *fiber.App {
	controller := &DeviceTokenController{
		Controller:         &controllers.Controller{},
		deviceService:      deviceService,
		devicetokenService: tokenService,
	}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Delete("/devicetoken/:id", func(ctx *fiber.Ctx) error {
		ctx.Locals("apikey", &apikey.APIKey{
			UserID:      7,
			Permissions: &permission.Permission{CanDeleteDeviceToken: true},
		})
		return controller.Destroy(ctx)
	})
	app.Post("/devicetoken", func(ctx *fiber.Ctx) error {
		ctx.Locals("apikey", &apikey.APIKey{
			UserID:      7,
			Permissions: &permission.Permission{CanStoreDeviceToken: true},
		})
		return controller.Store(ctx)
	})
	return app
}

func TestDeviceTokenDestroyRejectsCrossOwnerWithoutEnumeration(t *testing.T) {
	deviceService := &ownershipDeviceLookupStub{}
	tokenService := &ownershipTokenServiceStub{}
	app := newDeviceTokenOwnershipTestApp(deviceService, tokenService)

	response, err := app.Test(httptest.NewRequest(http.MethodDelete, "/devicetoken/token-owned-by-8", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), "token-owned-by-8")
	require.True(t, tokenService.scopedLookupCalled)
	require.False(t, tokenService.globalLookupCalled)
	require.Equal(t, uint(7), tokenService.gotOwnerID)
}

func TestDeviceTokenStoreRejectsCrossOwnerDeviceWithoutEnumeration(t *testing.T) {
	deviceService := &ownershipDeviceLookupStub{err: gorm.ErrRecordNotFound}
	tokenService := &ownershipTokenServiceStub{}
	app := newDeviceTokenOwnershipTestApp(deviceService, tokenService)
	request := httptest.NewRequest(http.MethodPost, "/devicetoken", bytes.NewBufferString(`{"device_id":"device-owned-by-8","user_id":"owner-8"}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), "device-owned-by-8")
	require.Nil(t, tokenService.created)
	require.Equal(t, uint(7), deviceService.gotOwnerID)
}

func TestDeviceTokenStoreUsesAuthenticatedOwnerAndScopedDevice(t *testing.T) {
	deviceService := &ownershipDeviceLookupStub{}
	tokenService := &ownershipTokenServiceStub{}
	app := newDeviceTokenOwnershipTestApp(deviceService, tokenService)
	request := httptest.NewRequest(http.MethodPost, "/devicetoken", bytes.NewBufferString(`{"device_id":"device-7","user_id":"owner-8"}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.NotNil(t, tokenService.created)
	require.Equal(t, uint(7), deviceService.gotOwnerID)
	require.Equal(t, uint(7), tokenService.created.UserID)
	require.Equal(t, uint(11), tokenService.created.DeviceID)
}
