package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/xraph/forgeui/cli"
	"github.com/xraph/forgeui/cli/templates"
	"github.com/xraph/forgeui/cli/util"
)

//nolint:gochecknoinits // init used for command registration
func init() {
	cli.RegisterCommand(GenerateCommand())
}

// GenerateCommand returns the generate command with subcommands
func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Short:   "Generate components, pages, and other code",
		Long:    `Generate boilerplate code for components, pages, and other ForgeUI artifacts.`,
		Usage:   "forgeui generate <type> <name> [flags]",
		Aliases: []string{"g"},
		Subcommands: []*cli.Command{
			ComponentCommand(),
			PageCommand(),
		},
	}
}

// ComponentCommand returns the component generation subcommand
func ComponentCommand() *cli.Command {
	return &cli.Command{
		Name:    "component",
		Short:   "Generate a new component",
		Long:    `Generate a new component with optional variants, props, and tests.`,
		Usage:   "forgeui generate component <name> [flags]",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			cli.StringFlag("type", "t", "Component type (basic, compound, form, layout, data)", "basic"),
			cli.StringFlag("dir", "d", "Output directory", "components"),
			cli.BoolFlag("with-variants", "", "Add CVA variants"),
			cli.BoolFlag("with-props", "", "Generate props struct"),
			cli.BoolFlag("with-test", "", "Generate test file"),
		},
		Run: runGenerateComponent,
	}
}

// PageCommand returns the page generation subcommand
func PageCommand() *cli.Command {
	return &cli.Command{
		Name:    "page",
		Short:   "Generate a new page",
		Long:    `Generate a new page handler with optional route registration.`,
		Usage:   "forgeui generate page <name> [flags]",
		Aliases: []string{"p"},
		Flags: []cli.Flag{
			cli.StringFlag("type", "t", "Page type (simple, dynamic, form, list, detail)", "simple"),
			cli.StringFlag("path", "", "Route path", ""),
			cli.StringFlag("dir", "d", "Output directory", "pages"),
			cli.BoolFlag("with-loader", "", "Add data loader"),
			cli.BoolFlag("with-meta", "", "Add SEO meta tags"),
		},
		Run: runGeneratePage,
	}
}

func runGenerateComponent(ctx *cli.Context) error {
	// Get component name
	if len(ctx.Args) == 0 {
		return errors.New("component name is required")
	}

	componentName := ctx.Args[0]
	if err := util.ValidateProjectName(componentName); err != nil {
		return fmt.Errorf("invalid component name: %w", err)
	}

	// Get component type
	componentType := ctx.GetString("type")
	outputDir := ctx.GetString("dir")
	withVariants := ctx.GetBool("with-variants")
	withProps := ctx.GetBool("with-props")
	withTest := ctx.GetBool("with-test")

	// Check if in ForgeUI project
	projectRoot, err := util.GetProjectRoot()
	if err != nil {
		return fmt.Errorf("not in a Go project: %w", err)
	}

	ctx.Printf("\n%sGenerating component...%s\n\n", util.ColorBlue, util.ColorReset)

	// Generate component
	spinner := util.NewSpinner(fmt.Sprintf("Creating %s component", componentName))
	spinner.Start()

	template, err := templates.GetComponentTemplate(componentType)
	if err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	componentDir := filepath.Join(projectRoot, outputDir, util.ToSnakeCase(componentName))
	if err := util.CreateDir(componentDir); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	opts := templates.ComponentOptions{
		Name:         componentName,
		Package:      util.ToSnakeCase(componentName),
		WithVariants: withVariants,
		WithProps:    withProps,
		WithTest:     withTest,
	}

	if err := template.Generate(componentDir, opts); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	spinner.Success("Component created")

	// Success message
	ctx.Printf("\n%s✓ Component created successfully!%s\n\n", util.ColorGreen, util.ColorReset)
	ctx.Printf("Location: %s\n", componentDir)
	ctx.Printf("  %s%s.go%s\n", util.ColorCyan, util.ToSnakeCase(componentName), util.ColorReset)

	if withTest {
		ctx.Printf("  %s%s_test.go%s\n", util.ColorCyan, util.ToSnakeCase(componentName), util.ColorReset)
	}

	ctx.Println()

	return nil
}

func runGeneratePage(ctx *cli.Context) error {
	// Get page name
	if len(ctx.Args) == 0 {
		return errors.New("page name is required")
	}

	pageName := ctx.Args[0]
	if err := util.ValidateProjectName(pageName); err != nil {
		return fmt.Errorf("invalid page name: %w", err)
	}

	// Get page type
	pageType := ctx.GetString("type")
	routePath := ctx.GetString("path")
	outputDir := ctx.GetString("dir")
	withLoader := ctx.GetBool("with-loader")
	withMeta := ctx.GetBool("with-meta")

	// Default route path
	if routePath == "" {
		routePath = "/" + util.ToSnakeCase(pageName)
	}

	// Check if in ForgeUI project
	projectRoot, err := util.GetProjectRoot()
	if err != nil {
		return fmt.Errorf("not in a Go project: %w", err)
	}

	ctx.Printf("\n%sGenerating page...%s\n\n", util.ColorBlue, util.ColorReset)

	// Generate page
	spinner := util.NewSpinner(fmt.Sprintf("Creating %s page", pageName))
	spinner.Start()

	template, err := templates.GetPageTemplate(pageType)
	if err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	pageDir := filepath.Join(projectRoot, outputDir)
	if err := util.CreateDir(pageDir); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	opts := templates.PageOptions{
		Name:       pageName,
		Package:    filepath.Base(outputDir),
		Path:       routePath,
		WithLoader: withLoader,
		WithMeta:   withMeta,
	}

	fileName := util.ToSnakeCase(pageName) + ".go"
	filePath := filepath.Join(pageDir, fileName)

	if err := template.Generate(filePath, opts); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	spinner.Success("Page created")

	// Success message
	ctx.Printf("\n%s✓ Page created successfully!%s\n\n", util.ColorGreen, util.ColorReset)
	ctx.Printf("Location: %s\n", filePath)
	ctx.Printf("Route: %s%s%s\n\n", util.ColorCyan, routePath, util.ColorReset)
	ctx.Printf("Next steps:\n")
	ctx.Printf("  1. Add route to your router: app.Router.Get(\"%s\", pages.%s)\n", routePath, util.ToPascalCase(pageName))
	ctx.Printf("  2. Run: %sforgeui dev%s\n\n", util.ColorCyan, util.ColorReset)

	return nil
}
