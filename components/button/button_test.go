package button

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/xraph/forgeui"
)

func TestButton(t *testing.T) {
	t.Run("renders default button", func(t *testing.T) {
		btn := Button(g.Text("Click me"))

		var buf bytes.Buffer
		if err := btn.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "<button") {
			t.Error("expected <button tag")
		}
		if !strings.Contains(html, "Click me") {
			t.Error("expected button text")
		}
		if !strings.Contains(html, "inline-flex") {
			t.Error("expected base classes")
		}
	})

	t.Run("renders with variant", func(t *testing.T) {
		btn := Button(
			g.Text("Delete"),
			WithVariant(forgeui.VariantDestructive),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "bg-destructive") {
			t.Error("expected destructive variant classes")
		}
	})

	t.Run("renders with size", func(t *testing.T) {
		btn := Button(
			g.Text("Large"),
			WithSize(forgeui.SizeLG),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "h-10") {
			t.Error("expected large size classes")
		}
	})

	t.Run("renders disabled button", func(t *testing.T) {
		btn := Button(
			g.Text("Disabled"),
			Disabled(),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "disabled") {
			t.Error("expected disabled attribute")
		}
	})

	t.Run("renders loading button", func(t *testing.T) {
		btn := Button(
			g.Text("Loading"),
			Loading(),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "disabled") {
			t.Error("expected disabled attribute when loading")
		}
		if !strings.Contains(html, "aria-busy") {
			t.Error("expected aria-busy attribute")
		}
		if !strings.Contains(html, "animate-spin") {
			t.Error("expected spinner animation")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		btn := Button(
			g.Text("Custom"),
			WithClass("custom-class"),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "custom-class") {
			t.Error("expected custom class")
		}
	})

	t.Run("renders with type", func(t *testing.T) {
		btn := Button(
			g.Text("Submit"),
			WithType("submit"),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, `type="submit"`) {
			t.Error("expected submit type")
		}
	})
}

func TestButton_Variants(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(g.Node, ...Option) g.Node
		wantText string
	}{
		{"primary", Primary, "bg-primary"},
		{"secondary", Secondary, "bg-secondary"},
		{"destructive", Destructive, "bg-destructive"},
		{"outline", Outline, "border"},
		{"ghost", Ghost, "hover:bg-accent"},
		{"link", Link, "underline-offset-4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btn := tt.fn(g.Text("Test"))

			var buf bytes.Buffer
			btn.Render(&buf)
			html := buf.String()

			if !strings.Contains(html, tt.wantText) {
				t.Errorf("expected %v class for %v variant", tt.wantText, tt.name)
			}
		})
	}
}

func TestIconButton(t *testing.T) {
	btn := IconButton(g.Text("X"))

	var buf bytes.Buffer
	btn.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "size-9") {
		t.Error("expected icon size classes")
	}
}

func TestButtonGroup(t *testing.T) {
	group := Group(
		[]GroupOption{WithGap("4")},
		Button(g.Text("Save")),
		Button(g.Text("Cancel")),
	)

	var buf bytes.Buffer
	if err := group.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "flex") {
		t.Error("expected flex class")
	}
	if !strings.Contains(html, "gap-4") {
		t.Error("expected gap class")
	}
	if !strings.Contains(html, "Save") || !strings.Contains(html, "Cancel") {
		t.Error("expected button content")
	}
}

func TestButton_AllOptions(t *testing.T) {
	btn := Button(
		g.Text("Full Options"),
		WithVariant(forgeui.VariantSecondary),
		WithSize(forgeui.SizeSM),
		WithType("reset"),
		WithClass("extra-class"),
		WithAttrs(g.Attr("data-test", "value")),
		Disabled(),
	)

	var buf bytes.Buffer
	btn.Render(&buf)
	html := buf.String()

	expected := []string{
		"bg-secondary",
		"h-8",
		`type="reset"`,
		"extra-class",
		`data-test="value"`,
		"disabled",
	}

	for _, exp := range expected {
		if !strings.Contains(html, exp) {
			t.Errorf("expected %v in button with all options", exp)
		}
	}
}
