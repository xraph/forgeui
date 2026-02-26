package htmx

import (
	"fmt"

	"github.com/a-h/templ"
)

// HxTrigger creates an hx-trigger attribute with a custom event.
//
// Example (in .templ files):
//
//	<div { htmx.HxTrigger("click")... } { htmx.HxGet("/api/data")... }>
func HxTrigger(event string) templ.Attributes {
	return templ.Attributes{"hx-trigger": event}
}

// HxTriggerClick creates an hx-trigger="click" attribute.
func HxTriggerClick() templ.Attributes {
	return templ.Attributes{"hx-trigger": "click"}
}

// HxTriggerChange creates an hx-trigger="change" attribute.
func HxTriggerChange() templ.Attributes {
	return templ.Attributes{"hx-trigger": "change"}
}

// HxTriggerSubmit creates an hx-trigger="submit" attribute.
func HxTriggerSubmit() templ.Attributes {
	return templ.Attributes{"hx-trigger": "submit"}
}

// HxTriggerLoad creates an hx-trigger="load" attribute.
// Triggers when the element is first loaded.
func HxTriggerLoad() templ.Attributes {
	return templ.Attributes{"hx-trigger": "load"}
}

// HxTriggerRevealed creates an hx-trigger="revealed" attribute.
// Triggers when the element is scrolled into the viewport.
func HxTriggerRevealed() templ.Attributes {
	return templ.Attributes{"hx-trigger": "revealed"}
}

// HxTriggerIntersect creates an hx-trigger="intersect" attribute.
// Triggers when the element intersects the viewport.
//
// Example (in .templ files):
//
//	<div { htmx.HxTriggerIntersect("once threshold:0.5")... } { htmx.HxGet("/api/lazy-load")... }>
func HxTriggerIntersect(options string) templ.Attributes {
	if options == "" {
		return templ.Attributes{"hx-trigger": "intersect"}
	}

	return templ.Attributes{"hx-trigger": "intersect " + options}
}

// HxTriggerEvery creates an hx-trigger with polling interval.
//
// Example (in .templ files):
//
//	<div { htmx.HxTriggerEvery("2s")... } { htmx.HxGet("/api/status")... }>
func HxTriggerEvery(duration string) templ.Attributes {
	return templ.Attributes{"hx-trigger": "every " + duration}
}

// HxTriggerMouseEnter creates an hx-trigger="mouseenter" attribute.
func HxTriggerMouseEnter() templ.Attributes {
	return templ.Attributes{"hx-trigger": "mouseenter"}
}

// HxTriggerMouseLeave creates an hx-trigger="mouseleave" attribute.
func HxTriggerMouseLeave() templ.Attributes {
	return templ.Attributes{"hx-trigger": "mouseleave"}
}

// HxTriggerThrottle creates an hx-trigger with throttling.
//
// Example (in .templ files):
//
//	<input { htmx.HxTriggerThrottle("keyup", "1s")... } { htmx.HxGet("/search")... }/>
func HxTriggerThrottle(event, delay string) templ.Attributes {
	return templ.Attributes{"hx-trigger": fmt.Sprintf("%s throttle:%s", event, delay)}
}

// HxTriggerDebounce creates an hx-trigger with debouncing.
//
// Example (in .templ files):
//
//	<input { htmx.HxTriggerDebounce("keyup", "500ms")... } { htmx.HxGet("/search")... }/>
func HxTriggerDebounce(event, delay string) templ.Attributes {
	return templ.Attributes{"hx-trigger": fmt.Sprintf("%s changed delay:%s", event, delay)}
}

// HxTriggerOnce creates an hx-trigger that fires only once.
//
// Example (in .templ files):
//
//	<div { htmx.HxTriggerOnce("click")... } { htmx.HxGet("/api/init")... }>
func HxTriggerOnce(event string) templ.Attributes {
	return templ.Attributes{"hx-trigger": event + " once"}
}

// HxTriggerFrom creates an hx-trigger from another element.
//
// Example (in .templ files):
//
//	<div id="target" { htmx.HxTriggerFrom("click from:#button")... } { htmx.HxGet("/api/data")... }>
func HxTriggerFrom(eventAndSelector string) templ.Attributes {
	return templ.Attributes{"hx-trigger": eventAndSelector}
}

// HxTriggerFilter creates an hx-trigger with an event filter.
//
// Example (in .templ files):
//
//	<input { htmx.HxTriggerFilter("keyup[key=='Enter']")... } { htmx.HxPost("/api/submit")... }/>
func HxTriggerFilter(eventAndFilter string) templ.Attributes {
	return templ.Attributes{"hx-trigger": eventAndFilter}
}

// HxTriggerQueue creates an hx-trigger with queue modifier.
//
// Example (in .templ files):
//
//	<button { htmx.HxTriggerQueue("click", "first")... } { htmx.HxPost("/api/action")... }>
func HxTriggerQueue(event, queueOption string) templ.Attributes {
	return templ.Attributes{"hx-trigger": fmt.Sprintf("%s queue:%s", event, queueOption)}
}

// HxTriggerTarget creates an hx-trigger with target modifier.
//
// Example (in .templ files):
//
//	<div { htmx.HxTriggerTarget("click", "#button")... } { htmx.HxGet("/api/data")... }>
func HxTriggerTarget(event, selector string) templ.Attributes {
	return templ.Attributes{"hx-trigger": fmt.Sprintf("%s target:%s", event, selector)}
}

// HxTriggerConsume creates an hx-trigger with consume modifier.
// Prevents the event from bubbling.
//
// Example (in .templ files):
//
//	<button { htmx.HxTriggerConsume("click")... } { htmx.HxPost("/api/action")... }>
func HxTriggerConsume(event string) templ.Attributes {
	return templ.Attributes{"hx-trigger": event + " consume"}
}
