package form

import (
	"bytes"
	"strings"
	"testing"

	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/components/input"
)

func TestField(t *testing.T) {
	t.Run("renders basic field", func(t *testing.T) {
		field := Field(
			"Email",
			input.Input(input.WithType("email")),
		)

		var buf bytes.Buffer
		field.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Email") {
			t.Error("expected field to contain label text")
		}
		if !strings.Contains(output, "space-y-2") {
			t.Error("expected field to have spacing classes")
		}
	})

	t.Run("renders field with ID", func(t *testing.T) {
		field := Field(
			"Email",
			input.Input(input.WithID("email")),
			WithID("email"),
		)

		var buf bytes.Buffer
		field.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "email") {
			t.Error("expected field to contain ID")
		}
	})

	t.Run("renders required field", func(t *testing.T) {
		field := Field(
			"Email",
			input.Input(),
			WithRequired(),
		)

		var buf bytes.Buffer
		field.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "*") {
			t.Error("expected required field to have asterisk")
		}
	})

	t.Run("renders field with description", func(t *testing.T) {
		field := Field(
			"Email",
			input.Input(),
			WithDescription("We will never share your email"),
		)

		var buf bytes.Buffer
		field.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "We will never share your email") {
			t.Error("expected field to contain description")
		}
		if !strings.Contains(output, "text-muted-foreground") {
			t.Error("expected description to have muted text color")
		}
	})

	t.Run("renders field with error", func(t *testing.T) {
		field := Field(
			"Email",
			input.Input(),
			WithError("Invalid email address"),
		)

		var buf bytes.Buffer
		field.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Invalid email address") {
			t.Error("expected field to contain error message")
		}
		if !strings.Contains(output, "text-destructive") {
			t.Error("expected error to have destructive color")
		}
		if !strings.Contains(output, "role=\"alert\"") {
			t.Error("expected error to have alert role")
		}
	})

	t.Run("renders field with custom class", func(t *testing.T) {
		field := Field(
			"Email",
			input.Input(),
			WithFieldClass("custom-field"),
		)

		var buf bytes.Buffer
		field.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-field") {
			t.Error("expected field to contain custom class")
		}
	})
}

func TestFieldLabel(t *testing.T) {
	t.Run("renders label", func(t *testing.T) {
		label := FieldLabel("Username", "username", false)

		var buf bytes.Buffer
		label.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Username") {
			t.Error("expected label to contain text")
		}
	})

	t.Run("renders required label", func(t *testing.T) {
		label := FieldLabel("Username", "username", true)

		var buf bytes.Buffer
		label.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "*") {
			t.Error("expected required label to have asterisk")
		}
	})

	t.Run("renders label with for attribute", func(t *testing.T) {
		label := FieldLabel("Username", "username", false)

		var buf bytes.Buffer
		label.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "for=\"username\"") {
			t.Error("expected label to have for attribute")
		}
	})
}

func TestFieldControl(t *testing.T) {
	t.Run("returns control as-is", func(t *testing.T) {
		control := html.Input(html.Type("text"))
		result := FieldControl(control)

		if result == nil {
			t.Error("expected FieldControl to return a non-nil node")
		}

		// Verify by rendering both and comparing output
		var controlBuf, resultBuf strings.Builder
		_ = control.Render(&controlBuf)
		_ = result.Render(&resultBuf)

		if controlBuf.String() != resultBuf.String() {
			t.Errorf("expected FieldControl to return control unchanged\ncontrol: %s\nresult:  %s",
				controlBuf.String(), resultBuf.String())
		}
	})
}

func TestFieldError(t *testing.T) {
	t.Run("renders error message", func(t *testing.T) {
		err := FieldError("This field is required")

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
	})

	t.Run("returns nil for empty error", func(t *testing.T) {
		err := FieldError("")

		if err != nil {
			t.Error("expected nil for empty error message")
		}
	})
}

func TestFieldDescription(t *testing.T) {
	t.Run("renders description", func(t *testing.T) {
		desc := FieldDescription("Helper text")

		var buf bytes.Buffer
		desc.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Helper text") {
			t.Error("expected description to contain text")
		}
		if !strings.Contains(output, "text-muted-foreground") {
			t.Error("expected description to have muted color")
		}
	})

	t.Run("returns nil for empty description", func(t *testing.T) {
		desc := FieldDescription("")

		if desc != nil {
			t.Error("expected nil for empty description")
		}
	})
}

func TestFieldIntegration(t *testing.T) {
	t.Run("renders complete field with all elements", func(t *testing.T) {
		field := Field(
			"Email Address",
			input.Input(
				input.WithType("email"),
				input.WithID("email"),
				input.WithPlaceholder("you@example.com"),
			),
			WithID("email"),
			WithName("email"),
			WithRequired(),
			WithDescription("We will never share your email with anyone else."),
			WithError("Please enter a valid email address"),
		)

		var buf bytes.Buffer
		field.Render(&buf)
		output := buf.String()

		// Check all elements are present
		if !strings.Contains(output, "Email Address") {
			t.Error("expected label text")
		}
		if !strings.Contains(output, "*") {
			t.Error("expected required indicator")
		}
		if !strings.Contains(output, "We will never share your email") {
			t.Error("expected description")
		}
		if !strings.Contains(output, "Please enter a valid email address") {
			t.Error("expected error message")
		}
		if !strings.Contains(output, "you@example.com") {
			t.Error("expected placeholder")
		}
	})
}
