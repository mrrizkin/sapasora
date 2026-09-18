package apikey

import (
	"sapasora/internal/modules/apikey"
	"sapasora/internal/modules/permission"
)

type APIKeyResponse *apikey.APIKey // @name APIKeyController.APIKeyResponse

type APIKeyListResponse *apikey.Pagination[*apikey.APIKey] // @name APIKeyController.APIKeyListResponse

type APIKeyStoreRequest struct {
	Name        string                `json:"name"`
	Permissions permission.Permission `json:"permissions"`
	UserID      string                `json:"user_id"`
} // @name APIKeyController.APIKeyStoreRequest

type APIKeyUpdateRequest struct {
	Name        string                `json:"name"`
	Permissions permission.Permission `json:"permissions"`
	UserID      string                `json:"user_id"`
} // @name APIKeyController.APIKeyUpdateRequest
