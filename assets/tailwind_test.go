package assets

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewTailwindProcessor(t *testing.T) {
	proc := NewTailwindProcessor()

	if proc == nil {
		t.Fatal("NewTailwindProcessor returned nil")
	}

	if proc.Name() != "TailwindCSS" {
		t.Errorf("Expected name 'TailwindCSS', got '%s'", proc.Name())
	}

	if len(proc.ContentPaths) == 0 {
		t.Error("Default content paths should not be empty")
	}

	if proc.OutputCSS == "" {
		t.Error("Default output CSS should not be empty")
	}
}

func TestTailwindProcessor_Name(t *testing.T) {
	proc := NewTailwindProcessor()

	if proc.Name() != "TailwindCSS" {
		t.Errorf("Expected name 'TailwindCSS', got '%s'", proc.Name())
	}
}

func TestTailwindProcessor_FileTypes(t *testing.T) {
	proc := NewTailwindProcessor()
	types := proc.FileTypes()

	if len(types) != 1 {
		t.Fatalf("Expected 1 file type, got %d", len(types))
	}

	if types[0] != ".css" {
		t.Errorf("Expected file type '.css', got '%s'", types[0])
	}
}

func TestTailwindProcessor_WithConfigPath(t *testing.T) {
	proc := NewTailwindProcessor()
	customPath := "/custom/tailwind.config.js"

	proc.WithConfigPath(customPath)

	if proc.ConfigPath != customPath {
		t.Errorf("Expected config path '%s', got '%s'", customPath, proc.ConfigPath)
	}
}

func TestTailwindProcessor_WithInputCSS(t *testing.T) {
	proc := NewTailwindProcessor()
	customInput := "/custom/input.css"

	proc.WithInputCSS(customInput)

	if proc.InputCSS != customInput {
		t.Errorf("Expected input CSS '%s', got '%s'", customInput, proc.InputCSS)
	}
}

func TestTailwindProcessor_WithOutputCSS(t *testing.T) {
	proc := NewTailwindProcessor()
	customOutput := "styles/main.css"

	proc.WithOutputCSS(customOutput)

	if proc.OutputCSS != customOutput {
		t.Errorf("Expected output CSS '%s', got '%s'", customOutput, proc.OutputCSS)
	}
}

func TestTailwindProcessor_WithContentPaths(t *testing.T) {
	proc := NewTailwindProcessor()
	customPaths := []string{"src/**/*.go", "components/**/*.go"}

	proc.WithContentPaths(customPaths)

	if len(proc.ContentPaths) != len(customPaths) {
		t.Errorf("Expected %d content paths, got %d", len(customPaths), len(proc.ContentPaths))
	}

	for i, path := range customPaths {
		if proc.ContentPaths[i] != path {
			t.Errorf("Expected content path '%s', got '%s'", path, proc.ContentPaths[i])
		}
	}
}

func TestTailwindProcessor_WithVerbose(t *testing.T) {
	proc := NewTailwindProcessor()

	if proc.Verbose {
		t.Error("Verbose should be false by default")
	}

	proc.WithVerbose(true)

	if !proc.Verbose {
		t.Error("Verbose should be true after setting")
	}
}

func TestTailwindProcessor_GenerateConfig(t *testing.T) {
	proc := NewTailwindProcessor()
	tempDir := t.TempDir()

	configPath, err := proc.generateConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to generate config: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	contentStr := string(content)

	// Should contain module.exports
	if !strings.Contains(contentStr, "module.exports") {
		t.Error("Config should contain module.exports")
	}

	// Should contain content paths
	if !strings.Contains(contentStr, "content:") {
		t.Error("Config should contain content paths")
	}

	// Should contain darkMode
	if !strings.Contains(contentStr, "darkMode") {
		t.Error("Config should contain darkMode setting")
	}
}

func TestTailwindProcessor_GenerateConfigWithCustomPaths(t *testing.T) {
	proc := NewTailwindProcessor()
	proc.WithContentPaths([]string{"src/**/*.go", "components/**/*.go"})

	tempDir := t.TempDir()

	configPath, err := proc.generateConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to generate config: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	contentStr := string(content)

	// Should contain custom paths
	if !strings.Contains(contentStr, "src/**/*.go") {
		t.Error("Config should contain custom content path 'src/**/*.go'")
	}

	if !strings.Contains(contentStr, "components/**/*.go") {
		t.Error("Config should contain custom content path 'components/**/*.go'")
	}
}

func TestTailwindProcessor_CreateInputCSS(t *testing.T) {
	proc := NewTailwindProcessor()
	tempDir := t.TempDir()

	inputPath, err := proc.createInputCSS(tempDir)
	if err != nil {
		t.Fatalf("Failed to create input CSS: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Error("Input CSS file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("Failed to read input CSS: %v", err)
	}

	contentStr := string(content)

	// Should contain Tailwind directives
	if !strings.Contains(contentStr, "@tailwind base") {
		t.Error("Input CSS should contain @tailwind base")
	}

	if !strings.Contains(contentStr, "@tailwind components") {
		t.Error("Input CSS should contain @tailwind components")
	}

	if !strings.Contains(contentStr, "@tailwind utilities") {
		t.Error("Input CSS should contain @tailwind utilities")
	}
}

func TestTailwindProcessor_GenerateCDNFallback(t *testing.T) {
	proc := NewTailwindProcessor()
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "app.css")

	err := proc.generateCDNFallback(outputPath)
	if err != nil {
		t.Fatalf("Failed to generate CDN fallback: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("CDN fallback file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read CDN fallback: %v", err)
	}

	contentStr := string(content)

	// Should contain CDN reference
	if !strings.Contains(contentStr, "cdn.tailwindcss.com") {
		t.Error("CDN fallback should reference Tailwind CDN")
	}

	// Should contain warning
	if !strings.Contains(contentStr, "CDN") {
		t.Error("CDN fallback should contain CDN reference")
	}
}

func TestTailwindProcessor_Process_CDNFallback(t *testing.T) {
	proc := NewTailwindProcessor()
	proc.UseCDN = true
	proc.ConfigPath = "/nonexistent/config.js" // Force CDN fallback

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	cfg := ProcessorConfig{
		InputDir:  tempDir,
		OutputDir: outputDir,
		IsDev:     true,
		Minify:    false,
	}

	ctx := context.Background()

	// This should fall back to CDN since Tailwind CLI might not be available
	err := proc.Process(ctx, cfg)

	// Should either succeed with CDN fallback or fail gracefully
	if err != nil && !strings.Contains(err.Error(), "tailwind") {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if output was created (either real CSS or CDN fallback)
	outputPath := filepath.Join(outputDir, proc.OutputCSS)
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, verify it has content
		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output: %v", err)
		}

		if len(content) == 0 {
			t.Error("Output file is empty")
		}
	}
}

func TestTailwindProcessor_Process_CreatesOutputDir(t *testing.T) {
	proc := NewTailwindProcessor()
	proc.UseCDN = true // Force CDN fallback to avoid Tailwind CLI dependency
	proc.OutputCSS = "nested/path/app.css"

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	cfg := ProcessorConfig{
		InputDir:  tempDir,
		OutputDir: outputDir,
		IsDev:     true,
	}

	ctx := context.Background()

	err := proc.Process(ctx, cfg)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Check nested directory was created
	nestedDir := filepath.Join(outputDir, "nested", "path")
	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("Nested output directory was not created")
	}
}

func TestTailwindProcessor_FluentAPI(t *testing.T) {
	proc := NewTailwindProcessor().
		WithConfigPath("/custom/config.js").
		WithInputCSS("/custom/input.css").
		WithOutputCSS("custom/output.css").
		WithContentPaths([]string{"**/*.go"}).
		WithVerbose(true)

	if proc.ConfigPath != "/custom/config.js" {
		t.Error("Fluent API failed to set ConfigPath")
	}

	if proc.InputCSS != "/custom/input.css" {
		t.Error("Fluent API failed to set InputCSS")
	}

	if proc.OutputCSS != "custom/output.css" {
		t.Error("Fluent API failed to set OutputCSS")
	}

	if !proc.Verbose {
		t.Error("Fluent API failed to set Verbose")
	}
}
