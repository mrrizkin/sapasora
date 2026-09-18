package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type controllerValidationPayload struct {
	Name string `validate:"required,min=3"`
}

func TestValidationHelpersUseDefaultValidatorWhenControllerIsUnwired(t *testing.T) {
	controller := &Controller{}

	require.Len(t, controller.Validate(controllerValidationPayload{}), 1)
	require.Error(t, controller.MustValidate(controllerValidationPayload{}))
	require.Empty(t, controller.Validate(controllerValidationPayload{Name: "valid"}))
}

func TestPublicIDParamRejectsUnsafePathValues(t *testing.T) {
	controller := &Controller{}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/:id", func(ctx *fiber.Ctx) error {
		if _, err := controller.PublicIDParam(ctx); err != nil {
			return err
		}
		return ctx.SendStatus(http.StatusNoContent)
	})

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/valid-id_123", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)

	response, err = app.Test(httptest.NewRequest(http.MethodGet, "/not.valid", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestBodyParserValidateMapsMalformedPayloadTo422(t *testing.T) {
	controller := &Controller{}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/payload", func(ctx *fiber.Ctx) error {
		var payload controllerValidationPayload
		return controller.BodyParserValidate(ctx, &payload)
	})

	request := httptest.NewRequest(http.MethodPost, "/payload", strings.NewReader(`{"name":`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)

	require.NoError(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
}
