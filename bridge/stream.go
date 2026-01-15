package bridge

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SSEHandler handles Server-Sent Events streaming
type SSEHandler struct {
	bridge   *Bridge
	security *Security
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(bridge *Bridge) *SSEHandler {
	return &SSEHandler{
		bridge:   bridge,
		security: NewSecurity(bridge.config),
	}
}

// ServeHTTP handles SSE requests
func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Create bridge context
	ctx := NewContext(r)

	// Get function name from query
	funcName := r.URL.Query().Get("method")
	if funcName == "" {
		h.sendError(w, flusher, "Method parameter required")
		return
	}

	// Get params from query
	paramsJSON := r.URL.Query().Get("params")

	var params json.RawMessage
	if paramsJSON != "" {
		params = json.RawMessage(paramsJSON)
	}

	// Get function
	fn, err := h.bridge.GetFunction(funcName)
	if err != nil {
		h.sendError(w, flusher, "Method not found")
		return
	}

	// Check authentication
	if err := h.security.CheckAuth(ctx, fn); err != nil {
		h.sendError(w, flusher, "Unauthorized")
		return
	}

	// Execute function and stream results
	h.streamExecution(w, flusher, ctx, funcName, params)
}

// streamExecution executes a function and streams the results
func (h *SSEHandler) streamExecution(w http.ResponseWriter, flusher http.Flusher, ctx Context, funcName string, params json.RawMessage) {
	// Create a channel for streaming
	streamChan := make(chan StreamChunk, 10)

	// Execute in goroutine
	go func() {
		defer close(streamChan)

		// Execute function
		result := h.bridge.execute(ctx, funcName, params)

		// Send result
		chunk := StreamChunk{
			Done: true,
		}

		if result.Error != nil {
			chunk.Error = result.Error
		} else {
			chunk.Data = result.Result
		}

		streamChan <- chunk
	}()

	// Stream chunks to client
	for chunk := range streamChan {
		data, err := json.Marshal(chunk)
		if err != nil {
			continue
		}

		_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()

		if chunk.Done {
			break
		}
	}
}

// sendError sends an error event
func (h *SSEHandler) sendError(w http.ResponseWriter, flusher http.Flusher, message string) {
	chunk := StreamChunk{
		Error: NewError(ErrCodeInternal, message),
		Done:  true,
	}

	data, err := json.Marshal(chunk)
	if err != nil {
		// Fallback to simple error message if marshal fails
		_, _ = fmt.Fprintf(w, "data: {\"error\":{\"code\":-32603,\"message\":\"%s\"},\"done\":true}\n\n", message)
	} else {
		_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	}

	flusher.Flush()
}

// StreamHandler returns an SSE handler for the bridge
func (b *Bridge) StreamHandler() http.Handler {
	return NewSSEHandler(b)
}

// StreamEvent sends an event to SSE clients
type StreamEvent struct {
	Event string `json:"event,omitempty"`
	Data  any    `json:"data"`
	ID    string `json:"id,omitempty"`
}

// WriteSSE writes an SSE event
func WriteSSE(w http.ResponseWriter, event StreamEvent) error {
	if event.Event != "" {
		_, _ = fmt.Fprintf(w, "event: %s\n", event.Event)
	}

	if event.ID != "" {
		_, _ = fmt.Fprintf(w, "id: %s\n", event.ID)
	}

	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}
