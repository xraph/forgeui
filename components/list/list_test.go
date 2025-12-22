package list

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
)

func TestList(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		children []g.Node
		wantRe   []string
	}{
		{
			name: "renders basic list",
			children: []g.Node{
				ListItem()(g.Text("Item 1")),
				ListItem()(g.Text("Item 2")),
			},
			wantRe: []string{
				`<ul`,
				`list-disc pl-5`,
				`Item 1`,
				`Item 2`,
			},
		},
		{
			name: "renders list without markers",
			opts: []Option{None()},
			children: []g.Node{
				ListItem()(g.Text("Item")),
			},
			wantRe: []string{
				`list-none`,
			},
		},
		{
			name: "renders list with icons variant",
			opts: []Option{Icons()},
			children: []g.Node{
				ListItem()(g.Text("Icon item")),
			},
			wantRe: []string{
				`list-none space-y-2`,
			},
		},
		{
			name: "renders spaced list",
			opts: []Option{Spaced()},
			children: []g.Node{
				ListItem()(g.Text("Spaced item")),
			},
			wantRe: []string{
				`space-y-2`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := List(tt.opts...)(tt.children...)
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

func TestOrderedList(t *testing.T) {
	var buf bytes.Buffer
	node := OrderedList()(
		ListItem()(g.Text("First")),
		ListItem()(g.Text("Second")),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "<ol") {
		t.Error("output should contain <ol")
	}
	if !strings.Contains(html, "list-decimal") {
		t.Error("output should contain list-decimal class")
	}
	if !strings.Contains(html, "First") {
		t.Error("output should contain 'First'")
	}
}

func TestListItem(t *testing.T) {
	tests := []struct {
		name     string
		opts     []ItemOption
		children []g.Node
		wantRe   []string
	}{
		{
			name: "renders simple list item",
			children: []g.Node{
				g.Text("Simple item"),
			},
			wantRe: []string{
				`<li`,
				`flex items-center gap-2 py-1`,
				`Simple item`,
			},
		},
		{
			name: "renders detailed list item",
			opts: []ItemOption{Detailed()},
			children: []g.Node{
				g.Text("Detailed item"),
			},
			wantRe: []string{
				`flex flex-col gap-1 py-2`,
			},
		},
		{
			name: "renders card style list item",
			opts: []ItemOption{CardStyle()},
			children: []g.Node{
				g.Text("Card item"),
			},
			wantRe: []string{
				`flex items-center gap-3 p-4 rounded-md border`,
			},
		},
		{
			name: "renders clickable list item",
			opts: []ItemOption{Clickable()},
			children: []g.Node{
				g.Text("Clickable"),
			},
			wantRe: []string{
				`cursor-pointer hover:bg-accent/50`,
			},
		},
		{
			name: "renders with onclick handler",
			opts: []ItemOption{WithItemOnClick("handleClick()")},
			children: []g.Node{
				g.Text("With handler"),
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
			node := ListItem(tt.opts...)(tt.children...)
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

func TestListVariants(t *testing.T) {
	variants := []struct {
		variant ListVariant
		want    string
	}{
		{VariantBullets, "list-disc"},
		{VariantNone, "list-none"},
		{VariantIcons, "list-none"},
	}

	for _, v := range variants {
		t.Run(string(v.variant), func(t *testing.T) {
			var buf bytes.Buffer
			node := List(WithVariant(v.variant))(
				ListItem()(g.Text("Test")),
			)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()
			if !strings.Contains(html, v.want) {
				t.Errorf("expected class %q in output\ngot: %s", v.want, html)
			}
		})
	}
}

func TestListItemVariants(t *testing.T) {
	variants := []struct {
		variant ItemVariant
		want    string
	}{
		{ItemSimple, "flex items-center gap-2 py-1"},
		{ItemDetailed, "flex flex-col gap-1 py-2"},
		{ItemCard, "flex items-center gap-3 p-4 rounded-md border"},
	}

	for _, v := range variants {
		t.Run(string(v.variant), func(t *testing.T) {
			var buf bytes.Buffer
			node := ListItem(WithItemVariant(v.variant))(
				g.Text("Test"),
			)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()
			if !strings.Contains(html, v.want) {
				t.Errorf("expected class %q in output\ngot: %s", v.want, html)
			}
		})
	}
}

