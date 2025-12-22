package assets

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Processor processes assets (CSS, JS, images, etc.)
// Implementations can handle different asset types and transformations.
type Processor interface {
	// Name returns the processor name for logging and debugging
	Name() string

	// Process executes the asset processing with the given configuration
	Process(ctx context.Context, cfg ProcessorConfig) error

	// FileTypes returns the file extensions this processor handles (e.g., [".css", ".js"])
	FileTypes() []string
}

// ProcessorConfig contains configuration for asset processing
type ProcessorConfig struct {
	// InputDir is the source directory for assets
	InputDir string

	// OutputDir is the destination directory for processed assets
	OutputDir string

	// IsDev indicates if running in development mode
	IsDev bool

	// Minify enables minification of assets
	Minify bool

	// SourceMaps enables source map generation
	SourceMaps bool

	// Watch enables file watching for automatic rebuilds
	Watch bool

	// CustomConfig allows processors to receive custom configuration
	CustomConfig map[string]any
}

// Pipeline orchestrates multiple asset processors in sequence.
// It ensures processors run in the correct order and handles errors gracefully.
type Pipeline struct {
	processors []Processor
	config     PipelineConfig
	mu         sync.RWMutex
	manager    *Manager
}

// PipelineConfig defines the pipeline's overall configuration
type PipelineConfig struct {
	// InputDir is the source directory for all assets
	InputDir string

	// OutputDir is the destination directory for processed assets
	OutputDir string

	// IsDev enables development mode (no minification, faster builds)
	IsDev bool

	// Minify enables minification of assets in production
	Minify bool

	// SourceMaps enables source map generation
	SourceMaps bool

	// Watch enables file watching for automatic rebuilds
	Watch bool

	// CleanOutput removes the output directory before building
	CleanOutput bool

	// Verbose enables detailed logging
	Verbose bool
}

// NewPipeline creates a new asset pipeline with the given configuration
func NewPipeline(cfg PipelineConfig, manager *Manager) *Pipeline {
	// Set defaults
	if cfg.InputDir == "" {
		cfg.InputDir = "assets"
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = "dist"
	}

	// In production, enable minification by default
	if !cfg.IsDev && !cfg.Minify {
		cfg.Minify = true
	}

	return &Pipeline{
		processors: make([]Processor, 0),
		config:     cfg,
		manager:    manager,
	}
}

// AddProcessor adds a processor to the pipeline.
// Processors are executed in the order they are added.
func (p *Pipeline) AddProcessor(processor Processor) *Pipeline {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.processors = append(p.processors, processor)

	return p
}

// Build executes all processors in sequence.
// If any processor fails, the build stops and returns the error.
func (p *Pipeline) Build(ctx context.Context) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Clean output directory if requested
	if p.config.CleanOutput {
		if err := p.cleanOutput(); err != nil {
			return fmt.Errorf("failed to clean output: %w", err)
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(p.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create processor config
	procConfig := ProcessorConfig{
		InputDir:   p.config.InputDir,
		OutputDir:  p.config.OutputDir,
		IsDev:      p.config.IsDev,
		Minify:     p.config.Minify,
		SourceMaps: p.config.SourceMaps,
		Watch:      p.config.Watch,
	}

	// Execute each processor
	for _, processor := range p.processors {
		if p.config.Verbose {
			fmt.Printf("[Pipeline] Running processor: %s\n", processor.Name())
		}

		if err := processor.Process(ctx, procConfig); err != nil {
			return fmt.Errorf("processor %s failed: %w", processor.Name(), err)
		}

		if p.config.Verbose {
			fmt.Printf("[Pipeline] Processor %s completed successfully\n", processor.Name())
		}
	}

	// Generate manifest for production builds
	if !p.config.IsDev && p.manager != nil {
		if err := p.generateManifest(); err != nil {
			return fmt.Errorf("failed to generate manifest: %w", err)
		}
	}

	return nil
}

// cleanOutput removes the output directory
func (p *Pipeline) cleanOutput() error {
	if _, err := os.Stat(p.config.OutputDir); err == nil {
		if err := os.RemoveAll(p.config.OutputDir); err != nil {
			return err
		}
	}

	return nil
}

// generateManifest creates an asset manifest file for production
func (p *Pipeline) generateManifest() error {
	if p.manager == nil {
		return nil
	}

	// Fingerprint all assets in output directory
	if err := p.manager.FingerprintAll(); err != nil {
		return fmt.Errorf("failed to fingerprint assets: %w", err)
	}

	// Save manifest
	manifestPath := filepath.Join(p.config.OutputDir, "manifest.json")
	if err := p.manager.SaveManifest(manifestPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	if p.config.Verbose {
		fmt.Printf("[Pipeline] Generated manifest: %s\n", manifestPath)
	}

	return nil
}

// ProcessorCount returns the number of processors in the pipeline
func (p *Pipeline) ProcessorCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.processors)
}

// Config returns the pipeline configuration
func (p *Pipeline) Config() PipelineConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.config
}

// Processors returns a copy of the processor list
func (p *Pipeline) Processors() []Processor {
	p.mu.RLock()
	defer p.mu.RUnlock()

	processors := make([]Processor, len(p.processors))
	copy(processors, p.processors)

	return processors
}
