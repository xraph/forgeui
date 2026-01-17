package assets

import (
	"net/http"
	"os"
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
		fullPath := filepath.Join(m.publicDir, actualPath)

		// Check if file exists
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)

			return
		}

		// Don't serve directories
		if info.IsDir() {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Set cache headers
		if !m.isDev && m.isFingerprinted(path) {
			// Fingerprinted assets can be cached for a long time
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Non-fingerprinted assets: moderate cache
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}

		// Serve the file
		http.ServeFile(w, r, fullPath)
	})
}
