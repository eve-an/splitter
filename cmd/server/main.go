package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eve-an/splitter/internal/config"
	"github.com/eve-an/splitter/internal/http"
	"github.com/eve-an/splitter/internal/http/handler"
	"github.com/eve-an/splitter/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config")
	}

	logger, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		log.Fatal("Error initializing logger")
	}

	featureHandler := handler.NewFeatureHandler(logger)

	router := http.NewRouter(logger, featureHandler)
	server := http.NewServer(config.ServerConifg, logger, router)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		server.Start()
	}()

	<-stop
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Info("server shutdown failed", slog.Any("error", err))
	}

	logger.Info("server stopped gracefully")
}
