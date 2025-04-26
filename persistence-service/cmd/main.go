// cmd/main.go
package main

import (
	"context"
	"github.com/afzalabbasi/message-service/persistence-service/internal/config"
	"github.com/afzalabbasi/message-service/persistence-service/internal/kafka"
	"github.com/afzalabbasi/message-service/persistence-service/internal/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repository
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize and start Kafka consumer
	consumer, err := kafka.NewConsumer(cfg, repo)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Create a context that is canceled when the program is interrupted
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming messages from Kafka in a goroutine
	go func() {
		if err := consumer.Consume(ctx); err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	log.Println("Persistence service started...")

	// Wait for interrupt signal to gracefully shut down the service
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down service...")

	// Cancel the context to signal to the consumer to stop
	cancel()
	log.Println("Service exiting")
}
