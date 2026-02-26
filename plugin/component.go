package plugin

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui"
)

// ComponentPlugin extends ForgeUI with new UI components.
type ComponentPlugin interface {
	Plugin

	// Components returns a map of component names to their constructors.
	Components() map[string]ComponentConstructor

	// CVAExtensions returns CVA configurations for component variants.
	CVAExtensions() map[string]*forgeui.CVA
}

// ComponentConstructor is a function that creates a component.
// It receives props (any type) and optional children components.
type ComponentConstructor func(props any, children ...templ.Component) templ.Component

// ComponentPluginBase provides a base implementation for component plugins.
type ComponentPluginBase struct {
	*PluginBase

	components map[string]ComponentConstructor
	cva        map[string]*forgeui.CVA
}

// NewComponentPluginBase creates a new ComponentPluginBase.
func NewComponentPluginBase(info PluginInfo, components map[string]ComponentConstructor) *ComponentPluginBase {
	return &ComponentPluginBase{
		PluginBase: NewPluginBase(info),
		components: components,
		cva:        make(map[string]*forgeui.CVA),
	}
}

// NewComponentPluginBaseWithCVA creates a ComponentPluginBase with CVA extensions.
func NewComponentPluginBaseWithCVA(
	info PluginInfo,
	components map[string]ComponentConstructor,
	cva map[string]*forgeui.CVA,
) *ComponentPluginBase {
	return &ComponentPluginBase{
		PluginBase: NewPluginBase(info),
		components: components,
		cva:        cva,
	}
}

// Components returns the component constructors.
func (c *ComponentPluginBase) Components() map[string]ComponentConstructor {
	return c.components
}

// CVAExtensions returns the CVA configurations.
func (c *ComponentPluginBase) CVAExtensions() map[string]*forgeui.CVA {
	return c.cva
}

// AddComponent adds a component constructor to the plugin.
func (c *ComponentPluginBase) AddComponent(name string, constructor ComponentConstructor) {
	if c.components == nil {
		c.components = make(map[string]ComponentConstructor)
	}

	c.components[name] = constructor
}

// AddCVA adds a CVA configuration for a component.
func (c *ComponentPluginBase) AddCVA(name string, cva *forgeui.CVA) {
	if c.cva == nil {
		c.cva = make(map[string]*forgeui.CVA)
	}

	c.cva[name] = cva
}
