package router

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// LayoutFunc wraps page content with a layout
type LayoutFunc func(ctx *PageContext, content g.Node) g.Node

// LayoutConfig holds layout configuration including parent relationship
type LayoutConfig struct {
	Fn     LayoutFunc
	Parent string
}

// LayoutOption configures a layout
type LayoutOption func(*LayoutConfig)

// WithParentLayout sets the parent layout for composition
func WithParentLayout(parent string) LayoutOption {
	return func(c *LayoutConfig) {
		c.Parent = parent
	}
}

// RegisterLayout registers a named layout with optional parent
func (r *Router) RegisterLayout(name string, fn LayoutFunc, opts ...LayoutOption) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.layouts == nil {
		r.layouts = make(map[string]LayoutFunc)
	}

	if r.layoutConfigs == nil {
		r.layoutConfigs = make(map[string]*LayoutConfig)
	}

	config := &LayoutConfig{
		Fn:     fn,
		Parent: "",
	}

	for _, opt := range opts {
		opt(config)
	}

	r.layouts[name] = fn
	r.layoutConfigs[name] = config
}

// SetDefaultLayout sets the default layout for all routes
func (r *Router) SetDefaultLayout(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaultLayout = name
}

// GetLayout retrieves a layout by name
func (r *Router) GetLayout(name string) (LayoutFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.layouts == nil {
		return nil, false
	}

	fn, ok := r.layouts[name]
	return fn, ok
}

// SetLayout sets the layout for a route
func (r *Route) SetLayout(name string) *Route {
	r.Layout = name
	return r
}

// NoLayout explicitly disables layout for a route
func (r *Route) NoLayout() *Route {
	r.Layout = "none"
	return r
}

// DefaultLayout provides a basic HTML5 layout
func DefaultLayout(ctx *PageContext, content g.Node) g.Node {
	title := "ForgeUI Application"
	if ctx.Meta != nil && ctx.Meta.Title != "" {
		title = ctx.Meta.Title
	}

	return html.Doctype(
		html.HTML(
			html.Lang("en"),
			html.Head(
				html.Meta(g.Attr("charset", "UTF-8")),
				html.Meta(
					g.Attr("name", "viewport"),
					g.Attr("content", "width=device-width, initial-scale=1.0"),
				),
				html.TitleEl(g.Text(title)),
				g.If(ctx.Meta != nil, g.Group(ctx.Meta.MetaTags())),
			),
			html.Body(
				content,
			),
		),
	)
}

// BlankLayout provides minimal layout (just wraps content in HTML structure)
func BlankLayout(ctx *PageContext, content g.Node) g.Node {
	return html.Doctype(
		html.HTML(
			html.Head(
				html.Meta(g.Attr("charset", "UTF-8")),
			),
			html.Body(
				content,
			),
		),
	)
}

// DashboardLayout provides a typical dashboard layout structure
func DashboardLayout(ctx *PageContext, content g.Node) g.Node {
	title := "Dashboard"
	if ctx.Meta != nil && ctx.Meta.Title != "" {
		title = ctx.Meta.Title
	}

	return html.Doctype(
		html.HTML(
			html.Lang("en"),
			html.Head(
				html.Meta(g.Attr("charset", "UTF-8")),
				html.Meta(
					g.Attr("name", "viewport"),
					g.Attr("content", "width=device-width, initial-scale=1.0"),
				),
				html.TitleEl(g.Text(title)),
			),
			html.Body(
				html.Class("dashboard-layout"),
				html.Div(
					html.Class("dashboard-container"),
					html.Nav(
						html.Class("dashboard-sidebar"),
						html.H2(g.Text("Dashboard")),
					),
					html.Main(
						html.Class("dashboard-content"),
						content,
					),
				),
			),
		),
	)
}

