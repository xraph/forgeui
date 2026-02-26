package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

// Middleware is a function that wraps a PageHandler.
type Middleware func(PageHandler) PageHandler

// Logger returns middleware that logs requests.
func Logger() Middleware {
	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			start := time.Now()

			// Process request
			comp, err := next(ctx)

			// Log request
			duration := time.Since(start)

			status := 200
			if err != nil {
				status = 500
			}

			log.Printf("[%s] %s %s - %d (%v)",
				ctx.Method(),
				ctx.Path(),
				ctx.ClientIP(),
				status,
				duration,
			)

			return comp, err
		}
	}
}

// Recovery returns middleware that recovers from panics.
func Recovery() Middleware {
	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (comp templ.Component, err error) {
			defer func() {
				if r := recover(); r != nil {
					// Log panic with stack trace
					log.Printf("PANIC: %v\n%s", r, debug.Stack())

					// Set error response
					ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

					comp = templ.Raw(fmt.Sprintf("<div><h1>500 - Internal Server Error</h1><p>The server encountered an unexpected error: %v</p></div>", r))
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()

			return next(ctx)
		}
	}
}

// CORS returns middleware that adds CORS headers.
func CORS(allowOrigin string) Middleware {
	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			// Set CORS headers
			ctx.SetHeader("Access-Control-Allow-Origin", allowOrigin)
			ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if ctx.Method() == http.MethodOptions {
				ctx.ResponseWriter.WriteHeader(http.StatusOK)
				return nil, nil
			}

			return next(ctx)
		}
	}
}

// RequestID returns middleware that adds a unique request ID.
func RequestID() Middleware {
	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			// Generate simple request ID (in production, use UUID or similar)
			requestID := strconv.FormatInt(time.Now().UnixNano(), 10)

			// Store in context
			ctx.Set("request_id", requestID)

			// Add to response header
			ctx.SetHeader("X-Request-ID", requestID)

			return next(ctx)
		}
	}
}

// BasicAuth returns middleware that requires HTTP basic authentication.
func BasicAuth(username, password string) Middleware {
	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			// Get credentials from request
			user, pass, ok := ctx.Request.BasicAuth()

			// Check credentials
			if !ok || user != username || pass != password {
				ctx.SetHeader("WWW-Authenticate", `Basic realm="Restricted"`)
				ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)

				return templ.Raw("401 - Unauthorized"), nil
			}

			return next(ctx)
		}
	}
}

// RequireMethod returns middleware that only allows specific HTTP methods.
func RequireMethod(methods ...string) Middleware {
	methodMap := make(map[string]bool)
	for _, m := range methods {
		methodMap[m] = true
	}

	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			if !methodMap[ctx.Method()] {
				ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)
				return templ.Raw(fmt.Sprintf("405 - Method Not Allowed: %s", ctx.Method())), nil
			}

			return next(ctx)
		}
	}
}

// Timeout returns middleware that enforces a timeout on requests.
func Timeout(duration time.Duration) Middleware {
	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			// Create timeout context
			timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), duration)
			defer cancel()

			// Update context
			ctx = ctx.WithContext(timeoutCtx)

			// Channel for result
			type result struct {
				comp templ.Component
				err  error
			}

			done := make(chan result, 1)

			// Execute handler in goroutine
			go func() {
				comp, err := next(ctx)
				done <- result{comp, err}
			}()

			// Wait for result or timeout
			select {
			case res := <-done:
				return res.comp, res.err
			case <-timeoutCtx.Done():
				ctx.ResponseWriter.WriteHeader(http.StatusRequestTimeout)
				return templ.Raw("408 - Request Timeout"), fmt.Errorf("request timeout after %v", duration)
			}
		}
	}
}

// SetHeader returns middleware that sets a response header.
func SetHeader(key, value string) Middleware {
	return func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			ctx.SetHeader(key, value)
			return next(ctx)
		}
	}
}

// Chain chains multiple middleware into one.
func Chain(middleware ...Middleware) Middleware {
	return func(next PageHandler) PageHandler {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}

		return next
	}
}
