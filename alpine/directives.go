// Package alpine provides Alpine.js integration helpers for ForgeUI.
//
// Alpine.js is a lightweight JavaScript framework for adding interactivity
// to server-rendered HTML. This package provides Go functions that generate
// Alpine directives as HTML attributes.
//
// Basic usage:
//
//	html.Div(
//	    alpine.XData(map[string]any{"count": 0}),
//	    alpine.XClick("count++"),
//	    alpine.XText("count"),
//	)
package alpine

import (
	"encoding/json"
	"fmt"

	g "github.com/maragudk/gomponents"
)

// XData creates an x-data attribute with the given state.
// Pass nil or an empty map for an empty x-data scope.
// Supports RawJS values for embedding raw JavaScript code (functions, getters).
//
// Example:
//
//	alpine.XData(map[string]any{
//	    "open": false,
//	    "count": 0,
//	    "increment": alpine.RawJS("function() { this.count++ }"),
//	})
func XData(state map[string]any) g.Node {
	if state == nil || len(state) == 0 {
		return g.Attr("x-data", "")
	}

	// Check if any value is RawJS - if so, use custom serialization
	hasRawJS := false
	for _, v := range state {
		if _, ok := v.(rawJS); ok {
			hasRawJS = true
			break
		}
	}

	if !hasRawJS {
		// Simple case: no RawJS, use standard JSON
		jsonData, _ := json.Marshal(state)
		return g.Attr("x-data", string(jsonData))
	}

	// Complex case: contains RawJS, use custom serialization
	return g.Attr("x-data", serializeXDataWithRawJS(state))
}

// serializeXDataWithRawJS creates a JavaScript object string that can contain
// raw JavaScript expressions (for functions, getters, etc.)
func serializeXDataWithRawJS(state map[string]any) string {
	if len(state) == 0 {
		return "{}"
	}

	result := "{"
	first := true

	for key, value := range state {
		if !first {
			result += ","
		}
		first = false

		// Add key (quote if necessary)
		if needsQuoting(key) {
			result += fmt.Sprintf("%q", key)
		} else {
			result += key
		}
		result += ":"

		// Add value
		switch v := value.(type) {
		case rawJS:
			result += v.code
		case string:
			result += fmt.Sprintf("%q", v)
		case bool:
			if v {
				result += "true"
			} else {
				result += "false"
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result += fmt.Sprintf("%d", v)
		case float32, float64:
			result += fmt.Sprintf("%v", v)
		case nil:
			result += "null"
		default:
			// For complex types, use JSON encoding
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				result += "null"
			} else {
				result += string(jsonBytes)
			}
		}
	}

	result += "}"
	return result
}

// needsQuoting returns true if the key needs to be quoted in JavaScript
func needsQuoting(key string) bool {
	if len(key) == 0 {
		return true
	}

	// Check for spaces or special characters
	for _, r := range key {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '$') {
			return true
		}
	}

	// First character cannot be a digit
	if key[0] >= '0' && key[0] <= '9' {
		return true
	}

	return false
}

// Data is an alias for XData for consistency.
func Data(state map[string]any) g.Node {
	return XData(state)
}

// XShow creates an x-show directive for conditional display.
// The element will be hidden/shown based on the expression.
//
// Example:
//
//	alpine.XShow("open")
//	alpine.XShow("count > 5")
func XShow(expr string) g.Node {
	return g.Attr("x-show", expr)
}

// XIf creates an x-if directive for conditional rendering.
// Should be used on <template> elements. The element will be
// added/removed from the DOM based on the expression.
//
// Example:
//
//	html.Template(
//	    alpine.XIf("isVisible"),
//	    html.Div(g.Text("Content")),
//	)
func XIf(expr string) g.Node {
	return g.Attr("x-if", expr)
}

// XFor creates an x-for directive for list rendering.
// Should be used on <template> elements.
//
// Example:
//
//	alpine.XFor("item in items")
//	alpine.XFor("(item, index) in items")
func XFor(expr string) g.Node {
	return g.Attr("x-for", expr)
}

// XForKeyed creates an x-for directive with a :key binding
// for optimized list rendering.
//
// Example:
//
//	html.Template(
//	    g.Group(alpine.XForKeyed("item in items", "item.id")),
//	    html.Div(alpine.XText("item.name")),
//	)
func XForKeyed(expr, key string) []g.Node {
	return []g.Node{
		g.Attr("x-for", expr),
		g.Attr(":key", key),
	}
}

// XCloak creates an x-cloak attribute to prevent flash of
// unstyled content before Alpine initializes.
// Add CSS: [x-cloak] { display: none !important; }
func XCloak() g.Node {
	return g.Attr("x-cloak", "")
}

// XIgnore creates an x-ignore attribute to prevent Alpine
// from initializing this element and its children.
func XIgnore() g.Node {
	return g.Attr("x-ignore", "")
}

// XBind creates a dynamic attribute binding using the :attr shorthand.
//
// Example:
//
//	alpine.XBind("disabled", "loading")
//	alpine.XBind("href", "'/users/' + userId")
func XBind(attr, expr string) g.Node {
	return g.Attr(":"+attr, expr)
}

// XBindClass creates a :class binding for dynamic classes.
//
// Example:
//
//	alpine.XBindClass("{'active': isActive, 'disabled': isDisabled}")
//	alpine.XBindClass("isActive ? 'text-green-500' : 'text-gray-500'")
func XBindClass(expr string) g.Node {
	return g.Attr(":class", expr)
}

// XBindStyle creates a :style binding for dynamic styles.
//
// Example:
//
//	alpine.XBindStyle("{'color': textColor, 'fontSize': size + 'px'}")
func XBindStyle(expr string) g.Node {
	return g.Attr(":style", expr)
}

// XBindDisabled creates a :disabled binding.
func XBindDisabled(expr string) g.Node {
	return g.Attr(":disabled", expr)
}

// XBindValue creates a :value binding.
func XBindValue(expr string) g.Node {
	return g.Attr(":value", expr)
}

// XModel creates an x-model directive for two-way data binding.
//
// Example:
//
//	html.Input(
//	    html.Type("text"),
//	    alpine.XModel("name"),
//	)
func XModel(expr string) g.Node {
	return g.Attr("x-model", expr)
}

// XModelDebounce creates an x-model directive with debouncing.
//
// Example:
//
//	alpine.XModelDebounce("searchQuery", 300)
func XModelDebounce(expr string, ms int) g.Node {
	return g.Attr(fmt.Sprintf("x-model.debounce.%dms", ms), expr)
}

// XModelNumber creates an x-model.number directive that
// automatically converts the value to a number.
func XModelNumber(expr string) g.Node {
	return g.Attr("x-model.number", expr)
}

// XModelLazy creates an x-model.lazy directive that only
// updates on the 'change' event instead of 'input'.
func XModelLazy(expr string) g.Node {
	return g.Attr("x-model.lazy", expr)
}

// XOn creates an event listener using the @event shorthand.
//
// Example:
//
//	alpine.XOn("click", "count++")
//	alpine.XOn("submit", "handleSubmit()")
func XOn(event, handler string) g.Node {
	return g.Attr("@"+event, handler)
}

// XClick creates a @click event handler.
func XClick(handler string) g.Node {
	return XOn("click", handler)
}

// XSubmit creates a @submit.prevent event handler
// (prevents default form submission).
func XSubmit(handler string) g.Node {
	return XOn("submit.prevent", handler)
}

// XInput creates an @input event handler.
func XInput(handler string) g.Node {
	return XOn("input", handler)
}

// XChange creates an @change event handler.
func XChange(handler string) g.Node {
	return XOn("change", handler)
}

// XKeydown creates a @keydown event handler with optional key modifier.
//
// Example:
//
//	alpine.XKeydown("escape", "close()")
//	alpine.XKeydown("enter", "submit()")
func XKeydown(key, handler string) g.Node {
	if key == "" {
		return XOn("keydown", handler)
	}
	return XOn("keydown."+key, handler)
}

// XKeyup creates a @keyup event handler with optional key modifier.
func XKeyup(key, handler string) g.Node {
	if key == "" {
		return XOn("keyup", handler)
	}
	return XOn("keyup."+key, handler)
}

// XClickOutside creates a @click.outside event handler
// that triggers when clicking outside the element.
func XClickOutside(handler string) g.Node {
	return XOn("click.outside", handler)
}

// XClickOnce creates a @click.once event handler
// that only triggers once.
func XClickOnce(handler string) g.Node {
	return XOn("click.once", handler)
}

// XClickPrevent creates a @click.prevent event handler
// that prevents the default action.
func XClickPrevent(handler string) g.Node {
	return XOn("click.prevent", handler)
}

// XClickStop creates a @click.stop event handler
// that stops event propagation.
func XClickStop(handler string) g.Node {
	return XOn("click.stop", handler)
}

// XText creates an x-text directive to set text content.
//
// Example:
//
//	html.Span(alpine.XText("user.name"))
func XText(expr string) g.Node {
	return g.Attr("x-text", expr)
}

// XHtml creates an x-html directive to set HTML content.
// ⚠️ WARNING: Be careful with x-html as it can introduce XSS vulnerabilities.
// Only use with trusted content.
//
// Example:
//
//	html.Div(alpine.XHtml("sanitizedContent"))
func XHtml(expr string) g.Node {
	return g.Attr("x-html", expr)
}

// XRef creates an x-ref attribute to register an element reference.
// Access via $refs.name in Alpine expressions.
//
// Example:
//
//	html.Input(alpine.XRef("emailInput"))
//	html.Button(alpine.XClick("$refs.emailInput.focus()"))
func XRef(name string) g.Node {
	return g.Attr("x-ref", name)
}

// XInit creates an x-init directive that runs when Alpine initializes.
//
// Example:
//
//	alpine.XInit("console.log('initialized')")
//	alpine.XInit("fetch('/api/data').then(r => r.json()).then(d => items = d)")
func XInit(expr string) g.Node {
	return g.Attr("x-init", expr)
}

// XInitFetch creates an x-init directive that fetches data from a URL.
//
// Example:
//
//	alpine.XInitFetch("/api/users", "users")
func XInitFetch(url, target string) g.Node {
	return g.Attr("x-init", fmt.Sprintf("%s = await (await fetch('%s')).json()", target, url))
}

// XEffect creates an x-effect directive that automatically
// re-runs when any referenced data changes.
//
// Example:
//
//	alpine.XEffect("console.log('count is now', count)")
func XEffect(expr string) g.Node {
	return g.Attr("x-effect", expr)
}

// XTeleport creates an x-teleport directive to move an element
// to a different part of the DOM.
//
// Example:
//
//	alpine.XTeleport("body")
//	alpine.XTeleport("#modal-container")
func XTeleport(selector string) g.Node {
	return g.Attr("x-teleport", selector)
}

// XId creates an x-id directive for generating unique IDs
// that are consistent across server/client.
//
// Example:
//
//	alpine.XId("['input']")
func XId(items string) g.Node {
	return g.Attr("x-id", items)
}

