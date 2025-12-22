package charts

// ChartOptions holds common configuration for all chart types.
type ChartOptions struct {
	// Title of the chart
	Title string

	// Width of the canvas (e.g., "400px", "100%")
	Width string

	// Height of the canvas (e.g., "300px")
	Height string

	// ShowLegend displays the legend
	ShowLegend bool

	// LegendPosition where to place the legend ("top", "bottom", "left", "right")
	LegendPosition string

	// ShowTooltip enables tooltips on hover
	ShowTooltip bool

	// MaintainAspectRatio maintains aspect ratio when resizing
	MaintainAspectRatio bool

	// Animated enables chart animations
	Animated bool
}

// DefaultOptions returns default chart options.
func DefaultOptions() ChartOptions {
	return ChartOptions{
		ShowLegend:          true,
		LegendPosition:      "top",
		ShowTooltip:         true,
		MaintainAspectRatio: true,
		Animated:            true,
	}
}

// DatasetConfig holds configuration for a dataset.
type DatasetConfig struct {
	Label           string
	Data            []float64
	BackgroundColor any // string or []string
	BorderColor     any // string or []string
	BorderWidth     int
	Fill            bool
	Tension         float64 // For line charts (curve smoothness)
}

// LineChartData holds data for line charts.
type LineChartData struct {
	Labels   []string
	Datasets []DatasetConfig
}

// BarChartData holds data for bar charts.
type BarChartData struct {
	Labels   []string
	Datasets []DatasetConfig
}

// PieChartData holds data for pie charts.
type PieChartData struct {
	Labels          []string
	Data            []float64
	BackgroundColor []string
	BorderColor     []string
	BorderWidth     int
}

// AreaChartData holds data for area charts.
type AreaChartData struct {
	Labels   []string
	Datasets []DatasetConfig
}

// DoughnutChartData holds data for doughnut charts.
type DoughnutChartData struct {
	Labels          []string
	Data            []float64
	BackgroundColor []string
	BorderColor     []string
	BorderWidth     int
}

// DefaultColors provides a default color palette for charts.
var DefaultColors = []string{
	"rgba(59, 130, 246, 0.8)",  // blue
	"rgba(34, 197, 94, 0.8)",   // green
	"rgba(234, 179, 8, 0.8)",   // yellow
	"rgba(239, 68, 68, 0.8)",   // red
	"rgba(168, 85, 247, 0.8)",  // purple
	"rgba(236, 72, 153, 0.8)",  // pink
	"rgba(14, 165, 233, 0.8)",  // sky
	"rgba(251, 146, 60, 0.8)",  // orange
}

// BorderColors provides border colors matching DefaultColors.
var BorderColors = []string{
	"rgb(59, 130, 246)",  // blue
	"rgb(34, 197, 94)",   // green
	"rgb(234, 179, 8)",   // yellow
	"rgb(239, 68, 68)",   // red
	"rgb(168, 85, 247)",  // purple
	"rgb(236, 72, 153)",  // pink
	"rgb(14, 165, 233)",  // sky
	"rgb(251, 146, 60)",  // orange
}

