package components

import (
	"fmt"
	"time"

	"agentmux/internal/core"
	"agentmux/internal/tui/styles"
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
	if len(sl.Sessions) == 0 {
		out := "\n"
		out += styles.Muted.Render("  No tmux sessions running.") + "\n\n"
		out += styles.HeaderDim.Render("  n") + styles.Muted.Render("  new session") + "\n"
		out += styles.HeaderDim.Render("  p") + styles.Muted.Render("  launch workspace") + "\n"
		return out
	}

	var out string
	for i, s := range sl.Sessions {
		if i >= height {
			break
		}

		// Build the line content
		dot := "  "
		if s.Attached {
			dot = styles.ListDot.Render("● ")
		}

		name := s.Name
		meta := fmt.Sprintf("  %dw", s.Windows)
		if !s.Created.IsZero() {
			meta += " " + relativeTime(s.Created)
		}

		if i == sl.Cursor {
			// Selected: highlighted full row with arrow
			line := "▸ " + dot + name + styles.ListMeta.Render(meta)
			out += styles.ListSelected.Width(width).Render(line) + "\n"
		} else {
			// Normal row
			line := "  " + dot + styles.ListNormal.Render(name) + styles.ListMeta.Render(meta)
			out += line + "\n"
		}
	}
	return out
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
