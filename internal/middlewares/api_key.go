package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/pkg/response"
)

// NewAPIKeyMiddleware memverifikasi API Key pada header X-API-Key, api-key, Authorization Bearer, atau query parameter.
// Jika validKey bernilai kosong (""), middleware ini dilewati (tidak ada verifikasi).
func NewAPIKeyMiddleware(validKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if validKey == "" {
			return c.Next()
		}

		// 1. Cek header X-API-Key (standar industri)
		clientKey := c.Get("X-API-Key")

		// 2. Cek alternatif header api-key (lowercase)
		if clientKey == "" {
			clientKey = c.Get("api-key")
		}

		// 3. Cek header Authorization: Bearer <key>
		if clientKey == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				clientKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 4. Cek query parameter ?api_key= (opsional, memudahkan testing cepat di browser)
		if clientKey == "" {
			clientKey = c.Query("api_key")
		}

		if clientKey == "" || clientKey != validKey {
			return response.JSONError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "API Key tidak valid atau belum disertakan di header 'X-API-Key'")
		}

		return c.Next()
	}
}
