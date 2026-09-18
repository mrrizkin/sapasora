package telegram

import (
	"context"
	"sapasora/internal/modules/device"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zelenin/go-tdlib/client"
)

type TelegramServiceImpl struct {
	tele *TDLib
}

// NewTelegramService creates a new Implementation of TelegramService
// @wired:provide
func NewTelegramService(tele *TDLib) TelegramService {
	return &TelegramServiceImpl{
		tele: tele,
	}
}

// CheckUser implements [TelegramService].
func (t *TelegramServiceImpl) CheckUser(
	ctx context.Context,
	device *device.Device,
	payload *CheckUserRequest,
) (*CheckUserResponse, error) {
	c, err := t.tele.GetClient(device.PublicID)
	if err != nil {
		return nil, err
	}

	var users []CheckUser
	for _, username := range payload.Username {
		user, err := c.GetUserByUsername(ctx, username)
		if err != nil {
			users = append(users, CheckUser{
				Query:        username,
				IsInTelegram: false,
				Username:     "",
				VerifiedName: "",
			})
			continue
		}

		users = append(users, CheckUser{
			Query:        username,
			IsInTelegram: true,
			Username:     user.Usernames.ActiveUsernames[0],
			VerifiedName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		})
	}

	for _, phone := range payload.Phone {
		user, err := c.GetUserByPhoneNumber(ctx, phone)
		if err != nil {
			users = append(users, CheckUser{
				Query:        phone,
				IsInTelegram: false,
				Username:     "",
				VerifiedName: "",
			})
			continue
		}

		users = append(users, CheckUser{
			Query:        phone,
			IsInTelegram: true,
			Username:     user.Usernames.ActiveUsernames[0],
			VerifiedName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		})
	}

	return &CheckUserResponse{
		Users: users,
	}, nil
}

// Connect implements [TelegramService].
func (t *TelegramServiceImpl) Connect(
	ctx context.Context,
	device *device.Device,
	payload *ConnectRequest,
) error {
	return t.tele.Connect(ctx, device)
}

// Disconnect implements [TelegramService].
func (t *TelegramServiceImpl) Disconnect(ctx context.Context, device *device.Device) error {
	c, err := t.tele.GetClient(device.PublicID)
	if err != nil {
		return err
	}

	return c.Disconnect(ctx)
}

// GetAvatar implements [TelegramService].
func (t *TelegramServiceImpl) GetAvatar(
	ctx context.Context,
	device *device.Device,
	payload *GetAvatarRequest,
) (*GetAvatarResponse, error) {
	c, err := t.tele.GetClient(device.PublicID)
	if err != nil {
		return nil, err
	}

	user, err := c.GetUserByPhoneNumber(ctx, payload.Phone)
	if err != nil {
		user, err = c.GetUserByUsername(ctx, payload.Username)
		if err != nil {
			return nil, err
		}
	}

	tele, err := c.TDLib()
	if err != nil {
		return nil, err
	}

	photos, err := tele.GetUserProfilePhotos(ctx, &client.GetUserProfilePhotosRequest{
		UserId: user.Id,
		Offset: 0,
		Limit:  1,
	})
	if err != nil {
		return nil, err
	}

	if len(photos.Photos) == 0 {
		return nil, errors.New("no avatar found")
	}

	avatar := photos.Photos[0]
	if avatar.Minithumbnail == nil {
		return nil, errors.New("no avatar found")
	}

	id, err := avatar.Id.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &GetAvatarResponse{
		URL:        byteToImageDataBase64(avatar.Minithumbnail.Data, "image/jpeg"),
		ID:         string(id),
		Type:       "image/jpeg",
		DirectPath: byteToImageDataBase64(avatar.Minithumbnail.Data, "image/jpeg"),
	}, nil
}

func byteToImageDataBase64(data []byte, imageType string) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	switch imageType {
	case "image/jpeg":
		return "data:image/jpeg;base64," + encoded
	case "image/png":
		return "data:image/png;base64," + encoded
	case "image/webp":
		return "data:image/webp;base64," + encoded
	default:
		return "data:image/png;base64," + encoded
	}
}

// GetContacts implements [TelegramService].
func (t *TelegramServiceImpl) GetContacts(
	ctx context.Context,
	device *device.Device,
) (*GetContactsResponse, error) {
	panic("unimplemented")
}

// GetStatus implements [TelegramService].
func (t *TelegramServiceImpl) GetStatus(
	ctx context.Context,
	device *device.Device,
) (*GetStatusResponse, error) {
	c, err := t.tele.GetClient(device.PublicID)
	if err != nil {
		return nil, err
	}

	isConnected := c.IsConnected(ctx)
	isLoggedIn := false
	if isConnected {
		isLoggedIn = c.IsLoggedIn(ctx)
	}

	if isConnected {
		t.tele.deviceService.SetDeviceStatusConnected(ctx, device)
	}

	return &GetStatusResponse{
		Connected: isConnected,
		LoggedIn:  isLoggedIn,
	}, nil
}

// GetUser implements [TelegramService].
func (t *TelegramServiceImpl) GetUser(
	ctx context.Context,
	device *device.Device,
	payload *GetUserRequest,
) (*GetUserResponse, error) {
	panic("unimplemented")
}

// Logout implements [TelegramService].
func (t *TelegramServiceImpl) Logout(ctx context.Context, device *device.Device) error {
	panic("unimplemented")
}

// SendAudio implements [TelegramService].
func (t *TelegramServiceImpl) SendAudio(
	ctx context.Context,
	device *device.Device,
	payload *SendAudioRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendButton implements [TelegramService].
func (t *TelegramServiceImpl) SendButton(
	ctx context.Context,
	device *device.Device,
	payload *SendButtonTextRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendChatPresence implements [TelegramService].
func (t *TelegramServiceImpl) SendChatPresence(
	ctx context.Context,
	device *device.Device,
	payload *ChatPresenceRequest,
) error {
	panic("unimplemented")
}

// SendContact implements [TelegramService].
func (t *TelegramServiceImpl) SendContact(
	ctx context.Context,
	device *device.Device,
	payload *SendContactRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendDocument implements [TelegramService].
func (t *TelegramServiceImpl) SendDocument(
	ctx context.Context,
	device *device.Device,
	payload *SendDocumentRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendImage implements [TelegramService].
func (t *TelegramServiceImpl) SendImage(
	ctx context.Context,
	device *device.Device,
	payload *SendImageRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendList implements [TelegramService].
func (t *TelegramServiceImpl) SendList(
	ctx context.Context,
	device *device.Device,
	payload *SendListRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendLocation implements [TelegramService].
func (t *TelegramServiceImpl) SendLocation(
	ctx context.Context,
	device *device.Device,
	payload *SendLocationRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendSticker implements [TelegramService].
func (t *TelegramServiceImpl) SendSticker(
	ctx context.Context,
	device *device.Device,
	payload *SendStickerRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}

// SendText implements [TelegramService].
func (t *TelegramServiceImpl) SendText(
	ctx context.Context,
	device *device.Device,
	payload *SendTextRequest,
) (*SendResponse, error) {
	c, err := t.tele.GetClient(device.PublicID)
	if err != nil {
		return nil, err
	}

	user, err := c.GetUserByPhoneNumber(ctx, payload.Phone)
	if err != nil {
		user, err = c.GetUserByUsername(ctx, payload.Username)
		if err != nil {
			return nil, err
		}
	}

	tele, err := c.TDLib()
	if err != nil {
		return nil, err
	}

	message, err := tele.SendMessage(ctx, &client.SendMessageRequest{
		ChatId: user.Id,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: payload.Body,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	sender, _ := t.getSender(ctx, tele)
	return &SendResponse{
		Timestamp: time.Unix(int64(message.Date), 0),
		ID:        strconv.Itoa(int(message.Id)),
		ServerID:  0,
		Sender:    sender,
	}, nil
}

func (t *TelegramServiceImpl) getSender(ctx context.Context, tele *client.Client) (string, error) {
	var sender string
	me, err := tele.GetMe(ctx)
	if err != nil {
		return sender, err
	}

	if me != nil {
		if me.Usernames != nil {
			if len(me.Usernames.ActiveUsernames) > 0 {
				sender = me.Usernames.ActiveUsernames[0]
			}
		}
	}

	if sender == "" {
		sender = strings.TrimSpace(me.FirstName + " " + me.LastName)
	}

	// fallback
	return sender, nil
}

// SendVideo implements [TelegramService].
func (t *TelegramServiceImpl) SendVideo(
	ctx context.Context,
	device *device.Device,
	payload *SendVideoRequest,
) (*SendResponse, error) {
	panic("unimplemented")
}
