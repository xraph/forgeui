package form

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
)

func TestDescription(t *testing.T) {
	t.Run("renders description text", func(t *testing.T) {
		desc := Description("Helper text for this field")

		var buf bytes.Buffer
		desc.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Helper text for this field") {
			t.Error("expected description to contain text")
		}
		if !strings.Contains(output, "text-muted-foreground") {
			t.Error("expected description to have muted color")
		}
	})

	t.Run("returns nil for empty description", func(t *testing.T) {
		desc := Description("")

		if desc != nil {
			t.Error("expected nil for empty description")
		}
	})

	t.Run("renders description with custom class", func(t *testing.T) {
		desc := Description(
			"Helper text",
			WithDescriptionClass("mt-4"),
		)

		var buf bytes.Buffer
		desc.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "mt-4") {
			t.Error("expected description to contain custom class")
		}
	})

	t.Run("renders description with custom attributes", func(t *testing.T) {
		desc := Description(
			"Helper text",
			WithDescriptionAttrs(g.Attr("data-testid", "field-description")),
		)

		var buf bytes.Buffer
		desc.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-testid=\"field-description\"") {
			t.Error("expected description to contain custom attribute")
		}
	})

	t.Run("has correct base classes", func(t *testing.T) {
		desc := Description("Helper text")

		var buf bytes.Buffer
		desc.Render(&buf)
		output := buf.String()

		expectedClasses := []string{
			"text-sm",
			"text-muted-foreground",
		}

		for _, class := range expectedClasses {
			if !strings.Contains(output, class) {
				t.Errorf("expected description to contain class: %s", class)
			}
		}
	})
}
