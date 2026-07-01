package components

import (
	"fmt"
	"strings"
	"time"

	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// Preview renders pane capture output.
type Preview struct {
	Content     string
	Title       string
	Dir         string
	Windows     []string
	Created     time.Time
	Attached    bool
	WindowCount int
	Error       string
}

// Render returns the preview panel content, fitting within width x height.
func (p *Preview) Render(width, height int) string {
	if width < 4 || height < 2 {
		return ""
	}

	var lines []string

	if p.Title != "" {
		titleText := styles.Truncate(p.Title, width)
		title := styles.PreviewTitle.Render(titleText)
		if p.Dir != "" {
			dirW := width - lipgloss.Width(titleText) - 2
			if dirW > 10 {
				title += "  " + styles.PreviewPath.Render(styles.Truncate(p.Dir, dirW))
			}
		}
		lines = append(lines, title)
	} else {
		return p.renderNoSelection(width, height)
	}

	if detail := p.detailLine(width); detail != "" {
		lines = append(lines, styles.PreviewDim.Render(detail))
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

func (p *Preview) detailLine(width int) string {
	var parts []string
	if p.Attached {
		parts = append(parts, "attached")
	} else if p.Title != "" {
		parts = append(parts, "detached")
	}

	windows := p.WindowCount
	if windows == 0 {
		windows = len(p.Windows)
	}
	if windows > 0 {
		parts = append(parts, fmt.Sprintf("%dw", windows))
	}
	if !p.Created.IsZero() {
		parts = append(parts, "created "+age(p.Created))
	}
	if len(parts) == 0 {
		return ""
	}
	return styles.Truncate(strings.Join(parts, " · "), width)
}

func age(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func joinPreviewLines(lines []string, height int) string {
	if len(lines) > height {
		lines = lines[:height]
	}
	return strings.Join(lines, "\n")
}
