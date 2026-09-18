// Package dashboard provides the Dashboard controllers
package dashboard

import (
	"sapasora/internal/app/http/controllers"

	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	*controllers.Controller
}

// NewDashboardController creates a new Dashboardcontrollers
// @wired:provide
func NewDashboardController(
	controller *controllers.Controller,
) *DashboardController {
	return &DashboardController{
		Controller: controller,
	}
}

func (c *DashboardController) Index(ctx *fiber.Ctx) error {
	return c.Inertia(ctx, "dashboard/index")
}
