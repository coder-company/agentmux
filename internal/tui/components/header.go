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
	brand := styles.HeaderBrand.Render("agentmux")
	section := styles.HeaderDim.Render(" sessions")

	var right string
	if h.SessionCount > 0 {
		right = styles.HeaderDim.Render(fmt.Sprintf("%d active", h.SessionCount))
	} else {
		right = styles.HeaderDim.Render("no active sessions")
	}

	if h.Mode != "" {
		right = styles.HeaderMode.Render(h.Mode) + "  " + right
	}

	// If terminal is too narrow, drop the right side
	left := brand + section
	brandW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	gap := h.Width - brandW - rightW - 2
	var row string
	if gap < 1 {
		row = left
	} else {
		row = left + lipgloss.NewStyle().Width(gap).Render("") + right
	}

	header := styles.HeaderStyle.Width(h.Width).Render(styles.Truncate(row, max(0, h.Width)))
	ruleW := h.Width
	if ruleW > 0 {
		rule := styles.Separator.Render(repeatChar('─', ruleW))
		return header + "\n" + rule
	}
	return header
}

func repeatChar(ch rune, n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
