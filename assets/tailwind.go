package assets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/xraph/forgeui/theme"
)

// TailwindVersion represents the Tailwind CSS major version.
type TailwindVersion int

const (
	// TailwindV3 uses JS-based config, @tailwind directives, and `npx tailwindcss` CLI.
	TailwindV3 TailwindVersion = 3

	// TailwindV4 uses CSS-based config, @import "tailwindcss", and `npx @tailwindcss/cli` CLI.
	TailwindV4 TailwindVersion = 4
)

// TailwindProcessor processes Tailwind CSS by scanning Go files for class usage
// and generating optimized CSS output. Supports both Tailwind v3 and v4.
type TailwindProcessor struct {
	// Version selects the Tailwind CSS major version (3 or 4). Default: TailwindV4.
	Version TailwindVersion

	// ConfigPath is the path to tailwind.config.js (v3 only, auto-generated if empty)
	ConfigPath string

	// InputCSS is the path to the input CSS file.
	// For v3: contains @tailwind directives (auto-generated if empty).
	// For v4: contains @import "tailwindcss" + @theme inline (auto-generated from themes if empty).
	InputCSS string

	// OutputCSS is the path where processed CSS will be written (relative to output directory)
	OutputCSS string

	// ContentPaths are glob patterns for files to scan for classes (v3 only).
	// Defaults to ["**/*.go"] to scan all Go files.
	ContentPaths []string

	// UseCDN falls back to CDN if Tailwind CLI is not available
	UseCDN bool

	// Verbose enables detailed logging
	Verbose bool

	// LightTheme is the light theme for generating v4 input CSS.
	// When set together with DarkTheme, the processor auto-generates a complete
	// Tailwind v4 input.css from theme tokens using theme.GenerateInputCSS().
	LightTheme *theme.Theme

	// DarkTheme is the dark theme for generating v4 input CSS.
	DarkTheme *theme.Theme
}

// NewTailwindProcessor creates a new Tailwind CSS processor with sensible defaults.
// Default version is TailwindV4 (CSS-based configuration).
func NewTailwindProcessor() *TailwindProcessor {
	return &TailwindProcessor{
		Version:      TailwindV4,
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

// Process executes the Tailwind CSS processing.
// For v4: generates CSS-based config, runs @tailwindcss/cli.
// For v3: generates JS config, runs tailwindcss CLI.
func (tp *TailwindProcessor) Process(ctx context.Context, cfg ProcessorConfig) error {
	if tp.Version == TailwindV3 {
		return tp.processV3(ctx, cfg)
	}

	return tp.processV4(ctx, cfg)
}

// processV4 handles Tailwind CSS v4 processing.
// v4 uses CSS-based configuration — no tailwind.config.js needed.
func (tp *TailwindProcessor) processV4(ctx context.Context, cfg ProcessorConfig) error {
	// Create input CSS if not provided
	inputCSS := tp.InputCSS
	if inputCSS == "" {
		var err error

		inputCSS, err = tp.createInputCSSv4(cfg.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to create v4 input CSS: %w", err)
		}

		defer func() { _ = os.Remove(inputCSS) }()
	}

	// Determine output path
	outputCSS := filepath.Join(cfg.OutputDir, tp.OutputCSS)

	// Ensure output directory exists
	outputDir := filepath.Dir(outputCSS)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if Tailwind v4 CLI is available
	cliCmd := tp.findTailwindV4CLI()
	if cliCmd == "" {
		if tp.UseCDN {
			return tp.generateCDNFallback(outputCSS)
		}

		return errors.New("tailwind v4 CLI not found (@tailwindcss/cli) and CDN fallback disabled. Install with: npm install -D @tailwindcss/cli")
	}

	// Build Tailwind v4 command: npx @tailwindcss/cli -i input.css -o output.css
	args := []string{
		cliCmd,
		"-i", inputCSS,
		"-o", outputCSS,
	}

	// Add minification in production
	if cfg.Minify {
		args = append(args, "--minify")
	}

	// Execute Tailwind v4 CLI
	cmd := exec.CommandContext(ctx, "npx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set working directory to the input CSS's parent directory so npm/node
	// can resolve the `tailwindcss` package (required by `@import "tailwindcss"`).
	cmd.Dir = filepath.Dir(inputCSS)

	if tp.Verbose {
		fmt.Printf("[Tailwind v4] Running: npx %s\n", strings.Join(args, " "))
		fmt.Printf("[Tailwind v4] Working directory: %s\n", cmd.Dir)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tailwind v4 build failed: %w", err)
	}

	if tp.Verbose {
		fmt.Printf("[Tailwind v4] Generated CSS: %s\n", outputCSS)
	}

	return nil
}

// processV3 handles Tailwind CSS v3 processing (backward compatibility).
// v3 uses JS-based config and @tailwind directives.
func (tp *TailwindProcessor) processV3(ctx context.Context, cfg ProcessorConfig) error {
	// Generate tailwind.config.js if not provided
	configPath := tp.ConfigPath
	if configPath == "" {
		var err error

		configPath, err = tp.generateConfig(cfg.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to generate tailwind config: %w", err)
		}

		defer func() { _ = os.Remove(configPath) }()
	}

	// Create input CSS if not provided
	inputCSS := tp.InputCSS
	if inputCSS == "" {
		var err error

		inputCSS, err = tp.createInputCSSv3(cfg.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to create input CSS: %w", err)
		}

		defer func() { _ = os.Remove(inputCSS) }()
	}

	// Determine output path
	outputCSS := filepath.Join(cfg.OutputDir, tp.OutputCSS)

	// Ensure output directory exists
	outputDir := filepath.Dir(outputCSS)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if Tailwind CLI is available
	if !tp.isTailwindV3Available() {
		if tp.UseCDN {
			return tp.generateCDNFallback(outputCSS)
		}

		return errors.New("tailwind CLI not found and CDN fallback disabled")
	}

	// Build Tailwind v3 CSS command
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

	// Execute Tailwind v3 CLI
	cmd := exec.CommandContext(ctx, "npx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if tp.Verbose {
		fmt.Printf("[Tailwind v3] Running: npx %s\n", strings.Join(args, " "))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tailwind v3 build failed: %w", err)
	}

	if tp.Verbose {
		fmt.Printf("[Tailwind v3] Generated CSS: %s\n", outputCSS)
	}

	return nil
}

// createInputCSSv4 generates a Tailwind v4 input CSS file.
// If themes are configured, uses theme.GenerateInputCSS() for a complete
// CSS-based configuration. Otherwise generates a minimal v4 input.
func (tp *TailwindProcessor) createInputCSSv4(outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	inputPath := filepath.Join(outputDir, "_tailwind_input.css")

	var content string
	if tp.LightTheme != nil && tp.DarkTheme != nil {
		// Generate complete input CSS from theme tokens
		content = theme.GenerateInputCSS(*tp.LightTheme, *tp.DarkTheme)
	} else {
		// Minimal v4 input — just the import
		content = "@import \"tailwindcss\";\n"
	}

	if err := os.WriteFile(inputPath, []byte(content), 0600); err != nil {
		return "", err
	}

	if tp.Verbose {
		if tp.LightTheme != nil && tp.DarkTheme != nil {
			fmt.Printf("[Tailwind v4] Generated theme-aware input CSS: %s\n", inputPath)
		} else {
			fmt.Printf("[Tailwind v4] Generated minimal input CSS: %s\n", inputPath)
		}
	}

	return inputPath, nil
}

// createInputCSSv3 creates a v3 input CSS file with @tailwind directives.
func (tp *TailwindProcessor) createInputCSSv3(outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	inputPath := filepath.Join(outputDir, "_tailwind_input.css")

	content := `@tailwind base;
@tailwind components;
@tailwind utilities;
`

	if err := os.WriteFile(inputPath, []byte(content), 0600); err != nil {
		return "", err
	}

	return inputPath, nil
}

// generateConfig creates a tailwind.config.js file (v3 only).
func (tp *TailwindProcessor) generateConfig(outputDir string) (string, error) {
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

	var contentJSONSb strings.Builder

	for i, path := range contentPaths {
		if i > 0 {
			contentJSONSb.WriteString(", ")
		}

		fmt.Fprintf(&contentJSONSb, `"%s"`, path)
	}

	contentJSON += contentJSONSb.String()
	contentJSON += "]"

	// Insert content paths into config
	config = strings.Replace(config, "module.exports = {",
		fmt.Sprintf("module.exports = {\n  content: %s,", contentJSON), 1)

	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		return "", err
	}

	if tp.Verbose {
		fmt.Printf("[Tailwind v3] Generated config: %s\n", configPath)
	}

	return configPath, nil
}

// findTailwindV4CLI detects the available Tailwind v4 CLI command.
// Returns the npx package name or empty string if not found.
// Checks in order: @tailwindcss/cli (v4), then tailwindcss (v4 compat).
func (tp *TailwindProcessor) findTailwindV4CLI() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try @tailwindcss/cli first (the official v4 package)
	cmd := exec.CommandContext(ctx, "npx", "@tailwindcss/cli", "--help")
	if err := cmd.Run(); err == nil {
		return "@tailwindcss/cli"
	}

	// Try tailwindcss (may be v4 installed under the old name)
	cmd = exec.CommandContext(ctx, "npx", "tailwindcss", "--help")
	if err := cmd.Run(); err == nil {
		return "tailwindcss"
	}

	return ""
}

// isTailwindV3Available checks if Tailwind v3 CLI is available.
func (tp *TailwindProcessor) isTailwindV3Available() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "npx", "tailwindcss", "--help")
	err := cmd.Run()

	return err == nil
}

// IsTailwindAvailable checks if any Tailwind CLI is available for the configured version.
func (tp *TailwindProcessor) IsTailwindAvailable() bool {
	if tp.Version == TailwindV4 {
		return tp.findTailwindV4CLI() != ""
	}

	return tp.isTailwindV3Available()
}

// generateCDNFallback creates a minimal CSS file that references Tailwind CDN
func (tp *TailwindProcessor) generateCDNFallback(outputPath string) error {
	var content string

	if tp.LightTheme != nil && tp.DarkTheme != nil {
		// Generate CSS variables even in CDN mode so themes work
		content = fmt.Sprintf(`/* Tailwind CSS via CDN - ForgeUI fallback */
/* Add to your HTML: <script src="https://cdn.tailwindcss.com"></script> */
/* For production, install: npm install -D @tailwindcss/cli */

%s
`, theme.GenerateCSS(*tp.LightTheme, *tp.DarkTheme))
	} else {
		content = `/* Tailwind CSS via CDN */
/* Add this to your HTML: <script src="https://cdn.tailwindcss.com"></script> */

/* ForgeUI: Tailwind CLI not found. Using CDN fallback. */
/* For production, install Tailwind: npm install -D @tailwindcss/cli */

:root {
  /* ForgeUI theme variables will be injected by theme system */
}
`
	}

	if err := os.WriteFile(outputPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write CDN fallback: %w", err)
	}

	fmt.Println("[Tailwind] WARNING: Tailwind CLI not found. Generated CDN fallback.")
	fmt.Println("[Tailwind] For production, install: npm install -D @tailwindcss/cli")

	return nil
}

// --- Option methods ---

// WithVersion sets the Tailwind CSS major version (TailwindV3 or TailwindV4).
func (tp *TailwindProcessor) WithVersion(version TailwindVersion) *TailwindProcessor {
	tp.Version = version
	return tp
}

// WithThemes sets the light and dark themes for generating v4 input CSS.
// When themes are set and no custom InputCSS is provided, the processor
// auto-generates a complete Tailwind v4 input.css from the theme tokens.
func (tp *TailwindProcessor) WithThemes(light, dark *theme.Theme) *TailwindProcessor {
	tp.LightTheme = light
	tp.DarkTheme = dark
	return tp
}

// WithConfigPath sets a custom tailwind.config.js path (v3 only)
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

// WithContentPaths sets custom content paths to scan for classes (v3 only)
func (tp *TailwindProcessor) WithContentPaths(paths []string) *TailwindProcessor {
	tp.ContentPaths = paths
	return tp
}

// WithVerbose enables verbose logging
func (tp *TailwindProcessor) WithVerbose(verbose bool) *TailwindProcessor {
	tp.Verbose = verbose
	return tp
}

// WithCDNFallback enables or disables CDN fallback when CLI is not available.
func (tp *TailwindProcessor) WithCDNFallback(useCDN bool) *TailwindProcessor {
	tp.UseCDN = useCDN
	return tp
}
