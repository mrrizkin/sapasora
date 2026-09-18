package telegram

import (
	"context"
	"errors"
	"path/filepath"
	"sapasora/internal/modules/device"
	"sapasora/platform/config"
	"sapasora/platform/logger"
	"sync"
	"time"

	"codeberg.org/mrrizkin/nihil"
	"github.com/zelenin/go-tdlib/client"
)

type TDLibDeviceInfo struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Events  nihil.NilString `json:"events"`
	Webhook nihil.NilString `json:"webhook"`
}

type ClientState string

const (
	ClientStateConnecting    ClientState = "connecting"
	ClientStateConnected     ClientState = "connected"
	ClientStateDisconnecting ClientState = "disconnecting"
	ClientStateDisconnected  ClientState = "disconnected"
	ClientStateError         ClientState = "error"
)

type Client struct {
	mu         sync.RWMutex
	client     *client.Client
	deviceInfo *TDLibDeviceInfo
	params     *client.SetTdlibParametersRequest
	state      ClientState

	log *logger.Logger

	// newClient is injectable so lifecycle behavior can be tested without a
	// running TDLib process. Production clients use client.NewClient.
	newClient func(client.AuthorizationStateHandler, client.Option) (*client.Client, error)
	getMe     func(context.Context) (*client.User, error)
	onClosed  func()
}

func NewClient(config config.Config, device *device.Device, log *logger.Logger) *Client {
	return &Client{
		params: &client.SetTdlibParametersRequest{
			UseTestDc:           false,
			DatabaseDirectory:   filepath.Join(".tdlib", device.PublicID, "database"),
			FilesDirectory:      filepath.Join(".tdlib", device.PublicID, "files"),
			UseFileDatabase:     true,
			UseChatInfoDatabase: true,
			UseMessageDatabase:  true,
			UseSecretChats:      false,
			ApiId:               int32(config.GetInt("telegram.api_id")),
			ApiHash:             config.GetString("telegram.api_hash"),
			SystemLanguageCode:  "en",
			DeviceModel:         "Server",
			SystemVersion:       "1.0.0",
			ApplicationVersion:  "1.0.0",
		},
		deviceInfo: &TDLibDeviceInfo{
			ID:      device.PublicID,
			Name:    device.Name,
			Events:  device.Events,
			Webhook: device.Webhook,
		},
		state: ClientStateDisconnected,
		log:   log,
		newClient: func(
			authorizer client.AuthorizationStateHandler,
			option client.Option,
		) (*client.Client, error) {
			return client.NewClient(authorizer, option)
		},
	}
}

func (c *Client) State() ClientState {
	if c == nil {
		return ClientStateDisconnected
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.state == "" {
		return ClientStateDisconnected
	}
	return c.state
}

func (c *Client) markConnectError() {
	c.mu.Lock()
	if c.state == ClientStateConnecting {
		c.state = ClientStateError
	}
	c.mu.Unlock()
}

func (c *Client) TDLib() (*client.Client, error) {
	if c == nil {
		return nil, errors.New("client not initialized")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.client == nil || c.state != ClientStateConnected {
		return nil, errors.New("client not initialized")
	}
	return c.client, nil
}

func (c *Client) Params() *client.SetTdlibParametersRequest {
	return c.params
}

// GetMe returns the authenticated Telegram identity used as the device JID.
// The hook keeps lifecycle tests independent from a native TDLib process.
func (c *Client) GetMe(ctx context.Context) (*client.User, error) {
	if c == nil {
		return nil, ErrNotConnected
	}
	if c.getMe != nil {
		return c.getMe(ctx)
	}
	tele, err := c.TDLib()
	if err != nil {
		return nil, err
	}
	return tele.GetMe(ctx)
}

func (c *Client) SetClosedHandler(handler func()) {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.onClosed = handler
	c.mu.Unlock()
}

func (c *Client) Connect(
	qrCodeHandler func(deviceID string, link string) error,
	timeout time.Duration,
) error {
	if c == nil {
		return ErrNotConnected
	}

	c.mu.Lock()
	if c.state == ClientStateConnected || c.state == ClientStateConnecting || c.state == ClientStateDisconnecting {
		c.mu.Unlock()
		return ErrClientAlreadyConnected
	}
	c.state = ClientStateConnecting
	factory := c.newClient
	if factory == nil {
		factory = func(
			authorizer client.AuthorizationStateHandler,
			option client.Option,
		) (*client.Client, error) {
			return client.NewClient(authorizer, option)
		}
	}
	c.mu.Unlock()

	start := time.Now()
	authorizer := client.QrAuthorizer(c.params, func(link string) error {
		if time.Since(start) > timeout {
			if err := qrCodeHandler(c.deviceInfo.ID, ""); err != nil {
				return err
			}
			return errors.New("timeout")
		}

		return qrCodeHandler(c.deviceInfo.ID, link)
	})

	eventHandler := client.WithResultHandler(client.NewCallbackResultHandler(c.EventHandler))

	tdLibClient, err := factory(authorizer, eventHandler)
	if err != nil {
		c.markConnectError()
		return err
	}

	if tdLibClient == nil {
		c.markConnectError()
		return errors.New("tdlib client is nil")
	}

	c.mu.Lock()
	if c.state != ClientStateConnecting {
		c.mu.Unlock()
		return ErrNotConnected
	}
	c.client = tdLibClient
	c.state = ClientStateConnected
	c.mu.Unlock()

	return nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	if c == nil {
		return nil
	}

	c.mu.Lock()
	if c.client == nil {
		c.state = ClientStateDisconnected
		c.mu.Unlock()
		return nil
	}
	tele := c.client
	c.state = ClientStateDisconnecting
	c.mu.Unlock()

	// Do not hold c.mu while Close waits for TDLib's closed authorization
	// update; EventHandler needs the same lock to consume that update.
	_, err := tele.Close(ctx)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		if c.client == tele {
			c.state = ClientStateError
		}
		return err
	}
	if c.client == tele {
		c.client = nil
		c.state = ClientStateDisconnected
	}
	return nil
}

func (c *Client) EventHandler(result client.Type) {
	update, ok := result.(*client.UpdateAuthorizationState)
	if !ok || update.AuthorizationState == nil {
		return
	}
	if update.AuthorizationState.AuthorizationStateConstructor() == client.ConstructorAuthorizationStateClosed {
		c.mu.Lock()
		c.client = nil
		c.state = ClientStateDisconnected
		onClosed := c.onClosed
		c.mu.Unlock()
		if onClosed != nil {
			onClosed()
		}
	}
}

func (c *Client) IsConnected(ctx context.Context) bool {
	return c != nil && c.State() == ClientStateConnected
}

func (c *Client) IsLoggedIn(ctx context.Context) bool {
	tele, err := c.TDLib()
	if err != nil {
		return false
	}
	authState, err := tele.GetAuthorizationState(ctx)
	if err != nil {
		return false
	}

	return authState.AuthorizationStateConstructor() == client.ConstructorAuthorizationStateReady
}

func (c *Client) GetUserByPhoneNumber(ctx context.Context, phone string) (*client.User, error) {
	tele, err := c.TDLib()
	if err != nil {
		return nil, ErrNotConnected
	}
	return tele.SearchUserByPhoneNumber(ctx, &client.SearchUserByPhoneNumberRequest{
		PhoneNumber: phone,
	})

}

func (c *Client) GetUserByUsername(ctx context.Context, username string) (*client.User, error) {
	tele, err := c.TDLib()
	if err != nil {
		return nil, ErrNotConnected
	}
	result, err := tele.SearchPublicChat(ctx, &client.SearchPublicChatRequest{
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	if result.Type.ChatTypeConstructor() != client.ConstructorChatTypePrivate {
		return nil, errors.New("user not found")
	}

	user, err := tele.GetUser(ctx, &client.GetUserRequest{
		UserId: result.Id,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
