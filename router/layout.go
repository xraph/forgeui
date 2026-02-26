package router

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

// LayoutFunc wraps page content with a layout.
type LayoutFunc func(ctx *PageContext, content templ.Component) templ.Component

// LayoutConfig holds layout configuration including parent relationship.
type LayoutConfig struct {
	Fn     LayoutFunc
	Parent string
}

// LayoutOption configures a layout.
type LayoutOption func(*LayoutConfig)

// WithParentLayout sets the parent layout for composition.
func WithParentLayout(parent string) LayoutOption {
	return func(c *LayoutConfig) {
		c.Parent = parent
	}
}

// RegisterLayout registers a named layout with optional parent.
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

// SetDefaultLayout sets the default layout for all routes.
func (r *Router) SetDefaultLayout(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.defaultLayout = name
}

// GetLayout retrieves a layout by name.
func (r *Router) GetLayout(name string) (LayoutFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.layouts == nil {
		return nil, false
	}

	fn, ok := r.layouts[name]

	return fn, ok
}

// SetLayout sets the layout for a route.
func (r *Route) SetLayout(name string) *Route {
	r.Layout = name
	return r
}

// NoLayout explicitly disables layout for a route.
func (r *Route) NoLayout() *Route {
	r.Layout = "none"
	return r
}

// DefaultLayout provides a basic HTML5 layout.
var DefaultLayout LayoutFunc = func(ctx *PageContext, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(renderCtx context.Context, w io.Writer) error {
		title := "ForgeUI Application"
		if ctx.Meta != nil && ctx.Meta.Title != "" {
			title = ctx.Meta.Title
		}

		if _, err := io.WriteString(w, `<!doctype html><html lang="en"><head>`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<meta charset="UTF-8"/>`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<meta name="viewport" content="width=device-width, initial-scale=1.0"/>`); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, `<title>%s</title>`, title); err != nil {
			return err
		}

		// Render meta tags
		if ctx.Meta != nil {
			if err := ctx.Meta.MetaTags().Render(renderCtx, w); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, `</head><body>`); err != nil {
			return err
		}

		// Render page content
		if err := content.Render(renderCtx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</body></html>`)
		return err
	})
}

// BlankLayout provides minimal layout (just wraps content in HTML structure).
var BlankLayout LayoutFunc = func(_ *PageContext, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(renderCtx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<!doctype html><html><head><meta charset="UTF-8"/></head><body>`); err != nil {
			return err
		}

		if err := content.Render(renderCtx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</body></html>`)
		return err
	})
}

// DashboardLayout provides a typical dashboard layout structure.
var DashboardLayout LayoutFunc = func(ctx *PageContext, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(renderCtx context.Context, w io.Writer) error {
		title := "Dashboard"
		if ctx.Meta != nil && ctx.Meta.Title != "" {
			title = ctx.Meta.Title
		}

		if _, err := io.WriteString(w, `<!doctype html><html lang="en"><head>`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<meta charset="UTF-8"/>`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<meta name="viewport" content="width=device-width, initial-scale=1.0"/>`); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, `<title>%s</title>`, title); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `</head><body class="dashboard-layout">`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<div class="dashboard-container">`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<nav class="dashboard-sidebar"><h2>Dashboard</h2></nav>`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<main class="dashboard-content">`); err != nil {
			return err
		}

		if err := content.Render(renderCtx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</main></div></body></html>`)
		return err
	})
}
