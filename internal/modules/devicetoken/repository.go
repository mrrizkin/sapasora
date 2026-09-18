package devicetoken

import (
	"context"

	"sapasora/platform/database"
)

type DeviceTokenRepositoryImpl struct {
	db *database.Database
}

// NewDeviceTokenRepository creates a new Implementation of DeviceTokenRepository
// @wired:provide
func NewDeviceTokenRepository(db *database.Database) DeviceTokenRepository {
	return &DeviceTokenRepositoryImpl{
		db: db,
	}
}

func (r *DeviceTokenRepositoryImpl) ListDeviceToken(
	ctx context.Context,
	page, limit int,
) (*Pagination[*DeviceToken], error) {

	qb := r.db.WithContext(ctx).
		Model(&DeviceToken{})

	var count int64
	total := qb.Count(&count)
	if total.Error != nil {
		return nil, total.Error
	}

	var devicetokens []*DeviceToken
	q := qb.Offset(page - 1).Limit(limit).Find(&devicetokens)
	if q.Error != nil {
		return nil, q.Error
	}

	return &Pagination[*DeviceToken]{
		Page:  page,
		Limit: limit,
		Total: count,
		Data:  devicetokens,
	}, nil
}

func (r *DeviceTokenRepositoryImpl) CreateDeviceToken(
	ctx context.Context,
	devicetoken *DeviceToken,
) error {
	if err := devicetoken.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(devicetoken).Error
}

func (r *DeviceTokenRepositoryImpl) GetDeviceToken(
	ctx context.Context,
	id int,
) (*DeviceToken, error) {
	var result DeviceToken
	q := r.db.WithContext(ctx).First(&result, id)
	return &result, q.Error
}

func (r *DeviceTokenRepositoryImpl) GetDeviceTokenByPublicID(
	ctx context.Context,
	publicID string,
) (*DeviceToken, error) {
	var result DeviceToken
	q := r.db.WithContext(ctx).Where("public_id = ?", publicID).First(&result)
	return &result, q.Error
}

func (r *DeviceTokenRepositoryImpl) GetDeviceTokenByPublicIDForUser(
	ctx context.Context,
	publicID string,
	userID uint,
) (*DeviceToken, error) {
	var result DeviceToken
	q := r.db.WithContext(ctx).
		Where("public_id = ? AND user_id = ?", publicID, userID).
		First(&result)
	return &result, q.Error
}

func (r *DeviceTokenRepositoryImpl) GetDeviceTokenByToken(
	ctx context.Context,
	token string,
) (*DeviceToken, error) {
	var result DeviceToken
	q := r.db.WithContext(ctx).Where("token = ?", token).First(&result)
	return &result, q.Error
}

func (r *DeviceTokenRepositoryImpl) GetDeviceTokenByDeviceID(
	ctx context.Context,
	deviceID string,
) (*DeviceToken, error) {
	var result DeviceToken
	q := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&result)
	return &result, q.Error
}

func (r *DeviceTokenRepositoryImpl) UpdateDeviceToken(
	ctx context.Context,
	devicetoken *DeviceToken,
) error {
	if err := devicetoken.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(devicetoken).Error
}

func (r *DeviceTokenRepositoryImpl) DeleteDeviceToken(
	ctx context.Context,
	devicetoken *DeviceToken,
) error {
	if err := devicetoken.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(devicetoken).Error
}
