package handlers

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var (
	// Connections holds WebSocket connections organized by session ID
	connections    = make(map[string]map[*websocket.Conn]bool)
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
	// Get session from the HTTP context (passed via Locals)
	sessionID := c.Locals("session_id")
	if sessionID == nil {
		log.Println("No session ID found for WebSocket connection")
		c.Close()
		return
	}

	sessionIDStr, ok := sessionID.(string)
	if !ok {
		log.Println("Invalid session ID type for WebSocket connection")
		c.Close()
		return
	}

	// Register new connection
	connectionsMux.Lock()
	if connections[sessionIDStr] == nil {
		connections[sessionIDStr] = make(map[*websocket.Conn]bool)
	}
	connections[sessionIDStr][c] = true
	connectionsMux.Unlock()

	log.Printf("WebSocket connection established for session: %s", sessionIDStr)

	defer func() {
		// Unregister connection on close
		connectionsMux.Lock()
		if connections[sessionIDStr] != nil {
			delete(connections[sessionIDStr], c)
			// Clean up empty session maps
			if len(connections[sessionIDStr]) == 0 {
				delete(connections, sessionIDStr)
			}
		}
		connectionsMux.Unlock()
		c.Close()
		log.Printf("WebSocket connection closed for session: %s", sessionIDStr)
	}()

	// Keep connection alive and handle incoming messages
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message for session %s: %v", sessionIDStr, err)
			}
			break
		}
	}
}

// BroadcastProgress sends a progress update to all connected clients (legacy function for backwards compatibility)
func BroadcastProgress(update interface{}) {
	BroadcastProgressToAll(update)
}

// BroadcastProgressToAll sends a progress update to all connected clients
func BroadcastProgressToAll(update interface{}) {
	data, err := json.Marshal(update)
	if err != nil {
		log.Printf("error marshaling progress update: %v", err)
		return
	}

	connectionsMux.Lock()
	defer connectionsMux.Unlock()

	for sessionID, sessionConns := range connections {
		for conn := range sessionConns {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("error writing message to session %s: %v", sessionID, err)
				conn.Close()
				delete(sessionConns, conn)
			}
		}
		// Clean up empty session maps
		if len(sessionConns) == 0 {
			delete(connections, sessionID)
		}
	}
}

// BroadcastProgressToSession sends a progress update to all connections of a specific session
func BroadcastProgressToSession(sessionID string, update interface{}) {
	data, err := json.Marshal(update)
	if err != nil {
		log.Printf("error marshaling progress update: %v", err)
		return
	}

	connectionsMux.Lock()
	defer connectionsMux.Unlock()

	sessionConns, exists := connections[sessionID]
	if !exists {
		log.Printf("no connections found for session %s", sessionID)
		return
	}

	for conn := range sessionConns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("error writing message to session %s: %v", sessionID, err)
			conn.Close()
			delete(sessionConns, conn)
		}
	}

	// Clean up empty session maps
	if len(sessionConns) == 0 {
		delete(connections, sessionID)
	}
}

// GetSessionIDFromContext helper function to get session ID from WebSocket context
func GetSessionIDFromContext(c *websocket.Conn) string {
	sessionID := c.Locals("session_id")
	if sessionID == nil {
		return ""
	}
	sessionIDStr, ok := sessionID.(string)
	if !ok {
		return ""
	}
	return sessionIDStr
}
