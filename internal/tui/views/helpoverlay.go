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

// Render returns the help overlay.
func (h *HelpOverlay) Render() string {
	w := h.Width * 45 / 100
	if w < 38 {
		w = 38
	}
	if w > 56 {
		w = 56
	}

	title := styles.OverlayTitle.Render("Keybindings")

	keyStyle := styles.FooterKey.Width(12)
	descStyle := styles.Muted

	section := func(name string, rows [][2]string) string {
		out := "\n" + styles.Bold.Render(name) + "\n"
		for _, r := range rows {
			out += keyStyle.Render(r[0]) + descStyle.Render(r[1]) + "\n"
		}
		return out
	}

	content := title
	content += section("Navigate", [][2]string{
		{"j / ↓", "next session"},
		{"k / ↑", "prev session"},
		{"g g", "first"},
		{"G", "last"},
	})
	content += section("Actions", [][2]string{
		{"⏎ enter", "attach"},
		{"n", "new session"},
		{"x", "kill (confirm)"},
		{"r", "rename"},
		{"R", "refresh"},
	})
	content += section("Views", [][2]string{
		{"/", "command palette"},
		{"p", "workspaces"},
		{"?", "this help"},
		{"esc", "close overlay"},
	})
	content += section("Quit", [][2]string{
		{"q", "exit"},
		{"ctrl+c", "exit"},
	})

	box := styles.Overlay.Width(w).Render(content)
	return lipgloss.Place(h.Width, h.Height,
		lipgloss.Center, lipgloss.Center, box)
}
