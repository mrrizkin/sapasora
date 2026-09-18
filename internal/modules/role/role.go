// Package role provides the Role modules
package role

import (
	"sapasora/internal/modules/permission"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint                   `json:"-"           gorm:"primarykey"`
	PublicID    string                 `json:"id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"-"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Permissions *permission.Permission `json:"permission"`
} // @name role.Role

type Pagination[T any] struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Data  []T   `json:"data"`
}

func (m *Role) TableName() string {
	return "m_roles"
}

func (m *Role) Valid() error {
	return nil
}
