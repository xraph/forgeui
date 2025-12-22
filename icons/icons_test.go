package icons

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestIcon(t *testing.T) {
	t.Run("renders basic icon", func(t *testing.T) {
		icon := Icon("M5 12h14")

		var buf bytes.Buffer
		if err := icon.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "<svg") {
			t.Error("expected svg element")
		}
		if !strings.Contains(html, "M5 12h14") {
			t.Error("expected path data")
		}
		if !strings.Contains(html, `viewBox="0 0 24 24"`) {
			t.Error("expected viewBox attribute")
		}
	})

	t.Run("renders with custom size", func(t *testing.T) {
		icon := Icon("M5 12h14", WithSize(32))

		var buf bytes.Buffer
		icon.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, `width="32"`) {
			t.Error("expected width=32")
		}
		if !strings.Contains(html, `height="32"`) {
			t.Error("expected height=32")
		}
	})

	t.Run("renders with custom color", func(t *testing.T) {
		icon := Icon("M5 12h14", WithColor("red"))

		var buf bytes.Buffer
		icon.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, `stroke="red"`) {
			t.Error("expected stroke=red")
		}
	})

	t.Run("renders with custom stroke width", func(t *testing.T) {
		icon := Icon("M5 12h14", WithStrokeWidth(3.0))

		var buf bytes.Buffer
		icon.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, `stroke-width="3"`) {
			t.Error("expected stroke-width=3")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		icon := Icon("M5 12h14", WithClass("custom-icon"))

		var buf bytes.Buffer
		icon.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "custom-icon") {
			t.Error("expected custom-icon class")
		}
	})

	t.Run("renders with all options", func(t *testing.T) {
		icon := Icon(
			"M5 12h14",
			WithSize(20),
			WithColor("blue"),
			WithStrokeWidth(2.5),
			WithClass("my-icon"),
			WithAttrs(g.Attr("data-test", "icon")),
		)

		var buf bytes.Buffer
		icon.Render(&buf)
		html := buf.String()

		expected := []string{
			`width="20"`,
			`height="20"`,
			`stroke="blue"`,
			`stroke-width="2.5"`,
			"my-icon",
			`data-test="icon"`,
		}

		for _, exp := range expected {
			if !strings.Contains(html, exp) {
				t.Errorf("expected %v in icon with all options", exp)
			}
		}
	})
}

func TestMultiPathIcon(t *testing.T) {
	t.Run("renders multiple paths", func(t *testing.T) {
		icon := MultiPathIcon([]string{
			"M18 6 6 18",
			"m6 6 12 12",
		})

		var buf bytes.Buffer
		icon.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "M18 6 6 18") {
			t.Error("expected first path")
		}
		if !strings.Contains(html, "m6 6 12 12") {
			t.Error("expected second path")
		}
	})
}

func TestLucideIcons(t *testing.T) {
	tests := []struct {
		name string
		icon func(...Option) g.Node
	}{
		{"Check", Check},
		{"X", X},
		{"ChevronDown", ChevronDown},
		{"ChevronUp", ChevronUp},
		{"ChevronLeft", ChevronLeft},
		{"ChevronRight", ChevronRight},
		{"Plus", Plus},
		{"Minus", Minus},
		{"Search", Search},
		{"Menu", Menu},
		{"User", User},
		{"Home", Home},
		{"Mail", Mail},
		{"AlertCircle", AlertCircle},
		{"Info", Info},
		{"CheckCircle", CheckCircle},
		{"XCircle", XCircle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon := tt.icon()

			var buf bytes.Buffer
			if err := icon.Render(&buf); err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			html := buf.String()
			// Verify it's a valid SVG with path data
			if !strings.Contains(html, "<svg") {
				t.Error("expected svg element")
			}
			if !strings.Contains(html, "<path") && !strings.Contains(html, "d=\"") {
				t.Error("expected svg path element with d attribute")
			}
			// Verify basic SVG attributes
			if !strings.Contains(html, "xmlns") {
				t.Error("expected xmlns attribute")
			}
			if !strings.Contains(html, "viewBox") {
				t.Error("expected viewBox attribute")
			}
		})
	}
}

func TestLucideIconsWithOptions(t *testing.T) {
	t.Run("renders with custom options", func(t *testing.T) {
		icon := Check(
			WithSize(32),
			WithColor("green"),
			WithClass("success-icon"),
		)

		var buf bytes.Buffer
		icon.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, `width="32"`) {
			t.Error("expected width=32")
		}
		if !strings.Contains(html, `stroke="green"`) {
			t.Error("expected stroke=green")
		}
		if !strings.Contains(html, "success-icon") {
			t.Error("expected success-icon class")
		}
	})
}

