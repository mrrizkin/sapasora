package device

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
	"sapasora/internal/modules/permission"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type ownershipDeviceServiceStub struct {
	devicemodule.DeviceService
	globalLookupCalled bool
	scopedLookupCalled bool
	gotOwnerID         uint
	created            *devicemodule.Device
	deleted            bool
}

func (s *ownershipDeviceServiceStub) GetDeviceByPublicID(_ context.Context, _ string) (*devicemodule.Device, error) {
	s.globalLookupCalled = true
	return nil, nil
}

func (s *ownershipDeviceServiceStub) GetDeviceByPublicIDForUser(_ context.Context, _ string, ownerID uint) (*devicemodule.Device, error) {
	s.scopedLookupCalled = true
	s.gotOwnerID = ownerID
	return nil, gorm.ErrRecordNotFound
}

func (s *ownershipDeviceServiceStub) CreateDevice(_ context.Context, created *devicemodule.Device) error {
	s.created = created
	return nil
}

func newDeviceOwnershipTestApp(stub *ownershipDeviceServiceStub) *fiber.App {
	controller := &DeviceController{
		Controller:    &controllers.Controller{},
		deviceService: stub,
	}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Delete("/device/:id", func(ctx *fiber.Ctx) error {
		ctx.Locals("apikey", &apikey.APIKey{
			UserID:      7,
			Permissions: &permission.Permission{CanDeleteDevice: true, CanStoreDevice: true},
		})
		return controller.Destroy(ctx)
	})
	app.Put("/device/:id", func(ctx *fiber.Ctx) error {
		ctx.Locals("apikey", &apikey.APIKey{
			UserID:      7,
			Permissions: &permission.Permission{CanUpdateDevice: true},
		})
		return controller.Update(ctx)
	})
	app.Post("/device", func(ctx *fiber.Ctx) error {
		ctx.Locals("apikey", &apikey.APIKey{
			UserID:      7,
			Permissions: &permission.Permission{CanStoreDevice: true},
		})
		return controller.Store(ctx)
	})
	return app
}

func TestDeviceDestroyRejectsCrossOwnerWithoutEnumeration(t *testing.T) {
	stub := &ownershipDeviceServiceStub{}
	app := newDeviceOwnershipTestApp(stub)

	response, err := app.Test(httptest.NewRequest(http.MethodDelete, "/device/device-owned-by-8", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), "device-owned-by-8")
	require.True(t, stub.scopedLookupCalled)
	require.False(t, stub.globalLookupCalled)
	require.Equal(t, uint(7), stub.gotOwnerID)
}

func TestDeviceUpdateRejectsCrossOwnerWithoutEnumeration(t *testing.T) {
	stub := &ownershipDeviceServiceStub{}
	app := newDeviceOwnershipTestApp(stub)
	request := httptest.NewRequest(http.MethodPut, "/device/device-owned-by-8", bytes.NewBufferString(`{"name":"new name","user_id":"owner-8"}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), "device-owned-by-8")
	require.True(t, stub.scopedLookupCalled)
	require.False(t, stub.globalLookupCalled)
	require.Equal(t, uint(7), stub.gotOwnerID)
}

func TestDeviceStoreUsesAuthenticatedOwnerInsteadOfRequestUserID(t *testing.T) {
	stub := &ownershipDeviceServiceStub{}
	app := newDeviceOwnershipTestApp(stub)
	request := httptest.NewRequest(http.MethodPost, "/device", bytes.NewBufferString(`{"name":"device","type":"whatsapp","user_id":"owner-8"}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.NotNil(t, stub.created)
	require.Equal(t, uint(7), stub.created.UserID)
}
