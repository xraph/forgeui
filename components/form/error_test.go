package form

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestError(t *testing.T) {
	t.Run("renders error message", func(t *testing.T) {
		err := Error("This field is required")

		var buf bytes.Buffer
		err.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "This field is required") {
			t.Error("expected error to contain message")
		}
		if !strings.Contains(output, "text-destructive") {
			t.Error("expected error to have destructive color")
		}
		if !strings.Contains(output, "role=\"alert\"") {
			t.Error("expected error to have alert role")
		}
		if !strings.Contains(output, "aria-live=\"polite\"") {
			t.Error("expected error to have aria-live attribute")
		}
	})

	t.Run("returns nil for empty error", func(t *testing.T) {
		err := Error("")

		if err != nil {
			t.Error("expected nil for empty error message")
		}
	})

	t.Run("renders error with custom class", func(t *testing.T) {
		err := Error(
			"Invalid input",
			WithErrorClass("mt-4"),
		)

		var buf bytes.Buffer
		err.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "mt-4") {
			t.Error("expected error to contain custom class")
		}
	})

	t.Run("renders error with custom attributes", func(t *testing.T) {
		err := Error(
			"Invalid input",
			WithErrorAttrs(g.Attr("data-testid", "error-message")),
		)

		var buf bytes.Buffer
		err.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-testid=\"error-message\"") {
			t.Error("expected error to contain custom attribute")
		}
	})

	t.Run("has correct base classes", func(t *testing.T) {
		err := Error("Error message")

		var buf bytes.Buffer
		err.Render(&buf)
		output := buf.String()

		expectedClasses := []string{
			"text-sm",
			"font-medium",
			"text-destructive",
		}

		for _, class := range expectedClasses {
			if !strings.Contains(output, class) {
				t.Errorf("expected error to contain class: %s", class)
			}
		}
	})
}
