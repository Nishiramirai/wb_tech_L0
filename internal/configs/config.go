package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type AppConfig struct {
	App
	Kafka
	Database
}

type App struct {
	Port      int
	CacheSize int
}

type Kafka struct {
	Brokers []string
	Topic   string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Возвращает конфиг приложения и нужных сервисов, считанных с переменных окружения или .env файла
func NewConfig(logger *zap.Logger) (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		logger.Info(".env file not found.", zap.Error(err))
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("APP_PORT is not defined or invalid: %w", err)
	}

	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		return nil, fmt.Errorf("CACHE_SIZE is not defined or invalid: %w", err)
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS is not defined ")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		return nil, fmt.Errorf("KAFKA_TOPIC is not defined")
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("POSTGRES_HOST is not defined")
	}

	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return nil, fmt.Errorf("POSTGRES_PORT is not defined: %w", err)
	}

	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("POSTGRES_USER is not defined")
	}

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD is not defined")
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		return nil, fmt.Errorf("POSTGRES_DB is not defined")
	}

	return &AppConfig{
		App: App{
			Port:      appPort,
			CacheSize: cacheSize,
		},
		Kafka: Kafka{
			Brokers: strings.Split(kafkaBrokers, ","),
			Topic:   kafkaTopic,
		},
		Database: Database{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  "disable",
		},
	}, nil
}
