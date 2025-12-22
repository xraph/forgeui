package cli

// FlagType represents the type of a command-line flag
type FlagType int

const (
	FlagTypeString FlagType = iota
	FlagTypeBool
	FlagTypeInt
)

// Flag represents a command-line flag
type Flag struct {
	// Name is the long flag name (e.g., "output")
	Name string

	// Short is the short flag name (e.g., "o")
	Short string

	// Type is the flag type (string, bool, int)
	Type FlagType

	// Usage is the flag description
	Usage string

	// Default is the default value if not provided
	Default any

	// Required indicates if this flag is required
	Required bool
}

// StringFlag creates a string flag
func StringFlag(name, short, usage string, defaultValue string) Flag {
	return Flag{
		Name:    name,
		Short:   short,
		Type:    FlagTypeString,
		Usage:   usage,
		Default: defaultValue,
	}
}

// BoolFlag creates a bool flag
func BoolFlag(name, short, usage string) Flag {
	return Flag{
		Name:    name,
		Short:   short,
		Type:    FlagTypeBool,
		Usage:   usage,
		Default: false,
	}
}

// IntFlag creates an int flag
func IntFlag(name, short, usage string, defaultValue int) Flag {
	return Flag{
		Name:    name,
		Short:   short,
		Type:    FlagTypeInt,
		Usage:   usage,
		Default: defaultValue,
	}
}
