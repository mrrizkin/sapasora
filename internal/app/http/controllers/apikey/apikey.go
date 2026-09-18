// Package apikey provides the APIKey controllers
package apikey

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/app/policies"
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/apikey"
	"sapasora/platform/satpam"
	"sapasora/platform/support/hash"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type APIKeyController struct {
	*controllers.Controller

	accountService account.AccountService
	apikeyService  apikey.APIKeyService
}

// NewAPIKeyController creates a new APIKeycontrollers
// @wired:provide
func NewAPIKeyController(
	controller *controllers.Controller,
	accountService account.AccountService,
	apikeyService apikey.APIKeyService,
) *APIKeyController {
	return &APIKeyController{
		Controller:     controller,
		accountService: accountService,
		apikeyService:  apikeyService,
	}
}

// List godoc
// @Summary      List apikey
// @Description  List apikey
// @Tags         APIKey
// @Produce      json
// @Param        page    query  int    false "Page"
// @Param        limit   query  int    false "Limit"
// @Success      200 {object} APIKeyListResponse
// @Router       /api/v1/api-key [get]
func (c *APIKeyController) List(ctx *fiber.Ctx) error {
    gate := satpam.New(&policies.CanListAPIKey{})
    subject, err := c.GetSubject(ctx, "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }

    gate.AuthorizeAllPermissions(subject)

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	apikeyList, err := c.apikeyService.ListAPIKey(ctx.Context(), page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(APIKeyListResponse(apikeyList))
}

// Get godoc
// @Summary      Get apikey
// @Description  Get apikey
// @Tags         APIKey
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} APIKeyResponse
// @Router       /api/v1/api-key/{id} [get]
func (c *APIKeyController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	apikey, err := c.apikeyService.GetAPIKeyByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

    gate := satpam.New(&policies.CanGetAPIKey{}).AddResource("apikey", apikey)
    subject, err := c.GetSubject(ctx, "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }

    gate.AuthorizeAllPermissions(subject)

	return ctx.JSON(APIKeyResponse(apikey))
}

// Store godoc
// @Summary      Store apikey
// @Description  Store apikey
// @Tags         APIKey
// @Accept       json
// @Produce      json
// @Param        apikey  body  APIKeyStoreRequest true "Body"
// @Success      200 {object} APIKeyResponse
// @Router       /api/v1/api-key [post]
func (c *APIKeyController) Store(ctx *fiber.Ctx) error {
    gate := satpam.New(&policies.CanStoreAPIKey{})
    subject, err := c.GetSubject(ctx, "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }
    gate.AuthorizeAllPermissions(subject)

	var payload APIKeyStoreRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	user, err := c.accountService.GetAccountByPublicID(ctx.Context(), payload.UserID)
	if err != nil {
		return err
	}

	apikey := apikey.APIKey{
		PublicID: hash.NanoID(),
		Name:     payload.Name,
		Key: fmt.Sprintf(
			"sk-dak-%s",
			hash.GenerateNanoID(
				"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
				42,
			),
		),
		Status:      apikey.APIKeyStatusActive,
		Permissions: &payload.Permissions,
		UserID:      user.ID,
	}

	if err := c.apikeyService.CreateAPIKey(ctx.Context(), &apikey); err != nil {
		return err
	}

	return ctx.JSON(APIKeyResponse(&apikey))
}

// Update godoc
// @Summary      Update apikey
// @Description  Update apikey
// @Tags         APIKey
// @Accept       json
// @Produce      json
// @Param        id  path  string true "ID"
// @Param        apikey  body  APIKeyUpdateRequest true "Body"
// @Success      200 {object} APIKeyResponse
// @Router       /api/v1/api-key/{id} [put]
func (c *APIKeyController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload APIKeyUpdateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	apikey, err := c.apikeyService.GetAPIKeyByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

    gate := satpam.New(&policies.CanUpdateAPIKey{}).AddResource("apikey", apikey)
    subject, err := c.GetSubject(ctx, "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }

    gate.AuthorizeAllPermissions(subject)

	apikey.Name = payload.Name
	apikey.Permissions = &payload.Permissions

	if err := ctx.BodyParser(apikey); err != nil {
		return err
	}

	if err := c.apikeyService.UpdateAPIKey(ctx.Context(), apikey); err != nil {
		return err
	}

	return ctx.JSON(APIKeyResponse(apikey))
}

// Destroy godoc
// @Summary      Destroy apikey
// @Description  Destroy apikey
// @Tags         APIKey
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} APIKeyResponse
// @Router       /api/v1/api-key/{id} [delete]
func (c *APIKeyController) Destroy(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	apikey, err := c.apikeyService.GetAPIKeyByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

    gate := satpam.New(&policies.CanDeleteAPIKey{}).AddResource("apikey", apikey)
    subject, err := c.GetSubject(ctx, "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }

    gate.AuthorizeAllPermissions(subject)

	if err := c.apikeyService.DeleteAPIKey(ctx.Context(), apikey); err != nil {
		return err
	}

	return ctx.JSON(APIKeyResponse(apikey))
}
