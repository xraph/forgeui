package bridge

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/a-h/templ"
)

// HTMXHandler serves bridge functions as HTTP endpoints returning HTML or JSON.
// Functions are mapped to /prefix/{functionName} where the prefix defaults to "/api/bridge/fn/".
// Supports GET and POST with parameters from query strings, form data, or JSON bodies.
type HTMXHandler struct {
	bridge   *Bridge
	security *Security
	prefix   string
}

// NewHTMXHandler creates a new HTMX handler with the given URL prefix
func NewHTMXHandler(bridge *Bridge, prefix string) *HTMXHandler {
	if prefix == "" {
		prefix = "/api/bridge/fn/"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	return &HTMXHandler{
		bridge:   bridge,
		security: NewSecurity(bridge.config),
		prefix:   prefix,
	}
}

// ServeHTTP handles HTTP requests for bridge functions
func (h *HTMXHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract function name from URL path
	funcName := strings.TrimPrefix(r.URL.Path, h.prefix)
	funcName = strings.TrimSuffix(funcName, "/")

	if funcName == "" {
		http.Error(w, "Function name required", http.StatusBadRequest)
		return
	}

	// Handle CORS
	if h.bridge.config.EnableCORS {
		h.handleCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Get function
	fn, err := h.bridge.GetFunction(funcName)
	if err != nil {
		http.Error(w, "Function not found", http.StatusNotFound)
		return
	}

	// Check allowed methods
	if !h.isMethodAllowed(fn, r.Method) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// CSRF check (skip for GET/HEAD - safe methods)
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		if csrfErr := h.security.CheckCSRF(r); csrfErr != nil {
			http.Error(w, "CSRF validation failed", http.StatusForbidden)
			return
		}
	}

	// Create bridge context
	ctx := NewContext(r)

	// Check authentication
	if authErr := h.security.CheckAuth(ctx, fn); authErr != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check rate limit
	rateLimitKey := getRateLimitKey(ctx)
	if rlErr := h.security.CheckRateLimit(rateLimitKey, fn); rlErr != nil {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Parse parameters from HTTP request
	var paramValue reflect.Value
	if fn.HasInput {
		var parseErr error
		paramValue, parseErr = parseHTTPParams(r, fn.InputType)
		if parseErr != nil {
			var bridgeErr *Error
			if errors.As(parseErr, &bridgeErr) {
				http.Error(w, bridgeErr.Message, http.StatusBadRequest)
			} else {
				http.Error(w, "Invalid parameters", http.StatusBadRequest)
			}
			return
		}
	}

	// Execute function directly (no JSON round-trip)
	result := h.bridge.executeDirect(ctx, fn, paramValue)

	if result.Error != nil {
		statusCode := h.errorToStatus(result.Error)
		http.Error(w, result.Error.Message, statusCode)
		return
	}

	// Set HTMX response headers
	h.setHTMXHeaders(w, fn)

	// Render response
	h.renderResponse(w, r, fn, result.Result)
}

// isMethodAllowed checks if the HTTP method is allowed for the function.
// If AllowedMethods is explicitly set, use that. Otherwise auto-detect:
// functions without input default to GET+POST, functions with input default to POST only.
func (h *HTMXHandler) isMethodAllowed(fn *Function, method string) bool {
	allowed := fn.AllowedMethods
	if len(allowed) == 0 {
		// Auto-detect based on signature
		if fn.HasInput {
			allowed = []string{http.MethodPost}
		} else {
			allowed = []string{http.MethodGet, http.MethodPost}
		}
	}

	for _, m := range allowed {
		if strings.EqualFold(m, method) {
			return true
		}
	}

	return false
}

// renderResponse renders the execution result as HTML or JSON
func (h *HTMXHandler) renderResponse(w http.ResponseWriter, r *http.Request, fn *Function, result any) {
	// Case 1: Function returns templ.Component directly
	if fn.ReturnsHTML {
		if comp, ok := result.(templ.Component); ok && comp != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if renderErr := comp.Render(r.Context(), w); renderErr != nil {
				http.Error(w, "Failed to render component", http.StatusInternalServerError)
			}
			return
		}
	}

	// Case 2: Function has a renderer option (data → HTML)
	if fn.Renderer != nil && result != nil {
		comp := fn.Renderer(result)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if renderErr := comp.Render(r.Context(), w); renderErr != nil {
			http.Error(w, "Failed to render component", http.StatusInternalServerError)
		}
		return
	}

	// Case 3: Void result (no output) - return empty 200
	if result == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Case 4: HTMX request but no HTML renderer - wrap JSON in <pre> as fallback
	if isHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprintf(w, "<pre>%s</pre>", jsonBytes)
		return
	}

	// Case 5: Non-HTMX request - return JSON
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// setHTMXHeaders sets HTMX-specific response headers from function configuration
func (h *HTMXHandler) setHTMXHeaders(w http.ResponseWriter, fn *Function) {
	if len(fn.HTMXTriggers) > 0 {
		w.Header().Set("HX-Trigger", strings.Join(fn.HTMXTriggers, ", "))
	}
	if fn.HTMXRedirect != "" {
		w.Header().Set("HX-Redirect", fn.HTMXRedirect)
	}
	if fn.HTMXReswap != "" {
		w.Header().Set("HX-Reswap", fn.HTMXReswap)
	}
	if fn.HTMXRetarget != "" {
		w.Header().Set("HX-Retarget", fn.HTMXRetarget)
	}
}

// errorToStatus maps bridge error codes to HTTP status codes
func (h *HTMXHandler) errorToStatus(err *Error) int {
	switch err.Code {
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeMethodNotFound:
		return http.StatusNotFound
	case ErrCodeInvalidParams, ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

// handleCORS sets CORS headers for the HTMX handler
func (h *HTMXHandler) handleCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	allowed := false
	for _, o := range h.bridge.config.AllowedOrigins {
		if o == "*" || o == origin {
			allowed = true
			break
		}
	}

	if !allowed {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, HX-Request, HX-Trigger, HX-Target, "+h.bridge.config.CSRFTokenHeader)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// isHTMXRequest checks if the request is from HTMX
func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// HTMXHandler returns an http.Handler for HTMX bridge requests
func (b *Bridge) HTMXHandler(prefix ...string) http.Handler {
	p := "/api/bridge/fn/"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return NewHTMXHandler(b, p)
}
