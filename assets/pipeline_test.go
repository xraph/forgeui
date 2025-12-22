package assets

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// Mock processor for testing
type mockProcessor struct {
	name      string
	fileTypes []string
	executed  bool
	shouldErr bool
}

func (m *mockProcessor) Name() string {
	return m.name
}

func (m *mockProcessor) Process(ctx context.Context, cfg ProcessorConfig) error {
	m.executed = true
	if m.shouldErr {
		return errors.New("mock processor error")
	}
	return nil
}

func (m *mockProcessor) FileTypes() []string {
	return m.fileTypes
}

func TestNewPipeline(t *testing.T) {
	cfg := PipelineConfig{
		InputDir:  "test-input",
		OutputDir: "test-output",
		IsDev:     true,
	}

	pipeline := NewPipeline(cfg, nil)

	if pipeline == nil {
		t.Fatal("NewPipeline returned nil")
	}

	if pipeline.ProcessorCount() != 0 {
		t.Errorf("New pipeline should have 0 processors, got %d", pipeline.ProcessorCount())
	}
}

func TestNewPipeline_Defaults(t *testing.T) {
	cfg := PipelineConfig{}
	pipeline := NewPipeline(cfg, nil)

	config := pipeline.Config()

	if config.InputDir != "assets" {
		t.Errorf("Expected default InputDir 'assets', got '%s'", config.InputDir)
	}

	if config.OutputDir != "dist" {
		t.Errorf("Expected default OutputDir 'dist', got '%s'", config.OutputDir)
	}
}

func TestNewPipeline_ProductionDefaults(t *testing.T) {
	cfg := PipelineConfig{
		IsDev: false,
	}
	pipeline := NewPipeline(cfg, nil)

	config := pipeline.Config()

	if !config.Minify {
		t.Error("Production builds should enable minification by default")
	}
}

func TestPipeline_AddProcessor(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{}, nil)

	proc1 := &mockProcessor{name: "proc1"}
	proc2 := &mockProcessor{name: "proc2"}

	pipeline.AddProcessor(proc1).AddProcessor(proc2)

	if pipeline.ProcessorCount() != 2 {
		t.Errorf("Expected 2 processors, got %d", pipeline.ProcessorCount())
	}
}

func TestPipeline_Build_Success(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	cfg := PipelineConfig{
		InputDir:  inputDir,
		OutputDir: outputDir,
		IsDev:     true,
	}

	pipeline := NewPipeline(cfg, nil)

	proc1 := &mockProcessor{name: "proc1"}
	proc2 := &mockProcessor{name: "proc2"}

	pipeline.AddProcessor(proc1).AddProcessor(proc2)

	ctx := context.Background()
	if err := pipeline.Build(ctx); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if !proc1.executed {
		t.Error("Processor 1 was not executed")
	}

	if !proc2.executed {
		t.Error("Processor 2 was not executed")
	}

	// Check output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("Output directory was not created")
	}
}

func TestPipeline_Build_ProcessorError(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	cfg := PipelineConfig{
		InputDir:  inputDir,
		OutputDir: outputDir,
		IsDev:     true,
	}

	pipeline := NewPipeline(cfg, nil)

	proc1 := &mockProcessor{name: "proc1"}
	proc2 := &mockProcessor{name: "proc2", shouldErr: true}
	proc3 := &mockProcessor{name: "proc3"}

	pipeline.AddProcessor(proc1).AddProcessor(proc2).AddProcessor(proc3)

	ctx := context.Background()
	err := pipeline.Build(ctx)

	if err == nil {
		t.Fatal("Build should have failed")
	}

	if !proc1.executed {
		t.Error("Processor 1 should have been executed")
	}

	if !proc2.executed {
		t.Error("Processor 2 should have been executed")
	}

	if proc3.executed {
		t.Error("Processor 3 should not have been executed after error")
	}
}

func TestPipeline_Build_CleanOutput(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create output directory with a file
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}
	testFile := filepath.Join(outputDir, "old-file.txt")
	if err := os.WriteFile(testFile, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := PipelineConfig{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		IsDev:       true,
		CleanOutput: true,
	}

	pipeline := NewPipeline(cfg, nil)

	ctx := context.Background()
	if err := pipeline.Build(ctx); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Old file should be removed
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Old file should have been removed")
	}

	// Output directory should still exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("Output directory should exist")
	}
}

func TestPipeline_Build_CreatesOutputDir(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output", "nested", "path")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	cfg := PipelineConfig{
		InputDir:  inputDir,
		OutputDir: outputDir,
		IsDev:     true,
	}

	pipeline := NewPipeline(cfg, nil)

	ctx := context.Background()
	if err := pipeline.Build(ctx); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Nested output directory should be created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("Nested output directory was not created")
	}
}

func TestPipeline_ProcessorCount(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{}, nil)

	if pipeline.ProcessorCount() != 0 {
		t.Errorf("Expected 0 processors, got %d", pipeline.ProcessorCount())
	}

	pipeline.AddProcessor(&mockProcessor{name: "proc1"})
	if pipeline.ProcessorCount() != 1 {
		t.Errorf("Expected 1 processor, got %d", pipeline.ProcessorCount())
	}

	pipeline.AddProcessor(&mockProcessor{name: "proc2"})
	if pipeline.ProcessorCount() != 2 {
		t.Errorf("Expected 2 processors, got %d", pipeline.ProcessorCount())
	}
}

func TestPipeline_Config(t *testing.T) {
	cfg := PipelineConfig{
		InputDir:  "custom-input",
		OutputDir: "custom-output",
		IsDev:     true,
		Minify:    false,
	}

	pipeline := NewPipeline(cfg, nil)
	returnedCfg := pipeline.Config()

	if returnedCfg.InputDir != cfg.InputDir {
		t.Errorf("Expected InputDir '%s', got '%s'", cfg.InputDir, returnedCfg.InputDir)
	}

	if returnedCfg.OutputDir != cfg.OutputDir {
		t.Errorf("Expected OutputDir '%s', got '%s'", cfg.OutputDir, returnedCfg.OutputDir)
	}

	if returnedCfg.IsDev != cfg.IsDev {
		t.Errorf("Expected IsDev %v, got %v", cfg.IsDev, returnedCfg.IsDev)
	}
}

func TestPipeline_Processors(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{}, nil)

	proc1 := &mockProcessor{name: "proc1"}
	proc2 := &mockProcessor{name: "proc2"}

	pipeline.AddProcessor(proc1).AddProcessor(proc2)

	processors := pipeline.Processors()

	if len(processors) != 2 {
		t.Errorf("Expected 2 processors, got %d", len(processors))
	}

	// Verify it's a copy (modifying returned slice shouldn't affect pipeline)
	processors[0] = &mockProcessor{name: "modified"}

	if pipeline.Processors()[0].Name() != "proc1" {
		t.Error("Modifying returned slice should not affect pipeline")
	}
}

func TestPipeline_ProcessorConfig(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	cfg := PipelineConfig{
		InputDir:   inputDir,
		OutputDir:  outputDir,
		IsDev:      false,
		Minify:     true,
		SourceMaps: true,
		Watch:      false,
	}

	pipeline := NewPipeline(cfg, nil)

	// Create a processor that verifies it receives correct config
	var receivedConfig ProcessorConfig
	verifyProc := &configCapturingProcessor{
		mockProcessor: mockProcessor{name: "verify"},
		capturedCfg:   &receivedConfig,
	}

	pipeline.AddProcessor(verifyProc)

	ctx := context.Background()
	if err := pipeline.Build(ctx); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify processor received correct config
	if receivedConfig.InputDir != inputDir {
		t.Errorf("Processor received wrong InputDir: %s", receivedConfig.InputDir)
	}

	if receivedConfig.OutputDir != outputDir {
		t.Errorf("Processor received wrong OutputDir: %s", receivedConfig.OutputDir)
	}

	if receivedConfig.IsDev != false {
		t.Error("Processor received wrong IsDev value")
	}

	if receivedConfig.Minify != true {
		t.Error("Processor received wrong Minify value")
	}

	if receivedConfig.SourceMaps != true {
		t.Error("Processor received wrong SourceMaps value")
	}
}

// configCapturingProcessor captures the config passed to Process
type configCapturingProcessor struct {
	mockProcessor
	capturedCfg *ProcessorConfig
}

func (c *configCapturingProcessor) Process(ctx context.Context, cfg ProcessorConfig) error {
	*c.capturedCfg = cfg
	return c.mockProcessor.Process(ctx, cfg)
}

func TestPipeline_ContextCancellation(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	cfg := PipelineConfig{
		InputDir:  inputDir,
		OutputDir: outputDir,
		IsDev:     true,
	}

	pipeline := NewPipeline(cfg, nil)

	// Create a processor that checks context
	ctxProc := &contextCheckingProcessor{
		mockProcessor: mockProcessor{name: "ctx-check"},
	}

	pipeline.AddProcessor(ctxProc)

	// Cancel context before build
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := pipeline.Build(ctx)
	if err == nil {
		t.Error("Build should fail with cancelled context")
	}
}

// contextCheckingProcessor checks if context is cancelled
type contextCheckingProcessor struct {
	mockProcessor
}

func (c *contextCheckingProcessor) Process(ctx context.Context, cfg ProcessorConfig) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

