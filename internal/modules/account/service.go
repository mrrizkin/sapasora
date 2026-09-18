package account

import (
	"context"
)

type AccountServiceImpl struct {
	repo AccountRepository
}

// NewAccountService creates a new Implementation of AccountService
// @wired:provide
func NewAccountService(repo AccountRepository) AccountService {
	return &AccountServiceImpl{
		repo: repo,
	}
}

func (s *AccountServiceImpl) ListAccount(
	ctx context.Context,
	search string,
	page, limit int,
) (*Pagination[*Account], error) {
	return s.repo.ListAccount(ctx, search, page, limit)
}

func (s *AccountServiceImpl) CreateAccount(ctx context.Context, account *Account) error {
	return s.repo.CreateAccount(ctx, account)
}

func (s *AccountServiceImpl) GetAccount(ctx context.Context, id int) (*Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *AccountServiceImpl) GetAccountByPublicID(
	ctx context.Context,
	publicID string,
) (*Account, error) {
	return s.repo.GetAccountByPublicID(ctx, publicID)
}

func (s *AccountServiceImpl) GetAccountByUsername(
	ctx context.Context,
	username string,
) (*Account, error) {
	return s.repo.GetAccountByUsername(ctx, username)
}

func (s *AccountServiceImpl) UpdateAccount(ctx context.Context, account *Account) error {
	return s.repo.UpdateAccount(ctx, account)
}

func (s *AccountServiceImpl) DeleteAccount(ctx context.Context, account *Account) error {
	return s.repo.DeleteAccount(ctx, account)
}
