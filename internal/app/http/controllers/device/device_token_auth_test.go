package device

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"sapasora/internal/app/http/controllers"
	devicemodule "sapasora/internal/modules/device"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestGetDeviceByTokenUsesMiddlewareLocalOnly(t *testing.T) {
	const token = "sk-dat-must-not-appear"

	resolved := &devicemodule.Device{
		PublicID: "device-1",
		Name:     "resolved-device",
	}
	controller := &DeviceController{Controller: &controllers.Controller{}}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/api/v1/device/by-token", func(c *fiber.Ctx) error {
		c.Locals("device", resolved)
		return controller.GetDeviceByToken(c)
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/device/by-token?token="+url.QueryEscape("client-supplied-token"),
		nil,
	)
	request.Header.Set("Referer", "https://client.example/device/by-token")

	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), token)
	require.Contains(t, string(body), "resolved-device")
	require.NotContains(t, request.URL.String(), token)
	require.NotContains(t, request.Referer(), token)
}

func TestGetDeviceByTokenRequiresMiddlewareLocal(t *testing.T) {
	controller := &DeviceController{Controller: &controllers.Controller{}}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/api/v1/device/by-token", controller.GetDeviceByToken)

	response, err := app.Test(httptest.NewRequest(
		http.MethodGet,
		"/api/v1/device/by-token",
		nil,
	))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
