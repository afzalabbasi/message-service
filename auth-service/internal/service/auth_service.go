// internal/service/auth_service.go
package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/afzalabbasi/message-service/auth-service/internal/config"
	"github.com/afzalabbasi/message-service/auth-service/internal/middleware"
	models "github.com/afzalabbasi/message-service/auth-service/internal/model"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// AuthService handles authentication business logic
type AuthService struct {
	db            *sql.DB
	jwtMiddleware *middleware.JWTMiddleware
	config        *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *config.Config) *AuthService {
	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		panic(err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Initialize database tables if they don't exist
	if err := initDB(db); err != nil {
		panic(err)
	}

	return &AuthService{
		db:            db,
		jwtMiddleware: middleware.NewJWTMiddleware(cfg.JWTSecret),
		config:        cfg,
	}
}

// Initialize database tables
func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
	_, err := db.Exec(query)
	return err
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req models.AuthRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2", req.Username, req.Email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new UUID for the user
	userID := uuid.New().String()

	// Insert the new user
	_, err = s.db.ExecContext(ctx,
		"INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, req.Username, req.Email, hashedPassword, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	// Generate a JWT token
	token, err := s.jwtMiddleware.GenerateToken(userID, req.Username, s.config.JWTExpiration)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:    token,
		Username: req.Username,
		UserID:   userID,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req models.AuthRequest) (*models.AuthResponse, error) {
	var user models.User

	// Query by username or email
	var err error
	if req.Username != "" {
		err = s.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username = $1", req.Username).
			Scan(&user.ID, &user.Username, &user.Password)
	} else {
		err = s.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE email = $1", req.Email).
			Scan(&user.ID, &user.Username, &user.Password)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verify password
	if !models.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate a JWT token
	token, err := s.jwtMiddleware.GenerateToken(user.ID, user.Username, s.config.JWTExpiration)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	}, nil
}
