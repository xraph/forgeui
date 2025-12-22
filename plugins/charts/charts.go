// Package charts provides data visualization components using Chart.js.
//
// The charts plugin provides beautiful, responsive charts that integrate
// seamlessly with ForgeUI's theme system.
//
// # Basic Usage
//
//	registry := plugin.NewRegistry()
//	registry.Use(charts.New())
//
// # Available Charts
//
//   - LineChart - Time series and trends
//   - BarChart - Categorical comparisons
//   - PieChart - Proportional data
//   - AreaChart - Filled line charts
//   - DoughnutChart - Hollow pie charts
//
// # Features
//
//   - Responsive sizing
//   - Theme integration (matches ForgeUI colors)
//   - Animation configuration
//   - Legend positioning
//   - Tooltip formatting
//   - Data point styling
package charts

import (
	"context"
	"encoding/json"

	g "github.com/maragudk/gomponents"

	"github.com/xraph/forgeui/plugin"
)

// Charts plugin implements Component and Alpine plugins.
type Charts struct {
	*plugin.ComponentPluginBase
	version string
}

// New creates a new Charts plugin.
func New() *Charts {
	c := &Charts{
		version: "4.4.1",
	}

	c.ComponentPluginBase = plugin.NewComponentPluginBase(
		plugin.PluginInfo{
			Name:        "charts",
			Version:     "1.0.0",
			Description: "Data visualization with Chart.js",
			Author:      "ForgeUI",
			License:     "MIT",
		},
		map[string]plugin.ComponentConstructor{
			"LineChart":     c.lineChart,
			"BarChart":      c.barChart,
			"PieChart":      c.pieChart,
			"AreaChart":     c.areaChart,
			"DoughnutChart": c.doughnutChart,
		},
	)

	return c
}

// Init initializes the charts plugin.
func (c *Charts) Init(ctx context.Context, registry *plugin.Registry) error {
	return c.ComponentPluginBase.Init(ctx, registry)
}

// Scripts returns Chart.js library.
func (c *Charts) Scripts() []plugin.Script {
	return []plugin.Script{
		{
			Name:     "chartjs",
			URL:      "https://cdn.jsdelivr.net/npm/chart.js@" + c.version,
			Priority: 100,
			Defer:    false,
		},
	}
}

// Directives returns custom Alpine directives.
func (c *Charts) Directives() []plugin.AlpineDirective {
	return nil
}

// Stores returns Alpine stores.
func (c *Charts) Stores() []plugin.AlpineStore {
	return nil
}

// Magics returns custom magic properties.
func (c *Charts) Magics() []plugin.AlpineMagic {
	return nil
}

// AlpineComponents returns Alpine.data components.
func (c *Charts) AlpineComponents() []plugin.AlpineComponent {
	return nil
}

// chartData converts data to JSON string.
func chartData(data any) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

// baseChartNode creates the base canvas element for charts.
func baseChartNode(chartType string, data any, options ChartOptions) g.Node {
	config := map[string]any{
		"type": chartType,
		"data": data,
		"options": map[string]any{
			"responsive":          true,
			"maintainAspectRatio": options.MaintainAspectRatio,
			"plugins": map[string]any{
				"legend": map[string]any{
					"display":  options.ShowLegend,
					"position": options.LegendPosition,
				},
				"tooltip": map[string]any{
					"enabled": options.ShowTooltip,
				},
			},
		},
	}

	if options.Title != "" {
		config["options"].(map[string]any)["plugins"].(map[string]any)["title"] = map[string]any{
			"display": true,
			"text":    options.Title,
		}
	}

	configJSON := chartData(config)

	return g.Group([]g.Node{
		g.El("canvas",
			g.Attr("x-data", "{ chart: null }"),
			g.Attr("x-init", `
				chart = new Chart($el, `+configJSON+`);
				$watch('data', value => {
					if (chart && value) {
						chart.data = value;
						chart.update();
					}
				});
			`),
			g.If(options.Width != "", g.Attr("width", options.Width)),
			g.If(options.Height != "", g.Attr("height", options.Height)),
		),
	})
}

