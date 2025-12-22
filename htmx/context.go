package htmx

import (
	"encoding/json"
	"net/http"
)

// IsHTMX checks if the request is from HTMX by looking for the HX-Request header.
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    if htmx.IsHTMX(r) {
//	        // Return partial HTML
//	        renderPartial(w)
//	    } else {
//	        // Return full page
//	        renderFullPage(w)
//	    }
//	}
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("Hx-Request") == "true"
}

// IsHTMXBoosted checks if the request is from an element using hx-boost.
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    if htmx.IsHTMXBoosted(r) {
//	        // Handle boosted navigation
//	    }
//	}
func IsHTMXBoosted(r *http.Request) bool {
	return r.Header.Get("Hx-Boosted") == "true"
}

// HTMXCurrentURL gets the current URL of the browser when the HTMX request was made.
//
// Example:
//
//	currentURL := htmx.HTMXCurrentURL(r)
func HTMXCurrentURL(r *http.Request) string {
	return r.Header.Get("Hx-Current-Url")
}

// HTMXHistoryRestoreRequest checks if this request is for history restoration
// after a miss in the local history cache.
//
// Example:
//
//	if htmx.HTMXHistoryRestoreRequest(r) {
//	    // Handle history restoration
//	}
func HTMXHistoryRestoreRequest(r *http.Request) bool {
	return r.Header.Get("Hx-History-Restore-Request") == "true"
}

// HTMXPrompt gets the user response to an hx-prompt.
//
// Example:
//
//	if prompt := htmx.HTMXPrompt(r); prompt != "" {
//	    // Use the prompt value
//	    processComment(prompt)
//	}
func HTMXPrompt(r *http.Request) string {
	return r.Header.Get("Hx-Prompt")
}

// HTMXTarget gets the ID of the target element if it exists.
//
// Example:
//
//	targetID := htmx.HTMXTarget(r)
func HTMXTarget(r *http.Request) string {
	return r.Header.Get("Hx-Target")
}

// HTMXTriggerName gets the name of the triggered element if it exists.
//
// Example:
//
//	triggerName := htmx.HTMXTriggerName(r)
func HTMXTriggerName(r *http.Request) string {
	return r.Header.Get("Hx-Trigger-Name")
}

// HTMXTrigger gets the ID of the triggered element if it exists.
//
// Example:
//
//	triggerID := htmx.HTMXTrigger(r)
func HTMXTrigger(r *http.Request) string {
	return r.Header.Get("Hx-Trigger")
}

// SetHTMXTrigger is a convenience function to set the HX-Trigger response header.
//
// Example:
//
//	htmx.SetHTMXTrigger(w, map[string]any{
//	    "itemUpdated": map[string]int{"id": 123},
//	})
func SetHTMXTrigger(w http.ResponseWriter, events map[string]any) {
	if jsonData, err := json.Marshal(events); err == nil {
		w.Header().Set("Hx-Trigger", string(jsonData))
	}
}

// SetHTMXLocation performs a client-side redirect without a full page reload.
//
// Example:
//
//	htmx.SetHTMXLocation(w, "/new-page")
func SetHTMXLocation(w http.ResponseWriter, path string) {
	w.Header().Set("Hx-Location", path)
}

// SetHTMXLocationWithContext performs a client-side redirect with context.
//
// Example:
//
//	htmx.SetHTMXLocationWithContext(w, map[string]any{
//	    "path": "/messages",
//	    "target": "#main",
//	    "swap": "innerHTML",
//	})
func SetHTMXLocationWithContext(w http.ResponseWriter, context map[string]any) {
	if jsonData, err := json.Marshal(context); err == nil {
		w.Header().Set("Hx-Location", string(jsonData))
	}
}

// SetHTMXRedirect performs a client-side redirect that does a full page reload.
//
// Example:
//
//	htmx.SetHTMXRedirect(w, "/login")
func SetHTMXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("Hx-Redirect", url)
}

// SetHTMXRefresh tells HTMX to do a full page refresh.
//
// Example:
//
//	htmx.SetHTMXRefresh(w)
func SetHTMXRefresh(w http.ResponseWriter) {
	w.Header().Set("Hx-Refresh", "true")
}

// SetHTMXReplaceURL replaces the current URL in the location bar.
//
// Example:
//
//	htmx.SetHTMXReplaceURL(w, "/new-url")
func SetHTMXReplaceURL(w http.ResponseWriter, url string) {
	w.Header().Set("Hx-Replace-Url", url)
}

// SetHTMXPushURL pushes a new URL into the browser history stack.
//
// Example:
//
//	htmx.SetHTMXPushURL(w, "/page/2")
func SetHTMXPushURL(w http.ResponseWriter, url string) {
	w.Header().Set("Hx-Push-Url", url)
}

// SetHTMXReswap allows you to specify how the response will be swapped.
//
// Example:
//
//	htmx.SetHTMXReswap(w, "outerHTML")
func SetHTMXReswap(w http.ResponseWriter, swapMethod string) {
	w.Header().Set("Hx-Reswap", swapMethod)
}

// SetHTMXRetarget allows you to specify a new target for the swap.
//
// Example:
//
//	htmx.SetHTMXRetarget(w, "#different-target")
func SetHTMXRetarget(w http.ResponseWriter, target string) {
	w.Header().Set("Hx-Retarget", target)
}

// SetHTMXReselect allows you to select a subset of the response to swap.
//
// Example:
//
//	htmx.SetHTMXReselect(w, "#content")
func SetHTMXReselect(w http.ResponseWriter, selector string) {
	w.Header().Set("Hx-Reselect", selector)
}
