package assets

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// Handler returns an http.Handler for serving static files
func (m *Manager) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get path and remove static path prefix
		path := strings.TrimPrefix(r.URL.Path, m.staticPath)

		// Validate path
		if !isValidPath(path) {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Remove fingerprint for lookup if present
		actualPath := m.stripFingerprint(path)

		// Try to open file from filesystem
		file, err := m.fileSystem.Open(actualPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer func() { _ = file.Close() }()

		// Get file info
		info, err := file.Stat()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Don't serve directories
		if info.IsDir() {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Set Content-Type based on file extension
		ext := filepath.Ext(actualPath)
		if contentType := mime.TypeByExtension(ext); contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}

		// Set cache headers
		if !m.isDev && m.isFingerprinted(path) {
			// Fingerprinted assets can be cached for a long time
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Non-fingerprinted assets: moderate cache
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}

		// Serve the file content
		http.ServeContent(w, r, filepath.Base(actualPath), info.ModTime(), file.(io.ReadSeeker))
	})
}
