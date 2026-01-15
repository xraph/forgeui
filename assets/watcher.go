package assets

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches files for changes and triggers callbacks.
// It uses fsnotify for cross-platform file system notifications.
type Watcher struct {
	watcher   *fsnotify.Watcher
	callbacks []WatchCallback
	debounce  time.Duration
	mu        sync.RWMutex
	patterns  []string
	verbose   bool
}

// WatchCallback is called when files change
type WatchCallback func(event fsnotify.Event) error

// NewWatcher creates a new file watcher
func NewWatcher() (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		watcher:   fsw,
		callbacks: make([]WatchCallback, 0),
		debounce:  300 * time.Millisecond, // Default 300ms debounce
		patterns:  make([]string, 0),
	}, nil
}

// AddPath adds a path to watch (file or directory)
func (w *Watcher) AddPath(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.watcher.Add(path); err != nil {
		return fmt.Errorf("failed to watch %s: %w", path, err)
	}

	if w.verbose {
		fmt.Printf("[Watcher] Watching: %s\n", path)
	}

	return nil
}

// AddPattern adds a glob pattern for file matching
func (w *Watcher) AddPattern(pattern string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.patterns = append(w.patterns, pattern)
}

// OnChange registers a callback for file changes
func (w *Watcher) OnChange(callback WatchCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.callbacks = append(w.callbacks, callback)
}

// SetDebounce sets the debounce duration for file changes
func (w *Watcher) SetDebounce(duration time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.debounce = duration
}

// SetVerbose enables verbose logging
func (w *Watcher) SetVerbose(verbose bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.verbose = verbose
}

// Start begins watching for file changes
func (w *Watcher) Start(ctx context.Context) error {
	// Debounce timer
	var (
		debounceTimer *time.Timer
		debounceMu    sync.Mutex
		lastEvent     fsnotify.Event
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			// Filter events
			if !w.shouldProcess(event) {
				continue
			}

			// Debounce rapid file changes
			debounceMu.Lock()

			lastEvent = event

			if debounceTimer != nil {
				debounceTimer.Stop()
			}

			// Capture event locally to avoid race condition
			// when the timer callback executes after the lock is released
			eventToProcess := lastEvent
			debounceTimer = time.AfterFunc(w.debounce, func() {
				w.processEvent(eventToProcess)
			})

			debounceMu.Unlock()

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}

			if w.verbose {
				fmt.Printf("[Watcher] Error: %v\n", err)
			}
		}
	}
}

// shouldProcess checks if an event should be processed
func (w *Watcher) shouldProcess(event fsnotify.Event) bool {
	// Only process write and create events
	if event.Op&fsnotify.Write != fsnotify.Write &&
		event.Op&fsnotify.Create != fsnotify.Create {
		return false
	}

	// Ignore common output/build directories to prevent reload loops
	// Only check relative path components, not absolute system paths
	path := filepath.Clean(event.Name)

	excludeDirs := []string{
		"dist", "build", "output", ".git", "node_modules",
		"vendor", ".cache",
	}

	// Check each path component
	parts := strings.Split(path, string(filepath.Separator)) //nolint:modernize // strings.SplitSeq not available in Go 1.24
	for _, part := range parts {
		if slices.Contains(excludeDirs, part) {
			return false
		}
	}

	// Ignore generated files that might be in static/ to prevent loops
	// Only watch source files, not output files
	base := filepath.Base(path)
	if strings.HasPrefix(base, "app.") || strings.HasPrefix(base, "bundle.") {
		// These are typically generated files
		return false
	}

	// Check patterns if any are defined
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.patterns) == 0 {
		return true
	}

	// Check if file matches any pattern
	for _, pattern := range w.patterns {
		matched, err := filepath.Match(pattern, filepath.Base(event.Name))
		if err == nil && matched {
			return true
		}
	}

	return false
}

// processEvent executes all callbacks for an event
func (w *Watcher) processEvent(event fsnotify.Event) {
	w.mu.RLock()
	callbacks := make([]WatchCallback, len(w.callbacks))
	copy(callbacks, w.callbacks)
	verbose := w.verbose
	w.mu.RUnlock()

	if verbose {
		fmt.Printf("[Watcher] File changed: %s\n", event.Name)
	}

	for _, callback := range callbacks {
		if err := callback(event); err != nil {
			if verbose {
				fmt.Printf("[Watcher] Callback error: %v\n", err)
			}
		}
	}
}

// Close stops the watcher and releases resources
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

// WatchDirectory recursively watches a directory and its subdirectories
func (w *Watcher) WatchDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories
		if info != nil {
			name := filepath.Base(path)
			if len(name) > 0 && name[0] == '.' {
				if info.IsDir() {
					return filepath.SkipDir
				}

				return nil
			}
		}

		// Watch directories only
		if info != nil && info.IsDir() {
			return w.AddPath(path)
		}

		return nil
	})
}
