package router

import (
	"context"
	"net/http"
	"strconv"
)

// PageContext wraps the HTTP request and response with additional utilities
type PageContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Params         Params
	values         map[string]interface{}
	LoadedData     any
	Meta           *RouteMeta
	app            interface{} // Reference to EnhancedApp (interface to avoid circular dependency)
}

// Param retrieves a path parameter by key
func (c *PageContext) Param(key string) string {
	if c.Params == nil {
		return ""
	}
	return c.Params[key]
}

// ParamInt retrieves a path parameter as an integer
func (c *PageContext) ParamInt(key string) (int, error) {
	val := c.Param(key)
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// ParamInt64 retrieves a path parameter as an int64
func (c *PageContext) ParamInt64(key string) (int64, error) {
	val := c.Param(key)
	if val == "" {
		return 0, nil
	}
	return strconv.ParseInt(val, 10, 64)
}

// Query retrieves a query parameter by key
func (c *PageContext) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// QueryDefault retrieves a query parameter with a default value
func (c *PageContext) QueryDefault(key, defaultVal string) string {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// QueryInt retrieves a query parameter as an integer
func (c *PageContext) QueryInt(key string) (int, error) {
	val := c.Query(key)
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// QueryInt64 retrieves a query parameter as an int64
func (c *PageContext) QueryInt64(key string) (int64, error) {
	val := c.Query(key)
	if val == "" {
		return 0, nil
	}
	return strconv.ParseInt(val, 10, 64)
}

// QueryBool retrieves a query parameter as a boolean
func (c *PageContext) QueryBool(key string) bool {
	val := c.Query(key)
	return val == "true" || val == "1" || val == "yes"
}

// Header retrieves a request header by key
func (c *PageContext) Header(key string) string {
	return c.Request.Header.Get(key)
}

// SetHeader sets a response header
func (c *PageContext) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// Cookie retrieves a cookie by name
func (c *PageContext) Cookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

// SetCookie sets a cookie in the response
func (c *PageContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

// Set stores a value in the context for the duration of the request
func (c *PageContext) Set(key string, value interface{}) {
	if c.values == nil {
		c.values = make(map[string]interface{})
	}
	c.values[key] = value
}

// Get retrieves a value from the context
func (c *PageContext) Get(key string) (interface{}, bool) {
	if c.values == nil {
		return nil, false
	}
	val, ok := c.values[key]
	return val, ok
}

// GetString retrieves a string value from the context
func (c *PageContext) GetString(key string) string {
	if val, ok := c.Get(key); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt retrieves an integer value from the context
func (c *PageContext) GetInt(key string) int {
	if val, ok := c.Get(key); ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return 0
}

// Context returns the request's context
func (c *PageContext) Context() context.Context {
	return c.Request.Context()
}

// WithContext returns a shallow copy of PageContext with a new context
func (c *PageContext) WithContext(ctx context.Context) *PageContext {
	copy := *c
	copy.Request = c.Request.WithContext(ctx)
	return &copy
}

// Method returns the HTTP method of the request
func (c *PageContext) Method() string {
	return c.Request.Method
}

// Path returns the URL path of the request
func (c *PageContext) Path() string {
	return c.Request.URL.Path
}

// Host returns the host of the request
func (c *PageContext) Host() string {
	return c.Request.Host
}

// IsSecure returns true if the request is over HTTPS
func (c *PageContext) IsSecure() bool {
	return c.Request.TLS != nil
}

// ClientIP returns the client's IP address (best effort)
func (c *PageContext) ClientIP() string {
	// Check X-Forwarded-For header
	if ip := c.Header("X-Forwarded-For"); ip != "" {
		return ip
	}
	// Check X-Real-IP header
	if ip := c.Header("X-Real-IP"); ip != "" {
		return ip
	}
	// Fallback to RemoteAddr
	return c.Request.RemoteAddr
}

// LoaderData returns the data loaded by the route's loader
func (c *PageContext) LoaderData() any {
	return c.LoadedData
}

// GetMeta returns the route's metadata
func (c *PageContext) GetMeta() *RouteMeta {
	return c.Meta
}

// App returns the application instance
// Returns interface{} to avoid circular dependency - cast to *forgeui.EnhancedApp in usage
func (c *PageContext) App() interface{} {
	return c.app
}
