package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eve-an/splitter/internal/cache"
	"github.com/eve-an/splitter/internal/config"
	"github.com/eve-an/splitter/internal/db"
	"github.com/eve-an/splitter/internal/feature"
	"github.com/eve-an/splitter/internal/http"
	"github.com/eve-an/splitter/internal/http/handler"
	"github.com/eve-an/splitter/internal/logger"
	"github.com/eve-an/splitter/internal/session"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	logger, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		log.Fatal("Error initializing logger")
	}

	database, err := db.New(config.Database)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	defer func() {
		logger.Info("database closed", slog.Any("error", database.Close()))
	}()

	featureRepo := feature.NewPostgresFeatureRepository(database.Pool, database.Queries)
	featureCache := cache.NewMemoryCache[*feature.Feature](time.Minute)

	eventRepo := feature.NewPostgresEventRepository(database.Queries)
	featureSvc := feature.NewService(featureRepo, eventRepo, featureCache)

	featureHandler := handler.NewFeatureHandler(logger, featureSvc)

	sessionSvc := session.NewService(time.Hour * 12)

	router := http.NewRouter(logger, featureHandler, sessionSvc, config.DefaultAuth) // TODO: support multi user auth
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
