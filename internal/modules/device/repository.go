package device

import (
	"context"
	"strings"
	"time"

	"sapasora/internal/modules/devicetoken"
	"sapasora/platform/database"
	"sapasora/platform/support/sql"
)

type DeviceRepositoryImpl struct {
	db *database.Database
}

// NewDeviceRepository creates a new Implementation of DeviceRepository
// @wired:provide
func NewDeviceRepository(db *database.Database) DeviceRepository {
	return &DeviceRepositoryImpl{
		db: db,
	}
}

func (r *DeviceRepositoryImpl) ListDevice(
	ctx context.Context,
	search string,
	page, limit int,
) (*Pagination[*Device], error) {

	wb := sql.NewLogicBuilder()

	search = strings.TrimSpace(search)
	if search != "" {
		wb.And("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	qb := r.db.WithContext(ctx).
		Model(&Device{})

	where, args := wb.GetLogic()
	if where != "" {
		qb = qb.Where(where, args...)
	}

	var count int64
	total := qb.Count(&count)
	if total.Error != nil {
		return nil, total.Error
	}

	var devices []*Device
	q := qb.Offset(page - 1).Limit(limit).Find(&devices)
	if q.Error != nil {
		return nil, q.Error
	}

	return &Pagination[*Device]{
		Page:  page,
		Limit: limit,
		Total: count,
		Data:  devices,
	}, nil
}

func (r *DeviceRepositoryImpl) CreateDevice(ctx context.Context, device *Device) error {
	if err := device.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *DeviceRepositoryImpl) GetDevice(ctx context.Context, id int) (*Device, error) {
	var result Device
	q := r.db.WithContext(ctx).
		First(&result, id)
	return &result, q.Error
}

func (r *DeviceRepositoryImpl) GetDeviceByPublicID(
	ctx context.Context,
	publicID string,
) (*Device, error) {
	var result Device
	q := r.db.WithContext(ctx).
		Where("public_id = ?", publicID).First(&result)
	return &result, q.Error
}

func (r *DeviceRepositoryImpl) GetDeviceByPublicIDForUser(
	ctx context.Context,
	publicID string,
	userID uint,
) (*Device, error) {
	var result Device
	q := r.db.WithContext(ctx).
		Where("public_id = ? AND user_id = ?", publicID, userID).
		First(&result)
	return &result, q.Error
}

// GetDeletedDeviceByPublicIDForUser is scoped to the owner and is used only to
// make a repeated DELETE a safe no-op. Normal lookups remain default-scoped.
func (r *DeviceRepositoryImpl) GetDeletedDeviceByPublicIDForUser(
	ctx context.Context,
	publicID string,
	userID uint,
) (*Device, error) {
	var result Device
	q := r.db.WithContext(ctx).
		Unscoped().
		Where("public_id = ? AND user_id = ?", publicID, userID).
		First(&result)
	return &result, q.Error
}

// GetDeviceByToken authenticates a token against its device owner relationship.
//
// The current schema has no workspace_id. Until a workspace relation exists, the
// user_id relationship between m_device_tokens and m_devices is the only owner
// scope available. Call GetDeviceByTokenForUser when the caller already carries
// an authenticated user scope.
func (r *DeviceRepositoryImpl) GetDeviceByToken(
	ctx context.Context,
	token string,
) (*Device, error) {
	return r.getDeviceByToken(ctx, token, nil)
}

// GetDeviceByTokenForUser applies the authenticated user's owner scope in
// addition to the token/device relationship predicates. It is intentionally an
// optional method so the existing DeviceRepository interface remains compatible.
func (r *DeviceRepositoryImpl) GetDeviceByTokenForUser(
	ctx context.Context,
	token string,
	userID uint,
) (*Device, error) {
	return r.getDeviceByToken(ctx, token, &userID)
}

func (r *DeviceRepositoryImpl) getDeviceByToken(
	ctx context.Context,
	token string,
	ownerID *uint,
) (*Device, error) {
	var result Device
	now := time.Now()
	qb := r.db.WithContext(ctx).
		Model(&Device{}).
		Joins(`JOIN m_device_tokens AS mdt ON mdt.device_id = m_devices.id
			AND mdt.deleted_at IS NULL
			AND mdt.status = ?
			AND (mdt.expired_at IS NULL OR mdt.expired_at > ?)
			AND mdt.user_id = m_devices.user_id`, devicetoken.DeviceTokenStatusActive.String(), now).
		Where("m_devices.deleted_at IS NULL").
		Where("m_devices.status = ?", DeviceStatusActive.String()).
		Where("mdt.token = ?", token)

	if ownerID != nil {
		qb = qb.Where("m_devices.user_id = ?", *ownerID)
	}

	q := qb.First(&result)
	return &result, q.Error
}

func (r *DeviceRepositoryImpl) UpdateDevice(ctx context.Context, device *Device) error {
	if err := device.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(device).Error
}

func (r *DeviceRepositoryImpl) DeleteDevice(ctx context.Context, device *Device) error {
	if err := device.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(device).Error
}

// listStartupDevices returns only devices that are eligible for provider
// startup. The schema currently has no auto_connect column, so startup must
// not invent or infer one. Until that schema decision is made, active is the
// explicit opt-in state and device.expired_at is the only existing device-level
// credential/session expiry boundary; no separate provider session field exists.
func (r *DeviceRepositoryImpl) listStartupDevices(
	ctx context.Context,
	deviceType DeviceType,
) ([]*Device, error) {
	var devices []*Device
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("m_devices.type = ?", deviceType).
		Where("m_devices.deleted_at IS NULL").
		Where("m_devices.status = ?", DeviceStatusActive.String()).
		Where("m_devices.expired_at IS NULL OR m_devices.expired_at > ?", now).
		Find(&devices).Error
	return devices, err
}

func (r *DeviceRepositoryImpl) GetAllWhatsappDevices(ctx context.Context) ([]*Device, error) {
	return r.listStartupDevices(ctx, DeviceTypeWhatsapp)
}

func (r *DeviceRepositoryImpl) GetAllTelegramDevices(ctx context.Context) ([]*Device, error) {
	return r.listStartupDevices(ctx, DeviceTypeTelegram)
}
