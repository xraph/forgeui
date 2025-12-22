package assets

// StyleOption is a functional option for configuring stylesheet elements
type StyleOption func(*styleConfig)

type styleConfig struct {
	media     string
	preload   bool
	integrity string
	crossOrigin string
}

// WithMedia sets the media attribute for a stylesheet
func WithMedia(media string) StyleOption {
	return func(c *styleConfig) {
		c.media = media
	}
}

// WithPreload enables preloading for a stylesheet
func WithPreload() StyleOption {
	return func(c *styleConfig) {
		c.preload = true
	}
}

// WithIntegrity sets the integrity hash (SRI) for a stylesheet
func WithIntegrity(hash string) StyleOption {
	return func(c *styleConfig) {
		c.integrity = hash
	}
}

// WithCrossOrigin sets the crossorigin attribute for a stylesheet
func WithCrossOrigin(value string) StyleOption {
	return func(c *styleConfig) {
		c.crossOrigin = value
	}
}

// ScriptOption is a functional option for configuring script elements
type ScriptOption func(*scriptConfig)

type scriptConfig struct {
	defer_      bool
	async       bool
	module      bool
	integrity   string
	crossOrigin string
	noModule    bool
}

// WithDefer sets the defer attribute for a script
func WithDefer() ScriptOption {
	return func(c *scriptConfig) {
		c.defer_ = true
	}
}

// WithAsync sets the async attribute for a script
func WithAsync() ScriptOption {
	return func(c *scriptConfig) {
		c.async = true
	}
}

// WithModule sets type="module" for a script
func WithModule() ScriptOption {
	return func(c *scriptConfig) {
		c.module = true
	}
}

// WithScriptIntegrity sets the integrity hash (SRI) for a script
func WithScriptIntegrity(hash string) ScriptOption {
	return func(c *scriptConfig) {
		c.integrity = hash
	}
}

// WithScriptCrossOrigin sets the crossorigin attribute for a script
func WithScriptCrossOrigin(value string) ScriptOption {
	return func(c *scriptConfig) {
		c.crossOrigin = value
	}
}

// WithNoModule sets the nomodule attribute for a script
func WithNoModule() ScriptOption {
	return func(c *scriptConfig) {
		c.noModule = true
	}
}

