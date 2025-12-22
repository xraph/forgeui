package htmx

import (
	"encoding/json"
	"fmt"

	g "maragu.dev/gomponents"
)

// HxGet creates an hx-get attribute for GET requests.
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/api/data"),
//	    g.Text("Load Data"),
//	)
func HxGet(url string) g.Node {
	return g.Attr("hx-get", url)
}

// HxPost creates an hx-post attribute for POST requests.
//
// Example:
//
//	html.Form(
//	    htmx.HxPost("/api/submit"),
//	    // form fields...
//	)
func HxPost(url string) g.Node {
	return g.Attr("hx-post", url)
}

// HxPut creates an hx-put attribute for PUT requests.
//
// Example:
//
//	html.Button(
//	    htmx.HxPut("/api/update/123"),
//	    g.Text("Update"),
//	)
func HxPut(url string) g.Node {
	return g.Attr("hx-put", url)
}

// HxPatch creates an hx-patch attribute for PATCH requests.
//
// Example:
//
//	html.Button(
//	    htmx.HxPatch("/api/partial-update/123"),
//	    g.Text("Patch"),
//	)
func HxPatch(url string) g.Node {
	return g.Attr("hx-patch", url)
}

// HxDelete creates an hx-delete attribute for DELETE requests.
//
// Example:
//
//	html.Button(
//	    htmx.HxDelete("/api/delete/123"),
//	    htmx.HxConfirm("Are you sure?"),
//	    g.Text("Delete"),
//	)
func HxDelete(url string) g.Node {
	return g.Attr("hx-delete", url)
}

// HxTarget creates an hx-target attribute to specify where to swap content.
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/api/data"),
//	    htmx.HxTarget("#results"),
//	    g.Text("Load"),
//	)
func HxTarget(selector string) g.Node {
	return g.Attr("hx-target", selector)
}

// HxSwap creates an hx-swap attribute with a swap strategy.
//
// Strategies:
//   - innerHTML (default): Replace the inner html of the target element
//   - outerHTML: Replace the entire target element
//   - beforebegin: Insert before the target element
//   - afterbegin: Insert before the first child of the target element
//   - beforeend: Insert after the last child of the target element
//   - afterend: Insert after the target element
//   - delete: Deletes the target element regardless of the response
//   - none: Does not append content
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/api/item"),
//	    htmx.HxSwap("outerHTML"),
//	    g.Text("Replace"),
//	)
func HxSwap(strategy string) g.Node {
	return g.Attr("hx-swap", strategy)
}

// HxSwapInnerHTML creates an hx-swap="innerHTML" attribute.
func HxSwapInnerHTML() g.Node {
	return g.Attr("hx-swap", "innerHTML")
}

// HxSwapOuterHTML creates an hx-swap="outerHTML" attribute.
func HxSwapOuterHTML() g.Node {
	return g.Attr("hx-swap", "outerHTML")
}

// HxSwapBeforeBegin creates an hx-swap="beforebegin" attribute.
func HxSwapBeforeBegin() g.Node {
	return g.Attr("hx-swap", "beforebegin")
}

// HxSwapAfterBegin creates an hx-swap="afterbegin" attribute.
func HxSwapAfterBegin() g.Node {
	return g.Attr("hx-swap", "afterbegin")
}

// HxSwapBeforeEnd creates an hx-swap="beforeend" attribute.
func HxSwapBeforeEnd() g.Node {
	return g.Attr("hx-swap", "beforeend")
}

// HxSwapAfterEnd creates an hx-swap="afterend" attribute.
func HxSwapAfterEnd() g.Node {
	return g.Attr("hx-swap", "afterend")
}

// HxSwapDelete creates an hx-swap="delete" attribute.
func HxSwapDelete() g.Node {
	return g.Attr("hx-swap", "delete")
}

// HxSwapNone creates an hx-swap="none" attribute.
func HxSwapNone() g.Node {
	return g.Attr("hx-swap", "none")
}

// HxIndicator creates an hx-indicator attribute to specify a loading indicator.
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/api/slow"),
//	    htmx.HxIndicator("#spinner"),
//	    g.Text("Load"),
//	)
//	html.Div(
//	    html.ID("spinner"),
//	    html.Class("htmx-indicator"),
//	    g.Text("Loading..."),
//	)
func HxIndicator(selector string) g.Node {
	return g.Attr("hx-indicator", selector)
}

// HxDisabledElt creates an hx-disabled-elt attribute to disable elements during requests.
//
// Example:
//
//	html.Button(
//	    htmx.HxPost("/api/submit"),
//	    htmx.HxDisabledElt("this"),
//	    g.Text("Submit"),
//	)
func HxDisabledElt(selector string) g.Node {
	return g.Attr("hx-disabled-elt", selector)
}

// HxSync creates an hx-sync attribute for request synchronization.
//
// Strategies:
//   - drop: Drop (ignore) the request if another is in flight
//   - abort: Abort the current request if a new one is triggered
//   - replace: Replace the current request if a new one is triggered
//   - queue: Queue requests
//   - queue first: Queue requests, but execute the first immediately
//   - queue last: Queue requests, but execute the last immediately
//   - queue all: Queue all requests
//
// Example:
//
//	html.Input(
//	    htmx.HxGet("/search"),
//	    htmx.HxSync("this", "replace"),
//	    html.Type("text"),
//	)
func HxSync(selector, strategy string) g.Node {
	return g.Attr("hx-sync", fmt.Sprintf("%s:%s", selector, strategy))
}

// HxBoost creates an hx-boost attribute for progressive enhancement of links/forms.
//
// Example:
//
//	html.Div(
//	    htmx.HxBoost(true),
//	    html.A(html.Href("/page1"), g.Text("Page 1")),
//	    html.A(html.Href("/page2"), g.Text("Page 2")),
//	)
func HxBoost(enabled bool) g.Node {
	if enabled {
		return g.Attr("hx-boost", "true")
	}

	return g.Attr("hx-boost", "false")
}

// HxPushURL creates an hx-push-url attribute for history management.
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/page/2"),
//	    htmx.HxPushURL(true),
//	    g.Text("Next Page"),
//	)
func HxPushURL(enabled bool) g.Node {
	if enabled {
		return g.Attr("hx-push-url", "true")
	}

	return g.Attr("hx-push-url", "false")
}

// HxPushURLWithPath creates an hx-push-url attribute with a custom path.
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/api/page/2"),
//	    htmx.HxPushURLWithPath("/page/2"),
//	    g.Text("Next Page"),
//	)
func HxPushURLWithPath(path string) g.Node {
	return g.Attr("hx-push-url", path)
}

// HxReplaceURL creates an hx-replace-url attribute to replace the current URL.
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/search?q=test"),
//	    htmx.HxReplaceURL(true),
//	    g.Text("Search"),
//	)
func HxReplaceURL(enabled bool) g.Node {
	if enabled {
		return g.Attr("hx-replace-url", "true")
	}

	return g.Attr("hx-replace-url", "false")
}

// HxReplaceURLWithPath creates an hx-replace-url attribute with a custom path.
func HxReplaceURLWithPath(path string) g.Node {
	return g.Attr("hx-replace-url", path)
}

// HxSelect creates an hx-select attribute for response filtering.
//
// Example:
//
//	html.Button(
//	    htmx.HxGet("/full-page"),
//	    htmx.HxSelect("#content"),
//	    htmx.HxTarget("#main"),
//	    g.Text("Load Content"),
//	)
func HxSelect(selector string) g.Node {
	return g.Attr("hx-select", selector)
}

// HxSelectOOB creates an hx-select-oob attribute for out-of-band swapping.
//
// Example:
//
//	htmx.HxSelectOOB("#notifications, #messages")
func HxSelectOOB(selector string) g.Node {
	return g.Attr("hx-select-oob", selector)
}

// HxHeaders creates an hx-headers attribute for custom headers.
//
// Example:
//
//	html.Button(
//	    htmx.HxPost("/api/data"),
//	    htmx.HxHeaders(map[string]string{
//	        "X-API-Key": "secret",
//	    }),
//	    g.Text("Submit"),
//	)
func HxHeaders(headers map[string]string) g.Node {
	jsonData, err := json.Marshal(headers)
	if err != nil {
		// Fallback to empty object on marshal error
		jsonData = []byte("{}")
	}

	return g.Attr("hx-headers", string(jsonData))
}

// HxVals creates an hx-vals attribute for extra values to submit.
//
// Example:
//
//	html.Button(
//	    htmx.HxPost("/api/data"),
//	    htmx.HxVals(map[string]any{
//	        "category": "urgent",
//	        "priority": 1,
//	    }),
//	    g.Text("Submit"),
//	)
func HxVals(values map[string]any) g.Node {
	jsonData, err := json.Marshal(values)
	if err != nil {
		// Fallback to empty object on marshal error
		jsonData = []byte("{}")
	}

	return g.Attr("hx-vals", string(jsonData))
}

// HxValsJS creates an hx-vals attribute with JavaScript evaluation.
//
// Example:
//
//	htmx.HxValsJS("js:{timestamp: Date.now()}")
func HxValsJS(jsExpr string) g.Node {
	return g.Attr("hx-vals", jsExpr)
}

// HxConfirm creates an hx-confirm attribute for confirmation dialogs.
//
// Example:
//
//	html.Button(
//	    htmx.HxDelete("/api/item/123"),
//	    htmx.HxConfirm("Are you sure you want to delete this item?"),
//	    g.Text("Delete"),
//	)
func HxConfirm(message string) g.Node {
	return g.Attr("hx-confirm", message)
}

// HxPrompt creates an hx-prompt attribute for user input.
//
// Example:
//
//	html.Button(
//	    htmx.HxPost("/api/comment"),
//	    htmx.HxPrompt("Enter your comment"),
//	    g.Text("Add Comment"),
//	)
func HxPrompt(message string) g.Node {
	return g.Attr("hx-prompt", message)
}

// HxParams creates an hx-params attribute to filter parameters.
//
// Values:
//   - "*": Include all parameters (default)
//   - "none": Include no parameters
//   - "param1,param2": Include only specified parameters
//   - "not param1,param2": Include all except specified parameters
//
// Example:
//
//	html.Form(
//	    htmx.HxPost("/api/submit"),
//	    htmx.HxParams("email,name"),
//	    // form fields...
//	)
func HxParams(params string) g.Node {
	return g.Attr("hx-params", params)
}

// HxExt creates an hx-ext attribute to load HTMX extensions.
//
// Example:
//
//	html.Div(
//	    htmx.HxExt("json-enc"),
//	    htmx.HxPost("/api/data"),
//	)
func HxExt(extension string) g.Node {
	return g.Attr("hx-ext", extension)
}

// HxInclude creates an hx-include attribute to include additional elements in requests.
//
// Example:
//
//	html.Button(
//	    htmx.HxPost("/api/submit"),
//	    htmx.HxInclude("#extra-data"),
//	    g.Text("Submit"),
//	)
func HxInclude(selector string) g.Node {
	return g.Attr("hx-include", selector)
}

// HxPreserve creates an hx-preserve attribute to preserve elements across requests.
//
// Example:
//
//	html.Video(
//	    htmx.HxPreserve(true),
//	    html.Src("/video.mp4"),
//	)
func HxPreserve(enabled bool) g.Node {
	if enabled {
		return g.Attr("hx-preserve", "true")
	}

	return g.Attr("hx-preserve", "false")
}

// HxEncoding creates an hx-encoding attribute to set request encoding.
//
// Example:
//
//	html.Form(
//	    htmx.HxPost("/api/upload"),
//	    htmx.HxEncoding("multipart/form-data"),
//	    // file input...
//	)
func HxEncoding(encoding string) g.Node {
	return g.Attr("hx-encoding", encoding)
}

// HxValidate creates an hx-validate attribute to force validation before submit.
//
// Example:
//
//	html.Input(
//	    html.Type("email"),
//	    html.Required(),
//	    htmx.HxValidate(true),
//	)
func HxValidate(enabled bool) g.Node {
	if enabled {
		return g.Attr("hx-validate", "true")
	}

	return g.Attr("hx-validate", "false")
}
