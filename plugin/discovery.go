package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

// Discover loads plugins from a directory (Go plugins).
//
// Note: This is optional and experimental. Go plugins have significant limitations:
//   - Only supported on Linux, FreeBSD, and macOS
//   - Must be built with the same Go version as the main program
//   - Must use the same versions of all dependencies
//   - Cannot be unloaded once loaded
//
// For production use, consider statically linking plugins instead.
//
// Plugin files must:
//   - Have a .so extension
//   - Export a variable named "Plugin" of type plugin.Plugin
//
// Example plugin:
//
//	package main
//
//	import "github.com/xraph/forgeui/plugin"
//
//	var Plugin = &MyPlugin{}
//
//	type MyPlugin struct {
//	    plugin.PluginBase
//	}
func (r *Registry) Discover(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".so" {
			return nil
		}

		p, err := plugin.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open plugin %s: %w", path, err)
		}

		sym, err := p.Lookup("Plugin")
		if err != nil {
			return fmt.Errorf("failed to lookup Plugin symbol in %s: %w", path, err)
		}

		forgePlugin, ok := sym.(Plugin)
		if !ok {
			return fmt.Errorf("%s does not implement Plugin interface", path)
		}

		if err := r.Register(forgePlugin); err != nil {
			return fmt.Errorf("failed to register plugin from %s: %w", path, err)
		}

		return nil
	})
}

// DiscoverSafe is like Discover but continues on error, collecting all errors.
// Returns all errors encountered during discovery.
func (r *Registry) DiscoverSafe(dir string) []error {
	var errs []error

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		if info.IsDir() || filepath.Ext(path) != ".so" {
			return nil
		}

		p, err := plugin.Open(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to open plugin %s: %w", path, err))
			return nil
		}

		sym, err := p.Lookup("Plugin")
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to lookup Plugin symbol in %s: %w", path, err))
			return nil
		}

		forgePlugin, ok := sym.(Plugin)
		if !ok {
			errs = append(errs, fmt.Errorf("%s does not implement Plugin interface", path))
			return nil
		}

		if err := r.Register(forgePlugin); err != nil {
			errs = append(errs, fmt.Errorf("failed to register plugin from %s: %w", path, err))
		}

		return nil
	})

	return errs
}
