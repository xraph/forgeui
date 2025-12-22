// Package primitives provides the Provider pattern for Alpine.js integration.
//
// The Provider pattern enables React-like context API for Alpine.js components,
// allowing parent components to share state and methods with deeply nested children
// without prop drilling.
package primitives

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ProviderProps defines the configuration for a Provider component.
type ProviderProps struct {
	// Name is the unique identifier for this provider (e.g., "sidebar", "theme")
	Name string

	// State is the initial state as a map of key-value pairs
	State map[string]any

	// Methods is raw JavaScript code defining provider methods
	// Methods will be merged with the state object and have access to 'this'
	// Example: "toggle() { this.open = !this.open; }"
	Methods string

	// OnInit is JavaScript code that runs when the provider initializes
	// Useful for setup logic, event listeners, etc.
	OnInit string

	// Class adds custom CSS classes to the provider wrapper
	Class string

	// Attrs adds custom HTML attributes to the provider wrapper
	Attrs []g.Node

	// Children are the child nodes that will have access to this provider
	Children []g.Node

	// Debug enables console logging for state changes
	Debug bool

	// Hooks define lifecycle hooks for the provider
	Hooks map[string]string
}

// ProviderOption is a functional option for configuring a Provider.
type ProviderOption func(*ProviderProps)

// WithProviderName sets the provider's unique identifier.
func WithProviderName(name string) ProviderOption {
	return func(p *ProviderProps) { p.Name = name }
}

// WithProviderState sets the provider's initial state.
func WithProviderState(state map[string]any) ProviderOption {
	return func(p *ProviderProps) { p.State = state }
}

// WithProviderMethods adds JavaScript methods to the provider.
// Methods have access to 'this' (the provider state) and Alpine magic properties.
//
// Example:
//
//	WithProviderMethods(`
//	    toggle() {
//	        this.open = !this.open;
//	        this.$dispatch('provider:toggled', { open: this.open });
//	    }
//	`)
func WithProviderMethods(methods string) ProviderOption {
	return func(p *ProviderProps) { p.Methods = methods }
}

// WithProviderInit sets initialization code for the provider.
func WithProviderInit(init string) ProviderOption {
	return func(p *ProviderProps) { p.OnInit = init }
}

// WithProviderClass adds custom CSS classes to the provider wrapper.
func WithProviderClass(class string) ProviderOption {
	return func(p *ProviderProps) { p.Class = class }
}

// WithProviderAttrs adds custom HTML attributes to the provider wrapper.
func WithProviderAttrs(attrs ...g.Node) ProviderOption {
	return func(p *ProviderProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// WithProviderChildren sets the children that will have access to this provider.
func WithProviderChildren(children ...g.Node) ProviderOption {
	return func(p *ProviderProps) { p.Children = append(p.Children, children...) }
}

// WithProviderDebug enables debug mode for the provider.
func WithProviderDebug(debug bool) ProviderOption {
	return func(p *ProviderProps) { p.Debug = debug }
}

// WithProviderHook adds a lifecycle hook to the provider.
// Common hooks: "onMount", "onUpdate", "onDestroy"
func WithProviderHook(name, code string) ProviderOption {
	return func(p *ProviderProps) {
		if p.Hooks == nil {
			p.Hooks = make(map[string]string)
		}

		p.Hooks[name] = code
	}
}

// defaultProviderProps returns default provider properties.
func defaultProviderProps() *ProviderProps {
	return &ProviderProps{
		State:    make(map[string]any),
		Hooks:    make(map[string]string),
		Debug:    false,
		Attrs:    []g.Node{},
		Children: []g.Node{},
	}
}

// rawAttr creates an attribute that escapes quotes but not JavaScript operators.
// rawHTMLAttr represents a raw HTML attribute that bypasses gomponents' HTML encoding.
// This is necessary for Alpine.js attributes containing JavaScript code.
type rawHTMLAttr struct {
	name  string
	value string
}

// Render implements the gomponents.Node interface for rawHTMLAttr.
func (r rawHTMLAttr) Render(w io.Writer) error {
	// Escape special characters for HTML attribute
	// Must escape & first to avoid double escaping
	val := strings.ReplaceAll(r.value, "&", "&amp;")
	val = strings.ReplaceAll(val, "<", "&lt;")
	val = strings.ReplaceAll(val, ">", "&gt;")
	val = strings.ReplaceAll(val, "\"", "&quot;")

	_, err := fmt.Fprintf(w, ` %s="%s"`, r.name, val)

	return err
}

// Type implements the gomponents.Node interface.
func (r rawHTMLAttr) Type() g.NodeType {
	return g.AttributeType
}

// Provider creates a provider context component that shares state with children.
//
// The Provider pattern enables deeply nested components to access shared state
// without prop drilling. It uses Alpine.js's x-data for scoped state and events
// for cross-component communication.
//
// Example:
//
//	primitives.Provider(
//	    primitives.WithProviderName("sidebar"),
//	    primitives.WithProviderState(map[string]any{
//	        "collapsed": false,
//	        "mobileOpen": false,
//	    }),
//	    primitives.WithProviderMethods(`
//	        toggle() {
//	            this.collapsed = !this.collapsed;
//	            this.$dispatch('sidebar:toggled', { collapsed: this.collapsed });
//	        }
//	    `),
//	    primitives.WithProviderChildren(
//	        // child components can access provider via $provider.sidebar
//	    ),
//	)
//
// Children can access provider state using:
// - Alpine expressions: $el.closest('[data-provider="sidebar"]').__x.$data.collapsed
// - Helper function: getProvider($el, 'sidebar')
// - Magic property (if registered globally): $provider.sidebar
func Provider(opts ...ProviderOption) g.Node {
	props := defaultProviderProps()
	for _, opt := range opts {
		opt(props)
	}

	if props.Name == "" {
		// Return children without provider wrapper if no name specified
		return g.Group(props.Children)
	}

	// Build x-data attribute value
	xDataValue := buildXDataValue(props)

	// Build x-init attribute value
	xInitValue := buildXInitValue(props)

	// Build attributes for provider wrapper
	attrs := []g.Node{
		g.Attr("data-provider", props.Name),
		// Use rawHTMLAttr to properly render Alpine.js attributes
		rawHTMLAttr{name: "x-data", value: xDataValue},
	}

	if xInitValue != "" {
		attrs = append(attrs, rawHTMLAttr{name: "x-init", value: xInitValue})
	}

	if props.Class != "" {
		attrs = append(attrs, html.Class(props.Class))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(props.Children),
	)
}

// buildXDataValue constructs the x-data attribute value from provider props.
// Generates JavaScript object literal syntax (not JSON) for Alpine.js compatibility.
func buildXDataValue(props *ProviderProps) string {
	if len(props.State) == 0 && props.Methods == "" {
		return "{}"
	}

	// Build state as JavaScript object literals (keys without quotes)
	var stateParts []string

	for key, value := range props.State {
		var jsValue string

		switch v := value.(type) {
		case string:
			// Escape backslashes first, then single quotes
			escaped := strings.ReplaceAll(v, `\`, `\\`)
			escaped = strings.ReplaceAll(escaped, `'`, `\'`)
			jsValue = fmt.Sprintf("'%s'", escaped)
		case bool:
			jsValue = strconv.FormatBool(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			jsValue = fmt.Sprintf("%v", v)
		case float32, float64:
			jsValue = fmt.Sprintf("%v", v)
		case nil:
			jsValue = "null"
		default:
			// Fallback to JSON for complex types, but convert to JS object literal
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				// Fallback to null on marshal error
				jsValue = "null"
			} else {
				jsValue = string(jsonBytes)
			}
		}
		// JavaScript object literal: key without quotes
		stateParts = append(stateParts, fmt.Sprintf("%s:%s", key, jsValue))
	}

	// If no methods, return just the state
	if props.Methods == "" {
		if len(stateParts) == 0 {
			return "{}"
		}

		return "{" + strings.Join(stateParts, ",") + "}"
	}

	// Add methods - clean up whitespace for HTML attributes
	methodsStr := strings.TrimSpace(props.Methods)
	// Remove any surrounding braces from methods if present
	methodsStr = strings.TrimPrefix(methodsStr, "{")
	methodsStr = strings.TrimSuffix(methodsStr, "}")
	// Remove newlines and tabs, collapse multiple spaces
	methodsStr = strings.ReplaceAll(methodsStr, "\n", " ")
	methodsStr = strings.ReplaceAll(methodsStr, "\r", " ")
	methodsStr = strings.ReplaceAll(methodsStr, "\t", " ")
	// Collapse multiple spaces to single space
	for strings.Contains(methodsStr, "  ") {
		methodsStr = strings.ReplaceAll(methodsStr, "  ", " ")
	}

	methodsStr = strings.TrimSpace(methodsStr)

	// Merge state and methods
	if len(stateParts) == 0 {
		return "{" + methodsStr + "}"
	}

	return "{" + strings.Join(stateParts, ",") + "," + methodsStr + "}"
}

// buildXInitValue constructs the x-init attribute value from provider props.
func buildXInitValue(props *ProviderProps) string {
	var initParts []string

	// Add debug logging if enabled
	if props.Debug {
		initParts = append(initParts, fmt.Sprintf(
			"console.log('[Provider:%s] Initialized with state:', this)",
			props.Name,
		))
	}

	// Add custom init code
	if props.OnInit != "" {
		// Normalize whitespace for HTML attributes
		initCode := strings.TrimSpace(props.OnInit)
		initCode = strings.ReplaceAll(initCode, "\n", " ")
		initCode = strings.ReplaceAll(initCode, "\r", " ")
		initCode = strings.ReplaceAll(initCode, "\t", " ")
		// Collapse multiple spaces
		for strings.Contains(initCode, "  ") {
			initCode = strings.ReplaceAll(initCode, "  ", " ")
		}

		initParts = append(initParts, strings.TrimSpace(initCode))
	}

	// Add lifecycle hooks
	if onMount, ok := props.Hooks["onMount"]; ok {
		initParts = append(initParts, strings.TrimSpace(onMount))
	}

	// Watch for state changes if debug mode is enabled
	if props.Debug {
		initParts = append(initParts, fmt.Sprintf(
			"$watch('$data', (value) => console.log('[Provider:%s] State changed:', value))",
			props.Name,
		))
	}

	// Add onUpdate hook as a watcher if present
	if onUpdate, ok := props.Hooks["onUpdate"]; ok {
		initParts = append(initParts, fmt.Sprintf(
			"$watch('$data', () => { %s })",
			strings.TrimSpace(onUpdate),
		))
	}

	if len(initParts) == 0 {
		return ""
	}

	return strings.Join(initParts, "; ")
}

// ProviderValue returns an Alpine expression to access a provider's state.
// Use this in x-bind, x-show, x-text, etc. to access provider state.
//
// Example:
//
//	g.Attr("x-show", primitives.ProviderValue("sidebar", "collapsed"))
//	// Generates: x-show="$el.closest('[data-provider=\"sidebar\"]').__x.$data.collapsed"
func ProviderValue(providerName, key string) string {
	return fmt.Sprintf("$el.closest('[data-provider=\"%s\"]').__x.$data.%s", providerName, key)
}

// ProviderMethod returns an Alpine expression to call a provider method.
//
// Example:
//
//	alpine.XClick(primitives.ProviderMethod("sidebar", "toggle"))
//	// Generates: @click="$el.closest('[data-provider=\"sidebar\"]').__x.$data.toggle()"
func ProviderMethod(providerName, method string, args ...string) string {
	argsStr := strings.Join(args, ", ")

	return fmt.Sprintf("$el.closest('[data-provider=\"%s\"]').__x.$data.%s(%s)",
		providerName, method, argsStr)
}

// ProviderDispatch returns an Alpine expression to dispatch a provider event.
// This is useful for cross-provider communication.
//
// Example:
//
//	alpine.XClick(primitives.ProviderDispatch("sidebar", "toggle", "{ open: true }"))
func ProviderDispatch(providerName, eventName, data string) string {
	if data == "" {
		data = "{}"
	}

	return fmt.Sprintf("$dispatch('provider:%s:%s', %s)", providerName, eventName, data)
}

// ProviderScriptUtilities returns a script tag with helper functions for working with providers.
// Include this once in your page head or before Alpine.js initialization.
//
// Provides global functions:
// - getProvider(el, name): Get provider state object
// - providerValue(el, name, key): Get specific provider state value
// - providerCall(el, name, method, ...args): Call provider method
// - providerDispatch(el, name, event, data): Dispatch provider event
func ProviderScriptUtilities() g.Node {
	return html.Script(g.Raw(`
// ForgeUI Provider Utilities
window.forgeui = window.forgeui || {};
window.forgeui.provider = {
  /**
   * Get a provider's state object
   * @param {HTMLElement} el - Element within provider context
   * @param {string} name - Provider name
   * @returns {Object} Provider state object
   */
  getProvider(el, name) {
    const providerEl = el.closest('[data-provider="' + name + '"]');
    if (!providerEl || !providerEl.__x) {
      console.warn('[ForgeUI] Provider "' + name + '" not found');
      return null;
    }
    return Alpine.raw(providerEl.__x.$data);
  },

  /**
   * Get a specific value from provider state
   * @param {HTMLElement} el - Element within provider context
   * @param {string} name - Provider name
   * @param {string} key - State key
   * @returns {*} State value
   */
  getValue(el, name, key) {
    const provider = this.getProvider(el, name);
    return provider ? provider[key] : undefined;
  },

  /**
   * Call a provider method
   * @param {HTMLElement} el - Element within provider context
   * @param {string} name - Provider name
   * @param {string} method - Method name
   * @param {...*} args - Method arguments
   * @returns {*} Method return value
   */
  call(el, name, method, ...args) {
    const provider = this.getProvider(el, name);
    if (provider && typeof provider[method] === 'function') {
      return provider[method](...args);
    }
    console.warn('[ForgeUI] Method "' + method + '" not found on provider "' + name + '"');
  },

  /**
   * Dispatch a provider event
   * @param {HTMLElement} el - Element to dispatch from
   * @param {string} name - Provider name
   * @param {string} event - Event name
   * @param {*} data - Event data
   */
  dispatch(el, name, event, data) {
    el.dispatchEvent(new CustomEvent('provider:' + name + ':' + event, {
      detail: data,
      bubbles: true,
      composed: true
    }));
  }
};

// Alpine magic property for easier access
document.addEventListener('alpine:init', () => {
  Alpine.magic('provider', (el) => ({
    get(name, key) {
      return window.forgeui.provider.getValue(el, name, key);
    },
    call(name, method, ...args) {
      return window.forgeui.provider.call(el, name, method, ...args);
    },
    dispatch(name, event, data) {
      window.forgeui.provider.dispatch(el, name, event, data);
    }
  }));
});
`))
}

// ProviderStack creates a nested stack of providers.
// This is useful when you need multiple providers in a component tree.
//
// Example:
//
//	primitives.ProviderStack(
//	    primitives.Provider(primitives.WithProviderName("theme"), ...),
//	    primitives.Provider(primitives.WithProviderName("sidebar"), ...),
//	    primitives.Provider(primitives.WithProviderName("auth"), ...),
//	)
func ProviderStack(providers ...g.Node) g.Node {
	if len(providers) == 0 {
		return g.Group(nil)
	}

	// Nest providers from outermost to innermost
	return g.Group(providers)
}
