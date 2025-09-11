package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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

func NewConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден. Используем переменные окружения напрямую.")
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("APP_PORT не задана или не число: %w", err)
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS не задана")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		return nil, fmt.Errorf("KAFKA_TOPIC не задана")
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("POSTGRES_HOST не задана")
	}

	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return nil, fmt.Errorf("POSTGRES_PORT не задана или не число: %w", err)
	}

	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("POSTGRES_USER не задана")
	}

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD не задана")
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		return nil, fmt.Errorf("POSTGRES_DB не задана")
	}

	return &AppConfig{
		App: App{
			Port: appPort,
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
