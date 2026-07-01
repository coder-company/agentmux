package components

import (
	"fmt"
	"strings"
	"time"

	"agentmux/internal/core"
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SessionList renders a list of sessions with a cursor.
type SessionList struct {
	Sessions []core.Session
	Cursor   int
}

// Selected returns the currently highlighted session, or nil.
func (sl *SessionList) Selected() *core.Session {
	if len(sl.Sessions) == 0 {
		return nil
	}
	if sl.Cursor >= len(sl.Sessions) {
		sl.Cursor = len(sl.Sessions) - 1
	}
	if sl.Cursor < 0 {
		sl.Cursor = 0
	}
	return &sl.Sessions[sl.Cursor]
}

// MoveUp moves the cursor up.
func (sl *SessionList) MoveUp() {
	if sl.Cursor > 0 {
		sl.Cursor--
	}
}

// MoveDown moves the cursor down.
func (sl *SessionList) MoveDown() {
	if sl.Cursor < len(sl.Sessions)-1 {
		sl.Cursor++
	}
}

// MoveTop jumps to the first item.
func (sl *SessionList) MoveTop() {
	sl.Cursor = 0
}

// MoveBottom jumps to the last item.
func (sl *SessionList) MoveBottom() {
	if len(sl.Sessions) > 0 {
		sl.Cursor = len(sl.Sessions) - 1
	}
}

// Render returns the rendered list.
func (sl *SessionList) Render(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	if len(sl.Sessions) == 0 {
		return renderEmptySessions(width, height)
	}

	var lines []string
	for i, s := range sl.Sessions {
		if len(lines) >= height {
			break
		}
		lines = append(lines, sl.renderRow(i, s, width))
	}

	if len(sl.Sessions) > height {
		more := fmt.Sprintf("  %d more below", len(sl.Sessions)-height)
		lines[height-1] = styles.PanelMeta.Render(styles.Truncate(more, width))
	}

	return strings.Join(lines, "\n")
}

func (sl *SessionList) renderRow(i int, s core.Session, width int) string {
	marker := " "
	if s.Attached {
		marker = "●"
	}

	meta := fmt.Sprintf("%dw", s.Windows)
	if !s.Created.IsZero() {
		meta += " " + relativeTime(s.Created)
	}

	cursor := "  "
	if i == sl.Cursor {
		cursor = "› "
	}

	prefix := cursor + marker + " "
	prefixW := lipgloss.Width(prefix)
	metaW := lipgloss.Width(meta)
	nameW := width - prefixW - metaW - 2
	if nameW < 8 {
		nameW = width - prefixW
		meta = ""
		metaW = 0
	}
	if nameW < 1 {
		nameW = 1
	}

	name := styles.Truncate(s.Name, nameW)
	gap := width - prefixW - lipgloss.Width(name) - metaW
	if gap < 1 {
		gap = 1
	}

	plain := styles.Truncate(prefix+name+strings.Repeat(" ", gap)+meta, width)
	if i == sl.Cursor {
		return styles.ListSelected.Width(width).Render(plain)
	}

	dot := " "
	if s.Attached {
		dot = styles.ListDot.Render("●")
	}

	line := "  " + dot + " " + styles.ListNormal.Render(name)
	if meta != "" {
		normalGap := width - lipgloss.Width(line) - metaW
		if normalGap < 1 {
			normalGap = 1
		}
		line += strings.Repeat(" ", normalGap) + styles.ListMeta.Render(meta)
	}
	return line
}

func renderEmptySessions(width, height int) string {
	rows := []string{
		"",
		styles.Bold.Render("  No tmux sessions"),
		styles.Muted.Render("  Create a detached session or launch a configured workspace."),
		"",
		actionHint("n", "new session"),
		actionHint("p", "workspace launcher"),
		actionHint("/", "command palette"),
	}

	if height < len(rows) {
		rows = rows[:height]
	}

	for i, row := range rows {
		if lipgloss.Width(row) > width {
			rows[i] = styles.Truncate(row, width)
		}
	}
	return strings.Join(rows, "\n")
}

func actionHint(key, desc string) string {
	return "  " + styles.FooterKey.Render(key) + "  " + styles.Muted.Render(desc)
}

func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
