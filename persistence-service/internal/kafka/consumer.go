// internal/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"github.com/afzalabbasi/message-service/persistence-service/internal/config"
	"github.com/afzalabbasi/message-service/persistence-service/internal/models"
	"github.com/afzalabbasi/message-service/persistence-service/internal/repository"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	reader *kafka.Reader
	repo   *repository.Repository
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, repo *repository.Repository) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.KafkaBrokers,
		Topic:       cfg.KafkaTopic,
		GroupID:     cfg.KafkaGroupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
		MaxWait:     500 * time.Millisecond,
	})

	return &Consumer{
		reader: reader,
		repo:   repo,
	}, nil
}

// Consume consumes messages from Kafka and persists them to the database
func (c *Consumer) Consume(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Use context with a timeout to make the read interruptible
			readCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			m, err := c.reader.ReadMessage(readCtx)
			cancel()

			if err != nil {
				// If context was canceled, return without error
				if ctx.Err() != nil {
					return nil
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			var kafkaMsg models.KafkaMessage
			if err := json.Unmarshal(m.Value, &kafkaMsg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Process message based on event type
			switch kafkaMsg.EventType {
			case "message_created":
				message := models.Message{
					ID:        kafkaMsg.MessageID,
					UserID:    kafkaMsg.UserID,
					Username:  kafkaMsg.Username,
					Content:   kafkaMsg.Content,
					RoomID:    kafkaMsg.RoomID,
					CreatedAt: kafkaMsg.Timestamp,
				}

				if err := c.repo.SaveMessage(ctx, message); err != nil {
					log.Printf("Error saving message: %v", err)
					continue
				}
				log.Printf("Message saved: %s", kafkaMsg.MessageID)

			case "message_updated":
				// Handle message updates if needed
				log.Printf("Message update not implemented yet: %s", kafkaMsg.MessageID)

			default:
				log.Printf("Unknown event type: %s", kafkaMsg.EventType)
			}
		}
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
