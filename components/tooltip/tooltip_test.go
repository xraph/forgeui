package tooltip_test

import (
	"bytes"
	"strings"
	"testing"

	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/components/tooltip"
)

// TestTooltip tests the Tooltip component
func TestTooltip(t *testing.T) {
	var buf bytes.Buffer
	node := tooltip.Tooltip(
		tooltip.TooltipProps{
			Position: forgeui.PositionTop,
			Delay:    200,
		},
		html.Button(html.Class("btn")),
		"Helpful tooltip",
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "x-data") {
		t.Error("missing Alpine data")
	}
	if !strings.Contains(html, "Helpful tooltip") {
		t.Error("missing tooltip content")
	}
	if !strings.Contains(html, "mouseenter") {
		t.Error("missing hover handler")
	}
}

// TestTooltipPositions tests all positions
func TestTooltipPositions(t *testing.T) {
	positions := []forgeui.Position{
		forgeui.PositionTop,
		forgeui.PositionRight,
		forgeui.PositionBottom,
		forgeui.PositionLeft,
	}

	for _, pos := range positions {
		t.Run(string(pos), func(t *testing.T) {
			var buf bytes.Buffer
			node := tooltip.Tooltip(
				tooltip.TooltipProps{Position: pos},
				html.Button(),
				"Test",
			)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}
		})
	}
}

