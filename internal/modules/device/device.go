// Package device provides the Device modules
package device

import (
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/devicetoken"
	"sapasora/internal/modules/permission"
	"time"

	"codeberg.org/mrrizkin/nihil"
	"gorm.io/gorm"
)

type Device struct {
	ID        uint           `json:"-"          gorm:"primarykey"`
	PublicID  string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Name        string                 `json:"name"`
	Type        DeviceType             `json:"type"`
	Webhook     nihil.NilString        `json:"webhook"`
	Jid         nihil.NilString        `json:"jid"`
	QRCode      nihil.NilString        `json:"qr_code"`
	Status      DeviceStatus           `json:"status"`
	AutoConnect bool                   `json:"auto_connect"`
	ExpiredAt   nihil.NilTime          `json:"-"`
	Events      nihil.NilString        `json:"events"`
	Permissions *permission.Permission `json:"permission"`

	UserID uint             `json:"-"`
	User   *account.Account `json:"user,omitempty" gorm:"->:all"`

	Tokens []*devicetoken.DeviceToken `json:"tokens,omitempty" gorm:"->:all"`
} // @name device.Device

type Pagination[T any] struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Data  []T   `json:"data"`
}

func (m *Device) TableName() string {
	return "m_devices"
}

func (m *Device) Valid() error {
	return nil
}
