package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync/atomic"

	"codeberg.org/mrrizkin/nihil"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var historySyncID int32

type Client struct {
	WAClient   *whatsmeow.Client
	DeviceInfo *WhatsmeowDeviceInfo
	Events     []string

	eventHandlerID uint32
	wMeow          *Whatsmeow
}

func NewClient(
	waClient *whatsmeow.Client,
	deviceInfo *WhatsmeowDeviceInfo,
	events []string,
	w *Whatsmeow,
) *Client {
	return &Client{
		WAClient:   waClient,
		DeviceInfo: deviceInfo,
		Events:     events,
		wMeow:      w,
	}
}

func (c *Client) SetClient(waClient *whatsmeow.Client) {
	c.WAClient = waClient
}

func (c *Client) SetEventHandlerID(eventHandlerID uint32) {
	c.eventHandlerID = eventHandlerID
}

func (c *Client) EventHandler(rawEvt any) {
	txtid := c.DeviceInfo.ID
	postmap := make(map[string]any)
	postmap["event"] = rawEvt
	dowebhook := 0
	path := ""

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(c.WAClient.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := c.WAClient.SendPresence(context.Background(), types.PresenceAvailable)
			if err != nil {
				c.wMeow.log.Warn("Failed to send available presence", "device", c.DeviceInfo.Name, "error", err)
			} else {
				c.wMeow.log.Info("Marked self as available", "device", c.DeviceInfo.Name)
			}
		}
	case *events.Connected, *events.PushNameSetting:
		if len(c.WAClient.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := c.WAClient.SendPresence(context.Background(), types.PresenceAvailable)
		if err != nil {
			c.wMeow.log.Warn("Failed to send available presence", "device", c.DeviceInfo.Name, "error", err)
		} else {
			c.wMeow.log.Info("Marked self as available", "device", c.DeviceInfo.Name)
		}
		c.wMeow.log.Info("Setting up status connection", "device", c.DeviceInfo.Name)
		err = c.wMeow.deviceService.SetDeviceStatusConnectedByPublicID(context.Background(), c.DeviceInfo.ID)
		if err != nil {
			c.wMeow.log.Error("Failed to set user connected", "device", c.DeviceInfo.Name, "error", err)
			return
		}
	case *events.PairSuccess:
		c.wMeow.log.Info("QR Pair Success", "device", c.DeviceInfo.Name)

		jid := evt.ID
		err = c.wMeow.deviceService.SetDeviceJIDByPublicID(context.Background(), c.DeviceInfo.ID, jid.String())
		if err != nil {
			c.wMeow.log.Error("Failed to set user jid", "device", c.DeviceInfo.Name, "error", err)
			return
		}

		c.DeviceInfo.Jid = nihil.String(jid.String())
	case *events.StreamReplaced:
		c.wMeow.log.Info("Received StreamReplaced event", "device", c.DeviceInfo.Name)
		return
	case *events.Message:
		postmap["type"] = "Message"
		dowebhook = 1
		// metaParts := []string{fmt.Sprintf("pushname: %s", evt.Info.PushName), fmt.Sprintf("timestamp: %s", evt.Info.Timestamp)}
		// if evt.Info.Type != "" {
		// 	metaParts = append(metaParts, fmt.Sprintf("type: %s", evt.Info.Type))
		// }
		// if evt.Info.Category != "" {
		// 	metaParts = append(metaParts, fmt.Sprintf("category: %s", evt.Info.Category))
		// }
		// if evt.IsViewOnce {
		// 	metaParts = append(metaParts, "view once")
		// 	metaParts = append(metaParts, "ephemeral")
		// }

		c.wMeow.log.Info("Message Received", "device", c.DeviceInfo.Name)

		// try to get Image if any
		img := evt.Message.GetImageMessage()
		if img != nil {

			// check/creates user directory for files
			userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
			_, err := os.Stat(userDirectory)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(userDirectory, 0751)
				if errDir != nil {
					c.wMeow.log.Error("Could not create user directory", "device", c.DeviceInfo.Name, "error", errDir)
					return
				}
			}

			downloadCtx, cancel := providerContext(context.Background())
			data, err := c.WAClient.Download(downloadCtx, img)
			cancel()
			if err != nil {
				c.wMeow.log.Error("Failed to download image", "device", c.DeviceInfo.Name, "error", err)
				return
			}
			path = fmt.Sprintf("%s/%s%s", userDirectory, evt.Info.ID, mediaExtension(img.GetMimetype()))
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				c.wMeow.log.Error("Failed to save image", "device", c.DeviceInfo.Name, "error", err)
				return
			}
			c.wMeow.log.Info("Image saved", "device", c.DeviceInfo.Name)
		}

		// try to get Audio if any
		audio := evt.Message.GetAudioMessage()
		if audio != nil {

			// check/creates user directory for files
			userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
			_, err := os.Stat(userDirectory)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(userDirectory, 0751)
				if errDir != nil {
					c.wMeow.log.Error("Could not create user directory", "device", c.DeviceInfo.Name, "error", errDir)
					return
				}
			}

			downloadCtx, cancel := providerContext(context.Background())
			data, err := c.WAClient.Download(downloadCtx, audio)
			cancel()
			if err != nil {
				c.wMeow.log.Error("Failed to download audio", "device", c.DeviceInfo.Name, "error", err)
				return
			}
			path = fmt.Sprintf("%s/%s%s", userDirectory, evt.Info.ID, mediaExtension(audio.GetMimetype()))
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				c.wMeow.log.Error("Failed to save audio", "device", c.DeviceInfo.Name, "error", err)
				return
			}
			c.wMeow.log.Info("Audio saved", "device", c.DeviceInfo.Name)
		}

		// try to get Document if any
		document := evt.Message.GetDocumentMessage()
		if document != nil {

			// check/creates user directory for files
			userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
			_, err := os.Stat(userDirectory)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(userDirectory, 0751)
				if errDir != nil {
					c.wMeow.log.Error("Could not create user directory", "device", c.DeviceInfo.Name, "error", errDir)
					return
				}
			}

			downloadCtx, cancel := providerContext(context.Background())
			data, err := c.WAClient.Download(downloadCtx, document)
			cancel()
			if err != nil {
				c.wMeow.log.Error("Failed to download document", "device", c.DeviceInfo.Name, "error", err)
				return
			}
			extension := filepath.Ext(document.GetFileName())
			if extension == "" {
				extension = mediaExtension(document.GetMimetype())
			}
			path = fmt.Sprintf("%s/%s%s", userDirectory, evt.Info.ID, extension)
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				c.wMeow.log.Error("Failed to save document", "device", c.DeviceInfo.Name, "error", err)
				return
			}
			c.wMeow.log.Info("Document saved", "device", c.DeviceInfo.Name)
		}
	case *events.Receipt:
		postmap["type"] = "ReadReceipt"
		dowebhook = 1
		switch evt.Type {
		case types.ReceiptTypeRead, types.ReceiptTypeReadSelf:
			c.wMeow.log.Info("Message was read", "device", c.DeviceInfo.Name)

			if evt.Type == types.ReceiptTypeRead {
				postmap["state"] = "Read"
			} else {
				postmap["state"] = "ReadSelf"
			}
		case types.ReceiptTypeDelivered:
			postmap["state"] = "Delivered"
			c.wMeow.log.Info("Message delivered", "device", c.DeviceInfo.Name)
		default:
			// Discard webhooks for inactive or other delivery types
			return
		}
	case *events.Presence:
		postmap["type"] = "Presence"
		dowebhook = 1
		if evt.Unavailable {
			postmap["state"] = "offline"
			if evt.LastSeen.IsZero() {
				c.wMeow.log.Info("User is now offline", "device", c.DeviceInfo.Name)
			} else {
				c.wMeow.log.Info("User is now offline", "device", c.DeviceInfo.Name)
			}
		} else {
			postmap["state"] = "online"
			c.wMeow.log.Info("User is now online", "device", c.DeviceInfo.Name)
		}
	case *events.HistorySync:
		postmap["type"] = "HistorySync"
		dowebhook = 1

		// check/creates user directory for files
		userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
		_, err := os.Stat(userDirectory)
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(userDirectory, 0751)
			if errDir != nil {
				c.wMeow.log.Error("Could not create user directory", "device", c.DeviceInfo.Name, "error", errDir)
				return
			}
		}

		id := atomic.AddInt32(&historySyncID, 1)
		fileName := fmt.Sprintf("%s/history-%d.json", userDirectory, id)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			c.wMeow.log.Error("Failed to open file to write history sync", "device", c.DeviceInfo.Name, "error", err)
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(evt.Data)
		if err != nil {
			c.wMeow.log.Error("Failed to write history sync", "device", c.DeviceInfo.Name, "error", err)
			return
		}
		c.wMeow.log.Info("Wrote history sync", "device", c.DeviceInfo.Name)
		_ = file.Close()
	case *events.AppState:
		c.wMeow.log.Info("App state event received", "device", c.DeviceInfo.Name)
	case *events.LoggedOut:
		c.wMeow.log.Info("Logged out", "device", c.DeviceInfo.Name)
		c.wMeow.SendKillChannel(c.DeviceInfo.ID)
		err := c.wMeow.deviceService.SetDeviceStatusDisconnectedByPublicID(context.Background(), c.DeviceInfo.ID)
		if err != nil {
			c.wMeow.log.Error("Failed to set user disconnected", "device", c.DeviceInfo.Name, "error", err)
			return
		}
	case *events.ChatPresence:
		postmap["type"] = "ChatPresence"
		dowebhook = 1
		c.wMeow.log.Info("Chat Presence received", "device", c.DeviceInfo.Name)
	case *events.CallOffer:
		c.wMeow.log.Info("Got call offer", "device", c.DeviceInfo.Name)
	case *events.CallAccept:
		c.wMeow.log.Info("Got call accept", "device", c.DeviceInfo.Name)
	case *events.CallTerminate:
		c.wMeow.log.Info("Got call terminate", "device", c.DeviceInfo.Name)
	case *events.CallOfferNotice:
		c.wMeow.log.Info("Got call offer notice", "device", c.DeviceInfo.Name)
	case *events.CallRelayLatency:
		c.wMeow.log.Info("Got call relay latency", "device", c.DeviceInfo.Name)
	default:
		c.wMeow.log.Warn("Unhandled event", "device", c.DeviceInfo.Name, "event", fmt.Sprintf("%+v", evt))
	}

	if dowebhook == 1 && c.DeviceInfo.Webhook.Valid && c.DeviceInfo.Webhook.String != "" {
		// call webhook
		webhookurl := c.DeviceInfo.Webhook.String

		if !slices.Contains(c.Events, postmap["type"].(string)) &&
			!slices.Contains(c.Events, "All") {
			c.wMeow.log.Warn(
				"Skipping webhook. Not subscribed for this type",
				"device",
				c.DeviceInfo.Name,
				"type",
				postmap["type"].(string),
			)
			return
		}

		c.wMeow.log.Info("Calling webhook", "device", c.DeviceInfo.Name)
		values, _ := json.Marshal(postmap)
		if path == "" {
			data := make(map[string]string)
			data["data"] = string(values)
			data["id"] = c.DeviceInfo.ID
			go c.wMeow.CallHook(webhookurl, data, c.DeviceInfo.ID)
		} else {
			data := make(map[string]string)
			data["data"] = string(values)
			data["id"] = c.DeviceInfo.ID
			go c.wMeow.CallHookFile(webhookurl, data, c.DeviceInfo.ID, path)
		}
	} else {
		c.wMeow.log.Warn("No webhook set for user", "device", c.DeviceInfo.Name)
	}
}
