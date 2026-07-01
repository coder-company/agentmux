package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestPreviewRenderFitsBounds(t *testing.T) {
	p := Preview{
		Title:   "session-with-a-long-name-that-needs-truncation",
		Dir:     "/very/long/path/that/should/not/escape/the/preview/panel",
		Windows: []string{"editor", "server", "logs-with-a-long-name"},
		Content: "short\n" + strings.Repeat("x", 80),
	}

	out := p.Render(30, 5)
	lines := strings.Split(out, "\n")
	if len(lines) > 5 {
		t.Fatalf("height = %d, want <= 5", len(lines))
	}
	for _, line := range lines {
		if width := lipgloss.Width(line); width > 30 {
			t.Fatalf("line width = %d, want <= 30: %q", width, line)
		}
	}
}

func TestPreviewNoSelectionFitsHeight(t *testing.T) {
	p := Preview{}
	out := p.Render(24, 3)
	if lines := strings.Split(out, "\n"); len(lines) > 3 {
		t.Fatalf("height = %d, want <= 3", len(lines))
	}
}
