package separator

import (
	"bytes"
	"strings"
	"testing"
)

func TestSeparator_Horizontal(t *testing.T) {
	sep := Separator()

	var buf bytes.Buffer
	if err := sep.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	html := buf.String()

	if !strings.Contains(html, "h-[1px]") {
		t.Error("expected horizontal height")
	}

	if !strings.Contains(html, "w-full") {
		t.Error("expected full width")
	}

	if !strings.Contains(html, `aria-orientation="horizontal"`) {
		t.Error("expected horizontal orientation")
	}
}

func TestSeparator_Vertical(t *testing.T) {
	sep := Separator(Vertical())

	var buf bytes.Buffer
	if err := sep.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	html := buf.String()

	if !strings.Contains(html, "h-full") {
		t.Error("expected full height")
	}

	if !strings.Contains(html, "w-[1px]") {
		t.Error("expected vertical width")
	}

	if !strings.Contains(html, `aria-orientation="vertical"`) {
		t.Error("expected vertical orientation")
	}
}
