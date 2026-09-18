package session

import (
	"context"
	"sapasora/platform/config"
	"sapasora/platform/database"
	"sapasora/platform/session/provider"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/fx"
)

type SessionProvider interface {
	Setup() (fiber.Storage, error)
}

type Session struct {
	*session.Store
	storage fiber.Storage
}

func NewSession(
	lc fx.Lifecycle,

	dbConfig *database.Config,
	sessionConfig config.Config,
) (*Session, error) {
	var driver SessionProvider

	switch sessionConfig.GetString("session.driver") {
	case "database":
		driver = provider.NewDatabase(dbConfig)
	case "file":
		driver = provider.NewFile()
	case "redis", "valkey", "memory":
		driver = provider.NewMemory(sessionConfig)
	default:
		driver = provider.NewFile()
	}

	storage, err := driver.Setup()
	if err != nil {
		return nil, err
	}

	cookieName := sessionConfig.GetString("session.cookie_name")
	sameSite := sessionConfig.GetString("session.same_site")
	secure := sessionConfig.GetBool("session.secure")
	httpOnly := sessionConfig.GetBool("session.http_only")
	store := Session{
		Store: session.New(session.Config{
			Storage:        storage,
			Expiration:     24 * time.Hour,
			KeyLookup:      fmt.Sprintf("cookie:%s_session_key", cookieName),
			CookieHTTPOnly: httpOnly,
			CookieSecure:   secure,
			CookieSameSite: sameSite,
		}),
		storage: storage,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return store.Stop()
		},
	})

	return &store, nil
}

func (s *Session) Stop() error {
	return s.storage.Close()
}
