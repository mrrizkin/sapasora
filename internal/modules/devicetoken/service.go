package devicetoken

import (
	"context"
)

type DeviceTokenServiceImpl struct {
	repo DeviceTokenRepository
}

// NewDeviceTokenService creates a new Implementation of DeviceTokenService
// @wired:provide
func NewDeviceTokenService(repo DeviceTokenRepository) DeviceTokenService {
	return &DeviceTokenServiceImpl{
		repo: repo,
	}
}

func (s *DeviceTokenServiceImpl) ListDeviceToken(
	ctx context.Context,
	page, limit int,
) (*Pagination[*DeviceToken], error) {
	return s.repo.ListDeviceToken(ctx, page, limit)
}

func (s *DeviceTokenServiceImpl) CreateDeviceToken(
	ctx context.Context,
	devicetoken *DeviceToken,
) error {
	return s.repo.CreateDeviceToken(ctx, devicetoken)
}

func (s *DeviceTokenServiceImpl) GetDeviceToken(ctx context.Context, id int) (*DeviceToken, error) {
	return s.repo.GetDeviceToken(ctx, id)
}

func (s *DeviceTokenServiceImpl) GetDeviceTokenByPublicID(
	ctx context.Context,
	publicID string,
) (*DeviceToken, error) {
	return s.repo.GetDeviceTokenByPublicID(ctx, publicID)
}

func (s *DeviceTokenServiceImpl) GetDeviceTokenByPublicIDForUser(
	ctx context.Context,
	publicID string,
	userID uint,
) (*DeviceToken, error) {
	return s.repo.GetDeviceTokenByPublicIDForUser(ctx, publicID, userID)
}

func (s *DeviceTokenServiceImpl) GetDeviceTokenByToken(
	ctx context.Context,
	token string,
) (*DeviceToken, error) {
	return s.repo.GetDeviceTokenByToken(ctx, token)
}

func (s *DeviceTokenServiceImpl) GetDeviceTokenByDeviceID(
	ctx context.Context,
	deviceID string,
) (*DeviceToken, error) {
	return s.repo.GetDeviceTokenByDeviceID(ctx, deviceID)
}

func (s *DeviceTokenServiceImpl) UpdateDeviceToken(
	ctx context.Context,
	devicetoken *DeviceToken,
) error {
	return s.repo.UpdateDeviceToken(ctx, devicetoken)
}

func (s *DeviceTokenServiceImpl) DeleteDeviceToken(
	ctx context.Context,
	devicetoken *DeviceToken,
) error {
	return s.repo.DeleteDeviceToken(ctx, devicetoken)
}
