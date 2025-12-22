package table

import (
	"bytes"
	"strings"
	"testing"
)

func TestDataTable(t *testing.T) {
	testData := []map[string]any{
		{"name": "John", "email": "john@example.com", "status": "active"},
		{"name": "Jane", "email": "jane@example.com", "status": "inactive"},
	}

	tests := []struct {
		name   string
		opts   []DataTableOption
		wantRe []string
	}{
		{
			name: "renders basic data table",
			opts: []DataTableOption{
				WithColumns(
					Column{Key: "name", Label: "Name"},
					Column{Key: "email", Label: "Email"},
				),
				WithData(testData),
			},
			wantRe: []string{
				`x-data`,
				`Name`,
				`Email`,
				`<template`,
				`x-for`,
			},
		},
		{
			name: "renders with sortable columns",
			opts: []DataTableOption{
				WithColumns(
					Column{Key: "name", Label: "Name", Sortable: true},
				),
				WithData(testData),
			},
			wantRe: []string{
				`sortBy`,
				`name`,
				`sortColumn`,
			},
		},
		{
			name: "renders with filterable columns",
			opts: []DataTableOption{
				WithColumns(
					Column{
						Key:        "status",
						Label:      "Status",
						Filterable: true,
						FilterOptions: []FilterOption{
							{Value: "active", Label: "Active"},
							{Value: "inactive", Label: "Inactive"},
						},
					},
				),
				WithData(testData),
			},
			wantRe: []string{
				`filters`,
				`<select`,
				`Active`,
				`Inactive`,
			},
		},
		{
			name: "renders with pagination",
			opts: []DataTableOption{
				WithColumns(
					Column{Key: "name", Label: "Name"},
				),
				WithData(testData),
				WithPagination(),
				WithPageSize(10),
			},
			wantRe: []string{
				`showPagination`,
				`true`,
				`Previous`,
				`Next`,
				`totalPages`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := DataTable(tt.opts...)
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

func TestDataTableColumn(t *testing.T) {
	tests := []struct {
		name   string
		column Column
		want   string
	}{
		{
			name: "basic column",
			column: Column{
				Key:   "name",
				Label: "Name",
			},
			want: "Name",
		},
		{
			name: "sortable column",
			column: Column{
				Key:      "email",
				Label:    "Email",
				Sortable: true,
			},
			want: "Email",
		},
		{
			name: "filterable column",
			column: Column{
				Key:        "status",
				Label:      "Status",
				Filterable: true,
			},
			want: "Status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.column.Label != tt.want {
				t.Errorf("column.Label = %q, want %q", tt.column.Label, tt.want)
			}
		})
	}
}

func TestHasFilterableColumns(t *testing.T) {
	tests := []struct {
		name    string
		columns []Column
		want    bool
	}{
		{
			name: "has filterable column",
			columns: []Column{
				{Key: "name", Label: "Name"},
				{Key: "status", Label: "Status", Filterable: true},
			},
			want: true,
		},
		{
			name: "no filterable columns",
			columns: []Column{
				{Key: "name", Label: "Name"},
				{Key: "email", Label: "Email"},
			},
			want: false,
		},
		{
			name:    "empty columns",
			columns: []Column{},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasFilterableColumns(tt.columns)
			if got != tt.want {
				t.Errorf("hasFilterableColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderDataTableHeaders(t *testing.T) {
	columns := []Column{
		{Key: "name", Label: "Name", Sortable: true},
		{Key: "email", Label: "Email"},
	}

	headers := renderDataTableHeaders(columns)

	if len(headers) != len(columns) {
		t.Errorf("expected %d headers, got %d", len(columns), len(headers))
	}

	var buf bytes.Buffer
	for _, header := range headers {
		if err := header.Render(&buf); err != nil {
			t.Fatalf("render error: %v", err)
		}
	}

	html := buf.String()

	if !strings.Contains(html, "Name") {
		t.Error("headers should contain 'Name'")
	}
	if !strings.Contains(html, "Email") {
		t.Error("headers should contain 'Email'")
	}
}

func TestRenderDataTableCells(t *testing.T) {
	columns := []Column{
		{Key: "name", Label: "Name"},
		{Key: "email", Label: "Email", Align: AlignRight},
	}

	cells := renderDataTableCells(columns)

	if len(cells) != len(columns) {
		t.Errorf("expected %d cells, got %d", len(columns), len(cells))
	}

	var buf bytes.Buffer
	for _, cell := range cells {
		if err := cell.Render(&buf); err != nil {
			t.Fatalf("render error: %v", err)
		}
	}

	html := buf.String()

	if !strings.Contains(html, "x-text") {
		t.Error("cells should contain x-text directive")
	}
	if !strings.Contains(html, "row.name") {
		t.Error("cells should reference row data")
	}
}

