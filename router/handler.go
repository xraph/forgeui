package router

import (
	"net/http"

	"github.com/a-h/templ"
)

// PageHandler is a function that handles a page request and returns a templ.Component for rendering.
type PageHandler func(*PageContext) (templ.Component, error)

// ErrorHandler handles errors that occur during request processing.
type ErrorHandler func(*PageContext, error) templ.Component

// HandlerFunc converts a standard http.HandlerFunc to a PageHandler.
// This is useful for integrating with existing HTTP handlers that don't return components.
func HandlerFunc(fn http.HandlerFunc) PageHandler {
	return func(ctx *PageContext) (templ.Component, error) {
		fn(ctx.ResponseWriter, ctx.Request)
		return nil, nil
	}
}
