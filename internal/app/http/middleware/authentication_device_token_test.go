package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	devicemodule "sapasora/internal/modules/device"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type deviceTokenAuthServiceStub struct {
	devicemodule.DeviceService
	device   *devicemodule.Device
	gotToken string
}

func (s *deviceTokenAuthServiceStub) GetDeviceByToken(
	_ context.Context,
	token string,
) (*devicemodule.Device, error) {
	s.gotToken = token
	if s.device == nil {
		return nil, errors.New("device token not found")
	}
	return s.device, nil
}

func TestAuthenticationMiddlewareRejectsInvalidDeviceTokenWithoutEchoingIt(t *testing.T) {
	const token = "sk-dat-invalid-must-not-echo"
	service := &deviceTokenAuthServiceStub{}
	middleware := NewAuthenticationMiddleware(nil, nil, nil, nil, service)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/api/v1/device/by-token", middleware.Handle)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/device/by-token", nil)
	request.Header.Set("Authorization", token)
	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), token)
}

func TestAuthenticationMiddlewareResolvesDeviceFromAuthorizationHeader(t *testing.T) {
	const token = "sk-dat-header-only"
	resolved := &devicemodule.Device{PublicID: "device-1"}
	service := &deviceTokenAuthServiceStub{device: resolved}
	middleware := NewAuthenticationMiddleware(nil, nil, nil, nil, service)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get(
		"/api/v1/device/by-token",
		middleware.Handle,
		func(c *fiber.Ctx) error {
			require.Same(t, resolved, c.Locals("device"))
			require.NotContains(t, c.OriginalURL(), token)
			require.NotContains(t, c.Get("Referer"), token)
			return c.SendStatus(http.StatusNoContent)
		},
	)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/device/by-token", nil)
	request.Header.Set("Authorization", token)
	request.Header.Set("Referer", "https://client.example/device/by-token")
	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
	require.Equal(t, token, service.gotToken)
}
