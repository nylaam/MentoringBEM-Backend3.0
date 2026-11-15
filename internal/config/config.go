package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	AppEnv          string
	APIKey          string
	SwaggerUser     string
	SwaggerPassword string
	DatabaseURL     string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file tidak ditemukan, membaca dari environment variables")
	}

	// Dukungan standard PaaS / Cloud deployment: PORT (Railway, Render, Fly.io, Cloud Run)
	// Jika PORT tidak ada, fallback ke APP_PORT, jika tidak ada fallback ke default port 6767.
	port := getEnv("PORT", getEnv("APP_PORT", "6767"))

	return &Config{
		AppPort:         port,
		AppEnv:          getEnv("APP_ENV", "development"),
		APIKey:          getEnv("API_KEY", ""),
		SwaggerUser:     getEnv("SWAGGER_USER", "admin"),
		SwaggerPassword: getEnv("SWAGGER_PASSWORD", "admin123"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		DBHost:          getEnv("DATABASE_HOST", "localhost"),
		DBPort:          getEnv("DATABASE_PORT", "5432"),
		DBUser:          getEnv("DATABASE_USER", "postgres"),
		DBPassword:      getEnv("DATABASE_PASSWORD", "postgres"),
		DBName:          getEnv("DATABASE_NAME", "toko_db"),
		DBSSLMode:       getEnv("DATABASE_SSLMODE", "disable"),
	}
}

// DSN mengembalikan connection string PostgreSQL.
// Jika DATABASE_URL tersedia (format umum di cloud/deploy seperti Render, Supabase, Neon),
// maka DATABASE_URL akan langsung digunakan.
func (c *Config) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
