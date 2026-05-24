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

// Render returns the string representation.
func (sl *SessionList) Render(width, height int) string {
	if len(sl.Sessions) == 0 {
		empty := styles.Muted.Render("No tmux sessions") + "\n\n"
		empty += styles.Subtle.Render("  n") + styles.Muted.Render("  create session") + "\n"
		empty += styles.Subtle.Render("  p") + styles.Muted.Render("  launch workspace") + "\n"
		return empty
	}

	var out string
	for i, s := range sl.Sessions {
		if i >= height {
			break
		}

		// Status indicator
		var indicator string
		if s.Attached {
			indicator = styles.Attached.Render("● ")
		} else {
			indicator = styles.Detached.Render("  ")
		}

		// Session name and metadata
		name := s.Name
		meta := fmt.Sprintf(" %d win", s.Windows)
		if !s.Created.IsZero() {
			meta += " · " + relativeTime(s.Created)
		}

		if i == sl.Cursor {
			label := indicator + name
			out += styles.Selected.Width(width).Render(label) + "\n"
		} else {
			nameStr := styles.Normal.Render(indicator + name)
			metaStr := styles.Muted.Render(meta)
			out += nameStr + metaStr + "\n"
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
