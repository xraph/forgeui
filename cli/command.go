package cli

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// Command represents a CLI command with flags, subcommands, and a run function.
type Command struct {
	// Name is the command name (e.g., "init", "generate")
	Name string

	// Short is a brief description of the command
	Short string

	// Long is a detailed description of the command
	Long string

	// Usage is the usage string (e.g., "forgeui init [project-name]")
	Usage string

	// Flags are the command-line flags for this command
	Flags []Flag

	// Run is the function to execute when the command is invoked
	Run func(*Context) error

	// Subcommands are nested commands (e.g., "plugin list", "plugin add")
	Subcommands []*Command

	// Hidden controls whether this command appears in help output
	Hidden bool

	// Aliases are alternative names for this command
	Aliases []string
}

// Execute runs the command with the given arguments
func (c *Command) Execute(args []string) error {
	ctx := &Context{
		Args:   args,
		Flags:  make(map[string]any),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	// Parse flags
	remaining, err := c.parseFlags(ctx, args)
	if err != nil {
		return err
	}

	ctx.Args = remaining

	// Check for help flag
	if help, ok := ctx.Flags["help"].(bool); ok && help {
		c.printHelp(ctx.Stdout)
		return nil
	}

	// Check for subcommand
	if len(remaining) > 0 && len(c.Subcommands) > 0 {
		subcmdName := remaining[0]
		for _, sub := range c.Subcommands {
			if sub.Name == subcmdName || contains(sub.Aliases, subcmdName) {
				return sub.Execute(remaining[1:])
			}
		}
	}

	// Run the command
	if c.Run == nil {
		// No run function, show help
		c.printHelp(ctx.Stdout)
		return nil
	}

	return c.Run(ctx)
}

// parseFlags parses command-line flags and returns remaining arguments
func (c *Command) parseFlags(ctx *Context, args []string) ([]string, error) {
	remaining := []string{}

	// Add default help flag
	hasHelp := false

	for _, f := range c.Flags {
		if f.Name == "help" {
			hasHelp = true
			break
		}
	}

	if !hasHelp {
		c.Flags = append(c.Flags, Flag{
			Name:  "help",
			Short: "h",
			Type:  FlagTypeBool,
			Usage: "Show help information",
		})
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Not a flag, add to remaining
		if !strings.HasPrefix(arg, "-") {
			remaining = append(remaining, arg)
			continue
		}

		// Parse flag
		flagName := strings.TrimLeft(arg, "-")
		isShort := strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--")

		// Find matching flag
		var flag *Flag

		for j := range c.Flags {
			if isShort && c.Flags[j].Short == flagName {
				flag = &c.Flags[j]
				break
			} else if !isShort && c.Flags[j].Name == flagName {
				flag = &c.Flags[j]
				break
			}
		}

		if flag == nil {
			return nil, fmt.Errorf("unknown flag: %s", arg)
		}

		// Parse flag value
		switch flag.Type {
		case FlagTypeBool:
			ctx.Flags[flag.Name] = true

		case FlagTypeString, FlagTypeInt:
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag %s requires a value", arg)
			}

			i++
			if flag.Type == FlagTypeString {
				ctx.Flags[flag.Name] = args[i]
			} else {
				var val int

				_, err := fmt.Sscanf(args[i], "%d", &val)
				if err != nil {
					return nil, fmt.Errorf("flag %s requires an integer value", arg)
				}

				ctx.Flags[flag.Name] = val
			}
		}
	}

	// Set default values for missing flags
	for _, flag := range c.Flags {
		if _, ok := ctx.Flags[flag.Name]; !ok && flag.Default != nil {
			ctx.Flags[flag.Name] = flag.Default
		}
	}

	return remaining, nil
}

// printHelp prints the command's help information
func (c *Command) printHelp(w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s\n\n", c.Short)

	if c.Long != "" {
		_, _ = fmt.Fprintf(w, "%s\n\n", c.Long)
	}

	// Usage
	usage := c.Usage
	if usage == "" {
		usage = "forgeui " + c.Name
		if len(c.Subcommands) > 0 {
			usage += " <command>"
		}

		if len(c.Flags) > 0 {
			usage += " [flags]"
		}
	}

	_, _ = fmt.Fprintf(w, "Usage:\n  %s\n\n", usage)

	// Subcommands
	if len(c.Subcommands) > 0 {
		_, _ = fmt.Fprintf(w, "Available Commands:\n")

		for _, sub := range c.Subcommands {
			if !sub.Hidden {
				_, _ = fmt.Fprintf(w, "  %-15s %s\n", sub.Name, sub.Short)
			}
		}

		_, _ = fmt.Fprintf(w, "\n")
	}

	// Flags
	if len(c.Flags) > 0 {
		_, _ = fmt.Fprintf(w, "Flags:\n")

		for _, flag := range c.Flags {
			shortFlag := ""
			if flag.Short != "" {
				shortFlag = fmt.Sprintf("-%s, ", flag.Short)
			}

			_, _ = fmt.Fprintf(w, "  %s--%-15s %s\n", shortFlag, flag.Name, flag.Usage)
		}

		_, _ = fmt.Fprintf(w, "\n")
	}

	// Footer
	_, _ = fmt.Fprintf(w, "Use \"forgeui %s <command> --help\" for more information about a command.\n", c.Name)
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
