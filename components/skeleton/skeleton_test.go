package skeleton

import (
	"bytes"
	"strings"
	"testing"
)

func TestSkeleton(t *testing.T) {
	skeleton := Skeleton()

	var buf bytes.Buffer
	if err := skeleton.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	html := buf.String()

	if !strings.Contains(html, "animate-pulse") {
		t.Error("expected animate-pulse class")
	}

	if !strings.Contains(html, "bg-muted") {
		t.Error("expected bg-muted class")
	}
}

func TestSkeleton_WithDimensions(t *testing.T) {
	skeleton := Skeleton(
		WithWidth("w-full"),
		WithHeight("h-12"),
	)

	var buf bytes.Buffer
	if err := skeleton.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	html := buf.String()

	if !strings.Contains(html, "w-full") {
		t.Error("expected w-full class")
	}

	if !strings.Contains(html, "h-12") {
		t.Error("expected h-12 class")
	}
}
