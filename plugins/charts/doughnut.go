package charts

import (
	g "github.com/maragudk/gomponents"
)

// doughnutChart creates a doughnut chart component (hollow pie chart).
func (c *Charts) doughnutChart(props any, children ...g.Node) g.Node {
	data, ok := props.(DoughnutChartData)
	if !ok {
		return g.Text("Invalid data for DoughnutChart")
	}

	opts := DefaultOptions()
	opts.Title = "Doughnut Chart"

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

	return baseChartNode("doughnut", chartData, opts)
}

// DoughnutChart creates a doughnut chart with the given data.
//
// Example:
//
//	charts.DoughnutChart(charts.DoughnutChartData{
//	    Labels: []string{"Development", "Marketing", "Sales", "Operations"},
//	    Data: []float64{35, 25, 20, 20},
//	    BackgroundColor: charts.DefaultColors,
//	    BorderColor: charts.BorderColors,
//	    BorderWidth: 2,
//	})
func DoughnutChart(data DoughnutChartData) g.Node {
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

	return baseChartNode("doughnut", chartData, opts)
}

