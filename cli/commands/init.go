package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/xraph/forgeui/cli"
	"github.com/xraph/forgeui/cli/templates"
	"github.com/xraph/forgeui/cli/util"
)

//nolint:gochecknoinits // init used for command registration
func init() {
	cli.RegisterCommand(InitCommand())
}

// InitCommand returns the init command
func InitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Short: "Initialize a new ForgeUI project",
		Long: `Initialize a new ForgeUI project with the specified template.

This command creates a new directory with a complete ForgeUI project structure,
including go.mod, main.go, and necessary directories.`,
		Usage: "forgeui init [project-name] [flags]",
		Flags: []cli.Flag{
			cli.StringFlag("template", "t", "Project template to use", "minimal"),
			cli.StringFlag("module", "m", "Go module path", ""),
			cli.BoolFlag("force", "f", "Force initialization even if directory exists"),
		},
		Run: runInit,
	}
}

func runInit(ctx *cli.Context) error {
	// Get project name
	var projectName string
	if len(ctx.Args) > 0 {
		projectName = ctx.Args[0]
	} else {
		var err error

		projectName, err = util.PromptWithDefault("Project name", "my-forgeui-app")
		if err != nil {
			return fmt.Errorf("failed to get project name: %w", err)
		}
	}

	// Validate project name
	if err := util.ValidateProjectName(projectName); err != nil {
		return err
	}

	// Get template
	templateName := ctx.GetString("template")
	if templateName == "" {
		options := []string{
			"minimal - Basic setup with one page",
			"standard - Full setup with router, assets, examples",
			"blog - Blog template with posts, tags",
			"dashboard - Admin dashboard with charts, tables",
			"api - API-first template with HTMX",
		}

		selected, err := util.Select("Select a project template:", options)
		if err != nil {
			return fmt.Errorf("failed to select template: %w", err)
		}

		templateNames := []string{"minimal", "standard", "blog", "dashboard", "api"}
		templateName = templateNames[selected]
	}

	// Get or generate module path
	modulePath := ctx.GetString("module")
	if modulePath == "" {
		defaultModule := "github.com/myuser/" + projectName
		modulePath, _ = util.PromptWithDefault("Go module path", defaultModule)
	}

	if err := util.ValidateGoModule(modulePath); err != nil {
		return err
	}

	// Create project directory
	projectDir := filepath.Join(".", projectName)
	if util.DirExists(projectDir) && !ctx.GetBool("force") {
		return fmt.Errorf("directory %s already exists (use --force to override)", projectName)
	}

	ctx.Printf("\n%sCreating ForgeUI project...%s\n\n", util.ColorBlue, util.ColorReset)

	// Create project structure
	spinner := util.NewSpinner("Creating project structure")
	spinner.Start()

	if err := createProjectStructure(projectDir); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	spinner.Success("Project structure created")

	// Generate template
	spinner = util.NewSpinner(fmt.Sprintf("Generating %s template", templateName))
	spinner.Start()

	template, err := templates.GetProjectTemplate(templateName)
	if err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	if err := template.Generate(projectDir, projectName, modulePath); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	spinner.Success("Template generated")

	// Initialize go.mod
	spinner = util.NewSpinner("Initializing Go module")
	spinner.Start()

	if err := initGoModule(projectDir, modulePath); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	spinner.Success("Go module initialized")

	// Install dependencies
	spinner = util.NewSpinner("Installing dependencies")
	spinner.Start()

	if err := installDependencies(projectDir); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	spinner.Success("Dependencies installed")

	// Success message
	ctx.Printf("\n%s✓ Project created successfully!%s\n\n", util.ColorGreen, util.ColorReset)
	ctx.Printf("Next steps:\n")
	ctx.Printf("  %scd %s%s\n", util.ColorCyan, projectName, util.ColorReset)
	ctx.Printf("  %sforgeui dev%s\n\n", util.ColorCyan, util.ColorReset)

	return nil
}

func createProjectStructure(projectDir string) error {
	dirs := []string{
		projectDir,
		filepath.Join(projectDir, "public"),
		filepath.Join(projectDir, "public", "css"),
		filepath.Join(projectDir, "public", "js"),
		filepath.Join(projectDir, "components"),
		filepath.Join(projectDir, "pages"),
	}

	for _, dir := range dirs {
		if err := util.CreateDir(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func initGoModule(projectDir, modulePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "mod", "init", modulePath)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	return nil
}

func installDependencies(projectDir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// go get forgeui
	cmd := exec.CommandContext(ctx, "go", "get", "github.com/xraph/forgeui@latest")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install forgeui: %w", err)
	}

	// go mod tidy
	cmd = exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tidy modules: %w", err)
	}

	return nil
}
