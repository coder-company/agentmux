package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestTruncateFitsTerminalWidth(t *testing.T) {
	tests := []struct {
		name  string
		value string
		width int
	}{
		{"ascii", "abcdef", 4},
		{"wide runes", "界界界", 3},
		{"zero", "abcdef", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.value, tt.width)
			if lipgloss.Width(got) > tt.width {
				t.Fatalf("width %d exceeds limit %d for %q", lipgloss.Width(got), tt.width, got)
			}
		})
	}
}

func TestPadRightFitsExactWidth(t *testing.T) {
	got := PadRight("abc", 8)
	if lipgloss.Width(got) != 8 {
		t.Fatalf("width = %d, want 8", lipgloss.Width(got))
	}

	got = PadRight("abcdefgh", 4)
	if lipgloss.Width(got) != 4 {
		t.Fatalf("truncated width = %d, want 4", lipgloss.Width(got))
	}
}
