package primitives

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// ProviderProps defines the configuration for a Provider component.
type ProviderProps struct {
	// Name is the unique identifier for this provider (e.g., "sidebar", "theme")
	Name string

	// State is the initial state as a map of key-value pairs
	State map[string]any

	// Methods is raw JavaScript code defining provider methods
	Methods string

	// OnInit is JavaScript code that runs when the provider initializes
	OnInit string

	// Class adds custom CSS classes to the provider wrapper
	Class string

	// Attributes adds custom HTML attributes to the provider wrapper
	Attributes templ.Attributes

	// Children are the child components that will have access to this provider
	Children []templ.Component

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
func WithProviderAttrs(attrs templ.Attributes) ProviderOption {
	return func(p *ProviderProps) {
		if p.Attributes == nil {
			p.Attributes = templ.Attributes{}
		}
		for k, v := range attrs {
			p.Attributes[k] = v
		}
	}
}

// WithProviderChildren sets the children that will have access to this provider.
func WithProviderChildren(children ...templ.Component) ProviderOption {
	return func(p *ProviderProps) { p.Children = append(p.Children, children...) }
}

// WithProviderDebug enables debug mode for the provider.
func WithProviderDebug(debug bool) ProviderOption {
	return func(p *ProviderProps) { p.Debug = debug }
}

// WithProviderHook adds a lifecycle hook to the provider.
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
		State:      make(map[string]any),
		Hooks:      make(map[string]string),
		Debug:      false,
		Attributes: templ.Attributes{},
		Children:   []templ.Component{},
	}
}

// Provider creates a provider context component that shares state with children.
func Provider(opts ...ProviderOption) templ.Component {
	props := defaultProviderProps()
	for _, opt := range opts {
		opt(props)
	}

	if props.Name == "" {
		// Return children without provider wrapper if no name specified
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			return renderChildren(ctx, w, props.Children)
		})
	}

	// Build x-data and x-init attribute values
	xDataValue := buildXDataValue(props)
	xInitValue := buildXInitValue(props)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Build merged attributes
		attrs := templ.Attributes{}
		attrs["data-provider"] = props.Name

		// x-data needs HTML-entity-safe rendering
		// We handle it via writeRawAttr below

		if props.Class != "" {
			attrs["class"] = props.Class
		}

		// Merge user attributes
		for k, v := range props.Attributes {
			attrs[k] = v
		}

		if _, err := io.WriteString(w, "<div"); err != nil {
			return err
		}
		if err := writeAttrs(w, attrs); err != nil {
			return err
		}
		// Write x-data with HTML-entity escaping
		if err := writeRawAttr(w, "x-data", xDataValue); err != nil {
			return err
		}
		if xInitValue != "" {
			if err := writeRawAttr(w, "x-init", xInitValue); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}

		if err := renderChildren(ctx, w, props.Children); err != nil {
			return err
		}

		_, err := io.WriteString(w, "</div>")
		return err
	})
}

// writeRawAttr writes an HTML attribute with proper entity escaping for Alpine.js.
func writeRawAttr(w io.Writer, name, value string) error {
	val := strings.ReplaceAll(value, "&", "&amp;")
	val = strings.ReplaceAll(val, "<", "&lt;")
	val = strings.ReplaceAll(val, ">", "&gt;")
	val = strings.ReplaceAll(val, "\"", "&quot;")

	_, err := fmt.Fprintf(w, ` %s="%s"`, name, val)
	return err
}

// buildXDataValue constructs the x-data attribute value from provider props.
func buildXDataValue(props *ProviderProps) string {
	if len(props.State) == 0 && props.Methods == "" {
		return "{}"
	}

	var stateParts []string

	for key, value := range props.State {
		var jsValue string

		switch v := value.(type) {
		case string:
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
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				jsValue = "null"
			} else {
				jsValue = string(jsonBytes)
			}
		}
		stateParts = append(stateParts, fmt.Sprintf("%s:%s", key, jsValue))
	}

	if props.Methods == "" {
		if len(stateParts) == 0 {
			return "{}"
		}
		return "{" + strings.Join(stateParts, ",") + "}"
	}

	methodsStr := strings.TrimSpace(props.Methods)
	methodsStr = strings.TrimPrefix(methodsStr, "{")
	methodsStr = strings.TrimSuffix(methodsStr, "}")
	methodsStr = strings.ReplaceAll(methodsStr, "\n", " ")
	methodsStr = strings.ReplaceAll(methodsStr, "\r", " ")
	methodsStr = strings.ReplaceAll(methodsStr, "\t", " ")
	for strings.Contains(methodsStr, "  ") {
		methodsStr = strings.ReplaceAll(methodsStr, "  ", " ")
	}
	methodsStr = strings.TrimSpace(methodsStr)

	if len(stateParts) == 0 {
		return "{" + methodsStr + "}"
	}

	return "{" + strings.Join(stateParts, ",") + "," + methodsStr + "}"
}

// buildXInitValue constructs the x-init attribute value from provider props.
func buildXInitValue(props *ProviderProps) string {
	var initParts []string

	if props.Debug {
		initParts = append(initParts, fmt.Sprintf(
			"console.log('[Provider:%s] Initialized with state:', this)",
			props.Name,
		))
	}

	if props.OnInit != "" {
		initCode := strings.TrimSpace(props.OnInit)
		initCode = strings.ReplaceAll(initCode, "\n", " ")
		initCode = strings.ReplaceAll(initCode, "\r", " ")
		initCode = strings.ReplaceAll(initCode, "\t", " ")
		for strings.Contains(initCode, "  ") {
			initCode = strings.ReplaceAll(initCode, "  ", " ")
		}
		initParts = append(initParts, strings.TrimSpace(initCode))
	}

	if onMount, ok := props.Hooks["onMount"]; ok {
		initParts = append(initParts, strings.TrimSpace(onMount))
	}

	if props.Debug {
		initParts = append(initParts, fmt.Sprintf(
			"$watch('$data', (value) => console.log('[Provider:%s] State changed:', value))",
			props.Name,
		))
	}

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
func ProviderValue(providerName, key string) string {
	return fmt.Sprintf("$el.closest('[data-provider=\"%s\"]').__x.$data.%s", providerName, key)
}

// ProviderMethod returns an Alpine expression to call a provider method.
func ProviderMethod(providerName, method string, args ...string) string {
	argsStr := strings.Join(args, ", ")
	return fmt.Sprintf("$el.closest('[data-provider=\"%s\"]').__x.$data.%s(%s)",
		providerName, method, argsStr)
}

// ProviderDispatch returns an Alpine expression to dispatch a provider event.
func ProviderDispatch(providerName, eventName, data string) string {
	if data == "" {
		data = "{}"
	}
	return fmt.Sprintf("$dispatch('provider:%s:%s', %s)", providerName, eventName, data)
}

// ProviderScriptUtilities returns a script tag with helper functions for working with providers.
func ProviderScriptUtilities() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script>
// ForgeUI Provider Utilities
window.forgeui = window.forgeui || {};
window.forgeui.provider = {
  getProvider(el, name) {
    const providerEl = el.closest('[data-provider="' + name + '"]');
    if (!providerEl || !providerEl.__x) {
      console.warn('[ForgeUI] Provider "' + name + '" not found');
      return null;
    }
    return Alpine.raw(providerEl.__x.$data);
  },
  getValue(el, name, key) {
    const provider = this.getProvider(el, name);
    return provider ? provider[key] : undefined;
  },
  call(el, name, method, ...args) {
    const provider = this.getProvider(el, name);
    if (provider && typeof provider[method] === 'function') {
      return provider[method](...args);
    }
    console.warn('[ForgeUI] Method "' + method + '" not found on provider "' + name + '"');
  },
  dispatch(el, name, event, data) {
    el.dispatchEvent(new CustomEvent('provider:' + name + ':' + event, {
      detail: data,
      bubbles: true,
      composed: true
    }));
  }
};
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
</script>`)
		return err
	})
}

// ProviderStack creates a nested stack of providers.
func ProviderStack(providers ...templ.Component) templ.Component {
	if len(providers) == 0 {
		return templ.NopComponent
	}
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return renderChildren(ctx, w, providers)
	})
}
