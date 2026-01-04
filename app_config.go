package forgeui

import (
	"github.com/xraph/forgeui/bridge"
	"github.com/xraph/forgeui/theme"
)

// AppConfig holds enhanced ForgeUI application configuration
type AppConfig struct {
	// Debug enables debug/development mode
	Debug bool

	// Theme is the active theme name (legacy support)
	Theme string

	// AssetPublicDir is the source directory for static assets
	AssetPublicDir string

	// AssetOutputDir is the output directory for processed assets
	AssetOutputDir string

	// AssetManifest is the path to a manifest file for production builds
	AssetManifest string

	// Bridge configuration (optional)
	BridgeConfig *bridge.Config
	EnableBridge bool

	// Theme configuration (enhanced)
	LightTheme *theme.Theme
	DarkTheme  *theme.Theme

	// DefaultLayout is the default layout name for all pages
	DefaultLayout string

	// StaticPath is the URL path for static assets
	StaticPath string

	// Component defaults (from legacy Config)
	DefaultSize    Size
	DefaultVariant Variant
	DefaultRadius  Radius
}

// DefaultAppConfig returns sensible defaults for ForgeUI application
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		Debug:          false,
		Theme:          "default",
		AssetPublicDir: "public",
		AssetOutputDir: "dist",
		AssetManifest:  "",
		EnableBridge:   false,
		BridgeConfig:   nil,
		LightTheme:     nil,
		DarkTheme:      nil,
		DefaultLayout:  "",
		StaticPath:     "/static",
		DefaultSize:    SizeMD,
		DefaultVariant: VariantDefault,
		DefaultRadius:  RadiusMD,
	}
}

// AppOption is a functional option for configuring the App
type AppOption func(*AppConfig)

// WithDev enables or disables development mode
func WithDev(dev bool) AppOption {
	return func(c *AppConfig) { c.Debug = dev }
}

// WithDebug enables or disables debug mode (alias for WithDev)
func WithDebug(debug bool) AppOption {
	return func(c *AppConfig) { c.Debug = debug }
}

// WithAssets sets the asset directories
func WithAssets(publicDir string, opts ...string) AppOption {
	return func(c *AppConfig) {
		c.AssetPublicDir = publicDir
		if len(opts) > 0 {
			c.AssetOutputDir = opts[0]
		}
	}
}

// WithAssetPublicDir sets the source directory for static assets
func WithAssetPublicDir(dir string) AppOption {
	return func(c *AppConfig) { c.AssetPublicDir = dir }
}

// WithAssetOutputDir sets the output directory for processed assets
func WithAssetOutputDir(dir string) AppOption {
	return func(c *AppConfig) { c.AssetOutputDir = dir }
}

// WithAssetManifest sets the path to a manifest file for production builds
func WithAssetManifest(path string) AppOption {
	return func(c *AppConfig) { c.AssetManifest = path }
}

// WithBridge enables and configures the bridge system
func WithBridge(opts ...bridge.ConfigOption) AppOption {
	return func(c *AppConfig) {
		c.EnableBridge = true

		bridgeConfig := bridge.DefaultConfig()
		for _, opt := range opts {
			opt(bridgeConfig)
		}

		c.BridgeConfig = bridgeConfig
	}
}

// WithThemes sets the light and dark themes
func WithThemes(light, dark *theme.Theme) AppOption {
	return func(c *AppConfig) {
		c.LightTheme = light
		c.DarkTheme = dark
	}
}

// WithDefaultLayout sets the default layout for all pages
func WithDefaultLayout(layout string) AppOption {
	return func(c *AppConfig) {
		c.DefaultLayout = layout
	}
}

// WithAppStaticPath sets the URL path for static assets
func WithAppStaticPath(path string) AppOption {
	return func(c *AppConfig) {
		c.StaticPath = path
	}
}

// WithStaticPath sets the static assets path (alias for WithAppStaticPath)
func WithStaticPath(path string) AppOption {
	return func(c *AppConfig) { c.StaticPath = path }
}

// WithThemeName sets the theme name (legacy support)
func WithThemeName(theme string) AppOption {
	return func(c *AppConfig) { c.Theme = theme }
}

// WithDefaultSize sets the default component size
func WithDefaultSize(size Size) AppOption {
	return func(c *AppConfig) { c.DefaultSize = size }
}

// WithDefaultVariant sets the default component variant
func WithDefaultVariant(variant Variant) AppOption {
	return func(c *AppConfig) { c.DefaultVariant = variant }
}

// WithDefaultRadius sets the default border radius
func WithDefaultRadius(radius Radius) AppOption {
	return func(c *AppConfig) { c.DefaultRadius = radius }
}

// WithAssetsPath sets the filesystem assets path (deprecated, use WithAssetPublicDir)
func WithAssetsPath(path string) AppOption {
	return func(c *AppConfig) { c.AssetPublicDir = path }
}

