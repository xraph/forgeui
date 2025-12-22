package charts

import (
	g "maragu.dev/gomponents"
)

// lineChart creates a line chart component.
func (c *Charts) lineChart(props any, children ...g.Node) g.Node {
	data, ok := props.(LineChartData)
	if !ok {
		return g.Text("Invalid data for LineChart")
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
//
// Example:
//
//	charts.LineChart(charts.LineChartData{
//	    Labels: []string{"Jan", "Feb", "Mar", "Apr"},
//	    Datasets: []charts.DatasetConfig{
//	        {
//	            Label: "Sales",
//	            Data: []float64{12, 19, 3, 5},
//	            BorderColor: "rgb(59, 130, 246)",
//	            BackgroundColor: "rgba(59, 130, 246, 0.2)",
//	            Tension: 0.4,
//	        },
//	    },
//	})
func LineChart(data LineChartData) g.Node {
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

