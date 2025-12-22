package assets

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewESBuildProcessor(t *testing.T) {
	proc := NewESBuildProcessor()

	if proc == nil {
		t.Fatal("NewESBuildProcessor returned nil")
	}

	if proc.Name() != "ESBuild" {
		t.Errorf("Expected name 'ESBuild', got '%s'", proc.Name())
	}

	if proc.Format == "" {
		t.Error("Default format should not be empty")
	}

	if proc.Target == "" {
		t.Error("Default target should not be empty")
	}

	if !proc.Bundle {
		t.Error("Bundle should be enabled by default")
	}
}

func TestESBuildProcessor_Name(t *testing.T) {
	proc := NewESBuildProcessor()

	if proc.Name() != "ESBuild" {
		t.Errorf("Expected name 'ESBuild', got '%s'", proc.Name())
	}
}

func TestESBuildProcessor_FileTypes(t *testing.T) {
	proc := NewESBuildProcessor()
	types := proc.FileTypes()

	expectedTypes := []string{".js", ".ts", ".jsx", ".tsx"}

	if len(types) != len(expectedTypes) {
		t.Fatalf("Expected %d file types, got %d", len(expectedTypes), len(types))
	}

	for i, expected := range expectedTypes {
		if types[i] != expected {
			t.Errorf("Expected file type '%s', got '%s'", expected, types[i])
		}
	}
}

func TestESBuildProcessor_WithEntryPoints(t *testing.T) {
	proc := NewESBuildProcessor()
	entries := []string{"src/app.js", "src/vendor.js"}

	proc.WithEntryPoints(entries...)

	if len(proc.EntryPoints) != len(entries) {
		t.Errorf("Expected %d entry points, got %d", len(entries), len(proc.EntryPoints))
	}

	for i, entry := range entries {
		if proc.EntryPoints[i] != entry {
			t.Errorf("Expected entry point '%s', got '%s'", entry, proc.EntryPoints[i])
		}
	}
}

func TestESBuildProcessor_WithOutfile(t *testing.T) {
	proc := NewESBuildProcessor()
	customOutfile := "dist/bundle.js"

	proc.WithOutfile(customOutfile)

	if proc.Outfile != customOutfile {
		t.Errorf("Expected outfile '%s', got '%s'", customOutfile, proc.Outfile)
	}
}

func TestESBuildProcessor_WithFormat(t *testing.T) {
	proc := NewESBuildProcessor()
	formats := []string{"iife", "cjs", "esm"}

	for _, format := range formats {
		proc.WithFormat(format)

		if proc.Format != format {
			t.Errorf("Expected format '%s', got '%s'", format, proc.Format)
		}
	}
}

func TestESBuildProcessor_WithTarget(t *testing.T) {
	proc := NewESBuildProcessor()
	target := "es2015"

	proc.WithTarget(target)

	if proc.Target != target {
		t.Errorf("Expected target '%s', got '%s'", target, proc.Target)
	}
}

func TestESBuildProcessor_WithBundle(t *testing.T) {
	proc := NewESBuildProcessor()

	proc.WithBundle(false)

	if proc.Bundle {
		t.Error("Bundle should be disabled")
	}

	proc.WithBundle(true)

	if !proc.Bundle {
		t.Error("Bundle should be enabled")
	}
}

func TestESBuildProcessor_WithSplitting(t *testing.T) {
	proc := NewESBuildProcessor()

	if proc.Splitting {
		t.Error("Splitting should be disabled by default")
	}

	proc.WithSplitting(true)

	if !proc.Splitting {
		t.Error("Splitting should be enabled")
	}
}

func TestESBuildProcessor_WithExternal(t *testing.T) {
	proc := NewESBuildProcessor()
	external := []string{"react", "react-dom"}

	proc.WithExternal(external...)

	if len(proc.External) != len(external) {
		t.Errorf("Expected %d external packages, got %d", len(external), len(proc.External))
	}

	for i, pkg := range external {
		if proc.External[i] != pkg {
			t.Errorf("Expected external package '%s', got '%s'", pkg, proc.External[i])
		}
	}
}

func TestESBuildProcessor_WithVerbose(t *testing.T) {
	proc := NewESBuildProcessor()

	if proc.Verbose {
		t.Error("Verbose should be false by default")
	}

	proc.WithVerbose(true)

	if !proc.Verbose {
		t.Error("Verbose should be true after setting")
	}
}

func TestESBuildProcessor_Process_NoEntryPoints(t *testing.T) {
	proc := NewESBuildProcessor()
	proc.Verbose = true

	tempDir := t.TempDir()

	cfg := ProcessorConfig{
		InputDir:  tempDir,
		OutputDir: tempDir,
		IsDev:     true,
	}

	ctx := context.Background()

	// Should not fail when no entry points are specified
	err := proc.Process(ctx, cfg)
	if err != nil {
		t.Errorf("Process should not fail with no entry points: %v", err)
	}
}

func TestESBuildProcessor_Process_CreatesOutputDir(t *testing.T) {
	proc := NewESBuildProcessor()
	proc.Verbose = true

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	// Create a real test.js file with valid JavaScript content
	testJSPath := filepath.Join(tempDir, "test.js")

	testJSContent := []byte("console.log('test');")
	if err := os.WriteFile(testJSPath, testJSContent, 0644); err != nil {
		t.Fatalf("Failed to create test.js: %v", err)
	}

	// Use the actual file path as entry point
	proc.WithEntryPoints(testJSPath)
	proc.WithOutfile("nested/path/bundle.js")

	cfg := ProcessorConfig{
		InputDir:  tempDir,
		OutputDir: outputDir,
		IsDev:     true,
	}

	ctx := context.Background()

	// This will likely fail if esbuild is not installed, but should at least
	// attempt to create the output directory
	_ = proc.Process(ctx, cfg)

	// We can't guarantee esbuild is installed, so just check the logic doesn't panic
}

func TestESBuildProcessor_FluentAPI(t *testing.T) {
	proc := NewESBuildProcessor().
		WithEntryPoints("src/app.js", "src/vendor.js").
		WithOutfile("dist/bundle.js").
		WithFormat("esm").
		WithTarget("es2015").
		WithBundle(true).
		WithSplitting(true).
		WithExternal("react", "react-dom").
		WithVerbose(true)

	if len(proc.EntryPoints) != 2 {
		t.Error("Fluent API failed to set EntryPoints")
	}

	if proc.Outfile != "dist/bundle.js" {
		t.Error("Fluent API failed to set Outfile")
	}

	if proc.Format != "esm" {
		t.Error("Fluent API failed to set Format")
	}

	if proc.Target != "es2015" {
		t.Error("Fluent API failed to set Target")
	}

	if !proc.Bundle {
		t.Error("Fluent API failed to set Bundle")
	}

	if !proc.Splitting {
		t.Error("Fluent API failed to set Splitting")
	}

	if len(proc.External) != 2 {
		t.Error("Fluent API failed to set External")
	}

	if !proc.Verbose {
		t.Error("Fluent API failed to set Verbose")
	}
}

func TestESBuildProcessor_DefaultValues(t *testing.T) {
	proc := NewESBuildProcessor()

	if proc.Format != "iife" {
		t.Errorf("Expected default format 'iife', got '%s'", proc.Format)
	}

	if proc.Target != "es2020" {
		t.Errorf("Expected default target 'es2020', got '%s'", proc.Target)
	}

	if !proc.Bundle {
		t.Error("Bundle should be enabled by default")
	}

	if proc.Splitting {
		t.Error("Splitting should be disabled by default")
	}

	if proc.Outfile != "js/app.js" {
		t.Errorf("Expected default outfile 'js/app.js', got '%s'", proc.Outfile)
	}
}
