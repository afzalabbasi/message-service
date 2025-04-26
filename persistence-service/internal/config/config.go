// internal/config/config.go
package config

import (
	"os"
	"strings"
)

// Config holds all configuration for the persistence service
type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
	PostgresURL  string
	LogLevel     string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersStr == "" {
		kafkaBrokersStr = "kafka:9092"
	}
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "messages"
	}

	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "persistence-service"
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgresql://postgres:postgres@postgres:5432/messaging?sslmode=disable"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		KafkaBrokers: kafkaBrokers,
		KafkaTopic:   kafkaTopic,
		KafkaGroupID: kafkaGroupID,
		PostgresURL:  postgresURL,
		LogLevel:     logLevel,
	}, nil
}
