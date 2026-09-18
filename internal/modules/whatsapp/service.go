package whatsapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sapasora/internal/modules/device"
	"sapasora/platform/logger"
	"sapasora/platform/support/arr"
	"slices"
	"strconv"
	"strings"
	"time"

	"codeberg.org/mrrizkin/nihil"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/proto/waVnameCert"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type WhatsappServiceImpl struct {
	wMeow *Whatsmeow
	log   *logger.Logger
}

// NewWhatsappService creates a new Implementation of WhatsappService
// @wired:provide
func NewWhatsappService(
	wMeow *Whatsmeow,
	log *logger.Logger,
) WhatsappService {
	return &WhatsappServiceImpl{
		wMeow: wMeow,
		log:   log,
	}
}

// CheckUser implements [WhatsappService].
func (w *WhatsappServiceImpl) CheckUser(
	ctx context.Context,
	device *device.Device,
	payload *CheckUserRequest,
) (*CheckUserResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	client, err := w.wMeow.GetClient(device.PublicID)
	if err != nil {
		return nil, err
	}

	whatsappUsers, err := client.IsOnWhatsApp(ctx, payload.Phone)
	if err != nil {
		return nil, err
	}

	users := make([]CheckUser, len(whatsappUsers))
	for i, user := range whatsappUsers {
		users[i] = CheckUser{
			Query:        user.Query,
			IsInWhatsapp: user.IsIn,
			JID:          user.JID.String(),
			VerifiedName: user.VerifiedName.Details.GetVerifiedName(),
		}
	}

	return &CheckUserResponse{
		Users: users,
	}, nil
}

// Connect implements [WhatsappService].
func (w *WhatsappServiceImpl) Connect(
	ctx context.Context,
	device *device.Device,
	payload *ConnectRequest,
) error {
	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		w.log.Warn("Device not found, creating new device", "device", device.Name)

		deviceInfo = NewDeviceInfoFromDevice(device)
		w.wMeow.deviceInfoStore.Set(device.PublicID, deviceInfo)
	}

	var subscribedEvents []string
	if len(payload.Subscribe) > 0 {
		for _, event := range payload.Subscribe {
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
	deviceInfo.Events = nihil.String(eventstring)

	w.log.Info("Attempting to connect to WhatsApp")
	w.wMeow.NewKillChannel(deviceInfo.ID)
	go w.wMeow.StartClient(context.WithoutCancel(ctx), deviceInfo, subscribedEvents)

	if payload.Immediate {
		return nil
	}

	w.log.Info("Waiting for WhatsApp connection")
	waitCtx, cancel := providerContext(ctx)
	defer cancel()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		client, err := w.wMeow.GetClient(deviceInfo.ID)
		if err == nil && client != nil && client.IsConnected() {
			return nil
		}
		select {
		case <-waitCtx.Done():
			return ErrFailedToConnect
		case <-ticker.C:
		}
	}
}

// Disconnect implements [WhatsappService].
func (w *WhatsappServiceImpl) Disconnect(
	ctx context.Context,
	device *device.Device,
) error {
	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return ErrDeviceNotFound
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return err
	}

	if !client.IsConnected() {
		w.log.Warn("Not connected to WhatsApp", "device", deviceInfo.Name)
		return nil
	}

	w.log.Info("Disconnecting from WhatsApp", "device", deviceInfo.Name)
	w.wMeow.SendKillChannel(deviceInfo.ID)

	return nil
}

// GetAvatar implements [WhatsappService].
func (w *WhatsappServiceImpl) GetAvatar(
	ctx context.Context,
	device *device.Device,
	payload *GetAvatarRequest,
) (*GetAvatarResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}
	if payload == nil {
		return nil, ErrInvalidPhoneNumber
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	jid, ok := w.wMeow.ParseJID(payload.Phone)
	if !ok {
		return nil, ErrInvalidPhoneNumber
	}

	var existingID string
	pic, err := client.GetProfilePictureInfo(ctx, jid,
		&whatsmeow.GetProfilePictureParams{
			Preview:    payload.Preview,
			ExistingID: existingID,
		})
	if err != nil {
		return nil, err
	}

	if pic == nil {
		return nil, ErrNoAvatarFound
	}

	return &GetAvatarResponse{
		URL:        pic.URL,
		ID:         pic.ID,
		Type:       pic.Type,
		DirectPath: pic.DirectPath,
	}, nil
}

// GetContacts implements [WhatsappService].
func (w *WhatsappServiceImpl) GetContacts(
	ctx context.Context,
	device *device.Device,
) (*GetContactsResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	result, err := client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	for jid, contact := range result {
		contacts = append(contacts, Contact{
			JID:          jid.String(),
			FullName:     contact.FullName,
			FirstName:    contact.FirstName,
			PushName:     contact.PushName,
			Found:        contact.Found,
			BusinessName: contact.BusinessName,

			RedactedPhone: contact.RedactedPhone,
		})
	}

	return &GetContactsResponse{
		Contacts: contacts,
	}, nil
}

// GetStatus implements [WhatsappService].
func (w *WhatsappServiceImpl) GetStatus(
	ctx context.Context,
	device *device.Device,
) (*GetStatusResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	isLoggedIn := client.IsLoggedIn()
	isConnected := client.IsConnected()

	if isConnected {
		w.wMeow.deviceService.SetDeviceStatusConnected(ctx, device)
	}

	return &GetStatusResponse{
		Connected: isConnected,
		LoggedIn:  isLoggedIn,
	}, nil
}

// GetUser implements [WhatsappService].
func (w *WhatsappServiceImpl) GetUser(
	ctx context.Context,
	device *device.Device,
	payload *GetUserRequest,
) (*GetUserResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	var jids []types.JID
	for _, phone := range payload.Phone {
		if jid, ok := w.wMeow.ParseJID(phone); ok {
			jids = append(jids, jid)
			continue
		}

		w.log.Warn("Invalid phone number", "phone", phone)
	}

	userInfo, err := client.GetUserInfo(ctx, jids)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	for jid, user := range userInfo {
		users = append(users, User{
			JID: jid.String(),
			VerifiedName: &VerifiedName{
				Serial:       user.VerifiedName.Details.Serial,
				Issuer:       user.VerifiedName.Details.Issuer,
				VerifiedName: user.VerifiedName.Details.VerifiedName,
				LocalizedNames: arr.Map(
					user.VerifiedName.Details.LocalizedNames,
					func(i *waVnameCert.LocalizedName) *string { return i.VerifiedName },
				),
				IssueTime: user.VerifiedName.Details.IssueTime,
			},
			Status:    user.Status,
			PictureID: user.PictureID,
			Devices:   arr.Map(user.Devices, func(i types.JID) string { return i.String() }),
		})
	}

	return &GetUserResponse{
		Users: users,
	}, nil
}

// Logout implements [WhatsappService].
func (w *WhatsappServiceImpl) Logout(
	ctx context.Context,
	device *device.Device,
) error {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return ErrDeviceNotFound
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return err
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	if !client.IsLoggedIn() {
		return ErrNotLoggedIn
	}

	if err := client.Logout(ctx); err != nil {
		return err
	}

	w.wMeow.SendKillChannel(deviceInfo.ID)

	return nil
}

// SendAudio implements [WhatsappService].
func (w *WhatsappServiceImpl) SendAudio(
	ctx context.Context,
	device *device.Device,
	payload *SendAudioRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	if payload == nil {
		return nil, ErrEmptyBody
	}

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	var uploaded whatsmeow.UploadResponse

	dataURL, err := decodeMediaDataURL(payload.Audio, "audio/ogg")
	if err != nil {
		return nil, errors.New(
			"audio data should be a valid base64 data URL with MIME type audio/ogg",
		)
	}

	filedata := dataURL.Data
	uploaded, err = client.Upload(ctx, filedata, whatsmeow.MediaAudio)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	ptt := true
	mime := "audio/ogg; codecs=opus"

	msg := &waE2E.Message{AudioMessage: &waE2E.AudioMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      &mime,
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
		PTT:           &ptt,
	}}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendButton implements [WhatsappService].
func (w *WhatsappServiceImpl) SendButton(
	ctx context.Context,
	device *device.Device,
	payload *SendButtonTextRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, ok := w.wMeow.ParseJID(payload.Phone)
	if !ok {
		return nil, ErrInvalidPhoneNumber
	}

	buttons := make([]*waE2E.ButtonsMessage_Button, len(payload.Buttons))
	for i, button := range payload.Buttons {
		buttons[i] = &waE2E.ButtonsMessage_Button{
			ButtonID: proto.String(button.ButtonID),
			ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{
				DisplayText: proto.String(button.ButtonText),
			},
			Type:           waE2E.ButtonsMessage_Button_RESPONSE.Enum(),
			NativeFlowInfo: &waE2E.ButtonsMessage_Button_NativeFlowInfo{},
		}
	}

	msg2 := &waE2E.ButtonsMessage{
		ContentText: proto.String(payload.Title),
		HeaderType:  waE2E.ButtonsMessage_EMPTY.Enum(),
		Buttons:     buttons,
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	response, err := client.SendMessage(ctx, recipient, &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				ButtonsMessage: msg2,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendChatPresence implements [WhatsappService].
func (w *WhatsappServiceImpl) SendChatPresence(
	ctx context.Context,
	device *device.Device,
	payload *ChatPresenceRequest,
) error {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return ErrDeviceNotFound
	}

	recipient, ok := w.wMeow.ParseJID(payload.Phone)
	if !ok {
		return ErrInvalidPhoneNumber
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return err
	}

	return client.SendChatPresence(
		ctx,
		recipient,
		types.ChatPresence(payload.State),
		types.ChatPresenceMedia(payload.Media),
	)
}

// SendContact implements [WhatsappService].
func (w *WhatsappServiceImpl) SendContact(
	ctx context.Context,
	device *device.Device,
	payload *SendContactRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	msg := &waE2E.Message{ContactMessage: &waE2E.ContactMessage{
		DisplayName: &payload.Name,
		Vcard:       &payload.Vcard,
	}}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendDocument implements [WhatsappService].
func (w *WhatsappServiceImpl) SendDocument(
	ctx context.Context,
	device *device.Device,
	payload *SendDocumentRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	if payload == nil {
		return nil, ErrEmptyBody
	}

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	if !validDocumentFilename(payload.FileName) {
		return nil, ErrInvalidFilename
	}

	var uploaded whatsmeow.UploadResponse
	dataURL, err := decodeMediaDataURL(payload.Document, "application/octet-stream")
	if err != nil {
		return nil, errors.New(
			"document data should be a valid base64 data URL with MIME type application/octet-stream",
		)
	}

	filedata := dataURL.Data
	uploaded, err = client.Upload(ctx, filedata, whatsmeow.MediaDocument)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	msg := &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{
		URL:           proto.String(uploaded.URL),
		FileName:      &payload.FileName,
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendImage implements [WhatsappService].
func (w *WhatsappServiceImpl) SendImage(
	ctx context.Context,
	device *device.Device,
	payload *SendImageRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	if payload == nil {
		return nil, ErrEmptyBody
	}

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	var uploaded whatsmeow.UploadResponse
	dataURL, err := decodeMediaDataURL(payload.Image, "image/*")
	if err != nil {
		return nil, errors.New(
			"image data should be a valid base64 data URL with an image MIME type",
		)
	}

	filedata := dataURL.Data
	uploaded, err = client.Upload(ctx, filedata, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	msg := &waE2E.Message{ImageMessage: &waE2E.ImageMessage{
		Caption:       proto.String(payload.Caption),
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendList implements [WhatsappService].
func (w *WhatsappServiceImpl) SendList(
	ctx context.Context,
	device *device.Device,
	payload *SendListRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, ok := w.wMeow.ParseJID(payload.Phone)
	if !ok {
		return nil, ErrInvalidPhoneNumber
	}

	sections := make([]*waE2E.ListMessage_Section, len(payload.Sections))
	for i, section := range payload.Sections {
		rows := make([]*waE2E.ListMessage_Row, len(section.Rows))
		for j, row := range section.Rows {
			var idText string
			if row.RowID == "" {
				idText = strconv.Itoa(j + 1)
			} else {
				idText = row.RowID
			}
			rows[j] = &waE2E.ListMessage_Row{
				RowID:       proto.String(idText),
				Title:       proto.String(row.Title),
				Description: proto.String(row.Description),
			}
		}
		sections[i] = &waE2E.ListMessage_Section{
			Title: proto.String(section.Title),
			Rows:  rows,
		}
	}

	msg1 := &waE2E.ListMessage{
		Title:       proto.String(payload.Title),
		Description: proto.String(payload.Description),
		ButtonText:  proto.String(payload.ButtonText),
		ListType:    waE2E.ListMessage_SINGLE_SELECT.Enum(),
		Sections:    sections,
		FooterText:  proto.String(payload.FooterText),
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	response, err := client.SendMessage(ctx, recipient, &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				ListMessage: msg1,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendLocation implements [WhatsappService].
func (w *WhatsappServiceImpl) SendLocation(
	ctx context.Context,
	device *device.Device,
	payload *SendLocationRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	msg := &waE2E.Message{LocationMessage: &waE2E.LocationMessage{
		DegreesLatitude:  &payload.Latitude,
		DegreesLongitude: &payload.Longitude,
		Name:             &payload.Name,
	}}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendSticker implements [WhatsappService].
func (w *WhatsappServiceImpl) SendSticker(
	ctx context.Context,
	device *device.Device,
	payload *SendStickerRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	if payload == nil {
		return nil, ErrEmptyBody
	}

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	var uploaded whatsmeow.UploadResponse
	dataURL, err := decodeMediaDataURL(payload.Sticker)
	if err != nil {
		return nil, errors.New(
			"sticker data should be a valid base64 data URL",
		)
	}

	filedata := dataURL.Data
	uploaded, err = client.Upload(ctx, filedata, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	msg := &waE2E.Message{StickerMessage: &waE2E.StickerMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
		PngThumbnail:  payload.PngThumbnail,
	}}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendText implements [WhatsappService].
func (w *WhatsappServiceImpl) SendText(
	ctx context.Context,
	device *device.Device,
	payload *SendTextRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &payload.Body,
		},
	}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}

// SendVideo implements [WhatsappService].
func (w *WhatsappServiceImpl) SendVideo(
	ctx context.Context,
	device *device.Device,
	payload *SendVideoRequest,
) (*SendResponse, error) {
	ctx, cancel := providerContext(ctx)
	defer cancel()

	if payload == nil {
		return nil, ErrEmptyBody
	}

	deviceInfo, found := w.wMeow.deviceInfoStore.Get(device.PublicID)
	if !found {
		return nil, ErrDeviceNotFound
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaID,
		payload.ContextInfo.Participant,
		w.wMeow,
	)
	if err != nil {
		return nil, err
	}

	var uploaded whatsmeow.UploadResponse
	dataURL, err := decodeMediaDataURL(payload.Video, "video/*")
	if err != nil {
		return nil, errors.New(
			"video data should be a valid base64 data URL with a video MIME type",
		)
	}

	client, err := w.wMeow.GetClient(deviceInfo.ID)
	if err != nil {
		return nil, err
	}

	filedata := dataURL.Data
	uploaded, err = client.Upload(ctx, filedata, whatsmeow.MediaVideo)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	msg := &waE2E.Message{VideoMessage: &waE2E.VideoMessage{
		Caption:       proto.String(payload.Caption),
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
		JPEGThumbnail: payload.JpegThumbnail,
	}}

	if payload.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*payload.ContextInfo.StanzaID),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	response, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch.Milliseconds(),
			Queue:    response.DebugTimings.Queue.Milliseconds(),

			Marshal:         response.DebugTimings.Marshal.Milliseconds(),
			GetParticipants: response.DebugTimings.GetParticipants.Milliseconds(),
			GetDevices:      response.DebugTimings.GetDevices.Milliseconds(),
			GroupEncrypt:    response.DebugTimings.GroupEncrypt.Milliseconds(),
			PeerEncrypt:     response.DebugTimings.PeerEncrypt.Milliseconds(),

			Send:  response.DebugTimings.Send.Milliseconds(),
			Resp:  response.DebugTimings.Resp.Milliseconds(),
			Retry: response.DebugTimings.Retry.Milliseconds(),
		},
		Sender: response.Sender.String(),
	}, nil
}
