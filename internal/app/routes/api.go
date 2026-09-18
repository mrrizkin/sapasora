package routes

import (
	"sapasora/internal/app/http/controllers/account"
	"sapasora/internal/app/http/controllers/apikey"
	"sapasora/internal/app/http/controllers/device"
	"sapasora/internal/app/http/controllers/devicetoken"
	"sapasora/internal/app/http/controllers/gateway"
	"sapasora/internal/app/http/controllers/role"
	"sapasora/internal/app/http/middleware"
	"sapasora/platform/server"

	"github.com/gofiber/fiber/v2"
)

// APIRouter is the router for the API, usually everything that is just the api and will be consumed by other applications
// @wired:provide(group=router)
func APIRouter(
	api *API,

	protected *middleware.AuthenticationMiddleware,

	accountController *account.AccountController,
	apiKeyController *apikey.APIKeyController,
	deviceController *device.DeviceController,
	deviceTokenController *devicetoken.DeviceTokenController,
	gatewayController *gateway.GatewayController,
	roleController *role.RoleController,
) server.Router {
	return server.NewRouter("/api", func(r fiber.Router) {

		v1 := r.Group("/v1").Name("api.v1.")
		v1.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "ok"})
		}).Name("health")

		// APIKey
		apiKey := v1.Group("/api-key", protected.Handle).Name("api-key.")
		apiKey.Get("/", apiKeyController.List).Name("list")
		apiKey.Get("/:id", apiKeyController.Get).Name("get")
		apiKey.Post("/", apiKeyController.Store).Name("store")
		apiKey.Put("/:id", apiKeyController.Update).Name("update")
		apiKey.Delete("/:id", apiKeyController.Destroy).Name("destroy")

		// Account
		account := v1.Group("/account", protected.Handle).Name("account.")
		account.Get("/", accountController.List).Name("list")
		account.Get("/:id", accountController.Get).Name("get")
		account.Post("/", accountController.Store).Name("store")
		account.Put("/:id", accountController.Update).Name("update")
		account.Put("/:id/password", accountController.UpdatePassword).Name("update_password")
		account.Delete("/:id", accountController.Destroy).Name("destroy")

		// Device
		device := v1.Group("/device", protected.Handle).Name("device.")
		device.Get("/", deviceController.List).Name("list")
		// Device-token authentication is header-only. The controller reads the
		// device resolved by AuthenticationMiddleware from c.Locals("device").
		device.Get("/by-token", deviceController.GetDeviceByToken).Name("get_device_by_token")
		device.Get("/:id", deviceController.Get).Name("get")
		device.Post("/", deviceController.Store).Name("store")
		device.Put("/:id", deviceController.Update).Name("update")
		device.Put("/:id/status", deviceController.UpdateStatus).Name("update_status")
		device.Delete("/:id", deviceController.Destroy).Name("destroy")

		// DeviceToken
		deviceToken := v1.Group("/devicetoken", protected.Handle).Name("devicetoken.")
		deviceToken.Get("/", deviceTokenController.List).Name("list")
		deviceToken.Get("/:id", deviceTokenController.Get).Name("get")
		deviceToken.Post("/", deviceTokenController.Store).Name("store")
		deviceToken.Delete("/:id", deviceTokenController.Destroy).Name("destroy")

		// Gateway
		gateway := v1.Group("/gateway", protected.Handle).Name("gateway.")
		gateway.Post("/check-user", gatewayController.CheckUser).Name("check_user")
		gateway.Post("/connect", gatewayController.Connect).Name("connect")
		gateway.Post("/disconnect", gatewayController.Disconnect).Name("disconnect")
		gateway.Post("/get-avatar", gatewayController.GetAvatar).Name("get_avatar")
		gateway.Post("/get-contacts", gatewayController.GetContacts).Name("get_contacts")
		gateway.Get("/get-qr", gatewayController.GetQR).Name("get_qr")
		gateway.Get("/get-status", gatewayController.GetStatus).Name("get_status")
		gateway.Post("/get-user", gatewayController.GetUser).Name("get_user")
		gateway.Post("/logout", gatewayController.Logout).Name("logout")
		gateway.Post("/send-audio", gatewayController.SendAudio).Name("send_audio")
		gateway.Post("/send-button", gatewayController.SendButton).Name("send_button")
		gateway.Post("/send-chat-presence", gatewayController.SendChatPresence).
			Name("send_chat_presence")
		gateway.Post("/send-contact", gatewayController.SendContact).Name("send_contact")
		gateway.Post("/send-document", gatewayController.SendDocument).Name("send_document")
		gateway.Post("/send-image", gatewayController.SendImage).Name("send_image")
		gateway.Post("/send-list", gatewayController.SendList).Name("send_list")
		gateway.Post("/send-location", gatewayController.SendLocation).Name("send_location")
		gateway.Post("/send-sticker", gatewayController.SendSticker).Name("send_sticker")
		gateway.Post("/send-text", gatewayController.SendText).Name("send_text")
		gateway.Post("/send-video", gatewayController.SendVideo).Name("send_video")

		// Role
		role := v1.Group("/role", protected.Handle).Name("role.")
		role.Get("/", roleController.List).Name("list")
		role.Get("/:id", roleController.Get).Name("get")
		role.Post("/", roleController.Store).Name("store")
		role.Put("/:id", roleController.Update).Name("update")
		role.Put("/:id/permissions", roleController.UpdatePermissions).Name("update_permissions")
		role.Delete("/:id", roleController.Destroy).Name("destroy")

		// Fallback
		r.Get("/*", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotFound)
		})

	}, api.PipeLine()...)
}
