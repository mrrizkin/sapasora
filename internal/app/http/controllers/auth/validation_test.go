package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sapasora/internal/app/http/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLoginRejectsInvalidPayloadWith422(t *testing.T) {
	controller := &AuthController{Controller: &controllers.Controller{}}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/login", controller.Login)

	for _, body := range []string{`{"username":"ab","password":""}`, `{"username":`} {
		request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response, err := app.Test(request)

		require.NoError(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	}
}
