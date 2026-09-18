// Package controllers is the controller for the application
package controllers

import (
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/apikey"
	"sapasora/internal/modules/device"
	"sapasora/platform/logger"
	"sapasora/platform/session"
	"sapasora/platform/support/http"
	"sapasora/platform/ui/inertia"
	"sapasora/platform/ui/view"
	"sapasora/platform/validator"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	valid   *validator.Validator
	view    *view.View
	inertia *inertia.Inertia

	Log     *logger.Logger
	Session *session.Session
}

// NewController creates a new controller for the application
// @wired:provide
func NewController(
	log *logger.Logger,
	session *session.Session,
	validator *validator.Validator,
	view *view.View,
	inertia *inertia.Inertia,
) *Controller {
	return &Controller{
		valid:   validator,
		view:    view,
		inertia: inertia,

		Log:     log,
		Session: session,
	}
}

func (c *Controller) GetSubject(ctx *fiber.Ctx, subjects ...string) (any, error) {
	for _, s := range subjects {
		switch s {
		case "account":
			if subject, ok := ctx.Locals(s).(*account.Account); ok {
				return subject, nil
			}
		case "apikey":
			if subject, ok := ctx.Locals(s).(*apikey.APIKey); ok {
				return subject, nil
			}
		case "device":
			if subject, ok := ctx.Locals(s).(*device.Device); ok {
				return subject, nil
			}
		}
	}
	return nil, fiber.ErrForbidden
}

func (c *Controller) View(ctx *fiber.Ctx, component templ.Component) error {
	return c.view.Render(ctx, component)
}

func (c *Controller) Inertia(ctx *fiber.Ctx, component string, props ...fiber.Map) error {
	return c.inertia.Render(ctx, component, props...)
}

func (c *Controller) InertiaRedirect(ctx *fiber.Ctx, url string, status ...int) error {
	return c.inertia.Redirect(ctx, url, status...)
}

func (c *Controller) InertiaLocation(ctx *fiber.Ctx, url string, status ...int) error {
	return c.inertia.Location(ctx, url, status...)
}

func (c *Controller) InertiaBack(ctx *fiber.Ctx, status ...int) error {
	return c.inertia.Back(ctx, status...)
}

func (c *Controller) BodyParserValidate(ctx *fiber.Ctx, out any) error {
	return c.valid.ParseBodyAndValidate(ctx, out)
}

func (c *Controller) QueryParserValidate(ctx *fiber.Ctx, out any) error {
	return http.ParseQueryParams(ctx, out)
}

func (c *Controller) Validate(data any) []validator.ErrorResponse {
	return c.valid.Validate(data)
}

func (c *Controller) FormatValidationErrors(errs []validator.ErrorResponse) []string {
	return c.valid.Format(errs)
}

func (c *Controller) MustValidate(data any) error {
	return c.valid.MustValidate(data)
}
