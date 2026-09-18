package account

import (
	"context"
	"strings"

	"sapasora/platform/database"
	"sapasora/platform/support/sql"
)

type AccountRepositoryImpl struct {
	db *database.Database
}

// NewAccountRepository creates a new Implementation of AccountRepository
// @wired:provide
func NewAccountRepository(db *database.Database) AccountRepository {
	return &AccountRepositoryImpl{
		db: db,
	}
}

func (r *AccountRepositoryImpl) ListAccount(
	ctx context.Context,
	search string,
	page, limit int,
) (*Pagination[*Account], error) {

	wb := sql.NewLogicBuilder()

	wb.And("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")

	qb := r.db.WithContext(ctx).
		Model(&Account{})

	where, args := wb.GetLogic()
	if where != "" {
		qb = qb.Where(where, args...)
	}

	var count int64
	total := qb.Count(&count)
	if total.Error != nil {
		return nil, total.Error
	}

	var accounts []*Account
	q := qb.Offset(page - 1).Limit(limit).Find(&accounts)
	if q.Error != nil {
		return nil, q.Error
	}

	return &Pagination[*Account]{
		Page:  page,
		Limit: limit,
		Total: count,
		Data:  accounts,
	}, nil
}

func (r *AccountRepositoryImpl) CreateAccount(ctx context.Context, account *Account) error {
	if err := account.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *AccountRepositoryImpl) GetAccount(ctx context.Context, id int) (*Account, error) {
	var result Account
	q := r.db.WithContext(ctx).First(&result, id)
	return &result, q.Error
}

func (r *AccountRepositoryImpl) GetAccountByPublicID(
	ctx context.Context,
	publicID string,
) (*Account, error) {
	var result Account
	q := r.db.WithContext(ctx).Where("public_id = ?", publicID).First(&result)
	return &result, q.Error
}

func (r *AccountRepositoryImpl) GetAccountByUsername(
	ctx context.Context,
	username string,
) (*Account, error) {
	var result Account
	q := r.db.WithContext(ctx).Where("username = ?", username).First(&result)
	return &result, q.Error
}

func (r *AccountRepositoryImpl) UpdateAccount(ctx context.Context, account *Account) error {
	if err := account.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *AccountRepositoryImpl) DeleteAccount(ctx context.Context, account *Account) error {
	if err := account.Valid(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(account).Error
}
