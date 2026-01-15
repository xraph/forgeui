package router

// GroupPageBuilder provides a fluent API for building and registering pages within a group
type GroupPageBuilder struct {
	group      *Group
	pattern    string
	method     string
	handler    PageHandler
	loader     LoaderFunc
	layout     string
	middleware []Middleware
	meta       *RouteMeta
	name       string
	noLayout   bool
}

// Handler sets the page handler function
func (gpb *GroupPageBuilder) Handler(handler PageHandler) *GroupPageBuilder {
	gpb.handler = handler
	return gpb
}

// Loader sets the data loader function
func (gpb *GroupPageBuilder) Loader(loader LoaderFunc) *GroupPageBuilder {
	gpb.loader = loader
	return gpb
}

// Layout sets the layout name for this page
func (gpb *GroupPageBuilder) Layout(layout string) *GroupPageBuilder {
	gpb.layout = layout
	gpb.noLayout = false

	return gpb
}

// NoLayout disables layout for this page
func (gpb *GroupPageBuilder) NoLayout() *GroupPageBuilder {
	gpb.noLayout = true
	gpb.layout = "none"

	return gpb
}

// Meta sets the page metadata (title and description)
func (gpb *GroupPageBuilder) Meta(title, description string) *GroupPageBuilder {
	if gpb.meta == nil {
		gpb.meta = &RouteMeta{}
	}

	gpb.meta.Title = title
	gpb.meta.Description = description

	return gpb
}

// MetaFull sets complete page metadata
func (gpb *GroupPageBuilder) MetaFull(meta *RouteMeta) *GroupPageBuilder {
	gpb.meta = meta
	return gpb
}

// Middleware adds middleware to this page
func (gpb *GroupPageBuilder) Middleware(middleware ...Middleware) *GroupPageBuilder {
	gpb.middleware = append(gpb.middleware, middleware...)
	return gpb
}

// Name sets a name for this route (for URL generation)
func (gpb *GroupPageBuilder) Name(name string) *GroupPageBuilder {
	gpb.name = name
	return gpb
}

// Method sets the HTTP method for this page
func (gpb *GroupPageBuilder) Method(method string) *GroupPageBuilder {
	gpb.method = method
	return gpb
}

// Register registers the page with the router through the group
func (gpb *GroupPageBuilder) Register() *Route {
	if gpb.handler == nil {
		panic("page handler is required")
	}

	// Create route based on method using group methods
	var route *Route

	switch gpb.method {
	case MethodGet:
		route = gpb.group.Get(gpb.pattern, gpb.handler)
	case MethodPost:
		route = gpb.group.Post(gpb.pattern, gpb.handler)
	case MethodPut:
		route = gpb.group.Put(gpb.pattern, gpb.handler)
	case MethodPatch:
		route = gpb.group.Patch(gpb.pattern, gpb.handler)
	case MethodDelete:
		route = gpb.group.Delete(gpb.pattern, gpb.handler)
	default:
		route = gpb.group.Get(gpb.pattern, gpb.handler)
	}

	// Apply loader if set
	if gpb.loader != nil {
		route.WithLoader(gpb.loader)
	}

	// Apply layout if set
	if gpb.noLayout {
		route.NoLayout()
	} else if gpb.layout != "" {
		route.SetLayout(gpb.layout)
	}

	// Apply metadata if set
	if gpb.meta != nil {
		route.WithMeta(gpb.meta)
	}

	// Apply middleware if set
	if len(gpb.middleware) > 0 {
		route.WithMiddleware(gpb.middleware...)
	}

	// Apply name if set
	if gpb.name != "" {
		gpb.group.router.Name(gpb.name, route)
	}

	return route
}

// GET is a convenience method to set method to GET
func (gpb *GroupPageBuilder) GET() *GroupPageBuilder {
	gpb.method = MethodGet
	return gpb
}

// POST is a convenience method to set method to POST
func (gpb *GroupPageBuilder) POST() *GroupPageBuilder {
	gpb.method = MethodPost
	return gpb
}

// PUT is a convenience method to set method to PUT
func (gpb *GroupPageBuilder) PUT() *GroupPageBuilder {
	gpb.method = MethodPut
	return gpb
}

// PATCH is a convenience method to set method to PATCH
func (gpb *GroupPageBuilder) PATCH() *GroupPageBuilder {
	gpb.method = MethodPatch
	return gpb
}

// DELETE is a convenience method to set method to DELETE
func (gpb *GroupPageBuilder) DELETE() *GroupPageBuilder {
	gpb.method = MethodDelete
	return gpb
}
