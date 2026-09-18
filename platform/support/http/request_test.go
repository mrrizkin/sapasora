package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestParseListQueryParamsUsesDefaultsAndBoundsPagination(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", func(c *fiber.Ctx) error {
		params, err := ParseListQueryParams(c)
		if err != nil {
			return err
		}
		return c.JSON(params)
	})

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	var params ListQueryParams
	require.NoError(t, json.NewDecoder(response.Body).Decode(&params))
	require.Equal(t, DefaultPage, params.Page)
	require.Equal(t, DefaultLimit, params.Limit)

	response, err = app.Test(httptest.NewRequest(http.MethodGet, "/?page=2&limit=25&search=device", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.NoError(t, json.NewDecoder(response.Body).Decode(&params))
	require.Equal(t, 2, params.Page)
	require.Equal(t, 25, params.Limit)
	require.Equal(t, "device", params.Search)

	for _, query := range []string{
		"?page=0",
		"?page=not-a-number",
		"?page=",
		"?limit=0",
		"?limit=101",
		"?limit=not-a-number",
		"?search=" + strings.Repeat("x", MaxFilterLength+1),
	} {
		response, err = app.Test(httptest.NewRequest(http.MethodGet, "/"+query, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, response.StatusCode, query)
	}
}

func TestParseQueryParamsRejectsExplicitEmptyNumericValue(t *testing.T) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, _ error) error {
			return c.SendStatus(http.StatusBadRequest)
		},
	})
	app.Get("/", func(c *fiber.Ctx) error {
		var params PaginationParams
		return ParseQueryParams(c, &params)
	})

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/?limit=", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestValidatePublicID(t *testing.T) {
	for _, publicID := range []string{"key-owned-by-8", "V1StGXR8_Z5jdHi6B-myT"} {
		require.NoError(t, ValidatePublicID(publicID))
	}
	for _, publicID := range []string{"", "../secret", "id with spaces", strings.Repeat("a", 256)} {
		require.Error(t, ValidatePublicID(publicID), publicID)
	}
}
