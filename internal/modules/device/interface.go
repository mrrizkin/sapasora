package device

import (
	"context"
)

type DeviceService interface {
	ListDevice(ctx context.Context, search string, page, limit int) (*Pagination[*Device], error)
	CreateDevice(ctx context.Context, device *Device) error
	GetDevice(ctx context.Context, id int) (*Device, error)
	GetDeviceByPublicID(ctx context.Context, publicID string) (*Device, error)
	GetDeviceByToken(ctx context.Context, token string) (*Device, error)
	UpdateDevice(ctx context.Context, device *Device) error
	DeleteDevice(ctx context.Context, device *Device) error

	GetAllWhatsappDevices(ctx context.Context) ([]*Device, error)
	GetAllTelegramDevices(ctx context.Context) ([]*Device, error)
	SetDeviceQRCode(ctx context.Context, device *Device, qrCode string) error
	SetDeviceQRCodeByPublicID(ctx context.Context, publicID string, qrCode string) error
	SetDeviceJID(ctx context.Context, device *Device, jid string) error
	SetDeviceJIDByPublicID(ctx context.Context, publicID string, jid string) error
	SetDeviceStatusConnected(ctx context.Context, device *Device) error
	SetDeviceStatusConnectedByPublicID(ctx context.Context, publicID string) error
	SetDeviceStatusDisconnected(ctx context.Context, device *Device) error
	SetDeviceStatusDisconnectedByPublicID(ctx context.Context, publicID string) error
}

type DeviceRepository interface {
	ListDevice(ctx context.Context, search string, page, limit int) (*Pagination[*Device], error)
	CreateDevice(ctx context.Context, device *Device) error
	GetDevice(ctx context.Context, id int) (*Device, error)
	GetDeviceByPublicID(ctx context.Context, publicID string) (*Device, error)
	GetDeviceByToken(ctx context.Context, token string) (*Device, error)
	UpdateDevice(ctx context.Context, device *Device) error
	DeleteDevice(ctx context.Context, device *Device) error

	GetAllWhatsappDevices(ctx context.Context) ([]*Device, error)
	GetAllTelegramDevices(ctx context.Context) ([]*Device, error)
}
