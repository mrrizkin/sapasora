// Package gateway provides the Gateway controllers
package gateway

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/app/policies"
	"sapasora/internal/modules/device"
	"sapasora/internal/modules/gateway"
	"sapasora/platform/satpam"
	"sapasora/platform/ui/inertia"
	"sapasora/resources/views"

	"github.com/gofiber/fiber/v2"
)

type GatewayController struct {
	*controllers.Controller

	deviceService  device.DeviceService
	gatewayService gateway.GatewayService
}

// NewGatewayController creates a new Gatewaycontrollers
// @wired:provide
func NewGatewayController(
	controller *controllers.Controller,
	deviceService device.DeviceService,
	gatewayService gateway.GatewayService,
) *GatewayController {
	return &GatewayController{
		Controller:     controller,
		deviceService:  deviceService,
		gatewayService: gatewayService,
	}
}

// CheckUser godoc
// @Summary      Check user
// @Description  Check user
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.CheckUserRequest true "Body"
// @Success      200 {object} gateway.CheckUserResponse
// @Router       /api/v1/gateway/check-user [post]
func (c *GatewayController) CheckUser(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanCheckUserGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.CheckUserRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.CheckUser(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// Connect godoc
// @Summary      Connect
// @Description  Connect
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.ConnectRequest true "Body"
// @Success      200 {object} GatewayResponse
// @Router       /api/v1/gateway/connect [post]
func (c *GatewayController) Connect(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanConnectGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.ConnectRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	err = c.gatewayService.Connect(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(GatewayResponse{
		Status:  "pending",
		Message: "Request accepted",
	})
}

// Disconnect godoc
// @Summary      Disconnect
// @Description  Disconnect
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Success      200 {object} GatewayResponse
// @Router       /api/v1/gateway/disconnect [post]
func (c *GatewayController) Disconnect(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanDisconnectGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	err = c.gatewayService.Disconnect(ctx.Context(), device)
	if err != nil {
		return err
	}

	return ctx.JSON(GatewayResponse{
		Status:  "pending",
		Message: "Request accepted",
	})
}

// GetAvatar godoc
// @Summary      Get avatar
// @Description  Get avatar
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.GetAvatarRequest true "Body"
// @Success      200 {object} gateway.GetAvatarResponse
// @Router       /api/v1/gateway/get-avatar [post]
func (c *GatewayController) GetAvatar(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanGetAvatarGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.GetAvatarRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.GetAvatar(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// GetContacts godoc
// @Summary      Get contacts
// @Description  Get contacts
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Success      200 {object} gateway.GetContactsResponse
// @Router       /api/v1/gateway/get-contacts [post]
func (c *GatewayController) GetContacts(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanGetContactsGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.GetContacts(ctx.Context(), device)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// GetQR godoc
// @Summary      Get QR
// @Description  Get QR
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Success      200 {object} gateway.GetQRResponse
// @Router       /api/v1/gateway/get-qr [get]
func (c *GatewayController) GetQR(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanGetQRGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	if !device.QRCode.Valid {
		return fiber.NewError(fiber.StatusInternalServerError, "QR Code not found")
	}

	if device.QRCode.String == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "QR Code not found")
	}

	if (ctx.Get("X-Requested-With") != "XMLHttpRequest" &&
		ctx.Get("Accept") != "application/json") ||
		inertia.IsInertiaRequest(ctx) {

		return c.View(ctx, views.QRCode(device.QRCode.String))
	}

	return ctx.JSON(gateway.GetQRResponse{
		QRCode: device.QRCode,
	})
}

// GetStatus godoc
// @Summary      Get status
// @Description  Get status
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Success      200 {object} gateway.GetStatusResponse
// @Router       /api/v1/gateway/get-status [get]
func (c *GatewayController) GetStatus(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanGetStatusGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.GetStatus(ctx.Context(), device)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// GetUser godoc
// @Summary      Get user
// @Description  Get user
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.GetUserRequest true "Body"
// @Success      200 {object} gateway.GetUserResponse
// @Router       /api/v1/gateway/get-user [post]
func (c *GatewayController) GetUser(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanGetUserGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.GetUserRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.GetUser(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// Logout godoc
// @Summary      Logout
// @Description  Logout
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Success      200 {object} GatewayResponse
// @Router       /api/v1/gateway/logout [post]
func (c *GatewayController) Logout(ctx *fiber.Ctx) error {
	device, ok := ctx.Locals("device").(*device.Device)
	gate := satpam.New(&policies.CanLogoutGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	err = c.gatewayService.Logout(ctx.Context(), device)
	if err != nil {
		return err
	}

	return ctx.JSON(GatewayResponse{
		Status:  "pending",
		Message: "Request accepted",
	})
}

// SendAudio godoc
// @Summary      Send audio
// @Description  Send audio
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendAudioRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-audio [post]
func (c *GatewayController) SendAudio(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendAudioGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendAudioRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendAudio(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendButton godoc
// @Summary      Send button
// @Description  Send button
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendButtonTextRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-button [post]
func (c *GatewayController) SendButton(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendButtonGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendButtonTextRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendButton(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendChatPresence godoc
// @Summary      Send chat presence
// @Description  Send chat presence
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.ChatPresenceRequest true "Body"
// @Success      200 {object} GatewayResponse
// @Router       /api/v1/gateway/send-chat-presence [post]
func (c *GatewayController) SendChatPresence(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendChatPresenceGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.ChatPresenceRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	err = c.gatewayService.SendChatPresence(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(GatewayResponse{
		Status:  "pending",
		Message: "Request accepted",
	})
}

// SendContact godoc
// @Summary      Send contact
// @Description  Send contact
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendContactRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-contact [post]
func (c *GatewayController) SendContact(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendContactGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendContactRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendContact(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendDocument godoc
// @Summary      Send document
// @Description  Send document
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendDocumentRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-document [post]
func (c *GatewayController) SendDocument(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendDocumentGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendDocumentRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendDocument(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendImage godoc
// @Summary      Send image
// @Description  Send image
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendImageRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-image [post]
func (c *GatewayController) SendImage(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendImageGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendImageRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendImage(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendList godoc
// @Summary      Send list
// @Description  Send list
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendListRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-list [post]
func (c *GatewayController) SendList(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendListGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendListRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendList(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendLocation godoc
// @Summary      Send location
// @Description  Send location
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendLocationRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-location [post]
func (c *GatewayController) SendLocation(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendLocationGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendLocationRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendLocation(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendSticker godoc
// @Summary      Send sticker
// @Description  Send sticker
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendStickerRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-sticker [post]
func (c *GatewayController) SendSticker(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendStickerGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendStickerRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendSticker(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendText godoc
// @Summary      Send text
// @Description  Send text
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendTextRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-text [post]
func (c *GatewayController) SendText(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendTextGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendTextRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendText(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

// SendVideo godoc
// @Summary      Send video
// @Description  Send video
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        device  body  gateway.SendVideoRequest true "Body"
// @Success      200 {object} gateway.SendResponse
// @Router       /api/v1/gateway/send-video [post]
func (c *GatewayController) SendVideo(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanSendVideoGateway{})
	subject, err := c.GetSubject(ctx, "device")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload gateway.SendVideoRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	device, ok := ctx.Locals("device").(*device.Device)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Device not found")
	}

	response, err := c.gatewayService.SendVideo(ctx.Context(), device, &payload)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
