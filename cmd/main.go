package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"orders-service/internal/app/service"
	"orders-service/internal/cache"
	"orders-service/internal/configs"
	"orders-service/internal/db"
	"orders-service/internal/http"
	"orders-service/internal/kafka"

	"go.uber.org/zap"
)

// Main довольно объемный вышел, не получилось пока graceful shutdown как-то красиво вынести
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := configs.NewConfig(logger)
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	connStr := db.CreateDBConnectionString(cfg.Database)
	database, err := db.NewDB(connStr)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	orderCache, err := cache.NewCache(context.Background(), database, cfg.CacheSize)
	if err != nil {
		logger.Fatal("Failed to create cache", zap.Error(err))
	}

	orderService := service.NewOrderService(database, orderCache)

	consumer, err := kafka.NewConsumer(cfg, orderService, logger)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", zap.Error(err))
	}
	defer consumer.Close()

	server, err := http.NewServer(orderService, logger)
	if err != nil {
		logger.Fatal("Failed to create HTTP server", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Start(ctx)
		logger.Info("Kafka consumer has finished its work.")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		go server.Start(cfg.App.Port)
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP server shutdown failed", zap.Error(err))
		}
	}()

	logger.Info("Application is running. Press Ctrl+C to stop.")

	<-ctx.Done()

	logger.Info("Received shutdown signal. Starting graceful shutdown...")

	stop()

	wg.Wait()

	logger.Info("Application gracefully stopped.")
}
