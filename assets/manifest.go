package assets

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Manifest represents a mapping of original asset paths to fingerprinted paths
type Manifest map[string]string

// LoadManifest loads a manifest from a JSON file
func LoadManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

// Save writes the manifest to a JSON file
func (manifest Manifest) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

// Get retrieves a fingerprinted path from the manifest
func (manifest Manifest) Get(path string) (string, bool) {
	fp, ok := manifest[path]
	return fp, ok
}

// Set adds or updates a mapping in the manifest
func (manifest Manifest) Set(path, fingerprinted string) {
	manifest[path] = fingerprinted
}

// GenerateManifest creates a manifest from all assets in a directory
func GenerateManifest(m *Manager) (Manifest, error) {
	manifest := make(Manifest)

	err := filepath.Walk(m.publicDir, func(path string, info os.FileInfo, err error) error {
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
		manifest[relPath] = fp

		return nil
	})
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
