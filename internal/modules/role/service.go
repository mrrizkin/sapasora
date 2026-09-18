package role

import (
	"context"
)

type RoleServiceImpl struct {
	repo RoleRepository
}

// NewRoleService creates a new Implementation of RoleService
// @wired:provide
func NewRoleService(repo RoleRepository) RoleService {
	return &RoleServiceImpl{
		repo: repo,
	}
}

func (s *RoleServiceImpl) ListRole(
	ctx context.Context,
	search string,
	page, limit int,
) (*Pagination[*Role], error) {
	return s.repo.ListRole(ctx, search, page, limit)
}

func (s *RoleServiceImpl) CreateRole(ctx context.Context, role *Role) error {
	return s.repo.CreateRole(ctx, role)
}

func (s *RoleServiceImpl) GetRole(ctx context.Context, id int) (*Role, error) {
	return s.repo.GetRole(ctx, id)
}

func (s *RoleServiceImpl) GetRoleByPublicID(
	ctx context.Context,
	publicID string,
) (*Role, error) {
	return s.repo.GetRoleByPublicID(ctx, publicID)
}

func (s *RoleServiceImpl) GetRoleByName(
	ctx context.Context,
	name string,
) (*Role, error) {
	return s.repo.GetRoleByName(ctx, name)
}

func (s *RoleServiceImpl) UpdateRole(ctx context.Context, role *Role) error {
	return s.repo.UpdateRole(ctx, role)
}

func (s *RoleServiceImpl) DeleteRole(ctx context.Context, role *Role) error {
	return s.repo.DeleteRole(ctx, role)
}
