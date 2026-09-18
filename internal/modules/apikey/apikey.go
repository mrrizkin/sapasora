// Package apikey provides the Apikey modules
package apikey

import (
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/permission"
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	ID          uint                   `json:"-"          gorm:"primarykey"`
	PublicID    string                 `json:"id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"-"`
	Name        string                 `json:"name"`
	Key         string                 `json:"key"`
	Status      APIKeyStatus           `json:"status"`
	Permissions *permission.Permission `json:"permission"`

	UserID uint             `json:"-"`
	User   *account.Account `json:"user,omitempty" gorm:"->:all"`
} // @name apikey.APIKey

type Pagination[T any] struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Data  []T   `json:"data"`
}

func (m *APIKey) TableName() string {
	return "m_api_keys"
}

func (m *APIKey) Valid() error {
	return nil
}
