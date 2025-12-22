package plugin

// AlpinePlugin extends Alpine.js functionality.
//
// An AlpinePlugin can provide:
//   - External scripts/libraries to load
//   - Custom Alpine directives (e.g., x-sortable)
//   - Global Alpine stores for state management
//   - Magic properties (e.g., $myPlugin)
//   - Alpine.data components
//
// Example:
//
//	type SortablePlugin struct {
//	    *PluginBase
//	}
//
//	func (p *SortablePlugin) Scripts() []Script {
//	    return []Script{
//	        {Name: "sortablejs", URL: "https://cdn.jsdelivr.net/.../Sortable.min.js"},
//	    }
//	}
//
//	func (p *SortablePlugin) Directives() []AlpineDirective {
//	    return []AlpineDirective{
//	        {Name: "sortable", Definition: "..."},
//	    }
//	}
type AlpinePlugin interface {
	Plugin

	// Scripts returns external scripts to load (libraries, dependencies).
	Scripts() []Script

	// Directives returns custom Alpine directives.
	Directives() []AlpineDirective

	// Stores returns Alpine stores to register globally.
	Stores() []AlpineStore

	// Magics returns custom magic properties.
	Magics() []AlpineMagic

	// AlpineComponents returns Alpine.data components.
	AlpineComponents() []AlpineComponent
}

// Script represents an external script to load.
//
// Scripts can be loaded from URLs or inlined. Priority determines load order
// (lower values load first).
//
// Example:
//
//	Script{
//	    Name:     "chartjs",
//	    URL:      "https://cdn.jsdelivr.net/npm/chart.js",
//	    Priority: 10,
//	    Defer:    true,
//	}
type Script struct {
	// Name is a unique identifier for the script
	Name string

	// URL is the script source (https://... or /path/to/script.js)
	URL string

	// Inline is JavaScript code to execute (alternative to URL)
	Inline string

	// Defer adds the defer attribute (loads after HTML parsing)
	Defer bool

	// Async adds the async attribute (loads asynchronously)
	Async bool

	// Priority determines load order (lower = loads first)
	// Default: 50. Use 1-10 for critical dependencies.
	Priority int

	// Module adds type="module" attribute
	Module bool

	// Integrity is the SRI hash for the script
	Integrity string

	// Crossorigin sets the crossorigin attribute
	Crossorigin string
}

// AlpineDirective represents a custom Alpine directive.
//
// Directives extend Alpine's x- attribute system. The definition is
// JavaScript code that Alpine will execute.
//
// Example:
//
//	AlpineDirective{
//	    Name: "sortable",
//	    Definition: `
//	        (el, { expression, modifiers }, { evaluate }) => {
//	            let options = expression ? evaluate(expression) : {};
//	            new Sortable(el, options);
//	        }
//	    `,
//	}
type AlpineDirective struct {
	// Name is the directive name (used as x-name)
	Name string

	// Definition is the JavaScript function implementing the directive
	// Signature: (el, { expression, modifiers }, { evaluate, effect, cleanup }) => {}
	Definition string
}

// AlpineStore represents a global Alpine store.
//
// Stores provide reactive state accessible via $store.storeName.
// The initial state and methods are combined into a single object.
//
// Example:
//
//	AlpineStore{
//	    Name: "notifications",
//	    InitialState: map[string]any{
//	        "items": []any{},
//	        "count": 0,
//	    },
//	    Methods: `
//	        add(item) {
//	            this.items.push(item);
//	            this.count++;
//	        },
//	        remove(id) {
//	            this.items = this.items.filter(i => i.id !== id);
//	            this.count--;
//	        }
//	    `,
//	}
type AlpineStore struct {
	// Name is the store identifier (accessed via $store.name)
	Name string

	// InitialState is the initial state as a map
	InitialState map[string]any

	// Methods is JavaScript code defining store methods
	// These are merged with the initial state
	Methods string
}

// AlpineMagic represents a custom magic property.
//
// Magic properties are accessed via $name and can return any value.
//
// Example:
//
//	AlpineMagic{
//	    Name: "clipboard",
//	    Definition: `
//	        (el) => {
//	            return {
//	                copy(text) {
//	                    navigator.clipboard.writeText(text);
//	                }
//	            }
//	        }
//	    `,
//	}
type AlpineMagic struct {
	// Name is the magic property name (accessed via $name)
	Name string

	// Definition is the JavaScript function that returns the magic value
	// Signature: (el, { Alpine }) => value
	Definition string
}

// AlpineComponent represents an Alpine.data component.
//
// Components encapsulate reactive state and methods that can be reused
// across the application.
//
// Example:
//
//	AlpineComponent{
//	    Name: "dropdown",
//	    Definition: `
//	        () => ({
//	            open: false,
//	            toggle() {
//	                this.open = !this.open;
//	            },
//	            close() {
//	                this.open = false;
//	            }
//	        })
//	    `,
//	}
type AlpineComponent struct {
	// Name is the component name (used with x-data="name")
	Name string

	// Definition is the JavaScript function returning the component state
	// Signature: (...args) => ({ state, methods })
	Definition string
}

// AlpinePluginBase provides default implementations for AlpinePlugin.
// Embed this to implement only the methods you need.
type AlpinePluginBase struct {
	*PluginBase
}

// NewAlpinePluginBase creates a new AlpinePluginBase.
func NewAlpinePluginBase(info PluginInfo) *AlpinePluginBase {
	return &AlpinePluginBase{
		PluginBase: NewPluginBase(info),
	}
}

// Scripts returns an empty slice by default.
func (a *AlpinePluginBase) Scripts() []Script {
	return nil
}

// Directives returns an empty slice by default.
func (a *AlpinePluginBase) Directives() []AlpineDirective {
	return nil
}

// Stores returns an empty slice by default.
func (a *AlpinePluginBase) Stores() []AlpineStore {
	return nil
}

// Magics returns an empty slice by default.
func (a *AlpinePluginBase) Magics() []AlpineMagic {
	return nil
}

// AlpineComponents returns an empty slice by default.
func (a *AlpinePluginBase) AlpineComponents() []AlpineComponent {
	return nil
}
