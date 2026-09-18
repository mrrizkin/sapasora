package apikey

import (
	"context"

	"sapasora/platform/database"
	"sapasora/platform/support/sql"
)

type APIKeyRepositoryImpl struct {
	db *database.Database
}

// NewAPIKeyRepository creates a new Implementation of APIKeyRepository
// @wired:provide
func NewAPIKeyRepository(db *database.Database) APIKeyRepository {
	return &APIKeyRepositoryImpl{
		db: db,
	}
}

func (r *APIKeyRepositoryImpl) ListAPIKey(
	ctx context.Context,
	page, limit int,
) (*Pagination[*APIKey], error) {

	wb := sql.NewLogicBuilder()

	qb := r.db.WithContext(ctx).
		Model(&APIKey{})

	where, args := wb.GetLogic()
	if where != "" {
		qb = qb.Where(where, args...)
	}

	var count int64
	total := qb.Count(&count)
	if total.Error != nil {
		return nil, total.Error
	}

	var result []*APIKey
	q := qb.Offset((page - 1) * limit).Limit(limit).Find(&result)
	if q.Error != nil {
		return nil, q.Error
	}

	return &Pagination[*APIKey]{
		Page:  page,
		Limit: limit,
		Total: count,
		Data:  result,
	}, q.Error
}

func (r *APIKeyRepositoryImpl) CreateAPIKey(ctx context.Context, apikey *APIKey) error {
	if err := apikey.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(apikey).Error
}

func (r *APIKeyRepositoryImpl) GetAPIKey(ctx context.Context, id int) (*APIKey, error) {
	var result APIKey
	q := r.db.WithContext(ctx).First(&result, id)
	return &result, q.Error
}

func (r *APIKeyRepositoryImpl) GetAPIKeyByPublicID(
	ctx context.Context,
	publicID string,
) (*APIKey, error) {
	var result APIKey
	q := r.db.WithContext(ctx).Where("public_id = ?", publicID).First(&result)
	return &result, q.Error
}

func (r *APIKeyRepositoryImpl) GetAPIKeyByKey(
	ctx context.Context,
	key string,
) (*APIKey, error) {
	var result APIKey
	q := r.db.WithContext(ctx).Where("key = ?", key).First(&result)
	return &result, q.Error
}

func (r *APIKeyRepositoryImpl) UpdateAPIKey(ctx context.Context, apikey *APIKey) error {
	if err := apikey.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(apikey).Error
}

func (r *APIKeyRepositoryImpl) DeleteAPIKey(ctx context.Context, apikey *APIKey) error {
	if err := apikey.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(apikey).Error
}
