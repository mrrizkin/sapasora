package whatsapp

import (
	"time"
)

type CheckUserRequest struct {
	Phone []string `json:"phone"`
} // @name whatsapp.CheckUserRequest

type CheckUser struct {
	Query        string `json:"query"`
	IsInWhatsapp bool   `json:"is_in_whatsapp"`
	JID          string `json:"jid"`
	VerifiedName string `json:"verified_name"`
} // @name whatsapp.CheckUser

type CheckUserResponse struct {
	Users []CheckUser `json:"users"`
} // @name whatsapp.CheckUserResponse

type ConnectRequest struct {
	Subscribe []string `json:"subscribe"`
	Immediate bool     `json:"immediate"`
} // @name whatsapp.ConnectRequest

type GetAvatarRequest struct {
	Phone   string `json:"phone"`
	Preview bool   `json:"preview"`
} // @name whatsapp.GetAvatarRequest

type GetAvatarResponse struct {
	URL        string `json:"url"`
	ID         string `json:"id"`
	Type       string `json:"type"`
	DirectPath string `json:"direct_path"`
} // @name whatsapp.GetAvatarResponse

type Contact struct {
	Found         bool   `json:"found"`
	JID           string `json:"jid"`
	FullName      string `json:"full_name"`
	FirstName     string `json:"first_name"`
	PushName      string `json:"push_name"`
	BusinessName  string `json:"business_name"`
	RedactedPhone string `json:"redacted_phone"`
} // @name whatsapp.Contact

type GetContactsResponse struct {
	Contacts []Contact `json:"contacts"`
} // @name whatsapp.GetContactsResponse

type GetStatusResponse struct {
	Connected bool `json:"connected"`
	LoggedIn  bool `json:"logged_in"`
} // @name whatsapp.GetStatusResponse

type GetUserRequest struct {
	Phone []string `json:"phone"`
} // @name whatsapp.GetUserRequest

type VerifiedName struct {
	Serial         *uint64   `json:"serial,omitempty"`
	Issuer         *string   `json:"issuer,omitempty"`
	VerifiedName   *string   `json:"name,omitempty"`
	LocalizedNames []*string `json:"localize_names,omitempty"`
	IssueTime      *uint64   `json:"issue_time,omitempty"`
} // @name whatsapp.VerifiedName

type User struct {
	JID          string        `json:"jid"`
	VerifiedName *VerifiedName `json:"verified_name"`
	Status       string        `json:"status"`
	PictureID    string        `json:"picture_id"`
	Devices      []string      `json:"devices"`
} // @name whatsapp.User

type GetUserResponse struct {
	Users []User `json:"users"`
} // @name whatsapp.GetUserResponse

type ContextInfo struct {
	StanzaID    *string `json:"stanza_id,omitempty"`
	Participant *string `json:"participant,omitempty"`
} // @name whatsapp.ContextInfo

type SendAudioRequest struct {
	Phone   string `json:"phone"`
	Audio   string `json:"audio"`
	Caption string `json:"caption"`
	ID      string `json:"id"`

	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendAudioRequest

type Button struct {
	ButtonID   string `json:"button_id"`
	ButtonText string `json:"button_text"`
} // @name whatsapp.Button

type SendButtonTextRequest struct {
	Phone   string   `json:"phone"`
	Title   string   `json:"title"`
	Buttons []Button `json:"buttons"`
	ID      string   `json:"id"`
} // @name whatsapp.SendButtonTextRequest

type ChatPresenceRequest struct {
	Phone string `json:"phone"`
	State string `json:"state"`
	Media string `json:"media"`
} // @name whatsapp.ChatPresenceRequest

type SendContactRequest struct {
	Phone       string      `json:"phone"`
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Vcard       string      `json:"vcard"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendContactRequest

type SendDocumentRequest struct {
	Phone       string      `json:"phone"`
	Document    string      `json:"document"`
	FileName    string      `json:"filename"`
	ID          string      `json:"id"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendDocumentRequest

type SendImageRequest struct {
	Phone       string      `json:"phone"`
	Image       string      `json:"image"`
	Caption     string      `json:"caption"`
	ID          string      `json:"id"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendImageRequest

type Row struct {
	RowID       string `json:"row_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
} // @name whatsapp.Row

type Section struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
} // @name whatsapp.Section

type SendListRequest struct {
	Phone       string    `json:"phone"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ButtonText  string    `json:"button_text"`
	FooterText  string    `json:"footer_text"`
	Sections    []Section `json:"sections"`
	ID          string    `json:"id"`
} // @name whatsapp.SendListRequest

type SendLocationRequest struct {
	Phone       string      `json:"phone"`
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendLocationRequest

type SendStickerRequest struct {
	Phone        string      `json:"phone"`
	Sticker      string      `json:"sticker"`
	ID           string      `json:"id"`
	PngThumbnail []byte      `json:"png_thumbnail"`
	ContextInfo  ContextInfo `json:"context_info"`
} // @name whatsapp.SendStickerRequest

type SendTextRequest struct {
	Phone       string      `json:"phone"`
	Body        string      `json:"body"`
	ID          string      `json:"id"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendTextRequest

type SendVideoRequest struct {
	Phone         string      `json:"phone"`
	Video         string      `json:"video"`
	Caption       string      `json:"caption"`
	ID            string      `json:"id"`
	JpegThumbnail []byte      `json:"jpeg_thumbnail"`
	ContextInfo   ContextInfo `json:"context_info"`
} // @name whatsapp.SendVideoRequest

type MessageDebugTimings struct {
	LIDFetch int64 `json:"lid_fetch"`
	Queue    int64 `json:"queue"`

	Marshal         int64 `json:"marshal"`
	GetParticipants int64 `json:"get_participants"`
	GetDevices      int64 `json:"get_devices"`
	GroupEncrypt    int64 `json:"group_encrypt"`
	PeerEncrypt     int64 `json:"peer_encrypt"`

	Send  int64 `json:"send"`
	Resp  int64 `json:"resp"`
	Retry int64 `json:"retry"`
} // @name whatsapp.MessageDebugTimings

type SendResponse struct {
	Timestamp    time.Time           `json:"timestamp"`
	ID           string              `json:"id"`
	ServerID     int                 `json:"server_id"`
	DebugTimings MessageDebugTimings `json:"debug_timings"`
	Sender       string              `json:"sender"`
} // @name whatsapp.SendResponse
