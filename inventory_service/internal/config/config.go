package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:Roshik8956@localhost:5432/products?sslmode=disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
