package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mentoring/toko-backend/domain/entity"
	"github.com/mentoring/toko-backend/internal/config"
	"github.com/mentoring/toko-backend/internal/infra/database"
	"github.com/mentoring/toko-backend/internal/infra/repository"
)

func main() {
	log.Println("[Seed] Memulai seeding data toko...")

	cfg := config.LoadConfig()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database Postgres: %v", err)
	}
	defer db.Close()

	repo := repository.NewProductRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sampleProducts := []*entity.Product{
		{ID: "prod-001", Name: "Laptop ASUS ROG Zephyrus", Price: 28_000_000, Stock: 15, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "prod-002", Name: "Mechanical Keyboard Keychron K2", Price: 1_450_000, Stock: 40, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "prod-003", Name: "Logitech MX Master 3S Mouse", Price: 1_200_000, Stock: 60, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "prod-004", Name: "Monitor Dell UltraSharp 27 Inch 4K", Price: 9_500_000, Stock: 20, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "prod-005", Name: "Enterprise AI Server Rack", Price: 65_000_000, Stock: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, p := range sampleProducts {
		if err := repo.Save(ctx, p); err != nil {
			log.Printf("Peringatan: Gagal seed produk %s: %v\n", p.Name, err)
		} else {
			fmt.Printf(" Berhasil seed produk: %s (Rp %d)\n", p.Name, p.Price)
		}
	}

	fmt.Println(" Seeding selesai!")
}
