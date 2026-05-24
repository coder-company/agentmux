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

	// If terminal is too narrow, drop the right side
	brandW := lipgloss.Width(brand)
	rightW := lipgloss.Width(right)
	gap := h.Width - brandW - rightW - 2
	var row string
	if gap < 1 {
		// Just show brand
		row = brand
	} else {
		row = brand + lipgloss.NewStyle().Width(gap).Render("") + right
	}

	header := styles.HeaderStyle.Width(h.Width).Render(row)
	ruleW := h.Width
	if ruleW > 0 {
		rule := styles.Separator.Render(repeatChar('─', ruleW))
		return header + "\n" + rule
	}
	return header
}

func repeatChar(ch rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
