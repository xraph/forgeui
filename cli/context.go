package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/xraph/forgeui/cli/util"
)

// Re-export color constants for convenience
const (
	ColorReset  = util.ColorReset
	ColorRed    = util.ColorRed
	ColorGreen  = util.ColorGreen
	ColorYellow = util.ColorYellow
	ColorBlue   = util.ColorBlue
	ColorPurple = util.ColorPurple
	ColorCyan   = util.ColorCyan
	ColorGray   = util.ColorGray
)

// Context holds the execution context for a command
type Context struct {
	// Args are the remaining positional arguments after flag parsing
	Args []string

	// Flags are the parsed command-line flags
	Flags map[string]any

	// Config is the loaded project configuration (if any)
	Config *Config

	// Stdout is the standard output writer
	Stdout io.Writer

	// Stderr is the standard error writer
	Stderr io.Writer

	// Stdin is the standard input reader
	Stdin io.Reader
}

// GetString returns a string flag value
func (c *Context) GetString(name string) string {
	if val, ok := c.Flags[name].(string); ok {
		return val
	}

	return ""
}

// GetBool returns a bool flag value
func (c *Context) GetBool(name string) bool {
	if val, ok := c.Flags[name].(bool); ok {
		return val
	}

	return false
}

// GetInt returns an int flag value
func (c *Context) GetInt(name string) int {
	if val, ok := c.Flags[name].(int); ok {
		return val
	}

	return 0
}

// Printf prints formatted output to stdout
func (c *Context) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(c.Stdout, format, args...)
}

// Println prints a line to stdout
func (c *Context) Println(args ...any) {
	_, _ = fmt.Fprintln(c.Stdout, args...)
}

// Errorf prints formatted error to stderr
func (c *Context) Errorf(format string, args ...any) {
	_, _ = fmt.Fprintf(c.Stderr, format, args...)
}

// Errorln prints an error line to stderr
func (c *Context) Errorln(args ...any) {
	_, _ = fmt.Fprintln(c.Stderr, args...)
}

// Success prints a success message in green
func (c *Context) Success(msg string) {
	_, _ = fmt.Fprintf(c.Stdout, "%s%s%s\n", ColorGreen, msg, ColorReset)
}

// Info prints an info message in blue
func (c *Context) Info(msg string) {
	_, _ = fmt.Fprintf(c.Stdout, "%s%s%s\n", ColorBlue, msg, ColorReset)
}

// Warning prints a warning message in yellow
func (c *Context) Warning(msg string) {
	_, _ = fmt.Fprintf(c.Stdout, "%s%s%s\n", ColorYellow, msg, ColorReset)
}

// ErrorMsg prints an error message in red
func (c *Context) ErrorMsg(msg string) {
	_, _ = fmt.Fprintf(c.Stderr, "%s%s%s\n", ColorRed, msg, ColorReset)
}

// LoadConfig loads the project configuration file
func (c *Context) LoadConfig() error {
	config, err := LoadConfig(".")
	if err != nil {
		// Config file is optional
		if os.IsNotExist(err) {
			c.Config = DefaultConfig()
			return nil
		}

		return err
	}

	c.Config = config

	return nil
}
