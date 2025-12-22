package table

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestTable(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		children []g.Node
		wantRe   []string
		wantNot  []string
	}{
		{
			name: "renders basic table",
			children: []g.Node{
				TableHeader()(
					TableRow()(
						TableHeaderCell()(g.Text("Name")),
					),
				),
			},
			wantRe: []string{
				`data-slot="table-container"`,
				`overflow-x-auto`,
				`data-slot="table"`,
				`w-full caption-bottom text-sm`,
			},
		},
		{
			name: "renders with custom class",
			opts: []Option{WithClass("custom-table")},
			children: []g.Node{
				TableBody()(),
			},
			wantRe: []string{
				`custom-table`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := Table(tt.opts...)(tt.children...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, pattern := range tt.wantRe {
				if !strings.Contains(html, pattern) && !matchesPattern(pattern, html) {
					t.Errorf("output does not contain pattern %q\ngot: %s", pattern, html)
				}
			}

			for _, notWant := range tt.wantNot {
				if strings.Contains(html, notWant) {
					t.Errorf("output should not contain %q\ngot: %s", notWant, html)
				}
			}
		})
	}
}

func TestTableHeader(t *testing.T) {
	tests := []struct {
		name     string
		opts     []HeaderOption
		children []g.Node
		wantRe   []string
	}{
		{
			name: "renders table header",
			children: []g.Node{
				TableRow()(
					TableHeaderCell()(g.Text("Column")),
				),
			},
			wantRe: []string{
				`<thead`,
				`data-slot="table-header"`,
				`[&amp;_tr]:border-b border-border`,
			},
		},
		{
			name: "renders sticky header",
			opts: []HeaderOption{StickyHeader()},
			children: []g.Node{
				TableRow()(),
			},
			wantRe: []string{
				`sticky top-0 z-10`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := TableHeader(tt.opts...)(tt.children...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, pattern := range tt.wantRe {
				if !strings.Contains(html, pattern) {
					t.Errorf("output does not contain pattern %q\ngot: %s", pattern, html)
				}
			}
		})
	}
}

func TestTableBody(t *testing.T) {
	var buf bytes.Buffer
	node := TableBody()(
		TableRow()(
			TableCell()(g.Text("Data")),
		),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "<tbody") {
		t.Error("output should contain <tbody")
	}
	if !strings.Contains(html, "tr:last-child") {
		t.Error("output should contain last-child border style")
	}
}

func TestTableRow(t *testing.T) {
	tests := []struct {
		name     string
		opts     []RowOption
		children []g.Node
		wantRe   []string
	}{
		{
			name: "renders basic row",
			children: []g.Node{
				TableCell()(g.Text("Cell")),
			},
			wantRe: []string{
				`<tr`,
				`border-border border-b transition-colors hover:bg-muted/50`,
			},
		},
		{
			name: "renders clickable row",
			opts: []RowOption{ClickableRow()},
			children: []g.Node{
				TableCell()(g.Text("Cell")),
			},
			wantRe: []string{
				`cursor-pointer`,
			},
		},
		{
			name: "renders with onclick handler",
			opts: []RowOption{WithOnClick("handleClick()")},
			children: []g.Node{
				TableCell()(g.Text("Cell")),
			},
			wantRe: []string{
				`onclick="handleClick()"`,
				`cursor-pointer`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := TableRow(tt.opts...)(tt.children...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, pattern := range tt.wantRe {
				if !strings.Contains(html, pattern) {
					t.Errorf("output does not contain pattern %q\ngot: %s", pattern, html)
				}
			}
		})
	}
}

func TestTableCell(t *testing.T) {
	tests := []struct {
		name     string
		opts     []CellOption
		children []g.Node
		wantRe   []string
	}{
		{
			name: "renders basic cell",
			children: []g.Node{
				g.Text("Cell content"),
			},
			wantRe: []string{
				`<td`,
				`data-slot="table-cell"`,
				`p-2 align-middle whitespace-nowrap`,
				`Cell content`,
			},
		},
		{
			name: "renders with center alignment",
			opts: []CellOption{WithAlign(AlignCenter)},
			children: []g.Node{
				g.Text("Centered"),
			},
			wantRe: []string{
				`text-center`,
			},
		},
		{
			name: "renders with right alignment",
			opts: []CellOption{WithAlign(AlignRight)},
			children: []g.Node{
				g.Text("Right"),
			},
			wantRe: []string{
				`text-right`,
			},
		},
		{
			name: "renders with width",
			opts: []CellOption{WithWidth("200px")},
			children: []g.Node{
				g.Text("Fixed width"),
			},
			wantRe: []string{
				`style="width: 200px"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := TableCell(tt.opts...)(tt.children...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, pattern := range tt.wantRe {
				if !strings.Contains(html, pattern) {
					t.Errorf("output does not contain pattern %q\ngot: %s", pattern, html)
				}
			}
		})
	}
}

func TestTableHeaderCell(t *testing.T) {
	tests := []struct {
		name     string
		opts     []CellOption
		children []g.Node
		wantRe   []string
	}{
		{
			name: "renders header cell",
			children: []g.Node{
				g.Text("Header"),
			},
			wantRe: []string{
				`<th`,
				`data-slot="table-head"`,
				`h-10 px-2`,
				`text-foreground`,
				`whitespace-nowrap`,
				`Header`,
			},
		},
		{
			name: "renders with alignment",
			opts: []CellOption{WithAlign(AlignCenter)},
			children: []g.Node{
				g.Text("Centered Header"),
			},
			wantRe: []string{
				`text-center`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := TableHeaderCell(tt.opts...)(tt.children...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, pattern := range tt.wantRe {
				if !strings.Contains(html, pattern) {
					t.Errorf("output does not contain pattern %q\ngot: %s", pattern, html)
				}
			}
		})
	}
}

// matchesPattern is a simple pattern matcher (for basic regex-like patterns)
func matchesPattern(pattern, text string) bool {
	// For now, just use strings.Contains
	// In a real implementation, you might use regexp
	return strings.Contains(text, pattern)
}
