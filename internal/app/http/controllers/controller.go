// Package controllers is the controller for the application
package controllers

import (
	"errors"

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
	"gorm.io/gorm"
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
		subject := ctx.Locals(s)
		if subject == nil {
			continue
		}

		switch subject.(type) {
		case *account.Account:
			return subject, nil
		case *apikey.APIKey:
			return subject, nil
		case *device.Device:
			return subject, nil
		}
	}
	return nil, fiber.ErrForbidden
}

// GetOwnerID returns the authenticated account ID from an account or API-key subject.
// Resource mutation handlers must use this value instead of owner fields from the request.
func (c *Controller) GetOwnerID(ctx *fiber.Ctx, subjects ...string) (uint, error) {
	subject, err := c.GetSubject(ctx, subjects...)
	if err != nil {
		return 0, err
	}

	switch subject := subject.(type) {
	case *account.Account:
		if subject.ID == 0 {
			return 0, fiber.ErrUnauthorized
		}
		return subject.ID, nil
	case *apikey.APIKey:
		if subject.UserID == 0 {
			return 0, fiber.ErrUnauthorized
		}
		return subject.UserID, nil
	default:
		return 0, fiber.ErrUnauthorized
	}
}

// OwnerLookupError prevents cross-owner resource lookups from exposing whether
// a public ID exists. Database failures still propagate normally.
func (c *Controller) OwnerLookupError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrForbidden
	}
	return err
}

func (c *Controller) View(ctx *fiber.Ctx, component templ.Component) error {
	return c.view.Render(ctx, component)
}

func (c *Controller) Inertia(ctx *fiber.Ctx, component string, props ...fiber.Map) error {
	// gonertia.Render only consumes the first Props map in the variadic slice
	// (see firstOr[Props](props, nil) in gonertia's response.go), so any
	// additional per-request data (like the authenticated user) must be
	// merged into a single map rather than appended as a separate element.
	merged := fiber.Map{}
	if acc, ok := ctx.Locals("account").(*account.Account); ok && acc != nil {
		merged["auth"] = fiber.Map{
			"user": fiber.Map{
				"id":       acc.PublicID,
				"name":     acc.Name,
				"username": acc.Username,
			},
		}
	}
	for _, p := range props {
		for k, v := range p {
			merged[k] = v
		}
	}
	return c.inertia.Render(ctx, component, merged)
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

func (c *Controller) validation() *validator.Validator {
	if c.valid == nil {
		// Keep lightweight controller tests safe when they construct a controller
		// without the production dependency graph.
		c.valid = validator.NewValidator()
	}
	return c.valid
}

func (c *Controller) BodyParserValidate(ctx *fiber.Ctx, out any) error {
	return c.validation().ParseBodyAndValidate(ctx, out)
}

func (c *Controller) QueryParserValidate(ctx *fiber.Ctx, out any) error {
	return http.ParseQueryParams(ctx, out)
}

func (c *Controller) ParseListQuery(ctx *fiber.Ctx) (http.ListQueryParams, error) {
	return http.ParseListQueryParams(ctx)
}

func (c *Controller) PublicIDParam(ctx *fiber.Ctx) (string, error) {
	id := ctx.Params("id")
	if err := http.ValidatePublicID(id); err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, "invalid id path parameter")
	}
	return id, nil
}

func (c *Controller) Validate(data any) []validator.ErrorResponse {
	return c.validation().Validate(data)
}

func (c *Controller) FormatValidationErrors(errs []validator.ErrorResponse) []string {
	return c.validation().Format(errs)
}

func (c *Controller) MustValidate(data any) error {
	return c.validation().MustValidate(data)
}
