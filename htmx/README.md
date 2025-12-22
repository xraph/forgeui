## HTMX Package

The `htmx` package provides seamless HTMX integration for ForgeUI, enabling modern, interactive user interfaces with server-side rendering.

## Features

- 🚀 **Full HTMX Support**: Complete attribute helpers for all HTMX features
- 🎯 **Type-Safe**: Go functions for all HTMX attributes
- 🔄 **Request Detection**: Server-side helpers to detect and handle HTMX requests
- 📡 **Response Control**: Programmatic control over HTMX response headers
- ⚡ **Progressive Enhancement**: Build with hx-boost for seamless navigation
- 🎨 **Loading States**: Built-in indicator and disabled element support

## Installation

```bash
go get github.com/xraph/forgeui/htmx
```

## Quick Start

### 1. Load HTMX

```go
html.Head(
    htmx.Scripts(),
    htmx.IndicatorCSS(),
)
```

### 2. Add HTMX Attributes

```go
html.Button(
    htmx.HxGet("/api/users"),
    htmx.HxTarget("#user-list"),
    htmx.HxSwap("innerHTML"),
    g.Text("Load Users"),
)
```

### 3. Server-Side Handling

```go
func handler(w http.ResponseWriter, r *http.Request) {
    if htmx.IsHTMX(r) {
        // Return partial HTML
        renderUserList(w, users)
    } else {
        // Return full page
        renderFullPage(w, users)
    }
}
```

## HTTP Methods

### Basic Requests

```go
// GET request
html.Button(
    htmx.HxGet("/api/data"),
    htmx.HxTarget("#results"),
)

// POST request
html.Form(
    htmx.HxPost("/api/submit"),
    // form fields...
)

// Other methods
htmx.HxPut("/api/update/123")
htmx.HxPatch("/api/partial-update/123")
htmx.HxDelete("/api/delete/123")
```

## Targeting & Swapping

### Target Selection

```go
htmx.HxTarget("#container")     // CSS selector
htmx.HxTarget("this")            // Current element
htmx.HxTarget("closest div")     // Closest parent div
htmx.HxTarget("next .item")      // Next sibling with class
htmx.HxTarget("previous input")  // Previous sibling input
```

### Swap Strategies

```go
// Replace strategies
htmx.HxSwapInnerHTML()    // Replace inner HTML (default)
htmx.HxSwapOuterHTML()    // Replace entire element

// Insert strategies
htmx.HxSwapBeforeBegin()  // Insert before element
htmx.HxSwapAfterBegin()   // Insert as first child
htmx.HxSwapBeforeEnd()    // Insert as last child
htmx.HxSwapAfterEnd()     // Insert after element

// Special strategies
htmx.HxSwapDelete()       // Delete target
htmx.HxSwapNone()         // Don't swap
```

## Triggers

### Basic Triggers

```go
html.Input(
    htmx.HxTriggerChange(),      // On change event
    htmx.HxGet("/search"),
)

html.Button(
    htmx.HxTriggerClick(),       // On click event
    htmx.HxPost("/action"),
)
```

### Advanced Triggers

```go
// Debouncing (wait for pause in events)
html.Input(
    htmx.HxTriggerDebounce("keyup", "500ms"),
    htmx.HxGet("/search"),
)

// Throttling (limit event frequency)
html.Div(
    htmx.HxTriggerThrottle("scroll", "1s"),
    htmx.HxGet("/more"),
)

// On load
html.Div(
    htmx.HxTriggerLoad(),
    htmx.HxGet("/initial-data"),
)

// When revealed (scrolled into view)
html.Div(
    htmx.HxTriggerRevealed(),
    htmx.HxGet("/lazy-load"),
)

// Fire only once
html.Button(
    htmx.HxTriggerOnce("click"),
    htmx.HxPost("/init"),
)

// Polling
html.Div(
    htmx.HxTriggerEvery("2s"),
    htmx.HxGet("/status"),
)
```

## Loading States

### Indicators

```go
html.Button(
    htmx.HxPost("/submit"),
    htmx.HxIndicator("#spinner"),
    g.Text("Submit"),
)

html.Div(
    html.ID("spinner"),
    html.Class("htmx-indicator"),
    icons.Loader(icons.WithClass("animate-spin")),
)
```

### Disabled Elements

```go
html.Button(
    htmx.HxPost("/submit"),
    htmx.HxDisabledElt("this"),  // Disable button during request
    g.Text("Submit"),
)
```

### Request Synchronization

```go
// Drop subsequent requests
html.Input(
    htmx.HxGet("/search"),
    htmx.HxSync("this", "drop"),
)

// Abort current request
html.Input(
    htmx.HxGet("/search"),
    htmx.HxSync("this", "abort"),
)

// Replace current request
html.Input(
    htmx.HxGet("/search"),
    htmx.HxSync("this", "replace"),
)
```

## Progressive Enhancement

### Boosted Navigation

```go
// Boost all links in a container
html.Nav(
    htmx.HxBoost(true),
    html.A(html.Href("/page1"), g.Text("Page 1")),
    html.A(html.Href("/page2"), g.Text("Page 2")),
)
```

### History Management

```go
// Push URL to history
html.Button(
    htmx.HxGet("/page/2"),
    htmx.HxPushURL(true),
    g.Text("Next Page"),
)

// Replace URL in history
html.Button(
    htmx.HxGet("/search?q=test"),
    htmx.HxReplaceURL(true),
    g.Text("Search"),
)

// Custom path in history
html.Button(
    htmx.HxGet("/api/page/2"),
    htmx.HxPushURLWithPath("/page/2"),
    g.Text("Next"),
)
```

## Server-Side

### Request Detection

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Check if request is from HTMX
    if htmx.IsHTMX(r) {
        // Return partial HTML
    }

    // Check if request is boosted
    if htmx.IsHTMXBoosted(r) {
        // Handle boosted request
    }

    // Get request details
    target := htmx.HTMXTarget(r)        // Target element
    trigger := htmx.HTMXTrigger(r)      // Triggering element
    prompt := htmx.HTMXPrompt(r)        // Prompt response
    currentURL := htmx.HTMXCurrentURL(r) // Current URL
}
```

### Response Headers

```go
// Trigger client-side events
htmx.SetHTMXTrigger(w, map[string]any{
    "showMessage": map[string]string{
        "level": "success",
        "text": "Saved successfully!",
    },
})

// Simple event
htmx.TriggerEvent(w, "refresh")

// Event with details
htmx.TriggerEventWithDetail(w, "itemUpdated", map[string]any{
    "id": 123,
})

// Client-side redirect
htmx.SetHTMXRedirect(w, "/login")

// Full page refresh
htmx.SetHTMXRefresh(w)

// Update URL
htmx.SetHTMXPushURL(w, "/new-url")
htmx.SetHTMXReplaceURL(w, "/current-url")

// Change swap target/strategy
htmx.SetHTMXRetarget(w, "#different-element")
htmx.SetHTMXReswap(w, "outerHTML")
```

### Response Header Helper

```go
htmx.SetResponseHeaders(w, htmx.ResponseHeaders{
    Trigger: map[string]any{
        "showToast": map[string]string{"msg": "Success!"},
    },
    Refresh: false,
    PushURL: "/page/2",
})
```

## Request Configuration

### Custom Headers

```go
html.Button(
    htmx.HxPost("/api/data"),
    htmx.HxHeaders(map[string]string{
        "X-API-Key": "secret",
        "X-Request-ID": "123",
    }),
)
```

### Extra Values

```go
html.Button(
    htmx.HxPost("/api/data"),
    htmx.HxVals(map[string]any{
        "category": "urgent",
        "priority": 1,
    }),
)

// JavaScript values
html.Button(
    htmx.HxPost("/api/data"),
    htmx.HxValsJS("js:{timestamp: Date.now()}"),
)
```

### Parameter Filtering

```go
// Include only specific params
html.Form(
    htmx.HxPost("/submit"),
    htmx.HxParams("email,name"),
)

// Include all params
htmx.HxParams("*")

// Include no params
htmx.HxParams("none")

// Exclude specific params
htmx.HxParams("not password,token")
```

### Confirmation & Prompts

```go
// Confirmation dialog
html.Button(
    htmx.HxDelete("/api/item/123"),
    htmx.HxConfirm("Are you sure you want to delete this?"),
    g.Text("Delete"),
)

// Prompt for input
html.Button(
    htmx.HxPost("/api/comment"),
    htmx.HxPrompt("Enter your comment"),
    g.Text("Add Comment"),
)
```

## Extensions

### Loading Extensions

```go
// Single extension
html.Head(
    htmx.Scripts(),
    htmx.ExtensionScript(htmx.ExtensionSSE),
)

// Multiple extensions
html.Head(
    g.Group(htmx.ScriptsWithExtensions(
        htmx.ExtensionSSE,
        htmx.ExtensionWebSockets,
        htmx.ExtensionJSONEnc,
    )),
)
```

### Available Extensions

- `ExtensionSSE` - Server-Sent Events
- `ExtensionWebSockets` - WebSocket support
- `ExtensionClassTools` - Advanced class manipulation
- `ExtensionPreload` - Preload content on hover
- `ExtensionHeadSupport` - Support for head merging
- `ExtensionResponseTargets` - Multiple targets from one response
- `ExtensionDebug` - Debug HTMX requests
- `ExtensionJSONEnc` - JSON encoding for requests
- `ExtensionMethodOverride` - HTTP method override
- `ExtensionMultiSwap` - Multiple swaps in one response

### Using Extensions

```go
html.Div(
    htmx.HxExt("json-enc"),
    htmx.HxPost("/api/data"),
    // Request body will be JSON encoded
)
```

## Advanced Features

### Out-of-Band Swaps

```go
html.Button(
    htmx.HxGet("/data"),
    htmx.HxSelectOOB("#notifications, #messages"),
)
```

### Response Filtering

```go
html.Button(
    htmx.HxGet("/full-page"),
    htmx.HxSelect("#content"),      // Select from response
    htmx.HxTarget("#main"),         // Insert into target
)
```

### Element Preservation

```go
// Preserve element during swaps
html.Video(
    htmx.HxPreserve(),
    html.Src("/video.mp4"),
)
```

### File Uploads

```go
html.Form(
    htmx.HxPost("/upload"),
    htmx.HxEncoding("multipart/form-data"),
    html.Input(html.Type("file"), html.Name("file")),
)
```

### Validation

```go
html.Input(
    html.Type("email"),
    html.Required(),
    htmx.HxValidate(true),  // Force validation before submit
)
```

## CSS Utilities

### Cloak (Hide Until Loaded)

```go
html.Head(
    htmx.CloakCSS(),
)

html.Div(
    htmx.HxCloak(),
    htmx.HxGet("/content"),
    htmx.HxTriggerLoad(),
)
```

### Indicator Styles

```go
html.Head(
    htmx.IndicatorCSS(),
)
```

## Best Practices

### 1. Progressive Enhancement

Start with working server-rendered forms, then enhance with HTMX:

```go
html.Form(
    html.Action("/submit"),
    html.Method("POST"),
    htmx.HxPost("/submit"),  // HTMX enhancement
    htmx.HxTarget("#results"),
    // form fields...
)
```

### 2. Proper Response Types

```go
func handler(w http.ResponseWriter, r *http.Request) {
    if htmx.IsHTMX(r) {
        // Return only what's needed
        w.Write([]byte("<div>Partial</div>"))
    } else {
        // Return full page for non-HTMX browsers
        renderFullPage(w)
    }
}
```

### 3. Loading States

Always provide feedback during requests:

```go
html.Button(
    htmx.HxPost("/submit"),
    htmx.HxIndicator("#spinner"),
    htmx.HxDisabledElt("this"),
    g.Text("Submit"),
)
```

### 4. Error Handling

Handle errors gracefully:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    if err := process(); err != nil {
        htmx.SetHTMXTrigger(w, map[string]any{
            "showError": map[string]string{
                "message": err.Error(),
            },
        })
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // success...
}
```

## Examples

See the [`example/`](../example/) directory for complete working examples including:
- Form submissions
- Infinite scroll
- Search with debouncing
- Modal dialogs
- Real-time updates

## Resources

- [HTMX Documentation](https://htmx.org/docs/)
- [HTMX Examples](https://htmx.org/examples/)
- [ForgeUI Bridge Integration](../bridge/)

## Performance Tips

1. **Use hx-sync** for search inputs to prevent race conditions
2. **Use hx-trigger modifiers** (debounce/throttle) to reduce server requests
3. **Cache responses** server-side for common requests
4. **Use out-of-band swaps** to update multiple elements efficiently
5. **Implement proper loading states** for better UX

