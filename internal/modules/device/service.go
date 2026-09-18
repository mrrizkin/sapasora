package device

import (
	"context"
	"errors"

	"codeberg.org/mrrizkin/nihil"
)

type DeviceServiceImpl struct {
	repo DeviceRepository
}

// NewDeviceService creates a new Implementation of DeviceService
// @wired:provide
func NewDeviceService(repo DeviceRepository) DeviceService {
	return &DeviceServiceImpl{
		repo: repo,
	}
}

func (s *DeviceServiceImpl) ListDevice(
	ctx context.Context,
	search string,
	page, limit int,
) (*Pagination[*Device], error) {
	return s.repo.ListDevice(ctx, search, page, limit)
}

func (s *DeviceServiceImpl) CreateDevice(ctx context.Context, device *Device) error {
	return s.repo.CreateDevice(ctx, device)
}

func (s *DeviceServiceImpl) GetDevice(ctx context.Context, id int) (*Device, error) {
	return s.repo.GetDevice(ctx, id)
}

func (s *DeviceServiceImpl) GetDeviceByPublicID(
	ctx context.Context,
	publicID string,
) (*Device, error) {
	return s.repo.GetDeviceByPublicID(ctx, publicID)
}

func (s *DeviceServiceImpl) GetDeviceByToken(
	ctx context.Context,
	token string,
) (*Device, error) {
	return s.repo.GetDeviceByToken(ctx, token)
}

// GetDeviceByTokenForUser delegates to the optional owner-scoped repository
// extension and fails closed when an implementation does not provide it.
func (s *DeviceServiceImpl) GetDeviceByTokenForUser(
	ctx context.Context,
	token string,
	userID uint,
) (*Device, error) {
	repo, ok := s.repo.(DeviceTokenOwnerScopedRepository)
	if !ok {
		return nil, errors.New("device token owner-scoped lookup is unavailable")
	}
	return repo.GetDeviceByTokenForUser(ctx, token, userID)
}

func (s *DeviceServiceImpl) UpdateDevice(ctx context.Context, device *Device) error {
	return s.repo.UpdateDevice(ctx, device)
}

func (s *DeviceServiceImpl) DeleteDevice(ctx context.Context, device *Device) error {
	return s.repo.DeleteDevice(ctx, device)
}

func (s *DeviceServiceImpl) GetAllWhatsappDevices(ctx context.Context) ([]*Device, error) {
	return s.repo.GetAllWhatsappDevices(ctx)
}

func (s *DeviceServiceImpl) GetAllTelegramDevices(ctx context.Context) ([]*Device, error) {
	return s.repo.GetAllTelegramDevices(ctx)
}

func (s *DeviceServiceImpl) SetDeviceQRCode(
	ctx context.Context,
	device *Device,
	qrCode string,
) error {
	device.QRCode = nihil.String(qrCode)
	return s.repo.UpdateDevice(ctx, device)
}

func (s *DeviceServiceImpl) SetDeviceQRCodeByPublicID(
	ctx context.Context,
	publicID string,
	qrCode string,
) error {
	device, err := s.repo.GetDeviceByPublicID(ctx, publicID)
	if err != nil {
		return err
	}

	return s.SetDeviceQRCode(ctx, device, qrCode)
}

func (s *DeviceServiceImpl) SetDeviceJID(ctx context.Context, device *Device, jid string) error {
	device.Jid = nihil.String(jid)
	return s.repo.UpdateDevice(ctx, device)
}

func (s *DeviceServiceImpl) SetDeviceJIDByPublicID(
	ctx context.Context,
	publicID string,
	jid string,
) error {
	device, err := s.repo.GetDeviceByPublicID(ctx, publicID)
	if err != nil {
		return err
	}

	return s.SetDeviceJID(ctx, device, jid)
}

func (s *DeviceServiceImpl) SetDeviceStatusConnected(ctx context.Context, device *Device) error {
	device.Status = DeviceStatusConnected
	return s.repo.UpdateDevice(ctx, device)
}

func (s *DeviceServiceImpl) SetDeviceStatusConnectedByPublicID(
	ctx context.Context,
	publicID string,
) error {
	device, err := s.repo.GetDeviceByPublicID(ctx, publicID)
	if err != nil {
		return err
	}

	return s.SetDeviceStatusConnected(ctx, device)
}

func (s *DeviceServiceImpl) SetDeviceStatusDisconnected(ctx context.Context, device *Device) error {
	device.Status = DeviceStatusDisconnected
	return s.repo.UpdateDevice(ctx, device)
}

func (s *DeviceServiceImpl) SetDeviceStatusDisconnectedByPublicID(
	ctx context.Context,
	publicID string,
) error {
	device, err := s.repo.GetDeviceByPublicID(ctx, publicID)
	if err != nil {
		return err
	}

	return s.SetDeviceStatusDisconnected(ctx, device)
}
