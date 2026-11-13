package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/mentoring/toko-backend/internal/app/order"
	"github.com/mentoring/toko-backend/internal/app/product"
	"github.com/mentoring/toko-backend/internal/config"
	"github.com/mentoring/toko-backend/internal/infra/database"
	"github.com/mentoring/toko-backend/internal/infra/repository"
	"github.com/mentoring/toko-backend/internal/infra/scheduler"
	"github.com/mentoring/toko-backend/internal/middlewares"
	"github.com/mentoring/toko-backend/pkg/response"
)

func main() {
	// 1. Load Configuration (12-Factor App)
	cfg := config.LoadConfig()
	log.Printf("[App] Memulai aplikasi di environment: %s", cfg.AppEnv)

	// 2. Inisialisasi Database Connection Pool (Infra)
	var db *sql.DB
	var err error
	db, err = database.NewPostgres(cfg)
	if err != nil {
		log.Printf("[Peringatan] Database Postgres tidak dapat dihubungi (%v). Berjalan dalam mode In-Memory Repository untuk hands-on lab!\n", err)
		db = nil
	} else {
		log.Println("[Database] Terhubung sukses ke PostgreSQL!")
		defer db.Close()
	}

	// 3. Inisialisasi Repositories (Mengimplementasikan Contracts / Interfaces)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// 4. Inisialisasi Services / Use Cases (Dependency Inversion)
	productSvc := product.NewProductService(productRepo)
	orderSvc := order.NewOrderService(orderRepo, productRepo)

	// 5. Inisialisasi Delivery Handlers (Go Fiber v3)
	productHandler := product.NewProductHandler(productSvc)
	orderHandler := order.NewOrderHandler(orderSvc)

	// Inisialisasi Background Scheduler (Infra)
	sched := scheduler.NewOrderScheduler()
	sched.Start(context.Background(), 5*time.Minute)
	defer sched.Stop()

	// 6. Setup Fiber v3 Engine & Centralized Error Handler (Slide 13)
	app := fiber.New(fiber.Config{
		AppName:      "Toko Backend API - Mentoring Hands-on (Fiber v3)",
		ErrorHandler: middlewares.CustomErrorHandler,
	})

	// Middlewares Stack
	app.Use(recover.New())
	app.Use(middlewares.NewRequestIDMiddleware())
	app.Use(middlewares.NewLoggerMiddleware())

	// Healthcheck endpoint (Slide 23)
	app.Get("/healthz", func(c fiber.Ctx) error {
		return response.JSONSuccess(c, fiber.StatusOK, fiber.Map{
			"status":  "healthy",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"version": "1.0.0",
		})
	})

	// 7. Route Registration
	api := app.Group("/api/v1")
	product.RegisterRoutes(api, productHandler)
	order.RegisterRoutes(api, orderHandler)

	// 8. Graceful Shutdown (Slide 23)
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
