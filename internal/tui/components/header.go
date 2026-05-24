package components

import (
	"fmt"

	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// Header renders the top bar.
type Header struct {
	Mode         string
	SessionCount int
	Width        int
}

// Render returns the header.
func (h *Header) Render() string {
	brand := styles.HeaderBrand.Render("⣿ agentmux")

	var right string
	if h.SessionCount > 0 {
		right = styles.HeaderDim.Render(fmt.Sprintf("%d sessions", h.SessionCount))
	} else {
		right = styles.HeaderDim.Render("no sessions")
	}

	if h.Mode != "" {
		right = styles.Muted.Render("["+h.Mode+"]") + "  " + right
	}

	gap := h.Width - lipgloss.Width(brand) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	row := brand + lipgloss.NewStyle().Width(gap).Render("") + right
	header := styles.HeaderStyle.Width(h.Width).Render(row)
	rule := styles.Separator.Width(h.Width).Render(repeatChar('─', h.Width))

	return header + "\n" + rule
}

func repeatChar(ch rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
