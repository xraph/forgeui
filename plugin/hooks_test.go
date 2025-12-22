package plugin

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestNewHookManager(t *testing.T) {
	m := NewHookManager()

	if m == nil {
		t.Fatal("NewHookManager() returned nil")
	}

	if m.hooks == nil {
		t.Error("hooks map not initialized")
	}
}

func TestHookOn(t *testing.T) {
	m := NewHookManager()
	called := false

	m.On("test", func(ctx *HookContext) error {
		called = true
		return nil
	})

	if !m.Has("test") {
		t.Error("hook not registered")
	}

	err := m.Trigger("test", &HookContext{Context: context.Background()})
	if err != nil {
		t.Errorf("Trigger() error = %v", err)
	}

	if !called {
		t.Error("hook handler not called")
	}
}

func TestHookOff(t *testing.T) {
	m := NewHookManager()

	m.On("test", func(ctx *HookContext) error {
		return nil
	})

	if !m.Has("test") {
		t.Error("hook not registered")
	}

	m.Off("test")

	if m.Has("test") {
		t.Error("hook not removed")
	}
}

func TestHookTrigger(t *testing.T) {
	m := NewHookManager()
	count := 0

	m.On("test", func(ctx *HookContext) error {
		count++
		return nil
	})

	m.On("test", func(ctx *HookContext) error {
		count++
		return nil
	})

	err := m.Trigger("test", &HookContext{Context: context.Background()})
	if err != nil {
		t.Errorf("Trigger() error = %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 handlers called, got %d", count)
	}
}

func TestHookTriggerError(t *testing.T) {
	m := NewHookManager()
	expectedErr := errors.New("test error")

	m.On("test", func(ctx *HookContext) error {
		return expectedErr
	})

	err := m.Trigger("test", &HookContext{Context: context.Background()})
	if err == nil {
		t.Error("expected error from trigger")
	}

	if err != expectedErr {
		t.Errorf("Trigger() error = %v, want %v", err, expectedErr)
	}
}

func TestHookTriggerStopsOnError(t *testing.T) {
	m := NewHookManager()
	count := 0

	m.On("test", func(ctx *HookContext) error {
		count++
		return errors.New("stop here")
	})

	m.On("test", func(ctx *HookContext) error {
		count++
		return nil
	})

	_ = m.Trigger("test", &HookContext{Context: context.Background()})

	if count != 1 {
		t.Errorf("expected 1 handler called, got %d", count)
	}
}

func TestHookTriggerNonexistent(t *testing.T) {
	m := NewHookManager()

	err := m.Trigger("nonexistent", &HookContext{Context: context.Background()})
	if err != nil {
		t.Errorf("Trigger() error = %v for nonexistent hook", err)
	}
}

func TestHookHas(t *testing.T) {
	m := NewHookManager()

	if m.Has("test") {
		t.Error("Has() returned true for nonexistent hook")
	}

	m.On("test", func(ctx *HookContext) error {
		return nil
	})

	if !m.Has("test") {
		t.Error("Has() returned false for existing hook")
	}
}

func TestHookCount(t *testing.T) {
	m := NewHookManager()

	if m.Count("test") != 0 {
		t.Errorf("Count() = %d, want 0", m.Count("test"))
	}

	m.On("test", func(ctx *HookContext) error {
		return nil
	})

	if m.Count("test") != 1 {
		t.Errorf("Count() = %d, want 1", m.Count("test"))
	}

	m.On("test", func(ctx *HookContext) error {
		return nil
	})

	if m.Count("test") != 2 {
		t.Errorf("Count() = %d, want 2", m.Count("test"))
	}
}

func TestHookContext(t *testing.T) {
	m := NewHookManager()
	ctx := context.Background()
	data := map[string]any{"key": "value"}

	m.On("test", func(hctx *HookContext) error {
		if hctx.Context != ctx {
			t.Error("context not passed correctly")
		}

		if hctx.Data["key"] != "value" {
			t.Error("data not passed correctly")
		}

		return nil
	})

	err := m.Trigger("test", &HookContext{
		Context: ctx,
		Data:    data,
	})
	if err != nil {
		t.Errorf("Trigger() error = %v", err)
	}
}

func TestHookConcurrency(t *testing.T) {
	m := NewHookManager()

	var wg sync.WaitGroup

	// Concurrent registrations
	for range 10 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			m.On("test", func(ctx *HookContext) error {
				return nil
			})
		}()
	}

	// Concurrent triggers
	for range 10 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_ = m.Trigger("test", &HookContext{Context: context.Background()})
		}()
	}

	wg.Wait()

	count := m.Count("test")
	if count != 10 {
		t.Errorf("expected 10 handlers, got %d", count)
	}
}

func TestHookConstants(t *testing.T) {
	// Verify hook constants are defined
	hooks := []string{
		HookBeforeInit,
		HookAfterInit,
		HookBeforeShutdown,
		HookAfterShutdown,
		HookBeforeRender,
		HookAfterRender,
		HookBeforeHead,
		HookAfterHead,
		HookBeforeBody,
		HookAfterBody,
		HookBeforeScripts,
		HookAfterScripts,
	}

	for _, hook := range hooks {
		if hook == "" {
			t.Errorf("hook constant is empty")
		}
	}
}
