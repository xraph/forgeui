// Package htmx provides HTMX integration for ForgeUI.
//
// HTMX allows you to access AJAX, CSS Transitions, WebSockets and Server Sent Events
// directly in HTML, using attributes. This enables building modern, interactive user
// interfaces while keeping your application architecture simple and maintainable.
//
// # Basic Usage
//
// Add HTMX attributes to your HTML elements:
//
//	html.Button(
//	    htmx.HxGet("/api/users"),
//	    htmx.HxTarget("#user-list"),
//	    htmx.HxSwap("innerHTML"),
//	    g.Text("Load Users"),
//	)
//
// # HTTP Methods
//
// HTMX supports all standard HTTP methods:
//
//	htmx.HxGet("/api/data")      // GET request
//	htmx.HxPost("/api/create")   // POST request
//	htmx.HxPut("/api/update/1")  // PUT request
//	htmx.HxPatch("/api/patch/1") // PATCH request
//	htmx.HxDelete("/api/del/1")  // DELETE request
//
// # Targeting and Swapping
//
// Control where and how responses are inserted:
//
//	html.Button(
//	    htmx.HxGet("/content"),
//	    htmx.HxTarget("#container"),
//	    htmx.HxSwap("beforeend"),  // Append to container
//	)
//
// Available swap strategies:
//   - innerHTML (default): Replace inner HTML
//   - outerHTML: Replace entire element
//   - beforebegin: Insert before element
//   - afterbegin: Insert as first child
//   - beforeend: Insert as last child
//   - afterend: Insert after element
//   - delete: Delete the target
//   - none: Don't swap
//
// # Triggers
//
// Specify when requests are made:
//
//	html.Input(
//	    htmx.HxGet("/search"),
//	    htmx.HxTriggerDebounce("keyup", "500ms"),
//	    htmx.HxTarget("#results"),
//	)
//
// Common triggers:
//   - click: On element click
//   - change: On input change
//   - submit: On form submit
//   - load: When element loads
//   - revealed: When scrolled into view
//   - intersect: When element intersects viewport
//
// # Server-Side Usage
//
// Detect HTMX requests and respond appropriately:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    if htmx.IsHTMX(r) {
//	        // Return partial HTML
//	        renderPartial(w, data)
//	    } else {
//	        // Return full page
//	        renderFullPage(w, data)
//	    }
//	}
//
// # Response Headers
//
// Control client behavior with response headers:
//
//	htmx.SetHTMXTrigger(w, map[string]any{
//	    "showMessage": map[string]string{
//	        "level": "success",
//	        "text": "Item saved!",
//	    },
//	})
//
// Available response headers:
//   - HX-Trigger: Trigger client-side events
//   - HX-Redirect: Perform client-side redirect
//   - HX-Refresh: Force full page refresh
//   - HX-Retarget: Change swap target
//   - HX-Reswap: Change swap strategy
//   - HX-Push-Url: Update browser URL
//   - HX-Replace-Url: Replace browser URL
//
// # Progressive Enhancement
//
// Use hx-boost for seamless navigation:
//
//	html.Div(
//	    htmx.HxBoost(true),
//	    html.A(html.Href("/page1"), g.Text("Page 1")),
//	    html.A(html.Href("/page2"), g.Text("Page 2")),
//	)
//
// # Loading States
//
// Show indicators during requests:
//
//	html.Button(
//	    htmx.HxPost("/api/submit"),
//	    htmx.HxIndicator("#spinner"),
//	    htmx.HxDisabledElt("this"),
//	    g.Text("Submit"),
//	)
//	html.Div(
//	    html.ID("spinner"),
//	    html.Class("htmx-indicator"),
//	    g.Text("Loading..."),
//	)
//
// # Extensions
//
// Load HTMX extensions for additional functionality:
//
//	htmx.ScriptsWithExtensions(
//	    htmx.ExtensionSSE,
//	    htmx.ExtensionWebSockets,
//	)
package htmx

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Default HTMX version
const DefaultVersion = "2.0.3"

// HTMX extension names
const (
	ExtensionSSE             = "sse"
	ExtensionWebSockets      = "ws"
	ExtensionClassTools      = "class-tools"
	ExtensionPreload         = "preload"
	ExtensionHeadSupport     = "head-support"
	ExtensionResponseTargets = "response-targets"
	ExtensionDebug           = "debug"
	ExtensionEventHeader     = "event-header"
	ExtensionIncludeVals     = "include-vals"
	ExtensionJSONEnc         = "json-enc"
	ExtensionMethodOverride  = "method-override"
	ExtensionMorphdom        = "morphdom-swap"
	ExtensionMultiSwap       = "multi-swap"
	ExtensionPathDeps        = "path-deps"
	ExtensionRestoreOnError  = "restored"
)

// Scripts returns a script tag that loads HTMX from CDN.
//
// Example:
//
//	html.Head(
//	    htmx.Scripts(),
//	)
func Scripts(version ...string) g.Node {
	ver := DefaultVersion
	if len(version) > 0 && version[0] != "" {
		ver = version[0]
	}

	return html.Script(
		html.Src("https://unpkg.com/htmx.org@"+ver),
		g.Attr("integrity", ""),
		g.Attr("crossorigin", "anonymous"),
	)
}

// ScriptsWithExtensions loads HTMX with the specified extensions.
//
// Example:
//
//	html.Head(
//	    g.Group(htmx.ScriptsWithExtensions(
//	        htmx.ExtensionSSE,
//	        htmx.ExtensionWebSockets,
//	    )),
//	)
func ScriptsWithExtensions(extensions ...string) []g.Node {
	nodes := []g.Node{Scripts()}

	for _, ext := range extensions {
		nodes = append(nodes, ExtensionScript(ext))
	}

	return nodes
}

// ExtensionScript returns a script tag for a specific HTMX extension.
//
// Example:
//
//	html.Head(
//	    htmx.Scripts(),
//	    htmx.ExtensionScript(htmx.ExtensionSSE),
//	)
func ExtensionScript(extension string) g.Node {
	return html.Script(
		html.Src(fmt.Sprintf("https://unpkg.com/htmx-ext-%s@2.0.0/ext/%s.js", extension, extension)),
		g.Attr("crossorigin", "anonymous"),
	)
}

// CloakCSS returns a style tag that prevents flash of unstyled content
// for elements with hx-cloak attribute.
//
// Add this to your page head and use HxCloak() on elements that should
// be hidden until HTMX initializes.
//
// Example:
//
//	html.Head(
//	    htmx.CloakCSS(),
//	)
//	html.Body(
//	    html.Div(
//	        htmx.HxCloak(),
//	        htmx.HxGet("/content"),
//	        htmx.HxTriggerLoad(),
//	    ),
//	)
func CloakCSS() g.Node {
	return html.StyleEl(g.Raw(`[hx-cloak] { display: none !important; }`))
}

// HxCloak creates an hx-cloak attribute to hide elements until HTMX processes them.
//
// Example:
//
//	html.Div(
//	    htmx.HxCloak(),
//	    htmx.HxGet("/content"),
//	    htmx.HxTriggerLoad(),
//	)
func HxCloak() g.Node {
	return g.Attr("hx-cloak", "")
}

// IndicatorCSS returns a style tag with default HTMX indicator styles.
//
// Elements with class .htmx-indicator will be hidden by default and shown
// during HTMX requests when targeted by hx-indicator.
//
// Example:
//
//	html.Head(
//	    htmx.IndicatorCSS(),
//	)
func IndicatorCSS() g.Node {
	return html.StyleEl(g.Raw(`
		.htmx-indicator {
			opacity: 0;
			transition: opacity 200ms ease-in;
		}
		.htmx-request .htmx-indicator,
		.htmx-request.htmx-indicator {
			opacity: 1;
		}
		.htmx-swapping {
			opacity: 0;
			transition: opacity 200ms ease-out;
		}
	`))
}

// ConfigMeta returns meta tags for HTMX configuration.
//
// Example:
//
//	html.Head(
//	    htmx.ConfigMeta(map[string]string{
//	        "defaultSwapStyle": "outerHTML",
//	        "timeout": "5000",
//	    }),
//	)
func ConfigMeta(config map[string]string) []g.Node {
	nodes := make([]g.Node, 0, len(config))
	for key, value := range config {
		nodes = append(nodes, html.Meta(
			html.Name("htmx-config:"+key),
			html.Content(value),
		))
	}

	return nodes
}

// HxDisinherit creates an hx-disinherit attribute that prevents inheritance
// of HTMX attributes from parent elements.
//
// Example:
//
//	html.Div(
//	    htmx.HxDisinherit("hx-select hx-get"),
//	    // child elements won't inherit hx-select or hx-get
//	)
func HxDisinherit(attributes string) g.Node {
	if attributes == "" {
		return g.Attr("hx-disinherit", "*")
	}

	return g.Attr("hx-disinherit", attributes)
}

// HxHistory creates an hx-history attribute to control history behavior.
//
// Example:
//
//	html.Div(
//	    htmx.HxHistory(false),
//	    htmx.HxBoost(true),
//	    // links won't be added to history
//	)
func HxHistory(enabled bool) g.Node {
	if enabled {
		return g.Attr("hx-history", "true")
	}

	return g.Attr("hx-history", "false")
}
