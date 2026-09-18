package telegram

import (
	"context"
	"sapasora/internal/modules/device"
)

type TelegramService interface {
	CheckUser(
		ctx context.Context,
		device *device.Device,
		payload *CheckUserRequest,
	) (*CheckUserResponse, error)
	Connect(ctx context.Context, device *device.Device, payload *ConnectRequest) error
	Disconnect(ctx context.Context, device *device.Device) error
	GetAvatar(
		ctx context.Context,
		device *device.Device,
		payload *GetAvatarRequest,
	) (*GetAvatarResponse, error)
	GetContacts(
		ctx context.Context,
		device *device.Device,
	) (*GetContactsResponse, error)
	GetStatus(
		ctx context.Context,
		device *device.Device,
	) (*GetStatusResponse, error)
	GetUser(
		ctx context.Context,
		device *device.Device,
		payload *GetUserRequest,
	) (*GetUserResponse, error)
	Logout(ctx context.Context, device *device.Device) error
	SendAudio(
		ctx context.Context,
		device *device.Device,
		payload *SendAudioRequest,
	) (*SendResponse, error)
	SendButton(
		ctx context.Context,
		device *device.Device,
		payload *SendButtonTextRequest,
	) (*SendResponse, error)
	SendChatPresence(ctx context.Context, device *device.Device, payload *ChatPresenceRequest) error
	SendContact(
		ctx context.Context,
		device *device.Device,
		payload *SendContactRequest,
	) (*SendResponse, error)
	SendDocument(
		ctx context.Context,
		device *device.Device,
		payload *SendDocumentRequest,
	) (*SendResponse, error)
	SendImage(
		ctx context.Context,
		device *device.Device,
		payload *SendImageRequest,
	) (*SendResponse, error)
	SendList(
		ctx context.Context,
		device *device.Device,
		payload *SendListRequest,
	) (*SendResponse, error)
	SendLocation(
		ctx context.Context,
		device *device.Device,
		payload *SendLocationRequest,
	) (*SendResponse, error)
	SendSticker(
		ctx context.Context,
		device *device.Device,
		payload *SendStickerRequest,
	) (*SendResponse, error)
	SendText(
		ctx context.Context,
		device *device.Device,
		payload *SendTextRequest,
	) (*SendResponse, error)
	SendVideo(
		ctx context.Context,
		device *device.Device,
		payload *SendVideoRequest,
	) (*SendResponse, error)
}
