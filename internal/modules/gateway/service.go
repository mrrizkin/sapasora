package gateway

import (
	"context"
	"sapasora/internal/modules/device"
	"sapasora/internal/modules/telegram"
	"sapasora/internal/modules/whatsapp"
	"sapasora/platform/support/arr"
	"errors"
)

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

// CheckUser implements [GatewayService].
func (g *GatewayServiceImpl) CheckUser(
	ctx context.Context,
	d *device.Device,
	payload *CheckUserRequest,
) (*CheckUserResponse, error) {
	switch d.Type {
	case device.DeviceTypeWhatsapp:
		response, err := g.whatsappService.CheckUser(ctx, d, &whatsapp.CheckUserRequest{
			Phone: payload.Phone,
		})
		if err != nil {
			return nil, err
		}

		return &CheckUserResponse{
			Users: arr.Map(response.Users, func(i whatsapp.CheckUser) CheckUser {
				return CheckUser{
					Query:        i.Query,
					IsInWhatsapp: i.IsInWhatsapp,
					JID:          i.JID,
					VerifiedName: i.VerifiedName,
				}
			}),
		}, nil
	case device.DeviceTypeTelegram:
		response, err := g.telegramService.CheckUser(ctx, d, &telegram.CheckUserRequest{
			Phone:    payload.Phone,
			Username: payload.Username,
		})
		if err != nil {
			return nil, err
		}
		return &CheckUserResponse{
			Users: arr.Map(response.Users, func(i telegram.CheckUser) CheckUser {
				return CheckUser{
					Query:        i.Query,
					IsInWhatsapp: i.IsInTelegram,
					JID:          i.Username,
					VerifiedName: i.VerifiedName,
				}
			}),
		}, nil
	default:
		return nil, errors.New("unknown device type")
	}
}

// Connect implements [GatewayService].
func (g *GatewayServiceImpl) Connect(
	ctx context.Context,
	d *device.Device,
	payload *ConnectRequest,
) error {
	switch d.Type {
	case device.DeviceTypeWhatsapp:
		return g.whatsappService.Connect(ctx, d, &whatsapp.ConnectRequest{
			Subscribe: payload.Subscribe,
			Immediate: payload.Immediate,
		})
	case device.DeviceTypeTelegram:
		return g.telegramService.Connect(ctx, d, &telegram.ConnectRequest{
			Subscribe: payload.Subscribe,
			Immediate: payload.Immediate,
		})
	default:
		return errors.New("unknown device type")
	}
}

// Disconnect implements [GatewayService].
func (g *GatewayServiceImpl) Disconnect(
	ctx context.Context,
	d *device.Device,
) error {
	switch d.Type {
	case device.DeviceTypeWhatsapp:
		return g.whatsappService.Disconnect(ctx, d)
	case device.DeviceTypeTelegram:
		return g.telegramService.Disconnect(ctx, d)
	default:
		return errors.New("unknown device type")
	}
}

// GetAvatar implements [GatewayService].
func (g *GatewayServiceImpl) GetAvatar(
	ctx context.Context,
	d *device.Device,
	payload *GetAvatarRequest,
) (*GetAvatarResponse, error) {
	switch d.Type {
	case device.DeviceTypeWhatsapp:
		response, err := g.whatsappService.GetAvatar(ctx, d, &whatsapp.GetAvatarRequest{
			Phone:   payload.Phone,
			Preview: payload.Preview,
		})
		if err != nil {
			return nil, err
		}

		return &GetAvatarResponse{
			ID:         response.ID,
			URL:        response.URL,
			Type:       response.Type,
			DirectPath: response.DirectPath,
		}, nil
	case device.DeviceTypeTelegram:
		response, err := g.telegramService.GetAvatar(ctx, d, &telegram.GetAvatarRequest{
			Phone:    payload.Phone,
			Username: payload.Username,
			Preview:  payload.Preview,
		})
		if err != nil {
			return nil, err
		}

		return &GetAvatarResponse{
			ID:         response.ID,
			URL:        response.URL,
			Type:       response.Type,
			DirectPath: response.DirectPath,
		}, nil
	default:
		return nil, errors.New("device type not supported")
	}
}

// GetContacts implements [GatewayService].
func (g *GatewayServiceImpl) GetContacts(
	ctx context.Context,
	device *device.Device,
) (*GetContactsResponse, error) {
	response, err := g.whatsappService.GetContacts(ctx, device)
	if err != nil {
		return nil, err
	}

	return &GetContactsResponse{
		Contacts: arr.Map(response.Contacts, func(i whatsapp.Contact) Contact {
			return Contact{
				JID:          i.JID,
				FullName:     i.FullName,
				FirstName:    i.FirstName,
				PushName:     i.PushName,
				Found:        i.Found,
				BusinessName: i.BusinessName,

				RedactedPhone: i.RedactedPhone,
			}
		}),
	}, nil
}

// GetStatus implements [GatewayService].
func (g *GatewayServiceImpl) GetStatus(
	ctx context.Context,
	d *device.Device,
) (*GetStatusResponse, error) {
	switch d.Type {
	case device.DeviceTypeWhatsapp:

		response, err := g.whatsappService.GetStatus(ctx, d)
		if err != nil {
			return nil, err
		}

		return &GetStatusResponse{
			Connected: response.Connected,
			LoggedIn:  response.LoggedIn,
		}, nil

	case device.DeviceTypeTelegram:

		response, err := g.telegramService.GetStatus(ctx, d)
		if err != nil {
			return nil, err
		}

		return &GetStatusResponse{
			Connected: response.Connected,
			LoggedIn:  response.LoggedIn,
		}, nil
	default:
		return nil, errors.New("unknown device type")
	}
}

// GetUser implements [GatewayService].
func (g *GatewayServiceImpl) GetUser(
	ctx context.Context,
	device *device.Device,
	payload *GetUserRequest,
) (*GetUserResponse, error) {
	response, err := g.whatsappService.GetUser(ctx, device, &whatsapp.GetUserRequest{
		Phone: payload.Phone,
	})
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{
		Users: arr.Map(response.Users, func(i whatsapp.User) User {
			return User{
				JID: i.JID,
				VerifiedName: &VerifiedName{
					Serial:       i.VerifiedName.Serial,
					Issuer:       i.VerifiedName.Issuer,
					VerifiedName: i.VerifiedName.VerifiedName,
					LocalizedNames: arr.Map(
						i.VerifiedName.LocalizedNames,
						func(i *string) *string { return i },
					),
					IssueTime: i.VerifiedName.IssueTime,
				},
				Status:    i.Status,
				PictureID: i.PictureID,
				Devices:   arr.Map(i.Devices, func(i string) string { return i }),
			}
		}),
	}, nil
}

// Logout implements [GatewayService].
func (g *GatewayServiceImpl) Logout(
	ctx context.Context,
	device *device.Device,
) error {
	return g.whatsappService.Logout(ctx, device)
}

// SendAudio implements [GatewayService].
func (g *GatewayServiceImpl) SendAudio(
	ctx context.Context,
	device *device.Device,
	payload *SendAudioRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendAudio(ctx, device, &whatsapp.SendAudioRequest{
		ID:      payload.ID,
		Phone:   payload.Phone,
		Audio:   payload.Audio,
		Caption: payload.Caption,
		ContextInfo: whatsapp.ContextInfo{
			StanzaID:    payload.ContextInfo.StanzaID,
			Participant: payload.ContextInfo.Participant,
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
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendButton implements [GatewayService].
func (g *GatewayServiceImpl) SendButton(
	ctx context.Context,
	device *device.Device,
	payload *SendButtonTextRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendButton(ctx, device, &whatsapp.SendButtonTextRequest{
		Phone: payload.Phone,
		Title: payload.Title,
		Buttons: arr.Map(
			payload.Buttons,
			func(i Button) whatsapp.Button {
				return whatsapp.Button{
					ButtonID:   i.ButtonID,
					ButtonText: i.ButtonText,
				}
			},
		),
		ID: payload.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendChatPresence implements [GatewayService].
func (g *GatewayServiceImpl) SendChatPresence(
	ctx context.Context,
	device *device.Device,
	payload *ChatPresenceRequest,
) error {
	return g.whatsappService.SendChatPresence(ctx, device, &whatsapp.ChatPresenceRequest{
		Phone: payload.Phone,
		State: payload.State,
		Media: payload.Media,
	})
}

// SendContact implements [GatewayService].
func (g *GatewayServiceImpl) SendContact(
	ctx context.Context,
	device *device.Device,
	payload *SendContactRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendContact(ctx, device, &whatsapp.SendContactRequest{
		Phone: payload.Phone,
		Name:  payload.Name,
		Vcard: payload.Vcard,
		ID:    payload.ID,
		ContextInfo: whatsapp.ContextInfo{
			StanzaID:    payload.ContextInfo.StanzaID,
			Participant: payload.ContextInfo.Participant,
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
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendDocument implements [GatewayService].
func (g *GatewayServiceImpl) SendDocument(
	ctx context.Context,
	device *device.Device,
	payload *SendDocumentRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendDocument(ctx, device, &whatsapp.SendDocumentRequest{
		Phone:    payload.Phone,
		Document: payload.Document,
		FileName: payload.FileName,
		ID:       payload.ID,
		ContextInfo: whatsapp.ContextInfo{
			StanzaID:    payload.ContextInfo.StanzaID,
			Participant: payload.ContextInfo.Participant,
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
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendImage implements [GatewayService].
func (g *GatewayServiceImpl) SendImage(
	ctx context.Context,
	device *device.Device,
	payload *SendImageRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendImage(ctx, device, &whatsapp.SendImageRequest{
		Phone:   payload.Phone,
		Image:   payload.Image,
		Caption: payload.Caption,
		ID:      payload.ID,
		ContextInfo: whatsapp.ContextInfo{
			StanzaID:    payload.ContextInfo.StanzaID,
			Participant: payload.ContextInfo.Participant,
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
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendList implements [GatewayService].
func (g *GatewayServiceImpl) SendList(
	ctx context.Context,
	device *device.Device,
	payload *SendListRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendList(ctx, device, &whatsapp.SendListRequest{
		Phone:       payload.Phone,
		Title:       payload.Title,
		Description: payload.Description,
		ButtonText:  payload.ButtonText,
		FooterText:  payload.FooterText,
		Sections: arr.Map(
			payload.Sections,
			func(i Section) whatsapp.Section {
				return whatsapp.Section{
					Title: i.Title,
					Rows: arr.Map(
						i.Rows,
						func(i Row) whatsapp.Row {
							return whatsapp.Row{
								RowID:       i.RowID,
								Title:       i.Title,
								Description: i.Description,
							}
						},
					),
				}
			},
		),
		ID: payload.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SendResponse{
		Timestamp: response.Timestamp,
		ID:        response.ID,
		ServerID:  response.ServerID,
		DebugTimings: MessageDebugTimings{
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendLocation implements [GatewayService].
func (g *GatewayServiceImpl) SendLocation(
	ctx context.Context,
	device *device.Device,
	payload *SendLocationRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendLocation(ctx, device, &whatsapp.SendLocationRequest{
		Phone:     payload.Phone,
		ID:        payload.ID,
		Name:      payload.Name,
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
		ContextInfo: whatsapp.ContextInfo{
			StanzaID:    payload.ContextInfo.StanzaID,
			Participant: payload.ContextInfo.Participant,
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
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendSticker implements [GatewayService].
func (g *GatewayServiceImpl) SendSticker(
	ctx context.Context,
	device *device.Device,
	payload *SendStickerRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendSticker(ctx, device, &whatsapp.SendStickerRequest{
		Phone:        payload.Phone,
		Sticker:      payload.Sticker,
		ID:           payload.ID,
		PngThumbnail: payload.PngThumbnail,
		ContextInfo: whatsapp.ContextInfo{
			StanzaID:    payload.ContextInfo.StanzaID,
			Participant: payload.ContextInfo.Participant,
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
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}

// SendText implements [GatewayService].
func (g *GatewayServiceImpl) SendText(
	ctx context.Context,
	d *device.Device,
	payload *SendTextRequest,
) (*SendResponse, error) {
	switch d.Type {
	case device.DeviceTypeWhatsapp:
		response, err := g.whatsappService.SendText(ctx, d, &whatsapp.SendTextRequest{
			Phone: payload.Phone,
			Body:  payload.Body,
			ID:    payload.ID,
			ContextInfo: whatsapp.ContextInfo{
				StanzaID:    payload.ContextInfo.StanzaID,
				Participant: payload.ContextInfo.Participant,
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
				LIDFetch: response.DebugTimings.LIDFetch,
				Queue:    response.DebugTimings.Queue,

				Marshal:         response.DebugTimings.Marshal,
				GetParticipants: response.DebugTimings.GetParticipants,
				GetDevices:      response.DebugTimings.GetDevices,
				GroupEncrypt:    response.DebugTimings.GroupEncrypt,
				PeerEncrypt:     response.DebugTimings.PeerEncrypt,

				Send:  response.DebugTimings.Send,
				Resp:  response.DebugTimings.Resp,
				Retry: response.DebugTimings.Retry,
			},
			Sender: response.Sender,
		}, nil

	case device.DeviceTypeTelegram:
		response, err := g.telegramService.SendText(ctx, d, &telegram.SendTextRequest{
			Phone:    payload.Phone,
			Username: payload.Username,
			Body:     payload.Body,
			ID:       payload.ID,
			ContextInfo: telegram.ContextInfo{
				StanzaID:    payload.ContextInfo.StanzaID,
				Participant: payload.ContextInfo.Participant,
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
				LIDFetch: response.DebugTimings.LIDFetch,
				Queue:    response.DebugTimings.Queue,

				Marshal:         response.DebugTimings.Marshal,
				GetParticipants: response.DebugTimings.GetParticipants,
				GetDevices:      response.DebugTimings.GetDevices,
				GroupEncrypt:    response.DebugTimings.GroupEncrypt,
				PeerEncrypt:     response.DebugTimings.PeerEncrypt,

				Send:  response.DebugTimings.Send,
				Resp:  response.DebugTimings.Resp,
				Retry: response.DebugTimings.Retry,
			},
			Sender: response.Sender,
		}, nil
	default:
		return nil, errors.New("unknown device type")
	}
}

// SendVideo implements [GatewayService].
func (g *GatewayServiceImpl) SendVideo(
	ctx context.Context,
	device *device.Device,
	payload *SendVideoRequest,
) (*SendResponse, error) {
	response, err := g.whatsappService.SendVideo(ctx, device, &whatsapp.SendVideoRequest{
		Phone:         payload.Phone,
		Video:         payload.Video,
		Caption:       payload.Caption,
		ID:            payload.ID,
		JpegThumbnail: payload.JpegThumbnail,
		ContextInfo: whatsapp.ContextInfo{
			StanzaID:    payload.ContextInfo.StanzaID,
			Participant: payload.ContextInfo.Participant,
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
			LIDFetch: response.DebugTimings.LIDFetch,
			Queue:    response.DebugTimings.Queue,

			Marshal:         response.DebugTimings.Marshal,
			GetParticipants: response.DebugTimings.GetParticipants,
			GetDevices:      response.DebugTimings.GetDevices,
			GroupEncrypt:    response.DebugTimings.GroupEncrypt,
			PeerEncrypt:     response.DebugTimings.PeerEncrypt,

			Send:  response.DebugTimings.Send,
			Resp:  response.DebugTimings.Resp,
			Retry: response.DebugTimings.Retry,
		},
		Sender: response.Sender,
	}, nil
}
