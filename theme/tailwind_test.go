package theme

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestTailwindConfigScript(t *testing.T) {
	script := TailwindConfigScript()

	var buf bytes.Buffer
	if err := script.Render(context.Background(), &buf); err != nil {
		t.Fatalf("failed to render script: %v", err)
	}

	output := buf.String()

	// Check that it renders a script tag
	if !strings.Contains(output, "<script") {
		t.Error("output should contain <script tag")
	}

	// Check that it contains the tailwind config
	if !strings.Contains(output, "tailwind.config") {
		t.Error("output should contain tailwind.config")
	}

	// Check that it contains OKLCH wrapper
	if !strings.Contains(output, "oklch(var(--background)") {
		t.Error("output should contain oklch wrapper for background color")
	}

	// Check that it contains alpha-value placeholder
	if !strings.Contains(output, "<alpha-value>") {
		t.Error("output should contain <alpha-value> placeholder")
	}

	t.Logf("Script output:\n%s", output)
}
