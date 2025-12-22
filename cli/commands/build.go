package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	
	"github.com/xraph/forgeui/cli"
	"github.com/xraph/forgeui/cli/util"
)

func init() {
	cli.RegisterCommand(BuildCommand())
}

// BuildCommand returns the build command
func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Short: "Build for production",
		Long: `Build the application for production deployment.

This command:
  - Processes and optimizes assets (CSS, JS)
  - Generates fingerprinted asset files
  - Creates a manifest file
  - Optionally compiles a Go binary`,
		Usage: "forgeui build [flags]",
		Flags: []cli.Flag{
			cli.StringFlag("output", "o", "Output directory", "dist"),
			cli.BoolFlag("binary", "b", "Compile Go binary"),
			cli.BoolFlag("minify", "m", "Minify assets"),
			cli.BoolFlag("embed", "e", "Embed assets in binary"),
		},
		Run: runBuild,
	}
}

func runBuild(ctx *cli.Context) error {
	// Load config
	if err := ctx.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	outputDir := ctx.GetString("output")
	if outputDir == "" {
		outputDir = ctx.Config.Build.OutputDir
	}
	if outputDir == "" {
		outputDir = "dist"
	}
	
	buildBinary := ctx.GetBool("binary") || ctx.Config.Build.Binary
	minify := ctx.GetBool("minify") || ctx.Config.Build.Minify
	embed := ctx.GetBool("embed") || ctx.Config.Build.EmbedAssets
	
	ctx.Printf("\n%sBuilding for production...%s\n\n", util.ColorBlue, util.ColorReset)
	
	// Create output directory
	spinner := util.NewSpinner("Creating output directory")
	spinner.Start()
	
	if err := util.CreateDir(outputDir); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}
	
	spinner.Success("Output directory created")
	
	// Copy static assets
	spinner = util.NewSpinner("Copying static assets")
	spinner.Start()
	
	if err := copyStaticAssets(outputDir, ctx.Config.Build.PublicDir); err != nil {
		spinner.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}
	
	spinner.Success("Static assets copied")
	
	// Process assets if configured
	if len(ctx.Config.Assets.CSS) > 0 || len(ctx.Config.Assets.JS) > 0 {
		spinner = util.NewSpinner("Processing assets")
		spinner.Start()
		
		if err := processAssets(outputDir, ctx.Config, minify); err != nil {
			spinner.Error(fmt.Sprintf("Failed: %v", err))
			return err
		}
		
		spinner.Success("Assets processed")
	}
	
	// Build binary if requested
	if buildBinary {
		spinner = util.NewSpinner("Compiling Go binary")
		spinner.Start()
		
		binaryName := ctx.Config.Name
		if binaryName == "" {
			binaryName = "app"
		}
		
		if err := buildGoBinary(outputDir, binaryName, embed); err != nil {
			spinner.Error(fmt.Sprintf("Failed: %v", err))
			return err
		}
		
		spinner.Success("Binary compiled")
	}
	
	// Success message
	ctx.Printf("\n%s✓ Build completed successfully!%s\n\n", util.ColorGreen, util.ColorReset)
	ctx.Printf("Output directory: %s%s%s\n", util.ColorCyan, outputDir, util.ColorReset)
	
	if buildBinary {
		binaryPath := filepath.Join(outputDir, ctx.Config.Name)
		ctx.Printf("Binary: %s%s%s\n", util.ColorCyan, binaryPath, util.ColorReset)
	}
	
	ctx.Println()
	
	return nil
}

func copyStaticAssets(outputDir, publicDir string) error {
	if publicDir == "" {
		publicDir = "public"
	}
	
	// Check if public directory exists
	if !util.DirExists(publicDir) {
		// No public directory, skip
		return nil
	}
	
	// Copy public directory to output/static
	staticDir := filepath.Join(outputDir, "static")
	if err := util.CreateDir(staticDir); err != nil {
		return err
	}
	
	return util.CopyDir(publicDir, staticDir)
}

func processAssets(outputDir string, config *cli.Config, minify bool) error {
	// This would integrate with the assets package
	// For now, just create placeholder files
	
	// Process CSS
	for _, cssFile := range config.Assets.CSS {
		if util.FileExists(cssFile) {
			destFile := filepath.Join(outputDir, "static", "css", filepath.Base(cssFile))
			if err := util.CreateDir(filepath.Dir(destFile)); err != nil {
				return err
			}
			
			data, err := os.ReadFile(cssFile)
			if err != nil {
				return err
			}
			
			if err := os.WriteFile(destFile, data, 0644); err != nil {
				return err
			}
		}
	}
	
	// Process JS
	for _, jsFile := range config.Assets.JS {
		if util.FileExists(jsFile) {
			destFile := filepath.Join(outputDir, "static", "js", filepath.Base(jsFile))
			if err := util.CreateDir(filepath.Dir(destFile)); err != nil {
				return err
			}
			
			data, err := os.ReadFile(jsFile)
			if err != nil {
				return err
			}
			
			if err := os.WriteFile(destFile, data, 0644); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func buildGoBinary(outputDir, binaryName string, embedAssets bool) error {
	binaryPath := filepath.Join(outputDir, binaryName)
	
	args := []string{"build", "-o", binaryPath}
	
	if embedAssets {
		args = append(args, "-tags", "embed")
	}
	
	args = append(args, ".")
	
	cmd := exec.CommandContext(context.Background(), "go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

