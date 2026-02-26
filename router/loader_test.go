package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
)

func TestLoaderExecution(t *testing.T) {
	r := New()

	// Create a loader
	loader := func(ctx context.Context, params Params) (any, error) {
		return "loaded data", nil
	}

	// Create a route with loader
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		data := ctx.LoaderData().(string)
		return templ.Raw(data), nil
	}).Loader(loader)

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "loaded data" {
		t.Errorf("Expected 'loaded data', got '%s'", w.Body.String())
	}
}

func TestLoaderWithParams(t *testing.T) {
	r := New()

	// Create a loader that uses params
	loader := func(ctx context.Context, params Params) (any, error) {
		id := params["id"]
		return "user:" + id, nil
	}

	// Create a route with loader
	r.Get("/users/:id", func(ctx *PageContext) (templ.Component, error) {
		data := ctx.LoaderData().(string)
		return templ.Raw(data), nil
	}).Loader(loader)

	// Test the route
	req := httptest.NewRequest(MethodGet, "/users/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Body.String() != "user:123" {
		t.Errorf("Expected 'user:123', got '%s'", w.Body.String())
	}
}

func TestLoaderError(t *testing.T) {
	r := New()

	// Create a loader that returns an error
	loader := func(ctx context.Context, params Params) (any, error) {
		return nil, Error404("Resource not found")
	}

	// Create a route with loader
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("should not reach here"), nil
	}).Loader(loader)

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestLoaderTimeout(t *testing.T) {
	r := New()

	// Create a loader that times out
	loader := func(ctx context.Context, params Params) (any, error) {
		time.Sleep(35 * time.Second) // Longer than LoaderTimeout
		return "data", nil
	}

	// Create a route with loader
	route := r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("content"), nil
	})
	route.LoaderFn = loader

	// Test the route with a short timeout context
	req := httptest.NewRequest(MethodGet, "/test", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should get a timeout error
	if w.Code != http.StatusRequestTimeout {
		t.Errorf("Expected status 408, got %d", w.Code)
	}
}

func TestLoaderDataAccess(t *testing.T) {
	r := New()

	// Create a loader
	type User struct {
		ID   int
		Name string
	}

	loader := func(ctx context.Context, params Params) (any, error) {
		return &User{ID: 1, Name: "John"}, nil
	}

	// Create a route with loader
	r.Get("/user", func(ctx *PageContext) (templ.Component, error) {
		user := ctx.LoaderData().(*User)
		return templ.Raw(user.Name), nil
	}).Loader(loader)

	// Test the route
	req := httptest.NewRequest(MethodGet, "/user", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Body.String() != "John" {
		t.Errorf("Expected 'John', got '%s'", w.Body.String())
	}
}

func TestLoaderWithNonLoaderError(t *testing.T) {
	r := New()

	// Create a loader that returns a regular error
	loader := func(ctx context.Context, params Params) (any, error) {
		return nil, errors.New("generic error")
	}

	// Create a route with loader
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("should not reach here"), nil
	}).Loader(loader)

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should get 500 error from error handler
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestNoLoader(t *testing.T) {
	r := New()

	// Create a route without loader
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		data := ctx.LoaderData()
		if data != nil {
			t.Error("Expected LoaderData to be nil when no loader is set")
		}

		return templ.Raw("content"), nil
	})

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
