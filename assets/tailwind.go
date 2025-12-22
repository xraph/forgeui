package assets

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xraph/forgeui/theme"
)

// TailwindProcessor processes Tailwind CSS by scanning Go files for class usage
// and generating optimized CSS output.
type TailwindProcessor struct {
	// ConfigPath is the path to tailwind.config.js (auto-generated if empty)
	ConfigPath string

	// InputCSS is the path to the input CSS file with @tailwind directives
	InputCSS string

	// OutputCSS is the path where processed CSS will be written
	OutputCSS string

	// ContentPaths are glob patterns for files to scan for classes
	// Defaults to ["**/*.go"] to scan all Go files
	ContentPaths []string

	// UseCDN falls back to CDN if Tailwind CLI is not available
	UseCDN bool

	// Verbose enables detailed logging
	Verbose bool
}

// NewTailwindProcessor creates a new Tailwind CSS processor with sensible defaults
func NewTailwindProcessor() *TailwindProcessor {
	return &TailwindProcessor{
		ContentPaths: []string{"**/*.go"},
		OutputCSS:    "css/app.css",
		UseCDN:       true,
	}
}

// Name returns the processor name
func (tp *TailwindProcessor) Name() string {
	return "TailwindCSS"
}

// FileTypes returns the file extensions this processor handles
func (tp *TailwindProcessor) FileTypes() []string {
	return []string{".css"}
}

// Process executes the Tailwind CSS processing
func (tp *TailwindProcessor) Process(ctx context.Context, cfg ProcessorConfig) error {
	// Generate tailwind.config.js if not provided
	configPath := tp.ConfigPath
	if configPath == "" {
		var err error
		configPath, err = tp.generateConfig(cfg.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to generate tailwind config: %w", err)
		}
		defer os.Remove(configPath) // Clean up after processing
	}

	// Create input CSS if not provided
	inputCSS := tp.InputCSS
	if inputCSS == "" {
		var err error
		inputCSS, err = tp.createInputCSS(cfg.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to create input CSS: %w", err)
		}
		defer os.Remove(inputCSS) // Clean up after processing
	}

	// Determine output path
	outputCSS := filepath.Join(cfg.OutputDir, tp.OutputCSS)

	// Ensure output directory exists
	outputDir := filepath.Dir(outputCSS)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if Tailwind CLI is available
	if !tp.isTailwindAvailable() {
		if tp.UseCDN {
			return tp.generateCDNFallback(outputCSS)
		}
		return fmt.Errorf("tailwind CLI not found and CDN fallback disabled")
	}

	// Build Tailwind CSS command
	args := []string{
		"tailwindcss",
		"-c", configPath,
		"-i", inputCSS,
		"-o", outputCSS,
	}

	// Add minification in production
	if cfg.Minify {
		args = append(args, "--minify")
	}

	// Execute Tailwind CLI
	cmd := exec.CommandContext(ctx, "npx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if tp.Verbose {
		fmt.Printf("[Tailwind] Running: npx %s\n", strings.Join(args, " "))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tailwind build failed: %w", err)
	}

	if tp.Verbose {
		fmt.Printf("[Tailwind] Generated CSS: %s\n", outputCSS)
	}

	return nil
}

// generateConfig creates a tailwind.config.js file
func (tp *TailwindProcessor) generateConfig(outputDir string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	configPath := filepath.Join(outputDir, "tailwind.config.js")

	// Use ForgeUI's theme configuration
	config := theme.GenerateTailwindConfig()

	// Add content paths
	contentPaths := tp.ContentPaths
	if len(contentPaths) == 0 {
		contentPaths = []string{"**/*.go"}
	}

	// Build content array
	contentJSON := "["
	for i, path := range contentPaths {
		if i > 0 {
			contentJSON += ", "
		}
		contentJSON += fmt.Sprintf(`"%s"`, path)
	}
	contentJSON += "]"

	// Insert content paths into config
	config = strings.Replace(config, "module.exports = {", 
		fmt.Sprintf("module.exports = {\n  content: %s,", contentJSON), 1)

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return "", err
	}

	if tp.Verbose {
		fmt.Printf("[Tailwind] Generated config: %s\n", configPath)
	}

	return configPath, nil
}

// createInputCSS creates an input CSS file with Tailwind directives
func (tp *TailwindProcessor) createInputCSS(outputDir string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	inputPath := filepath.Join(outputDir, "input.css")

	content := `@tailwind base;
@tailwind components;
@tailwind utilities;
`

	if err := os.WriteFile(inputPath, []byte(content), 0644); err != nil {
		return "", err
	}

	return inputPath, nil
}

// isTailwindAvailable checks if Tailwind CLI is available
func (tp *TailwindProcessor) isTailwindAvailable() bool {
	cmd := exec.Command("npx", "tailwindcss", "--help")
	err := cmd.Run()
	return err == nil
}

// generateCDNFallback creates a minimal CSS file that references Tailwind CDN
func (tp *TailwindProcessor) generateCDNFallback(outputPath string) error {
	content := `/* Tailwind CSS via CDN */
/* Add this to your HTML: <script src="https://cdn.tailwindcss.com"></script> */

/* ForgeUI: Tailwind CLI not found. Using CDN fallback. */
/* For production, install Tailwind: npm install -D tailwindcss */

:root {
  /* ForgeUI theme variables will be injected by theme system */
}
`

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write CDN fallback: %w", err)
	}

	fmt.Println("[Tailwind] WARNING: Tailwind CLI not found. Generated CDN fallback.")
	fmt.Println("[Tailwind] Install with: npm install -D tailwindcss")
	fmt.Println("[Tailwind] Or use CDN: <script src=\"https://cdn.tailwindcss.com\"></script>")

	return nil
}

// WithConfigPath sets a custom tailwind.config.js path
func (tp *TailwindProcessor) WithConfigPath(path string) *TailwindProcessor {
	tp.ConfigPath = path
	return tp
}

// WithInputCSS sets a custom input CSS file
func (tp *TailwindProcessor) WithInputCSS(path string) *TailwindProcessor {
	tp.InputCSS = path
	return tp
}

// WithOutputCSS sets the output CSS file path (relative to output directory)
func (tp *TailwindProcessor) WithOutputCSS(path string) *TailwindProcessor {
	tp.OutputCSS = path
	return tp
}

// WithContentPaths sets custom content paths to scan for classes
func (tp *TailwindProcessor) WithContentPaths(paths []string) *TailwindProcessor {
	tp.ContentPaths = paths
	return tp
}

// WithVerbose enables verbose logging
func (tp *TailwindProcessor) WithVerbose(verbose bool) *TailwindProcessor {
	tp.Verbose = verbose
	return tp
}

