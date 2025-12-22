package plugin

import (
	"errors"
	"fmt"
	"maps"
	"sync"
)

// Registry manages plugin registration and lifecycle.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
	order   []string // initialization order
	hooks   *HookManager

	// Type-specific registries (Phase 12)
	components map[string]ComponentPlugin
	alpine     map[string]AlpinePlugin
	themes     map[string]ThemePlugin
	middleware []MiddlewarePlugin // Sorted by priority
}

// NewRegistry creates a new plugin registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins:    make(map[string]Plugin),
		hooks:      NewHookManager(),
		components: make(map[string]ComponentPlugin),
		alpine:     make(map[string]AlpinePlugin),
		themes:     make(map[string]ThemePlugin),
		middleware: []MiddlewarePlugin{},
	}
}

// Register adds a plugin to the registry.
// Automatically detects plugin type and stores in appropriate registry.
func (r *Registry) Register(p Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := p.Name()
	if name == "" {
		return errors.New("plugin name cannot be empty")
	}

	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	// Store in main registry
	r.plugins[name] = p

	// Type detection: store in specialized registries
	if cp, ok := p.(ComponentPlugin); ok {
		r.components[name] = cp
	}

	if ap, ok := p.(AlpinePlugin); ok {
		r.alpine[name] = ap
	}

	if tp, ok := p.(ThemePlugin); ok {
		r.themes[name] = tp
	}

	if mp, ok := p.(MiddlewarePlugin); ok {
		r.middleware = append(r.middleware, mp)
		// Sort middleware by priority after adding
		r.sortMiddleware()
	}

	return nil
}

// Get retrieves a plugin by name.
func (r *Registry) Get(name string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.plugins[name]

	return p, ok
}

// All returns all registered plugins.
func (r *Registry) All() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		result = append(result, p)
	}

	return result
}

// Use is a convenience method for chaining registrations.
func (r *Registry) Use(plugins ...Plugin) *Registry {
	for _, p := range plugins {
		_ = r.Register(p)
	}

	return r
}

// Unregister removes a plugin.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	delete(r.plugins, name)

	return nil
}

// Has checks if a plugin is registered.
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.plugins[name]

	return ok
}

// Count returns the number of registered plugins.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.plugins)
}

// Hooks returns the hook manager for this registry.
func (r *Registry) Hooks() *HookManager {
	return r.hooks
}

// ResolveDependencies checks all dependencies are satisfied.
func (r *Registry) ResolveDependencies() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, p := range r.plugins {
		for _, dep := range p.Dependencies() {
			target, ok := r.plugins[dep.Name]
			if !ok {
				if dep.Optional {
					continue
				}

				return fmt.Errorf("plugin %s requires %s which is not registered",
					name, dep.Name)
			}

			if !dep.Satisfies(target.Version()) {
				return fmt.Errorf("plugin %s requires %s %s but %s is registered",
					name, dep.Name, dep.Version, target.Version())
			}
		}
	}

	return nil
}

// sortMiddleware sorts middleware plugins by priority (lower = first).
// Must be called with lock held.
func (r *Registry) sortMiddleware() {
	if len(r.middleware) <= 1 {
		return
	}

	// Simple bubble sort (fine for small plugin counts)
	for i := range len(r.middleware) - 1 {
		for j := range len(r.middleware) - i - 1 {
			if r.middleware[j].Priority() > r.middleware[j+1].Priority() {
				r.middleware[j], r.middleware[j+1] = r.middleware[j+1], r.middleware[j]
			}
		}
	}
}

// GetComponentPlugin retrieves a component plugin by name.
func (r *Registry) GetComponentPlugin(name string) (ComponentPlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.components[name]

	return p, ok
}

// GetAlpinePlugin retrieves an Alpine plugin by name.
func (r *Registry) GetAlpinePlugin(name string) (AlpinePlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.alpine[name]

	return p, ok
}

// GetThemePlugin retrieves a theme plugin by name.
func (r *Registry) GetThemePlugin(name string) (ThemePlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.themes[name]

	return p, ok
}

// CollectScripts collects all scripts from Alpine plugins.
// Scripts are returned in priority order (lower priority first).
func (r *Registry) CollectScripts() []Script {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var scripts []Script
	for _, ap := range r.alpine {
		scripts = append(scripts, ap.Scripts()...)
	}

	// Sort by priority
	if len(scripts) > 1 {
		for i := range len(scripts) - 1 {
			for j := range len(scripts) - i - 1 {
				pi := scripts[j].Priority
				pj := scripts[j+1].Priority

				if pi == 0 {
					pi = 50
				}

				if pj == 0 {
					pj = 50
				}

				if pi > pj {
					scripts[j], scripts[j+1] = scripts[j+1], scripts[j]
				}
			}
		}
	}

	return scripts
}

// CollectDirectives collects all Alpine directives from Alpine plugins.
func (r *Registry) CollectDirectives() []AlpineDirective {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var directives []AlpineDirective
	for _, ap := range r.alpine {
		directives = append(directives, ap.Directives()...)
	}

	return directives
}

// CollectStores collects all Alpine stores from Alpine plugins.
func (r *Registry) CollectStores() []AlpineStore {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var stores []AlpineStore
	for _, ap := range r.alpine {
		stores = append(stores, ap.Stores()...)
	}

	return stores
}

// CollectMagics collects all Alpine magic properties from Alpine plugins.
func (r *Registry) CollectMagics() []AlpineMagic {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var magics []AlpineMagic
	for _, ap := range r.alpine {
		magics = append(magics, ap.Magics()...)
	}

	return magics
}

// CollectAlpineComponents collects all Alpine.data components from Alpine plugins.
func (r *Registry) CollectAlpineComponents() []AlpineComponent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var components []AlpineComponent
	for _, ap := range r.alpine {
		components = append(components, ap.AlpineComponents()...)
	}

	return components
}

// CollectComponents collects all component constructors from component plugins.
// Returns a map of component name to constructor.
func (r *Registry) CollectComponents() map[string]ComponentConstructor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]ComponentConstructor)

	for _, cp := range r.components {
		maps.Copy(result, cp.Components())
	}

	return result
}

// CollectMiddleware returns all middleware plugins in priority order.
// Middleware is already sorted by priority in the registry.
func (r *Registry) CollectMiddleware() []MiddlewarePlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]MiddlewarePlugin, len(r.middleware))
	copy(result, r.middleware)

	return result
}
