// Package permission provides the Permission modules
package permission

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Permission struct {
	// devices
	CanListDevice   bool `json:"can_list_device"`
	CanGetDevice    bool `json:"can_get_device"`
	CanUpdateDevice bool `json:"can_update_device"`
	CanStoreDevice  bool `json:"can_store_device"`
	CanDeleteDevice bool `json:"can_delete_device"`

	// accounts
	CanListAccount           bool `json:"can_list_account"`
	CanGetAccount            bool `json:"can_get_account"`
	CanUpdateAccount         bool `json:"can_update_account"`
	CanUpdatePasswordAccount bool `json:"can_update_password_account"`
	CanStoreAccount          bool `json:"can_store_account"`
	CanDeleteAccount         bool `json:"can_delete_account"`

	// roles
	CanListRole             bool `json:"can_list_role"`
	CanGetRole              bool `json:"can_get_role"`
	CanUpdateRole           bool `json:"can_update_role"`
	CanUpdatePermissionRole bool `json:"can_update_permission_role"`
	CanStoreRole            bool `json:"can_store_role"`
	CanDeleteRole           bool `json:"can_delete_role"`

	// apikeys
	CanListAPIKey   bool `json:"can_list_api_key"`
	CanGetAPIKey    bool `json:"can_get_api_key"`
	CanUpdateAPIKey bool `json:"can_update_api_key"`
	CanStoreAPIKey  bool `json:"can_store_api_key"`
	CanDeleteAPIKey bool `json:"can_delete_api_key"`

	// gateways
	CanCheckUserGateway        bool `json:"can_check_user_gateway"`
	CanConnectGateway          bool `json:"can_connect_gateway"`
	CanDisconnectGateway       bool `json:"can_disconnect_gateway"`
	CanGetAvatarGateway        bool `json:"can_get_avatar_gateway"`
	CanGetContactsGateway      bool `json:"can_get_contacts_gateway"`
	CanGetQRGateway            bool `json:"can_get_qr_gateway"`
	CanGetStatusGateway        bool `json:"can_get_status_gateway"`
	CanGetUserGateway          bool `json:"can_get_user_gateway"`
	CanLogoutGateway           bool `json:"can_logout_gateway"`
	CanSendAudioGateway        bool `json:"can_send_audio_gateway"`
	CanSendButtonGateway       bool `json:"can_send_button_gateway"`
	CanSendChatPresenceGateway bool `json:"can_send_chat_presence_gateway"`
	CanSendContactGateway      bool `json:"can_send_contact_gateway"`
	CanSendDocumentGateway     bool `json:"can_send_document_gateway"`
	CanSendImageGateway        bool `json:"can_send_image_gateway"`
	CanSendListGateway         bool `json:"can_send_list_gateway"`
	CanSendLocationGateway     bool `json:"can_send_location_gateway"`
	CanSendStickerGateway      bool `json:"can_send_sticker_gateway"`
	CanSendTextGateway         bool `json:"can_send_text_gateway"`
	CanSendVideoGateway        bool `json:"can_send_video_gateway"`

	// devicetokens
	CanListDeviceToken   bool `json:"can_list_device_token"`
	CanGetDeviceToken    bool `json:"can_get_device_token"`
	CanStoreDeviceToken  bool `json:"can_store_device_token"`
	CanDeleteDeviceToken bool `json:"can_delete_device_token"`
} // @name permission.Permission

func (p *Permission) Scan(value any) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to scan Permission: invalid type")
	}

	return json.Unmarshal(bytes, p)
}

func (p Permission) Value() (driver.Value, error) {
	if p == (Permission{}) {
		return nil, nil
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return string(bytes), nil
}
