package forgeui

import (
	"github.com/xraph/forgeui/router"
)

// PageBuilder provides a fluent API for building and registering pages
type PageBuilder struct {
	app        *EnhancedApp
	pattern    string
	method     string
	handler    router.PageHandler
	loader     router.LoaderFunc
	layout     string
	middleware []router.Middleware
	meta       *router.RouteMeta
	name       string
	noLayout   bool
}

// Handler sets the page handler function
func (pb *PageBuilder) Handler(handler router.PageHandler) *PageBuilder {
	pb.handler = handler
	return pb
}

// Loader sets the data loader function
func (pb *PageBuilder) Loader(loader router.LoaderFunc) *PageBuilder {
	pb.loader = loader
	return pb
}

// Layout sets the layout name for this page
func (pb *PageBuilder) Layout(layout string) *PageBuilder {
	pb.layout = layout
	pb.noLayout = false

	return pb
}

// NoLayout disables layout for this page
func (pb *PageBuilder) NoLayout() *PageBuilder {
	pb.noLayout = true
	pb.layout = "none"

	return pb
}

// Meta sets the page metadata (title and description)
func (pb *PageBuilder) Meta(title, description string) *PageBuilder {
	if pb.meta == nil {
		pb.meta = &router.RouteMeta{}
	}

	pb.meta.Title = title
	pb.meta.Description = description

	return pb
}

// MetaFull sets complete page metadata
func (pb *PageBuilder) MetaFull(meta *router.RouteMeta) *PageBuilder {
	pb.meta = meta
	return pb
}

// Middleware adds middleware to this page
func (pb *PageBuilder) Middleware(middleware ...router.Middleware) *PageBuilder {
	pb.middleware = append(pb.middleware, middleware...)
	return pb
}

// Name sets a name for this route (for URL generation)
func (pb *PageBuilder) Name(name string) *PageBuilder {
	pb.name = name
	return pb
}

// Method sets the HTTP method for this page
func (pb *PageBuilder) Method(method string) *PageBuilder {
	pb.method = method
	return pb
}

// Register registers the page with the router
func (pb *PageBuilder) Register() *router.Route {
	if pb.handler == nil {
		panic("page handler is required")
	}

	// Create route based on method
	var route *router.Route

	switch pb.method {
	case "GET":
		route = pb.app.router.Get(pb.pattern, pb.handler)
	case "POST":
		route = pb.app.router.Post(pb.pattern, pb.handler)
	case "PUT":
		route = pb.app.router.Put(pb.pattern, pb.handler)
	case "PATCH":
		route = pb.app.router.Patch(pb.pattern, pb.handler)
	case "DELETE":
		route = pb.app.router.Delete(pb.pattern, pb.handler)
	default:
		route = pb.app.router.Get(pb.pattern, pb.handler)
	}

	// Apply loader if set
	if pb.loader != nil {
		route.WithLoader(pb.loader)
	}

	// Apply layout if set
	if pb.noLayout {
		route.NoLayout()
	} else if pb.layout != "" {
		route.SetLayout(pb.layout)
	}

	// Apply metadata if set
	if pb.meta != nil {
		route.WithMeta(pb.meta)
	}

	// Apply middleware if set
	if len(pb.middleware) > 0 {
		route.WithMiddleware(pb.middleware...)
	}

	// Apply name if set
	if pb.name != "" {
		pb.app.router.Name(pb.name, route)
	}

	return route
}

// GET is a convenience method to set method to GET
func (pb *PageBuilder) GET() *PageBuilder {
	pb.method = "GET"
	return pb
}

// POST is a convenience method to set method to POST
func (pb *PageBuilder) POST() *PageBuilder {
	pb.method = "POST"
	return pb
}

// PUT is a convenience method to set method to PUT
func (pb *PageBuilder) PUT() *PageBuilder {
	pb.method = "PUT"
	return pb
}

// PATCH is a convenience method to set method to PATCH
func (pb *PageBuilder) PATCH() *PageBuilder {
	pb.method = "PATCH"
	return pb
}

// DELETE is a convenience method to set method to DELETE
func (pb *PageBuilder) DELETE() *PageBuilder {
	pb.method = "DELETE"
	return pb
}
