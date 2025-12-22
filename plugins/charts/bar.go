package charts

import (
	g "maragu.dev/gomponents"
)

// barChart creates a bar chart component.
func (c *Charts) barChart(props any, children ...g.Node) g.Node {
	data, ok := props.(BarChartData)
	if !ok {
		return g.Text("Invalid data for BarChart")
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
//
// Example:
//
//	charts.BarChart(charts.BarChartData{
//	    Labels: []string{"Q1", "Q2", "Q3", "Q4"},
//	    Datasets: []charts.DatasetConfig{
//	        {
//	            Label: "Revenue",
//	            Data: []float64{45000, 52000, 48000, 61000},
//	            BackgroundColor: charts.DefaultColors[0],
//	            BorderColor: charts.BorderColors[0],
//	            BorderWidth: 1,
//	        },
//	    },
//	})
func BarChart(data BarChartData) g.Node {
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

