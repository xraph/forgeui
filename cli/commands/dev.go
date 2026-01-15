package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/xraph/forgeui/cli"
	"github.com/xraph/forgeui/cli/util"
)

//nolint:gochecknoinits // init used for command registration
func init() {
	cli.RegisterCommand(DevCommand())
}

// DevCommand returns the dev command
func DevCommand() *cli.Command {
	return &cli.Command{
		Name:  "dev",
		Short: "Start development server with hot reload",
		Long: `Start the development server with hot reload enabled.

The dev server watches for file changes and automatically rebuilds assets.
It starts your Go application and provides a live development experience.`,
		Usage: "forgeui dev [flags]",
		Flags: []cli.Flag{
			cli.IntFlag("port", "p", "Port to listen on", 3000),
			cli.StringFlag("host", "h", "Host to bind to", "localhost"),
			cli.BoolFlag("open", "o", "Open browser automatically"),
		},
		Run: runDev,
	}
}

func runDev(ctx *cli.Context) error {
	// Load config
	if err := ctx.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get port from flag or config
	port := ctx.GetInt("port")
	if port == 0 {
		port = ctx.Config.Dev.Port
	}

	if port == 0 {
		port = 3000
	}

	host := ctx.GetString("host")
	if host == "" {
		host = ctx.Config.Dev.Host
	}

	if host == "" {
		host = "localhost"
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	ctx.Printf("\n%sStarting development server...%s\n\n", util.ColorBlue, util.ColorReset)

	// Check if main.go exists
	if !util.FileExists("main.go") {
		return errors.New("main.go not found - are you in a ForgeUI project?")
	}

	// Create context for graceful shutdown
	srvCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the Go application
	cmd := exec.CommandContext(srvCtx, "go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Wait a moment for server to start
	time.Sleep(500 * time.Millisecond)

	// Print success message
	url := "http://" + addr

	ctx.Printf("%s✓ Server started%s\n\n", util.ColorGreen, util.ColorReset)
	ctx.Printf("  %sLocal:%s   %s\n", util.ColorBold, util.ColorReset, url)
	ctx.Printf("  %sPress Ctrl+C to stop%s\n\n", util.ColorGray, util.ColorReset)

	// Open browser if requested
	if ctx.GetBool("open") || ctx.Config.Dev.OpenBrowser {
		if err := openBrowser(url); err != nil {
			// Log browser open error but don't fail
			ctx.Printf("%sWarning: Failed to open browser: %v%s\n", util.ColorYellow, err, util.ColorReset)
		}
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt
	<-sigChan

	ctx.Printf("\n%sStopping server...%s\n", util.ColorYellow, util.ColorReset)

	// Cancel context to stop the server
	cancel()

	// Wait for process to exit
	if err := cmd.Wait(); err != nil {
		// Process may have exited with error, which is expected
		_ = err // Ignore error as we're shutting down
	}

	ctx.Printf("%s✓ Server stopped%s\n", util.ColorGreen, util.ColorReset)

	return nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch {
	case commandExists("xdg-open"):
		cmd = exec.CommandContext(ctx, "xdg-open", url)
	case commandExists("open"):
		cmd = exec.CommandContext(ctx, "open", url)
	case commandExists("start"):
		cmd = exec.CommandContext(ctx, "cmd", "/c", "start", url)
	default:
		return errors.New("no browser opener found")
	}

	return cmd.Start()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

