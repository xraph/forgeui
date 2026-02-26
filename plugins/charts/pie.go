package charts

import (
	"github.com/a-h/templ"
)

// pieChart creates a pie chart component.
func (c *Charts) pieChart(props any, children ...templ.Component) templ.Component {
	data, ok := props.(PieChartData)
	if !ok {
		return textComponent("Invalid data for PieChart")
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
func PieChart(data PieChartData) templ.Component {
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
