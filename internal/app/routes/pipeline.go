package routes

import (
	"fmt"
	"time"

	"sapasora/platform/config"
	"sapasora/platform/session"
	"sapasora/platform/support/arr"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

type Web struct {
	config  config.Config
	session *session.Session
}

// NewWeb create a new pipeline for the web
// @wired:provide
func NewWeb(config config.Config, session *session.Session) *Web {
	return &Web{
		config:  config,
		session: session,
	}
}

func (s *Web) PipeLine() []fiber.Handler {
	csrfKey := s.config.GetString("security.csrf.key", "X-CSRF-Token")
	cookieName := s.config.GetString("security.csrf.cookie_name", "fiber_csrf_token")
	sameSite := s.config.GetString("security.csrf.same_site", "Lax")
	secure := s.config.GetBool("security.csrf.secure", false)
	useSession := s.config.GetBool("security.csrf.session", true)
	expiration := s.config.GetDuration("security.csrf.expiration", 3600*time.Second)

	csrfConfig := csrf.Config{
		// The header is the submitted token. The cookie is only the readable
		// double-submit source; it is never accepted as the request token.
		KeyLookup:         fmt.Sprintf("header:%s", csrfKey),
		CookieName:        cookieName,
		CookiePath:        "/",
		CookieSameSite:    sameSite,
		CookieSecure:      secure,
		CookieSessionOnly: false,
		CookieHTTPOnly:    false,
		SingleUseToken:    false,
		Expiration:        expiration,
		KeyGenerator:      utils.UUIDv4,
		ErrorHandler:      csrf.ConfigDefault.ErrorHandler,
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
	}
	if useSession {
		csrfConfig.Session = s.session.Store
	}

	return arr.List(
		csrf.New(csrfConfig),
		cors.New(),
	)
}

type API struct {
}

// NewAPI create a new pipeline for the API
// @wired:provide
func NewAPI() *API {
	return &API{}
}

func (s *API) PipeLine() []fiber.Handler {

	return arr.List(
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		}),
	)
}

type Channel struct {
}

// NewChannel create a new pipeline for the Channel
// @wired:provide
func NewChannel() *Channel {
	return &Channel{}
}

func (s *Channel) PipeLine() []fiber.Handler {
	return arr.List(
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		}),
		func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("websocket_allowed", true)
				return c.Next()
			}

			return fiber.ErrUpgradeRequired
		},
	)
}
