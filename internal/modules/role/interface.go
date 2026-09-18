package role

import (
	"context"
)

type RoleService interface {
	ListRole(ctx context.Context, search string, page, limit int) (*Pagination[*Role], error)
	CreateRole(ctx context.Context, role *Role) error
	GetRole(ctx context.Context, id int) (*Role, error)
	GetRoleByPublicID(ctx context.Context, publicID string) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, role *Role) error
}

type RoleRepository interface {
	ListRole(ctx context.Context, search string, page, limit int) (*Pagination[*Role], error)
	CreateRole(ctx context.Context, role *Role) error
	GetRole(ctx context.Context, id int) (*Role, error)
	GetRoleByPublicID(ctx context.Context, publicID string) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, role *Role) error
}
