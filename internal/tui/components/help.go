package components

import (
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// HelpBinding is a key/description pair for the footer.
type HelpBinding struct {
	Key  string
	Desc string
}

// RenderFooter renders a compact keybinding footer bar.
func RenderFooter(bindings []HelpBinding, width int) string {
	var parts []string
	for _, b := range bindings {
		entry := styles.FooterKey.Render(b.Key) + " " + styles.FooterDesc.Render(b.Desc)
		parts = append(parts, entry)
	}

	sep := styles.FooterSep.Render(" │ ")
	var row string
	for i, p := range parts {
		if i > 0 {
			row += sep
		}
		// Stop adding if we'd overflow
		if lipgloss.Width(row)+lipgloss.Width(sep)+lipgloss.Width(p) > width-4 {
			break
		}
		row += p
	}

	return styles.Footer.Width(width).Render(row)
}

// SessionsFooter is the default footer for the sessions view.
func SessionsFooter(width int) string {
	return RenderFooter([]HelpBinding{
		{"↑↓/jk", "nav"},
		{"enter", "attach"},
		{"n", "new"},
		{"x", "kill"},
		{"r", "rename"},
		{"/", "search"},
		{"p", "projects"},
		{"?", "help"},
		{"q", "quit"},
	}, width)
}
