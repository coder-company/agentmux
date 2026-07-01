package views

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestFuzzyContains(t *testing.T) {
	tests := []struct {
		haystack string
		needle   string
		want     bool
	}{
		{"new session", "new", true},
		{"new session", "ns", true},
		{"new session", "nss", true},
		{"new session", "newsession", true},
		{"new session", "xyz", false},
		{"kill", "kl", true},
		{"kill", "ki", true},
		{"kill", "kx", false},
		{"", "", true},
		{"abc", "", true},
		{"", "a", false},
		{"attach: main", "main", true},
		{"launch: api", "api", true},
		{"launch: api", "la", true},
	}
	for _, tt := range tests {
		t.Run(tt.haystack+"/"+tt.needle, func(t *testing.T) {
			got := fuzzyContains(tt.haystack, tt.needle)
			if got != tt.want {
				t.Errorf("fuzzyContains(%q, %q) = %v, want %v",
					tt.haystack, tt.needle, got, tt.want)
			}
		})
	}
}

func TestPaletteSetQuery(t *testing.T) {
	actions := []Action{
		{Name: "New Session", Desc: "Create session"},
		{Name: "Kill Session", Desc: "Destroy session"},
		{Name: "Refresh", Desc: "Reload list"},
		{Name: "Quit", Desc: "Exit"},
	}
	p := NewPalette(actions)

	// Empty query shows all
	p.SetQuery("")
	if len(p.Filtered) != 4 {
		t.Errorf("empty query: got %d filtered, want 4", len(p.Filtered))
	}

	// Filter by name
	p.SetQuery("ses")
	if len(p.Filtered) != 2 {
		t.Errorf("query 'ses': got %d filtered, want 2", len(p.Filtered))
	}

	// Filter by desc
	p.SetQuery("exit")
	if len(p.Filtered) != 1 {
		t.Errorf("query 'exit': got %d filtered, want 1", len(p.Filtered))
	}
	if p.Filtered[0].Name != "Quit" {
		t.Errorf("expected Quit, got %q", p.Filtered[0].Name)
	}

	// No matches
	p.SetQuery("zzzzz")
	if len(p.Filtered) != 0 {
		t.Errorf("query 'zzzzz': got %d filtered, want 0", len(p.Filtered))
	}

	// Cursor resets on query change
	p.SetQuery("")
	p.Cursor = 3
	p.SetQuery("q")
	if p.Cursor != 0 {
		t.Errorf("cursor should reset to 0 on query change, got %d", p.Cursor)
	}
}

func TestPaletteNavigation(t *testing.T) {
	actions := []Action{
		{Name: "A"}, {Name: "B"}, {Name: "C"},
	}
	p := NewPalette(actions)
	p.SetQuery("")

	if p.Cursor != 0 {
		t.Fatalf("initial cursor should be 0")
	}

	p.MoveDown()
	if p.Cursor != 1 {
		t.Errorf("after MoveDown: got %d, want 1", p.Cursor)
	}

	p.MoveDown()
	p.MoveDown() // should clamp
	if p.Cursor != 2 {
		t.Errorf("after clamp: got %d, want 2", p.Cursor)
	}

	p.MoveUp()
	if p.Cursor != 1 {
		t.Errorf("after MoveUp: got %d, want 1", p.Cursor)
	}

	// Selected
	sel := p.Selected()
	if sel == nil || sel.Name != "B" {
		t.Errorf("expected B, got %v", sel)
	}
}

func TestPaletteTypeAndBackspace(t *testing.T) {
	p := NewPalette([]Action{{Name: "Test"}})
	p.TypeChar('h')
	p.TypeChar('i')
	if p.Query != "hi" {
		t.Errorf("query: got %q, want 'hi'", p.Query)
	}

	p.Backspace()
	if p.Query != "h" {
		t.Errorf("after backspace: got %q, want 'h'", p.Query)
	}

	p.Backspace()
	if p.Query != "" {
		t.Errorf("after double backspace: got %q, want ''", p.Query)
	}

	p.Backspace() // should not panic
	if p.Query != "" {
		t.Errorf("extra backspace: got %q, want ''", p.Query)
	}
}

func TestPaletteVisibleActionsKeepsCursorInWindow(t *testing.T) {
	var actions []Action
	for i := 0; i < 20; i++ {
		actions = append(actions, Action{Name: fmt.Sprintf("Action %02d", i)})
	}
	p := NewPalette(actions)
	p.Cursor = 14

	start, visible := p.visibleActions(5)
	if start > p.Cursor || start+len(visible) <= p.Cursor {
		t.Fatalf("cursor %d not visible in window %d-%d", p.Cursor, start, start+len(visible)-1)
	}
}

func TestPaletteRenderFitsScreenWidth(t *testing.T) {
	p := NewPalette([]Action{
		{Name: "Long action name that should fit", Desc: "Long action description that should also be constrained", Key: "n"},
	})
	p.Width = 60
	p.Height = 20

	out := p.Render()
	if width := lipgloss.Width(out); width > 60 {
		t.Fatalf("rendered width = %d, want <= 60", width)
	}
}
