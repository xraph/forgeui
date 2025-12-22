package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the ForgeUI project configuration
type Config struct {
	Name    string      `json:"name"`
	Version string      `json:"version"`
	Dev     DevConfig   `json:"dev"`
	Build   BuildConfig `json:"build"`
	Assets  AssetsConfig `json:"assets"`
	Plugins []string    `json:"plugins"`
	Router  RouterConfig `json:"router"`
}

// DevConfig holds development server configuration
type DevConfig struct {
	Port        int    `json:"port"`
	Host        string `json:"host"`
	AutoReload  bool   `json:"auto_reload"`
	OpenBrowser bool   `json:"open_browser"`
}

// BuildConfig holds build configuration
type BuildConfig struct {
	OutputDir    string `json:"output_dir"`
	PublicDir    string `json:"public_dir"`
	Minify       bool   `json:"minify"`
	Binary       bool   `json:"binary"`
	EmbedAssets  bool   `json:"embed_assets"`
}

// AssetsConfig holds asset configuration
type AssetsConfig struct {
	CSS []string `json:"css"`
	JS  []string `json:"js"`
}

// RouterConfig holds router configuration
type RouterConfig struct {
	BasePath string `json:"base_path"`
	NotFound string `json:"not_found"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Name:    "forgeui-app",
		Version: "1.0.0",
		Dev: DevConfig{
			Port:        3000,
			Host:        "localhost",
			AutoReload:  true,
			OpenBrowser: false,
		},
		Build: BuildConfig{
			OutputDir:   "dist",
			PublicDir:   "public",
			Minify:      true,
			Binary:      false,
			EmbedAssets: true,
		},
		Assets: AssetsConfig{
			CSS: []string{"public/css/app.css"},
			JS:  []string{"public/js/app.js"},
		},
		Plugins: []string{},
		Router: RouterConfig{
			BasePath: "/",
			NotFound: "pages/404.go",
		},
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(dir string) (*Config, error) {
	// Try .forgeui.json first
	configPath := filepath.Join(dir, ".forgeui.json")
	if _, err := os.Stat(configPath); err == nil {
		return loadJSONConfig(configPath)
	}
	
	// Try forgeui.json
	configPath = filepath.Join(dir, "forgeui.json")
	if _, err := os.Stat(configPath); err == nil {
		return loadJSONConfig(configPath)
	}
	
	return nil, os.ErrNotExist
}

// loadJSONConfig loads a JSON configuration file
func loadJSONConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	return config, nil
}

// Save saves the configuration to a file
func (c *Config) Save(dir string) error {
	configPath := filepath.Join(dir, ".forgeui.json")
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	
	return nil
}

