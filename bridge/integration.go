package bridge

import (
	"context"
	"net/http"
)

// Integration provides ForgeUI integration helpers
type Integration struct {
	bridge *Bridge
}

// NewIntegration creates a new ForgeUI integration
func NewIntegration(bridge *Bridge) *Integration {
	return &Integration{bridge: bridge}
}

// RegisterHTTPRoutes registers bridge routes with standard http.ServeMux
func (i *Integration) RegisterHTTPRoutes(mux *http.ServeMux) {
	// Main RPC endpoint
	mux.Handle("/api/bridge/call", i.bridge.Handler())

	// WebSocket endpoint
	mux.Handle("/api/bridge/ws", NewWSHandler(i.bridge))

	// SSE streaming endpoint
	mux.Handle("/api/bridge/stream/", i.bridge.StreamHandler())

	// Introspection endpoint
	mux.Handle("/api/bridge/functions", i.bridge.IntrospectionHandler())
}

// BridgeMiddleware creates a standard http.Handler middleware
func BridgeMiddleware(bridge *Bridge) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add bridge context
			bridgeCtx := NewContext(r)

			// Store in request context
			ctx := r.Context()
			ctx = WithBridgeContext(ctx, bridgeCtx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// Context keys
type contextKey string

const bridgeContextKey contextKey = "bridge_context"

// WithBridgeContext adds a bridge context to a context.Context
func WithBridgeContext(ctx context.Context, bridgeCtx Context) context.Context {
	return context.WithValue(ctx, bridgeContextKey, bridgeCtx)
}

// GetBridgeContext retrieves a bridge context from context.Context
func GetBridgeContext(ctx context.Context) (Context, bool) {
	bridgeCtx, ok := ctx.Value(bridgeContextKey).(Context)
	return bridgeCtx, ok
}
