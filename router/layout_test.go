package router

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestLayoutRegistration(t *testing.T) {
	r := New()

	// Register a layout
	r.RegisterLayout("test", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, "<div>Layout: "); err != nil {
				return err
			}
			if err := content.Render(tCtx, w); err != nil {
				return err
			}
			_, err := io.WriteString(w, "</div>")
			return err
		})
	})

	// Verify layout was registered
	layout, ok := r.GetLayout("test")
	if !ok {
		t.Error("Expected layout to be registered")
	}

	if layout == nil {
		t.Error("Expected layout function to not be nil")
	}
}

func TestLayoutApplication(t *testing.T) {
	r := New()

	// Register a layout
	r.RegisterLayout("wrapper", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, "<div>["); err != nil {
				return err
			}
			if err := content.Render(tCtx, w); err != nil {
				return err
			}
			_, err := io.WriteString(w, "]</div>")
			return err
		})
	})

	// Create a route with layout
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("content"), nil
	}).SetLayout("wrapper")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "[") || !strings.Contains(w.Body.String(), "content") || !strings.Contains(w.Body.String(), "]") {
		t.Errorf("Expected output to contain '[', 'content', and ']', got '%s'", w.Body.String())
	}
}

func TestDefaultLayout(t *testing.T) {
	r := New()

	// Register and set default layout
	r.RegisterLayout("default", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, "Default: "); err != nil {
				return err
			}
			return content.Render(tCtx, w)
		})
	})
	r.SetDefaultLayout("default")

	// Create a route without explicit layout
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("content"), nil
	})

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "Default: ") || !strings.Contains(w.Body.String(), "content") {
		t.Errorf("Expected output to contain 'Default: ' and 'content', got '%s'", w.Body.String())
	}
}

func TestLayoutOverride(t *testing.T) {
	r := New()

	// Register layouts
	r.RegisterLayout("default", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, "Default: "); err != nil {
				return err
			}
			return content.Render(tCtx, w)
		})
	})
	r.RegisterLayout("custom", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, "Custom: "); err != nil {
				return err
			}
			return content.Render(tCtx, w)
		})
	})
	r.SetDefaultLayout("default")

	// Create a route with custom layout
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("content"), nil
	}).SetLayout("custom")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "Custom: ") || !strings.Contains(w.Body.String(), "content") {
		t.Errorf("Expected output to contain 'Custom: ' and 'content', got '%s'", w.Body.String())
	}
}

func TestNoLayout(t *testing.T) {
	r := New()

	// Register and set default layout
	r.RegisterLayout("default", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, "Default: "); err != nil {
				return err
			}
			return content.Render(tCtx, w)
		})
	})
	r.SetDefaultLayout("default")

	// Create a route with no layout
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("content"), nil
	}).NoLayout()

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Body.String() != "content" {
		t.Errorf("Expected 'content', got '%s'", w.Body.String())
	}
}

func TestDefaultLayoutFunction(t *testing.T) {
	ctx := &PageContext{
		Request: httptest.NewRequest(MethodGet, "/test", nil),
		Meta: &RouteMeta{
			Title: "Test Page",
		},
	}

	content := templ.Raw("Page content")
	result := DefaultLayout(ctx, content)

	if result == nil {
		t.Error("Expected layout to return a component")
	}

	// Render and check output
	var buf strings.Builder

	_ = result.Render(context.Background(), &buf)
	output := buf.String()

	if !strings.Contains(output, "Test Page") {
		t.Error("Expected layout to contain title")
	}

	if !strings.Contains(output, "Page content") {
		t.Error("Expected layout to contain content")
	}
}

func TestBlankLayoutFunction(t *testing.T) {
	ctx := &PageContext{
		Request: httptest.NewRequest(MethodGet, "/test", nil),
	}
	content := templ.Raw("Content")
	result := BlankLayout(ctx, content)

	if result == nil {
		t.Error("Expected layout to return a component")
	}

	// Render and check output
	var buf strings.Builder

	_ = result.Render(context.Background(), &buf)
	output := buf.String()

	if !strings.Contains(output, "Content") {
		t.Error("Expected layout to contain content")
	}
}

func TestDashboardLayoutFunction(t *testing.T) {
	ctx := &PageContext{
		Request: httptest.NewRequest(MethodGet, "/test", nil),
		Meta: &RouteMeta{
			Title: "Dashboard",
		},
	}

	content := templ.Raw("Dashboard content")
	result := DashboardLayout(ctx, content)

	if result == nil {
		t.Error("Expected layout to return a component")
	}

	// Render and check output
	var buf strings.Builder

	_ = result.Render(context.Background(), &buf)
	output := buf.String()

	if !strings.Contains(output, "Dashboard") {
		t.Error("Expected layout to contain title")
	}

	if !strings.Contains(output, "Dashboard content") {
		t.Error("Expected layout to contain content")
	}

	if !strings.Contains(output, "dashboard-layout") {
		t.Error("Expected layout to have dashboard-layout class")
	}
}
