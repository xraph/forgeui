package primitives

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
)

func TestIntegration_ComplexLayout(t *testing.T) {
	// Test a complex layout combining multiple primitives
	layout := Container(
		VStack("8",
			Text(
				TextAs("h1"),
				TextSize("text-4xl"),
				TextWeight("font-bold"),
				TextChildren(g.Text("Welcome")),
			),
			Flex(
				FlexJustify("between"),
				FlexAlign("center"),
				FlexChildren(
					Box(
						WithPadding("p-4"),
						WithBackground("bg-blue-500"),
						WithRounded("rounded-lg"),
						WithChildren(g.Text("Box 1")),
					),
					Spacer(),
					Box(
						WithPadding("p-4"),
						WithBackground("bg-green-500"),
						WithRounded("rounded-lg"),
						WithChildren(g.Text("Box 2")),
					),
				),
			),
			Grid(
				GridCols(2),
				GridColsMD(4),
				GridGap("4"),
				GridChildren(
					Box(WithChildren(g.Text("Item 1"))),
					Box(WithChildren(g.Text("Item 2"))),
					Box(WithChildren(g.Text("Item 3"))),
					Box(WithChildren(g.Text("Item 4"))),
				),
			),
		),
	)

	var buf bytes.Buffer
	if err := layout.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()

	// Verify key structural elements
	expected := []string{
		"container",
		"flex-col",
		"<h1",
		"text-4xl",
		"justify-between",
		"flex-1",
		"grid",
		"grid-cols-2",
		"md:grid-cols-4",
	}

	for _, exp := range expected {
		if !strings.Contains(html, exp) {
			t.Errorf("expected %v in complex layout", exp)
		}
	}
}

func TestBox_WithAllOptions(t *testing.T) {
	box := Box(
		WithAs("article"),
		WithClass("custom-class"),
		WithPadding("p-8"),
		WithMargin("m-4"),
		WithBackground("bg-gray-100"),
		WithRounded("rounded-xl"),
		WithShadow("shadow-lg"),
		WithWidth("w-full"),
		WithHeight("h-64"),
		WithAttrs(g.Attr("data-test", "value")),
		WithChildren(g.Text("Content")),
	)

	var buf bytes.Buffer
	box.Render(&buf)
	html := buf.String()

	expected := []string{
		"<article",
		"custom-class",
		"p-8",
		"m-4",
		"bg-gray-100",
		"rounded-xl",
		"shadow-lg",
		"w-full",
		"h-64",
		`data-test="value"`,
		"Content",
	}

	for _, exp := range expected {
		if !strings.Contains(html, exp) {
			t.Errorf("expected %v in box with all options", exp)
		}
	}
}

func TestFlex_AllVariants(t *testing.T) {
	tests := []struct {
		name string
		opts []FlexOption
		want []string
	}{
		{
			name: "wrapped flex with all options",
			opts: []FlexOption{
				FlexDirection("col"),
				FlexWrap("wrap"),
				FlexJustify("evenly"),
				FlexAlign("baseline"),
				FlexGap("6"),
				FlexClass("custom"),
			},
			want: []string{"flex-col", "flex-wrap", "justify-evenly", "items-baseline", "gap-6", "custom"},
		},
		{
			name: "nowrap flex",
			opts: []FlexOption{
				FlexWrap("nowrap"),
			},
			want: []string{"flex-nowrap"},
		},
		{
			name: "wrap-reverse flex",
			opts: []FlexOption{
				FlexWrap("wrap-reverse"),
			},
			want: []string{"flex-wrap-reverse"},
		},
		{
			name: "various justify options",
			opts: []FlexOption{
				FlexJustify("around"),
			},
			want: []string{"justify-around"},
		},
		{
			name: "stretch align",
			opts: []FlexOption{
				FlexAlign("stretch"),
			},
			want: []string{"items-stretch"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flex := Flex(tt.opts...)

			var buf bytes.Buffer
			flex.Render(&buf)
			html := buf.String()

			for _, want := range tt.want {
				if !strings.Contains(html, want) {
					t.Errorf("expected %v in flex variant %s", want, tt.name)
				}
			}
		})
	}
}

func TestGrid_ResponsiveLayout(t *testing.T) {
	grid := Grid(
		GridCols(1),
		GridColsSM(2),
		GridColsMD(3),
		GridColsLG(4),
		GridColsXL(6),
		GridGap("8"),
		GridClass("custom-grid"),
		GridAttrs(g.Attr("data-grid", "responsive")),
		GridChildren(
			Box(WithChildren(g.Text("Item"))),
		),
	)

	var buf bytes.Buffer
	grid.Render(&buf)
	html := buf.String()

	expected := []string{
		"grid-cols-1",
		"sm:grid-cols-2",
		"md:grid-cols-3",
		"lg:grid-cols-4",
		"xl:grid-cols-6",
		"gap-8",
		"custom-grid",
		`data-grid="responsive"`,
	}

	for _, exp := range expected {
		if !strings.Contains(html, exp) {
			t.Errorf("expected %v in responsive grid", exp)
		}
	}
}

func TestText_AllOptions(t *testing.T) {
	text := Text(
		TextAs("span"),
		TextSize("text-xl"),
		TextWeight("font-semibold"),
		TextColor("text-red-500"),
		TextAlign("text-right"),
		TextClass("custom-text"),
		TextAttrs(g.Attr("data-text", "value")),
		TextChildren(g.Text("Hello World")),
	)

	var buf bytes.Buffer
	text.Render(&buf)
	html := buf.String()

	expected := []string{
		"<span",
		"text-xl",
		"font-semibold",
		"text-red-500",
		"text-right",
		"custom-text",
		`data-text="value"`,
		"Hello World",
	}

	for _, exp := range expected {
		if !strings.Contains(html, exp) {
			t.Errorf("expected %v in text with all options", exp)
		}
	}
}

func TestStack_WithMultipleChildren(t *testing.T) {
	items := []g.Node{
		g.Text("Item 1"),
		g.Text("Item 2"),
		g.Text("Item 3"),
	}

	vstack := VStack("4", items...)

	var buf bytes.Buffer
	vstack.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "Item 1") || !strings.Contains(html, "Item 2") || !strings.Contains(html, "Item 3") {
		t.Error("expected all items in vstack")
	}
}

func TestCenter_EmptyChildren(t *testing.T) {
	center := Center()

	var buf bytes.Buffer
	if err := center.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "justify-center") {
		t.Error("expected justify-center even with no children")
	}
}
