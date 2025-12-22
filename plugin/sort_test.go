package plugin

import (
	"testing"
)

func TestTopologicalSortEmpty(t *testing.T) {
	r := NewRegistry()

	sorted, err := r.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	if len(sorted) != 0 {
		t.Errorf("expected empty result, got %d plugins", len(sorted))
	}
}

func TestTopologicalSortNoDeps(t *testing.T) {
	r := NewRegistry()

	p1 := newMockPlugin("A", "1.0.0", nil)
	p2 := newMockPlugin("B", "1.0.0", nil)
	p3 := newMockPlugin("C", "1.0.0", nil)

	_ = r.Register(p1)
	_ = r.Register(p2)
	_ = r.Register(p3)

	sorted, err := r.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	if len(sorted) != 3 {
		t.Errorf("expected 3 plugins, got %d", len(sorted))
	}
}

func TestTopologicalSortLinear(t *testing.T) {
	r := NewRegistry()

	// C depends on B, B depends on A
	// Expected order: A, B, C
	pA := newMockPlugin("A", "1.0.0", nil)
	pB := newMockPlugin("B", "1.0.0", []Dependency{
		{Name: "A", Version: ">=1.0.0"},
	})
	pC := newMockPlugin("C", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0"},
	})

	_ = r.Register(pA)
	_ = r.Register(pB)
	_ = r.Register(pC)

	sorted, err := r.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	if len(sorted) != 3 {
		t.Fatalf("expected 3 plugins, got %d", len(sorted))
	}

	// Verify order
	names := make([]string, len(sorted))
	for i, p := range sorted {
		names[i] = p.Name()
	}

	if names[0] != "A" {
		t.Errorf("expected A first, got %s", names[0])
	}

	if names[1] != "B" {
		t.Errorf("expected B second, got %s", names[1])
	}

	if names[2] != "C" {
		t.Errorf("expected C third, got %s", names[2])
	}
}

func TestTopologicalSortDiamond(t *testing.T) {
	r := NewRegistry()

	// D depends on B and C
	// B depends on A
	// C depends on A
	// Expected order: A, then B and C (any order), then D
	pA := newMockPlugin("A", "1.0.0", nil)
	pB := newMockPlugin("B", "1.0.0", []Dependency{
		{Name: "A", Version: ">=1.0.0"},
	})
	pC := newMockPlugin("C", "1.0.0", []Dependency{
		{Name: "A", Version: ">=1.0.0"},
	})
	pD := newMockPlugin("D", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0"},
		{Name: "C", Version: ">=1.0.0"},
	})

	_ = r.Register(pA)
	_ = r.Register(pB)
	_ = r.Register(pC)
	_ = r.Register(pD)

	sorted, err := r.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	if len(sorted) != 4 {
		t.Fatalf("expected 4 plugins, got %d", len(sorted))
	}

	// Build position map
	pos := make(map[string]int)
	for i, p := range sorted {
		pos[p.Name()] = i
	}

	// A must come before B, C, and D
	if pos["A"] >= pos["B"] || pos["A"] >= pos["C"] || pos["A"] >= pos["D"] {
		t.Error("A must come before B, C, and D")
	}

	// B and C must come before D
	if pos["B"] >= pos["D"] || pos["C"] >= pos["D"] {
		t.Error("B and C must come before D")
	}
}

func TestTopologicalSortCircular(t *testing.T) {
	r := NewRegistry()

	// A depends on B, B depends on A (circular)
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0"},
	})
	pB := newMockPlugin("B", "1.0.0", []Dependency{
		{Name: "A", Version: ">=1.0.0"},
	})

	_ = r.Register(pA)
	_ = r.Register(pB)

	_, err := r.TopologicalSort()
	if err == nil {
		t.Error("expected error for circular dependency")
	}
}

func TestTopologicalSortCircularComplex(t *testing.T) {
	r := NewRegistry()

	// A -> B -> C -> A (circular)
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0"},
	})
	pB := newMockPlugin("B", "1.0.0", []Dependency{
		{Name: "C", Version: ">=1.0.0"},
	})
	pC := newMockPlugin("C", "1.0.0", []Dependency{
		{Name: "A", Version: ">=1.0.0"},
	})

	_ = r.Register(pA)
	_ = r.Register(pB)
	_ = r.Register(pC)

	_, err := r.TopologicalSort()
	if err == nil {
		t.Error("expected error for circular dependency")
	}
}

func TestTopologicalSortOptionalDeps(t *testing.T) {
	r := NewRegistry()

	// A optionally depends on B (B not registered)
	// Should still work
	pA := newMockPlugin("A", "1.0.0", []Dependency{
		{Name: "B", Version: ">=1.0.0", Optional: true},
	})

	_ = r.Register(pA)

	sorted, err := r.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	if len(sorted) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(sorted))
	}
}
