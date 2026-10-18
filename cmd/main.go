package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/aleksandrboro/hookflow/internal/config"
	"github.com/aleksandrboro/hookflow/internal/handler"
	"github.com/aleksandrboro/hookflow/internal/router"
	"github.com/aleksandrboro/hookflow/internal/server"
	"github.com/aleksandrboro/hookflow/internal/service"
	"github.com/aleksandrboro/hookflow/internal/storage"
)

const (
	levelInfo  = "info"
	levelDebug = "debug"
	levelError = "error"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.LogLevel)

	log.Info("config", slog.Any("cfg", cfg))

	store := storage.NewStore()

	service := service.NewJobService(store)

	handler := handler.NewJobHandler(service, log)

	router := router.NewRouter(handler)

	server, err := server.NewServer(log, cfg.Server, router)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := server.Start(); err != nil {
			log.Error("failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	<-sign

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	cancel()
}

func setupLogger(loglevel string) *slog.Logger {
	var log *slog.Logger

	switch loglevel {
	case levelInfo:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case levelDebug:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case levelError:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
