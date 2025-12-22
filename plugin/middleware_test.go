package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewarePluginBase(t *testing.T) {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "true")
			next.ServeHTTP(w, r)
		})
	}

	plugin := NewMiddlewarePluginBase(
		PluginInfo{Name: "test-middleware", Version: "1.0.0"},
		middleware,
		10,
	)

	if plugin.Name() != "test-middleware" {
		t.Errorf("expected name 'test-middleware', got %s", plugin.Name())
	}

	if plugin.Priority() != 10 {
		t.Errorf("expected priority 10, got %d", plugin.Priority())
	}

	if plugin.Middleware() == nil {
		t.Error("expected middleware function, got nil")
	}
}

func TestMiddlewarePluginDefaultPriority(t *testing.T) {
	plugin := NewMiddlewarePluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		nil,
		0, // Zero should default to 50
	)

	if plugin.Priority() != 50 {
		t.Errorf("expected default priority 50, got %d", plugin.Priority())
	}
}

func TestMiddlewarePluginExecution(t *testing.T) {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "executed")
			next.ServeHTTP(w, r)
		})
	}

	plugin := NewMiddlewarePluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		middleware,
		10,
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := plugin.Middleware()(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	if w.Header().Get("X-Middleware") != "executed" {
		t.Error("middleware not executed")
	}
}

func TestMiddlewarePluginRegistration(t *testing.T) {
	registry := NewRegistry()

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	plugin := NewMiddlewarePluginBase(
		PluginInfo{Name: "middleware", Version: "1.0.0"},
		middleware,
		10,
	)

	err := registry.Register(plugin)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Should be in middleware collection
	mws := registry.CollectMiddleware()
	if len(mws) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(mws))
	}

	if mws[0].Name() != "middleware" {
		t.Errorf("expected name 'middleware', got %s", mws[0].Name())
	}
}

func TestCollectMiddlewarePriority(t *testing.T) {
	registry := NewRegistry()

	mw1 := NewMiddlewarePluginBase(
		PluginInfo{Name: "mw1", Version: "1.0.0"},
		func(next http.Handler) http.Handler { return next },
		50,
	)

	mw2 := NewMiddlewarePluginBase(
		PluginInfo{Name: "mw2", Version: "1.0.0"},
		func(next http.Handler) http.Handler { return next },
		10,
	)

	mw3 := NewMiddlewarePluginBase(
		PluginInfo{Name: "mw3", Version: "1.0.0"},
		func(next http.Handler) http.Handler { return next },
		30,
	)

	_ = registry.Register(mw1)
	_ = registry.Register(mw2)
	_ = registry.Register(mw3)

	mws := registry.CollectMiddleware()
	if len(mws) != 3 {
		t.Fatalf("expected 3 middleware, got %d", len(mws))
	}

	// Should be sorted by priority (lower first)
	if mws[0].Name() != "mw2" {
		t.Errorf("expected mw2 first (priority 10), got %s", mws[0].Name())
	}
	if mws[1].Name() != "mw3" {
		t.Errorf("expected mw3 second (priority 30), got %s", mws[1].Name())
	}
	if mws[2].Name() != "mw1" {
		t.Errorf("expected mw1 third (priority 50), got %s", mws[2].Name())
	}
}

func TestMiddlewarePluginLifecycle(t *testing.T) {
	registry := NewRegistry()
	ctx := context.Background()

	plugin := NewMiddlewarePluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		func(next http.Handler) http.Handler { return next },
		10,
	)

	_ = registry.Register(plugin)

	err := registry.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	err = registry.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
}

