// Package swagger provides swagger for the application
package swagger

import (
	"sapasora/platform/config"
	swagger_view "sapasora/platform/swagger/view"
	"sapasora/platform/ui/view"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Swagger struct {
	config config.Config
	view   *view.View
}

type SwaggerIn struct {
	fx.In

	Config config.Config
	View   *view.View
}

func NewSwagger(in SwaggerIn) *Swagger {
	return &Swagger{
		config: in.Config,
		view:   in.View,
	}
}

func (s *Swagger) Register(router fiber.Router) {
	router.Get("/docs/api", s.SwaggerUI)
}

func (s *Swagger) SwaggerUI(c *fiber.Ctx) error {
	return s.view.Render(
		c,
		swagger_view.Swagger(s.config.GetString("app.swagger.path", "/docs/v3/openapi.json")),
	)
}
