// Package devicetoken provides the DeviceToken controllers
package devicetoken

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/app/policies"
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/device"
	"sapasora/internal/modules/devicetoken"
	"sapasora/platform/satpam"
	"sapasora/platform/support/hash"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type DeviceTokenController struct {
	*controllers.Controller

	accountService     account.AccountService
	deviceService      device.DeviceService
	devicetokenService devicetoken.DeviceTokenService
}

// NewDeviceTokenController creates a new DeviceTokencontrollers
// @wired:provide
func NewDeviceTokenController(
	controller *controllers.Controller,
	accountService account.AccountService,
	deviceService device.DeviceService,
	devicetokenService devicetoken.DeviceTokenService,
) *DeviceTokenController {
	return &DeviceTokenController{
		Controller:         controller,
		accountService:     accountService,
		deviceService:      deviceService,
		devicetokenService: devicetokenService,
	}
}

// List godoc
// @Summary      List devicetoken
// @Description  List devicetoken
// @Tags         DeviceToken
// @Produce      json
// @Param        page    query  int    false "Page"
// @Param        limit   query  int    false "Limit"
// @Success      200 {object} DeviceTokenListResponse
// @Router       /api/v1/devicetoken [get]
func (c *DeviceTokenController) List(ctx *fiber.Ctx) error {
    gate := satpam.New(&policies.CanListDeviceToken{})
    subject, err := c.GetSubject(ctx, "apikey", "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrForbidden
    }
    gate.AuthorizeAllPermissions(subject)

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	devicetokenList, err := c.devicetokenService.ListDeviceToken(ctx.Context(), page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(DeviceTokenListResponse(devicetokenList))
}

// Get godoc
// @Summary      Get devicetoken
// @Description  Get devicetoken
// @Tags         DeviceToken
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} DeviceTokenResponse
// @Router       /api/v1/devicetoken/{id} [get]
func (c *DeviceTokenController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	devicetoken, err := c.devicetokenService.GetDeviceTokenByPublicID(ctx.Context(), id)

    gate := satpam.New(&policies.CanGetDeviceToken{}).AddResource("devicetoken", devicetoken)
    subject, err := c.GetSubject(ctx, "apikey", "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }

    gate.AuthorizeAllPermissions(subject)

	if err != nil {
		return err
	}

	return ctx.JSON(DeviceTokenResponse(devicetoken))
}

// Store godoc
// @Summary      Store devicetoken
// @Description  Store devicetoken
// @Tags         DeviceToken
// @Accept       json
// @Produce      json
// @Param        devicetoken  body  DeviceTokenStoreRequest true "Body"
// @Success      200 {object} DeviceTokenResponse
// @Router       /api/v1/devicetoken [post]
func (c *DeviceTokenController) Store(ctx *fiber.Ctx) error {
    gate := satpam.New(&policies.CanStoreDeviceToken{})
    subject, err := c.GetSubject(ctx, "apikey", "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }
    gate.AuthorizeAllPermissions(subject)

	var payload DeviceTokenStoreRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, err := c.deviceService.GetDeviceByPublicID(ctx.Context(), payload.DeviceID)
	if err != nil {
		return err
	}

	user, err := c.accountService.GetAccountByPublicID(ctx.Context(), payload.UserID)
	if err != nil {
		return err
	}

	devicetoken := devicetoken.DeviceToken{
		PublicID: hash.NanoID(),
		Token: fmt.Sprintf(
			"sk-dat-%s",
			hash.GenerateNanoID(
				"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
				42,
			),
		),
		Status:    devicetoken.DeviceTokenStatusActive,
		ExpiredAt: payload.ExpiredAt,
		DeviceID:  device.ID,
		UserID:    user.ID,
	}

	if err := c.devicetokenService.CreateDeviceToken(ctx.Context(), &devicetoken); err != nil {
		return err
	}

	return ctx.JSON(DeviceTokenResponse(&devicetoken))
}

// Destroy godoc
// @Summary      Destroy devicetoken
// @Description  Destroy devicetoken
// @Tags         DeviceToken
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} DeviceTokenResponse
// @Router       /api/v1/devicetoken/{id} [delete]
func (c *DeviceTokenController) Destroy(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	devicetoken, err := c.devicetokenService.GetDeviceTokenByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

    gate := satpam.New(&policies.CanDeleteDeviceToken{}).AddResource("devicetoken", devicetoken)
    subject, err := c.GetSubject(ctx, "apikey", "account")
    if err != nil {
        return err
    }
    if subject == nil {
        return fiber.ErrUnauthorized
    }

    gate.AuthorizeAllPermissions(subject)

	if err := c.devicetokenService.DeleteDeviceToken(ctx.Context(), devicetoken); err != nil {
		return err
	}

	return ctx.JSON(DeviceTokenResponse(devicetoken))
}
