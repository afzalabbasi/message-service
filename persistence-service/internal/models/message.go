// internal/models/message.go
package models

import (
	"time"
)

// Message represents a chat message
type Message struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Username  string    `db:"username"`
	Content   string    `db:"content"`
	RoomID    string    `db:"room_id"`
	CreatedAt time.Time `db:"created_at"`
}

// KafkaMessage represents a message that is consumed from Kafka
type KafkaMessage struct {
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"` // e.g., "message_created", "message_updated"
}
