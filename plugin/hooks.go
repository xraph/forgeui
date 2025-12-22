package plugin

import (
	"context"
	"sync"

	g "github.com/maragudk/gomponents"
)

// Hook names for plugin lifecycle and rendering.
//
// Available hooks:
//
// Lifecycle hooks:
//   - before_init: Before plugin initialization
//   - after_init: After plugin initialization
//   - before_shutdown: Before plugin shutdown
//   - after_shutdown: After plugin shutdown
//
// Render hooks:
//   - before_render: Before page render
//   - after_render: After page render
//   - before_head: Before <head> content
//   - after_head: After <head> content
//   - before_body: Before <body> content
//   - after_body: After </body> (scripts area)
//   - before_scripts: Before script tags
//   - after_scripts: After script tags
const (
	HookBeforeInit     = "before_init"
	HookAfterInit      = "after_init"
	HookBeforeShutdown = "before_shutdown"
	HookAfterShutdown  = "after_shutdown"
	HookBeforeRender   = "before_render"
	HookAfterRender    = "after_render"
	HookBeforeHead     = "before_head"
	HookAfterHead      = "after_head"
	HookBeforeBody     = "before_body"
	HookAfterBody      = "after_body"
	HookBeforeScripts  = "before_scripts"
	HookAfterScripts   = "after_scripts"
)

// HookFunc is a hook handler function.
type HookFunc func(ctx *HookContext) error

// HookContext provides context to hook handlers.
type HookContext struct {
	Context context.Context
	Data    map[string]any
	Nodes   []g.Node // For render hooks
}

// HookManager manages hook registration and execution.
type HookManager struct {
	mu    sync.RWMutex
	hooks map[string][]HookFunc
}

// NewHookManager creates a new hook manager.
func NewHookManager() *HookManager {
	return &HookManager{
		hooks: make(map[string][]HookFunc),
	}
}

// On registers a hook handler.
func (m *HookManager) On(hook string, fn HookFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks[hook] = append(m.hooks[hook], fn)
}

// Off removes all handlers for a hook.
func (m *HookManager) Off(hook string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.hooks, hook)
}

// Trigger executes all handlers for a hook.
// Returns the first error encountered.
func (m *HookManager) Trigger(hook string, ctx *HookContext) error {
	m.mu.RLock()
	handlers := make([]HookFunc, len(m.hooks[hook]))
	copy(handlers, m.hooks[hook])
	m.mu.RUnlock()

	for _, fn := range handlers {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Has checks if a hook has any handlers registered.
func (m *HookManager) Has(hook string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.hooks[hook]) > 0
}

// Count returns the number of handlers for a hook.
func (m *HookManager) Count(hook string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.hooks[hook])
}

