package components

import (
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// HelpBinding is a key/description pair.
type HelpBinding struct {
	Key  string
	Desc string
}

// RenderFooter renders a keybinding footer that fits in the given width.
func RenderFooter(bindings []HelpBinding, width int) string {
	sep := styles.HeaderDim.Render(" · ")

	var row string
	for i, b := range bindings {
		entry := styles.FooterKey.Render(b.Key) + " " + styles.FooterDesc.Render(b.Desc)
		if i > 0 {
			candidate := row + sep + entry
			if lipgloss.Width(candidate) > width-2 {
				break
			}
			row = candidate
		} else {
			row = entry
		}
	}

	return styles.FooterStyle.Width(width).Render(row)
}

// SessionsFooter returns the main view footer.
func SessionsFooter(width int) string {
	return RenderFooter([]HelpBinding{
		{"j/k", "move"},
		{"⏎", "attach"},
		{"n", "new"},
		{"x", "kill"},
		{"r", "rename"},
		{"/", "palette"},
		{"p", "workspaces"},
		{"?", "help"},
		{"q", "quit"},
	}, width)
}
