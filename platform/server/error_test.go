package server

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	platformconfig "sapasora/platform/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestDefaultErrorHandlerSanitizesProductionErrors(t *testing.T) {
	cfg := platformconfig.NewConfigManager()
	cfg.Set("app.env", "production")

	app := fiber.New(fiber.Config{ErrorHandler: DefaultErrorHandler(cfg)})
	app.Get("/failure", func(*fiber.Ctx) error {
		return errors.New("database password=super-secret at postgres://user:pass@example.test/db")
	})

	req := httptest.NewRequest("GET", "/failure", nil)
	req.Header.Set("Accept", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), "super-secret")
	require.NotContains(t, string(body), "postgres://")
	require.NotContains(t, string(body), "database password")
	require.Contains(t, string(body), "Internal Server Error")
}

func TestDefaultErrorHandlerRedactsSensitiveClientError(t *testing.T) {
	cfg := platformconfig.NewConfigManager()
	cfg.Set("app.env", "production")

	app := fiber.New(fiber.Config{ErrorHandler: DefaultErrorHandler(cfg)})
	app.Get("/failure", func(*fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "invalid token=secret-value")
	})

	req := httptest.NewRequest("GET", "/failure", nil)
	req.Header.Set("Accept", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), "secret-value")
	require.Contains(t, string(body), "Bad Request")
}
