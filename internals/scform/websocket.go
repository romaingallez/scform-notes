package scform

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// WebSocket is a custom websocket that uses gobwas/ws as the transport layer.
type WebSocket struct {
	conn net.Conn
}

// NewWebSocket creates a new WebSocket connection
func NewWebSocket(u string) *WebSocket {
	log.Printf("Attempting to establish WebSocket connection to: %s", u)
	conn, _, _, err := ws.Dial(context.Background(), u)
	if err != nil {
		DebugLog("WebSocket connection failed: %v", err)
		return nil
	}
	DebugLog("WebSocket connection established successfully")
	return &WebSocket{conn}
}

// Send sends data through the WebSocket connection
func (w *WebSocket) Send(b []byte) error {
	if w == nil || w.conn == nil {
		return errors.New("websocket connection is not established")
	}
	// debugLog("Sending WebSocket message, length: %d bytes", len(b))
	err := wsutil.WriteClientText(w.conn, b)
	if err != nil {
		DebugLog("Error sending WebSocket message: %v", err)
		return err
	}
	// debugLog("WebSocket message sent successfully")
	return nil
}

// Read reads data from the WebSocket connection
func (w *WebSocket) Read() ([]byte, error) {
	if w == nil || w.conn == nil {
		return nil, errors.New("websocket connection is not established")
	}
	// debugLog("Reading from WebSocket connection...")
	data, err := wsutil.ReadServerText(w.conn)
	if err != nil {
		DebugLog("Error reading from WebSocket: %v", err)
		return nil, err
	}
	// debugLog("Received WebSocket message, length: %d bytes", len(data))
	return data, nil
}

// Close closes the WebSocket connection
func (w *WebSocket) Close() error {
	if w == nil || w.conn == nil {
		return nil
	}
	return w.conn.Close()
}
