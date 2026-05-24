package views

import (
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// HelpOverlay shows a keybinding reference.
type HelpOverlay struct {
	Width  int
	Height int
}

type helpSection struct {
	title    string
	bindings []helpRow
}

type helpRow struct {
	key  string
	desc string
}

// Render returns the help overlay.
func (h *HelpOverlay) Render() string {
	w := h.Width * 50 / 100
	if w < 40 {
		w = 40
	}
	if w > 64 {
		w = 64
	}

	sections := []helpSection{
		{
			title: "Navigation",
			bindings: []helpRow{
				{"↑/↓, j/k", "Move cursor"},
				{"g g", "Jump to top"},
				{"G", "Jump to bottom"},
			},
		},
		{
			title: "Sessions",
			bindings: []helpRow{
				{"enter", "Attach to session"},
				{"n", "New session"},
				{"x", "Kill session"},
				{"r", "Rename session"},
				{"R", "Refresh list"},
			},
		},
		{
			title: "Views",
			bindings: []helpRow{
				{"/", "Command palette"},
				{"p", "Workspace launcher"},
				{"?", "Toggle this help"},
				{"esc", "Close overlay"},
			},
		},
		{
			title: "General",
			bindings: []helpRow{
				{"q, Ctrl+C", "Quit"},
			},
		},
	}

	title := styles.PanelTitle.Render("Keybindings")
	var content string
	for _, sec := range sections {
		content += "\n" + styles.Bold.Render(sec.title) + "\n"
		for _, b := range sec.bindings {
			key := styles.FooterKey.Width(14).Render(b.key)
			desc := styles.Muted.Render(b.desc)
			content += key + desc + "\n"
		}
	}

	box := styles.Overlay.Width(w).Render(title + content)
	return lipgloss.Place(h.Width, h.Height,
		lipgloss.Center, lipgloss.Center, box)
}
