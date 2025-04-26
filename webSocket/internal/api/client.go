// internal/api/client.go
package api

import (
	"encoding/json"
	"github.com/afzalabbasi/message-service/webSocket/internal/models"

	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 4096
)

// Client represents a WebSocket client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan models.Message
	roomID   string
	userID   string
	username string
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse message content
		var messageContent struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(data, &messageContent); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		// Create a new message
		message := models.Message{
			ID:        uuid.New().String(),
			UserID:    c.userID,
			Username:  c.username,
			Content:   messageContent.Content,
			RoomID:    c.roomID,
			CreatedAt: time.Now(),
		}

		// Publish message to Kafka
		if err := c.hub.PublishMessage(message); err != nil {
			log.Printf("error publishing message: %v", err)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write message to WebSocket
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling message: %v", err)
				return
			}
			w.Write(messageBytes)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
