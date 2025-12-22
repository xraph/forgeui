package plugin

import (
	g "github.com/maragudk/gomponents"
	"github.com/xraph/forgeui"
)

// ComponentPlugin extends ForgeUI with new UI components.
//
// A ComponentPlugin provides:
//   - Component constructors for creating new components
//   - Optional CVA extensions for variant styling
//
// Example:
//
//	type ChartPlugin struct {
//	    *ComponentPluginBase
//	}
//
//	func NewChartPlugin() *ChartPlugin {
//	    return &ChartPlugin{
//	        ComponentPluginBase: NewComponentPluginBase(
//	            PluginInfo{Name: "charts", Version: "1.0.0"},
//	            map[string]ComponentConstructor{
//	                "LineChart": lineChartConstructor,
//	                "BarChart":  barChartConstructor,
//	            },
//	        ),
//	    }
//	}
type ComponentPlugin interface {
	Plugin

	// Components returns a map of component names to their constructors.
	// Component names should be CamelCase and unique.
	Components() map[string]ComponentConstructor

	// CVAExtensions returns CVA configurations for component variants.
	// The keys should match component names from Components().
	CVAExtensions() map[string]*forgeui.CVA
}

// ComponentConstructor is a function that creates a component.
// It receives props (any type) and optional children nodes.
//
// Example:
//
//	func lineChartConstructor(props any, children ...g.Node) g.Node {
//	    opts := props.(*ChartOptions)
//	    return html.Div(
//	        html.Class("chart-container"),
//	        g.Attr("data-chart-type", "line"),
//	        g.Attr("data-chart-data", opts.DataJSON()),
//	        g.Group(children),
//	    )
//	}
type ComponentConstructor func(props any, children ...g.Node) g.Node

// ComponentPluginBase provides a base implementation for component plugins.
// Embed this in your plugin to inherit default behavior.
//
// Example:
//
//	type MyPlugin struct {
//	    *ComponentPluginBase
//	}
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
// This can be called during plugin initialization.
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

