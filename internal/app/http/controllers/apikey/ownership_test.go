package apikey

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"sapasora/internal/app/http/controllers"
	"sapasora/internal/modules/account"
	apikeymodule "sapasora/internal/modules/apikey"
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type ownershipAPIKeyServiceStub struct {
	apikeymodule.APIKeyService
	globalLookupCalled bool
	scopedLookupCalled bool
	gotOwnerID         uint
	created            *apikeymodule.APIKey
}

func (s *ownershipAPIKeyServiceStub) GetAPIKeyByPublicID(_ context.Context, _ string) (*apikeymodule.APIKey, error) {
	s.globalLookupCalled = true
	return nil, nil
}

func (s *ownershipAPIKeyServiceStub) GetAPIKeyByPublicIDForUser(_ context.Context, _ string, ownerID uint) (*apikeymodule.APIKey, error) {
	s.scopedLookupCalled = true
	s.gotOwnerID = ownerID
	return nil, gorm.ErrRecordNotFound
}

func (s *ownershipAPIKeyServiceStub) CreateAPIKey(_ context.Context, created *apikeymodule.APIKey) error {
	s.created = created
	return nil
}

func newAPIKeyOwnershipTestApp(service apikeymodule.APIKeyService) *fiber.App {
	controller := &APIKeyController{
		Controller:    &controllers.Controller{},
		apikeyService: service,
	}
	accountSubject := &account.Account{
		ID: 7,
		Role: &role.Role{Permissions: &permission.Permission{
			CanStoreAPIKey:  true,
			CanDeleteAPIKey: true,
		}},
	}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Delete("/api-key/:id", func(ctx *fiber.Ctx) error {
		ctx.Locals("account", accountSubject)
		return controller.Destroy(ctx)
	})
	app.Post("/api-key", func(ctx *fiber.Ctx) error {
		ctx.Locals("account", accountSubject)
		return controller.Store(ctx)
	})
	return app
}

func TestAPIKeyDestroyRejectsCrossOwnerWithoutEnumeration(t *testing.T) {
	service := &ownershipAPIKeyServiceStub{}
	app := newAPIKeyOwnershipTestApp(service)

	response, err := app.Test(httptest.NewRequest(http.MethodDelete, "/api-key/key-owned-by-8", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotContains(t, string(body), "key-owned-by-8")
	require.True(t, service.scopedLookupCalled)
	require.False(t, service.globalLookupCalled)
	require.Equal(t, uint(7), service.gotOwnerID)
}

func TestAPIKeyStoreUsesAuthenticatedOwnerInsteadOfRequestUserID(t *testing.T) {
	service := &ownershipAPIKeyServiceStub{}
	app := newAPIKeyOwnershipTestApp(service)
	request := httptest.NewRequest(http.MethodPost, "/api-key", bytes.NewBufferString(`{"name":"key","user_id":"owner-8","permissions":{}}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.NotNil(t, service.created)
	require.Equal(t, uint(7), service.created.UserID)
}
