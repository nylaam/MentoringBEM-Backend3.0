package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/mentoring/toko-backend/docs"
	"github.com/mentoring/toko-backend/internal/app/product"
	"github.com/mentoring/toko-backend/internal/config"
	"github.com/mentoring/toko-backend/internal/infra/database"
	"github.com/mentoring/toko-backend/internal/infra/repository"
	"github.com/mentoring/toko-backend/internal/middlewares"
	"github.com/mentoring/toko-backend/pkg/response"
)

func main() {
	// 1. Load Configuration (12-Factor App)
	cfg := config.LoadConfig()
	log.Printf("[App] Memulai aplikasi di environment: %s (Port: %s)", cfg.AppEnv, cfg.AppPort)

	// 2. Inisialisasi Database Connection Pool (Infra)
	var db *sql.DB
	var err error
	db, err = database.NewPostgres(cfg)
	if err != nil {
		log.Printf("[Peringatan] Database Postgres tidak dapat dihubungi (%v). Berjalan dalam mode In-Memory Repository untuk hands-on lab / workshop!\n", err)
		db = nil
	} else {
		log.Println("[Database] Terhubung sukses ke PostgreSQL!")
		defer db.Close()
	}

	// 3. Inisialisasi Repositories (Mengimplementasikan Contracts / Interfaces)
	productRepo := repository.NewProductRepository(db)

	// 4. Inisialisasi Services / Use Cases (Dependency Inversion)
	productSvc := product.NewProductService(productRepo)

	// 5. Inisialisasi Delivery Handlers (Go Fiber v3)
	productHandler := product.NewProductHandler(productSvc)

	// 6. Setup Fiber v3 Engine & Centralized Error Handler
	app := fiber.New(fiber.Config{
		AppName:      "Toko Backend API - Product Service (Fiber v3)",
		ErrorHandler: middlewares.CustomErrorHandler,
	})

	// Middlewares Stack
	app.Use(recover.New())

	// CORS Middleware agar Frontend (Next.js, React, Vue, Vite, dll.) dapat mengakses API tanpa terblokir browser
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-API-Key", "api-key"},
	}))

	app.Use(middlewares.NewRequestIDMiddleware())
	app.Use(middlewares.NewLoggerMiddleware())

	// Root Welcome Endpoint
	app.Get("/", func(c fiber.Ctx) error {
		return response.JSONSuccess(c, fiber.StatusOK, fiber.Map{
			"message": "Selamat datang di Toko Backend Product API (Mentoring & Workshop)",
			"health":  "/healthz",
			"docs":    "/docs",
			"api":     "/api/v1/products",
			"version": "1.0.0",
		})
	})

	// Healthcheck endpoint
	app.Get("/healthz", func(c fiber.Ctx) error {
		dbStatus := "connected"
		if db == nil {
			dbStatus = "in-memory-fallback"
		}
		return response.JSONSuccess(c, fiber.StatusOK, fiber.Map{
			"status":   "healthy",
			"database": dbStatus,
			"time":     time.Now().UTC().Format(time.RFC3339),
			"version":  "1.0.0",
		})
	})

	// 7. Route Registration
	// Register Swagger UI Documentation tunggal (/docs) dengan proteksi Basic Auth & API Key dari .env
	docs.RegisterSwaggerRoutes(app, cfg.SwaggerUser, cfg.SwaggerPassword, cfg.APIKey)
	log.Printf("[Docs] Swagger UI aktif di http://localhost%s/docs (Kredensial dari .env: User '%s', Password Protected)\n", fmt.Sprintf(":%s", cfg.AppPort), cfg.SwaggerUser)

	// Business API (/api/v1)
	api := app.Group("/api/v1")

	// Pasang proteksi API Key pada /api/v1 jika API_KEY dikonfigurasi di .env
	if cfg.APIKey != "" {
		log.Println("[Security] Proteksi API Key aktif untuk endpoint /api/v1 (kunci dibaca dari .env)")
		api.Use(middlewares.NewAPIKeyMiddleware(cfg.APIKey))
	} else {
		log.Println("[Security] API_KEY tidak disetel, endpoint /api/v1 bersifat publik (tanpa autentikasi API Key)")
	}

	product.RegisterRoutes(api, productHandler)

	// 8. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("[Server] Menerima sinyal shutdown, mematikan server secara graceful...")
		if err := app.Shutdown(); err != nil {
			log.Printf("[Server] Error saat shutdown: %v\n", err)
		}
	}()

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("[Server] Server mendengarkan di http://localhost%s\n", addr)
	if err := app.Listen(addr); err != nil {
		log.Printf("[Server] Server berhenti: %v\n", err)
	}
}
