package device

import (
	"sapasora/internal/modules/device"
	"sapasora/internal/modules/permission"

	"codeberg.org/mrrizkin/nihil"
)

type DeviceResponse *device.Device // @name DeviceController.DeviceResponse

type DeviceListResponse *device.Pagination[*device.Device] // @name DeviceController.DeviceListResponse

type DeviceStoreRequest struct {
	Name        string                 `json:"name"`
	Type        device.DeviceType      `json:"type"`
	Webhook     nihil.NilString        `json:"webhook"`
	Events      nihil.NilString        `json:"events"`
	ExpiredAt   nihil.NilTime          `json:"expired_at"`
	UserID      string                 `json:"user_id"`
	Permissions *permission.Permission `json:"permission"`
} // @name DeviceController.DeviceStoreRequest

type DeviceUpdateRequest struct {
	Name        string                 `json:"name"`
	Type        device.DeviceType      `json:"type"`
	Webhook     nihil.NilString        `json:"webhook"`
	Events      nihil.NilString        `json:"events"`
	ExpiredAt   nihil.NilTime          `json:"expired_at"`
	UserID      string                 `json:"user_id"`
	Permissions *permission.Permission `json:"permission"`
} // @name DeviceController.DeviceUpdateRequest

type DeviceStatusRequest struct {
	Status device.DeviceStatus `json:"status"`
} // @name DeviceController.DeviceStatusRequest
