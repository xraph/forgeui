// Package charts provides data visualization components using Chart.js.
package charts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/a-h/templ"

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
func chartDataJSON(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}

	return string(jsonData)
}

// baseChartNode creates the base canvas element for charts.
func baseChartNode(chartType string, data any, options ChartOptions) templ.Component {
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

	configJSON := chartDataJSON(config)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<canvas x-data="{ chart: null }" x-init="`); err != nil {
			return err
		}

		initCode := fmt.Sprintf(`chart = new Chart($el, %s); $watch('data', value => { if (chart && value) { chart.data = value; chart.update(); } });`, configJSON)
		if _, err := io.WriteString(w, initCode); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `"`); err != nil {
			return err
		}

		if options.Width != "" {
			if _, err := fmt.Fprintf(w, ` width="%s"`, options.Width); err != nil {
				return err
			}
		}

		if options.Height != "" {
			if _, err := fmt.Fprintf(w, ` height="%s"`, options.Height); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, `></canvas>`)
		return err
	})
}

// textComponent creates a simple text component for error messages.
func textComponent(text string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	})
}
