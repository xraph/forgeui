package forgeui

// Config holds ForgeUI application configuration
type Config struct {
	// Debug enables debug mode
	Debug bool

	// Theme is the active theme name
	Theme string

	// StaticPath is the URL path for static assets
	StaticPath string

	// AssetsPath is the filesystem path for assets (deprecated, use AssetPublicDir)
	AssetsPath string

	// AssetPublicDir is the source directory for static assets
	AssetPublicDir string

	// AssetOutputDir is the output directory for processed assets
	AssetOutputDir string

	// AssetManifest is the path to a manifest file for production builds
	AssetManifest string

	// DefaultSize is the default component size
	DefaultSize Size

	// DefaultVariant is the default component variant
	DefaultVariant Variant

	// DefaultRadius is the default border radius
	DefaultRadius Radius
}

// DefaultConfig returns sensible defaults for ForgeUI configuration
func DefaultConfig() *Config {
	return &Config{
		Debug:          false,
		Theme:          "default",
		StaticPath:     "/static",
		AssetsPath:     "public",
		AssetPublicDir: "public",
		AssetOutputDir: "dist",
		AssetManifest:  "",
		DefaultSize:    SizeMD,
		DefaultVariant: VariantDefault,
		DefaultRadius:  RadiusMD,
	}
}

// ConfigOption is a functional option for Config
type ConfigOption func(*Config)

// WithDebug enables or disables debug mode
func WithDebug(debug bool) ConfigOption {
	return func(c *Config) { c.Debug = debug }
}

// WithThemeName sets the theme name
func WithThemeName(theme string) ConfigOption {
	return func(c *Config) { c.Theme = theme }
}

// WithStaticPath sets the static assets path
func WithStaticPath(path string) ConfigOption {
	return func(c *Config) { c.StaticPath = path }
}

// WithAssetsPath sets the filesystem assets path
func WithAssetsPath(path string) ConfigOption {
	return func(c *Config) { c.AssetsPath = path }
}

// WithDefaultSize sets the default component size
func WithDefaultSize(size Size) ConfigOption {
	return func(c *Config) { c.DefaultSize = size }
}

// WithDefaultVariant sets the default component variant
func WithDefaultVariant(variant Variant) ConfigOption {
	return func(c *Config) { c.DefaultVariant = variant }
}

// WithDefaultRadius sets the default border radius
func WithDefaultRadius(radius Radius) ConfigOption {
	return func(c *Config) { c.DefaultRadius = radius }
}

// WithAssetPublicDir sets the source directory for static assets
func WithAssetPublicDir(dir string) ConfigOption {
	return func(c *Config) { c.AssetPublicDir = dir }
}

// WithAssetOutputDir sets the output directory for processed assets
func WithAssetOutputDir(dir string) ConfigOption {
	return func(c *Config) { c.AssetOutputDir = dir }
}

// WithAssetManifest sets the path to a manifest file for production builds
func WithAssetManifest(path string) ConfigOption {
	return func(c *Config) { c.AssetManifest = path }
}
