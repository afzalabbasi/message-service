// cmd/main.go
package main

import (
	"context"
	"github.com/afzalabbasi/message-service/webSocket/internal/api"
	"github.com/afzalabbasi/message-service/webSocket/internal/config"
	kafka "github.com/afzalabbasi/message-service/webSocket/internal/kakfa"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up Kafka consumer and producer
	kafkaConsumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	kafkaProducer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Initialize WebSocket hub
	hub := api.NewHub(kafkaProducer)
	go hub.Run()

	// Start consuming messages from Kafka in a goroutine
	go func() {
		if err := kafkaConsumer.Consume(hub); err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	// Initialize HTTP handler
	handler := api.NewHandler(hub, cfg)

	// Configure server
	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      handler.SetupRoutes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting WebSocket service on %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
