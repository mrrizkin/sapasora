package telegram

import (
	"context"
	"errors"
	"sync"

	"sapasora/internal/modules/device"
)

// FakeAdapter is a deterministic TelegramService adapter for unit and contract
// tests. It records the operation name and returns an optional configured error;
// it never starts TDLib or performs network I/O.
//
// The fake intentionally returns empty successful responses because callers
// under test should assert dispatch and error handling rather than provider
// payload details. Tests that need provider payloads should use a focused fake
// for that use case.
type FakeAdapter struct {
	mu    sync.Mutex
	calls []string
	Err   error
}

var _ TelegramService = (*FakeAdapter)(nil)

// LastCall returns the most recently invoked operation, or an empty string.
func (f *FakeAdapter) LastCall() string {
	if f == nil {
		return ""
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.calls) == 0 {
		return ""
	}
	return f.calls[len(f.calls)-1]
}

// Calls returns a copy of all operations observed by the fake.
func (f *FakeAdapter) Calls() []string {
	if f == nil {
		return nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.calls...)
}

func (f *FakeAdapter) record(operation string) error {
	if f == nil {
		return errors.New("fake Telegram adapter is nil")
	}
	f.mu.Lock()
	f.calls = append(f.calls, operation)
	err := f.Err
	f.mu.Unlock()
	return err
}

func (f *FakeAdapter) CheckUser(context.Context, *device.Device, *CheckUserRequest) (*CheckUserResponse, error) {
	return &CheckUserResponse{}, f.record("check_user")
}

func (f *FakeAdapter) Connect(context.Context, *device.Device, *ConnectRequest) error {
	return f.record("connect")
}

func (f *FakeAdapter) Disconnect(context.Context, *device.Device) error {
	return f.record("disconnect")
}

func (f *FakeAdapter) GetAvatar(context.Context, *device.Device, *GetAvatarRequest) (*GetAvatarResponse, error) {
	return &GetAvatarResponse{}, f.record("get_avatar")
}

func (f *FakeAdapter) GetContacts(context.Context, *device.Device) (*GetContactsResponse, error) {
	return &GetContactsResponse{}, f.record("get_contacts")
}

func (f *FakeAdapter) GetStatus(context.Context, *device.Device) (*GetStatusResponse, error) {
	return &GetStatusResponse{}, f.record("get_status")
}

func (f *FakeAdapter) GetUser(context.Context, *device.Device, *GetUserRequest) (*GetUserResponse, error) {
	return &GetUserResponse{}, f.record("get_user")
}

func (f *FakeAdapter) Logout(context.Context, *device.Device) error {
	return f.record("logout")
}

func (f *FakeAdapter) SendAudio(context.Context, *device.Device, *SendAudioRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_audio")
}

func (f *FakeAdapter) SendButton(context.Context, *device.Device, *SendButtonTextRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_button")
}

func (f *FakeAdapter) SendChatPresence(context.Context, *device.Device, *ChatPresenceRequest) error {
	return f.record("send_chat_presence")
}

func (f *FakeAdapter) SendContact(context.Context, *device.Device, *SendContactRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_contact")
}

func (f *FakeAdapter) SendDocument(context.Context, *device.Device, *SendDocumentRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_document")
}

func (f *FakeAdapter) SendImage(context.Context, *device.Device, *SendImageRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_image")
}

func (f *FakeAdapter) SendList(context.Context, *device.Device, *SendListRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_list")
}

func (f *FakeAdapter) SendLocation(context.Context, *device.Device, *SendLocationRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_location")
}

func (f *FakeAdapter) SendSticker(context.Context, *device.Device, *SendStickerRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_sticker")
}

func (f *FakeAdapter) SendText(context.Context, *device.Device, *SendTextRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_text")
}

func (f *FakeAdapter) SendVideo(context.Context, *device.Device, *SendVideoRequest) (*SendResponse, error) {
	return &SendResponse{}, f.record("send_video")
}
