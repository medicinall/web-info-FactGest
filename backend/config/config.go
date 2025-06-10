package config

import (
	"os"
)

type Config struct {
	Port     string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "3306"),
		DBUser:   getEnv("DB_USER", "factugest"),
		DBPass:   getEnv("DB_PASS", "factugest123"),
		DBName:   getEnv("DB_NAME", "factugest_db"),
		JWTSecret: getEnv("JWT_SECRET", "factugest-secret-key-2024"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

