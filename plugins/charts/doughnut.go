package charts

import (
	"github.com/a-h/templ"
)

// doughnutChart creates a doughnut chart component (hollow pie chart).
func (c *Charts) doughnutChart(props any, children ...templ.Component) templ.Component {
	data, ok := props.(DoughnutChartData)
	if !ok {
		return textComponent("Invalid data for DoughnutChart")
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
func DoughnutChart(data DoughnutChartData) templ.Component {
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
