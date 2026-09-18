package telegram

import (
	"context"
	"errors"
	"path/filepath"
	"sapasora/internal/modules/device"
	"sapasora/platform/config"
	"sapasora/platform/logger"
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

type Client struct {
	client     *client.Client
	deviceInfo *TDLibDeviceInfo
	params     *client.SetTdlibParametersRequest

	log *logger.Logger
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
		log: log,
	}
}

func (c *Client) TDLib() (*client.Client, error) {
	if c.client == nil {
		return nil, errors.New("client not initialized")
	}
	return c.client, nil
}

func (c *Client) Params() *client.SetTdlibParametersRequest {
	return c.params
}

func (c *Client) Connect(
	qrCodeHandler func(deviceID string, link string) error,
	timeout time.Duration,
) error {
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

	tdLibClient, err := client.NewClient(authorizer, eventHandler)
	if err != nil {
		return err
	}

	if tdLibClient == nil {
		return errors.New("tdlib client is nil")
	}

	c.client = tdLibClient

	return nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	if c == nil || c.client == nil {
		return nil
	}
	_, err := c.client.Close(ctx)
	return err
}

func (c *Client) EventHandler(result client.Type) {
	// switch result.GetType() {
	// default:
	// 	c.log.Warn("Unknown event", "event", result.GetType())
	// }
}

func (c *Client) IsConnected(ctx context.Context) bool {
	return c != nil && c.client != nil
}

func (c *Client) IsLoggedIn(ctx context.Context) bool {
	if c == nil || c.client == nil {
		return false
	}
	authState, err := c.client.GetAuthorizationState(ctx)
	if err != nil {
		return false
	}

	return authState.AuthorizationStateConstructor() == client.ConstructorAuthorizationStateReady
}

func (c *Client) GetUserByPhoneNumber(ctx context.Context, phone string) (*client.User, error) {
	if c == nil || c.client == nil {
		return nil, ErrNotConnected
	}
	return c.client.SearchUserByPhoneNumber(ctx, &client.SearchUserByPhoneNumberRequest{
		PhoneNumber: phone,
	})

}

func (c *Client) GetUserByUsername(ctx context.Context, username string) (*client.User, error) {
	if c == nil || c.client == nil {
		return nil, ErrNotConnected
	}
	result, err := c.client.SearchPublicChat(ctx, &client.SearchPublicChatRequest{
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	if result.Type.ChatTypeConstructor() != client.ConstructorChatTypePrivate {
		return nil, errors.New("user not found")
	}

	user, err := c.client.GetUser(ctx, &client.GetUserRequest{
		UserId: result.Id,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
