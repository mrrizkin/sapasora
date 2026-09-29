package apikey

import (
	"context"
)

type APIKeyService interface {
	ListAPIKey(ctx context.Context, page, limit int) (*Pagination[*APIKey], error)
	// ListAPIKeyForUser scopes the listing to API keys owned by userID.
	// Callers listing keys on behalf of a specific account subject must use
	// this instead of ListAPIKey, which returns every key in the system
	// regardless of owner.
	ListAPIKeyForUser(ctx context.Context, userID uint, page, limit int) (*Pagination[*APIKey], error)
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
	ListAPIKeyForUser(ctx context.Context, userID uint, page, limit int) (*Pagination[*APIKey], error)
	CreateAPIKey(ctx context.Context, apikey *APIKey) error
	GetAPIKey(ctx context.Context, id int) (*APIKey, error)
	GetAPIKeyByPublicID(ctx context.Context, publicID string) (*APIKey, error)
	GetAPIKeyByPublicIDForUser(ctx context.Context, publicID string, userID uint) (*APIKey, error)
	GetAPIKeyByKey(ctx context.Context, key string) (*APIKey, error)
	UpdateAPIKey(ctx context.Context, apikey *APIKey) error
	DeleteAPIKey(ctx context.Context, apikey *APIKey) error
}
