package gateway

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"sapasora/internal/modules/device"
	"sapasora/internal/modules/telegram"
	"sapasora/internal/modules/whatsapp"
	"sapasora/platform/support/arr"
)

var (
	// ErrUnsupportedCapability is the stable sentinel for a capability that a
	// provider or device type cannot serve.
	ErrUnsupportedCapability = errors.New("unsupported capability")
	// ErrUnknownDeviceType is returned when a gateway request has no known
	// provider to dispatch to.
	ErrUnknownDeviceType = errors.New("unknown device type")
	// ErrAdapterUnavailable is returned when the selected provider was not
	// configured. It is preferable to a nil-interface panic in the gateway.
	ErrAdapterUnavailable = errors.New("provider adapter unavailable")
)

// CapabilityError keeps the capability error machine-checkable while retaining
// useful provider/capability context for API clients and logs.
type CapabilityError struct {
	Capability string
	DeviceType device.DeviceType
}

func (e *CapabilityError) Error() string {
	return fmt.Sprintf(
		"capability %q is not supported for device type %q",
		e.Capability,
		deviceTypeName(e.DeviceType),
	)
}

func (e *CapabilityError) Unwrap() error { return ErrUnsupportedCapability }

func deviceTypeName(deviceType device.DeviceType) string {
	switch deviceType {
	case device.DeviceTypeWhatsapp:
		return "whatsapp"
	case device.DeviceTypeTelegram:
		return "telegram"
	default:
		return fmt.Sprintf("unknown(%d)", deviceType)
	}
}

func unsupportedCapability(capability string, deviceType device.DeviceType) error {
	return &CapabilityError{Capability: capability, DeviceType: deviceType}
}

func adapterUnavailable(provider string) error {
	return fmt.Errorf("%w: %s", ErrAdapterUnavailable, provider)
}

func unknownDeviceType(deviceType device.DeviceType) error {
	return fmt.Errorf("%w: %w", ErrUnknownDeviceType, unsupportedCapability("gateway", deviceType))
}

type GatewayServiceImpl struct {
	whatsappService whatsapp.WhatsappService
	telegramService telegram.TelegramService
}

// NewGatewayService creates a new Implementation of GatewayService
// @wired:provide
func NewGatewayService(
	whatsappService whatsapp.WhatsappService,
	telegramService telegram.TelegramService,
) GatewayService {
	return &GatewayServiceImpl{
		whatsappService: whatsappService,
		telegramService: telegramService,
	}
}

func isNilAdapter(adapter any) bool {
	if adapter == nil {
		return true
	}
	value := reflect.ValueOf(adapter)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func (g *GatewayServiceImpl) whatsapp() (whatsapp.WhatsappService, error) {
	if g == nil || isNilAdapter(g.whatsappService) {
		return nil, adapterUnavailable("whatsapp")
	}
	return g.whatsappService, nil
}

func (g *GatewayServiceImpl) telegram() (telegram.TelegramService, error) {
	if g == nil || isNilAdapter(g.telegramService) {
		return nil, adapterUnavailable("telegram")
	}
	return g.telegramService, nil
}

func dispatchDevice(d *device.Device) (device.DeviceType, error) {
	if d == nil {
		return 0, fmt.Errorf("%w: nil device", ErrUnknownDeviceType)
	}
	return d.Type, nil
}

func mapWhatsappUser(i whatsapp.User) User {
	return User{
		JID:          i.JID,
		VerifiedName: mapWhatsappVerifiedName(i.VerifiedName),
		Status:       i.Status,
		PictureID:    i.PictureID,
		Devices:      i.Devices,
	}
}

func mapWhatsappVerifiedName(i *whatsapp.VerifiedName) *VerifiedName {
	if i == nil {
		return nil
	}
	return &VerifiedName{
		Serial:         i.Serial,
		Issuer:         i.Issuer,
		VerifiedName:   i.VerifiedName,
		LocalizedNames: arr.Map(i.LocalizedNames, func(name *string) *string { return name }),
		IssueTime:      i.IssueTime,
	}
}

func mapTelegramUser(i telegram.User) User {
	return User{
		JID:          i.Username,
		VerifiedName: mapTelegramVerifiedName(i.VerifiedName),
		Status:       i.Status,
		PictureID:    i.PictureID,
		Devices:      i.Devices,
	}
}

func mapTelegramVerifiedName(i *telegram.VerifiedName) *VerifiedName {
	if i == nil {
		return nil
	}
	return &VerifiedName{
		Serial:         i.Serial,
		Issuer:         i.Issuer,
		VerifiedName:   i.VerifiedName,
		LocalizedNames: arr.Map(i.LocalizedNames, func(name *string) *string { return name }),
		IssueTime:      i.IssueTime,
	}
}

func mapWhatsappContextInfo(i ContextInfo) whatsapp.ContextInfo {
	return whatsapp.ContextInfo{StanzaID: i.StanzaID, Participant: i.Participant}
}

func mapTelegramContextInfo(i ContextInfo) telegram.ContextInfo {
	return telegram.ContextInfo{StanzaID: i.StanzaID, Participant: i.Participant}
}

func mapWhatsappSendResponse(i *whatsapp.SendResponse) *SendResponse {
	return &SendResponse{
		Timestamp: i.Timestamp,
		ID:        i.ID,
		ServerID:  i.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch:        i.DebugTimings.LIDFetch,
			Queue:           i.DebugTimings.Queue,
			Marshal:         i.DebugTimings.Marshal,
			GetParticipants: i.DebugTimings.GetParticipants,
			GetDevices:      i.DebugTimings.GetDevices,
			GroupEncrypt:    i.DebugTimings.GroupEncrypt,
			PeerEncrypt:     i.DebugTimings.PeerEncrypt,
			Send:            i.DebugTimings.Send,
			Resp:            i.DebugTimings.Resp,
			Retry:           i.DebugTimings.Retry,
		},
		Sender: i.Sender,
	}
}

func mapTelegramSendResponse(i *telegram.SendResponse) *SendResponse {
	return &SendResponse{
		Timestamp: i.Timestamp,
		ID:        i.ID,
		ServerID:  i.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch:        i.DebugTimings.LIDFetch,
			Queue:           i.DebugTimings.Queue,
			Marshal:         i.DebugTimings.Marshal,
			GetParticipants: i.DebugTimings.GetParticipants,
			GetDevices:      i.DebugTimings.GetDevices,
			GroupEncrypt:    i.DebugTimings.GroupEncrypt,
			PeerEncrypt:     i.DebugTimings.PeerEncrypt,
			Send:            i.DebugTimings.Send,
			Resp:            i.DebugTimings.Resp,
			Retry:           i.DebugTimings.Retry,
		},
		Sender: i.Sender,
	}
}

// CheckUser implements [GatewayService].
func (g *GatewayServiceImpl) CheckUser(ctx context.Context, d *device.Device, payload *CheckUserRequest) (*CheckUserResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.CheckUser(ctx, d, &whatsapp.CheckUserRequest{Phone: payload.Phone})
		if err != nil {
			return nil, err
		}
		return &CheckUserResponse{Users: arr.Map(response.Users, func(i whatsapp.CheckUser) CheckUser {
			return CheckUser{Query: i.Query, IsInWhatsapp: i.IsInWhatsapp, JID: i.JID, VerifiedName: i.VerifiedName}
		})}, nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.CheckUser(ctx, d, &telegram.CheckUserRequest{Phone: payload.Phone, Username: payload.Username})
		if err != nil {
			return nil, err
		}
		return &CheckUserResponse{Users: arr.Map(response.Users, func(i telegram.CheckUser) CheckUser {
			return CheckUser{Query: i.Query, IsInWhatsapp: i.IsInTelegram, JID: i.Username, VerifiedName: i.VerifiedName}
		})}, nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// Connect implements [GatewayService].
func (g *GatewayServiceImpl) Connect(ctx context.Context, d *device.Device, payload *ConnectRequest) error {
	t, err := dispatchDevice(d)
	if err != nil {
		return err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return err
		}
		return a.Connect(ctx, d, &whatsapp.ConnectRequest{Subscribe: payload.Subscribe, Immediate: payload.Immediate})
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return err
		}
		return a.Connect(ctx, d, &telegram.ConnectRequest{Subscribe: payload.Subscribe, Immediate: payload.Immediate})
	default:
		return unknownDeviceType(t)
	}
}

// Disconnect implements [GatewayService].
func (g *GatewayServiceImpl) Disconnect(ctx context.Context, d *device.Device) error {
	t, err := dispatchDevice(d)
	if err != nil {
		return err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return err
		}
		return a.Disconnect(ctx, d)
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return err
		}
		return a.Disconnect(ctx, d)
	default:
		return unknownDeviceType(t)
	}
}

// GetAvatar implements [GatewayService].
func (g *GatewayServiceImpl) GetAvatar(ctx context.Context, d *device.Device, payload *GetAvatarRequest) (*GetAvatarResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.GetAvatar(ctx, d, &whatsapp.GetAvatarRequest{Phone: payload.Phone, Preview: payload.Preview})
		if err != nil {
			return nil, err
		}
		return &GetAvatarResponse{ID: response.ID, URL: response.URL, Type: response.Type, DirectPath: response.DirectPath}, nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.GetAvatar(ctx, d, &telegram.GetAvatarRequest{Phone: payload.Phone, Username: payload.Username, Preview: payload.Preview})
		if err != nil {
			return nil, err
		}
		return &GetAvatarResponse{ID: response.ID, URL: response.URL, Type: response.Type, DirectPath: response.DirectPath}, nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// GetContacts implements [GatewayService].
func (g *GatewayServiceImpl) GetContacts(ctx context.Context, d *device.Device) (*GetContactsResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.GetContacts(ctx, d)
		if err != nil {
			return nil, err
		}
		return &GetContactsResponse{Contacts: arr.Map(response.Contacts, func(i whatsapp.Contact) Contact {
			return Contact{JID: i.JID, FullName: i.FullName, FirstName: i.FirstName, PushName: i.PushName, Found: i.Found, BusinessName: i.BusinessName, RedactedPhone: i.RedactedPhone}
		})}, nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.GetContacts(ctx, d)
		if err != nil {
			return nil, err
		}
		return &GetContactsResponse{Contacts: arr.Map(response.Contacts, func(i telegram.Contact) Contact {
			return Contact{JID: i.Username, FullName: i.FullName, FirstName: i.FirstName, PushName: i.PushName, Found: i.Found, BusinessName: i.BusinessName, RedactedPhone: i.RedactedPhone}
		})}, nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// GetStatus implements [GatewayService].
func (g *GatewayServiceImpl) GetStatus(ctx context.Context, d *device.Device) (*GetStatusResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.GetStatus(ctx, d)
		if err != nil {
			return nil, err
		}
		return &GetStatusResponse{Connected: response.Connected, LoggedIn: response.LoggedIn}, nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.GetStatus(ctx, d)
		if err != nil {
			return nil, err
		}
		return &GetStatusResponse{Connected: response.Connected, LoggedIn: response.LoggedIn}, nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// GetUser implements [GatewayService].
func (g *GatewayServiceImpl) GetUser(ctx context.Context, d *device.Device, payload *GetUserRequest) (*GetUserResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.GetUser(ctx, d, &whatsapp.GetUserRequest{Phone: payload.Phone})
		if err != nil {
			return nil, err
		}
		return &GetUserResponse{Users: arr.Map(response.Users, mapWhatsappUser)}, nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.GetUser(ctx, d, &telegram.GetUserRequest{Phone: payload.Phone, Username: payload.Username})
		if err != nil {
			return nil, err
		}
		return &GetUserResponse{Users: arr.Map(response.Users, mapTelegramUser)}, nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// Logout implements [GatewayService].
func (g *GatewayServiceImpl) Logout(ctx context.Context, d *device.Device) error {
	t, err := dispatchDevice(d)
	if err != nil {
		return err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return err
		}
		return a.Logout(ctx, d)
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return err
		}
		return a.Logout(ctx, d)
	default:
		return unknownDeviceType(t)
	}
}

// SendAudio implements [GatewayService].
func (g *GatewayServiceImpl) SendAudio(ctx context.Context, d *device.Device, payload *SendAudioRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendAudio(ctx, d, &whatsapp.SendAudioRequest{Phone: payload.Phone, Audio: payload.Audio, Caption: payload.Caption, ID: payload.ID, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendAudio(ctx, d, &telegram.SendAudioRequest{Phone: payload.Phone, Username: payload.Username, Audio: payload.Audio, Caption: payload.Caption, ID: payload.ID, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendButton implements [GatewayService].
func (g *GatewayServiceImpl) SendButton(ctx context.Context, d *device.Device, payload *SendButtonTextRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendButton(ctx, d, &whatsapp.SendButtonTextRequest{Phone: payload.Phone, Title: payload.Title, Buttons: arr.Map(payload.Buttons, func(i Button) whatsapp.Button { return whatsapp.Button{ButtonID: i.ButtonID, ButtonText: i.ButtonText} }), ID: payload.ID})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendButton(ctx, d, &telegram.SendButtonTextRequest{Phone: payload.Phone, Username: payload.Username, Title: payload.Title, Buttons: arr.Map(payload.Buttons, func(i Button) telegram.Button { return telegram.Button{ButtonID: i.ButtonID, ButtonText: i.ButtonText} }), ID: payload.ID})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendChatPresence implements [GatewayService].
func (g *GatewayServiceImpl) SendChatPresence(ctx context.Context, d *device.Device, payload *ChatPresenceRequest) error {
	t, err := dispatchDevice(d)
	if err != nil {
		return err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return err
		}
		return a.SendChatPresence(ctx, d, &whatsapp.ChatPresenceRequest{Phone: payload.Phone, State: payload.State, Media: payload.Media})
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return err
		}
		return a.SendChatPresence(ctx, d, &telegram.ChatPresenceRequest{Phone: payload.Phone, Username: payload.Username, State: payload.State, Media: payload.Media})
	default:
		return unknownDeviceType(t)
	}
}

// SendContact implements [GatewayService].
func (g *GatewayServiceImpl) SendContact(ctx context.Context, d *device.Device, payload *SendContactRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendContact(ctx, d, &whatsapp.SendContactRequest{Phone: payload.Phone, ID: payload.ID, Name: payload.Name, Vcard: payload.Vcard, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendContact(ctx, d, &telegram.SendContactRequest{Phone: payload.Phone, Username: payload.Username, ID: payload.ID, Name: payload.Name, Vcard: payload.Vcard, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendDocument implements [GatewayService].
func (g *GatewayServiceImpl) SendDocument(ctx context.Context, d *device.Device, payload *SendDocumentRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendDocument(ctx, d, &whatsapp.SendDocumentRequest{Phone: payload.Phone, Document: payload.Document, FileName: payload.FileName, ID: payload.ID, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendDocument(ctx, d, &telegram.SendDocumentRequest{Phone: payload.Phone, Username: payload.Username, Document: payload.Document, FileName: payload.FileName, ID: payload.ID, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendImage implements [GatewayService].
func (g *GatewayServiceImpl) SendImage(ctx context.Context, d *device.Device, payload *SendImageRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendImage(ctx, d, &whatsapp.SendImageRequest{Phone: payload.Phone, Image: payload.Image, Caption: payload.Caption, ID: payload.ID, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendImage(ctx, d, &telegram.SendImageRequest{Phone: payload.Phone, Username: payload.Username, Image: payload.Image, Caption: payload.Caption, ID: payload.ID, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendList implements [GatewayService].
func (g *GatewayServiceImpl) SendList(ctx context.Context, d *device.Device, payload *SendListRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		sections := arr.Map(payload.Sections, func(i Section) whatsapp.Section {
			return whatsapp.Section{Title: i.Title, Rows: arr.Map(i.Rows, func(i Row) whatsapp.Row {
				return whatsapp.Row{RowID: i.RowID, Title: i.Title, Description: i.Description}
			})}
		})
		response, err := a.SendList(ctx, d, &whatsapp.SendListRequest{Phone: payload.Phone, Title: payload.Title, Description: payload.Description, ButtonText: payload.ButtonText, FooterText: payload.FooterText, Sections: sections, ID: payload.ID})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		sections := arr.Map(payload.Sections, func(i Section) telegram.Section {
			return telegram.Section{Title: i.Title, Rows: arr.Map(i.Rows, func(i Row) telegram.Row {
				return telegram.Row{RowID: i.RowID, Title: i.Title, Description: i.Description}
			})}
		})
		response, err := a.SendList(ctx, d, &telegram.SendListRequest{Phone: payload.Phone, Username: payload.Username, Title: payload.Title, Description: payload.Description, ButtonText: payload.ButtonText, FooterText: payload.FooterText, Sections: sections, ID: payload.ID})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendLocation implements [GatewayService].
func (g *GatewayServiceImpl) SendLocation(ctx context.Context, d *device.Device, payload *SendLocationRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendLocation(ctx, d, &whatsapp.SendLocationRequest{Phone: payload.Phone, ID: payload.ID, Name: payload.Name, Latitude: payload.Latitude, Longitude: payload.Longitude, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendLocation(ctx, d, &telegram.SendLocationRequest{Phone: payload.Phone, Username: payload.Username, ID: payload.ID, Name: payload.Name, Latitude: payload.Latitude, Longitude: payload.Longitude, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendSticker implements [GatewayService].
func (g *GatewayServiceImpl) SendSticker(ctx context.Context, d *device.Device, payload *SendStickerRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendSticker(ctx, d, &whatsapp.SendStickerRequest{Phone: payload.Phone, Sticker: payload.Sticker, ID: payload.ID, PngThumbnail: payload.PngThumbnail, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendSticker(ctx, d, &telegram.SendStickerRequest{Phone: payload.Phone, Username: payload.Username, Sticker: payload.Sticker, ID: payload.ID, PngThumbnail: payload.PngThumbnail, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendText implements [GatewayService].
func (g *GatewayServiceImpl) SendText(ctx context.Context, d *device.Device, payload *SendTextRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendText(ctx, d, &whatsapp.SendTextRequest{Phone: payload.Phone, Body: payload.Body, ID: payload.ID, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendText(ctx, d, &telegram.SendTextRequest{Phone: payload.Phone, Username: payload.Username, Body: payload.Body, ID: payload.ID, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}

// SendVideo implements [GatewayService].
func (g *GatewayServiceImpl) SendVideo(ctx context.Context, d *device.Device, payload *SendVideoRequest) (*SendResponse, error) {
	t, err := dispatchDevice(d)
	if err != nil {
		return nil, err
	}
	switch t {
	case device.DeviceTypeWhatsapp:
		a, err := g.whatsapp()
		if err != nil {
			return nil, err
		}
		response, err := a.SendVideo(ctx, d, &whatsapp.SendVideoRequest{Phone: payload.Phone, Video: payload.Video, Caption: payload.Caption, ID: payload.ID, JpegThumbnail: payload.JpegThumbnail, ContextInfo: mapWhatsappContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapWhatsappSendResponse(response), nil
	case device.DeviceTypeTelegram:
		a, err := g.telegram()
		if err != nil {
			return nil, err
		}
		response, err := a.SendVideo(ctx, d, &telegram.SendVideoRequest{Phone: payload.Phone, Username: payload.Username, Video: payload.Video, Caption: payload.Caption, ID: payload.ID, JpegThumbnail: payload.JpegThumbnail, ContextInfo: mapTelegramContextInfo(payload.ContextInfo)})
		if err != nil {
			return nil, err
		}
		return mapTelegramSendResponse(response), nil
	default:
		return nil, unknownDeviceType(t)
	}
}
