package devicetoken

import (
	"sapasora/internal/modules/devicetoken"

	"codeberg.org/mrrizkin/nihil"
)

type DeviceTokenResponse *devicetoken.DeviceToken // @name DeviceTokenController.DeviceTokenResponse

type DeviceTokenListResponse *devicetoken.Pagination[*devicetoken.DeviceToken] // @name DeviceTokenController.DeviceTokenListResponse

type DeviceTokenStoreRequest struct {
	ExpiredAt nihil.NilTime `json:"expired_at"`
	DeviceID  string        `json:"device_id" validate:"required,min=1,max=255"`
	UserID    string        `json:"user_id"`
} // @name DeviceTokenController.DeviceTokenStoreRequest

type DeviceTokenUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
} // @name DeviceTokenController.DeviceTokenUpdateRequest

type DeviceTokenUpdatePermissionsRequest struct {
	Permissions string `json:"permissions" validate:"omitempty,max=10000"`
} // @name DeviceTokenController.DeviceTokenUpdatePermissionsRequest
