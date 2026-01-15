package main

import (
	"os"

	"github.com/xraph/forgeui/cli"

	// Import commands to register them
	_ "github.com/xraph/forgeui/cli/commands"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
