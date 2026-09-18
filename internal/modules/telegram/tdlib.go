package telegram

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sapasora/internal/modules/device"
	"sapasora/internal/modules/providerstartup"
	"sapasora/platform/config"
	"sapasora/platform/logger"
	"strconv"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
	"go.uber.org/fx"
)

const (
	reconnectMaxAttempts  = 3
	reconnectInitialDelay = time.Second
	reconnectMaxDelay     = 5 * time.Second
)

type TDLib struct {
	clientStore *Store[*Client]
	log         *logger.Logger
	config      config.Config

	deviceService device.DeviceService

	startupMu          sync.RWMutex
	lastStartupError   error
	startupWG          sync.WaitGroup
	startupConcurrency int
	connectWG          sync.WaitGroup
	lifecycleMu        sync.Mutex
	stopping           bool
	stopStartup        context.CancelFunc
	metrics            *providerstartup.LifecycleMetrics
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

		deviceService:      deviceService,
		startupConcurrency: cfg.GetInt("provider.startup_concurrency", providerstartup.DefaultStartupConcurrency),
		stopStartup:        cancel,
		metrics:            providerstartup.NewLifecycleMetrics(),
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

	providerstartup.Run(ctx, devices, t.startupConcurrency,
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

func (t *TDLib) beginConnect() error {
	t.lifecycleMu.Lock()
	defer t.lifecycleMu.Unlock()
	if t.stopping {
		return ErrProviderStopping
	}
	t.connectWG.Add(1)
	return nil
}

func (t *TDLib) endConnect() {
	t.connectWG.Done()
}

func (t *TDLib) connect(ctx context.Context, device *device.Device, waitForConnection bool) error {
	if err := t.beginConnect(); err != nil {
		return err
	}
	release := true
	defer func() {
		if release {
			t.endConnect()
		}
	}()

	t.log.Info("Connect to TDLib", "device", device.Name, "wait_for_connection", waitForConnection)

	client, err := t.GetClient(device.PublicID)
	if err != nil {
		t.log.Info("Creating new TDLib client", "device", device.Name)
		client = NewClient(t.config, device, t.log)
		t.clientStore.Set(device.PublicID, client)
	} else if state := client.State(); state == ClientStateConnected || state == ClientStateConnecting {
		t.log.Info("Already connected to TDLib", "device", device.Name, "state", state)
		if state == ClientStateConnected {
			return t.syncConnectedDevice(ctx, device, client)
		}
		return nil
	}
	client.SetClosedHandler(func() {
		if err := t.deviceService.SetDeviceStatusDisconnectedByPublicID(context.Background(), device.PublicID); err != nil {
			t.log.Error("Failed to persist TDLib disconnected status", "device", device.Name, "error", err)
		}
	})
	connect := func() error {
		t.metrics.RecordStartupAttempt("telegram")
		if err := client.Connect(t.qrHandler, 30*time.Second); err != nil {
			t.metrics.RecordStartupFailure("telegram")
			return err
		}
		if err := t.syncConnectedDevice(ctx, device, client); err != nil {
			// Do not leave a usable provider session behind when its durable
			// identity/status could not be persisted.
			if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
				t.log.Error("Failed to rollback TDLib connection", "device", device.Name, "error", disconnectErr)
			}
			t.metrics.RecordStartupFailure("telegram")
			return err
		}
		t.metrics.RecordStartupSuccess("telegram")
		return nil
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
	go t.superviseAsyncConnect(device, connectWithBackoff)
	release = false

	return nil
}

// superviseAsyncConnect is the single error boundary for manual connections
// started in the background. Errors must be recorded before the goroutine
// exits; otherwise Connect would report success while its provider failure
// became invisible to lifecycle supervision. A panic from provider code is
// converted to a stable error so it cannot terminate the process.
func (t *TDLib) superviseAsyncConnect(device *device.Device, connect func() error) {
	defer t.endConnect()
	defer func() {
		if recover() != nil {
			t.recordStartupError(device, ErrAsyncConnectionPanic)
			if t.log != nil {
				t.log.Error("Failed to connect to TDLib", "device", deviceName(device), "error", ErrAsyncConnectionPanic)
			}
		}
	}()

	if err := connect(); err != nil {
		t.recordStartupError(device, err)
		if t.log != nil {
			t.log.Error("Failed to connect to TDLib", "device", deviceName(device), "error", err)
		}
	}
}

func deviceName(d *device.Device) string {
	if d == nil {
		return ""
	}
	return d.Name
}

// Stop disconnects every client that was started by this provider. A failed
// client does not prevent the remaining clients from being stopped.
func (t *TDLib) Stop(ctx context.Context) error {
	t.lifecycleMu.Lock()
	t.stopping = true
	t.lifecycleMu.Unlock()

	// Wait for manual Connect calls as well as startup workers. Otherwise a
	// late async connection can install a live TDLib client after shutdown has
	// already started closing the current store.
	t.connectWG.Wait()

	if t.clientStore == nil {
		return nil
	}

	var stopErr error
	for _, client := range t.clientStore.Values() {
		if client == nil {
			continue
		}
		wasConnected := client.IsConnected(ctx)
		if err := client.Disconnect(ctx); err != nil {
			if stopErr == nil {
				stopErr = err
			}
			t.log.Error("Failed to stop TDLib client", "error", err)
			continue
		}
		if wasConnected {
			t.metrics.RecordDisconnect("telegram")
		}
		if err := t.markDisconnected(ctx, client); err != nil {
			if stopErr == nil {
				stopErr = err
			}
			t.log.Error("Failed to persist TDLib disconnected status", "error", err)
		}
	}
	return stopErr
}

// Disconnect stops one Telegram client and records the lifecycle transition.
func (t *TDLib) Disconnect(ctx context.Context, deviceID string) error {
	client, err := t.GetClient(deviceID)
	if err != nil {
		return err
	}
	wasConnected := client.IsConnected(ctx)
	if err := client.Disconnect(ctx); err != nil {
		return err
	}
	if wasConnected {
		t.metrics.RecordDisconnect("telegram")
	}
	if err := t.markDisconnected(ctx, client); err != nil {
		return err
	}
	return nil
}

type telegramLifecycleDeviceService interface {
	SetDeviceJID(context.Context, *device.Device, string) error
	SetDeviceStatusConnected(context.Context, *device.Device) error
}

// syncConnectedDevice persists the provider-owned Telegram identity and
// connected state together with the lifecycle transition. Telegram uses its
// numeric account ID as the stable JID value.
func (t *TDLib) syncConnectedDevice(ctx context.Context, d *device.Device, client *Client) error {
	if t.deviceService == nil {
		return errors.New("telegram device service is not configured")
	}
	return syncConnectedDevice(ctx, t.deviceService, d, client)
}

func syncConnectedDevice(
	ctx context.Context,
	service telegramLifecycleDeviceService,
	d *device.Device,
	client *Client,
) error {
	if service == nil {
		return errors.New("telegram device service is not configured")
	}
	if d == nil || client == nil {
		return errors.New("telegram device identity is not initialized")
	}
	me, err := client.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("resolve Telegram identity: %w", err)
	}
	if me == nil || me.Id <= 0 {
		return errors.New("resolve Telegram identity: empty account ID")
	}
	if err := service.SetDeviceJID(ctx, d, strconv.FormatInt(me.Id, 10)); err != nil {
		return fmt.Errorf("persist Telegram identity: %w", err)
	}
	if err := service.SetDeviceStatusConnected(ctx, d); err != nil {
		return fmt.Errorf("persist Telegram connected status: %w", err)
	}
	return nil
}

func (t *TDLib) markDisconnected(ctx context.Context, client *Client) error {
	if t.deviceService == nil || client == nil || client.deviceInfo == nil {
		return nil
	}
	if err := t.deviceService.SetDeviceStatusDisconnectedByPublicID(ctx, client.deviceInfo.ID); err != nil {
		return fmt.Errorf("persist Telegram disconnected status: %w", err)
	}
	return nil
}

// LifecycleMetrics returns the provider lifecycle counters for operational
// inspection and metrics export.
func (t *TDLib) LifecycleMetrics() *providerstartup.LifecycleMetrics {
	return t.metrics
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
