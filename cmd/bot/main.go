package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bot-architecture/internal/bot"
	"bot-architecture/internal/config"
)

func main() {
	slog.Info("starting bot")
	cfg, err := config.Load("config.toml")
	if err != nil {
		slog.Error("failed to load configuration", "err", err)
		os.Exit(1)
	}

	app, err := bot.New(cfg)
	if err != nil {
		slog.Error("failed to create bot", "err", err)
		os.Exit(1)
	}

	go func() {
		if startErr := app.SetupWebhook(); startErr != nil {
			slog.Error("bot exited with error", "err", startErr)
			os.Exit(1)
		}
	}()

	// wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	slog.Info("waiting for interrupt signal")

	<-sigCh

	slog.Info("shutting down bot")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Shutdown(shutdownCtx)
}
