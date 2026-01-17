package assets

import (
	"testing"
)

func TestManager_CustomStaticPath(t *testing.T) {
	tests := []struct {
		name           string
		staticPath     string
		assetPath      string
		expectedPrefix string
	}{
		{
			name:           "default static path",
			staticPath:     "",
			assetPath:      "app.css",
			expectedPrefix: "/static/",
		},
		{
			name:           "custom assets path",
			staticPath:     "/assets",
			assetPath:      "app.css",
			expectedPrefix: "/assets/",
		},
		{
			name:           "custom public path",
			staticPath:     "/public",
			assetPath:      "app.css",
			expectedPrefix: "/public/",
		},
		{
			name:           "path with trailing slash",
			staticPath:     "/custom/",
			assetPath:      "app.css",
			expectedPrefix: "/custom/",
		},
		{
			name:           "path without leading slash",
			staticPath:     "files",
			assetPath:      "app.css",
			expectedPrefix: "/files/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(Config{
				StaticPath: tt.staticPath,
				IsDev:      true,
			})

			url := m.URL(tt.assetPath)

			expected := tt.expectedPrefix + tt.assetPath
			if url != expected {
				t.Errorf("URL() = %q, want %q", url, expected)
			}
		})
	}
}

func TestManager_StaticPath_Production(t *testing.T) {
	m := NewManager(Config{
		StaticPath: "/assets",
		IsDev:      false,
	})

	// Create a test fingerprint
	m.mu.Lock()
	m.fingerprints["app.css"] = "app.abc123.css"
	m.mu.Unlock()

	url := m.URL("app.css")

	expected := "/assets/app.abc123.css"
	if url != expected {
		t.Errorf("URL() = %q, want %q", url, expected)
	}
}
