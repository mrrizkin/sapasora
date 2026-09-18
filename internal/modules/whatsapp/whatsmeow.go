package whatsapp

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"sapasora/internal/modules/device"
	"sapasora/platform/config"
	"sapasora/platform/database"
	"sapasora/platform/logger"
	"slices"
	"strings"
	"sync"
	"time"

	"codeberg.org/mrrizkin/nihil"
	"github.com/go-resty/resty/v2"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

type WhatsmeowDeviceInfo struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Events  nihil.NilString `json:"events"`
	Jid     nihil.NilString `json:"jid"`
	Webhook nihil.NilString `json:"webhook"`
}

func NewDeviceInfoFromDevice(device *device.Device) *WhatsmeowDeviceInfo {
	return &WhatsmeowDeviceInfo{
		ID:      device.PublicID,
		Name:    device.Name,
		Events:  device.Events,
		Jid:     device.Jid,
		Webhook: device.Webhook,
	}
}

type Whatsmeow struct {
	clientStore     *Store[*whatsmeow.Client]
	clientHTTP      *Store[*resty.Client]
	deviceInfoStore *Store[*WhatsmeowDeviceInfo]
	killchannel     *Store[(chan bool)]

	lifecycleMu sync.Mutex

	container *sqlstore.Container
	log       *logger.Logger

	deviceService device.DeviceService
}

var MessageTypes = []string{
	"Message",
	"ReadReceipt",
	"Presence",
	"HistorySync",
	"ChatPresence",
	"All",
}

// NewWhatsmeow creates a new Whatsmeow instance.
// @wired:provide
func NewWhatsmeow(
	logger *logger.Logger,
	config config.Config,

	deviceService device.DeviceService,
) (*Whatsmeow, error) {
	dbConfig := database.NewConfig(config)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	container, err := sqlstore.New(
		ctx,
		dbConfig.Protocol(),
		dbConfig.DSN(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Whatsmeow{
		clientStore:     NewStore[*whatsmeow.Client](),
		clientHTTP:      NewStore[*resty.Client](),
		deviceInfoStore: NewStore[*WhatsmeowDeviceInfo](),
		killchannel:     NewStore[(chan bool)](),

		container: container,
		log:       logger.Scope("whatsmeow"),

		deviceService: deviceService,
	}, nil
}

// ConnectDevices to Whatsmeow Websocket on server startup if last state was connected
func (w *Whatsmeow) ConnectDevices(ctx context.Context) {
	devices, err := w.deviceService.GetAllWhatsappDevices(ctx)
	if err != nil {
		w.log.Error("Failed to get devices", "error", err)
		return
	}

	for _, device := range devices {
		w.log.Info("Connect to Whatsmeow on startup", "device", device.Name)

		deviceInfo := NewDeviceInfoFromDevice(device)
		w.deviceInfoStore.Set(device.PublicID, deviceInfo)

		// Gets and set subscription to webhook events
		var events []string
		if deviceInfo.Events.Valid {
			events = strings.Split(deviceInfo.Events.String, ",")
		}

		var subscribedEvents []string
		if len(events) > 0 {
			for _, event := range events {
				if !slices.Contains(MessageTypes, event) {
					w.log.Warn("Message type discarded", "device", deviceInfo.Name, "type", event)
					continue
				}
				if !slices.Contains(subscribedEvents, event) {
					subscribedEvents = append(subscribedEvents, event)
				}
			}
		} else {
			subscribedEvents = append(subscribedEvents, "All")
		}

		eventstring := strings.Join(subscribedEvents, ",")
		w.log.Info("Attempt to connect", "device", deviceInfo.Name, "events", eventstring)
		w.NewKillChannel(deviceInfo.ID)
		go w.StartClient(ctx, deviceInfo, subscribedEvents)
	}
}

func (w *Whatsmeow) StartClient(
	ctx context.Context,
	deviceInfo *WhatsmeowDeviceInfo,
	subscriptions []string,
) {
	w.log.Info("Starting websocket connection to Whatsmeow", "device", deviceInfo.Name)

	var deviceStore *store.Device
	var err error

	if client, ok := w.clientStore.Get(deviceInfo.ID); ok && client != nil {
		isConnected := client.IsConnected()
		if isConnected {
			w.log.Info("Already connected to Whatsmeow", "device", deviceInfo.Name)
			return
		}
	}

	// Capture the channel for this connection attempt before provider setup. A reconnect
	// replaces the channel, so an older attempt must not start listening on the new one.
	killchannel, ok := w.killchannel.Get(deviceInfo.ID)
	if !ok || killchannel == nil {
		w.log.Error("Kill channel not found", "device", deviceInfo.Name)
		return
	}

	if deviceInfo.Jid.Valid {
		jid, valid := w.ParseJID(deviceInfo.Jid.String)
		if !valid {
			w.log.Warn("Invalid stored jid. Creating new device", "device", deviceInfo.Name)
			deviceStore = w.container.NewDevice()
		} else {
			deviceCtx, cancel := providerContext(ctx)
			deviceStore, err = w.container.GetDevice(deviceCtx, jid)
			cancel()
			if err != nil {
				w.log.Error("Failed to load WhatsApp device store", "device", deviceInfo.Name, "error", err)
				return
			}
		}
	} else {
		w.log.Warn("No jid found. Creating new device", "device", deviceInfo.Name)
		deviceStore = w.container.NewDevice()
	}

	if deviceStore == nil {
		w.log.Warn("No store found. Creating new one", "device", deviceInfo.Name)
		deviceStore = w.container.NewDevice()
	}

	// osName := "Sapasora"
	// store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_CHROME.Enum()
	// store.DeviceProps.Os = &osName

	whatsmeowClient := whatsmeow.NewClient(deviceStore, nil)

	w.lifecycleMu.Lock()
	currentKillChannel, current := w.killchannel.Get(deviceInfo.ID)
	if !current || currentKillChannel != killchannel {
		w.lifecycleMu.Unlock()
		return
	}
	w.clientStore.Set(deviceInfo.ID, whatsmeowClient)
	w.lifecycleMu.Unlock()

	client := NewClient(
		whatsmeowClient,
		deviceInfo,
		subscriptions,
		w,
	)
	client.SetEventHandlerID(client.WAClient.AddEventHandler(client.EventHandler))

	httpClient := resty.New()
	httpClient.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	httpClient.SetTimeout(5 * time.Second)
	httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	w.lifecycleMu.Lock()
	currentKillChannel, current = w.killchannel.Get(deviceInfo.ID)
	if !current || currentKillChannel != killchannel {
		w.lifecycleMu.Unlock()
		return
	}
	w.clientHTTP.Set(deviceInfo.ID, httpClient)
	w.lifecycleMu.Unlock()

	if whatsmeowClient.Store.ID == nil {
		// No ID stored, new login

		qrChan, err := whatsmeowClient.GetQRChannel(context.WithoutCancel(ctx))
		if err != nil {
			// This error means that we're already logged in, so ignore it.
			if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
				w.log.Error("Failed to get QR channel", "device", deviceInfo.Name, "error", err)
			}
		} else {
			go func() {
				connectErr := connectWithTimeout(ctx, whatsmeowClient)
				if connectErr != nil {
					w.log.Error("Failed to connect to WhatsMeow", "device", deviceInfo.Name, "error", connectErr)
					return
				}
				for evt := range qrChan {
					switch evt.Event {
					case "code":
						// Store encoded/embeded base64 QR on database for retrieval
						image, _ := qrcode.Encode(evt.Code, qrcode.Medium, 256)
						base64qrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
						err := w.deviceService.SetDeviceQRCodeByPublicID(ctx, deviceInfo.ID, base64qrcode)
						if err != nil {
							w.log.Error("Failed to set QR code", "device", deviceInfo.Name, "error", err)
						}
					case "timeout":
						// Clear QR code from DB on timeout
						err := w.deviceService.SetDeviceQRCodeByPublicID(ctx, deviceInfo.ID, "")
						if err != nil {
							w.log.Error("Failed to set QR code", "device", deviceInfo.Name, "error", err)
						}
						w.log.Warn("QR timeout killing channel", "device", deviceInfo.Name)
						sendKillSignal(killchannel)
					case "success":
						w.log.Info("QR pairing ok!", "device", deviceInfo.Name)

						// Clear QR code after pairing
						err := w.deviceService.SetDeviceQRCodeByPublicID(ctx, deviceInfo.ID, "")
						if err != nil {
							w.log.Error("Failed to set QR code", "device", deviceInfo.Name, "error", err)
						}
					default:
						w.log.Info("Login event", "device", deviceInfo.Name)
					}
				}
			}()
		}
	} else {
		// Already logged in, just connect
		w.log.Info("Already logged in, just connect", "device", deviceInfo.Name)
		err = connectWithTimeout(ctx, whatsmeowClient)
		if err != nil {
			w.log.Error("Failed to connect to WhatsMeow", "device", deviceInfo.Name, "error", err)
			return
		}
	}

	// Keep connected client live until disconnected/killed
	for {
		select {
		case <-killchannel:
			w.log.Info("Received kill signal", "device", deviceInfo.Name)
			whatsmeowClient.Disconnect()
			w.cleanupSession(deviceInfo.ID, whatsmeowClient, killchannel)
			err := w.deviceService.SetDeviceStatusDisconnectedByPublicID(ctx, deviceInfo.ID)
			if err != nil {
				w.log.Error(
					"Failed to set user disconnected",
					"device",
					deviceInfo.Name,
					"error",
					err,
				)
			}
			return
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func connectWithTimeout(ctx context.Context, client *whatsmeow.Client) error {
	connectCtx, cancel := providerContext(ctx)
	defer cancel()

	result := make(chan error, 1)
	go func() {
		result <- client.Connect()
	}()

	select {
	case err := <-result:
		return err
	case <-connectCtx.Done():
		client.Disconnect()
		return connectCtx.Err()
	}
}

func (w *Whatsmeow) ParseJID(jid string) (types.JID, bool) {
	recipient, ok := parseJID(jid)
	if !ok && w != nil && w.log != nil {
		w.log.Warn("Bad jid format, return empty", "jid", jid)
	}
	return recipient, ok
}

func parseJID(jid string) (types.JID, bool) {
	if jid == "" {
		return types.NewJID("", types.DefaultUserServer), false
	}
	if jid[0] == '+' {
		jid = jid[1:]
	}

	// Basic only digit check for recipient phone number, we want to remove @server and .session.
	phonenumber := strings.Split(jid, "@")[0]
	phonenumber = strings.Split(phonenumber, ".")[0]
	phonenumber = strings.Split(phonenumber, ":")[0]
	if phonenumber == "" {
		return types.NewJID("", types.DefaultUserServer), false
	}
	for _, c := range phonenumber {
		if c < '0' || c > '9' {
			return types.NewJID("", types.DefaultUserServer), false
		}
	}

	if !strings.ContainsRune(jid, '@') {
		return types.NewJID(jid, types.DefaultUserServer), true
	}

	recipient, err := types.ParseJID(jid)
	if err != nil || recipient.User == "" {
		return recipient, false
	}
	return recipient, true
}

func (w *Whatsmeow) CallHook(myurl string, payload map[string]string, deviceID string) {
	w.log.Debug("Sending POST", "url", myurl)
	if client, ok := w.clientHTTP.Get(deviceID); ok {
		_, err := client.R().SetFormData(payload).Post(myurl)

		if err != nil {
			w.log.Error("Failed to send webhook", "url", myurl, "error", err)
		}
	}
}

func (w *Whatsmeow) CallHookFile(
	myurl string,
	payload map[string]string,
	deviceID string,
	path string,
) {
	w.log.Debug("Sending POST", "url", myurl)
	if client, ok := w.clientHTTP.Get(deviceID); ok {
		_, err := client.R().SetFiles(map[string]string{
			"file": path,
		}).SetFormData(payload).Post(myurl)

		if err != nil {
			w.log.Error("Failed to send webhook", "url", myurl, "error", err)
		}
	}
}

func (w *Whatsmeow) GetClient(deviceID string) (*whatsmeow.Client, error) {
	if client, ok := w.clientStore.Get(deviceID); ok {
		if client == nil {
			return nil, ErrNotConnected
		}

		return client, nil
	}

	return nil, ErrNoSession
}

func (w *Whatsmeow) NewKillChannel(deviceID string) {
	w.lifecycleMu.Lock()
	defer w.lifecycleMu.Unlock()

	if w.killchannel == nil {
		w.killchannel = NewStore[chan bool]()
	}
	if killchannel, ok := w.killchannel.Get(deviceID); ok {
		if w.clientStore != nil {
			if client, found := w.clientStore.Get(deviceID); found && client != nil && client.IsConnected() {
				return
			}
		}
		sendKillSignal(killchannel)
	}
	w.killchannel.Set(deviceID, make(chan bool, 1))
}

func (w *Whatsmeow) SendKillChannel(deviceID string) {
	if w.killchannel == nil {
		return
	}
	if killchannel, ok := w.killchannel.Get(deviceID); ok {
		sendKillSignal(killchannel)
	}
}

func sendKillSignal(killchannel chan bool) {
	if killchannel == nil {
		return
	}
	select {
	case killchannel <- true:
	default:
	}
}

func (w *Whatsmeow) cleanupSession(
	deviceID string,
	client *whatsmeow.Client,
	killchannel chan bool,
) {
	w.lifecycleMu.Lock()
	defer w.lifecycleMu.Unlock()

	if current, ok := w.clientStore.Get(deviceID); ok && current == client {
		w.clientStore.Delete(deviceID)
		w.clientHTTP.Delete(deviceID)
	}
	if current, ok := w.killchannel.Get(deviceID); ok && current == killchannel {
		w.killchannel.Delete(deviceID)
	}
}

func validateMessageFields(
	phone string,
	stanzaid *string,
	participant *string,
	wa *Whatsmeow,
) (types.JID, error) {
	recipient, ok := wa.ParseJID(phone)
	if !ok {
		return types.NewJID("", types.DefaultUserServer), ErrInvalidPhoneNumber
	}

	if stanzaid != nil {
		if participant == nil {
			return types.NewJID(
				"",
				types.DefaultUserServer,
			), ErrMissingParticipant
		}
	}

	if participant != nil {
		if stanzaid == nil {
			return types.NewJID(
				"",
				types.DefaultUserServer,
			), ErrMissingStanzaID
		}
	}

	return recipient, nil
}
