package assets

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var fingerprintRegex = regexp.MustCompile(`^(.+)\.([a-f0-9]{8})(\.[^.]+)$`)

// fingerprint generates a content-based fingerprint for an asset
func (m *Manager) fingerprint(path string) string {
	// Validate path to prevent directory traversal
	if !isValidPath(path) {
		return path
	}

	fullPath := filepath.Join(m.publicDir, path)
	f, err := os.Open(fullPath)
	if err != nil {
		return path
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return path
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))[:8]

	// Split into name and extension
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]

	return fmt.Sprintf("%s.%s%s", base, hash, ext)
}

// stripFingerprint removes the fingerprint hash from a path
func (m *Manager) stripFingerprint(path string) string {
	matches := fingerprintRegex.FindStringSubmatch(path)
	if len(matches) == 4 {
		// matches[1] is the base name
		// matches[2] is the hash
		// matches[3] is the extension
		return matches[1] + matches[3]
	}
	return path
}

// isFingerprinted checks if a path contains a fingerprint hash
func (m *Manager) isFingerprinted(path string) bool {
	return fingerprintRegex.MatchString(path)
}

// isValidPath validates that a path doesn't contain directory traversal attempts
func isValidPath(path string) bool {
	// Reject paths containing ".."
	if strings.Contains(path, "..") {
		return false
	}

	// Reject absolute paths
	if filepath.IsAbs(path) {
		return false
	}

	// Clean the path and ensure it doesn't escape
	cleaned := filepath.Clean(path)
	if strings.HasPrefix(cleaned, "..") {
		return false
	}

	return true
}

// FingerprintAll generates fingerprints for all assets in the public directory
func (m *Manager) FingerprintAll() error {
	return filepath.Walk(m.publicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(m.publicDir, path)
		if err != nil {
			return err
		}

		// Generate fingerprint
		fp := m.fingerprint(relPath)

		// Cache it
		m.mu.Lock()
		m.fingerprints[relPath] = fp
		m.mu.Unlock()

		return nil
	})
}

