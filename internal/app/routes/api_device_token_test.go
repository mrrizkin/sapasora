package routes

import (
	"testing"

	"sapasora/internal/app/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestAPIRouterUsesHeaderOnlyDeviceTokenEndpoint(t *testing.T) {
	apiRouter := APIRouter(
		&API{},
		&middleware.AuthenticationMiddleware{},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	apiRouter.Handle(app.Group(apiRouter.Path, apiRouter.Pipeline...))

	var paths []string
	for _, route := range app.GetRoutes(true) {
		if route.Method == fiber.MethodGet {
			paths = append(paths, route.Path)
		}
	}

	require.Contains(t, paths, "/api/v1/device/by-token")
	require.NotContains(t, paths, "/api/v1/device/:token/token")
}
