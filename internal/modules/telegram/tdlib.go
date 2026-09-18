package telegram

import (
	"context"
	"sapasora/internal/modules/device"
	"sapasora/platform/config"
	"sapasora/platform/logger"
	"encoding/base64"
	"time"

	"github.com/skip2/go-qrcode"
)

type TDLib struct {
	clientStore *Store[*Client]
	log         *logger.Logger
	config      config.Config

	deviceService device.DeviceService
}

// NewTDLib creates a new TDLib instance
// @wired:provide
func NewTDLib(log *logger.Logger, cfg config.Config, deviceService device.DeviceService) *TDLib {
	return &TDLib{
		clientStore: NewStore[*Client](),
		log:         log,
		config:      cfg,

		deviceService: deviceService,
	}
}

// ConnectDevices to TDLib on server startup if last state was connected
func (t *TDLib) ConnectDevices(ctx context.Context) {
	devices, err := t.deviceService.GetAllTelegramDevices(ctx)
	if err != nil {
		t.log.Error("Failed to get devices", "error", err)
		return
	}

	for _, device := range devices {
		t.Connect(ctx, device)
	}
}

func (t *TDLib) Connect(ctx context.Context, device *device.Device) error {
	t.log.Info("Connect to TDLib", "device", device.Name)

	if client, err := t.GetClient(device.PublicID); err == nil {
		isConnected := client.IsConnected(ctx)
		if isConnected {
			t.log.Info("Already connected to TDLib", "device", device.Name)
			return nil
		}
	}

	t.log.Info("Creating new TDLib client", "device", device.Name)

	client := NewClient(t.config, device, t.log)
	t.clientStore.Set(device.PublicID, client)
	go client.Connect(t.qrHandler, 30*time.Second)

	return nil
}

// GetClient returns the TDLib client for the given device
func (t *TDLib) GetClient(deviceID string) (*Client, error) {
	if client, ok := t.clientStore.Get(deviceID); ok {
		if client == nil {
			return nil, ErrNotConnected
		}

		return client, nil
	}

	return nil, ErrNoSession
}

// qrHandler handles the QR code
func (t *TDLib) qrHandler(deviceID, link string) error {
	if link == "" {
		return t.deviceService.SetDeviceQRCodeByPublicID(context.Background(), deviceID, "")
	}

	image, _ := qrcode.Encode(link, qrcode.Medium, 256)
	base64qrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
	return t.deviceService.SetDeviceQRCodeByPublicID(
		context.Background(),
		deviceID,
		base64qrcode,
	)
}
