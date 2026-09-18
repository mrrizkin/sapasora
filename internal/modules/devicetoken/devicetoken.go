// Package devicetoken provides the Devicetoken modules
package devicetoken

import (
	"sapasora/internal/modules/account"
	"time"

	"codeberg.org/mrrizkin/nihil"
	"gorm.io/gorm"
)

type DeviceToken struct {
	ID        uint              `json:"-"          gorm:"primarykey"`
	PublicID  string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"-"`
	Token     string            `json:"token"`
	Status    DeviceTokenStatus `json:"status"`
	ExpiredAt nihil.NilTime     `json:"-"`

	DeviceID uint `json:"-"`

	UserID uint             `json:"-"`
	User   *account.Account `json:"user,omitempty" gorm:"->:all"`
} // @name devicetoken.DeviceToken

type Pagination[T any] struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Data  []T   `json:"data"`
}

func (m *DeviceToken) TableName() string {
	return "m_device_tokens"
}

func (m *DeviceToken) Valid() error {
	return nil
}
