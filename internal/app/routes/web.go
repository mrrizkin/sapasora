package routes

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/app/http/controllers/account"
	"sapasora/internal/app/http/controllers/auth"
	"sapasora/internal/app/http/controllers/dashboard"
	"sapasora/internal/app/http/controllers/device"
	"sapasora/internal/app/http/controllers/role"
	"sapasora/internal/app/http/middleware"
	"sapasora/platform/server"
	"sapasora/platform/support/arr"
	"sapasora/platform/swagger"

	"github.com/gofiber/fiber/v2"
)

// WebRouter is the router for the web, usually everything that is not an API or will be consumed by this application only
// @wired:provide(group=router)
func WebRouter(
	web *Web,

	swagger *swagger.Swagger,

	inertia *middleware.InertiaMiddleware,
	protected *middleware.AuthenticationMiddleware,

	authController *auth.AuthController,
	accountController *account.AccountController,
	dashboardController *dashboard.DashboardController,
	deviceController *device.DeviceController,
	roleController *role.RoleController,
	welcomeController *controllers.WelcomeController,
) server.Router {
	return server.NewRouter("/", func(router fiber.Router) {
		server.Group(router, arr.List(inertia.Handle), func(r fiber.Router) {
			r.Get("/", func(c *fiber.Ctx) error {
				return c.Redirect("/dashboard")
			}).Name("index")

			r.Get("/dashboard", protected.Handle, dashboardController.Index).Name("dashboard")

			device := r.Group("/devices", protected.Handle).Name("device.")
			device.Get("/", deviceController.Index).Name("index")
			device.Get("/create", deviceController.Create).Name("create")
			device.Get("/:id/show", deviceController.Show).Name("show")
			device.Get("/:id/edit", deviceController.Edit).Name("edit")

			auth := r.Group("/auth").Name("auth.")
			auth.Get("/", authController.Index).Name("index")
			auth.Post("/login", authController.Login).Name("login")
			auth.Delete("/logout", authController.Logout).Name("logout")
		})

		swagger.Register(router)
	}, web.PipeLine()...)
}

func dummy(c *fiber.Ctx) error {
	return c.SendString("dummy")
}
