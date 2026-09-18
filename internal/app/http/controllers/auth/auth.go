// Package auth provides the Auth controllers
package auth

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/role"
	"sapasora/platform/support/hash"
	"sapasora/platform/ui/inertia"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	*controllers.Controller

	accountService account.AccountService
	roleService    role.RoleService
}

// NewAuthController creates a new Authcontrollers
// @wired:provide
func NewAuthController(
	controller *controllers.Controller,
	accountService account.AccountService,
	roleService role.RoleService,
) *AuthController {
	return &AuthController{
		Controller:     controller,
		accountService: accountService,
		roleService:    roleService,
	}
}

func (c *AuthController) Index(ctx *fiber.Ctx) error {
	return c.Inertia(ctx, "auth/index")
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var payload LoginRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	account, err := c.accountService.GetAccountByUsername(ctx.Context(), payload.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "username or password is incorrect")
	}

	if ok := hash.Argon2Verify(payload.Password, account.HashedPassword); !ok {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "username or password is incorrect")
	}

	role, err := c.roleService.GetRole(ctx.Context(), int(account.RoleID))
	if err != nil {
		return err
	}

	account.Role = role

	session, err := c.Session.Get(ctx)
	if err != nil {
		return err
	}

	accountMarshal, err := json.Marshal(account)
	if err != nil {
		return err
	}

	session.Set("account", string(accountMarshal))
	if err := session.Save(); err != nil {
		return err
	}

	if inertia.IsInertiaRequest(ctx) {
		return ctx.Redirect("/dashboard")
	}

	return ctx.JSON((account))
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	session, err := c.Session.Get(ctx)
	if err != nil {
		return err
	}

	account := session.Get("account")
	if account == nil {
		return ctx.JSON((nil))
	}

	session.Delete("account")
	session.Destroy()

	return ctx.JSON(account)
}
