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
	w := h.Width * 52 / 100
	if w < 44 {
		w = 44
	}
	if w > 68 {
		w = 68
	}
	if h.Width > 0 && w > h.Width-4 {
		w = h.Width - 4
	}
	if w < 26 {
		w = 26
	}
	innerW := w - 6

	title := styles.OverlayTitle.Render("Keybindings")

	keyStyle := styles.FooterKey.Width(11)
	descStyle := styles.Muted

	section := func(name string, rows [][2]string) string {
		out := "\n" + styles.PanelMeta.Render(name) + "\n"
		for _, r := range rows {
			line := keyStyle.Render(r[0]) + descStyle.Render(styles.Truncate(r[1], innerW-11))
			out += line + "\n"
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
	content += section("Shape", [][2]string{
		{"/", "filter sessions live"},
		{"esc", "clear active filter"},
		{"s", "cycle sort"},
		{"tab", "cycle split/list/preview"},
	})
	content += section("Actions", [][2]string{
		{"⏎ enter", "attach"},
		{"n", "new session"},
		{"x", "kill (confirm)"},
		{"r", "rename"},
		{"R", "refresh"},
	})
	content += section("Views", [][2]string{
		{":", "command palette"},
		{"ctrl+p", "command palette"},
		{"p", "workspaces"},
		{"?", "this help"},
		{"esc", "close overlay"},
	})
	content += section("Quit", [][2]string{
		{"q", "exit"},
		{"ctrl+c", "exit"},
	})
	content += "\n" + styles.HeaderDim.Render(styles.Truncate("esc, ?, or q closes this panel", innerW))

	box := styles.Overlay.Width(w).Render(content)
	return lipgloss.Place(h.Width, h.Height,
		lipgloss.Center, lipgloss.Center, box)
}
