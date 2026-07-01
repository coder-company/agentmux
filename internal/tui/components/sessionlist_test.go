package components

import (
	"strings"
	"testing"
	"time"

	"agentmux/internal/core"

	"github.com/charmbracelet/lipgloss"
)

func TestSessionListNavigation(t *testing.T) {
	sl := &SessionList{
		Sessions: []core.Session{
			{Name: "a"}, {Name: "b"}, {Name: "c"},
		},
	}

	if sl.Cursor != 0 {
		t.Fatalf("initial cursor: %d", sl.Cursor)
	}

	sl.MoveDown()
	if sl.Cursor != 1 {
		t.Errorf("after MoveDown: %d", sl.Cursor)
	}

	sl.MoveDown()
	sl.MoveDown() // clamp
	if sl.Cursor != 2 {
		t.Errorf("after clamp down: %d", sl.Cursor)
	}

	sl.MoveUp()
	if sl.Cursor != 1 {
		t.Errorf("after MoveUp: %d", sl.Cursor)
	}

	sl.MoveTop()
	if sl.Cursor != 0 {
		t.Errorf("after MoveTop: %d", sl.Cursor)
	}

	sl.MoveBottom()
	if sl.Cursor != 2 {
		t.Errorf("after MoveBottom: %d", sl.Cursor)
	}
}

func TestSessionListEmpty(t *testing.T) {
	sl := &SessionList{}

	if sel := sl.Selected(); sel != nil {
		t.Errorf("expected nil selected on empty list")
	}

	// Should not panic
	sl.MoveUp()
	sl.MoveDown()
	sl.MoveTop()
	sl.MoveBottom()
}

func TestSessionListSelected(t *testing.T) {
	sl := &SessionList{
		Sessions: []core.Session{
			{Name: "first"}, {Name: "second"},
		},
		Cursor: 1,
	}

	sel := sl.Selected()
	if sel == nil || sel.Name != "second" {
		t.Errorf("expected 'second', got %v", sel)
	}

	// Out-of-bounds cursor gets clamped
	sl.Cursor = 99
	sel = sl.Selected()
	if sel == nil || sel.Name != "second" {
		t.Errorf("expected clamp to 'second', got %v", sel)
	}

	sl.Cursor = -1
	sel = sl.Selected()
	if sel == nil || sel.Name != "first" {
		t.Errorf("expected clamp to 'first', got %v", sel)
	}
}

func TestRelativeTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		t    time.Time
		want string
	}{
		{now.Add(-30 * time.Second), "now"},
		{now.Add(-5 * time.Minute), "5m"},
		{now.Add(-3 * time.Hour), "3h"},
		{now.Add(-48 * time.Hour), "2d"},
	}

	for _, tt := range tests {
		got := relativeTime(tt.t)
		if got != tt.want {
			t.Errorf("relativeTime(%v ago): got %q, want %q",
				time.Since(tt.t), got, tt.want)
		}
	}
}

func TestSessionListRenderFitsWidth(t *testing.T) {
	sl := &SessionList{
		Sessions: []core.Session{
			{Name: "a-very-long-session-name-that-would-overflow", Windows: 12, Attached: true},
			{Name: "second", Windows: 1},
		},
	}

	out := sl.Render(24, 4)
	for _, line := range strings.Split(out, "\n") {
		if width := lipgloss.Width(line); width > 24 {
			t.Fatalf("line width = %d, want <= 24: %q", width, line)
		}
	}
}
