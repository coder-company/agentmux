package components

import (
	"strings"

	"agentmux/internal/tui/styles"
)

// Preview renders pane capture output.
type Preview struct {
	Content string
	Title   string
	Dir     string
	Error   string
}

// Render returns the preview panel content.
func (p *Preview) Render(width, height int) string {
	if width < 4 || height < 2 {
		return ""
	}

	// Title line
	var titleLine string
	if p.Title != "" {
		titleLine = styles.PanelTitle.Render(p.Title)
		if p.Dir != "" {
			titleLine += " " + styles.Muted.Render(p.Dir)
		}
	} else {
		titleLine = styles.PanelTitleDim.Render("Preview")
	}

	if p.Error != "" {
		return titleLine + "\n\n" + styles.Warning.Render("⚠ ") + styles.Muted.Render(p.Error)
	}

	if p.Content == "" {
		hint := styles.Muted.Render("No pane output captured.")
		return titleLine + "\n\n" + hint
	}

	// Trim trailing blank lines
	content := strings.TrimRight(p.Content, "\n \t")
	if content == "" {
		return titleLine + "\n\n" + styles.Muted.Render("(empty)")
	}

	lines := strings.Split(content, "\n")
	maxLines := height - 2
	if maxLines < 1 {
		maxLines = 1
	}
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	// Truncate long lines
	for i, line := range lines {
		runes := []rune(line)
		if len(runes) > width-1 {
			lines[i] = string(runes[:width-2]) + "…"
		}
	}

	return titleLine + "\n" + strings.Join(lines, "\n")
}
