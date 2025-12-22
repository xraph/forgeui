package bridge

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// WSHandler handles WebSocket connections
type WSHandler struct {
	bridge      *Bridge
	security    *Security
	connections sync.Map // map[string]*wsConnection
}

// wsConnection represents a WebSocket connection
type wsConnection struct {
	conn   *websocket.Conn
	ctx    Context
	send   chan []byte
	userID string
}

// NewWSHandler creates a new WebSocket handler
func NewWSHandler(bridge *Bridge) *WSHandler {
	return &WSHandler{
		bridge:   bridge,
		security: NewSecurity(bridge.config),
	}
}

// ServeHTTP handles WebSocket upgrade requests
func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: h.bridge.config.AllowedOrigins,
	})
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create bridge context
	ctx := NewContext(r)

	// Create connection
	wsConn := &wsConnection{
		conn: conn,
		ctx:  ctx,
		send: make(chan []byte, 256),
	}

	// Register connection
	connID := generateRequestID()
	h.connections.Store(connID, wsConn)

	// Handle connection
	go h.handleConnection(connID, wsConn)
}

// handleConnection manages a WebSocket connection
func (h *WSHandler) handleConnection(connID string, wsConn *wsConnection) {
	defer func() {
		h.connections.Delete(connID)
		wsConn.conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	// Start write pump
	go h.writePump(wsConn)

	// Read messages
	for {
		_, message, err := wsConn.conn.Read(context.Background())
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		// Process message
		h.processMessage(wsConn, message)
	}
}

// processMessage processes an incoming WebSocket message
func (h *WSHandler) processMessage(wsConn *wsConnection, message []byte) {
	var req Request
	if err := json.Unmarshal(message, &req); err != nil {
		h.sendError(wsConn, nil, ErrInvalidRequest)
		return
	}

	// Get function
	fn, err := h.bridge.GetFunction(req.Method)
	if err != nil {
		h.sendError(wsConn, req.ID, ErrMethodNotFound)
		return
	}

	// Check authentication
	if err := h.security.CheckAuth(wsConn.ctx, fn); err != nil {
		if bridgeErr, ok := err.(*Error); ok {
			h.sendError(wsConn, req.ID, bridgeErr)
		} else {
			h.sendError(wsConn, req.ID, ErrUnauthorized)
		}
		return
	}

	// Execute function
	result := h.bridge.execute(wsConn.ctx, req.Method, req.Params)

	// Build response
	resp := Response{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	if result.Error != nil {
		resp.Error = result.Error
	} else {
		resp.Result = result.Result
	}

	// Send response
	data, _ := json.Marshal(resp)
	wsConn.send <- data
}

// sendError sends an error response
func (h *WSHandler) sendError(wsConn *wsConnection, id any, err *Error) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}

	data, _ := json.Marshal(resp)
	wsConn.send <- data
}

// writePump sends messages to the WebSocket
func (h *WSHandler) writePump(wsConn *wsConnection) {
	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-wsConn.send:
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := wsConn.conn.Write(ctx, websocket.MessageText, message)
			cancel()

			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			// Send ping
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := wsConn.conn.Ping(ctx)
			cancel()

			if err != nil {
				return
			}
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *WSHandler) Broadcast(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	h.connections.Range(func(key, value any) bool {
		if wsConn, ok := value.(*wsConnection); ok {
			select {
			case wsConn.send <- data:
			default:
				// Channel full, skip
			}
		}
		return true
	})
}

// SendToUser sends a message to a specific user
func (h *WSHandler) SendToUser(userID string, event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	h.connections.Range(func(key, value any) bool {
		if wsConn, ok := value.(*wsConnection); ok {
			if wsConn.userID == userID {
				select {
				case wsConn.send <- data:
				default:
					// Channel full, skip
				}
			}
		}
		return true
	})
}

