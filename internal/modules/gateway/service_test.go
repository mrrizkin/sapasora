package gateway

import (
	"context"
	"errors"
	"testing"

	"sapasora/internal/modules/device"
	"sapasora/internal/modules/telegram"
	"sapasora/internal/modules/whatsapp"
)

type fakeWhatsappService struct{ called string }

func (f *fakeWhatsappService) mark(name string) { f.called = name }
func (f *fakeWhatsappService) CheckUser(context.Context, *device.Device, *whatsapp.CheckUserRequest) (*whatsapp.CheckUserResponse, error) {
	f.mark("check_user")
	return &whatsapp.CheckUserResponse{}, nil
}
func (f *fakeWhatsappService) Connect(context.Context, *device.Device, *whatsapp.ConnectRequest) error {
	f.mark("connect")
	return nil
}
func (f *fakeWhatsappService) Disconnect(context.Context, *device.Device) error {
	f.mark("disconnect")
	return nil
}
func (f *fakeWhatsappService) GetAvatar(context.Context, *device.Device, *whatsapp.GetAvatarRequest) (*whatsapp.GetAvatarResponse, error) {
	f.mark("get_avatar")
	return &whatsapp.GetAvatarResponse{}, nil
}
func (f *fakeWhatsappService) GetContacts(context.Context, *device.Device) (*whatsapp.GetContactsResponse, error) {
	f.mark("get_contacts")
	return &whatsapp.GetContactsResponse{}, nil
}
func (f *fakeWhatsappService) GetStatus(context.Context, *device.Device) (*whatsapp.GetStatusResponse, error) {
	f.mark("get_status")
	return &whatsapp.GetStatusResponse{}, nil
}
func (f *fakeWhatsappService) GetUser(context.Context, *device.Device, *whatsapp.GetUserRequest) (*whatsapp.GetUserResponse, error) {
	f.mark("get_user")
	return &whatsapp.GetUserResponse{}, nil
}
func (f *fakeWhatsappService) Logout(context.Context, *device.Device) error {
	f.mark("logout")
	return nil
}
func (f *fakeWhatsappService) SendAudio(context.Context, *device.Device, *whatsapp.SendAudioRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_audio")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendButton(context.Context, *device.Device, *whatsapp.SendButtonTextRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_button")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendChatPresence(context.Context, *device.Device, *whatsapp.ChatPresenceRequest) error {
	f.mark("send_chat_presence")
	return nil
}
func (f *fakeWhatsappService) SendContact(context.Context, *device.Device, *whatsapp.SendContactRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_contact")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendDocument(context.Context, *device.Device, *whatsapp.SendDocumentRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_document")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendImage(context.Context, *device.Device, *whatsapp.SendImageRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_image")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendList(context.Context, *device.Device, *whatsapp.SendListRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_list")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendLocation(context.Context, *device.Device, *whatsapp.SendLocationRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_location")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendSticker(context.Context, *device.Device, *whatsapp.SendStickerRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_sticker")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendText(context.Context, *device.Device, *whatsapp.SendTextRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_text")
	return &whatsapp.SendResponse{}, nil
}
func (f *fakeWhatsappService) SendVideo(context.Context, *device.Device, *whatsapp.SendVideoRequest) (*whatsapp.SendResponse, error) {
	f.mark("send_video")
	return &whatsapp.SendResponse{}, nil
}

type fakeTelegramService struct{ called string }

func (f *fakeTelegramService) mark(name string) { f.called = name }
func (f *fakeTelegramService) CheckUser(context.Context, *device.Device, *telegram.CheckUserRequest) (*telegram.CheckUserResponse, error) {
	f.mark("check_user")
	return &telegram.CheckUserResponse{}, nil
}
func (f *fakeTelegramService) Connect(context.Context, *device.Device, *telegram.ConnectRequest) error {
	f.mark("connect")
	return nil
}
func (f *fakeTelegramService) Disconnect(context.Context, *device.Device) error {
	f.mark("disconnect")
	return nil
}
func (f *fakeTelegramService) GetAvatar(context.Context, *device.Device, *telegram.GetAvatarRequest) (*telegram.GetAvatarResponse, error) {
	f.mark("get_avatar")
	return &telegram.GetAvatarResponse{}, nil
}
func (f *fakeTelegramService) GetContacts(context.Context, *device.Device) (*telegram.GetContactsResponse, error) {
	f.mark("get_contacts")
	return &telegram.GetContactsResponse{}, nil
}
func (f *fakeTelegramService) GetStatus(context.Context, *device.Device) (*telegram.GetStatusResponse, error) {
	f.mark("get_status")
	return &telegram.GetStatusResponse{}, nil
}
func (f *fakeTelegramService) GetUser(context.Context, *device.Device, *telegram.GetUserRequest) (*telegram.GetUserResponse, error) {
	f.mark("get_user")
	return &telegram.GetUserResponse{}, nil
}
func (f *fakeTelegramService) Logout(context.Context, *device.Device) error {
	f.mark("logout")
	return nil
}
func (f *fakeTelegramService) SendAudio(context.Context, *device.Device, *telegram.SendAudioRequest) (*telegram.SendResponse, error) {
	f.mark("send_audio")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendButton(context.Context, *device.Device, *telegram.SendButtonTextRequest) (*telegram.SendResponse, error) {
	f.mark("send_button")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendChatPresence(context.Context, *device.Device, *telegram.ChatPresenceRequest) error {
	f.mark("send_chat_presence")
	return nil
}
func (f *fakeTelegramService) SendContact(context.Context, *device.Device, *telegram.SendContactRequest) (*telegram.SendResponse, error) {
	f.mark("send_contact")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendDocument(context.Context, *device.Device, *telegram.SendDocumentRequest) (*telegram.SendResponse, error) {
	f.mark("send_document")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendImage(context.Context, *device.Device, *telegram.SendImageRequest) (*telegram.SendResponse, error) {
	f.mark("send_image")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendList(context.Context, *device.Device, *telegram.SendListRequest) (*telegram.SendResponse, error) {
	f.mark("send_list")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendLocation(context.Context, *device.Device, *telegram.SendLocationRequest) (*telegram.SendResponse, error) {
	f.mark("send_location")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendSticker(context.Context, *device.Device, *telegram.SendStickerRequest) (*telegram.SendResponse, error) {
	f.mark("send_sticker")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendText(context.Context, *device.Device, *telegram.SendTextRequest) (*telegram.SendResponse, error) {
	f.mark("send_text")
	return &telegram.SendResponse{}, nil
}
func (f *fakeTelegramService) SendVideo(context.Context, *device.Device, *telegram.SendVideoRequest) (*telegram.SendResponse, error) {
	f.mark("send_video")
	return &telegram.SendResponse{}, nil
}

type gatewayOperation struct {
	name string
	run  func(*GatewayServiceImpl, *device.Device) error
}

func gatewayOperations() []gatewayOperation {
	return []gatewayOperation{
		{"check_user", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.CheckUser(context.Background(), d, &CheckUserRequest{})
			return err
		}},
		{"connect", func(g *GatewayServiceImpl, d *device.Device) error {
			return g.Connect(context.Background(), d, &ConnectRequest{})
		}},
		{"disconnect", func(g *GatewayServiceImpl, d *device.Device) error { return g.Disconnect(context.Background(), d) }},
		{"get_avatar", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.GetAvatar(context.Background(), d, &GetAvatarRequest{})
			return err
		}},
		{"get_contacts", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.GetContacts(context.Background(), d)
			return err
		}},
		{"get_status", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.GetStatus(context.Background(), d)
			return err
		}},
		{"get_user", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.GetUser(context.Background(), d, &GetUserRequest{})
			return err
		}},
		{"logout", func(g *GatewayServiceImpl, d *device.Device) error { return g.Logout(context.Background(), d) }},
		{"send_audio", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendAudio(context.Background(), d, &SendAudioRequest{})
			return err
		}},
		{"send_button", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendButton(context.Background(), d, &SendButtonTextRequest{})
			return err
		}},
		{"send_chat_presence", func(g *GatewayServiceImpl, d *device.Device) error {
			return g.SendChatPresence(context.Background(), d, &ChatPresenceRequest{})
		}},
		{"send_contact", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendContact(context.Background(), d, &SendContactRequest{})
			return err
		}},
		{"send_document", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendDocument(context.Background(), d, &SendDocumentRequest{})
			return err
		}},
		{"send_image", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendImage(context.Background(), d, &SendImageRequest{})
			return err
		}},
		{"send_list", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendList(context.Background(), d, &SendListRequest{})
			return err
		}},
		{"send_location", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendLocation(context.Background(), d, &SendLocationRequest{})
			return err
		}},
		{"send_sticker", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendSticker(context.Background(), d, &SendStickerRequest{})
			return err
		}},
		{"send_text", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendText(context.Background(), d, &SendTextRequest{})
			return err
		}},
		{"send_video", func(g *GatewayServiceImpl, d *device.Device) error {
			_, err := g.SendVideo(context.Background(), d, &SendVideoRequest{})
			return err
		}},
	}
}

func TestGatewayDispatchesEveryOperationByDeviceType(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		deviceType  device.DeviceType
		makeAdapter func() (GatewayService, func() string)
	}{
		{"whatsapp", device.DeviceTypeWhatsapp, func() (GatewayService, func() string) {
			adapter := &fakeWhatsappService{}
			return NewGatewayService(adapter, nil), func() string { return adapter.called }
		}},
		{"telegram", device.DeviceTypeTelegram, func() (GatewayService, func() string) {
			adapter := &fakeTelegramService{}
			return NewGatewayService(nil, adapter), func() string { return adapter.called }
		}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			gateway, called := testCase.makeAdapter()
			implementation := gateway.(*GatewayServiceImpl)
			for _, operation := range gatewayOperations() {
				t.Run(operation.name, func(t *testing.T) {
					if err := operation.run(implementation, &device.Device{Type: testCase.deviceType}); err != nil {
						t.Fatal(err)
					}
					if got := called(); got != operation.name {
						t.Fatalf("adapter operation = %q, want %q", got, operation.name)
					}
				})
			}
		})
	}
}

func TestGatewayReturnsStableErrorsForUnknownDeviceAndNilAdapter(t *testing.T) {
	unknown := &device.Device{Type: device.DeviceType(99)}
	_, err := (&GatewayServiceImpl{}).GetContacts(context.Background(), unknown)
	if !errors.Is(err, ErrUnknownDeviceType) || !errors.Is(err, ErrUnsupportedCapability) {
		t.Fatalf("unknown device error = %v, want both stable sentinels", err)
	}

	_, err = (&GatewayServiceImpl{}).GetContacts(context.Background(), &device.Device{Type: device.DeviceTypeTelegram})
	if !errors.Is(err, ErrAdapterUnavailable) {
		t.Fatalf("nil adapter error = %v, want ErrAdapterUnavailable", err)
	}
}
