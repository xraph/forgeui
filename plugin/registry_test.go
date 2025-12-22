package plugin

import (
	"sync"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()

	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	if r.plugins == nil {
		t.Error("plugins map not initialized")
	}

	if r.hooks == nil {
		t.Error("hooks not initialized")
	}

	if r.Count() != 0 {
		t.Errorf("expected 0 plugins, got %d", r.Count())
	}
}

func TestRegister(t *testing.T) {
	r := NewRegistry()
	p := newMockPlugin("test", "1.0.0", nil)

	err := r.Register(p)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if r.Count() != 1 {
		t.Errorf("expected 1 plugin, got %d", r.Count())
	}

	if !r.Has("test") {
		t.Error("plugin not found after registration")
	}
}

func TestRegisterEmptyName(t *testing.T) {
	r := NewRegistry()
	p := newMockPlugin("", "1.0.0", nil)

	err := r.Register(p)
	if err == nil {
		t.Error("expected error for empty plugin name")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	p1 := newMockPlugin("test", "1.0.0", nil)
	p2 := newMockPlugin("test", "2.0.0", nil)

	_ = r.Register(p1)

	err := r.Register(p2)
	if err == nil {
		t.Error("expected error for duplicate plugin")
	}
}

func TestGet(t *testing.T) {
	r := NewRegistry()
	p := newMockPlugin("test", "1.0.0", nil)
	_ = r.Register(p)

	got, ok := r.Get("test")
	if !ok {
		t.Fatal("Get() returned false for existing plugin")
	}

	if got.Name() != "test" {
		t.Errorf("Get() name = %v, want test", got.Name())
	}
}

func TestGetNotFound(t *testing.T) {
	r := NewRegistry()

	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("Get() returned true for nonexistent plugin")
	}
}

func TestAll(t *testing.T) {
	r := NewRegistry()
	p1 := newMockPlugin("test1", "1.0.0", nil)
	p2 := newMockPlugin("test2", "1.0.0", nil)

	_ = r.Register(p1)
	_ = r.Register(p2)

	all := r.All()
	if len(all) != 2 {
		t.Errorf("All() returned %d plugins, want 2", len(all))
	}
}

func TestUse(t *testing.T) {
	r := NewRegistry()
	p1 := newMockPlugin("test1", "1.0.0", nil)
	p2 := newMockPlugin("test2", "1.0.0", nil)

	r.Use(p1, p2)

	if r.Count() != 2 {
		t.Errorf("Use() registered %d plugins, want 2", r.Count())
	}
}

func TestUnregister(t *testing.T) {
	r := NewRegistry()
	p := newMockPlugin("test", "1.0.0", nil)
	_ = r.Register(p)

	err := r.Unregister("test")
	if err != nil {
		t.Fatalf("Unregister() error = %v", err)
	}

	if r.Has("test") {
		t.Error("plugin still exists after unregister")
	}
}

func TestUnregisterNotFound(t *testing.T) {
	r := NewRegistry()

	err := r.Unregister("nonexistent")
	if err == nil {
		t.Error("expected error for unregistering nonexistent plugin")
	}
}

func TestHas(t *testing.T) {
	r := NewRegistry()
	p := newMockPlugin("test", "1.0.0", nil)
	_ = r.Register(p)

	if !r.Has("test") {
		t.Error("Has() returned false for existing plugin")
	}

	if r.Has("nonexistent") {
		t.Error("Has() returned true for nonexistent plugin")
	}
}

func TestCount(t *testing.T) {
	r := NewRegistry()

	if r.Count() != 0 {
		t.Errorf("Count() = %d, want 0", r.Count())
	}

	p1 := newMockPlugin("test1", "1.0.0", nil)
	_ = r.Register(p1)

	if r.Count() != 1 {
		t.Errorf("Count() = %d, want 1", r.Count())
	}

	p2 := newMockPlugin("test2", "1.0.0", nil)
	_ = r.Register(p2)

	if r.Count() != 2 {
		t.Errorf("Count() = %d, want 2", r.Count())
	}
}

func TestHooks(t *testing.T) {
	r := NewRegistry()
	hooks := r.Hooks()

	if hooks == nil {
		t.Error("Hooks() returned nil")
	}

	if hooks != r.hooks {
		t.Error("Hooks() returned different instance")
	}
}

func TestResolveDependencies(t *testing.T) {
	r := NewRegistry()

	// Plugin A depends on B
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0"},
	})
	pB := newMockPlugin("B", "1.0.0", nil)

	_ = r.Register(pA)
	_ = r.Register(pB)

	err := r.ResolveDependencies()
	if err != nil {
		t.Errorf("ResolveDependencies() error = %v", err)
	}
}

func TestResolveDependenciesMissing(t *testing.T) {
	r := NewRegistry()

	// Plugin A depends on B, but B is not registered
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0"},
	})

	_ = r.Register(pA)

	err := r.ResolveDependencies()
	if err == nil {
		t.Error("expected error for missing dependency")
	}
}

func TestResolveDependenciesOptional(t *testing.T) {
	r := NewRegistry()

	// Plugin A optionally depends on B
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0", Optional: true},
	})

	_ = r.Register(pA)

	err := r.ResolveDependencies()
	if err != nil {
		t.Errorf("ResolveDependencies() error = %v for optional dependency", err)
	}
}

func TestResolveDependenciesVersionMismatch(t *testing.T) {
	r := NewRegistry()

	// Plugin A requires B >= 2.0.0, but B is 1.0.0
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=2.0.0"},
	})
	pB := newMockPlugin("B", "1.0.0", nil)

	_ = r.Register(pA)
	_ = r.Register(pB)

	err := r.ResolveDependencies()
	if err == nil {
		t.Error("expected error for version mismatch")
	}
}

func TestRegistryConcurrency(t *testing.T) {
	r := NewRegistry()

	var wg sync.WaitGroup

	// Concurrent registrations
	for i := range 10 {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			p := newMockPlugin(string(rune('A'+n)), "1.0.0", nil)
			_ = r.Register(p)
		}(i)
	}

	// Concurrent reads
	for range 10 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_ = r.Count()
			_ = r.All()
		}()
	}

	wg.Wait()

	if r.Count() != 10 {
		t.Errorf("expected 10 plugins after concurrent operations, got %d", r.Count())
	}
}
