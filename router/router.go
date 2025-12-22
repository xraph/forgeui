package router

import (
	"net/http"
	"sort"
	"sync"

	g "github.com/maragudk/gomponents"
)

// Router handles HTTP routing for ForgeUI applications
type Router struct {
	mu            sync.RWMutex
	routes        []*Route
	namedRoutes   map[string]*Route
	notFound      PageHandler
	errorHandler  ErrorHandler
	middleware    []Middleware
	basePath      string
	layouts       map[string]LayoutFunc
	layoutConfigs map[string]*LayoutConfig
	errorPages    map[int]PageHandler
	defaultLayout string
	app           interface{} // Reference to EnhancedApp (interface to avoid circular dependency)
}

// RouterOption is a functional option for configuring the Router
type RouterOption func(*Router)

// New creates a new router with optional configuration
func New(opts ...RouterOption) *Router {
	r := &Router{
		routes:       make([]*Route, 0),
		namedRoutes:  make(map[string]*Route),
		middleware:   make([]Middleware, 0),
		notFound:     defaultNotFound,
		errorHandler: defaultError,
		layouts:      make(map[string]LayoutFunc),
		errorPages:   make(map[int]PageHandler),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// WithBasePath sets a base path for all routes
func WithBasePath(path string) RouterOption {
	return func(r *Router) {
		r.basePath = path
	}
}

// WithNotFound sets a custom 404 handler
func WithNotFound(handler PageHandler) RouterOption {
	return func(r *Router) {
		r.notFound = handler
	}
}

// WithErrorHandler sets a custom error handler
func WithErrorHandler(handler ErrorHandler) RouterOption {
	return func(r *Router) {
		r.errorHandler = handler
	}
}

// Use adds global middleware to the router
func (r *Router) Use(middleware ...Middleware) *Router {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middleware = append(r.middleware, middleware...)
	return r
}

// Handle registers a route with the given method, pattern, and handler
func (r *Router) Handle(method, pattern string, handler PageHandler) *Route {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Prepend base path if set
	if r.basePath != "" {
		pattern = r.basePath + pattern
	}

	route := newRoute(pattern, method, handler)
	r.routes = append(r.routes, route)

	// Sort routes by priority (lower priority number = higher precedence)
	sort.Slice(r.routes, func(i, j int) bool {
		return r.routes[i].priority < r.routes[j].priority
	})

	return route
}

// Get registers a GET route
func (r *Router) Get(pattern string, handler PageHandler) *Route {
	return r.Handle(MethodGet, pattern, handler)
}

// Post registers a POST route
func (r *Router) Post(pattern string, handler PageHandler) *Route {
	return r.Handle(MethodPost, pattern, handler)
}

// Put registers a PUT route
func (r *Router) Put(pattern string, handler PageHandler) *Route {
	return r.Handle(MethodPut, pattern, handler)
}

// Patch registers a PATCH route
func (r *Router) Patch(pattern string, handler PageHandler) *Route {
	return r.Handle(MethodPatch, pattern, handler)
}

// Delete registers a DELETE route
func (r *Router) Delete(pattern string, handler PageHandler) *Route {
	return r.Handle(MethodDelete, pattern, handler)
}

// Options registers an OPTIONS route
func (r *Router) Options(pattern string, handler PageHandler) *Route {
	return r.Handle(MethodOptions, pattern, handler)
}

// Head registers a HEAD route
func (r *Router) Head(pattern string, handler PageHandler) *Route {
	return r.Handle(MethodHead, pattern, handler)
}

// Match registers a route that matches multiple HTTP methods
func (r *Router) Match(methods []string, pattern string, handler PageHandler) []*Route {
	routes := make([]*Route, 0, len(methods))
	for _, method := range methods {
		route := r.Handle(method, pattern, handler)
		routes = append(routes, route)
	}
	return routes
}

// Name registers a named route for URL generation
func (r *Router) Name(name string, route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()
	route.Name = name
	r.namedRoutes[name] = route
}

// findRoute finds a matching route for the given method and path
func (r *Router) findRoute(method, path string) (*Route, Params) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Try to match routes in priority order
	for _, route := range r.routes {
		// Check method match
		if route.Method != method {
			continue
		}

		// Check path match
		if params, ok := route.Match(path); ok {
			return route, params
		}
	}

	return nil, nil
}

// SetApp sets the app reference for PageContext
func (r *Router) SetApp(app interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.app = app
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Find matching route
	route, params := r.findRoute(req.Method, req.URL.Path)

	// Create page context
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
		Params:         params,
		values:         make(map[string]interface{}),
		app:            r.app,
	}

	var node g.Node
	var err error

	if route == nil {
		// No route found - call 404 handler
		node, err = r.notFound(ctx)
	} else {
		// Set metadata in context
		ctx.Meta = route.Metadata

		// Execute loader if present
		if route.LoaderFn != nil {
			data, loaderErr := route.executeLoader(req.Context(), params)
			if loaderErr != nil {
				// Handle loader error
				if le, ok := loaderErr.(*LoaderError); ok {
					errorHandler := r.getErrorPage(le.Status)
					node, err = errorHandler(ctx)
					if node != nil {
						if ctx.ResponseWriter.Header().Get("Content-Type") == "" {
							ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
						}
						_ = node.Render(ctx.ResponseWriter)
					}
					return
				}
				err = loaderErr
			} else {
				ctx.LoadedData = data
			}
		}

		// Build middleware chain
		handler := route.Handler

		// Apply route-specific middleware (in reverse order)
		for i := len(route.Middleware) - 1; i >= 0; i-- {
			handler = route.Middleware[i](handler)
		}

		// Apply global middleware (in reverse order)
		for i := len(r.middleware) - 1; i >= 0; i-- {
			handler = r.middleware[i](handler)
		}

		// Execute handler
		if err == nil {
			node, err = handler(ctx)
		}

		// Apply layout if configured (with composition support)
		if node != nil && err == nil {
			layoutName := route.Layout
			if layoutName == "" && r.defaultLayout != "" {
				layoutName = r.defaultLayout
			}

			if layoutName != "" && layoutName != "none" {
				node = r.applyLayoutChain(ctx, node, layoutName)
			}
		}
	}

	// Handle errors
	if err != nil {
		node = r.errorHandler(ctx, err)
	}

	// Render node if present
	if node != nil {
		if ctx.ResponseWriter.Header().Get("Content-Type") == "" {
			ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		_ = node.Render(ctx.ResponseWriter)
	}
}

// applyLayoutChain applies layouts in composition order (child -> parent -> root)
func (r *Router) applyLayoutChain(ctx *PageContext, content g.Node, layoutName string) g.Node {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Build layout chain from child to root
	var chain []LayoutFunc
	currentLayout := layoutName
	visited := make(map[string]bool)

	for currentLayout != "" {
		// Prevent infinite loops
		if visited[currentLayout] {
			break
		}
		visited[currentLayout] = true

		// Get layout function
		layoutFn, ok := r.layouts[currentLayout]
		if !ok {
			break
		}

		// Add to chain
		chain = append(chain, layoutFn)

		// Get parent layout
		if config, ok := r.layoutConfigs[currentLayout]; ok && config.Parent != "" {
			currentLayout = config.Parent
		} else {
			break
		}
	}

	// Apply layouts in reverse order (innermost first)
	result := content
	for i := len(chain) - 1; i >= 0; i-- {
		result = chain[i](ctx, result)
	}

	return result
}

// defaultNotFound is the default 404 handler
func defaultNotFound(ctx *PageContext) (g.Node, error) {
	ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
	return g.Text("404 Not Found"), nil
}

// defaultError is the default error handler
func defaultError(ctx *PageContext, err error) g.Node {
	ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	return g.Text("500 Internal Server Error: " + err.Error())
}
