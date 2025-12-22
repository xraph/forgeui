package charts

import (
	g "github.com/maragudk/gomponents"
)

// pieChart creates a pie chart component.
func (c *Charts) pieChart(props any, children ...g.Node) g.Node {
	data, ok := props.(PieChartData)
	if !ok {
		return g.Text("Invalid data for PieChart")
	}

	opts := DefaultOptions()
	opts.Title = "Pie Chart"

	chartData := map[string]any{
		"labels": data.Labels,
		"datasets": []map[string]any{
			{
				"data":            data.Data,
				"backgroundColor": data.BackgroundColor,
				"borderColor":     data.BorderColor,
				"borderWidth":     data.BorderWidth,
			},
		},
	}

	return baseChartNode("pie", chartData, opts)
}

// PieChart creates a pie chart with the given data.
//
// Example:
//
//	charts.PieChart(charts.PieChartData{
//	    Labels: []string{"Chrome", "Firefox", "Safari", "Edge"},
//	    Data: []float64{55.2, 18.7, 15.3, 10.8},
//	    BackgroundColor: charts.DefaultColors,
//	    BorderColor: charts.BorderColors,
//	    BorderWidth: 1,
//	})
func PieChart(data PieChartData) g.Node {
	opts := DefaultOptions()
	
	chartData := map[string]any{
		"labels": data.Labels,
		"datasets": []map[string]any{
			{
				"data":            data.Data,
				"backgroundColor": data.BackgroundColor,
				"borderColor":     data.BorderColor,
				"borderWidth":     data.BorderWidth,
			},
		},
	}

	return baseChartNode("pie", chartData, opts)
}

