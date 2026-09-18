package gateway

import (
	"time"

	"codeberg.org/mrrizkin/nihil"
)

type GetQRResponse struct {
	QRCode nihil.NilString `json:"qr_code"`
} // @name gateway.GetQRResponse

type CheckUserRequest struct {
	Phone    []string `json:"phone" validate:"required_without=Username,dive,min=1,max=255"`
	Username []string `json:"username" validate:"required_without=Phone,dive,min=1,max=255"`
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
	Subscribe []string `json:"subscribe" validate:"omitempty,max=100,dive,min=1,max=100"`
	Immediate bool     `json:"immediate"`
} // @name whatsapp.ConnectRequest

type GetAvatarRequest struct {
	Phone    string `json:"phone" validate:"required_without=Username,max=255"`
	Username string `json:"username" validate:"required_without=Phone,max=255"`
	Preview  bool   `json:"preview"`
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
	Phone    []string `json:"phone" validate:"required_without=Username,dive,min=1,max=255"`
	Username []string `json:"username" validate:"required_without=Phone,dive,min=1,max=255"`
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
	Phone    string `json:"phone" validate:"required_without=Username,max=255"`
	Username string `json:"username" validate:"required_without=Phone,max=255"`
	Audio    string `json:"audio" validate:"required"`
	Caption  string `json:"caption" validate:"omitempty,max=4096"`
	ID       string `json:"id" validate:"omitempty,max=255"`

	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendAudioRequest

type Button struct {
	ButtonID   string `json:"button_id" validate:"required,max=255"`
	ButtonText string `json:"button_text" validate:"required,max=255"`
} // @name whatsapp.Button

type SendButtonTextRequest struct {
	Phone    string   `json:"phone" validate:"required_without=Username,max=255"`
	Username string   `json:"username" validate:"required_without=Phone,max=255"`
	Title    string   `json:"title" validate:"required,max=4096"`
	Buttons  []Button `json:"buttons" validate:"required,min=1,max=10,dive"`
	ID       string   `json:"id" validate:"omitempty,max=255"`
} // @name whatsapp.SendButtonTextRequest

type ChatPresenceRequest struct {
	Phone    string `json:"phone" validate:"required_without=Username,max=255"`
	Username string `json:"username" validate:"required_without=Phone,max=255"`
	State    string `json:"state" validate:"required,oneof=composing paused"`
	Media    string `json:"media" validate:"omitempty,oneof=audio"`
} // @name whatsapp.ChatPresenceRequest

type SendContactRequest struct {
	Phone       string      `json:"phone" validate:"required_without=Username,max=255"`
	Username    string      `json:"username" validate:"required_without=Phone,max=255"`
	ID          string      `json:"id" validate:"omitempty,max=255"`
	Name        string      `json:"name" validate:"required,max=255"`
	Vcard       string      `json:"vcard" validate:"required"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendContactRequest

type SendDocumentRequest struct {
	Phone       string      `json:"phone" validate:"required_without=Username,max=255"`
	Username    string      `json:"username" validate:"required_without=Phone,max=255"`
	Document    string      `json:"document" validate:"required"`
	FileName    string      `json:"filename" validate:"required,max=255"`
	ID          string      `json:"id" validate:"omitempty,max=255"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendDocumentRequest

type SendImageRequest struct {
	Phone       string      `json:"phone" validate:"required_without=Username,max=255"`
	Username    string      `json:"username" validate:"required_without=Phone,max=255"`
	Image       string      `json:"image" validate:"required"`
	Caption     string      `json:"caption" validate:"omitempty,max=4096"`
	ID          string      `json:"id" validate:"omitempty,max=255"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendImageRequest

type Row struct {
	RowID       string `json:"row_id" validate:"omitempty,max=255"`
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
} // @name whatsapp.Row

type Section struct {
	Title string `json:"title" validate:"required,max=255"`
	Rows  []Row  `json:"rows" validate:"required,min=1,max=10,dive"`
} // @name whatsapp.Section

type SendListRequest struct {
	Phone       string    `json:"phone" validate:"required_without=Username,max=255"`
	Username    string    `json:"username" validate:"required_without=Phone,max=255"`
	Title       string    `json:"title" validate:"required,max=255"`
	Description string    `json:"description" validate:"required,max=4096"`
	ButtonText  string    `json:"button_text" validate:"required,max=255"`
	FooterText  string    `json:"footer_text" validate:"omitempty,max=255"`
	Sections    []Section `json:"sections" validate:"required,min=1,max=10,dive"`
	ID          string    `json:"id" validate:"omitempty,max=255"`
} // @name whatsapp.SendListRequest

type SendLocationRequest struct {
	Phone       string      `json:"phone" validate:"required_without=Username,max=255"`
	Username    string      `json:"username" validate:"required_without=Phone,max=255"`
	ID          string      `json:"id" validate:"omitempty,max=255"`
	Name        string      `json:"name" validate:"omitempty,max=255"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendLocationRequest

type SendStickerRequest struct {
	Phone        string      `json:"phone" validate:"required_without=Username,max=255"`
	Username     string      `json:"username" validate:"required_without=Phone,max=255"`
	Sticker      string      `json:"sticker" validate:"required"`
	ID           string      `json:"id" validate:"omitempty,max=255"`
	PngThumbnail []byte      `json:"png_thumbnail"`
	ContextInfo  ContextInfo `json:"context_info"`
} // @name whatsapp.SendStickerRequest

type SendTextRequest struct {
	Phone       string      `json:"phone" validate:"required_without=Username,max=255"`
	Username    string      `json:"username" validate:"required_without=Phone,max=255"`
	Body        string      `json:"body" validate:"required,max=4096"`
	ID          string      `json:"id" validate:"omitempty,max=255"`
	ContextInfo ContextInfo `json:"context_info"`
} // @name whatsapp.SendTextRequest

type SendVideoRequest struct {
	Phone         string      `json:"phone" validate:"required_without=Username,max=255"`
	Username      string      `json:"username" validate:"required_without=Phone,max=255"`
	Video         string      `json:"video" validate:"required"`
	Caption       string      `json:"caption" validate:"omitempty,max=4096"`
	ID            string      `json:"id" validate:"omitempty,max=255"`
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
