package plugin

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestInitialize(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p1 := newMockPlugin("A", "1.0.0", nil)
	p2 := newMockPlugin("B", "1.0.0", nil)

	_ = r.Register(p1)
	_ = r.Register(p2)

	err := r.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if !p1.initCalled {
		t.Error("plugin A Init not called")
	}

	if !p2.initCalled {
		t.Error("plugin B Init not called")
	}
}

func TestInitializeWithDependencies(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	// B depends on A
	pA := newMockPlugin("A", "1.0.0", nil)
	pB := newMockPlugin("B", "1.0.0", []Dependency{
		{Name: "A", Version: ">=1.0.0"},
	})

	_ = r.Register(pA)
	_ = r.Register(pB)

	err := r.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if !pA.initCalled {
		t.Error("plugin A Init not called")
	}

	if !pB.initCalled {
		t.Error("plugin B Init not called")
	}
}

func TestInitializeError(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p := newMockPlugin("test", "1.0.0", nil)
	p.initError = errors.New("init failed")

	_ = r.Register(p)

	err := r.Initialize(ctx)
	if err == nil {
		t.Error("expected error from Initialize")
	}
}

func TestInitializeDependencyError(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	// A depends on B, but B is not registered
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0"},
	})

	_ = r.Register(pA)

	err := r.Initialize(ctx)
	if err == nil {
		t.Error("expected error for missing dependency")
	}
}

func TestInitializeHooks(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	beforeCalled := false
	afterCalled := false

	r.Hooks().On(HookBeforeInit, func(hctx *HookContext) error {
		beforeCalled = true
		return nil
	})

	r.Hooks().On(HookAfterInit, func(hctx *HookContext) error {
		afterCalled = true
		return nil
	})

	p := newMockPlugin("test", "1.0.0", nil)
	_ = r.Register(p)

	err := r.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if !beforeCalled {
		t.Error("before_init hook not called")
	}

	if !afterCalled {
		t.Error("after_init hook not called")
	}
}

func TestShutdown(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p1 := newMockPlugin("A", "1.0.0", nil)
	p2 := newMockPlugin("B", "1.0.0", nil)

	_ = r.Register(p1)
	_ = r.Register(p2)
	_ = r.Initialize(ctx)

	err := r.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	if !p1.shutdownCalled {
		t.Error("plugin A Shutdown not called")
	}

	if !p2.shutdownCalled {
		t.Error("plugin B Shutdown not called")
	}
}

func TestShutdownReverseOrder(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	var (
		order []string
		mu    sync.Mutex
	)

	// B depends on A, so init order is A, B
	// Shutdown order should be B, A
	pA := &mockPluginWithOrder{
		mockPlugin: newMockPlugin("A", "1.0.0", nil),
		order:      &order,
		mu:         &mu,
	}

	pB := &mockPluginWithOrder{
		mockPlugin: newMockPlugin("B", "1.0.0", []Dependency{
			{Name: "A", Version: ">=1.0.0"},
		}),
		order: &order,
		mu:    &mu,
	}

	_ = r.Register(pA)
	_ = r.Register(pB)
	_ = r.Initialize(ctx)

	err := r.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	if len(order) != 2 {
		t.Fatalf("expected 2 shutdowns, got %d", len(order))
	}

	if order[0] != "B" {
		t.Errorf("expected B to shutdown first, got %s", order[0])
	}

	if order[1] != "A" {
		t.Errorf("expected A to shutdown second, got %s", order[1])
	}
}

type mockPluginWithOrder struct {
	*mockPlugin

	order *[]string
	mu    *sync.Mutex
}

func (m *mockPluginWithOrder) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	*m.order = append(*m.order, m.Name())
	m.mu.Unlock()

	return m.mockPlugin.Shutdown(ctx)
}

func TestShutdownContinuesOnError(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p1 := newMockPlugin("A", "1.0.0", nil)
	p1.shutdownError = errors.New("shutdown failed")

	p2 := newMockPlugin("B", "1.0.0", nil)

	_ = r.Register(p1)
	_ = r.Register(p2)
	_ = r.Initialize(ctx)

	err := r.Shutdown(ctx)
	if err == nil {
		t.Error("expected error from Shutdown")
	}

	// Both should be called despite error
	if !p1.shutdownCalled {
		t.Error("plugin A Shutdown not called")
	}

	if !p2.shutdownCalled {
		t.Error("plugin B Shutdown not called")
	}
}

func TestShutdownHooks(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	beforeCalled := false
	afterCalled := false

	r.Hooks().On(HookBeforeShutdown, func(hctx *HookContext) error {
		beforeCalled = true
		return nil
	})

	r.Hooks().On(HookAfterShutdown, func(hctx *HookContext) error {
		afterCalled = true
		return nil
	})

	p := newMockPlugin("test", "1.0.0", nil)
	_ = r.Register(p)
	_ = r.Initialize(ctx)

	err := r.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	if !beforeCalled {
		t.Error("before_shutdown hook not called")
	}

	if !afterCalled {
		t.Error("after_shutdown hook not called")
	}
}

func TestShutdownWithoutInitialize(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p := newMockPlugin("test", "1.0.0", nil)
	_ = r.Register(p)

	// Shutdown without Initialize should not error
	err := r.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	if p.shutdownCalled {
		t.Error("plugin Shutdown should not be called without Initialize")
	}
}
