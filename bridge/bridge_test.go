package bridge

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	b := New()

	if b == nil {
		t.Fatal("New() returned nil")
	}

	if b.config == nil {
		t.Fatal("Bridge config is nil")
	}

	if b.config.Timeout != 30*time.Second {
		t.Errorf("Default timeout = %v, want 30s", b.config.Timeout)
	}

	if b.config.MaxBatchSize != 10 {
		t.Errorf("Default MaxBatchSize = %d, want 10", b.config.MaxBatchSize)
	}

	if !b.config.EnableCSRF {
		t.Error("Default EnableCSRF = false, want true")
	}
}

func TestNew_WithOptions(t *testing.T) {
	b := New(
		WithTimeout(60*time.Second),
		WithMaxBatchSize(20),
		WithCSRF(false),
		WithDefaultRateLimit(100),
	)

	if b.config.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s", b.config.Timeout)
	}

	if b.config.MaxBatchSize != 20 {
		t.Errorf("MaxBatchSize = %d, want 20", b.config.MaxBatchSize)
	}

	if b.config.EnableCSRF {
		t.Error("EnableCSRF = true, want false")
	}

	if b.config.DefaultRateLimit != 100 {
		t.Errorf("DefaultRateLimit = %d, want 100", b.config.DefaultRateLimit)
	}
}

func TestBridge_Register(t *testing.T) {
	b := New()

	err := b.Register("testFunc", validHandler)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if b.FunctionCount() != 1 {
		t.Errorf("FunctionCount() = %d, want 1", b.FunctionCount())
	}

	if !b.HasFunction("testFunc") {
		t.Error("HasFunction(testFunc) = false, want true")
	}
}

func TestBridge_Register_Duplicate(t *testing.T) {
	b := New()

	err := b.Register("testFunc", validHandler)
	if err != nil {
		t.Fatalf("First Register() error = %v", err)
	}

	err = b.Register("testFunc", validHandler)
	if err == nil {
		t.Error("Second Register() error = nil, want error")
	}
}

func TestBridge_Register_WithOptions(t *testing.T) {
	b := New()

	err := b.Register("testFunc", validHandler,
		RequireAuth(),
		WithFunctionTimeout(10*time.Second),
		WithRateLimit(50),
	)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	fn, err := b.GetFunction("testFunc")
	if err != nil {
		t.Fatalf("GetFunction() error = %v", err)
	}

	if !fn.RequireAuth {
		t.Error("Function RequireAuth = false, want true")
	}

	if fn.Timeout != 10*time.Second {
		t.Errorf("Function Timeout = %v, want 10s", fn.Timeout)
	}

	if fn.RateLimit != 50 {
		t.Errorf("Function RateLimit = %d, want 50", fn.RateLimit)
	}
}

func TestBridge_Unregister(t *testing.T) {
	b := New()

	b.Register("testFunc", validHandler)

	err := b.Unregister("testFunc")
	if err != nil {
		t.Fatalf("Unregister() error = %v", err)
	}

	if b.HasFunction("testFunc") {
		t.Error("HasFunction(testFunc) = true after Unregister, want false")
	}

	// Try to unregister non-existent function
	err = b.Unregister("nonexistent")
	if err == nil {
		t.Error("Unregister(nonexistent) error = nil, want error")
	}
}

func TestBridge_GetFunction(t *testing.T) {
	b := New()

	b.Register("testFunc", validHandler)

	fn, err := b.GetFunction("testFunc")
	if err != nil {
		t.Fatalf("GetFunction() error = %v", err)
	}

	if fn.Name != "testFunc" {
		t.Errorf("Function Name = %s, want testFunc", fn.Name)
	}

	// Try to get non-existent function
	_, err = b.GetFunction("nonexistent")
	if err == nil {
		t.Error("GetFunction(nonexistent) error = nil, want error")
	}
}

func TestBridge_ListFunctions(t *testing.T) {
	b := New()

	b.Register("func1", validHandler)
	b.Register("func2", validHandler)
	b.Register("func3", validHandler)

	names := b.ListFunctions()
	if len(names) != 3 {
		t.Errorf("len(ListFunctions()) = %d, want 3", len(names))
	}

	// Check all names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, expected := range []string{"func1", "func2", "func3"} {
		if !nameMap[expected] {
			t.Errorf("ListFunctions() missing %s", expected)
		}
	}
}

func TestBridge_ConcurrentAccess(t *testing.T) {
	b := New()

	// Register initial function
	b.Register("func1", validHandler)

	// Test concurrent read and write
	done := make(chan bool)

	// Concurrent readers
	for range 10 {
		go func() {
			for range 100 {
				b.HasFunction("func1")
				b.ListFunctions()
			}

			done <- true
		}()
	}

	// Concurrent writers
	for i := range 5 {
		go func(n int) {
			for j := range 20 {
				name := fmt.Sprintf("func%d_%d", n, j)
				b.Register(name, validHandler)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for range 15 {
		<-done
	}
}
