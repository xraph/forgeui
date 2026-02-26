// Package htmx provides HTMX integration for ForgeUI.
//
// HTMX allows you to access AJAX, CSS Transitions, WebSockets and Server Sent Events
// directly in HTML, using attributes. This enables building modern, interactive user
// interfaces while keeping your application architecture simple and maintainable.
//
// # Basic Usage
//
// Add HTMX attributes to your templ components:
//
//	<button { htmx.HxGet("/api/users")... } { htmx.HxTarget("#user-list")... } { htmx.HxSwap("innerHTML")... }>
//	    Load Users
//	</button>
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
//	<button { htmx.HxGet("/content")... } { htmx.HxTarget("#container")... } { htmx.HxSwap("beforeend")... }>
//	    Load
//	</button>
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
//	<input { htmx.HxGet("/search")... } { htmx.HxTriggerDebounce("keyup", "500ms")... } { htmx.HxTarget("#results")... }/>
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
// # Progressive Enhancement
//
// Use hx-boost for seamless navigation:
//
//	<div { htmx.HxBoost(true)... }>
//	    <a href="/page1">Page 1</a>
//	    <a href="/page2">Page 2</a>
//	</div>
//
// # Loading States
//
// Show indicators during requests:
//
//	<button { htmx.HxPost("/api/submit")... } { htmx.HxIndicator("#spinner")... } { htmx.HxDisabledElt("this")... }>
//	    Submit
//	</button>
//	<div id="spinner" class="htmx-indicator">Loading...</div>
//
// # Extensions
//
// Load HTMX extensions for additional functionality:
//
//	@htmx.ScriptsWithExtensions(htmx.ExtensionSSE, htmx.ExtensionWebSockets)
package htmx

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

// DefaultVersion is the default HTMX version.
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

// Scripts returns a templ.Component that renders the HTMX script tag from CDN.
//
// Example (in .templ files):
//
//	@htmx.Scripts()
func Scripts(version ...string) templ.Component {
	ver := DefaultVersion
	if len(version) > 0 && version[0] != "" {
		ver = version[0]
	}

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<script src="https://unpkg.com/htmx.org@%s" crossorigin="anonymous"></script>`, ver)
		return err
	})
}

// ScriptsWithExtensions returns a templ.Component that loads HTMX with the specified extensions.
//
// Example (in .templ files):
//
//	@htmx.ScriptsWithExtensions(htmx.ExtensionSSE, htmx.ExtensionWebSockets)
func ScriptsWithExtensions(extensions ...string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := Scripts().Render(ctx, w); err != nil {
			return err
		}

		for _, ext := range extensions {
			if err := ExtensionScript(ext).Render(ctx, w); err != nil {
				return err
			}
		}

		return nil
	})
}

// ExtensionScript returns a templ.Component for a specific HTMX extension script tag.
//
// Example (in .templ files):
//
//	@htmx.ExtensionScript(htmx.ExtensionSSE)
func ExtensionScript(extension string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<script src="https://unpkg.com/htmx-ext-%s@2.0.0/ext/%s.js" crossorigin="anonymous"></script>`, extension, extension)
		return err
	})
}

// CloakCSS returns a templ.Component with CSS to prevent flash of unstyled content
// for elements with hx-cloak attribute.
//
// Example (in .templ files):
//
//	@htmx.CloakCSS()
func CloakCSS() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<style>[hx-cloak] { display: none !important; }</style>`)
		return err
	})
}

// HxCloak creates an hx-cloak attribute to hide elements until HTMX processes them.
//
// Example (in .templ files):
//
//	<div { htmx.HxCloak()... } { htmx.HxGet("/content")... }>
func HxCloak() templ.Attributes {
	return templ.Attributes{"hx-cloak": ""}
}

// IndicatorCSS returns a templ.Component with default HTMX indicator styles.
//
// Example (in .templ files):
//
//	@htmx.IndicatorCSS()
func IndicatorCSS() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<style>
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
	</style>`)
		return err
	})
}

// ConfigMeta returns a templ.Component that renders meta tags for HTMX configuration.
//
// Example (in .templ files):
//
//	@htmx.ConfigMeta(map[string]string{"defaultSwapStyle": "outerHTML"})
func ConfigMeta(config map[string]string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		for key, value := range config {
			if _, err := fmt.Fprintf(w, `<meta name="htmx-config:%s" content="%s"/>`, key, value); err != nil {
				return err
			}
		}
		return nil
	})
}

// HxDisinherit creates an hx-disinherit attribute that prevents inheritance
// of HTMX attributes from parent elements.
//
// Example (in .templ files):
//
//	<div { htmx.HxDisinherit("hx-select hx-get")... }>
func HxDisinherit(attributes string) templ.Attributes {
	if attributes == "" {
		return templ.Attributes{"hx-disinherit": "*"}
	}

	return templ.Attributes{"hx-disinherit": attributes}
}

// HxHistory creates an hx-history attribute to control history behavior.
//
// Example (in .templ files):
//
//	<div { htmx.HxHistory(false)... } { htmx.HxBoost(true)... }>
func HxHistory(enabled bool) templ.Attributes {
	if enabled {
		return templ.Attributes{"hx-history": "true"}
	}

	return templ.Attributes{"hx-history": "false"}
}
