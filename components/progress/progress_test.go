package progress

import (
	"bytes"
	"strings"
	"testing"
)

func TestProgress(t *testing.T) {
	progress := Progress(WithValue(50))

	var buf bytes.Buffer
	if err := progress.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	html := buf.String()

	if !strings.Contains(html, `role="progressbar"`) {
		t.Error("expected role attribute")
	}

	if !strings.Contains(html, `aria-valuenow="50"`) {
		t.Error("expected aria-valuenow attribute")
	}

	if !strings.Contains(html, "width: 50%") {
		t.Error("expected 50% width style")
	}
}

func TestProgress_Bounds(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  string
	}{
		{"below zero", -10, `aria-valuenow="0"`},
		{"zero", 0, `aria-valuenow="0"`},
		{"middle", 50, `aria-valuenow="50"`},
		{"hundred", 100, `aria-valuenow="100"`},
		{"above hundred", 150, `aria-valuenow="100"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := Progress(WithValue(tt.value))

			var buf bytes.Buffer
			if err := progress.Render(&buf); err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("expected %v in progress", tt.want)
			}
		})
	}
}
