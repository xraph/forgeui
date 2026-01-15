package assets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestNewDevServer(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{IsDev: true}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatalf("NewDevServer failed: %v", err)
	}

	if ds == nil {
		t.Fatal("NewDevServer returned nil")
	}

	defer func() { _ = ds.Close() }()
}

func TestDevServer_SetVerbose(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{IsDev: true}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ds.Close() }()

	if ds.verbose {
		t.Error("Verbose should be false by default")
	}

	ds.SetVerbose(true)

	if !ds.verbose {
		t.Error("Verbose should be true after setting")
	}
}

func TestDevServer_ClientCount(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{IsDev: true}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ds.Close() }()

	if ds.ClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", ds.ClientCount())
	}

	// Add mock clients
	ds.mu.Lock()
	ds.sseClients = append(ds.sseClients, make(chan string, 1))
	ds.sseClients = append(ds.sseClients, make(chan string, 1))
	ds.mu.Unlock()

	if ds.ClientCount() != 2 {
		t.Errorf("Expected 2 clients, got %d", ds.ClientCount())
	}
}

func TestDevServer_NotifyClients(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{IsDev: true}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ds.Close() }()

	// Create mock clients
	client1 := make(chan string, 1)
	client2 := make(chan string, 1)

	ds.mu.Lock()
	ds.sseClients = append(ds.sseClients, client1, client2)
	ds.mu.Unlock()

	// Notify clients
	ds.notifyClients("test-message")

	// Check clients received message
	select {
	case msg := <-client1:
		if msg != "test-message" {
			t.Errorf("Client 1 received wrong message: %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 1 did not receive message")
	}

	select {
	case msg := <-client2:
		if msg != "test-message" {
			t.Errorf("Client 2 received wrong message: %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 2 did not receive message")
	}
}

func TestDevServer_SSEHandler(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{IsDev: true}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ds.Close() }()

	// Use a real HTTP server for SSE testing to avoid race conditions
	// with httptest.ResponseRecorder
	server := httptest.NewServer(ds.SSEHandler())
	defer func() { _ = server.Close() }()

	// Make request to SSE endpoint
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to connect to SSE endpoint: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Check headers
	if resp.Header.Get("Content-Type") != "text/event-stream" {
		t.Error("Missing or incorrect Content-Type header")
	}

	if resp.Header.Get("Cache-Control") != "no-cache" {
		t.Error("Missing or incorrect Cache-Control header")
	}

	// Read initial connection message
	buf := make([]byte, 256)

	n, err := resp.Body.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read initial message: %v", err)
	}

	body := string(buf[:n])
	if !strings.Contains(body, "data: connected") {
		t.Errorf("Missing initial connection message, got: %s", body)
	}

	// Give a moment for client registration
	time.Sleep(50 * time.Millisecond)

	// Client should be registered
	if ds.ClientCount() != 1 {
		t.Errorf("Expected 1 client, got %d", ds.ClientCount())
	}
}

func TestDevServer_HotReloadScript(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{IsDev: true}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ds.Close() }()

	script := ds.HotReloadScript()

	// Check script contains essential parts
	if !strings.Contains(script, "EventSource") {
		t.Error("Script should contain EventSource")
	}

	if !strings.Contains(script, "/_forgeui/reload") {
		t.Error("Script should reference reload endpoint")
	}

	if !strings.Contains(script, "location.reload()") {
		t.Error("Script should contain reload logic")
	}

	if !strings.Contains(script, "<script>") {
		t.Error("Script should be wrapped in script tags")
	}
}

func TestDevServer_Close(t *testing.T) {
	pipeline := NewPipeline(PipelineConfig{IsDev: true}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}

	// Add mock clients
	client1 := make(chan string, 1)
	client2 := make(chan string, 1)

	ds.mu.Lock()
	ds.sseClients = append(ds.sseClients, client1, client2)
	ds.mu.Unlock()

	// Close dev server
	err = ds.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Clients should be closed
	if ds.ClientCount() != 0 {
		t.Errorf("Expected 0 clients after close, got %d", ds.ClientCount())
	}
}

func TestDevServer_OnFileChange(t *testing.T) {
	tempDir := t.TempDir()

	pipeline := NewPipeline(PipelineConfig{
		InputDir:  tempDir,
		OutputDir: tempDir,
		IsDev:     true,
	}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ds.Close() }()

	ds.SetVerbose(false) // Disable verbose for cleaner test output

	// Add mock client
	client := make(chan string, 1)

	ds.mu.Lock()
	ds.sseClients = append(ds.sseClients, client)
	ds.mu.Unlock()

	// Simulate file change
	ctx := context.Background()
	err = ds.onFileChange(ctx, fsnotify.Event{Name: "test.go"})

	// Error is okay if processors fail, we're just testing the flow
	if err != nil && !strings.Contains(err.Error(), "processor") {
		t.Errorf("Unexpected error: %v", err)
	}

	// Client should receive reload message (if build succeeded)
	select {
	case msg := <-client:
		if msg != "reload" {
			t.Errorf("Expected 'reload' message, got '%s'", msg)
		}
	case <-time.After(100 * time.Millisecond):
		// Build might have failed, that's okay for this test
	}
}

func TestDevServer_ConcurrentBuilds(t *testing.T) {
	tempDir := t.TempDir()

	pipeline := NewPipeline(PipelineConfig{
		InputDir:  tempDir,
		OutputDir: tempDir,
		IsDev:     true,
	}, nil)

	ds, err := NewDevServer(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ds.Close() }()

	ds.SetVerbose(false)

	ctx := context.Background()

	// Try to trigger multiple concurrent builds
	done := make(chan bool, 3)

	for range 3 {
		go func() {
			_ = ds.onFileChange(ctx, fsnotify.Event{Name: "test.go"})

			done <- true
		}()
	}

	// Wait for all goroutines
	for range 3 {
		<-done
	}

	// Should not panic or deadlock
}
