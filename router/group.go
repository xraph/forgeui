package router

// Group represents a group of routes with shared configuration
// Renamed from RouteGroup for cleaner API
type Group struct {
	router     *Router
	prefix     string
	layout     string
	middleware []Middleware
}

// RouteGroup is an alias for Group (backward compatibility)
type RouteGroup = Group

// GroupOption configures a Group
type GroupOption func(*Group)

// GroupLayout sets the default layout for all routes in the group
func GroupLayout(layout string) GroupOption {
	return func(g *Group) {
		g.layout = layout
	}
}

// GroupMiddleware adds middleware to all routes in the group
func GroupMiddleware(middleware ...Middleware) GroupOption {
	return func(g *Group) {
		g.middleware = append(g.middleware, middleware...)
	}
}

// Group creates a new route group with a common prefix and options
func (r *Router) Group(prefix string, opts ...GroupOption) *Group {
	g := &Group{
		router:     r,
		prefix:     prefix,
		middleware: make([]Middleware, 0),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Layout sets the layout for this group (fluent API)
func (g *Group) Layout(layout string) *Group {
	g.layout = layout
	return g
}

// Middleware adds middleware to this group (fluent API)
func (g *Group) Middleware(middleware ...Middleware) *Group {
	g.middleware = append(g.middleware, middleware...)
	return g
}

// Get registers a GET route in the group
func (g *Group) Get(pattern string, handler PageHandler) *Route {
	route := g.router.Get(g.prefix+pattern, handler)
	g.applyGroupConfig(route)
	return route
}

// Post registers a POST route in the group
func (g *Group) Post(pattern string, handler PageHandler) *Route {
	route := g.router.Post(g.prefix+pattern, handler)
	g.applyGroupConfig(route)
	return route
}

// Put registers a PUT route in the group
func (g *Group) Put(pattern string, handler PageHandler) *Route {
	route := g.router.Put(g.prefix+pattern, handler)
	g.applyGroupConfig(route)
	return route
}

// Patch registers a PATCH route in the group
func (g *Group) Patch(pattern string, handler PageHandler) *Route {
	route := g.router.Patch(g.prefix+pattern, handler)
	g.applyGroupConfig(route)
	return route
}

// Delete registers a DELETE route in the group
func (g *Group) Delete(pattern string, handler PageHandler) *Route {
	route := g.router.Delete(g.prefix+pattern, handler)
	g.applyGroupConfig(route)
	return route
}

// Options registers an OPTIONS route in the group
func (g *Group) Options(pattern string, handler PageHandler) *Route {
	route := g.router.Options(g.prefix+pattern, handler)
	g.applyGroupConfig(route)
	return route
}

// Head registers a HEAD route in the group
func (g *Group) Head(pattern string, handler PageHandler) *Route {
	route := g.router.Head(g.prefix+pattern, handler)
	g.applyGroupConfig(route)
	return route
}

// Group creates a nested route group
func (g *Group) Group(prefix string, opts ...GroupOption) *Group {
	nested := &Group{
		router:     g.router,
		prefix:     g.prefix + prefix,
		layout:     g.layout, // Inherit parent layout
		middleware: make([]Middleware, len(g.middleware)),
	}

	// Copy parent middleware
	copy(nested.middleware, g.middleware)

	// Apply nested group options
	for _, opt := range opts {
		opt(nested)
	}

	return nested
}

// applyGroupConfig applies group configuration to a route
func (g *Group) applyGroupConfig(route *Route) {
	// Apply group middleware (prepend to route middleware)
	if len(g.middleware) > 0 {
		route.Middleware = append(g.middleware, route.Middleware...)
	}

	// Apply group layout if set and route doesn't have one
	if g.layout != "" && route.Layout == "" {
		route.Layout = g.layout
	}
}
