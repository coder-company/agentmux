package components

import (
	"strings"

	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
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

	if p.Title != "" {
		title := styles.PreviewTitle.Render(styles.Truncate(p.Title, width))
		if p.Dir != "" {
			dirW := width - lipgloss.Width(p.Title) - 2
			if dirW > 10 {
				title += "  " + styles.PreviewPath.Render(styles.Truncate(p.Dir, dirW))
			}
		}
		lines = append(lines, title)
	} else {
		return p.renderNoSelection(width, height)
	}

	if len(p.Windows) > 0 {
		winLine := styles.PreviewDim.Render(styles.Truncate("windows  "+strings.Join(p.Windows, ", "), width))
		lines = append(lines, winLine)
	}

	lines = append(lines, "")

	if p.Error != "" {
		lines = append(lines,
			styles.ErrorText.Render(styles.Truncate("  "+p.Error, width)),
			styles.Muted.Render(styles.Truncate("  The session is still selectable; pane capture failed.", width)),
		)
		return joinPreviewLines(lines, height)
	}

	if strings.TrimSpace(p.Content) == "" {
		lines = append(lines,
			styles.Muted.Render("  No pane output yet."),
			styles.Subtle.Render("  Attach to inspect the live session."),
		)
		return joinPreviewLines(lines, height)
	}

	content := strings.TrimRight(p.Content, "\n \t")
	contentLines := strings.Split(content, "\n")
	available := height - len(lines)
	if available < 1 {
		available = 1
	}
	if len(contentLines) > available {
		contentLines = contentLines[len(contentLines)-available:]
	}

	for i, line := range contentLines {
		contentLines[i] = styles.PreviewContent.Render(styles.Truncate(line, width))
	}

	lines = append(lines, contentLines...)
	return joinPreviewLines(lines, height)
}

func (p *Preview) renderNoSelection(width, height int) string {
	lines := []string{
		styles.PreviewTitle.Render("No session selected"),
		"",
		styles.Muted.Render(styles.Truncate("  Move through the list to preview pane output.", width)),
		"",
		previewHint("j/k", "move"),
		previewHint("enter", "attach"),
		previewHint("n", "new session"),
	}
	return joinPreviewLines(lines, height)
}

func previewHint(key, desc string) string {
	return "  " + styles.FooterKey.Render(key) + "  " + styles.Muted.Render(desc)
}

func joinPreviewLines(lines []string, height int) string {
	if len(lines) > height {
		lines = lines[:height]
	}
	return strings.Join(lines, "\n")
}
