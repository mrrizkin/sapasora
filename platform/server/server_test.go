package server

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	platformconfig "sapasora/platform/config"
	"sapasora/platform/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	cfg := platformconfig.NewConfigManager()
	cfg.Set("app.name", "test")
	cfg.Set("app.env", "production")
	cfg.Set("server.read_timeout", 2*time.Second)
	cfg.Set("server.write_timeout", 3*time.Second)
	cfg.Set("server.idle_timeout", 4*time.Second)
	cfg.Set("server.body_limit", 10)
	cfg.Set("app.log.console", true)
	cfg.Set("app.log.file", false)
	cfg.Set("app.log.level", "disable")

	return NewServer(nil, ServerIn{
		Config: cfg,
		Logger: logger.NewLogger(cfg),
	})
}

func TestNewServerAppliesLimitsAndTimeouts(t *testing.T) {
	srv := newTestServer(t)
	cfg := srv.App().Config()

	require.Equal(t, 2*time.Second, cfg.ReadTimeout)
	require.Equal(t, 3*time.Second, cfg.WriteTimeout)
	require.Equal(t, 4*time.Second, cfg.IdleTimeout)
	require.Equal(t, 10, cfg.BodyLimit)
}

func TestNewServerAddsSecurityHeadersWithoutChangingResponse(t *testing.T) {
	srv := newTestServer(t)
	srv.App().Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	resp, err := srv.App().Test(httptest.NewRequest("GET", "/health", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	require.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "strict-origin-when-cross-origin", resp.Header.Get("Referrer-Policy"))
	require.Contains(t, resp.Header.Get("Content-Security-Policy"), "default-src 'self'")
	require.Contains(t, resp.Header.Get("Content-Security-Policy"), "object-src 'none'")
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "{\"ok\":true}", string(body))
}

func TestNewServerRejectsOversizedRequestBody(t *testing.T) {
	srv := newTestServer(t)
	srv.App().Post("/body", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("POST", "/body", bytes.NewBufferString("01234567890"))
	resp, err := srv.App().Test(req)
	require.Error(t, err)
	if resp != nil {
		require.Equal(t, fiber.StatusRequestEntityTooLarge, resp.StatusCode)
	}
}
