package assets

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// DevServer provides development features like hot reload and file watching.
// It watches for file changes and notifies connected browsers via Server-Sent Events (SSE).
type DevServer struct {
	pipeline   *Pipeline
	watcher    *Watcher
	sseClients []chan string
	mu         sync.RWMutex
	verbose    bool
	building   bool
	buildMu    sync.Mutex
}

// NewDevServer creates a new development server
func NewDevServer(pipeline *Pipeline) (*DevServer, error) {
	watcher, err := NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &DevServer{
		pipeline:   pipeline,
		watcher:    watcher,
		sseClients: make([]chan string, 0),
	}, nil
}

// Start begins the development server with file watching and hot reload
func (ds *DevServer) Start(ctx context.Context) error {
	// Configure watcher
	ds.watcher.SetDebounce(500 * time.Millisecond)
	ds.watcher.SetVerbose(ds.verbose)

	// Watch patterns for Go files (Tailwind classes), CSS, and JS
	ds.watcher.AddPattern("*.go")
	ds.watcher.AddPattern("*.css")
	ds.watcher.AddPattern("*.js")
	ds.watcher.AddPattern("*.ts")

	// Add callback to rebuild on changes
	ds.watcher.OnChange(func(event fsnotify.Event) error {
		return ds.onFileChange(ctx, event)
	})

	// Watch project directories
	if err := ds.watchProjectDirs(); err != nil {
		return fmt.Errorf("failed to setup watchers: %w", err)
	}

	fmt.Println("[DevServer] Hot reload enabled - watching Go files for changes")
	fmt.Println("[DevServer] Edit any .go file to trigger reload")

	// Start watcher in background
	go func() {
		if err := ds.watcher.Start(ctx); err != nil {
			fmt.Printf("[DevServer] Watcher error: %v\n", err)
		}
	}()

	return nil
}

// watchProjectDirs sets up watchers for common project directories
func (ds *DevServer) watchProjectDirs() error {
	// Only watch source directories, NOT output directories
	// Watching output directories causes infinite reload loops

	// Watch current directory for .go files only (non-recursive)
	// This watches the example app source files
	cwd, _ := os.Getwd()
	if ds.verbose {
		fmt.Printf("[DevServer] Watching current directory for Go files: %s\n", cwd)
	}

	if err := ds.watcher.AddPath("."); err != nil {
		if ds.verbose {
			fmt.Printf("[DevServer] Warning: Could not watch current directory: %v\n", err)
		}
	}

	// Note: We deliberately DON'T watch static/ directory because:
	// 1. That's where output files go (app.css, app.js)
	// 2. Watching it causes infinite reload loops
	// 3. Source CSS/JS files should be elsewhere (if any)

	return nil
}

// onFileChange handles file change events
func (ds *DevServer) onFileChange(ctx context.Context, event fsnotify.Event) error {
	// Prevent concurrent builds
	ds.buildMu.Lock()

	if ds.building {
		ds.buildMu.Unlock()
		return nil
	}

	ds.building = true
	ds.buildMu.Unlock()

	defer func() {
		ds.buildMu.Lock()
		ds.building = false
		ds.buildMu.Unlock()
	}()

	if ds.verbose {
		fmt.Printf("[DevServer] Rebuilding due to: %s\n", event.Name)
	}

	// Rebuild pipeline
	if err := ds.pipeline.Build(ctx); err != nil {
		if ds.verbose {
			fmt.Printf("[DevServer] Build failed: %v\n", err)
		}

		return err
	}

	if ds.verbose {
		fmt.Println("[DevServer] Build successful, reloading browsers...")
	}

	// Notify all SSE clients
	ds.notifyClients("reload")

	return nil
}

// notifyClients sends a message to all connected SSE clients
func (ds *DevServer) notifyClients(message string) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	for _, client := range ds.sseClients {
		select {
		case client <- message:
		default:
			// Client channel full, skip
		}
	}
}

// SSEHandler returns an HTTP handler for Server-Sent Events
func (ds *DevServer) SSEHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create client channel
		clientChan := make(chan string, 10)

		// Register client
		ds.mu.Lock()
		ds.sseClients = append(ds.sseClients, clientChan)
		clientIndex := len(ds.sseClients) - 1
		ds.mu.Unlock()

		// Cleanup on disconnect
		defer func() {
			ds.mu.Lock()
			// Remove client from list
			ds.sseClients = append(ds.sseClients[:clientIndex], ds.sseClients[clientIndex+1:]...)
			ds.mu.Unlock()
			close(clientChan)
		}()

		// Send initial connection message
		fmt.Fprintf(w, "data: connected\n\n")
		w.(http.Flusher).Flush()

		if ds.verbose {
			fmt.Println("[DevServer] SSE client connected")
		}

		// Listen for messages or client disconnect
		for {
			select {
			case <-r.Context().Done():
				if ds.verbose {
					fmt.Println("[DevServer] SSE client disconnected")
				}

				return

			case msg := <-clientChan:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				w.(http.Flusher).Flush()
			}
		}
	}
}

// HotReloadScript returns the client-side JavaScript for hot reload
func (ds *DevServer) HotReloadScript() string {
	return `<script>
(function() {
  const es = new EventSource('/_forgeui/reload');
  
  es.onmessage = function(event) {
    if (event.data === 'reload') {
      console.log('[ForgeUI] Reloading page...');
      location.reload();
    }
  };
  
  es.onerror = function() {
    console.log('[ForgeUI] Hot reload disconnected, retrying...');
    setTimeout(() => location.reload(), 1000);
  };
  
  console.log('[ForgeUI] Hot reload connected');
})();
</script>`
}

// SetVerbose enables verbose logging
func (ds *DevServer) SetVerbose(verbose bool) {
	ds.verbose = verbose
	if ds.watcher != nil {
		ds.watcher.SetVerbose(verbose)
	}
}

// Close stops the dev server and releases resources
func (ds *DevServer) Close() error {
	// Close all SSE clients
	ds.mu.Lock()

	for _, client := range ds.sseClients {
		close(client)
	}

	ds.sseClients = nil
	ds.mu.Unlock()

	// Close watcher
	if ds.watcher != nil {
		return ds.watcher.Close()
	}

	return nil
}

// ClientCount returns the number of connected SSE clients
func (ds *DevServer) ClientCount() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return len(ds.sseClients)
}
