// internal/models/message.go
package models

import (
	"time"
)

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
}

// KafkaMessage represents a message that is published/consumed to/from Kafka
type KafkaMessage struct {
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"` // e.g., "message_created", "message_updated"
}
