package alpine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/a-h/templ"
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
func XData(state map[string]any) templ.Attributes {
	if len(state) == 0 {
		return templ.Attributes{"x-data": ""}
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
		jsonData, err := json.Marshal(state)
		if err != nil {
			// Fallback to empty object on marshal error
			return templ.Attributes{"x-data": "{}"}
		}

		return templ.Attributes{"x-data": string(jsonData)}
	}

	// Complex case: contains RawJS, use custom serialization
	return templ.Attributes{"x-data": serializeXDataWithRawJS(state)}
}

// serializeXDataWithRawJS creates a JavaScript object string that can contain
// raw JavaScript expressions (for functions, getters, etc.)
func serializeXDataWithRawJS(state map[string]any) string {
	if len(state) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{")

	first := true

	var resultSb72 strings.Builder

	for key, value := range state {
		if !first {
			resultSb72.WriteString(",")
		}

		first = false

		// Add key (quote if necessary)
		if needsQuoting(key) {
			resultSb72.WriteString(fmt.Sprintf("%q", key))
		} else {
			resultSb72.WriteString(key)
		}

		resultSb72.WriteString(":")

		// Add value
		switch v := value.(type) {
		case rawJS:
			result.WriteString(v.code)
		case string:
			result.WriteString(fmt.Sprintf("%q", v))
		case bool:
			if v {
				result.WriteString("true")
			} else {
				result.WriteString("false")
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result.WriteString(fmt.Sprintf("%d", v))
		case float32, float64:
			result.WriteString(fmt.Sprintf("%v", v))
		case nil:
			result.WriteString("null")
		default:
			// For complex types, use JSON encoding
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				result.WriteString("null")
			} else {
				result.Write(jsonBytes)
			}
		}
	}

	result.WriteString(resultSb72.String())

	result.WriteString("}")

	return result.String()
}

// needsQuoting returns true if the key needs to be quoted in JavaScript
func needsQuoting(key string) bool {
	if len(key) == 0 {
		return true
	}

	// Check for spaces or special characters
	for _, r := range key {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' && r != '$' {
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
func Data(state map[string]any) templ.Attributes {
	return XData(state)
}

// XShow creates an x-show directive for conditional display.
// The element will be hidden/shown based on the expression.
//
// Example:
//
//	alpine.XShow("open")
//	alpine.XShow("count > 5")
func XShow(expr string) templ.Attributes {
	return templ.Attributes{"x-show": expr}
}

// XIf creates an x-if directive for conditional rendering.
// Should be used on <template> elements. The element will be
// added/removed from the DOM based on the expression.
func XIf(expr string) templ.Attributes {
	return templ.Attributes{"x-if": expr}
}

// XFor creates an x-for directive for list rendering.
// Should be used on <template> elements.
//
// Example:
//
//	alpine.XFor("item in items")
//	alpine.XFor("(item, index) in items")
func XFor(expr string) templ.Attributes {
	return templ.Attributes{"x-for": expr}
}

// XForKeyed creates an x-for directive with a :key binding
// for optimized list rendering.
//
// Example (in .templ files):
//
//	<template { alpine.XForKeyed("item in items", "item.id")... }>
func XForKeyed(expr, key string) templ.Attributes {
	return templ.Attributes{
		"x-for": expr,
		":key":  key,
	}
}

// XCloak creates an x-cloak attribute to prevent flash of
// unstyled content before Alpine initializes.
// Add CSS: [x-cloak] { display: none !important; }
func XCloak() templ.Attributes {
	return templ.Attributes{"x-cloak": ""}
}

// XIgnore creates an x-ignore attribute to prevent Alpine
// from initializing this element and its children.
func XIgnore() templ.Attributes {
	return templ.Attributes{"x-ignore": ""}
}

// XBind creates a dynamic attribute binding using the :attr shorthand.
//
// Example:
//
//	alpine.XBind("disabled", "loading")
//	alpine.XBind("href", "'/users/' + userId")
func XBind(attr, expr string) templ.Attributes {
	return templ.Attributes{":" + attr: expr}
}

// XBindClass creates a :class binding for dynamic classes.
//
// Example:
//
//	alpine.XBindClass("{'active': isActive, 'disabled': isDisabled}")
//	alpine.XBindClass("isActive ? 'text-green-500' : 'text-gray-500'")
func XBindClass(expr string) templ.Attributes {
	return templ.Attributes{":class": expr}
}

// XBindStyle creates a :style binding for dynamic styles.
//
// Example:
//
//	alpine.XBindStyle("{'color': textColor, 'fontSize': size + 'px'}")
func XBindStyle(expr string) templ.Attributes {
	return templ.Attributes{":style": expr}
}

// XBindDisabled creates a :disabled binding.
func XBindDisabled(expr string) templ.Attributes {
	return templ.Attributes{":disabled": expr}
}

// XBindValue creates a :value binding.
func XBindValue(expr string) templ.Attributes {
	return templ.Attributes{":value": expr}
}

// XModel creates an x-model directive for two-way data binding.
//
// Example (in .templ files):
//
//	<input type="text" { alpine.XModel("name")... }/>
func XModel(expr string) templ.Attributes {
	return templ.Attributes{"x-model": expr}
}

// XModelDebounce creates an x-model directive with debouncing.
//
// Example:
//
//	alpine.XModelDebounce("searchQuery", 300)
func XModelDebounce(expr string, ms int) templ.Attributes {
	return templ.Attributes{fmt.Sprintf("x-model.debounce.%dms", ms): expr}
}

// XModelNumber creates an x-model.number directive that
// automatically converts the value to a number.
func XModelNumber(expr string) templ.Attributes {
	return templ.Attributes{"x-model.number": expr}
}

// XModelLazy creates an x-model.lazy directive that only
// updates on the 'change' event instead of 'input'.
func XModelLazy(expr string) templ.Attributes {
	return templ.Attributes{"x-model.lazy": expr}
}

// XOn creates an event listener using the @event shorthand.
//
// Example:
//
//	alpine.XOn("click", "count++")
//	alpine.XOn("submit", "handleSubmit()")
func XOn(event, handler string) templ.Attributes {
	return templ.Attributes{"@" + event: handler}
}

// XClick creates a @click event handler.
func XClick(handler string) templ.Attributes {
	return XOn("click", handler)
}

// XSubmit creates a @submit.prevent event handler
// (prevents default form submission).
func XSubmit(handler string) templ.Attributes {
	return XOn("submit.prevent", handler)
}

// XInput creates an @input event handler.
func XInput(handler string) templ.Attributes {
	return XOn("input", handler)
}

// XChange creates an @change event handler.
func XChange(handler string) templ.Attributes {
	return XOn("change", handler)
}

// XKeydown creates a @keydown event handler with optional key modifier.
//
// Example:
//
//	alpine.XKeydown("escape", "close()")
//	alpine.XKeydown("enter", "submit()")
func XKeydown(key, handler string) templ.Attributes {
	if key == "" {
		return XOn("keydown", handler)
	}

	return XOn("keydown."+key, handler)
}

// XKeyup creates a @keyup event handler with optional key modifier.
func XKeyup(key, handler string) templ.Attributes {
	if key == "" {
		return XOn("keyup", handler)
	}

	return XOn("keyup."+key, handler)
}

// XClickOutside creates a @click.outside event handler
// that triggers when clicking outside the element.
func XClickOutside(handler string) templ.Attributes {
	return XOn("click.outside", handler)
}

// XClickOnce creates a @click.once event handler
// that only triggers once.
func XClickOnce(handler string) templ.Attributes {
	return XOn("click.once", handler)
}

// XClickPrevent creates a @click.prevent event handler
// that prevents the default action.
func XClickPrevent(handler string) templ.Attributes {
	return XOn("click.prevent", handler)
}

// XClickStop creates a @click.stop event handler
// that stops event propagation.
func XClickStop(handler string) templ.Attributes {
	return XOn("click.stop", handler)
}

// XText creates an x-text directive to set text content.
//
// Example (in .templ files):
//
//	<span { alpine.XText("user.name")... }></span>
func XText(expr string) templ.Attributes {
	return templ.Attributes{"x-text": expr}
}

// XHtml creates an x-html directive to set HTML content.
// WARNING: Be careful with x-html as it can introduce XSS vulnerabilities.
// Only use with trusted content.
func XHtml(expr string) templ.Attributes {
	return templ.Attributes{"x-html": expr}
}

// XRef creates an x-ref attribute to register an element reference.
// Access via $refs.name in Alpine expressions.
func XRef(name string) templ.Attributes {
	return templ.Attributes{"x-ref": name}
}

// XInit creates an x-init directive that runs when Alpine initializes.
//
// Example:
//
//	alpine.XInit("console.log('initialized')")
func XInit(expr string) templ.Attributes {
	return templ.Attributes{"x-init": expr}
}

// XInitFetch creates an x-init directive that fetches data from a URL.
//
// Example:
//
//	alpine.XInitFetch("/api/users", "users")
func XInitFetch(url, target string) templ.Attributes {
	return templ.Attributes{"x-init": fmt.Sprintf("%s = await (await fetch('%s')).json()", target, url)}
}

// XEffect creates an x-effect directive that automatically
// re-runs when any referenced data changes.
func XEffect(expr string) templ.Attributes {
	return templ.Attributes{"x-effect": expr}
}

// XTeleport creates an x-teleport directive to move an element
// to a different part of the DOM.
//
// Example:
//
//	alpine.XTeleport("body")
//	alpine.XTeleport("#modal-container")
func XTeleport(selector string) templ.Attributes {
	return templ.Attributes{"x-teleport": selector}
}

// XId creates an x-id directive for generating unique IDs
// that are consistent across server/client.
func XId(items string) templ.Attributes {
	return templ.Attributes{"x-id": items}
}
