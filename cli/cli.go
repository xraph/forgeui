// Package cli provides a command-line interface framework for ForgeUI.
//
// The CLI supports commands, subcommands, flags, and interactive prompts.
// It is designed to be lightweight with zero external dependencies.
package cli

import (
	"fmt"
	"os"
)

// Version is the CLI version
const Version = "1.0.0"

var rootCmd *Command

func init() {
	rootCmd = &Command{
		Name:  "forgeui",
		Short: "ForgeUI - A modern Go UI framework CLI",
		Long: `ForgeUI CLI provides tools for creating, developing, and building ForgeUI applications.

Features:
  - Project initialization with templates
  - Component and page generation
  - Development server with hot reload
  - Production builds with asset optimization
  - Plugin management`,
		Usage: "forgeui <command> [flags]",
	}
	
	// Register built-in commands (will be added by command files)
	rootCmd.Subcommands = []*Command{}
}

// Execute runs the CLI application
func Execute() error {
	args := os.Args[1:]
	
	// Check for version flag
	for _, arg := range args {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("ForgeUI CLI v%s\n", Version)
			return nil
		}
	}
	
	if len(args) == 0 {
		rootCmd.printHelp(os.Stdout)
		return nil
	}
	
	return rootCmd.Execute(args)
}

// RegisterCommand adds a command to the root command
func RegisterCommand(cmd *Command) {
	rootCmd.Subcommands = append(rootCmd.Subcommands, cmd)
}

// GetRootCommand returns the root command (for testing)
func GetRootCommand() *Command {
	return rootCmd
}

