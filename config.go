package forgeui

// Config holds ForgeUI application configuration
// Deprecated: Use AppConfig instead. Config is kept for backward compatibility.
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
// Deprecated: Use DefaultAppConfig instead
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
// Note: ConfigOption is now an alias for AppOption for unified configuration
type ConfigOption = AppOption

// Migration Guide:
// ----------------
// Old way (deprecated):
//   app := forgeui.New(forgeui.WithDebug(true), forgeui.WithStaticPath("/assets"))
//
// New way (recommended):
//   app := forgeui.New(forgeui.WithDev(true), forgeui.WithStaticPath("/assets"))
//
// The API is backward compatible, but new features (bridge, themes) require AppConfig.
// All With* functions now work with both Config and AppConfig.
