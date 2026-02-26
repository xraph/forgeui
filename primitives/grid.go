package primitives

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"

	"github.com/xraph/forgeui"
)

// GridProps defines properties for the Grid component
type GridProps struct {
	Cols       int    // number of columns
	ColsSM     int    // columns at sm breakpoint
	ColsMD     int    // columns at md breakpoint
	ColsLG     int    // columns at lg breakpoint
	ColsXL     int    // columns at xl breakpoint
	Gap        string // gap size
	Class      string
	Children   []templ.Component
	Attributes templ.Attributes
}

// GridOption is a functional option for configuring Grid
type GridOption func(*GridProps)

// GridCols sets the number of columns
func GridCols(cols int) GridOption {
	return func(p *GridProps) { p.Cols = cols }
}

// GridColsSM sets columns at sm breakpoint
func GridColsSM(cols int) GridOption {
	return func(p *GridProps) { p.ColsSM = cols }
}

// GridColsMD sets columns at md breakpoint
func GridColsMD(cols int) GridOption {
	return func(p *GridProps) { p.ColsMD = cols }
}

// GridColsLG sets columns at lg breakpoint
func GridColsLG(cols int) GridOption {
	return func(p *GridProps) { p.ColsLG = cols }
}

// GridColsXL sets columns at xl breakpoint
func GridColsXL(cols int) GridOption {
	return func(p *GridProps) { p.ColsXL = cols }
}

// GridGap sets the gap
func GridGap(gap string) GridOption {
	return func(p *GridProps) { p.Gap = gap }
}

// GridClass adds custom classes
func GridClass(class string) GridOption {
	return func(p *GridProps) { p.Class = class }
}

// GridChildren adds child components
func GridChildren(children ...templ.Component) GridOption {
	return func(p *GridProps) { p.Children = append(p.Children, children...) }
}

// GridAttrs adds custom attributes
func GridAttrs(attrs templ.Attributes) GridOption {
	return func(p *GridProps) {
		if p.Attributes == nil {
			p.Attributes = templ.Attributes{}
		}
		for k, v := range attrs {
			p.Attributes[k] = v
		}
	}
}

// gridClasses computes the CSS classes for a Grid component.
func gridClasses(props *GridProps) string {
	classes := []string{"grid"}

	if props.Cols > 0 {
		classes = append(classes, fmt.Sprintf("grid-cols-%d", props.Cols))
	}
	if props.ColsSM > 0 {
		classes = append(classes, fmt.Sprintf("sm:grid-cols-%d", props.ColsSM))
	}
	if props.ColsMD > 0 {
		classes = append(classes, fmt.Sprintf("md:grid-cols-%d", props.ColsMD))
	}
	if props.ColsLG > 0 {
		classes = append(classes, fmt.Sprintf("lg:grid-cols-%d", props.ColsLG))
	}
	if props.ColsXL > 0 {
		classes = append(classes, fmt.Sprintf("xl:grid-cols-%d", props.ColsXL))
	}

	if props.Gap != "" {
		classes = append(classes, "gap-"+props.Gap)
	}
	if props.Class != "" {
		classes = append(classes, props.Class)
	}

	return forgeui.CN(classes...)
}

// Grid creates a CSS Grid container.
func Grid(opts ...GridOption) templ.Component {
	props := &GridProps{
		Cols: 1,
		Gap:  "4",
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := gridClasses(props)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := writeOpenTag(w, "div", classes, props.Attributes); err != nil {
			return err
		}
		if err := renderChildren(ctx, w, props.Children); err != nil {
			return err
		}
		return writeCloseTag(w, "div")
	})
}
