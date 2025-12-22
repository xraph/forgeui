package bridge

import "encoding/json"

// Request represents a JSON-RPC 2.0 request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any         `json:"id,omitempty"`
	Result  any         `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// BatchRequest represents multiple requests in a single call
type BatchRequest []Request

// BatchResponse represents multiple responses
type BatchResponse []Response

// Error represents a JSON-RPC 2.0 error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Standard JSON-RPC 2.0 error codes
const (
	ErrCodeParseError     = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternal       = -32603
)

// Custom error codes
const (
	ErrCodeUnauthorized = -32001
	ErrCodeRateLimit    = -32002
	ErrCodeTimeout      = -32003
	ErrCodeBadRequest   = -32004
	ErrCodeForbidden    = -32005
)

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new bridge error
func NewError(code int, message string, data ...any) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}
	if len(data) > 0 {
		err.Data = data[0]
	}
	return err
}

// Event represents a server-initiated event (for WebSocket/SSE)
type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// StreamChunk represents a chunk of streaming data
type StreamChunk struct {
	Data  any  `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
	Done  bool `json:"done"`
}

// Common error constructors
var (
	ErrParseError     = NewError(ErrCodeParseError, "Parse error")
	ErrInvalidRequest = NewError(ErrCodeInvalidRequest, "Invalid request")
	ErrMethodNotFound = NewError(ErrCodeMethodNotFound, "Method not found")
	ErrInvalidParams  = NewError(ErrCodeInvalidParams, "Invalid params")
	ErrInternal       = NewError(ErrCodeInternal, "Internal error")
	ErrUnauthorized   = NewError(ErrCodeUnauthorized, "Unauthorized")
	ErrRateLimit      = NewError(ErrCodeRateLimit, "Rate limit exceeded")
	ErrTimeout        = NewError(ErrCodeTimeout, "Request timeout")
	ErrBadRequest     = NewError(ErrCodeBadRequest, "Bad request")
	ErrForbidden      = NewError(ErrCodeForbidden, "Forbidden")
)

