package device

import (
	"context"
	"strings"

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

func (r *DeviceRepositoryImpl) GetDeviceByToken(
	ctx context.Context,
	token string,
) (*Device, error) {
	var result Device
	q := r.db.WithContext(ctx).
		Joins("JOIN m_device_tokens AS mdt ON mdt.device_id = m_devices.id").
		Where("mdt.token = ?", token).First(&result)
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

func (r *DeviceRepositoryImpl) GetAllWhatsappDevices(ctx context.Context) ([]*Device, error) {
	var devices []*Device
	err := r.db.WithContext(ctx).Where("type = ?", DeviceTypeWhatsapp).Find(&devices).Error
	return devices, err
}

func (r *DeviceRepositoryImpl) GetAllTelegramDevices(ctx context.Context) ([]*Device, error) {
	var devices []*Device
	err := r.db.WithContext(ctx).Where("type = ?", DeviceTypeTelegram).Find(&devices).Error
	return devices, err
}
