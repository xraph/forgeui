package router

import (
	"net/http"

	g "github.com/maragudk/gomponents"
)

// PageHandler is a function that handles a page request and returns a Node for rendering
type PageHandler func(*PageContext) (g.Node, error)

// ErrorHandler handles errors that occur during request processing
type ErrorHandler func(*PageContext, error) g.Node

// HandlerFunc converts a standard http.HandlerFunc to a PageHandler
// This is useful for integrating with existing HTTP handlers that don't return nodes
func HandlerFunc(fn http.HandlerFunc) PageHandler {
	return func(ctx *PageContext) (g.Node, error) {
		fn(ctx.ResponseWriter, ctx.Request)
		return nil, nil
	}
}

