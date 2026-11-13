package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	AppEnv     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file tidak ditemukan, membaca dari environment variables")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBHost:     getEnv("DATABASE_HOST", "localhost"),
		DBPort:     getEnv("DATABASE_PORT", "5432"),
		DBUser:     getEnv("DATABASE_USER", "postgres"),
		DBPassword: getEnv("DATABASE_PASSWORD", "postgres"),
		DBName:     getEnv("DATABASE_NAME", "toko_db"),
		DBSSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
