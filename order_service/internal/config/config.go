package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
}

func LoadConfig() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8082"))
	return &Config{
		Port:        port,
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:Roshik8956@localhost:5432/orders?sslmode=disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
