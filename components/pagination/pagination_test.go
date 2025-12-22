package pagination_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui/components/pagination"
	g "maragu.dev/gomponents"
)

func TestPagination(t *testing.T) {
	tests := []struct {
		name string
		opts []pagination.Option
		want []string
	}{
		{
			name: "renders default pagination",
			opts: []pagination.Option{
				pagination.WithTotalPages(5),
			},
			want: []string{
				`x-data`,
				`currentPage`,
				`totalPages`,
				`role="navigation"`,
				`aria-label="Pagination"`,
			},
		},
		{
			name: "renders with current page",
			opts: []pagination.Option{
				pagination.WithCurrentPage(3),
				pagination.WithTotalPages(10),
			},
			want: []string{
				`currentPage`,
				`3`,
			},
		},
		{
			name: "renders with first/last buttons",
			opts: []pagination.Option{
				pagination.WithTotalPages(5),
				pagination.WithShowFirstLast(true),
			},
			want: []string{
				`First`,
				`Last`,
			},
		},
		{
			name: "renders with prev/next buttons",
			opts: []pagination.Option{
				pagination.WithTotalPages(5),
				pagination.WithShowPrevNext(true),
			},
			want: []string{
				`Previous`,
				`Next`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := pagination.Pagination(tt.opts...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, want := range tt.want {
				if !strings.Contains(html, want) {
					t.Errorf("output missing %q\ngot: %s", want, html)
				}
			}
		})
	}
}

func TestPaginationButtons(t *testing.T) {
	tests := []struct {
		name string
		node g.Node
		want []string
	}{
		{
			name: "first button",
			node: pagination.FirstButton(),
			want: []string{
				`First`,
				`:disabled`,
				`currentPage === 1`,
			},
		},
		{
			name: "last button",
			node: pagination.LastButton(),
			want: []string{
				`Last`,
				`:disabled`,
				`currentPage === totalPages`,
			},
		},
		{
			name: "prev button",
			node: pagination.PrevButton(),
			want: []string{
				`Previous`,
				`:disabled`,
			},
		},
		{
			name: "next button",
			node: pagination.NextButton(),
			want: []string{
				`Next`,
				`:disabled`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tt.node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, want := range tt.want {
				if !strings.Contains(html, want) {
					t.Errorf("output missing %q\ngot: %s", want, html)
				}
			}
		})
	}
}
