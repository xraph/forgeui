package bridge

import (
	"sync"
)

// HookType identifies different hook points
type HookType string

const (
	// BeforeCall is called before function execution
	BeforeCall HookType = "before_call"

	// AfterCall is called after function execution
	AfterCall HookType = "after_call"

	// OnError is called when an error occurs
	OnError HookType = "on_error"

	// OnSuccess is called when a function succeeds
	OnSuccess HookType = "on_success"
)

// Hook is a function called at specific points in the request lifecycle
type Hook func(ctx Context, data HookData)

// HookData contains information passed to hooks
type HookData struct {
	FunctionName string
	Params       any
	Result       any
	Error        error
	Duration     int64 // Execution time in microseconds
}

// HookManager manages hooks
type HookManager struct {
	mu    sync.RWMutex
	hooks map[HookType][]Hook
}

// NewHookManager creates a new hook manager
func NewHookManager() *HookManager {
	return &HookManager{
		hooks: make(map[HookType][]Hook),
	}
}

// Register adds a hook for the given type
func (hm *HookManager) Register(hookType HookType, hook Hook) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.hooks[hookType] = append(hm.hooks[hookType], hook)
}

// Trigger executes all hooks of the given type
func (hm *HookManager) Trigger(hookType HookType, ctx Context, data HookData) {
	hm.mu.RLock()
	hooks := hm.hooks[hookType]
	hm.mu.RUnlock()

	for _, hook := range hooks {
		// Execute hook in a goroutine to avoid blocking
		// Use a copy of the hook to avoid closure issues
		h := hook

		go func() {
			defer func() {
				// Recover from panics in hooks
				if r := recover(); r != nil {
					// Log panic but don't propagate
					// In production, you'd want to log this properly
					_ = r // Silence unused variable warning
				}
			}()

			h(ctx, data)
		}()
	}
}

// Clear removes all hooks of the given type
func (hm *HookManager) Clear(hookType HookType) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	delete(hm.hooks, hookType)
}

// ClearAll removes all hooks
func (hm *HookManager) ClearAll() {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.hooks = make(map[HookType][]Hook)
}

// Count returns the number of hooks for a given type
func (hm *HookManager) Count(hookType HookType) int {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	return len(hm.hooks[hookType])
}
