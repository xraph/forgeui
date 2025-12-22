package popover_test

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/popover"
)

// TestPopover tests the Popover component
func TestPopover(t *testing.T) {
	positions := []forgeui.Position{
		forgeui.PositionTop,
		forgeui.PositionRight,
		forgeui.PositionBottom,
		forgeui.PositionLeft,
	}

	for _, pos := range positions {
		t.Run(string(pos), func(t *testing.T) {
			var buf bytes.Buffer
			node := popover.Popover(
				popover.PopoverProps{Position: pos},
				button.Button(g.Text("Open")),
				html.Div(g.Text("Popover content")),
			)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()
			if !strings.Contains(html, "x-data") {
				t.Error("missing Alpine data")
			}
			if !strings.Contains(html, "x-show") {
				t.Error("missing x-show directive")
			}
		})
	}
}

// TestPopoverAlignment tests different alignments
func TestPopoverAlignment(t *testing.T) {
	alignments := []forgeui.Align{
		forgeui.AlignStart,
		forgeui.AlignCenter,
		forgeui.AlignEnd,
	}

	for _, align := range alignments {
		t.Run(string(align), func(t *testing.T) {
			var buf bytes.Buffer
			node := popover.Popover(
				popover.PopoverProps{
					Position: forgeui.PositionBottom,
					Align:    align,
				},
				button.Button(g.Text("Open")),
				html.Div(g.Text("Content")),
			)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}
		})
	}
}

