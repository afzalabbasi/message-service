// internal/kafka/producer.go
package kafka

import (
	"context"
	"encoding/json"
	"github.com/afzalabbasi/message-service/webSocket/internal/config"
	"github.com/afzalabbasi/message-service/webSocket/internal/models"
	"github.com/segmentio/kafka-go"

	"log"
)

// Producer represents a Kafka producer
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.Config) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaBrokers...),
		Topic:        cfg.KafkaProducerTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return &Producer{
		writer: writer,
	}, nil
}

// PublishMessage publishes a message to Kafka
func (p *Producer) PublishMessage(message models.Message) error {
	kafkaMsg := models.KafkaMessage{
		MessageID: message.ID,
		UserID:    message.UserID,
		Username:  message.Username,
		Content:   message.Content,
		RoomID:    message.RoomID,
		Timestamp: message.CreatedAt,
		EventType: "message_created",
	}

	value, err := json.Marshal(kafkaMsg)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(kafkaMsg.RoomID),
			Value: value,
		},
	)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}

	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.writer.Close()
}
