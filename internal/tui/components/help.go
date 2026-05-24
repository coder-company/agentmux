package components

import "agentmux/internal/tui/styles"

// HelpBar renders the bottom help line.
type HelpBar struct {
	Bindings []HelpBinding
}

// HelpBinding is a key/description pair.
type HelpBinding struct {
	Key  string
	Desc string
}

// DefaultHelp returns the standard help bindings.
func DefaultHelp() *HelpBar {
	return &HelpBar{
		Bindings: []HelpBinding{
			{"enter", "attach"},
			{"n", "new"},
			{"k", "kill"},
			{"r", "rename"},
			{"/", "search"},
			{"p", "projects"},
			{"q", "quit"},
		},
	}
}

// Render returns the help bar string.
func (h *HelpBar) Render(width int) string {
	var parts string
	for i, b := range h.Bindings {
		if i > 0 {
			parts += "  "
		}
		parts += styles.Status.Render(b.Key) + " " + styles.Muted.Render(b.Desc)
	}
	return parts
}
