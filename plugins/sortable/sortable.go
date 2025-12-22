// Package sortable provides drag-and-drop list reordering using SortableJS.
//
// The sortable plugin adds an x-sortable directive that makes lists draggable
// with smooth animations and callbacks for server synchronization.
//
// # Basic Usage
//
//	registry := plugin.NewRegistry()
//	registry.Use(sortable.New())
//
// # Features
//
//   - Drag-and-drop list reordering
//   - Handle selector configuration
//   - Animation options
//   - Callback hooks (onStart, onEnd, onUpdate)
//   - Server sync via HTMX
package sortable

import (
	"context"

	"github.com/xraph/forgeui/plugin"
)

// Sortable plugin implements Alpine plugin.
type Sortable struct {
	*plugin.PluginBase

	version string
}

// New creates a new Sortable plugin.
func New() *Sortable {
	return &Sortable{
		PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
			Name:        "sortable",
			Version:     "1.0.0",
			Description: "Drag-and-drop list reordering with SortableJS",
			Author:      "ForgeUI",
			License:     "MIT",
		}),
		version: "1.15.2",
	}
}

// Init initializes the sortable plugin.
func (s *Sortable) Init(ctx context.Context, registry *plugin.Registry) error {
	return nil
}

// Shutdown cleanly shuts down the plugin.
func (s *Sortable) Shutdown(ctx context.Context) error {
	return nil
}

// Scripts returns SortableJS library and directive.
func (s *Sortable) Scripts() []plugin.Script {
	return []plugin.Script{
		{
			Name:     "sortablejs",
			URL:      "https://cdn.jsdelivr.net/npm/sortablejs@" + s.version + "/Sortable.min.js",
			Priority: 100,
			Defer:    false,
		},
	}
}

// Directives returns the x-sortable directive.
func (s *Sortable) Directives() []plugin.AlpineDirective {
	return []plugin.AlpineDirective{
		{
			Name: "sortable",
			Definition: `
Alpine.directive('sortable', (el, { expression }, { evaluate, effect }) => {
	let sortable = null;
	
	// Parse options from expression
	const getOptions = () => {
		const options = expression ? evaluate(expression) : {};
		return {
			animation: options.animation || 150,
			handle: options.handle || null,
			draggable: options.draggable || null,
			ghostClass: options.ghostClass || 'sortable-ghost',
			chosenClass: options.chosenClass || 'sortable-chosen',
			dragClass: options.dragClass || 'sortable-drag',
			disabled: options.disabled || false,
			
			onStart: (evt) => {
				if (options.onStart) {
					evaluate(options.onStart + '($event)', { $event: evt });
				}
			},
			
			onEnd: (evt) => {
				if (options.onEnd) {
					evaluate(options.onEnd + '($event)', { $event: evt });
				}
			},
			
			onUpdate: (evt) => {
				// Get all items in new order
				const items = Array.from(el.children).map(child => {
					return child.dataset.id || child.dataset.value || child.textContent.trim();
				});
				
				// Dispatch custom event with new order
				el.dispatchEvent(new CustomEvent('sortable-update', {
					detail: { order: items, event: evt },
					bubbles: true
				}));
				
				if (options.onUpdate) {
					evaluate(options.onUpdate + '($event)', { $event: evt });
				}
				
				// If HTMX endpoint specified, sync to server
				if (options.syncUrl) {
					fetch(options.syncUrl, {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ order: items })
					});
				}
			}
		};
	};
	
	// Initialize Sortable
	effect(() => {
		if (sortable) {
			sortable.destroy();
		}
		sortable = Sortable.create(el, getOptions());
	});
	
	// Cleanup
	return () => {
		if (sortable) {
			sortable.destroy();
		}
	};
});
			`,
		},
	}
}

// Stores returns Alpine stores.
func (s *Sortable) Stores() []plugin.AlpineStore {
	return nil
}

// Magics returns custom magic properties.
func (s *Sortable) Magics() []plugin.AlpineMagic {
	return nil
}

// AlpineComponents returns Alpine.data components.
func (s *Sortable) AlpineComponents() []plugin.AlpineComponent {
	return nil
}
