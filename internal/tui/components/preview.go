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

// Render returns the preview panel content.
func (p *Preview) Render(width, height int) string {
	if width < 4 || height < 2 {
		return ""
	}

	// Title line with session name and directory
	var titleLine string
	if p.Title != "" {
		titleLine = styles.PanelTitle.Render(p.Title)
		if p.Dir != "" {
			titleLine += " " + styles.Muted.Render(p.Dir)
		}
	} else {
		titleLine = styles.PanelTitleDim.Render("Preview")
	}

	// Window list
	var windowLine string
	if len(p.Windows) > 0 {
		names := strings.Join(p.Windows, " │ ")
		windowLine = styles.Subtle.Render("windows: " + names)
	}

	if p.Error != "" {
		out := titleLine + "\n"
		if windowLine != "" {
			out += windowLine + "\n"
		}
		out += "\n" + styles.Warning.Render("⚠ ") + styles.Muted.Render(p.Error)
		return out
	}

	if p.Content == "" {
		out := titleLine + "\n"
		if windowLine != "" {
			out += windowLine + "\n"
		}
		out += "\n" + styles.Muted.Render("No pane output captured.")
		return out
	}

	// Trim trailing blank lines
	content := strings.TrimRight(p.Content, "\n \t")
	if content == "" {
		out := titleLine + "\n"
		if windowLine != "" {
			out += windowLine + "\n"
		}
		out += "\n" + styles.Muted.Render("(empty)")
		return out
	}

	lines := strings.Split(content, "\n")
	headerLines := 2 // title + blank
	if windowLine != "" {
		headerLines = 3
	}
	maxLines := height - headerLines
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

	out := titleLine + "\n"
	if windowLine != "" {
		out += windowLine + "\n"
	}
	out += strings.Join(lines, "\n")
	return out
}
