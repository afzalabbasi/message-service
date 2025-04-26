// internal/api/hub.go
package api

import (
	kafka "github.com/afzalabbasi/message-service/webSocket/internal/kakfa"
	"github.com/afzalabbasi/message-service/webSocket/internal/models"
	"log"
	"sync"
)

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients by room
	clients map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Kafka producer for publishing messages
	kafkaProducer *kafka.Producer

	// Mutex for thread-safe access to clients map
	mu sync.RWMutex
}

// NewHub creates a new hub
func NewHub(kafkaProducer *kafka.Producer) *Hub {
	return &Hub{
		clients:       make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		kafkaProducer: kafkaProducer,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// Initialize room if it doesn't exist
			if _, ok := h.clients[client.roomID]; !ok {
				h.clients[client.roomID] = make(map[*Client]bool)
			}
			h.clients[client.roomID][client] = true
			h.mu.Unlock()
			log.Printf("Client connected: %s in room %s", client.userID, client.roomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.roomID]; ok {
				if _, ok := h.clients[client.roomID][client]; ok {
					delete(h.clients[client.roomID], client)
					close(client.send)
					log.Printf("Client disconnected: %s from room %s", client.userID, client.roomID)

					// Remove room if empty
					if len(h.clients[client.roomID]) == 0 {
						delete(h.clients, client.roomID)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends a message to all clients in a room
func (h *Hub) Broadcast(message models.Message, roomID string) {
	h.mu.RLock()
	if clients, ok := h.clients[roomID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				h.mu.RUnlock()
				h.mu.Lock()
				delete(h.clients[roomID], client)
				h.mu.Unlock()
				h.mu.RLock()
			}
		}
	}
	h.mu.RUnlock()
}

// PublishMessage publishes a message to Kafka
func (h *Hub) PublishMessage(message models.Message) error {
	return h.kafkaProducer.PublishMessage(message)
}
