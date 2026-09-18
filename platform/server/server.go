// Package server provides server for the application
package server

import (
	"fmt"
	"time"

	"sapasora/platform/config"
	"sapasora/platform/logger"
	"sapasora/platform/server/security"
	"sapasora/platform/session"
	"sapasora/platform/support/debug"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewServer),
)

type Server struct {
	app *fiber.App

	logger *logger.Logger
}

type ServerIn struct {
	fx.In

	Config config.Config

	Session *session.Session
	Logger  *logger.Logger

	CustomErrorHandler []fiber.ErrorHandler `group:"custom_error_handler"`

	Routers []Router `group:"router"`
}

func NewServer(lc fx.Lifecycle, in ServerIn) *Server {
	var errorHandler fiber.ErrorHandler
	for _, handler := range in.CustomErrorHandler {
		errorHandler = handler
		break
	}

	if errorHandler == nil {
		errorHandler = DefaultErrorHandler(in.Config)
	}

	app := fiber.New(fiber.Config{
		Prefork:               in.Config.GetBool("app.prefork", false),
		AppName:               in.Config.GetString("app.name", "application"),
		ReadTimeout:           in.Config.GetDuration("server.read_timeout", 15*time.Second),
		WriteTimeout:          in.Config.GetDuration("server.write_timeout", 30*time.Second),
		IdleTimeout:           in.Config.GetDuration("server.idle_timeout", 60*time.Second),
		BodyLimit:             in.Config.GetInt("server.body_limit", 16*1024*1024),
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})

	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy: security.DefaultCSP(),
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		HSTSPreloadEnabled:    true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		PermissionPolicy:      "camera=(), microphone=(), geolocation=(), payment=(), usb=()",
	}))

	app.Static("/", "public")
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: in.Logger.GetLogger(),
	}))
	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			if stackFrames, err := debug.StackTrace(); err == nil {
				c.Locals("stack_trace", stackFrames)
			}
			in.Logger.Error(fmt.Sprintf("panic: %v\n", e))
		},
	}))
	app.Use(idempotency.New())

	(Routers)(in.Routers).Register(app)

	return &Server{
		app:    app,
		logger: in.Logger,
	}
}

func (s *Server) App() *fiber.App {
	return s.app
}
