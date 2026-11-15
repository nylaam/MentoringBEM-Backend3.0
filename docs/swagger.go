package docs

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
)

const swaggerUITemplate = `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <title>Toko Backend - Swagger API Documentation</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5/favicon-32x32.png" sizes="32x32" />
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; }
    .topbar { display: none; }
    .swagger-ui .info { margin: 25px 0; }
    .swagger-ui .info .title { font-size: 28px; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/docs/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        persistAuthorization: true
      });
      if ("%s" !== "") {
        ui.preauthorizeApiKey("ApiKeyAuth", "%s");
      }
      window.ui = ui;
    };
  </script>
</body>
</html>`

// RegisterSwaggerRoutes mendaftarkan endpoint tunggal Swagger UI dan spesifikasi OpenAPI dengan proteksi Basic Auth.
// Kredensial login (username & password) serta API Key dibaca langsung dari environment (.env).
func RegisterSwaggerRoutes(router fiber.Router, username, password, apiKey string) {
	docsGroup := router.Group("/docs", basicauth.New(basicauth.Config{
		Users: map[string]string{
			username: password,
		},
		Realm: "Toko Backend Swagger Documentation",
	}))

	renderedHTML := fmt.Sprintf(swaggerUITemplate, apiKey, apiKey)

	docsGroup.Get("", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(renderedHTML)
	})

	docsGroup.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(renderedHTML)
	})

	docsGroup.Get("/openapi.json", func(c fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Send(OpenAPISpec)
	})
}
