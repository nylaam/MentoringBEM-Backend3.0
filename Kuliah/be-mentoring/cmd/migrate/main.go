package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mentoring/toko-backend/internal/config"
)

func main() {
	var direction string
	flag.StringVar(&direction, "direction", "up", "Arah migrasi: 'up' untuk migrasi maju, 'down' untuk rollback")
	flag.Parse()

	// Cek jika ada argumen positional: go run ./cmd/migrate down
	if flag.NArg() > 0 {
		direction = flag.Arg(0)
	}

	direction = strings.ToLower(direction)
	if direction != "up" && direction != "down" {
		log.Fatalf("Arah migrasi tidak valid: '%s'. Gunakan 'up' atau 'down'", direction)
	}

	log.Printf("[Migrate] Memulai database migration (%s)...\n", strings.ToUpper(direction))

	cfg := config.LoadConfig()
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal koneksi ke database Postgres: %v", err)
	}

	// Baca seluruh file migrasi sesuai arah
	pattern := fmt.Sprintf("migrations/*%s.sql", direction)
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Gagal membaca folder migrations: %v", err)
	}

	if len(files) == 0 {
		log.Printf("Tidak ditemukan file migrasi dengan pola: %s\n", pattern)
		return
	}

	// Jika up: urutkan naik (000001 -> 000002 -> 000003)
	// Jika down: urutkan turun (000003 -> 000002 -> 000001)
	if direction == "up" {
		sort.Strings(files)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(files)))
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Gagal membaca file %s: %v", file, err)
		}

		log.Printf("[Migrate] Menjalankan %s ...", file)
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("Gagal mengeksekusi migrasi %s: %v", file, err)
		}
		fmt.Printf(" Selesai: %s\n", file)
	}

	fmt.Printf("\n Migrasi database (%s) berhasil dijalankan! (%d file dieksekusi)\n", strings.ToUpper(direction), len(files))
}
