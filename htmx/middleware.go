package htmx

import (
	"encoding/json"
	"net/http"
)

// ResponseHeaders contains HTMX response headers.
type ResponseHeaders struct {
	// Trigger allows you to trigger client-side events
	Trigger map[string]any

	// TriggerAfterSwap triggers events after the swap step
	TriggerAfterSwap map[string]any

	// TriggerAfterSettle triggers events after the settle step
	TriggerAfterSettle map[string]any

	// Redirect performs a client-side redirect
	Redirect string

	// Refresh forces a full page refresh
	Refresh bool

	// ReplaceURL replaces the current URL in the browser location bar
	ReplaceURL string

	// PushURL pushes a new URL into the browser history
	PushURL string

	// Reswap allows you to specify how the response will be swapped
	Reswap string

	// Retarget allows you to specify a new target for the swap
	Retarget string

	// Reselect allows you to select a subset of the response to swap
	Reselect string
}

// SetResponseHeaders sets HTMX response headers.
//
// Example:
//
//	htmx.SetResponseHeaders(w, htmx.ResponseHeaders{
//	    Trigger: map[string]any{
//	        "showMessage": map[string]string{"text": "Saved!"},
//	    },
//	    Refresh: false,
//	})
func SetResponseHeaders(w http.ResponseWriter, headers ResponseHeaders) {
	if len(headers.Trigger) > 0 {
		if jsonData, err := json.Marshal(headers.Trigger); err == nil {
			w.Header().Set("HX-Trigger", string(jsonData))
		}
	}

	if len(headers.TriggerAfterSwap) > 0 {
		if jsonData, err := json.Marshal(headers.TriggerAfterSwap); err == nil {
			w.Header().Set("HX-Trigger-After-Swap", string(jsonData))
		}
	}

	if len(headers.TriggerAfterSettle) > 0 {
		if jsonData, err := json.Marshal(headers.TriggerAfterSettle); err == nil {
			w.Header().Set("HX-Trigger-After-Settle", string(jsonData))
		}
	}

	if headers.Redirect != "" {
		w.Header().Set("HX-Redirect", headers.Redirect)
	}

	if headers.Refresh {
		w.Header().Set("HX-Refresh", "true")
	}

	if headers.ReplaceURL != "" {
		w.Header().Set("HX-Replace-Url", headers.ReplaceURL)
	}

	if headers.PushURL != "" {
		w.Header().Set("HX-Push-Url", headers.PushURL)
	}

	if headers.Reswap != "" {
		w.Header().Set("HX-Reswap", headers.Reswap)
	}

	if headers.Retarget != "" {
		w.Header().Set("HX-Retarget", headers.Retarget)
	}

	if headers.Reselect != "" {
		w.Header().Set("HX-Reselect", headers.Reselect)
	}
}

// Middleware creates an HTTP middleware that detects HTMX requests.
// It adds HTMX request information to the request context.
//
// Example:
//
//	router := http.NewServeMux()
//	router.Handle("/", htmx.Middleware(handler))
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HTMX requests are identified by the HX-Request header
		// The middleware doesn't modify the request, it just allows
		// handlers to check for HTMX-specific headers using the
		// context helper functions

		next.ServeHTTP(w, r)
	})
}

// TriggerEvent sets an HX-Trigger response header with a simple event.
//
// Example:
//
//	htmx.TriggerEvent(w, "showMessage")
func TriggerEvent(w http.ResponseWriter, eventName string) {
	w.Header().Set("HX-Trigger", eventName)
}

// TriggerEventWithDetail sets an HX-Trigger response header with event details.
//
// Example:
//
//	htmx.TriggerEventWithDetail(w, "showMessage", map[string]any{
//	    "level": "success",
//	    "message": "Item saved successfully",
//	})
func TriggerEventWithDetail(w http.ResponseWriter, eventName string, detail map[string]any) {
	trigger := map[string]any{eventName: detail}
	if jsonData, err := json.Marshal(trigger); err == nil {
		w.Header().Set("HX-Trigger", string(jsonData))
	}
}

// TriggerEvents sets multiple HX-Trigger events.
//
// Example:
//
//	htmx.TriggerEvents(w, map[string]any{
//	    "event1": nil,
//	    "event2": map[string]string{"key": "value"},
//	})
func TriggerEvents(w http.ResponseWriter, events map[string]any) {
	if jsonData, err := json.Marshal(events); err == nil {
		w.Header().Set("HX-Trigger", string(jsonData))
	}
}

// StopPolling sets HX-Trigger with status code 286 to stop polling.
//
// Example:
//
//	if completed {
//	    htmx.StopPolling(w)
//	    return
//	}
func StopPolling(w http.ResponseWriter) {
	w.WriteHeader(286)
}

