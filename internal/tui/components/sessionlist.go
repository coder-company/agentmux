package components

import (
	"fmt"

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

// Render returns the string representation.
func (sl *SessionList) Render(width int) string {
	if len(sl.Sessions) == 0 {
		return styles.Muted.Render("  No tmux sessions.\n\n  Press n to create one,\n  or p to launch a workspace.")
	}

	var out string
	for i, s := range sl.Sessions {
		indicator := "  "
		if s.Attached {
			indicator = "● "
		}

		label := fmt.Sprintf("%s%s (%d win)", indicator, s.Name, s.Windows)

		if i == sl.Cursor {
			out += styles.Selected.Width(width-2).Render(label) + "\n"
		} else {
			out += styles.Normal.Width(width-2).Render(label) + "\n"
		}
	}
	return out
}
