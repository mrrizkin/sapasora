package middleware

import (
	"encoding/json"
	"strings"

	"sapasora/internal/modules/account"
	"sapasora/internal/modules/apikey"
	"sapasora/internal/modules/device"
	"sapasora/platform/logger"
	"sapasora/platform/session"
	"sapasora/platform/ui/inertia"

	"github.com/gofiber/fiber/v2"
)

type AuthenticationMiddleware struct {
	log           *logger.Logger
	inertia       *inertia.Inertia
	session       *session.Session
	apikeyService apikey.APIKeyService
	deviceService device.DeviceService
}

// NewAuthenticationMiddleware creates a new middleware for the authentication
// @wired:provide
func NewAuthenticationMiddleware(
	log *logger.Logger,
	inertia *inertia.Inertia,
	session *session.Session,
	apikeyService apikey.APIKeyService,
	deviceService device.DeviceService,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		log:           log,
		inertia:       inertia,
		session:       session,
		apikeyService: apikeyService,
		deviceService: deviceService,
	}
}

func (m *AuthenticationMiddleware) Handle(c *fiber.Ctx) error {
	if apiKey := c.Get("Authorization"); apiKey != "" {
		if strings.HasPrefix(apiKey, "sk-dat-") { // secret key device access token
			device, err := m.deviceService.GetDeviceByToken(c.Context(), apiKey)
			if err != nil {
				// Do not log the credential or distinguish revoked/expired/missing
				// tokens; this is an audit signal without enumeration leakage.
				if m.log != nil {
					m.log.Warn(
						"Device token authentication failed",
						"auth_method", "device_token",
						"reason", "invalid_or_inactive_or_expired",
					)
				}
				return fiber.ErrUnauthorized
			}

			c.Locals("device", device)
			return c.Next()
		}

		if strings.HasPrefix(apiKey, "sk-dak-") { // secret key sapasora access key
			key, err := m.apikeyService.GetAPIKeyByKey(c.Context(), apiKey)
			if err != nil {
				// Never include the supplied credential or an error that may
				// contain it in logs.
				if m.log != nil {
					m.log.Error(
						"API key authentication failed",
						"auth_method", "api_key",
						"reason", "invalid_or_inactive",
					)
				}
				return fiber.ErrUnauthorized
			}

			c.Locals("apikey", key)
			return c.Next()
		}
	}

	if session, err := m.session.Get(c); err == nil {
		if accountJSON, ok := session.Get("account").(string); ok {
			var account account.Account
			if err := json.Unmarshal([]byte(accountJSON), &account); err == nil {
				c.Locals("account", &account)
				return c.Next()
			}
		}
	}

	if inertia.IsInertiaRequest(c) {
		return m.inertia.Redirect(c, "/auth")
	}

	return fiber.ErrUnauthorized
}
