package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/xraph/forgeui/cli/templates"
	"github.com/xraph/forgeui/cli/util"
)

// TestInitWorkflow tests the complete init workflow
func TestInitWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	projectName := "test-project"
	projectDir := filepath.Join(tmpDir, projectName)
	modulePath := "github.com/test/test-project"

	// Create project structure
	if err := util.CreateDir(projectDir); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	if err := util.CreateDir(filepath.Join(projectDir, "public")); err != nil {
		t.Fatalf("Failed to create public dir: %v", err)
	}

	if err := util.CreateDir(filepath.Join(projectDir, "components")); err != nil {
		t.Fatalf("Failed to create components dir: %v", err)
	}

	if err := util.CreateDir(filepath.Join(projectDir, "pages")); err != nil {
		t.Fatalf("Failed to create pages dir: %v", err)
	}

	// Generate minimal template
	template := &templates.MinimalTemplate{}
	if err := template.Generate(projectDir, projectName, modulePath); err != nil {
		t.Fatalf("Failed to generate template: %v", err)
	}

	// Check files exist
	expectedFiles := []string{
		"main.go",
		"home.templ",
		".forgeui.json",
		".gitignore",
		"README.md",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(projectDir, file)
		if !util.FileExists(filePath) {
			t.Errorf("Expected file not created: %s", file)
		}
	}

	// Initialize go module (without network access)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "mod", "init", modulePath)

	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		t.Logf("go mod init failed (expected in test environment): %v", err)
	}

	// Verify go.mod exists
	if util.FileExists(filepath.Join(projectDir, "go.mod")) {
		t.Log("go.mod created successfully")
	}
}

// TestGenerateComponentWorkflow tests component generation
func TestGenerateComponentWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	// Setup fake project
	if err := util.CreateFile(filepath.Join(tmpDir, "go.mod"), "module github.com/test/test\n\ngo 1.25.0\n"); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	componentDir := filepath.Join(tmpDir, "components", "test_button")

	// Generate basic component
	template, err := templates.GetComponentTemplate("basic")
	if err != nil {
		t.Fatalf("Failed to get component template: %v", err)
	}

	opts := templates.ComponentOptions{
		Name:      "TestButton",
		Package:   "test_button",
		WithProps: true,
		WithTest:  true,
	}

	if err := template.Generate(componentDir, opts); err != nil {
		t.Fatalf("Failed to generate component: %v", err)
	}

	// Check files exist — templ file for component, props file for Go types
	if !util.FileExists(filepath.Join(componentDir, "test_button.templ")) {
		t.Error("Component templ file not created")
	}

	if !util.FileExists(filepath.Join(componentDir, "test_button_props.go")) {
		t.Error("Props file not created")
	}

	if !util.FileExists(filepath.Join(componentDir, "test_button_test.go")) {
		t.Error("Test file not created")
	}

	// Read and verify templ file
	data, err := os.ReadFile(filepath.Join(componentDir, "test_button.templ"))
	if err != nil {
		t.Fatalf("Failed to read component templ file: %v", err)
	}

	content := string(data)
	if content == "" {
		t.Error("Component templ file is empty")
	}

	// Check for expected content in templ file
	expectedStrings := []string{
		"package test_button",
		"templ TestButton",
	}

	for _, expected := range expectedStrings {
		if !containsString(content, expected) {
			t.Errorf("Component templ file missing expected string: %s", expected)
		}
	}

	// Check props file has expected content
	propsData, err := os.ReadFile(filepath.Join(componentDir, "test_button_props.go"))
	if err != nil {
		t.Fatalf("Failed to read props file: %v", err)
	}

	propsContent := string(propsData)
	if !containsString(propsContent, "TestButtonProps") {
		t.Error("Props file missing TestButtonProps struct")
	}
}

// TestGeneratePageWorkflow tests page generation
func TestGeneratePageWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	// Setup fake project
	if err := util.CreateFile(filepath.Join(tmpDir, "go.mod"), "module github.com/test/test\n\ngo 1.25.0\n"); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	pagesDir := filepath.Join(tmpDir, "pages")
	if err := util.CreateDir(pagesDir); err != nil {
		t.Fatalf("Failed to create pages dir: %v", err)
	}

	// Generate simple page
	template, err := templates.GetPageTemplate("simple")
	if err != nil {
		t.Fatalf("Failed to get page template: %v", err)
	}

	opts := templates.PageOptions{
		Name:    "About",
		Package: "pages",
		Path:    "/about",
	}

	pageFile := filepath.Join(pagesDir, "about.go")
	if err := template.Generate(pageFile, opts); err != nil {
		t.Fatalf("Failed to generate page: %v", err)
	}

	// Check handler Go file exists
	if !util.FileExists(pageFile) {
		t.Error("Page handler file not created")
	}

	// Check templ file exists
	templFile := filepath.Join(pagesDir, "about.templ")
	if !util.FileExists(templFile) {
		t.Error("Page templ file not created")
	}

	// Read and verify handler Go file
	data, err := os.ReadFile(pageFile)
	if err != nil {
		t.Fatalf("Failed to read page handler file: %v", err)
	}

	content := string(data)
	if content == "" {
		t.Error("Page handler file is empty")
	}

	// Check for expected content in handler file
	expectedStrings := []string{
		"package pages",
		"func About",
		"router.PageContext",
	}

	for _, expected := range expectedStrings {
		if !containsString(content, expected) {
			t.Errorf("Page file missing expected string: %s", expected)
		}
	}

	// Read and verify templ file
	templData, err := os.ReadFile(templFile)
	if err != nil {
		t.Fatalf("Failed to read page templ file: %v", err)
	}

	templContent := string(templData)
	if !containsString(templContent, "templ AboutPage") {
		t.Error("Page templ file missing templ component definition")
	}
}

// TestConfigWorkflow tests configuration save/load
func TestConfigWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config
	config := DefaultConfig()
	config.Name = "test-project"
	config.Plugins = []string{"toast", "charts"}

	// Save config
	if err := config.Save(tmpDir); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config
	if loaded.Name != config.Name {
		t.Errorf("Config name = %v, want %v", loaded.Name, config.Name)
	}

	if len(loaded.Plugins) != len(config.Plugins) {
		t.Errorf("Config plugins length = %v, want %v", len(loaded.Plugins), len(config.Plugins))
	}

	for i, plugin := range config.Plugins {
		if loaded.Plugins[i] != plugin {
			t.Errorf("Plugin %d = %v, want %v", i, loaded.Plugins[i], plugin)
		}
	}
}

// containsString checks if s contains substr (local helper for string containment)
func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexInString(s, substr) >= 0))
}

func indexInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
