package htmx

import (
	"encoding/json"
	"fmt"

	"github.com/a-h/templ"
)

// HxGet creates an hx-get attribute for GET requests.
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/api/data")... }>Load Data</button>
func HxGet(url string) templ.Attributes {
	return templ.Attributes{"hx-get": url}
}

// HxPost creates an hx-post attribute for POST requests.
//
// Example (in .templ files):
//
//	<form { htmx.HxPost("/api/submit")... }>
func HxPost(url string) templ.Attributes {
	return templ.Attributes{"hx-post": url}
}

// HxPut creates an hx-put attribute for PUT requests.
//
// Example (in .templ files):
//
//	<button { htmx.HxPut("/api/update/123")... }>Update</button>
func HxPut(url string) templ.Attributes {
	return templ.Attributes{"hx-put": url}
}

// HxPatch creates an hx-patch attribute for PATCH requests.
//
// Example (in .templ files):
//
//	<button { htmx.HxPatch("/api/partial-update/123")... }>Patch</button>
func HxPatch(url string) templ.Attributes {
	return templ.Attributes{"hx-patch": url}
}

// HxDelete creates an hx-delete attribute for DELETE requests.
//
// Example (in .templ files):
//
//	<button { htmx.HxDelete("/api/delete/123")... } { htmx.HxConfirm("Are you sure?")... }>Delete</button>
func HxDelete(url string) templ.Attributes {
	return templ.Attributes{"hx-delete": url}
}

// HxTarget creates an hx-target attribute to specify where to swap content.
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/api/data")... } { htmx.HxTarget("#results")... }>Load</button>
func HxTarget(selector string) templ.Attributes {
	return templ.Attributes{"hx-target": selector}
}

// HxSwap creates an hx-swap attribute with a swap strategy.
//
// Strategies: innerHTML, outerHTML, beforebegin, afterbegin, beforeend, afterend, delete, none
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/api/item")... } { htmx.HxSwap("outerHTML")... }>Replace</button>
func HxSwap(strategy string) templ.Attributes {
	return templ.Attributes{"hx-swap": strategy}
}

// HxSwapInnerHTML creates an hx-swap="innerHTML" attribute.
func HxSwapInnerHTML() templ.Attributes {
	return templ.Attributes{"hx-swap": "innerHTML"}
}

// HxSwapOuterHTML creates an hx-swap="outerHTML" attribute.
func HxSwapOuterHTML() templ.Attributes {
	return templ.Attributes{"hx-swap": "outerHTML"}
}

// HxSwapBeforeBegin creates an hx-swap="beforebegin" attribute.
func HxSwapBeforeBegin() templ.Attributes {
	return templ.Attributes{"hx-swap": "beforebegin"}
}

// HxSwapAfterBegin creates an hx-swap="afterbegin" attribute.
func HxSwapAfterBegin() templ.Attributes {
	return templ.Attributes{"hx-swap": "afterbegin"}
}

// HxSwapBeforeEnd creates an hx-swap="beforeend" attribute.
func HxSwapBeforeEnd() templ.Attributes {
	return templ.Attributes{"hx-swap": "beforeend"}
}

// HxSwapAfterEnd creates an hx-swap="afterend" attribute.
func HxSwapAfterEnd() templ.Attributes {
	return templ.Attributes{"hx-swap": "afterend"}
}

// HxSwapDelete creates an hx-swap="delete" attribute.
func HxSwapDelete() templ.Attributes {
	return templ.Attributes{"hx-swap": "delete"}
}

// HxSwapNone creates an hx-swap="none" attribute.
func HxSwapNone() templ.Attributes {
	return templ.Attributes{"hx-swap": "none"}
}

// HxIndicator creates an hx-indicator attribute to specify a loading indicator.
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/api/slow")... } { htmx.HxIndicator("#spinner")... }>Load</button>
func HxIndicator(selector string) templ.Attributes {
	return templ.Attributes{"hx-indicator": selector}
}

// HxDisabledElt creates an hx-disabled-elt attribute to disable elements during requests.
//
// Example (in .templ files):
//
//	<button { htmx.HxPost("/api/submit")... } { htmx.HxDisabledElt("this")... }>Submit</button>
func HxDisabledElt(selector string) templ.Attributes {
	return templ.Attributes{"hx-disabled-elt": selector}
}

// HxSync creates an hx-sync attribute for request synchronization.
//
// Strategies: drop, abort, replace, queue, queue first, queue last, queue all
//
// Example (in .templ files):
//
//	<input { htmx.HxGet("/search")... } { htmx.HxSync("this", "replace")... } type="text"/>
func HxSync(selector, strategy string) templ.Attributes {
	return templ.Attributes{"hx-sync": fmt.Sprintf("%s:%s", selector, strategy)}
}

// HxBoost creates an hx-boost attribute for progressive enhancement of links/forms.
//
// Example (in .templ files):
//
//	<div { htmx.HxBoost(true)... }>
func HxBoost(enabled bool) templ.Attributes {
	if enabled {
		return templ.Attributes{"hx-boost": "true"}
	}

	return templ.Attributes{"hx-boost": "false"}
}

// HxPushURL creates an hx-push-url attribute for history management.
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/page/2")... } { htmx.HxPushURL(true)... }>Next Page</button>
func HxPushURL(enabled bool) templ.Attributes {
	if enabled {
		return templ.Attributes{"hx-push-url": "true"}
	}

	return templ.Attributes{"hx-push-url": "false"}
}

// HxPushURLWithPath creates an hx-push-url attribute with a custom path.
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/api/page/2")... } { htmx.HxPushURLWithPath("/page/2")... }>Next</button>
func HxPushURLWithPath(path string) templ.Attributes {
	return templ.Attributes{"hx-push-url": path}
}

// HxReplaceURL creates an hx-replace-url attribute to replace the current URL.
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/search?q=test")... } { htmx.HxReplaceURL(true)... }>Search</button>
func HxReplaceURL(enabled bool) templ.Attributes {
	if enabled {
		return templ.Attributes{"hx-replace-url": "true"}
	}

	return templ.Attributes{"hx-replace-url": "false"}
}

// HxReplaceURLWithPath creates an hx-replace-url attribute with a custom path.
func HxReplaceURLWithPath(path string) templ.Attributes {
	return templ.Attributes{"hx-replace-url": path}
}

// HxSelect creates an hx-select attribute for response filtering.
//
// Example (in .templ files):
//
//	<button { htmx.HxGet("/full-page")... } { htmx.HxSelect("#content")... }>Load Content</button>
func HxSelect(selector string) templ.Attributes {
	return templ.Attributes{"hx-select": selector}
}

// HxSelectOOB creates an hx-select-oob attribute for out-of-band swapping.
//
// Example:
//
//	htmx.HxSelectOOB("#notifications, #messages")
func HxSelectOOB(selector string) templ.Attributes {
	return templ.Attributes{"hx-select-oob": selector}
}

// HxHeaders creates an hx-headers attribute for custom headers.
//
// Example (in .templ files):
//
//	<button { htmx.HxPost("/api/data")... } { htmx.HxHeaders(map[string]string{"X-API-Key": "secret"})... }>
func HxHeaders(headers map[string]string) templ.Attributes {
	jsonData, err := json.Marshal(headers)
	if err != nil {
		jsonData = []byte("{}")
	}

	return templ.Attributes{"hx-headers": string(jsonData)}
}

// HxVals creates an hx-vals attribute for extra values to submit.
//
// Example (in .templ files):
//
//	<button { htmx.HxPost("/api/data")... } { htmx.HxVals(map[string]any{"category": "urgent"})... }>
func HxVals(values map[string]any) templ.Attributes {
	jsonData, err := json.Marshal(values)
	if err != nil {
		jsonData = []byte("{}")
	}

	return templ.Attributes{"hx-vals": string(jsonData)}
}

// HxValsJS creates an hx-vals attribute with JavaScript evaluation.
//
// Example:
//
//	htmx.HxValsJS("js:{timestamp: Date.now()}")
func HxValsJS(jsExpr string) templ.Attributes {
	return templ.Attributes{"hx-vals": jsExpr}
}

// HxConfirm creates an hx-confirm attribute for confirmation dialogs.
//
// Example (in .templ files):
//
//	<button { htmx.HxDelete("/api/item/123")... } { htmx.HxConfirm("Are you sure?")... }>Delete</button>
func HxConfirm(message string) templ.Attributes {
	return templ.Attributes{"hx-confirm": message}
}

// HxPrompt creates an hx-prompt attribute for user input.
//
// Example (in .templ files):
//
//	<button { htmx.HxPost("/api/comment")... } { htmx.HxPrompt("Enter your comment")... }>Add</button>
func HxPrompt(message string) templ.Attributes {
	return templ.Attributes{"hx-prompt": message}
}

// HxParams creates an hx-params attribute to filter parameters.
//
// Values: "*" (all), "none", "param1,param2" (include), "not param1,param2" (exclude)
//
// Example (in .templ files):
//
//	<form { htmx.HxPost("/api/submit")... } { htmx.HxParams("email,name")... }>
func HxParams(params string) templ.Attributes {
	return templ.Attributes{"hx-params": params}
}

// HxExt creates an hx-ext attribute to load HTMX extensions.
//
// Example (in .templ files):
//
//	<div { htmx.HxExt("json-enc")... } { htmx.HxPost("/api/data")... }>
func HxExt(extension string) templ.Attributes {
	return templ.Attributes{"hx-ext": extension}
}

// HxInclude creates an hx-include attribute to include additional elements in requests.
//
// Example (in .templ files):
//
//	<button { htmx.HxPost("/api/submit")... } { htmx.HxInclude("#extra-data")... }>Submit</button>
func HxInclude(selector string) templ.Attributes {
	return templ.Attributes{"hx-include": selector}
}

// HxPreserve creates an hx-preserve attribute to preserve elements across requests.
//
// Example (in .templ files):
//
//	<video { htmx.HxPreserve(true)... } src="/video.mp4"></video>
func HxPreserve(enabled bool) templ.Attributes {
	if enabled {
		return templ.Attributes{"hx-preserve": "true"}
	}

	return templ.Attributes{"hx-preserve": "false"}
}

// HxEncoding creates an hx-encoding attribute to set request encoding.
//
// Example (in .templ files):
//
//	<form { htmx.HxPost("/api/upload")... } { htmx.HxEncoding("multipart/form-data")... }>
func HxEncoding(encoding string) templ.Attributes {
	return templ.Attributes{"hx-encoding": encoding}
}

// HxValidate creates an hx-validate attribute to force validation before submit.
//
// Example (in .templ files):
//
//	<input type="email" required { htmx.HxValidate(true)... }/>
func HxValidate(enabled bool) templ.Attributes {
	if enabled {
		return templ.Attributes{"hx-validate": "true"}
	}

	return templ.Attributes{"hx-validate": "false"}
}
