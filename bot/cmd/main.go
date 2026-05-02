package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"receptor/bot/internal/client"
	"receptor/bot/internal/config"
	"receptor/bot/internal/handler"
	"receptor/bot/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	backendHTTP := &http.Client{Timeout: 25 * time.Second}
	backendClient, err := client.NewBackendClient(cfg.BackendURL, backendHTTP)
	if err != nil {
		log.Fatalf("backend client: %v", err)
	}

	botService, err := service.NewBotService(backendClient)
	if err != nil {
		log.Fatalf("bot service: %v", err)
	}

	telegramHTTP := &http.Client{Timeout: 35 * time.Second}
	store, err := handler.NewSessionStore(cfg.SessionStorePath)
	if err != nil {
		log.Fatalf("session store: %v", err)
	}

	tgHandler, err := handler.NewTelegramHandler(cfg.TelegramToken, telegramHTTP, botService, store)
	if err != nil {
		log.Fatalf("telegram handler: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := tgHandler.RunLongPolling(ctx); err != nil && err != context.Canceled {
		log.Fatalf("bot stopped: %v", err)
	}
}

