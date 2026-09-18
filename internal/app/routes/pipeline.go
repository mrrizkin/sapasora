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
	"github.com/gofiber/fiber/v2/middleware/helmet"
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
	cookieName := s.config.GetString("security.csrf.cookie_name", "fiber")
	sameSite := s.config.GetString("security.csrf.same_site", "Lax")
	secure := s.config.GetBool("security.csrf.secure", false)
	httpOnly := s.config.GetBool("security.csrf.http_only", true)
	expiration := s.config.GetDuration("security.csrf.expiration", 3600*time.Second)

	return arr.List(
		csrf.New(csrf.Config{
			KeyLookup:         fmt.Sprintf("cookie:%s", cookieName),
			CookieName:        cookieName,
			CookieSameSite:    sameSite,
			CookieSecure:      secure,
			CookieSessionOnly: true,
			CookieHTTPOnly:    httpOnly,
			SingleUseToken:    true,
			Expiration:        expiration,
			KeyGenerator:      utils.UUIDv4,
			ErrorHandler:      csrf.ConfigDefault.ErrorHandler,
			Extractor:         csrf.CsrfFromCookie(cookieName),
			Session:           s.session.Store,
			SessionKey:        "fiber.csrf.token",
			HandlerContextKey: "fiber.csrf.handler",
		}),
		cors.New(),
		helmet.New(),
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
			AllowHeaders: "Origin, Content-Type, Accept",
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
