// internal/api/handler.go
package api

import (
	"encoding/json"
	"github.com/afzalabbasi/message-service/webSocket/internal/config"
	"github.com/afzalabbasi/message-service/webSocket/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// JWTClaims defines the JWT claims structure
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Handler handles HTTP requests for the WebSocket service
type Handler struct {
	hub      *Hub
	config   *config.Config
	upgrader websocket.Upgrader
}

// NewHandler creates a new Handler
func NewHandler(hub *Hub, cfg *config.Config) *Handler {
	return &Handler{
		hub:    hub,
		config: cfg,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development; restrict in production
			},
		},
	}
}

// SetupRoutes sets up the routes for the WebSocket service
func (h *Handler) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Routes
	r.Get("/ws/{roomID}", h.handleWebSocket)
	r.Get("/health", h.healthCheck)

	return r
}

// validateToken validates a JWT token
func (h *Handler) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// handleWebSocket handles WebSocket connections
func (h *Handler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	// Get token from query param
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Authentication token is required", http.StatusUnauthorized)
		return
	}

	// Validate token
	claims, err := h.validateToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create a new client
	client := &Client{
		hub:      h.hub,
		conn:     conn,
		send:     make(chan models.Message, 256),
		roomID:   roomID,
		userID:   claims.UserID,
		username: claims.Username,
	}

	// Register client with hub
	h.hub.register <- client

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// healthCheck handles health checks
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
