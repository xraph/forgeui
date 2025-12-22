// Package toast provides a toast notification plugin for ForgeUI.
//
// The toast plugin uses Alpine.js stores for state management and provides
// a clean API for showing temporary notification messages.
//
// # Basic Usage
//
//	registry := plugin.NewRegistry()
//	registry.Use(toast.New())
//
// # Features
//
//   - Multiple variants (info, success, warning, error)
//   - Auto-dismiss with configurable timeout
//   - Queue management (max 5 visible)
//   - Position options (top-right, bottom-center, etc.)
//   - Custom styling per variant
package toast

import (
	"context"

	"github.com/xraph/forgeui/plugin"
)

// Toast plugin implements Alpine and Component plugins.
type Toast struct {
	*plugin.PluginBase

	config Config
}

// Config holds toast plugin configuration.
type Config struct {
	// Position determines where toasts appear
	Position string

	// MaxVisible is the maximum number of toasts to show at once
	MaxVisible int

	// DefaultTimeout is the default auto-dismiss time in milliseconds
	DefaultTimeout int
}

// DefaultConfig returns the default toast configuration.
func DefaultConfig() Config {
	return Config{
		Position:       "top-right",
		MaxVisible:     5,
		DefaultTimeout: 5000,
	}
}

// New creates a new Toast plugin with default config.
func New(opts ...Option) *Toast {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(&config)
	}

	return &Toast{
		PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
			Name:        "toast",
			Version:     "1.0.0",
			Description: "Toast notification system with Alpine.js",
			Author:      "ForgeUI",
			License:     "MIT",
		}),
		config: config,
	}
}

// Option configures the toast plugin.
type Option func(*Config)

// WithPosition sets the toast position.
func WithPosition(position string) Option {
	return func(c *Config) {
		c.Position = position
	}
}

// WithMaxVisible sets the maximum visible toasts.
func WithMaxVisible(max int) Option {
	return func(c *Config) {
		c.MaxVisible = max
	}
}

// WithDefaultTimeout sets the default timeout in milliseconds.
func WithDefaultTimeout(ms int) Option {
	return func(c *Config) {
		c.DefaultTimeout = ms
	}
}

// Init initializes the toast plugin.
func (t *Toast) Init(ctx context.Context, registry *plugin.Registry) error {
	return nil
}

// Shutdown cleanly shuts down the plugin.
func (t *Toast) Shutdown(ctx context.Context) error {
	return nil
}

// Scripts returns scripts needed for the toast plugin.
func (t *Toast) Scripts() []plugin.Script {
	return nil
}

// Directives returns custom Alpine directives.
func (t *Toast) Directives() []plugin.AlpineDirective {
	return nil
}

// Stores returns Alpine stores for toast management.
func (t *Toast) Stores() []plugin.AlpineStore {
	return []plugin.AlpineStore{
		{
			Name: "toasts",
			InitialState: map[string]any{
				"items":          []any{},
				"position":       t.config.Position,
				"maxVisible":     t.config.MaxVisible,
				"defaultTimeout": t.config.DefaultTimeout,
			},
			Methods: `
				show(message, variant = 'info', timeout = null) {
					const id = Date.now() + Math.random();
					const toast = {
						id,
						message,
						variant,
						visible: true
					};
					
					this.items.push(toast);
					
					// Limit to max visible
					if (this.items.length > this.maxVisible) {
						this.items.shift();
					}
					
					// Auto dismiss
					const dismissTimeout = timeout !== null ? timeout : this.defaultTimeout;
					if (dismissTimeout > 0) {
						setTimeout(() => this.dismiss(id), dismissTimeout);
					}
					
					return id;
				},
				
				info(message, timeout = null) {
					return this.show(message, 'info', timeout);
				},
				
				success(message, timeout = null) {
					return this.show(message, 'success', timeout);
				},
				
				warning(message, timeout = null) {
					return this.show(message, 'warning', timeout);
				},
				
				error(message, timeout = null) {
					return this.show(message, 'error', timeout);
				},
				
				dismiss(id) {
					const index = this.items.findIndex(t => t.id === id);
					if (index !== -1) {
						this.items.splice(index, 1);
					}
				},
				
				clear() {
					this.items = [];
				}
			`,
		},
	}
}

// Magics returns custom magic properties.
func (t *Toast) Magics() []plugin.AlpineMagic {
	return nil
}

// AlpineComponents returns Alpine.data components.
func (t *Toast) AlpineComponents() []plugin.AlpineComponent {
	return nil
}
