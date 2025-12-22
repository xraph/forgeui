package router

// RouteBuilder provides a fluent API for building routes
type RouteBuilder struct {
	router     *Router
	route      *Route
	registered bool
}

// newBuilder creates a new RouteBuilder
func newBuilder(router *Router, route *Route) *RouteBuilder {
	return &RouteBuilder{
		router:     router,
		route:      route,
		registered: false,
	}
}

// Name sets the name of the route for URL generation
func (b *RouteBuilder) Name(name string) *RouteBuilder {
	b.route.Name = name
	b.router.Name(name, b.route)
	return b
}

// Middleware adds middleware to this specific route
func (b *RouteBuilder) Middleware(middleware ...Middleware) *RouteBuilder {
	b.route.Middleware = append(b.route.Middleware, middleware...)
	return b
}

// Handler sets the handler function (optional, used if handler not set during registration)
func (b *RouteBuilder) Handler(handler PageHandler) *RouteBuilder {
	b.route.Handler = handler
	return b
}

// Build finalizes the route (no-op, for API consistency)
func (b *RouteBuilder) Build() *Route {
	return b.route
}

