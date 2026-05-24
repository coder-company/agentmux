package components

import (
	"strings"

	"agentmux/internal/tui/styles"
)

// Preview renders pane capture output.
type Preview struct {
	Content string
	Title   string
	Error   string
}

// Render returns the preview panel content.
func (p *Preview) Render(width, height int) string {
	title := styles.Title.Render(p.Title)

	if p.Error != "" {
		return title + "\n" + styles.Muted.Render("⚠ "+p.Error)
	}

	if p.Content == "" {
		return title + "\n" + styles.Muted.Render("No output captured.\nSelect a session to preview pane content.")
	}

	// Trim trailing blank lines
	content := strings.TrimRight(p.Content, "\n ")

	lines := strings.Split(content, "\n")
	maxLines := height - 3
	if maxLines < 1 {
		maxLines = 1
	}
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	// Truncate long lines to fit width
	for i, line := range lines {
		runes := []rune(line)
		if len(runes) > width {
			lines[i] = string(runes[:width-1]) + "…"
		}
	}

	return title + "\n" + strings.Join(lines, "\n")
}
