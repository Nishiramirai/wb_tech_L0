package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App      `yaml:"app"`
	Kafka    `yaml:"kafka"`
	Database `yaml:"database"`
}

type App struct {
	Port      int `yaml:"port"`
	CacheSize int `yaml:"cache_size"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func NewConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func JoinBrokers(brokers []string) string {
	var result string
	for i, broker := range brokers {
		result += broker
		if i < len(brokers)-1 {
			result += ","
		}
	}
	return result
}

func CreateDBConnectionString(dbCfg Database) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.SSLMode)
}
