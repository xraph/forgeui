package plugin

import (
	"context"
	"fmt"
)

// Initialize initializes all plugins in dependency order.
//
// The initialization process:
// 1. Resolves all dependencies
// 2. Performs topological sort to determine order
// 3. Initializes each plugin in order
// 4. Triggers before_init and after_init hooks
//
// If any plugin fails to initialize, the process stops and returns an error.
func (r *Registry) Initialize(ctx context.Context) error {
	// Resolve dependencies first
	if err := r.ResolveDependencies(); err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}

	// Get sorted order
	sorted, err := r.TopologicalSort()
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}

	// Initialize each plugin
	for _, p := range sorted {
		// Trigger before_init hook
		if err := r.hooks.Trigger(HookBeforeInit, &HookContext{
			Context: ctx,
			Data:    map[string]any{"plugin": p.Name()},
		}); err != nil {
			return fmt.Errorf("before_init hook failed for %s: %w", p.Name(), err)
		}

		if err := p.Init(ctx, r); err != nil {
			return fmt.Errorf("failed to initialize plugin %s: %w", p.Name(), err)
		}

		// Trigger after_init hook
		if err := r.hooks.Trigger(HookAfterInit, &HookContext{
			Context: ctx,
			Data:    map[string]any{"plugin": p.Name()},
		}); err != nil {
			return fmt.Errorf("after_init hook failed for %s: %w", p.Name(), err)
		}
	}

	// Store initialization order for shutdown
	r.mu.Lock()
	r.order = make([]string, len(sorted))
	for i, p := range sorted {
		r.order[i] = p.Name()
	}
	r.mu.Unlock()

	return nil
}

// Shutdown shuts down all plugins in reverse initialization order.
//
// Unlike initialization, shutdown attempts to shut down all plugins even if
// some fail. All errors are collected and returned as a single error.
//
// This ensures that resources are cleaned up as much as possible.
func (r *Registry) Shutdown(ctx context.Context) error {
	r.mu.RLock()
	order := make([]string, len(r.order))
	copy(order, r.order)
	r.mu.RUnlock()

	var errs []error

	// Shutdown in reverse initialization order
	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]

		r.mu.RLock()
		p, ok := r.plugins[name]
		r.mu.RUnlock()

		if !ok {
			continue
		}

		// Trigger before_shutdown hook
		if err := r.hooks.Trigger(HookBeforeShutdown, &HookContext{
			Context: ctx,
			Data:    map[string]any{"plugin": name},
		}); err != nil {
			errs = append(errs, fmt.Errorf("before_shutdown hook failed for %s: %w", name, err))
		}

		if err := p.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("plugin %s: %w", name, err))
		}

		// Trigger after_shutdown hook
		if err := r.hooks.Trigger(HookAfterShutdown, &HookContext{
			Context: ctx,
			Data:    map[string]any{"plugin": name},
		}); err != nil {
			errs = append(errs, fmt.Errorf("after_shutdown hook failed for %s: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

