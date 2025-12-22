package bridge

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Middleware wraps an http.Handler with additional functionality
type Middleware func(http.Handler) http.Handler

// Chain combines multiple middleware
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			log.Printf(
				"%s %s %d %s",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				time.Since(start),
			)
		})
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Panic recovered: %v", err)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					resp := Response{
						JSONRPC: "2.0",
						Error:   ErrInternal,
					}

					_ = json.NewEncoder(w).Encode(resp)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				// Check if origin is allowed
				allowed := false

				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == "*" || allowedOrigin == origin {
						allowed = true
						break
					}
				}

				if allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Set("Access-Control-Max-Age", "86400")
				}
			}

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adds a unique request ID
func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}

			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter

	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
