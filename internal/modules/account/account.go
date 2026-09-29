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

// SessionAccount is the shape persisted into the session store. It mirrors
// Account but keeps the primary key (ID), which Account hides from JSON via
// `json:"-"` to prevent raw DB IDs leaking into API responses. Session
// storage is not an API response, and the app relies on the ID surviving the
// round-trip for every owner-scoped mutation (see Controller.GetOwnerID), so
// it is serialized through this dedicated type instead.
type SessionAccount struct {
	ID       uint       `json:"id"`
	PublicID string     `json:"public_id"`
	Name     string     `json:"name"`
	Username string     `json:"username"`
	RoleID   uint       `json:"role_id"`
	Role     *role.Role `json:"role,omitempty"`
}

// ToSession converts an Account into its session-storable representation.
func ToSession(a *Account) SessionAccount {
	return SessionAccount{
		ID:       a.ID,
		PublicID: a.PublicID,
		Name:     a.Name,
		Username: a.Username,
		RoleID:   a.RoleID,
		Role:     a.Role,
	}
}

// FromSession reconstructs an Account from its session-storable representation.
func FromSession(s SessionAccount) *Account {
	return &Account{
		ID:       s.ID,
		PublicID: s.PublicID,
		Name:     s.Name,
		Username: s.Username,
		RoleID:   s.RoleID,
		Role:     s.Role,
	}
}
