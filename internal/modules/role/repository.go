package role

import (
	"context"
	"strings"

	"sapasora/platform/database"
	"sapasora/platform/support/sql"
)

type RoleRepositoryImpl struct {
	db *database.Database
}

// NewRoleRepository creates a new Implementation of RoleRepository
// @wired:provide
func NewRoleRepository(db *database.Database) RoleRepository {
	return &RoleRepositoryImpl{
		db: db,
	}
}

func (r *RoleRepositoryImpl) ListRole(
	ctx context.Context,
	search string,
	page, limit int,
) (*Pagination[*Role], error) {

	wb := sql.NewLogicBuilder()

	wb.And("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")

	qb := r.db.WithContext(ctx).
		Model(&Role{})

	where, args := wb.GetLogic()
	if where != "" {
		qb = qb.Where(where, args...)
	}

	var count int64
	total := qb.Count(&count)
	if total.Error != nil {
		return nil, total.Error
	}

	var roles []*Role
	q := qb.Offset(page - 1).Limit(limit).Find(&roles)
	if q.Error != nil {
		return nil, q.Error
	}

	return &Pagination[*Role]{
		Page:  page,
		Limit: limit,
		Total: count,
		Data:  roles,
	}, nil
}

func (r *RoleRepositoryImpl) CreateRole(ctx context.Context, role *Role) error {
	if err := role.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepositoryImpl) GetRole(ctx context.Context, id int) (*Role, error) {
	var result Role
	q := r.db.WithContext(ctx).First(&result, id)
	return &result, q.Error
}

func (r *RoleRepositoryImpl) GetRoleByPublicID(
	ctx context.Context,
	publicID string,
) (*Role, error) {
	var result Role
	q := r.db.WithContext(ctx).Where("public_id = ?", publicID).First(&result)
	return &result, q.Error
}

func (r *RoleRepositoryImpl) GetRoleByName(
	ctx context.Context,
	name string,
) (*Role, error) {
	var result Role
	q := r.db.WithContext(ctx).Where("name = ?", name).First(&result)
	return &result, q.Error
}

func (r *RoleRepositoryImpl) UpdateRole(ctx context.Context, role *Role) error {
	if err := role.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *RoleRepositoryImpl) DeleteRole(ctx context.Context, role *Role) error {
	if err := role.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(role).Error
}
