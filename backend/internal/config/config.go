package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds runtime settings from the environment (see backend_spec §2).
type Config struct {
	DBURL     string
	JWTSecret string
	Port      string
	MLURL     string
}

// Load reads optional .env then required variables.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBURL:     strings.TrimSpace(os.Getenv("DB_URL")),
		JWTSecret: strings.TrimSpace(os.Getenv("JWT_SECRET")),
		Port:      strings.TrimSpace(os.Getenv("PORT")),
		MLURL:     strings.TrimSpace(os.Getenv("ML_URL")),
	}
	if cfg.DBURL == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg, nil
}
