// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the WebSocket service
type Config struct {
	ServerAddress      string
	KafkaBrokers       []string
	KafkaConsumerTopic string
	KafkaProducerTopic string
	KafkaGroupID       string
	AuthServiceURL     string
	JWTSecret          string
	LogLevel           string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8082"
	}

	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersStr == "" {
		kafkaBrokersStr = "kafka:9092"
	}
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	kafkaConsumerTopic := os.Getenv("KAFKA_CONSUMER_TOPIC")
	if kafkaConsumerTopic == "" {
		kafkaConsumerTopic = "messages"
	}

	kafkaProducerTopic := os.Getenv("KAFKA_PRODUCER_TOPIC")
	if kafkaProducerTopic == "" {
		kafkaProducerTopic = "messages"
	}

	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "websocket-service"
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service:8081"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		ServerAddress:      serverAddr,
		KafkaBrokers:       kafkaBrokers,
		KafkaConsumerTopic: kafkaConsumerTopic,
		KafkaProducerTopic: kafkaProducerTopic,
		KafkaGroupID:       kafkaGroupID,
		AuthServiceURL:     authServiceURL,
		JWTSecret:          jwtSecret,
		LogLevel:           logLevel,
	}, nil
}
