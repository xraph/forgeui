package charts

import (
	"github.com/a-h/templ"
)

// areaChart creates an area chart component (filled line chart).
func (c *Charts) areaChart(props any, children ...templ.Component) templ.Component {
	data, ok := props.(AreaChartData)
	if !ok {
		return textComponent("Invalid data for AreaChart")
	}

	opts := DefaultOptions()
	opts.Title = "Area Chart"

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
					"fill":            true,
					"tension":         ds.Tension,
				}
			}

			return datasets
		}(),
	}

	return baseChartNode("line", chartData, opts)
}

// AreaChart creates an area chart with the given data.
func AreaChart(data AreaChartData) templ.Component {
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
					"fill":            true,
					"tension":         ds.Tension,
				}
			}

			return datasets
		}(),
	}

	return baseChartNode("line", chartData, opts)
}
