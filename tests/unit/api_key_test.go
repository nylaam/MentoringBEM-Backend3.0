package unit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/internal/middlewares"
	"github.com/mentoring/toko-backend/pkg/response"
	"github.com/stretchr/testify/assert"
)

func setupAPIKeyTestApp(apiKey string) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.CustomErrorHandler,
	})

	api := app.Group("/api/v1")
	if apiKey != "" {
		api.Use(middlewares.NewAPIKeyMiddleware(apiKey))
	}

	api.Get("/test", func(c fiber.Ctx) error {
		return response.JSONSuccess(c, fiber.StatusOK, fiber.Map{"authenticated": true})
	})

	return app
}

func TestAPIKeyMiddleware_Success_XAPIKey(t *testing.T) {
	app := setupAPIKeyTestApp("secret-123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("X-API-Key", "secret-123")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAPIKeyMiddleware_Success_BearerAuth(t *testing.T) {
	app := setupAPIKeyTestApp("secret-123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("Authorization", "Bearer secret-123")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAPIKeyMiddleware_Success_QueryParam(t *testing.T) {
	app := setupAPIKeyTestApp("secret-123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test?api_key=secret-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAPIKeyMiddleware_Unauthorized_MissingKey(t *testing.T) {
	app := setupAPIKeyTestApp("secret-123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAPIKeyMiddleware_Unauthorized_WrongKey(t *testing.T) {
	app := setupAPIKeyTestApp("secret-123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("X-API-Key", "wrong-key")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAPIKeyMiddleware_DisabledWhenKeyEmpty(t *testing.T) {
	app := setupAPIKeyTestApp("")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
