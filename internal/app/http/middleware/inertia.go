package middleware

import (
	"sapasora/platform/session"
	"sapasora/platform/ui/inertia"

	"github.com/gofiber/fiber/v2"
)

type InertiaMiddleware struct {
	inertia *inertia.Inertia
}

// NewInertiaMiddleware creates a new middleware for the inertia
// @wired:provide
func NewInertiaMiddleware(
	session *session.Session,
	inertia *inertia.Inertia,
) *InertiaMiddleware {
	return &InertiaMiddleware{
		inertia: inertia,
	}
}

func (m *InertiaMiddleware) Handle(c *fiber.Ctx) error {
	return m.inertia.Middleware(c)
}
