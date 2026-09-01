package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	GRPCPort    string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	cfg := Config{
		Port:        getEnv("PORT", ":8080"),
		GRPCPort:    getEnv("GRPC_PORT", ":9090"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/go101"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key"),
	}

	log.Printf("loaded config: http=%s grpc=%s", cfg.Port, cfg.GRPCPort)
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
