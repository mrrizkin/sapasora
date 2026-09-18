package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type WelcomeController struct {
	*Controller
}

type WelcomeControllerIn struct {
	fx.In

	Controller *Controller
}

// NewWelcomeController creates a new controller for the welcome page
// @wired:provide
func NewWelcomeController(
	controller *Controller,
) *WelcomeController {
	return &WelcomeController{
		Controller: controller,
	}
}

func (c *WelcomeController) Welcome(ctx *fiber.Ctx) error {
	return c.Inertia(ctx, "welcome")
}
