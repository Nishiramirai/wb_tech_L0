package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"orders-service/internal/app/service"
	"orders-service/internal/cache"
	"orders-service/internal/configs"
	"orders-service/internal/db"
	"orders-service/internal/http"
	"orders-service/internal/kafka"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := configs.NewConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	connStr := configs.CreateDBConnectionString(cfg.Database)
	database, err := db.NewDB(connStr)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orderCache, err := cache.NewCache(ctx, database, cfg.App.CacheSize)
	if err != nil {
		logger.Fatal("Failed to create cache", zap.Error(err))
	}

	orderService := service.NewOrderService(database, orderCache)

	consumer, err := kafka.NewConsumer(cfg, orderService, logger)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", zap.Error(err))
	}
	go consumer.Start(ctx)

	server, err := http.NewServer(orderService, logger)
	if err != nil {
		logger.Fatal("Failed to create HTTP server", zap.Error(err))
	}
	go server.Start(cfg.App.Port)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	logger.Info("Shutting down the service...")

	cancel()
	consumer.Close()

	logger.Info("Service gracefully stopped.")
}
