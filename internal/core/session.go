package core

import "time"

// Session represents a tmux session.
type Session struct {
	Name      string
	Windows   int
	Created   time.Time
	Attached  bool
	LastUsed  time.Time
	Directory string
}
