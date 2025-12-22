package assets

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ESBuildProcessor bundles and minifies JavaScript files using esbuild.
// It's optional for SSR applications but useful for custom JavaScript.
type ESBuildProcessor struct {
	// EntryPoints are the JavaScript entry files to bundle
	EntryPoints []string

	// Outfile is the output bundle file path (relative to output directory)
	Outfile string

	// Format specifies the output format ("iife", "cjs", "esm")
	Format string

	// Target specifies the JavaScript target (e.g., "es2020")
	Target string

	// Bundle enables bundling of dependencies
	Bundle bool

	// Splitting enables code splitting (requires format: "esm")
	Splitting bool

	// External are packages to exclude from bundling
	External []string

	// Verbose enables detailed logging
	Verbose bool
}

// NewESBuildProcessor creates a new ESBuild processor with sensible defaults
func NewESBuildProcessor() *ESBuildProcessor {
	return &ESBuildProcessor{
		Format:  "iife",
		Target:  "es2020",
		Bundle:  true,
		Outfile: "js/app.js",
	}
}

// Name returns the processor name
func (ep *ESBuildProcessor) Name() string {
	return "ESBuild"
}

// FileTypes returns the file extensions this processor handles
func (ep *ESBuildProcessor) FileTypes() []string {
	return []string{".js", ".ts", ".jsx", ".tsx"}
}

// Process executes the ESBuild bundling
func (ep *ESBuildProcessor) Process(ctx context.Context, cfg ProcessorConfig) error {
	// Check if esbuild is available
	if !ep.isESBuildAvailable() {
		if ep.Verbose {
			fmt.Println("[ESBuild] ESBuild not found, skipping JavaScript bundling")
			fmt.Println("[ESBuild] Install with: npm install -D esbuild")
		}

		return nil // Don't fail, just skip
	}

	// Validate entry points
	if len(ep.EntryPoints) == 0 {
		if ep.Verbose {
			fmt.Println("[ESBuild] No entry points specified, skipping")
		}

		return nil
	}

	// Determine output path
	outfile := filepath.Join(cfg.OutputDir, ep.Outfile)

	// Ensure output directory exists
	outDir := filepath.Dir(outfile)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build esbuild command
	args := []string{"esbuild"}

	// Add entry points
	args = append(args, ep.EntryPoints...)

	// Add output file
	args = append(args, "--outfile="+outfile)

	// Add format
	if ep.Format != "" {
		args = append(args, "--format="+ep.Format)
	}

	// Add target
	if ep.Target != "" {
		args = append(args, "--target="+ep.Target)
	}

	// Add bundle flag
	if ep.Bundle {
		args = append(args, "--bundle")
	}

	// Add minification in production
	if cfg.Minify {
		args = append(args, "--minify")
	}

	// Add source maps
	if cfg.SourceMaps {
		args = append(args, "--sourcemap")
	}

	// Add code splitting
	if ep.Splitting {
		args = append(args, "--splitting")
		// Splitting requires outdir instead of outfile
		args = append(args, "--outdir="+filepath.Dir(outfile))
	}

	// Add external packages
	for _, ext := range ep.External {
		args = append(args, "--external:"+ext)
	}

	// Execute esbuild
	cmd := exec.CommandContext(ctx, "npx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if ep.Verbose {
		fmt.Printf("[ESBuild] Running: npx %s\n", strings.Join(args, " "))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("esbuild failed: %w", err)
	}

	if ep.Verbose {
		fmt.Printf("[ESBuild] Generated bundle: %s\n", outfile)
	}

	return nil
}

// isESBuildAvailable checks if esbuild CLI is available
func (ep *ESBuildProcessor) isESBuildAvailable() bool {
	cmd := exec.Command("npx", "esbuild", "--version")
	err := cmd.Run()

	return err == nil
}

// WithEntryPoints sets the JavaScript entry points
func (ep *ESBuildProcessor) WithEntryPoints(entries ...string) *ESBuildProcessor {
	ep.EntryPoints = entries
	return ep
}

// WithOutfile sets the output bundle file path
func (ep *ESBuildProcessor) WithOutfile(path string) *ESBuildProcessor {
	ep.Outfile = path
	return ep
}

// WithFormat sets the output format
func (ep *ESBuildProcessor) WithFormat(format string) *ESBuildProcessor {
	ep.Format = format
	return ep
}

// WithTarget sets the JavaScript target
func (ep *ESBuildProcessor) WithTarget(target string) *ESBuildProcessor {
	ep.Target = target
	return ep
}

// WithBundle enables or disables bundling
func (ep *ESBuildProcessor) WithBundle(bundle bool) *ESBuildProcessor {
	ep.Bundle = bundle
	return ep
}

// WithSplitting enables code splitting
func (ep *ESBuildProcessor) WithSplitting(splitting bool) *ESBuildProcessor {
	ep.Splitting = splitting
	return ep
}

// WithExternal sets packages to exclude from bundling
func (ep *ESBuildProcessor) WithExternal(external ...string) *ESBuildProcessor {
	ep.External = external
	return ep
}

// WithVerbose enables verbose logging
func (ep *ESBuildProcessor) WithVerbose(verbose bool) *ESBuildProcessor {
	ep.Verbose = verbose
	return ep
}
