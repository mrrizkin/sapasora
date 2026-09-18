package apikey

import (
	"context"
)

type APIKeyService interface {
	ListAPIKey(ctx context.Context, page, limit int) (*Pagination[*APIKey], error)
	CreateAPIKey(ctx context.Context, apikey *APIKey) error
	GetAPIKey(ctx context.Context, id int) (*APIKey, error)
	GetAPIKeyByPublicID(ctx context.Context, publicID string) (*APIKey, error)
	GetAPIKeyByPublicIDForUser(ctx context.Context, publicID string, userID uint) (*APIKey, error)
	GetAPIKeyByKey(ctx context.Context, key string) (*APIKey, error)
	UpdateAPIKey(ctx context.Context, apikey *APIKey) error
	DeleteAPIKey(ctx context.Context, apikey *APIKey) error
}

type APIKeyRepository interface {
	ListAPIKey(ctx context.Context, page, limit int) (*Pagination[*APIKey], error)
	CreateAPIKey(ctx context.Context, apikey *APIKey) error
	GetAPIKey(ctx context.Context, id int) (*APIKey, error)
	GetAPIKeyByPublicID(ctx context.Context, publicID string) (*APIKey, error)
	GetAPIKeyByPublicIDForUser(ctx context.Context, publicID string, userID uint) (*APIKey, error)
	GetAPIKeyByKey(ctx context.Context, key string) (*APIKey, error)
	UpdateAPIKey(ctx context.Context, apikey *APIKey) error
	DeleteAPIKey(ctx context.Context, apikey *APIKey) error
}
