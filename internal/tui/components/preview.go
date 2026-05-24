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
	Windows []string
	Error   string
}

// Render returns the preview panel content, fitting within width x height.
func (p *Preview) Render(width, height int) string {
	if width < 4 || height < 2 {
		return ""
	}

	var lines []string

	// Title
	if p.Title != "" {
		title := styles.PreviewTitle.Render(p.Title)
		if p.Dir != "" {
			title += "  " + styles.PreviewPath.Render(p.Dir)
		}
		lines = append(lines, title)
	} else {
		lines = append(lines, styles.PreviewPath.Render("no session selected"))
	}

	// Windows
	if len(p.Windows) > 0 {
		winLine := styles.PreviewDim.Render("win: " + strings.Join(p.Windows, ", "))
		lines = append(lines, winLine)
	}

	lines = append(lines, "") // spacer

	// Error state
	if p.Error != "" {
		lines = append(lines, styles.Muted.Render("  "+p.Error))
		return strings.Join(lines, "\n")
	}

	// Empty state
	if strings.TrimSpace(p.Content) == "" {
		lines = append(lines, styles.Muted.Render("  (no output)"))
		return strings.Join(lines, "\n")
	}

	// Content — last N lines that fit
	content := strings.TrimRight(p.Content, "\n \t")
	contentLines := strings.Split(content, "\n")
	available := height - len(lines)
	if available < 1 {
		available = 1
	}
	if len(contentLines) > available {
		contentLines = contentLines[len(contentLines)-available:]
	}

	// Truncate long lines
	for i, line := range contentLines {
		runes := []rune(line)
		if len(runes) > width {
			contentLines[i] = string(runes[:width-1]) + "…"
		}
	}

	lines = append(lines, contentLines...)
	return strings.Join(lines, "\n")
}
