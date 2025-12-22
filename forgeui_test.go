package forgeui_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
)

func TestIntegration_NodeWithCVA(t *testing.T) {
	// Test that Node and CVA work together
	cva := forgeui.NewCVA("btn", "rounded").
		Variant("size", map[string][]string{
			"sm": {"h-8", "px-3", "text-sm"},
			"lg": {"h-12", "px-6", "text-lg"},
		}).
		Default("size", "sm")

	classes := cva.Classes(map[string]string{"size": "lg"})

	node := forgeui.El("button").
		Class(classes).
		Attr("type", "button").
		Children(g.Text("Click me")).
		Build()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()

	// Verify all expected parts are present
	expected := []string{
		"<button",
		"btn",
		"rounded",
		"h-12",
		"px-6",
		"text-lg",
		`type="button"`,
		"Click me",
	}

	for _, want := range expected {
		if !strings.Contains(html, want) {
			t.Errorf("expected %v in HTML output", want)
		}
	}
}

func TestIntegration_AppInitialization(t *testing.T) {
	// Test that App initializes correctly with config
	app := forgeui.New(
		forgeui.WithDebug(true),
		forgeui.WithThemeName("dark"),
		forgeui.WithStaticPath("/public"),
	)

	if err := app.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	cfg := app.Config()
	if !cfg.Debug {
		t.Error("expected Debug to be true")
	}

	if cfg.Theme != "dark" {
		t.Errorf("expected Theme to be 'dark', got %v", cfg.Theme)
	}

	if !app.IsDev() {
		t.Error("IsDev() should return true when Debug is true")
	}
}

func TestIntegration_UtilsWithCVA(t *testing.T) {
	// Test that utility functions work with CVA
	cva := forgeui.NewCVA("base").
		Variant("active", map[string][]string{
			"true":  {"bg-blue-500"},
			"false": {"bg-gray-500"},
		}).
		Default("active", "false")

	isActive := true
	classes := cva.Classes(map[string]string{
		"active": forgeui.IfElse(isActive, "true", "false"),
	})

	finalClasses := forgeui.CN(classes, forgeui.If(isActive, "font-bold"))

	if !strings.Contains(finalClasses, "bg-blue-500") {
		t.Error("expected bg-blue-500 for active state")
	}

	if !strings.Contains(finalClasses, "font-bold") {
		t.Error("expected font-bold for active state")
	}
}

func TestIntegration_ComplexComponent(t *testing.T) {
	// Test building a complex component with all features
	buttonCVA := forgeui.NewCVA(
		"inline-flex", "items-center", "justify-center", "rounded-md",
		"font-medium", "transition-colors",
	).
		Variant("variant", map[string][]string{
			"default":     {"bg-primary", "text-primary-foreground"},
			"destructive": {"bg-destructive", "text-destructive-foreground"},
		}).
		Variant("size", map[string][]string{
			"sm": {"h-9", "px-3", "text-sm"},
			"lg": {"h-11", "px-8", "text-lg"},
		}).
		Default("variant", "default").
		Default("size", "sm")

	isDisabled := false
	variant := "destructive"
	size := "lg"

	classes := buttonCVA.Classes(map[string]string{
		"variant": variant,
		"size":    size,
	})

	finalClasses := forgeui.CN(classes, forgeui.If(isDisabled, "opacity-50 cursor-not-allowed"))

	button := forgeui.El("button").
		Class(finalClasses).
		Attr("type", "button").
		Children(g.Text("Delete")).
		Build()

	var buf bytes.Buffer
	if err := button.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()

	expected := []string{
		"inline-flex",
		"bg-destructive",
		"h-11",
		"px-8",
		"Delete",
	}

	for _, want := range expected {
		if !strings.Contains(html, want) {
			t.Errorf("expected %v in HTML output, got: %s", want, html)
		}
	}

	// Should not contain disabled classes
	if strings.Contains(html, "opacity-50") {
		t.Error("should not contain opacity-50 when not disabled")
	}
}
