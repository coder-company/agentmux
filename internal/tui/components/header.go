package components

import (
	"fmt"

	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// Header renders the top bar with app name, mode, and session count.
type Header struct {
	Mode         string
	SessionCount int
	Width        int
}

// Render returns the header bar.
func (h *Header) Render() string {
	brand := styles.HeaderBrand.Render("⣿ agentmux")

	mode := ""
	if h.Mode != "" {
		mode = styles.Muted.Render(" · ") + styles.HeaderMeta.Render(h.Mode)
	}

	count := styles.HeaderMeta.Render(fmt.Sprintf("%d sessions", h.SessionCount))

	left := brand + mode
	right := count

	gap := h.Width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}
	spacer := lipgloss.NewStyle().Width(gap).Render("")

	row := lipgloss.JoinHorizontal(lipgloss.Center, left, spacer, right)
	return styles.Header.Width(h.Width).Render(row)
}
