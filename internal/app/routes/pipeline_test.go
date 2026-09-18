package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appconfig "sapasora/platform/config"
	appsession "sapasora/platform/session"

	"github.com/gofiber/fiber/v2"
	fibersession "github.com/gofiber/fiber/v2/middleware/session"
	fiberstorage "github.com/gofiber/storage/memory/v2"
	"github.com/stretchr/testify/require"
)

const testCSRFHeader = "X-Test-CSRF"

func newCSRFTestApp(t *testing.T, expiration time.Duration) *fiber.App {
	t.Helper()

	cfg := appconfig.NewConfigManager()
	cfg.Set("security.csrf.key", testCSRFHeader)
	cfg.Set("security.csrf.cookie_name", "test_csrf")
	cfg.Set("security.csrf.same_site", "Lax")
	cfg.Set("security.csrf.secure", false)
	cfg.Set("security.csrf.session", true)
	cfg.Set("security.csrf.expiration", expiration)

	store := fibersession.New(fibersession.Config{Storage: fiberstorage.New()})
	web := NewWeb(cfg, &appsession.Session{Store: store})

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	for _, handler := range web.PipeLine() {
		app.Use(handler)
	}
	app.Get("/form", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Post("/mutate", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	return app
}

func csrfBootstrap(t *testing.T, app *fiber.App) (string, string) {
	t.Helper()

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/form", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)

	cookies := response.Cookies()
	var cookieHeader []string
	var csrfToken string
	for _, cookie := range cookies {
		cookieHeader = append(cookieHeader, cookie.Name+"="+cookie.Value)
		if cookie.Name == "test_csrf" {
			csrfToken = cookie.Value
			require.Equal(t, "/", cookie.Path)
			require.False(t, cookie.HttpOnly, "the CSRF cookie must be readable by same-origin JavaScript")
		}
	}
	require.NotEmpty(t, csrfToken)

	return strings.Join(cookieHeader, "; "), csrfToken
}

func csrfRequest(t *testing.T, app *fiber.App, cookies, token string) *http.Response {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, "/mutate", nil)
	request.Header.Set("Cookie", cookies)
	if token != "" {
		request.Header.Set(testCSRFHeader, token)
	}

	response, err := app.Test(request)
	require.NoError(t, err)
	return response
}

func TestWebCSRFValidHeaderAndCookie(t *testing.T) {
	app := newCSRFTestApp(t, time.Hour)
	cookies, token := csrfBootstrap(t, app)

	response := csrfRequest(t, app, cookies, token)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
}

func TestWebCSRFMissingToken(t *testing.T) {
	app := newCSRFTestApp(t, time.Hour)
	cookies, _ := csrfBootstrap(t, app)

	response := csrfRequest(t, app, cookies, "")
	require.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestWebCSRFWrongToken(t *testing.T) {
	app := newCSRFTestApp(t, time.Hour)
	cookies, _ := csrfBootstrap(t, app)

	response := csrfRequest(t, app, cookies, "wrong-token")
	require.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestWebCSRFExpiredToken(t *testing.T) {
	app := newCSRFTestApp(t, time.Second)
	cookies, token := csrfBootstrap(t, app)
	time.Sleep(1100 * time.Millisecond)

	response := csrfRequest(t, app, cookies, token)
	require.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestAPIPipelineDoesNotRequireCSRF(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	api := NewAPI()
	for _, handler := range api.PipeLine() {
		app.Use(handler)
	}
	app.Post("/signed-webhook", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNoContent)
	})

	response, err := app.Test(httptest.NewRequest(http.MethodPost, "/signed-webhook", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
}
