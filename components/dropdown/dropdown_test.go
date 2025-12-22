package dropdown_test

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/dropdown"
)

// TestDropdown tests the basic Dropdown component
func TestDropdown(t *testing.T) {
	var buf bytes.Buffer

	node := dropdown.Dropdown(
		dropdown.DropdownProps{
			Position: forgeui.PositionBottom,
			Align:    forgeui.AlignStart,
		},
		button.Button(g.Text("Menu")),
		html.Div(g.Text("Menu content")),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "x-data") {
		t.Error("missing Alpine data")
	}

	if !strings.Contains(html, "role=\"menu\"") {
		t.Error("missing menu role")
	}
}

// TestDropdownMenu tests the DropdownMenu component
func TestDropdownMenu(t *testing.T) {
	var buf bytes.Buffer

	node := dropdown.DropdownMenu(
		dropdown.DropdownMenuTrigger(button.Button(g.Text("Options"))),
		dropdown.DropdownMenuContent(
			dropdown.DropdownMenuItem(g.Text("Profile")),
			dropdown.DropdownMenuSeparator(),
			dropdown.DropdownMenuItem(g.Text("Settings")),
		),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "x-data") {
		t.Error("missing Alpine data")
	}

	if !strings.Contains(html, "menuitem") {
		t.Error("missing menu items")
	}
}

// TestDropdownMenuCheckbox tests checkbox menu items
func TestDropdownMenuCheckbox(t *testing.T) {
	var buf bytes.Buffer

	node := dropdown.DropdownMenuCheckboxItem("check1", "Enable feature", true)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "menuitemcheckbox") {
		t.Error("missing checkbox role")
	}

	if !strings.Contains(html, "Enable feature") {
		t.Error("missing label")
	}
}

// TestDropdownMenuRadio tests radio menu items
func TestDropdownMenuRadio(t *testing.T) {
	var buf bytes.Buffer

	node := dropdown.DropdownMenuRadioGroup(
		"option1",
		dropdown.DropdownMenuRadioItem("option1", "Option 1"),
		dropdown.DropdownMenuRadioItem("option2", "Option 2"),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "menuitemradio") {
		t.Error("missing radio role")
	}
}

// TestContextMenu tests the ContextMenu component
func TestContextMenu(t *testing.T) {
	var buf bytes.Buffer

	node := dropdown.ContextMenu(
		html.Div(g.Text("Right-click me")),
		dropdown.ContextMenuContent(
			dropdown.ContextMenuItem(g.Text("Copy")),
			dropdown.ContextMenuItem(g.Text("Paste")),
		),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "contextmenu") {
		t.Error("missing contextmenu event handler")
	}

	if !strings.Contains(html, "x-data") {
		t.Error("missing Alpine data")
	}
}
