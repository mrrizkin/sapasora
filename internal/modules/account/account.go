// Package account provides the Account modules
package account

import (
	"sapasora/internal/modules/role"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        uint           `json:"-"          gorm:"primarykey"`
	PublicID  string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Name           string `json:"name"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"`

	RoleID uint       `json:"-"`
	Role   *role.Role `json:"role,omitempty" gorm:"->all"`
} // @name account.Account

type Pagination[T any] struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Data  []T   `json:"data"`
}

func (m *Account) TableName() string {
	return "m_users"
}

func (m *Account) Valid() error {
	return nil
}
