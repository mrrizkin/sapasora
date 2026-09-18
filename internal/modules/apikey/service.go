package apikey

import (
	"context"
)

type APIKeyServiceImpl struct {
    repo APIKeyRepository
}

// NewAPIKeyService creates a new Implementation of APIKeyService
// @wired:provide
func NewAPIKeyService(repo APIKeyRepository) APIKeyService {
    return &APIKeyServiceImpl{
        repo: repo,
    }
}

func (s *APIKeyServiceImpl) ListAPIKey(ctx context.Context, page, limit int) (*Pagination[*APIKey], error) {
    return s.repo.ListAPIKey(ctx, page, limit)
}

func (s *APIKeyServiceImpl) CreateAPIKey(ctx context.Context, apikey *APIKey) error {
    return s.repo.CreateAPIKey(ctx, apikey)
}

func (s *APIKeyServiceImpl) GetAPIKey(ctx context.Context, id int) (*APIKey, error) {
    return s.repo.GetAPIKey(ctx, id)
}

func (s *APIKeyServiceImpl) GetAPIKeyByPublicID(ctx context.Context, publicID string) (*APIKey, error) {
    return s.repo.GetAPIKeyByPublicID(ctx, publicID)
}

func (s *APIKeyServiceImpl) GetAPIKeyByKey(ctx context.Context, key string) (*APIKey, error) {
    return s.repo.GetAPIKeyByKey(ctx, key)
}

func (s *APIKeyServiceImpl) UpdateAPIKey(ctx context.Context, apikey *APIKey) error {
    return s.repo.UpdateAPIKey(ctx, apikey)
}

func (s *APIKeyServiceImpl) DeleteAPIKey(ctx context.Context, apikey *APIKey) error {
    return s.repo.DeleteAPIKey(ctx, apikey)
}
