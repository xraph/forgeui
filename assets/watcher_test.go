package assets

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestNewWatcher(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	if w == nil {
		t.Fatal("NewWatcher returned nil")
	}
}

func TestWatcher_AddPath(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	tempDir := t.TempDir()

	err = w.AddPath(tempDir)
	if err != nil {
		t.Errorf("AddPath failed: %v", err)
	}
}

func TestWatcher_AddPattern(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	w.AddPattern("*.go")
	w.AddPattern("*.css")

	if len(w.patterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(w.patterns))
	}
}

func TestWatcher_OnChange(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	w.OnChange(func(event fsnotify.Event) error {
		return nil
	})

	if len(w.callbacks) != 1 {
		t.Errorf("Expected 1 callback, got %d", len(w.callbacks))
	}
}

func TestWatcher_SetDebounce(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	customDebounce := 500 * time.Millisecond
	w.SetDebounce(customDebounce)

	if w.debounce != customDebounce {
		t.Errorf("Expected debounce %v, got %v", customDebounce, w.debounce)
	}
}

func TestWatcher_SetVerbose(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	if w.verbose {
		t.Error("Verbose should be false by default")
	}

	w.SetVerbose(true)

	if !w.verbose {
		t.Error("Verbose should be true after setting")
	}
}

func TestWatcher_FileChange(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	// Disable debouncing for more predictable test timing
	w.SetDebounce(0)

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// Create initial file
	if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	// Add watcher
	if err := w.AddPath(tempDir); err != nil {
		t.Fatal(err)
	}

	// Set up callback with channel to signal when watcher is ready
	changed := make(chan bool, 1)
	started := make(chan bool, 1)

	w.OnChange(func(event fsnotify.Event) error {
		if filepath.Base(event.Name) == "test.txt" {
			select {
			case changed <- true:
			default:
			}
		}

		return nil
	})

	// Start watcher
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	go func() {
		started <- true

		_ = w.Start(ctx)
	}()

	// Wait for watcher goroutine to start
	<-started

	// Give fsnotify additional time to initialize watches (especially on Linux)
	// Linux/Ubuntu often needs more time for inotify setup
	time.Sleep(500 * time.Millisecond)

	// Modify file
	if err := os.WriteFile(testFile, []byte("modified"), 0644); err != nil {
		t.Fatal(err)
	}

	// Wait for change notification with increased timeout for slower CI environments
	// fsnotify behavior varies significantly across OSes:
	// - macOS: typically fast (<100ms)
	// - Linux: can be slower, especially with inotify (100-500ms)
	// - Windows: ReadDirectoryChangesW can also have latency
	select {
	case <-changed:
		// Success
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for file change notification")
	}
}

func TestWatcher_PatternFiltering(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	// Only watch .go files
	w.AddPattern("*.go")

	// Test shouldProcess
	goEvent := fsnotify.Event{
		Name: "/path/to/file.go",
		Op:   fsnotify.Write,
	}

	if !w.shouldProcess(goEvent) {
		t.Error("Should process .go file")
	}

	txtEvent := fsnotify.Event{
		Name: "/path/to/file.txt",
		Op:   fsnotify.Write,
	}

	if w.shouldProcess(txtEvent) {
		t.Error("Should not process .txt file")
	}
}

func TestWatcher_IgnoreNonWriteEvents(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	// Test different event types
	events := []struct {
		op         fsnotify.Op
		shouldProc bool
	}{
		{fsnotify.Write, true},
		{fsnotify.Create, true},
		{fsnotify.Remove, false},
		{fsnotify.Rename, false},
		{fsnotify.Chmod, false},
	}

	for _, tc := range events {
		event := fsnotify.Event{
			Name: "/path/to/file.txt",
			Op:   tc.op,
		}

		result := w.shouldProcess(event)
		if result != tc.shouldProc {
			t.Errorf("Event %v: expected shouldProcess=%v, got %v", tc.op, tc.shouldProc, result)
		}
	}
}

func TestWatcher_Close(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}

	err = w.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestWatcher_WatchDirectory(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	tempDir := t.TempDir()

	// Create nested directory structure
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create hidden directory (should be skipped)
	hiddenDir := filepath.Join(tempDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}

	err = w.WatchDirectory(tempDir)
	if err != nil {
		t.Errorf("WatchDirectory failed: %v", err)
	}
}

func TestWatcher_MultipleCallbacks(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	callback1Called := false
	callback2Called := false

	w.OnChange(func(event fsnotify.Event) error {
		callback1Called = true
		return nil
	})

	w.OnChange(func(event fsnotify.Event) error {
		callback2Called = true
		return nil
	})

	// Simulate event processing
	event := fsnotify.Event{
		Name: "/path/to/file.txt",
		Op:   fsnotify.Write,
	}

	w.processEvent(event)

	if !callback1Called {
		t.Error("Callback 1 was not called")
	}

	if !callback2Called {
		t.Error("Callback 2 was not called")
	}
}

func TestWatcher_ContextCancellation(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	tempDir := t.TempDir()
	if err := w.AddPath(tempDir); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)

	go func() {
		done <- w.Start(ctx)
	}()

	// Cancel context
	cancel()

	// Wait for Start to return
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for context cancellation")
	}
}
