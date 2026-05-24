package components

import (
	"strings"

	"agentmux/internal/tui/styles"
)

// Preview renders pane capture output.
type Preview struct {
	Content string
	Title   string
}

// Render returns the preview panel content.
func (p *Preview) Render(width, height int) string {
	title := styles.Title.Render(p.Title)

	if p.Content == "" {
		return title + "\n" + styles.Muted.Render("No preview available")
	}

	lines := strings.Split(p.Content, "\n")
	if len(lines) > height-2 {
		lines = lines[len(lines)-(height-2):]
	}

	content := strings.Join(lines, "\n")
	return title + "\n" + styles.Muted.Render(content)
}
