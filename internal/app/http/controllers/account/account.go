// Package account provides the Account controllers
package account

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/app/policies"
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/role"
	"sapasora/platform/satpam"
	"sapasora/platform/support/hash"

	"github.com/gofiber/fiber/v2"
)

type AccountController struct {
	*controllers.Controller

	accountService account.AccountService
	roleService    role.RoleService
}

// NewAccountController creates a new Accountcontrollers
// @wired:provide
func NewAccountController(
	controller *controllers.Controller,
	accountService account.AccountService,
	roleService role.RoleService,
) *AccountController {
	return &AccountController{
		Controller:     controller,
		accountService: accountService,
		roleService:    roleService,
	}
}

// List godoc
// @Summary      List account
// @Description  List account
// @Tags         Account
// @Produce      json
// @Param        page    query  int    false "Page"
// @Param        limit   query  int    false "Limit"
// @Param        search  query  string false "Search"
// @Success      200 {object} AccountListResponse
// @Security     Authorization
// @Router       /api/v1/account [get]
func (c *AccountController) List(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanListAccount{})
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrForbidden
	}

	gate.AuthorizeAllPermissions(subject)

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	search := ctx.Query("search")

	accountList, err := c.accountService.ListAccount(ctx.Context(), search, page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(AccountListResponse(accountList))
}

// Get godoc
// @Summary      Get account
// @Description  Get account
// @Tags         Account
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} AccountResponse
// @Security     Authorization
// @Router       /api/v1/account/{id} [get]
func (c *AccountController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	account, err := c.accountService.GetAccountByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanGetAccount{}).AddResource("account", account)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	return ctx.JSON(AccountResponse(account))
}

// Store godoc
// @Summary      Store account
// @Description  Store account
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        account  body  AccountStoreRequest true "Body"
// @Success      200 {object} AccountResponse
// @Security     Authorization
// @Router       /api/v1/account [post]
func (c *AccountController) Store(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanStoreAccount{})
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload AccountStoreRequest
	if err := c.BodyParserValidate(ctx, &payload); err != nil {
		return err
	}

	role, err := c.roleService.GetRoleByPublicID(ctx.Context(), payload.RoleID)
	if err != nil {
		return err
	}

	hashedPassword, err := hash.Argon2(payload.Password)
	if err != nil {
		return err
	}

	account := account.Account{
		PublicID:       hash.NanoID(),
		Name:           payload.Name,
		Username:       payload.Username,
		HashedPassword: hashedPassword,
		RoleID:         role.ID,
	}

	if err := c.accountService.CreateAccount(ctx.Context(), &account); err != nil {
		return err
	}

	return ctx.JSON(AccountResponse(&account))
}

// Update godoc
// @Summary      Update account
// @Description  Update account
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        id  path  string true "ID"
// @Param        account  body  AccountUpdateRequest true "Body"
// @Success      200 {object} AccountResponse
// @Security     Authorization
// @Router       /api/v1/account/{id} [put]
func (c *AccountController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload AccountUpdateRequest
	if err := c.BodyParserValidate(ctx, &payload); err != nil {
		return err
	}

	account, err := c.accountService.GetAccountByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanUpdateAccount{}).AddResource("account", account)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	role, err := c.roleService.GetRoleByPublicID(ctx.Context(), payload.RoleID)
	if err != nil {
		return err
	}

	account.Name = payload.Name
	account.Username = payload.Username
	account.RoleID = role.ID

	if err := c.accountService.UpdateAccount(ctx.Context(), account); err != nil {
		return err
	}

	return ctx.JSON(AccountResponse(account))
}

// UpdatePassword godoc
// @Summary      Update account password
// @Description  Update account password
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        id  path  string true "ID"
// @Param        account  body  AccountUpdatePasswordRequest true "Body"
// @Success      200 {object} AccountResponse
// @Security     Authorization
// @Router       /api/v1/account/{id}/password [put]
func (c *AccountController) UpdatePassword(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload AccountUpdatePasswordRequest
	if err := c.BodyParserValidate(ctx, &payload); err != nil {
		return err
	}

	account, err := c.accountService.GetAccountByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanUpdatePasswordAccount{}).AddResource("account", account)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	if ok := hash.Argon2Verify(payload.OldPassword, account.HashedPassword); !ok {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Incorrect password")
	}

	hashedPassword, err := hash.Argon2(payload.NewPassword)
	if err != nil {
		return err
	}

	account.HashedPassword = hashedPassword

	if err := c.accountService.UpdateAccount(ctx.Context(), account); err != nil {
		return err
	}

	return ctx.JSON(AccountResponse(account))
}

// Destroy godoc
// @Summary      Destroy account
// @Description  Destroy account
// @Tags         Account
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} AccountResponse
// @Security     Authorization
// @Router       /api/v1/account/{id} [delete]
func (c *AccountController) Destroy(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	account, err := c.accountService.GetAccountByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanDeleteAccount{}).AddResource("account", account)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	if err := c.accountService.DeleteAccount(ctx.Context(), account); err != nil {
		return err
	}

	return ctx.JSON(AccountResponse(account))
}
