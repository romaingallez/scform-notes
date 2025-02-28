package handlers

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var (
	// Connections holds all active WebSocket connections
	connections    = make(map[*websocket.Conn]bool)
	connectionsMux sync.Mutex
)

// ProgressUpdate represents a progress update message
type ProgressUpdate struct {
	Status   string  `json:"status"`
	Message  string  `json:"message"`
	Progress float64 `json:"progress"`
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(c *websocket.Conn) {
	// Register new connection
	connectionsMux.Lock()
	connections[c] = true
	connectionsMux.Unlock()

	defer func() {
		// Unregister connection on close
		connectionsMux.Lock()
		delete(connections, c)
		connectionsMux.Unlock()
		c.Close()
	}()

	// Keep connection alive and handle incoming messages
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
	}
}

// BroadcastProgress sends a progress update to all connected clients
func BroadcastProgress(update interface{}) {
	data, err := json.Marshal(update)
	if err != nil {
		log.Printf("error marshaling progress update: %v", err)
		return
	}

	connectionsMux.Lock()
	defer connectionsMux.Unlock()

	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("error writing message: %v", err)
			conn.Close()
			delete(connections, conn)
		}
	}
}
