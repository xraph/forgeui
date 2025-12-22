package forgeui

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestNode_El(t *testing.T) {
	node := El("div").Build()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "<div") {
		t.Errorf("expected <div tag, got %s", html)
	}
}

func TestNode_Class(t *testing.T) {
	node := El("div").
		Class("class-1", "class-2").
		Class("class-3").
		Build()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "class-1") || !strings.Contains(html, "class-2") || !strings.Contains(html, "class-3") {
		t.Errorf("expected all classes, got %s", html)
	}
}

func TestNode_Attr(t *testing.T) {
	node := El("input").
		Attr("type", "text").
		Attr("placeholder", "Enter text").
		Build()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `type="text"`) {
		t.Errorf("expected type attribute, got %s", html)
	}

	if !strings.Contains(html, `placeholder="Enter text"`) {
		t.Errorf("expected placeholder attribute, got %s", html)
	}
}

func TestNode_Children(t *testing.T) {
	node := El("div").
		Children(g.Text("Hello"), g.Text(" World")).
		Build()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Hello World") {
		t.Errorf("expected 'Hello World', got %s", html)
	}
}

func TestNode_ChainedMethods(t *testing.T) {
	node := El("button").
		Class("btn", "btn-primary").
		Attr("type", "button").
		Attr("disabled", "").
		Children(g.Text("Click me")).
		Build()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()

	tests := []string{
		"<button",
		"btn",
		"btn-primary",
		`type="button"`,
		`disabled=""`,
		"Click me",
	}

	for _, want := range tests {
		if !strings.Contains(html, want) {
			t.Errorf("expected %v in %v", want, html)
		}
	}
}

func TestNode_EmptyClasses(t *testing.T) {
	node := El("div").Build()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()
	// Should render div without class attribute if no classes
	if strings.Contains(html, "class=") {
		t.Errorf("expected no class attribute, got %s", html)
	}
}
