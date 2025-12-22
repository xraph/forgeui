package modal_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/modal"
)

// TestModal tests the base Modal component
func TestModal(t *testing.T) {
	tests := []struct {
		name    string
		props   modal.ModalProps
		trigger g.Node
		content []g.Node
		wantRe  string
	}{
		{
			name: "renders modal with default props",
			props: modal.ModalProps{
				Size:          forgeui.SizeMD,
				CloseOnEscape: true,
				ShowClose:     true,
			},
			trigger: button.Button(g.Text("Open")),
			content: []g.Node{g.Text("Modal content")},
			wantRe:  `x-data.*open.*false`,
		},
		{
			name: "renders with large size",
			props: modal.ModalProps{
				Size: forgeui.SizeLG,
			},
			trigger: button.Button(g.Text("Open")),
			content: []g.Node{g.Text("Content")},
			wantRe:  `max-w-lg`,
		},
		{
			name: "renders with escape handler",
			props: modal.ModalProps{
				CloseOnEscape: true,
			},
			trigger: button.Button(g.Text("Open")),
			content: []g.Node{g.Text("Content")},
			wantRe:  `keydown\.escape\.window`,
		},
		{
			name: "renders with close button",
			props: modal.ModalProps{
				ShowClose: true,
			},
			trigger: button.Button(g.Text("Open")),
			content: []g.Node{g.Text("Content")},
			wantRe:  `aria-label="Close"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := modal.Modal(tt.props, tt.trigger, tt.content...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()
			if tt.wantRe != "" {
				re := regexp.MustCompile(tt.wantRe)
				if !re.MatchString(html) {
					t.Errorf("output does not match pattern %q\ngot: %s", tt.wantRe, html)
				}
			}
		})
	}
}

// TestModalHeader tests the ModalHeader component
func TestModalHeader(t *testing.T) {
	var buf bytes.Buffer
	node := modal.ModalHeader("Test Title", "Test description")
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Test Title") {
		t.Error("missing title")
	}
	if !strings.Contains(html, "Test description") {
		t.Error("missing description")
	}
}

// TestDialog tests the Dialog component
func TestDialog(t *testing.T) {
	var buf bytes.Buffer
	node := modal.Dialog(
		modal.DialogTrigger(button.Button(g.Text("Open"))),
		modal.DialogContent(
			modal.DialogHeader(
				modal.DialogTitle("Title"),
				modal.DialogDescription("Description"),
			),
		),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "x-data") {
		t.Error("missing Alpine data")
	}
	if !strings.Contains(html, "Title") {
		t.Error("missing title")
	}
}

// TestAlertDialog tests the AlertDialog component
func TestAlertDialog(t *testing.T) {
	var buf bytes.Buffer
	node := modal.AlertDialog(
		modal.AlertDialogTrigger(button.Button(g.Text("Delete"))),
		modal.AlertDialogContent(
			modal.AlertDialogHeader(
				modal.AlertDialogTitle("Are you sure?"),
			),
		),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "alertdialog") {
		t.Error("missing alertdialog role")
	}
}

// TestDrawer tests the Drawer component
func TestDrawer(t *testing.T) {
	sides := []forgeui.Side{
		forgeui.SideLeft,
		forgeui.SideRight,
		forgeui.SideTop,
		forgeui.SideBottom,
	}

	for _, side := range sides {
		t.Run(string(side), func(t *testing.T) {
			var buf bytes.Buffer
			node := modal.Drawer(
				modal.DrawerProps{Side: side},
				button.Button(g.Text("Open")),
				g.Text("Content"),
			)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()
			if !strings.Contains(html, "x-data") {
				t.Error("missing Alpine data")
			}
		})
	}
}

// TestSheet tests the Sheet component
func TestSheet(t *testing.T) {
	var buf bytes.Buffer
	node := modal.Sheet(
		modal.SheetTrigger(button.Button(g.Text("Open"))),
		modal.SheetContent(
			forgeui.SideRight,
			modal.SheetHeader(
				modal.SheetTitle("Title"),
			),
		),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "x-data") {
		t.Error("missing Alpine data")
	}
}

