package plugin

import "net/http"

// MiddlewarePlugin provides HTTP middleware for the application.
//
// Middleware plugins can intercept and modify HTTP requests/responses.
// They are executed in priority order (lower priority = executes first).
//
// Example:
//
//	type HTMXPlugin struct {
//	    *PluginBase
//	}
//
//	func (p *HTMXPlugin) Middleware() func(http.Handler) http.Handler {
//	    return func(next http.Handler) http.Handler {
//	        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	            // Detect HTMX requests
//	            if r.Header.Get("HX-Request") == "true" {
//	                ctx := context.WithValue(r.Context(), htmxRequestKey, true)
//	                r = r.WithContext(ctx)
//	            }
//	            next.ServeHTTP(w, r)
//	        })
//	    }
//	}
//
//	func (p *HTMXPlugin) Priority() int {
//	    return 10 // Execute early
//	}
type MiddlewarePlugin interface {
	Plugin

	// Middleware returns the HTTP middleware function.
	// The middleware should call next.ServeHTTP(w, r) to continue the chain.
	Middleware() func(http.Handler) http.Handler

	// Priority determines execution order.
	// Lower values execute first (e.g., 1 = first, 100 = last).
	// Default priority if not specified: 50.
	Priority() int
}

// MiddlewarePluginBase provides a base implementation for middleware plugins.
type MiddlewarePluginBase struct {
	*PluginBase

	middleware func(http.Handler) http.Handler
	priority   int
}

// NewMiddlewarePluginBase creates a new MiddlewarePluginBase.
func NewMiddlewarePluginBase(
	info PluginInfo,
	middleware func(http.Handler) http.Handler,
	priority int,
) *MiddlewarePluginBase {
	return &MiddlewarePluginBase{
		PluginBase: NewPluginBase(info),
		middleware: middleware,
		priority:   priority,
	}
}

// Middleware returns the middleware function.
func (m *MiddlewarePluginBase) Middleware() func(http.Handler) http.Handler {
	return m.middleware
}

// Priority returns the execution priority.
func (m *MiddlewarePluginBase) Priority() int {
	if m.priority == 0 {
		return 50 // Default priority
	}

	return m.priority
}
