// internal/repository/repository.go
package repository

import (
	"context"
	"database/sql"
	"github.com/afzalabbasi/message-service/persistence-service/internal/config"
	"github.com/afzalabbasi/message-service/persistence-service/internal/models"
	_ "github.com/lib/pq"
)

// Repository handles database operations
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Repository
func NewRepository(cfg *config.Config) (*Repository, error) {
	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize database tables if they don't exist
	if err := initDB(db); err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// Initialize database tables
func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		username VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		room_id VARCHAR(36) NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id);
	CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
	CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
	`
	_, err := db.Exec(query)
	return err
}

// SaveMessage saves a message to the database
func (r *Repository) SaveMessage(ctx context.Context, message models.Message) error {
	query := `
	INSERT INTO messages (id, user_id, username, content, room_id, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO NOTHING
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		message.ID,
		message.UserID,
		message.Username,
		message.Content,
		message.RoomID,
		message.CreatedAt,
	)
	return err
}

// GetMessagesByRoom retrieves messages for a specific room
func (r *Repository) GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]models.Message, error) {
	query := `
	SELECT id, user_id, username, content, room_id, created_at
	FROM messages
	WHERE room_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.RoomID,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}
