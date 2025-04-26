// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all configuration for the authentication service
type Config struct {
	ServerAddress string
	JWTSecret     string
	JWTExpiration time.Duration
	PostgresURL   string
	LogLevel      string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8081"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	jwtExpStr := os.Getenv("JWT_EXPIRATION")
	jwtExp := 24 * time.Hour
	if jwtExpStr != "" {
		var err error
		jwtExp, err = time.ParseDuration(jwtExpStr)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRATION format: %v", err)
		}
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
		ServerAddress: serverAddr,
		JWTSecret:     jwtSecret,
		JWTExpiration: jwtExp,
		PostgresURL:   postgresURL,
		LogLevel:      logLevel,
	}, nil
}
