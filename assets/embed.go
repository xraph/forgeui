package assets

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

// EmbeddedManager extends Manager to work with embed.FS
type EmbeddedManager struct {
	*Manager
	embedFS fs.FS
}

// NewEmbeddedManager creates a manager that serves assets from an embedded filesystem
func NewEmbeddedManager(embedFS fs.FS, cfg Config) *EmbeddedManager {
	m := NewManager(cfg)
	return &EmbeddedManager{
		Manager: m,
		embedFS: embedFS,
	}
}

// EmbeddedHandler returns an http.Handler for serving embedded static files
func (m *EmbeddedManager) EmbeddedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get path and remove /static/ prefix
		path := strings.TrimPrefix(r.URL.Path, "/static/")

		// Validate path
		if !isValidPath(path) {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Remove fingerprint for lookup if present
		actualPath := m.stripFingerprint(path)

		// Try to open file from embedded filesystem
		file, err := m.embedFS.Open(actualPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

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

// fingerprintEmbedded generates a fingerprint for an embedded asset
func (m *EmbeddedManager) fingerprintEmbedded(path string) string {
	// Validate path
	if !isValidPath(path) {
		return path
	}

	file, err := m.embedFS.Open(path)
	if err != nil {
		return path
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return path
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))[:8]

	// Split into name and extension
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]

	return fmt.Sprintf("%s.%s%s", base, hash, ext)
}

// FingerprintAllEmbedded generates fingerprints for all embedded assets
func (m *EmbeddedManager) FingerprintAllEmbedded() error {
	return fs.WalkDir(m.embedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and the root
		if d.IsDir() || path == "." {
			return nil
		}

		// Generate fingerprint
		fp := m.fingerprintEmbedded(path)

		// Cache it
		m.mu.Lock()
		m.fingerprints[path] = fp
		m.mu.Unlock()

		return nil
	})
}

