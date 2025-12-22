package charts

import (
	g "maragu.dev/gomponents"
)

// areaChart creates an area chart component (filled line chart).
func (c *Charts) areaChart(props any, children ...g.Node) g.Node {
	data, ok := props.(AreaChartData)
	if !ok {
		return g.Text("Invalid data for AreaChart")
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
					"fill":            true, // Area charts are always filled
					"tension":         ds.Tension,
				}
			}

			return datasets
		}(),
	}

	return baseChartNode("line", chartData, opts)
}

// AreaChart creates an area chart with the given data.
//
// Example:
//
//	charts.AreaChart(charts.AreaChartData{
//	    Labels: []string{"Week 1", "Week 2", "Week 3", "Week 4"},
//	    Datasets: []charts.DatasetConfig{
//	        {
//	            Label: "Active Users",
//	            Data: []float64{1200, 1900, 1500, 2100},
//	            BorderColor: "rgb(59, 130, 246)",
//	            BackgroundColor: "rgba(59, 130, 246, 0.3)",
//	            BorderWidth: 2,
//	            Tension: 0.4,
//	        },
//	    },
//	})
func AreaChart(data AreaChartData) g.Node {
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
