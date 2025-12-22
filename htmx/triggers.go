package htmx

import (
	"fmt"

	g "maragu.dev/gomponents"
)

// HxTrigger creates an hx-trigger attribute with a custom event.
//
// Example:
//
//	html.Div(
//	    htmx.HxTrigger("click"),
//	    htmx.HxGet("/api/data"),
//	)
func HxTrigger(event string) g.Node {
	return g.Attr("hx-trigger", event)
}

// HxTriggerClick creates an hx-trigger="click" attribute.
func HxTriggerClick() g.Node {
	return g.Attr("hx-trigger", "click")
}

// HxTriggerChange creates an hx-trigger="change" attribute.
func HxTriggerChange() g.Node {
	return g.Attr("hx-trigger", "change")
}

// HxTriggerSubmit creates an hx-trigger="submit" attribute.
func HxTriggerSubmit() g.Node {
	return g.Attr("hx-trigger", "submit")
}

// HxTriggerLoad creates an hx-trigger="load" attribute.
// Triggers when the element is first loaded.
func HxTriggerLoad() g.Node {
	return g.Attr("hx-trigger", "load")
}

// HxTriggerRevealed creates an hx-trigger="revealed" attribute.
// Triggers when the element is scrolled into the viewport.
func HxTriggerRevealed() g.Node {
	return g.Attr("hx-trigger", "revealed")
}

// HxTriggerIntersect creates an hx-trigger="intersect" attribute.
// Triggers when the element intersects the viewport.
//
// Example:
//
//	html.Div(
//	    htmx.HxTriggerIntersect("once threshold:0.5"),
//	    htmx.HxGet("/api/lazy-load"),
//	)
func HxTriggerIntersect(options string) g.Node {
	if options == "" {
		return g.Attr("hx-trigger", "intersect")
	}
	return g.Attr("hx-trigger", fmt.Sprintf("intersect %s", options))
}

// HxTriggerEvery creates an hx-trigger with polling interval.
//
// Example:
//
//	html.Div(
//	    htmx.HxTriggerEvery("2s"),
//	    htmx.HxGet("/api/status"),
//	)
func HxTriggerEvery(duration string) g.Node {
	return g.Attr("hx-trigger", fmt.Sprintf("every %s", duration))
}

// HxTriggerMouseEnter creates an hx-trigger="mouseenter" attribute.
func HxTriggerMouseEnter() g.Node {
	return g.Attr("hx-trigger", "mouseenter")
}

// HxTriggerMouseLeave creates an hx-trigger="mouseleave" attribute.
func HxTriggerMouseLeave() g.Node {
	return g.Attr("hx-trigger", "mouseleave")
}

// HxTriggerThrottle creates an hx-trigger with throttling.
//
// Example:
//
//	html.Input(
//	    htmx.HxTriggerThrottle("keyup", "1s"),
//	    htmx.HxGet("/search"),
//	)
func HxTriggerThrottle(event, delay string) g.Node {
	return g.Attr("hx-trigger", fmt.Sprintf("%s throttle:%s", event, delay))
}

// HxTriggerDebounce creates an hx-trigger with debouncing.
//
// Example:
//
//	html.Input(
//	    htmx.HxTriggerDebounce("keyup", "500ms"),
//	    htmx.HxGet("/search"),
//	)
func HxTriggerDebounce(event, delay string) g.Node {
	return g.Attr("hx-trigger", fmt.Sprintf("%s changed delay:%s", event, delay))
}

// HxTriggerOnce creates an hx-trigger that fires only once.
//
// Example:
//
//	html.Div(
//	    htmx.HxTriggerOnce("click"),
//	    htmx.HxGet("/api/init"),
//	)
func HxTriggerOnce(event string) g.Node {
	return g.Attr("hx-trigger", fmt.Sprintf("%s once", event))
}

// HxTriggerFrom creates an hx-trigger from another element.
//
// Example:
//
//	html.Div(
//	    html.ID("target"),
//	    htmx.HxTriggerFrom("click from:#button"),
//	    htmx.HxGet("/api/data"),
//	)
func HxTriggerFrom(eventAndSelector string) g.Node {
	return g.Attr("hx-trigger", eventAndSelector)
}

// HxTriggerFilter creates an hx-trigger with an event filter.
//
// Example:
//
//	html.Input(
//	    htmx.HxTriggerFilter("keyup[key=='Enter']"),
//	    htmx.HxPost("/api/submit"),
//	)
func HxTriggerFilter(eventAndFilter string) g.Node {
	return g.Attr("hx-trigger", eventAndFilter)
}

// HxTriggerQueue creates an hx-trigger with queue modifier.
//
// Example:
//
//	html.Button(
//	    htmx.HxTriggerQueue("click", "first"),
//	    htmx.HxPost("/api/action"),
//	)
func HxTriggerQueue(event, queueOption string) g.Node {
	return g.Attr("hx-trigger", fmt.Sprintf("%s queue:%s", event, queueOption))
}

// HxTriggerTarget creates an hx-trigger with target modifier.
//
// Example:
//
//	html.Div(
//	    htmx.HxTriggerTarget("click", "#button"),
//	    htmx.HxGet("/api/data"),
//	)
func HxTriggerTarget(event, selector string) g.Node {
	return g.Attr("hx-trigger", fmt.Sprintf("%s target:%s", event, selector))
}

// HxTriggerConsume creates an hx-trigger with consume modifier.
// Prevents the event from bubbling.
//
// Example:
//
//	html.Button(
//	    htmx.HxTriggerConsume("click"),
//	    htmx.HxPost("/api/action"),
//	)
func HxTriggerConsume(event string) g.Node {
	return g.Attr("hx-trigger", fmt.Sprintf("%s consume", event))
}

