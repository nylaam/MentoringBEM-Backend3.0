package unit_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/docs"
	"github.com/stretchr/testify/assert"
)

func setupSwaggerTestApp(user, pass, apiKey string) *fiber.App {
	app := fiber.New()
	docs.RegisterSwaggerRoutes(app, user, pass, apiKey)
	return app
}

func basicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestSwagger_Unauthorized_WithoutCredentials(t *testing.T) {
	app := setupSwaggerTestApp("admin", "admin123", "secret-key")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	assert.Contains(t, strings.ToLower(resp.Header.Get("Www-Authenticate")), "basic")
}

func TestSwagger_Unauthorized_WrongPassword(t *testing.T) {
	app := setupSwaggerTestApp("admin", "admin123", "secret-key")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	req.Header.Set("Authorization", basicAuthHeader("admin", "wrongpassword"))

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestSwagger_Success_WithCredentials(t *testing.T) {
	app := setupSwaggerTestApp("admin", "admin123", "secret-key")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	req.Header.Set("Authorization", basicAuthHeader("admin", "admin123"))

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/html")

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(bodyBytes), "swagger-ui")
	assert.Contains(t, string(bodyBytes), "secret-key")
}

func TestSwagger_OpenAPISpec_Success(t *testing.T) {
	app := setupSwaggerTestApp("admin", "admin123", "secret-key")

	req := httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil)
	req.Header.Set("Authorization", basicAuthHeader("admin", "admin123"))

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var specMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &specMap)
	assert.NoError(t, err)
	assert.Equal(t, "3.0.3", specMap["openapi"])
}
