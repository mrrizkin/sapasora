package devicetoken

import (
	"context"
)

type DeviceTokenService interface {
	ListDeviceToken(ctx context.Context, page, limit int) (*Pagination[*DeviceToken], error)
	CreateDeviceToken(ctx context.Context, devicetoken *DeviceToken) error
	GetDeviceToken(ctx context.Context, id int) (*DeviceToken, error)
	GetDeviceTokenByPublicID(ctx context.Context, publicID string) (*DeviceToken, error)
	GetDeviceTokenByPublicIDForUser(ctx context.Context, publicID string, userID uint) (*DeviceToken, error)
	GetDeviceTokenByToken(ctx context.Context, token string) (*DeviceToken, error)
	GetDeviceTokenByDeviceID(ctx context.Context, deviceID string) (*DeviceToken, error)
	ListDeviceTokensByDeviceID(ctx context.Context, deviceID uint) ([]*DeviceToken, error)
	UpdateDeviceToken(ctx context.Context, devicetoken *DeviceToken) error
	DeleteDeviceToken(ctx context.Context, devicetoken *DeviceToken) error
}

type DeviceTokenRepository interface {
	ListDeviceToken(ctx context.Context, page, limit int) (*Pagination[*DeviceToken], error)
	CreateDeviceToken(ctx context.Context, devicetoken *DeviceToken) error
	GetDeviceToken(ctx context.Context, id int) (*DeviceToken, error)
	GetDeviceTokenByPublicID(ctx context.Context, publicID string) (*DeviceToken, error)
	GetDeviceTokenByPublicIDForUser(ctx context.Context, publicID string, userID uint) (*DeviceToken, error)
	GetDeviceTokenByToken(ctx context.Context, token string) (*DeviceToken, error)
	GetDeviceTokenByDeviceID(ctx context.Context, deviceID string) (*DeviceToken, error)
	ListDeviceTokensByDeviceID(ctx context.Context, deviceID uint) ([]*DeviceToken, error)
	UpdateDeviceToken(ctx context.Context, devicetoken *DeviceToken) error
	DeleteDeviceToken(ctx context.Context, devicetoken *DeviceToken) error
}
