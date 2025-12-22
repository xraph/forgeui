package bridge

import (
	"fmt"
	"sync"
	"time"
)

// Bridge manages function registration and execution
type Bridge struct {
	mu        sync.RWMutex
	functions map[string]*Function
	config    *Config
	hooks     *HookManager
}

// Config holds bridge configuration
type Config struct {
	// Timeout is the default timeout for function execution
	Timeout time.Duration

	// MaxBatchSize is the maximum number of requests in a batch
	MaxBatchSize int

	// EnableCSRF enables CSRF token validation
	EnableCSRF bool

	// EnableCORS enables CORS headers
	EnableCORS bool

	// AllowedOrigins is the list of allowed origins for CORS
	AllowedOrigins []string

	// DefaultRateLimit is the default rate limit (requests per minute)
	DefaultRateLimit int

	// EnableCache enables result caching
	EnableCache bool

	// CSRFTokenHeader is the header name for CSRF token
	CSRFTokenHeader string

	// CSRFCookieName is the cookie name for CSRF token
	CSRFCookieName string
}

// DefaultConfig returns the default bridge configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:          30 * time.Second,
		MaxBatchSize:     10,
		EnableCSRF:       true,
		EnableCORS:       true,
		AllowedOrigins:   []string{"*"},
		DefaultRateLimit: 60,
		EnableCache:      true,
		CSRFTokenHeader:  "X-CSRF-Token",
		CSRFCookieName:   "csrf_token",
	}
}

// ConfigOption configures the Bridge
type ConfigOption func(*Config)

// WithTimeout sets the default timeout
func WithTimeout(d time.Duration) ConfigOption {
	return func(c *Config) {
		c.Timeout = d
	}
}

// WithMaxBatchSize sets the maximum batch size
func WithMaxBatchSize(size int) ConfigOption {
	return func(c *Config) {
		c.MaxBatchSize = size
	}
}

// WithCSRF enables or disables CSRF protection
func WithCSRF(enabled bool) ConfigOption {
	return func(c *Config) {
		c.EnableCSRF = enabled
	}
}

// WithCORS enables or disables CORS
func WithCORS(enabled bool) ConfigOption {
	return func(c *Config) {
		c.EnableCORS = enabled
	}
}

// WithAllowedOrigins sets the allowed origins for CORS
func WithAllowedOrigins(origins ...string) ConfigOption {
	return func(c *Config) {
		c.AllowedOrigins = origins
	}
}

// WithDefaultRateLimit sets the default rate limit
func WithDefaultRateLimit(rpm int) ConfigOption {
	return func(c *Config) {
		c.DefaultRateLimit = rpm
	}
}

// WithCache enables or disables caching
func WithCache(enabled bool) ConfigOption {
	return func(c *Config) {
		c.EnableCache = enabled
	}
}

// New creates a new Bridge with the given configuration
func New(opts ...ConfigOption) *Bridge {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	return &Bridge{
		functions: make(map[string]*Function),
		config:    config,
		hooks:     NewHookManager(),
	}
}

// Register registers a new function
// Expected signature: func(Context, InputType) (OutputType, error)
func (b *Bridge) Register(name string, handler any, opts ...FunctionOption) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check if function already exists
	if _, exists := b.functions[name]; exists {
		return fmt.Errorf("function %s already registered", name)
	}

	// Analyze function signature
	fn, err := analyzeFunction(handler)
	if err != nil {
		return fmt.Errorf("invalid function signature: %w", err)
	}

	fn.Name = name

	// Apply function options
	for _, opt := range opts {
		opt(fn)
	}

	// Store function
	b.functions[name] = fn

	return nil
}

// Unregister removes a function
func (b *Bridge) Unregister(name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.functions[name]; !exists {
		return fmt.Errorf("function %s not found", name)
	}

	delete(b.functions, name)

	return nil
}

// GetFunction retrieves a registered function
func (b *Bridge) GetFunction(name string) (*Function, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	fn, exists := b.functions[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	return fn, nil
}

// ListFunctions returns all registered function names
func (b *Bridge) ListFunctions() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	names := make([]string, 0, len(b.functions))
	for name := range b.functions {
		names = append(names, name)
	}

	return names
}

// GetConfig returns the bridge configuration
func (b *Bridge) GetConfig() *Config {
	return b.config
}

// GetHooks returns the hook manager
func (b *Bridge) GetHooks() *HookManager {
	return b.hooks
}

// FunctionCount returns the number of registered functions
func (b *Bridge) FunctionCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.functions)
}

// HasFunction checks if a function is registered
func (b *Bridge) HasFunction(name string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, exists := b.functions[name]

	return exists
}
