package validator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codeberg.org/mrrizkin/nihil"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type validationPayload struct {
	Name    string          `validate:"required,min=3,max=10"`
	Mode    string          `validate:"oneof=fast safe"`
	Email   string          `validate:"omitempty,email"`
	Webhook nihil.NilString `validate:"omitempty,url"`
}

func TestValidatorReportsTagsAndAcceptsNullableURL(t *testing.T) {
	v := NewValidator()

	require.Empty(t, v.Validate(validationPayload{
		Name:    "valid",
		Mode:    "fast",
		Email:   "person@example.test",
		Webhook: nihil.String("https://example.test/hook"),
	}))

	errs := v.Validate(validationPayload{Name: "x", Mode: "unknown", Email: "not-an-email", Webhook: nihil.String("not-a-url")})
	require.Len(t, errs, 4)
	require.ElementsMatch(t, []string{"min", "oneof", "email", "url"}, []string{errs[0].Tag, errs[1].Tag, errs[2].Tag, errs[3].Tag})
}

func TestParseBodyAndValidateMapsMalformedAndInvalidPayloadsTo422(t *testing.T) {
	v := NewValidator()
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/payload", func(ctx *fiber.Ctx) error {
		var payload validationPayload
		if err := v.ParseBodyAndValidate(ctx, &payload); err != nil {
			return err
		}
		return ctx.SendStatus(http.StatusNoContent)
	})

	for _, body := range []string{`{"name":`, `{"name":"x","mode":"fast"}`} {
		request := httptest.NewRequest(http.MethodPost, "/payload", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response, err := app.Test(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	}

	request := httptest.NewRequest(http.MethodPost, "/payload", strings.NewReader(`{"name":"valid","mode":"safe","webhook":"https://example.test/hook"}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
}
