// internal/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"github.com/afzalabbasi/message-service/webSocket/internal/api"
	"github.com/afzalabbasi/message-service/webSocket/internal/config"
	"github.com/afzalabbasi/message-service/webSocket/internal/models"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.KafkaBrokers,
		Topic:       cfg.KafkaConsumerTopic,
		GroupID:     cfg.KafkaGroupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
		MaxWait:     500 * time.Millisecond,
	})

	return &Consumer{
		reader: reader,
	}, nil
}

// Consume consumes messages from Kafka and sends them to the WebSocket hub
func (c *Consumer) Consume(hub *api.Hub) error {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var kafkaMsg models.KafkaMessage
		if err := json.Unmarshal(m.Value, &kafkaMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Convert KafkaMessage to Message and broadcast to all clients in the room
		message := models.Message{
			ID:        kafkaMsg.MessageID,
			UserID:    kafkaMsg.UserID,
			Username:  kafkaMsg.Username,
			Content:   kafkaMsg.Content,
			RoomID:    kafkaMsg.RoomID,
			CreatedAt: kafkaMsg.Timestamp,
		}

		hub.Broadcast(message, kafkaMsg.RoomID)
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
