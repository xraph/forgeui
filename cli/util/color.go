package util

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"

	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldBlue   = "\033[1;34m"
	ColorBoldPurple = "\033[1;35m"
	ColorBoldCyan   = "\033[1;36m"
	ColorBoldWhite  = "\033[1;37m"
)

var (
	// NoColor disables color output
	NoColor = os.Getenv("NO_COLOR") != ""
)

// Colorize wraps text in ANSI color codes
func Colorize(color, text string) string {
	if NoColor {
		return text
	}

	return color + text + ColorReset
}

// Red returns text in red
func Red(text string) string {
	return Colorize(ColorRed, text)
}

// Green returns text in green
func Green(text string) string {
	return Colorize(ColorGreen, text)
}

// Yellow returns text in yellow
func Yellow(text string) string {
	return Colorize(ColorYellow, text)
}

// Blue returns text in blue
func Blue(text string) string {
	return Colorize(ColorBlue, text)
}

// Purple returns text in purple
func Purple(text string) string {
	return Colorize(ColorPurple, text)
}

// Cyan returns text in cyan
func Cyan(text string) string {
	return Colorize(ColorCyan, text)
}

// Gray returns text in gray
func Gray(text string) string {
	return Colorize(ColorGray, text)
}

// Bold returns text in bold
func Bold(text string) string {
	if NoColor {
		return text
	}

	return "\033[1m" + text + "\033[22m"
}

// Success prints a success message
func Success(msg string) {
	fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, msg)
}

// Error prints an error message
func Error(msg string) {
	fmt.Printf("%s✗%s %s\n", ColorRed, ColorReset, msg)
}

// Info prints an info message
func Info(msg string) {
	fmt.Printf("%sℹ%s %s\n", ColorBlue, ColorReset, msg)
}

// Warning prints a warning message
func Warning(msg string) {
	fmt.Printf("%s⚠%s %s\n", ColorYellow, ColorReset, msg)
}
