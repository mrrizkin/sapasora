package devicetoken

import (
	"sapasora/internal/modules/devicetoken"

	"codeberg.org/mrrizkin/nihil"
)

type DeviceTokenResponse *devicetoken.DeviceToken // @name DeviceTokenController.DeviceTokenResponse

type DeviceTokenListResponse *devicetoken.Pagination[*devicetoken.DeviceToken] // @name DeviceTokenController.DeviceTokenListResponse

type DeviceTokenStoreRequest struct {
	ExpiredAt nihil.NilTime `json:"expired_at"`
	DeviceID  string        `json:"device_id"`
	UserID    string        `json:"user_id"`
} // @name DeviceTokenController.DeviceTokenStoreRequest

type DeviceTokenUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
} // @name DeviceTokenController.DeviceTokenUpdateRequest

type DeviceTokenUpdatePermissionsRequest struct {
	Permissions string `json:"permissions"`
} // @name DeviceTokenController.DeviceTokenUpdatePermissionsRequest
