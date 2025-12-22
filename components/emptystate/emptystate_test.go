package emptystate

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func TestEmptyState(t *testing.T) {
	tests := []struct {
		name   string
		opts   []Option
		wantRe []string
	}{
		{
			name: "renders basic empty state",
			opts: []Option{
				WithTitle("No data"),
			},
			wantRe: []string{
				`<div`,
				`flex flex-col items-center justify-center`,
				`No data`,
			},
		},
		{
			name: "renders with title and description",
			opts: []Option{
				WithTitle("No results found"),
				WithDescription("Try adjusting your filters"),
			},
			wantRe: []string{
				`No results found`,
				`Try adjusting your filters`,
				`text-2xl font-bold`,
				`text-sm text-muted-foreground`,
			},
		},
		{
			name: "renders with icon",
			opts: []Option{
				WithIcon(html.Div(g.Text("Icon"))),
				WithTitle("Empty"),
			},
			wantRe: []string{
				`Icon`,
				`text-muted-foreground/60`,
			},
		},
		{
			name: "renders with action button",
			opts: []Option{
				WithTitle("No items"),
				WithAction(html.Button(g.Text("Add Item"))),
			},
			wantRe: []string{
				`Add Item`,
				`<button`,
			},
		},
		{
			name: "renders all options",
			opts: []Option{
				WithIcon(html.Span(g.Text("📦"))),
				WithTitle("No products"),
				WithDescription("Start by adding your first product"),
				WithAction(html.Button(g.Text("Add Product"))),
			},
			wantRe: []string{
				`📦`,
				`No products`,
				`Start by adding your first product`,
				`Add Product`,
			},
		},
		{
			name: "renders with custom class",
			opts: []Option{
				WithTitle("Empty"),
				WithClass("custom-empty"),
			},
			wantRe: []string{
				`custom-empty`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := EmptyState(tt.opts...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			htmlStr := buf.String()

			for _, pattern := range tt.wantRe {
				if !strings.Contains(htmlStr, pattern) {
					t.Errorf("output does not contain pattern %q\ngot: %s", pattern, htmlStr)
				}
			}
		})
	}
}

func TestEmptyStateWithoutOptions(t *testing.T) {
	var buf bytes.Buffer
	node := EmptyState()
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Should still render the container
	if !strings.Contains(html, "flex flex-col items-center justify-center") {
		t.Error("should render container classes")
	}
}

func TestEmptyStateStructure(t *testing.T) {
	var buf bytes.Buffer
	node := EmptyState(
		WithIcon(html.Div(g.Text("Icon"))),
		WithTitle("Title"),
		WithDescription("Description"),
		WithAction(html.Button(g.Text("Action"))),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Verify structure elements exist in order
	iconIndex := strings.Index(html, "Icon")
	titleIndex := strings.Index(html, "Title")
	descIndex := strings.Index(html, "Description")
	actionIndex := strings.Index(html, "Action")

	if iconIndex == -1 || titleIndex == -1 || descIndex == -1 || actionIndex == -1 {
		t.Error("all elements should be present")
	}

	// Verify order: icon -> title -> description -> action
	if !(iconIndex < titleIndex && titleIndex < descIndex && descIndex < actionIndex) {
		t.Error("elements should be in correct order")
	}
}

