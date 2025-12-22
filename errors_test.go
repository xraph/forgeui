package forgeui

import (
	"errors"
	"testing"
)

func TestComponentError(t *testing.T) {
	t.Run("with wrapped error", func(t *testing.T) {
		innerErr := errors.New("inner error")
		err := &ComponentError{
			Component: "Button",
			Message:   "render failed",
			Err:       innerErr,
		}

		expected := "Button: render failed: inner error"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}

		if unwrapped := err.Unwrap(); !errors.Is(unwrapped, innerErr) {
			t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
		}
	})

	t.Run("without wrapped error", func(t *testing.T) {
		err := &ComponentError{
			Component: "Card",
			Message:   "invalid props",
		}

		expected := "Card: invalid props"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}

		if unwrapped := err.Unwrap(); unwrapped != nil {
			t.Errorf("Unwrap() should return nil when no wrapped error")
		}
	})
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "variant",
		Message: "unknown variant type",
	}

	expected := "validation error on variant: unknown variant type"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestPluginError(t *testing.T) {
	t.Run("with wrapped error", func(t *testing.T) {
		innerErr := errors.New("init failed")
		err := &PluginError{
			Plugin:  "toast-plugin",
			Message: "initialization error",
			Err:     innerErr,
		}

		expected := "plugin toast-plugin: initialization error: init failed"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}

		if unwrapped := err.Unwrap(); !errors.Is(unwrapped, innerErr) {
			t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
		}
	})

	t.Run("without wrapped error", func(t *testing.T) {
		err := &PluginError{
			Plugin:  "chart-plugin",
			Message: "not found",
		}

		expected := "plugin chart-plugin: not found"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}
	})
}
