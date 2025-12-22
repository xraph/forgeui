package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/xraph/forgeui/cli"
	"github.com/xraph/forgeui/cli/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//nolint:gochecknoinits // init used for command registration
func init() {
	cli.RegisterCommand(PluginCommand())
}

// PluginCommand returns the plugin command with subcommands
func PluginCommand() *cli.Command {
	return &cli.Command{
		Name:  "plugin",
		Short: "Manage ForgeUI plugins",
		Long:  `Manage ForgeUI plugins - list, add, remove, and get information about plugins.`,
		Usage: "forgeui plugin <command> [args] [flags]",
		Subcommands: []*cli.Command{
			PluginListCommand(),
			PluginAddCommand(),
			PluginRemoveCommand(),
			PluginInfoCommand(),
		},
	}
}

// PluginListCommand lists installed plugins
func PluginListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Short: "List installed plugins",
		Long:  `List all plugins currently installed in the project.`,
		Usage: "forgeui plugin list",
		Run:   runPluginList,
	}
}

// PluginAddCommand adds a plugin
func PluginAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Short: "Add a plugin to the project",
		Long:  `Add a ForgeUI plugin to the project and update configuration.`,
		Usage: "forgeui plugin add <plugin-name>",
		Run:   runPluginAdd,
	}
}

// PluginRemoveCommand removes a plugin
func PluginRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Short:   "Remove a plugin from the project",
		Long:    `Remove a ForgeUI plugin from the project configuration.`,
		Usage:   "forgeui plugin remove <plugin-name>",
		Aliases: []string{"rm"},
		Run:     runPluginRemove,
	}
}

// PluginInfoCommand shows plugin information
func PluginInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "info",
		Short: "Show plugin information",
		Long:  `Display detailed information about a specific plugin.`,
		Usage: "forgeui plugin info <plugin-name>",
		Run:   runPluginInfo,
	}
}

func runPluginList(ctx *cli.Context) error {
	// Load config
	if err := ctx.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(ctx.Config.Plugins) == 0 {
		ctx.Printf("%sNo plugins installed%s\n", util.ColorGray, util.ColorReset)
		ctx.Printf("\nInstall a plugin with: %sforgeui plugin add <plugin-name>%s\n", util.ColorCyan, util.ColorReset)

		return nil
	}

	ctx.Printf("\n%sInstalled Plugins:%s\n\n", util.ColorBold, util.ColorReset)

	for i, plugin := range ctx.Config.Plugins {
		ctx.Printf("  %d. %s%s%s\n", i+1, util.ColorGreen, plugin, util.ColorReset)
	}

	ctx.Println()

	return nil
}

func runPluginAdd(ctx *cli.Context) error {
	if len(ctx.Args) == 0 {
		return errors.New("plugin name is required")
	}

	pluginName := ctx.Args[0]

	// Load config
	if err := ctx.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if already installed
	if slices.Contains(ctx.Config.Plugins, pluginName) {
		return fmt.Errorf("plugin %s is already installed", pluginName)
	}

	ctx.Printf("\n%sAdding plugin: %s%s\n\n", util.ColorBlue, pluginName, util.ColorReset)

	// Install the plugin package
	spinner := util.NewSpinner("Installing plugin package")
	spinner.Start()

	cmdCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pluginPath := "github.com/xraph/forgeui/plugins/" + pluginName
	cmd := exec.CommandContext(cmdCtx, "go", "get", pluginPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return fmt.Errorf("failed to install plugin: %w", err)
	}

	spinner.Success("Plugin package installed")

	// Add to config
	spinner = util.NewSpinner("Updating configuration")
	spinner.Start()

	ctx.Config.Plugins = append(ctx.Config.Plugins, pluginName)

	if err := ctx.Config.Save("."); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return fmt.Errorf("failed to save config: %w", err)
	}

	spinner.Success("Configuration updated")

	// Success message
	ctx.Printf("\n%s✓ Plugin added successfully!%s\n\n", util.ColorGreen, util.ColorReset)
	ctx.Printf("Plugin: %s%s%s\n\n", util.ColorCyan, pluginName, util.ColorReset)
	ctx.Printf("Next steps:\n")
	ctx.Printf("  1. Import the plugin in your code:\n")
	ctx.Printf("     %simport \"github.com/xraph/forgeui/plugins/%s\"%s\n", util.ColorGray, pluginName, util.ColorReset)
	ctx.Printf("  2. Register it with your ForgeUI app\n\n")

	return nil
}

func runPluginRemove(ctx *cli.Context) error {
	if len(ctx.Args) == 0 {
		return errors.New("plugin name is required")
	}

	pluginName := ctx.Args[0]

	// Load config
	if err := ctx.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find and remove plugin
	found := false
	newPlugins := []string{}

	for _, p := range ctx.Config.Plugins {
		if p == pluginName {
			found = true
		} else {
			newPlugins = append(newPlugins, p)
		}
	}

	if !found {
		return fmt.Errorf("plugin %s is not installed", pluginName)
	}

	ctx.Config.Plugins = newPlugins

	// Save config
	spinner := util.NewSpinner("Removing plugin")
	spinner.Start()

	if err := ctx.Config.Save("."); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return fmt.Errorf("failed to save config: %w", err)
	}

	spinner.Success("Plugin removed")

	ctx.Printf("\n%s✓ Plugin removed successfully!%s\n\n", util.ColorGreen, util.ColorReset)
	ctx.Printf("Note: You may want to remove the import and registration code manually.\n\n")

	return nil
}

func runPluginInfo(ctx *cli.Context) error {
	if len(ctx.Args) == 0 {
		return errors.New("plugin name is required")
	}

	pluginName := ctx.Args[0]

	// Built-in plugins info
	plugins := map[string]struct {
		Description string
		Features    []string
	}{
		"toast": {
			Description: "Notification system with customizable toasts",
			Features:    []string{"Success, error, warning, info notifications", "Auto-dismiss", "Position control", "Alpine.js integration"},
		},
		"sortable": {
			Description: "Drag-and-drop sorting with Sortable.js",
			Features:    []string{"Drag and drop lists", "Multiple lists", "Touch support", "Animation"},
		},
		"charts": {
			Description: "Data visualization with Chart.js",
			Features:    []string{"Line, bar, pie, area, doughnut charts", "Responsive", "Customizable options", "Interactive"},
		},
		"analytics": {
			Description: "Analytics tracking integration",
			Features:    []string{"Page view tracking", "Event tracking", "Custom properties", "Multiple providers"},
		},
		"seo": {
			Description: "SEO optimization tools",
			Features:    []string{"Meta tags", "Open Graph", "Twitter Cards", "Structured data"},
		},
		"htmxplugin": {
			Description: "HTMX integration wrapper",
			Features:    []string{"HTMX as Alpine plugin", "Easy integration", "Enhanced attributes"},
		},
	}

	info, ok := plugins[pluginName]
	if !ok {
		ctx.Printf("%sNo information available for plugin: %s%s\n", util.ColorYellow, pluginName, util.ColorReset)
		ctx.Printf("\nTry: %sforgeui plugin list%s\n", util.ColorCyan, util.ColorReset)

		return nil
	}

	ctx.Printf("\n%s%s%s\n", util.ColorBold, cases.Title(language.Und).String(pluginName), util.ColorReset)
	ctx.Printf("%s\n\n", info.Description)

	ctx.Printf("%sFeatures:%s\n", util.ColorBold, util.ColorReset)

	for _, feature := range info.Features {
		ctx.Printf("  • %s\n", feature)
	}

	ctx.Printf("\n%sInstallation:%s\n", util.ColorBold, util.ColorReset)
	ctx.Printf("  %sforgeui plugin add %s%s\n\n", util.ColorCyan, pluginName, util.ColorReset)

	return nil
}
