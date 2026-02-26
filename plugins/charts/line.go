package charts

import (
	"github.com/a-h/templ"
)

// lineChart creates a line chart component.
func (c *Charts) lineChart(props any, children ...templ.Component) templ.Component {
	data, ok := props.(LineChartData)
	if !ok {
		return textComponent("Invalid data for LineChart")
	}

	opts := DefaultOptions()
	opts.Title = "Line Chart"

	chartData := map[string]any{
		"labels": data.Labels,
		"datasets": func() []map[string]any {
			datasets := make([]map[string]any, len(data.Datasets))
			for i, ds := range data.Datasets {
				datasets[i] = map[string]any{
					"label":           ds.Label,
					"data":            ds.Data,
					"borderColor":     ds.BorderColor,
					"backgroundColor": ds.BackgroundColor,
					"borderWidth":     ds.BorderWidth,
					"fill":            ds.Fill,
					"tension":         ds.Tension,
				}
			}

			return datasets
		}(),
	}

	return baseChartNode("line", chartData, opts)
}

// LineChart creates a line chart with the given data.
func LineChart(data LineChartData) templ.Component {
	opts := DefaultOptions()

	chartData := map[string]any{
		"labels": data.Labels,
		"datasets": func() []map[string]any {
			datasets := make([]map[string]any, len(data.Datasets))
			for i, ds := range data.Datasets {
				datasets[i] = map[string]any{
					"label":           ds.Label,
					"data":            ds.Data,
					"borderColor":     ds.BorderColor,
					"backgroundColor": ds.BackgroundColor,
					"borderWidth":     ds.BorderWidth,
					"fill":            ds.Fill,
					"tension":         ds.Tension,
				}
			}

			return datasets
		}(),
	}

	return baseChartNode("line", chartData, opts)
}
