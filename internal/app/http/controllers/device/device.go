// Package device provides the Device controllers
package device

import (
	"errors"
	"fmt"

	"sapasora/internal/app/http/controllers"
	"sapasora/internal/app/policies"
	"sapasora/internal/modules/account"
	devicemodule "sapasora/internal/modules/device"
	devicelifecycle "sapasora/internal/modules/devicelifecycle"
	"sapasora/internal/modules/permission"
	"sapasora/platform/satpam"
	"sapasora/platform/support/hash"
	"sapasora/platform/ui/inertia"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DeviceController struct {
	*controllers.Controller

	accountService        account.AccountService
	deviceService         devicemodule.DeviceService
	deviceDeletionService devicelifecycle.DeviceDeletionService
}

// NewDeviceController creates a new Devicecontrollers
// @wired:provide
func NewDeviceController(
	controller *controllers.Controller,
	accountService account.AccountService,
	deviceService devicemodule.DeviceService,
	deviceDeletionService devicelifecycle.DeviceDeletionService,
) *DeviceController {
	return &DeviceController{
		Controller:            controller,
		accountService:        accountService,
		deviceService:         deviceService,
		deviceDeletionService: deviceDeletionService,
	}
}

// Index show the index page
func (c *DeviceController) Index(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanListDevice{})
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	gate.AuthorizeAllPermissions(subject)

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	search := ctx.Query("search")

	deviceList, err := c.deviceService.ListDevice(ctx.Context(), search, page, limit)
	if err != nil {
		return err
	}

	return c.Inertia(ctx, "device/index", fiber.Map{
		"devices": DeviceListResponse(deviceList),
		"filters": fiber.Map{
			"search": search,
		},
		"pagination": fiber.Map{
			"pageIndex": page - 1,
			"pageSize":  limit,
		},
	})
}

// Create show the form to create a new device
func (c *DeviceController) Create(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanStoreDevice{})
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	gate.AuthorizeAllPermissions(subject)
	return c.Inertia(ctx, "device/create")
}

// Show display the device detail
func (c *DeviceController) Show(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	device, err := c.deviceService.GetDeviceByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanGetDevice{}).AddResource("device", device)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	gate.AuthorizeAllPermissions(subject)

	if device.User == nil {
		user, err := c.accountService.GetAccount(ctx.Context(), int(device.UserID))
		if err != nil {
			return err
		}
		device.User = user
	}

	return c.Inertia(ctx, "device/show", fiber.Map{
		"device": device,
	})
}

// Edit show the form to edit the device
func (c *DeviceController) Edit(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	device, err := c.deviceService.GetDeviceByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanUpdateDevice{}).AddResource("device", device)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	gate.AuthorizeAllPermissions(subject)

	return c.Inertia(ctx, "device/edit", fiber.Map{
		"device": device,
	})
}

// List godoc
// @Summary      List device
// @Description  List device
// @Tags         Device
// @Produce      json
// @Param        page    query  int    false "Page"
// @Param        limit   query  int    false "Limit"
// @Param        search  query  string false "Search"
// @Success      200 {object} DeviceListResponse
// @Router       /api/v1/device [get]
func (c *DeviceController) List(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanListDevice{})
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
	search := ctx.Query("search")

	deviceList, err := c.deviceService.ListDevice(ctx.Context(), search, page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(DeviceListResponse(deviceList))
}

// Get godoc
// @Summary      Get device
// @Description  Get device
// @Tags         Device
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} DeviceResponse
// @Router       /api/v1/device/{id} [get]
func (c *DeviceController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	device, err := c.deviceService.GetDeviceByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanGetDevice{}).AddResource("device", device)
	subject, err := c.GetSubject(ctx, "apikey", "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	return ctx.JSON(DeviceResponse(device))
}

// GetDeviceByToken godoc
// @Summary      Get device for the authenticated device token
// @Description  Returns the device resolved by the Authorization header. The device token is never accepted in the URL.
// @Tags         Device
// @Produce      json
// @Security     X-API-KEY
// @Success      200 {object} DeviceResponse
// @Failure      401 {object} map[string]string
// @Router       /api/v1/device/by-token [get]
func (c *DeviceController) GetDeviceByToken(ctx *fiber.Ctx) error {
	// AuthenticationMiddleware is the only component allowed to resolve a
	// device token. Do not read a token from path, query, or request body.
	device, ok := ctx.Locals("device").(*devicemodule.Device)
	if !ok || device == nil {
		return fiber.ErrUnauthorized
	}

	return ctx.JSON(DeviceResponse(device))
}

// Store godoc
// @Summary      Store device
// @Description  Store device
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        device  body  DeviceStoreRequest true "Body"
// @Success      200 {object} DeviceResponse
// @Router       /api/v1/device [post]
func (c *DeviceController) Store(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanStoreDevice{})
	subject, err := c.GetSubject(ctx, "apikey", "account")
	if err != nil {
		return err
	}
	gate.AuthorizeAllPermissions(subject)

	ownerID, err := c.GetOwnerID(ctx, "apikey", "account")
	if err != nil {
		return err
	}

	var payload DeviceStoreRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	if payload.Permissions == nil {
		payload.Permissions = &permission.Permission{
			CanCheckUserGateway:        true,
			CanConnectGateway:          true,
			CanDisconnectGateway:       true,
			CanGetAvatarGateway:        true,
			CanGetContactsGateway:      true,
			CanGetQRGateway:            true,
			CanGetStatusGateway:        true,
			CanGetUserGateway:          true,
			CanLogoutGateway:           true,
			CanSendAudioGateway:        true,
			CanSendButtonGateway:       true,
			CanSendChatPresenceGateway: true,
			CanSendContactGateway:      true,
			CanSendDocumentGateway:     true,
			CanSendImageGateway:        true,
			CanSendListGateway:         true,
			CanSendLocationGateway:     true,
			CanSendStickerGateway:      true,
			CanSendTextGateway:         true,
			CanSendVideoGateway:        true,
		}
	}

	device := devicemodule.Device{
		PublicID:    hash.NanoID(),
		Name:        payload.Name,
		Type:        payload.Type,
		Status:      devicemodule.DeviceStatusInactive,
		UserID:      ownerID,
		Events:      payload.Events,
		ExpiredAt:   payload.ExpiredAt,
		Webhook:     payload.Webhook,
		Permissions: payload.Permissions,
	}

	if err := c.deviceService.CreateDevice(ctx.Context(), &device); err != nil {
		return err
	}

	if inertia.IsInertiaRequest(ctx) {
		return c.InertiaRedirect(ctx, fmt.Sprintf("/devices/%s/show", device.PublicID))
	}

	return ctx.JSON(DeviceResponse(&device))
}

// Update godoc
// @Summary      Update device
// @Description  Update device
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        id  path  string true "ID"
// @Param        device  body  DeviceUpdateRequest true "Body"
// @Success      200 {object} DeviceResponse
// @Router       /api/v1/device/{id} [put]
func (c *DeviceController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload DeviceUpdateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	subject, err := c.GetSubject(ctx, "apikey", "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	ownerID, err := c.GetOwnerID(ctx, "apikey", "account")
	if err != nil {
		return err
	}

	device, err := c.deviceService.GetDeviceByPublicIDForUser(ctx.Context(), id, ownerID)
	if err != nil {
		return c.OwnerLookupError(err)
	}

	gate := satpam.New(&policies.CanUpdateDevice{}).AddResource("device", device)
	gate.AuthorizeAllPermissions(subject)

	if payload.Permissions == nil {
		payload.Permissions = &permission.Permission{
			CanCheckUserGateway:        true,
			CanConnectGateway:          true,
			CanDisconnectGateway:       true,
			CanGetAvatarGateway:        true,
			CanGetContactsGateway:      true,
			CanGetQRGateway:            true,
			CanGetStatusGateway:        true,
			CanGetUserGateway:          true,
			CanLogoutGateway:           true,
			CanSendAudioGateway:        true,
			CanSendButtonGateway:       true,
			CanSendChatPresenceGateway: true,
			CanSendContactGateway:      true,
			CanSendDocumentGateway:     true,
			CanSendImageGateway:        true,
			CanSendListGateway:         true,
			CanSendLocationGateway:     true,
			CanSendStickerGateway:      true,
			CanSendTextGateway:         true,
			CanSendVideoGateway:        true,
		}
	}

	device.Name = payload.Name
	device.Type = payload.Type
	device.Webhook = payload.Webhook
	device.Events = payload.Events
	device.ExpiredAt = payload.ExpiredAt
	device.Permissions = payload.Permissions

	if err := c.deviceService.UpdateDevice(ctx.Context(), device); err != nil {
		return err
	}

	return ctx.JSON(DeviceResponse(device))
}

// UpdateStatus godoc
// @Summary      Update device status
// @Description  Update device status
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        id  path  string true "ID"
// @Param        status  body  DeviceStatusRequest true "Body"
// @Success      200 {object} DeviceResponse
// @Router       /api/v1/device/{id}/status [put]
func (c *DeviceController) UpdateStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload DeviceStatusRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	subject, err := c.GetSubject(ctx, "apikey", "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	ownerID, err := c.GetOwnerID(ctx, "apikey", "account")
	if err != nil {
		return err
	}

	device, err := c.deviceService.GetDeviceByPublicIDForUser(ctx.Context(), id, ownerID)
	if err != nil {
		return c.OwnerLookupError(err)
	}

	gate := satpam.New(&policies.CanUpdateDevice{}).AddResource("device", device)
	gate.AuthorizeAllPermissions(subject)

	if payload.Status == device.Status {
		return ctx.JSON(DeviceResponse(device))
	}

	device.Status = payload.Status

	if err := c.deviceService.UpdateDevice(ctx.Context(), device); err != nil {
		return err
	}

	if inertia.IsInertiaRequest(ctx) {
		return c.InertiaRedirect(ctx, "device/show")
	}

	return ctx.JSON(DeviceResponse(device))
}

// Destroy godoc
// @Summary      Destroy device
// @Description  Destroy device
// @Tags         Device
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} DeviceResponse
// @Router       /api/v1/device/{id} [delete]
func (c *DeviceController) Destroy(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	subject, err := c.GetSubject(ctx, "apikey", "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	ownerID, err := c.GetOwnerID(ctx, "apikey", "account")
	if err != nil {
		return err
	}

	device, err := c.deviceService.GetDeviceByPublicIDForUser(ctx.Context(), id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if lookup, ok := c.deviceService.(devicemodule.DeletedDeviceOwnerScopedService); ok {
				deletedDevice, lookupErr := lookup.GetDeletedDeviceByPublicIDForUser(ctx.Context(), id, ownerID)
				if lookupErr == nil && deletedDevice != nil && deletedDevice.DeletedAt.Valid {
					gate := satpam.New(&policies.CanDeleteDevice{}).AddResource("device", deletedDevice)
					gate.AuthorizeAllPermissions(subject)
					return ctx.JSON(DeviceResponse(deletedDevice))
				}
			}
		}
		return c.OwnerLookupError(err)
	}

	gate := satpam.New(&policies.CanDeleteDevice{}).AddResource("device", device)
	gate.AuthorizeAllPermissions(subject)

	if c.deviceDeletionService != nil {
		if err := c.deviceDeletionService.DeleteDevice(ctx.Context(), device); err != nil {
			return err
		}
	} else if err := c.deviceService.DeleteDevice(ctx.Context(), device); err != nil {
		// Keep manually constructed controller test doubles compatible. Production
		// wiring always supplies the lifecycle coordinator above.
		return err
	}

	return ctx.JSON(DeviceResponse(device))
}
