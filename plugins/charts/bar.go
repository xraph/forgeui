package charts

import (
	"github.com/a-h/templ"
)

// barChart creates a bar chart component.
func (c *Charts) barChart(props any, children ...templ.Component) templ.Component {
	data, ok := props.(BarChartData)
	if !ok {
		return textComponent("Invalid data for BarChart")
	}

	opts := DefaultOptions()
	opts.Title = "Bar Chart"

	chartData := map[string]any{
		"labels": data.Labels,
		"datasets": func() []map[string]any {
			datasets := make([]map[string]any, len(data.Datasets))
			for i, ds := range data.Datasets {
				datasets[i] = map[string]any{
					"label":           ds.Label,
					"data":            ds.Data,
					"backgroundColor": ds.BackgroundColor,
					"borderColor":     ds.BorderColor,
					"borderWidth":     ds.BorderWidth,
				}
			}

			return datasets
		}(),
	}

	return baseChartNode("bar", chartData, opts)
}

// BarChart creates a bar chart with the given data.
func BarChart(data BarChartData) templ.Component {
	opts := DefaultOptions()

	chartData := map[string]any{
		"labels": data.Labels,
		"datasets": func() []map[string]any {
			datasets := make([]map[string]any, len(data.Datasets))
			for i, ds := range data.Datasets {
				datasets[i] = map[string]any{
					"label":           ds.Label,
					"data":            ds.Data,
					"backgroundColor": ds.BackgroundColor,
					"borderColor":     ds.BorderColor,
					"borderWidth":     ds.BorderWidth,
				}
			}

			return datasets
		}(),
	}

	return baseChartNode("bar", chartData, opts)
}
