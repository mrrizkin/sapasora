package account

import (
	"context"
)

type AccountService interface {
	ListAccount(ctx context.Context, search string, page, limit int) (*Pagination[*Account], error)
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, id int) (*Account, error)
	GetAccountByPublicID(ctx context.Context, publicID string) (*Account, error)
	GetAccountByUsername(ctx context.Context, username string) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, account *Account) error
}

type AccountRepository interface {
	ListAccount(ctx context.Context, search string, page, limit int) (*Pagination[*Account], error)
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, id int) (*Account, error)
	GetAccountByPublicID(ctx context.Context, publicID string) (*Account, error)
	GetAccountByUsername(ctx context.Context, username string) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, account *Account) error
}
