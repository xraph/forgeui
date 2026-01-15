package bridge

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// HTTPHandler implements http.Handler for bridge requests
type HTTPHandler struct {
	bridge   *Bridge
	security *Security
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(bridge *Bridge) *HTTPHandler {
	return &HTTPHandler{
		bridge:   bridge,
		security: NewSecurity(bridge.config),
	}
}

// ServeHTTP handles HTTP requests
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Handle CORS
	if h.bridge.config.EnableCORS {
		h.handleCORS(w, r)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Only accept POST requests
	if r.Method != http.MethodPost {
		h.writeError(w, nil, ErrInvalidRequest)
		return
	}

	// Check CSRF
	if err := h.security.CheckCSRF(r); err != nil {
		var bridgeErr *Error
		if errors.As(err, &bridgeErr) {
			h.writeError(w, nil, bridgeErr)
		} else {
			h.writeError(w, nil, ErrBadRequest)
		}

		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, nil, NewError(ErrCodeBadRequest, "Failed to read request body"))
		return
	}
	defer func() { _ = r.Body.Close() }()

	// Create bridge context
	ctx := NewContext(r)

	// Try to parse as single request first
	var singleReq Request
	if err := json.Unmarshal(body, &singleReq); err == nil && singleReq.Method != "" {
		// Single request
		h.handleSingleRequest(w, ctx, singleReq)
		return
	}

	// Try to parse as batch request
	var batchReq BatchRequest
	if err := json.Unmarshal(body, &batchReq); err == nil && len(batchReq) > 0 {
		// Batch request
		h.handleBatchRequest(w, ctx, batchReq)
		return
	}

	// Invalid request format
	h.writeError(w, nil, ErrInvalidRequest)
}

// handleSingleRequest handles a single RPC request
func (h *HTTPHandler) handleSingleRequest(w http.ResponseWriter, ctx Context, req Request) {
	// Validate JSON-RPC version
	if req.JSONRPC != "2.0" {
		h.writeError(w, req.ID, ErrInvalidRequest)
		return
	}

	// Get function
	fn, err := h.bridge.GetFunction(req.Method)
	if err != nil {
		h.writeError(w, req.ID, ErrMethodNotFound)
		return
	}

	// Check authentication
	if err := h.security.CheckAuth(ctx, fn); err != nil {
		var bridgeErr *Error
		if errors.As(err, &bridgeErr) {
			h.writeError(w, req.ID, bridgeErr)
		} else {
			h.writeError(w, req.ID, ErrUnauthorized)
		}

		return
	}

	// Check rate limit
	rateLimitKey := getRateLimitKey(ctx)
	if err := h.security.CheckRateLimit(rateLimitKey, fn); err != nil {
		var bridgeErr *Error
		if errors.As(err, &bridgeErr) {
			h.writeError(w, req.ID, bridgeErr)
		} else {
			h.writeError(w, req.ID, ErrRateLimit)
		}

		return
	}

	// Execute function
	result := h.bridge.execute(ctx, req.Method, req.Params)

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

	// Write response
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Log encoding error - response already sent
		_ = err // Error already logged by potential middleware
	}
}

// handleBatchRequest handles a batch of RPC requests
func (h *HTTPHandler) handleBatchRequest(w http.ResponseWriter, ctx Context, batch BatchRequest) {
	// Check batch size
	if len(batch) > h.bridge.config.MaxBatchSize {
		h.writeError(w, nil, NewError(ErrCodeBadRequest, "Batch size exceeds maximum"))
		return
	}

	// Convert to Request slice
	requests := make([]Request, len(batch))
	copy(requests, batch)

	// Execute batch
	responses := h.bridge.CallBatch(ctx, requests)

	// Write response
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responses); err != nil {
		// Log encoding error - response already sent
		_ = err // Error already logged by potential middleware
	}
}

// writeError writes an error response
func (h *HTTPHandler) writeError(w http.ResponseWriter, id any, err *Error) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}

	w.WriteHeader(http.StatusOK) // JSON-RPC errors use 200 OK

	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		// Log encoding error - response already sent
		_ = encodeErr // Error already logged by potential middleware
	}
}

// handleCORS sets CORS headers
func (h *HTTPHandler) handleCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	// Check if origin is allowed
	allowed := false

	for _, allowedOrigin := range h.bridge.config.AllowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			allowed = true
			break
		}
	}

	if !allowed {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+h.bridge.config.CSRFTokenHeader)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// Handler returns an http.Handler for the bridge
func (b *Bridge) Handler() http.Handler {
	return NewHTTPHandler(b)
}
