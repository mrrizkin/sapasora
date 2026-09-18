package telegram

import (
	"context"
	"encoding/base64"
	"sapasora/internal/modules/device"
	"sapasora/internal/modules/providerstartup"
	"sapasora/platform/config"
	"sapasora/platform/logger"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
	"go.uber.org/fx"
)

const (
	startupConcurrency    = 4
	reconnectMaxAttempts  = 3
	reconnectInitialDelay = time.Second
	reconnectMaxDelay     = 5 * time.Second
)

type TDLib struct {
	clientStore *Store[*Client]
	log         *logger.Logger
	config      config.Config

	deviceService device.DeviceService

	startupMu        sync.RWMutex
	lastStartupError error
	startupWG        sync.WaitGroup
	stopStartup      context.CancelFunc
}

// NewTDLib creates a new TDLib instance
// @wired:provide
func NewTDLib(
	log *logger.Logger,
	cfg config.Config,
	deviceService device.DeviceService,
	lc fx.Lifecycle,
) *TDLib {
	runCtx, cancel := context.WithCancel(context.Background())
	t := &TDLib{
		clientStore: NewStore[*Client](),
		log:         log,
		config:      cfg,

		deviceService: deviceService,
		stopStartup:   cancel,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			t.startupWG.Add(1)
			go func() {
				defer t.startupWG.Done()
				t.ConnectDevices(runCtx)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			t.stopStartup()
			t.startupWG.Wait()
			return t.Stop(ctx)
		},
	})

	return t
}

// ConnectDevices to TDLib on server startup if last state was connected
func (t *TDLib) ConnectDevices(ctx context.Context) {
	devices, err := t.deviceService.GetAllTelegramDevices(ctx)
	if err != nil {
		t.recordStartupError(nil, err)
		return
	}

	providerstartup.Run(ctx, devices, startupConcurrency,
		func(ctx context.Context, d *device.Device) error {
			return t.connect(ctx, d, true)
		},
		func(d *device.Device, err error) {
			t.recordStartupError(d, err)
		},
	)
}

func (t *TDLib) Connect(ctx context.Context, device *device.Device) error {
	return t.connect(ctx, device, false)
}

func (t *TDLib) connect(ctx context.Context, device *device.Device, waitForConnection bool) error {
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
	connect := func() error {
		return client.Connect(t.qrHandler, 30*time.Second)
	}
	connectWithBackoff := func() error {
		return providerstartup.Retry(
			ctx,
			reconnectMaxAttempts,
			reconnectInitialDelay,
			reconnectMaxDelay,
			connect,
		)
	}
	if waitForConnection {
		if err := connectWithBackoff(); err != nil {
			t.recordStartupError(device, err)
			return err
		}
		return nil
	}
	go func() {
		if err := connectWithBackoff(); err != nil {
			t.recordStartupError(device, err)
			t.log.Error("Failed to connect to TDLib", "device", device.Name, "error", err)
		}
	}()

	return nil
}

// Stop disconnects every client that was started by this provider. A failed
// client does not prevent the remaining clients from being stopped.
func (t *TDLib) Stop(ctx context.Context) error {
	var stopErr error
	for _, client := range t.clientStore.Values() {
		if err := client.Disconnect(ctx); err != nil {
			if stopErr == nil {
				stopErr = err
			}
			t.log.Error("Failed to stop TDLib client", "error", err)
		}
	}
	return stopErr
}

// LastStartupError returns the most recent per-device startup error, if any.
func (t *TDLib) LastStartupError() error {
	t.startupMu.RLock()
	defer t.startupMu.RUnlock()
	return t.lastStartupError
}

func (t *TDLib) recordStartupError(d *device.Device, err error) {
	if err == nil {
		return
	}
	t.startupMu.Lock()
	t.lastStartupError = err
	t.startupMu.Unlock()
	if d == nil {
		t.log.Error("TDLib startup failed", "error", err)
	}
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
