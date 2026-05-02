package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	BackendURL    string
	SessionStorePath string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		TelegramToken: strings.TrimSpace(os.Getenv("TELEGRAM_TOKEN")),
		BackendURL:    strings.TrimSpace(os.Getenv("BACKEND_URL")),
		SessionStorePath: strings.TrimSpace(os.Getenv("SESSION_STORE_PATH")),
	}

	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_TOKEN is required")
	}
	if cfg.BackendURL == "" {
		return nil, fmt.Errorf("BACKEND_URL is required")
	}

	return cfg, nil
}

