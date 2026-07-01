package views

import (
	"testing"
	"time"

	"agentmux/internal/core"
)

func TestSessionFilteringMatchesNameAndDirectory(t *testing.T) {
	v := &SessionsView{
		AllSessions: []core.Session{
			{Name: "api", Directory: "/srv/api"},
			{Name: "worker", Directory: "/srv/jobs"},
			{Name: "notes", Directory: "/home/me/docs"},
		},
	}

	v.SetFilter("srv")
	if got := v.VisibleCount(); got != 2 {
		t.Fatalf("visible = %d, want 2", got)
	}

	v.SetFilter("note")
	if got := v.VisibleCount(); got != 1 {
		t.Fatalf("visible = %d, want 1", got)
	}
	if sel := v.List.Selected(); sel == nil || sel.Name != "notes" {
		t.Fatalf("selected = %#v, want notes", sel)
	}
}

func TestSessionSortModes(t *testing.T) {
	now := time.Now()
	v := &SessionsView{
		AllSessions: []core.Session{
			{Name: "beta", Windows: 2, Created: now.Add(-2 * time.Hour)},
			{Name: "alpha", Windows: 5, Created: now.Add(-1 * time.Hour)},
			{Name: "gamma", Windows: 1, Created: now.Add(-3 * time.Hour), Attached: true},
		},
	}

	v.applySessions("")
	if got := v.List.Sessions[0].Name; got != "alpha" {
		t.Fatalf("newest first = %q, want alpha", got)
	}

	v.CycleSort()
	if got := v.List.Sessions[0].Name; got != "alpha" {
		t.Fatalf("name first = %q, want alpha", got)
	}

	v.CycleSort()
	if got := v.List.Sessions[0].Name; got != "alpha" {
		t.Fatalf("windows first = %q, want alpha", got)
	}

	v.CycleSort()
	if got := v.List.Sessions[0].Name; got != "gamma" {
		t.Fatalf("attached first = %q, want gamma", got)
	}
}

func TestSessionLayoutCycle(t *testing.T) {
	v := &SessionsView{}

	if v.LayoutLabel() != "split" {
		t.Fatalf("initial layout = %q, want split", v.LayoutLabel())
	}
	v.CycleLayout()
	if v.LayoutLabel() != "list" {
		t.Fatalf("second layout = %q, want list", v.LayoutLabel())
	}
	v.CycleLayout()
	if v.LayoutLabel() != "preview" {
		t.Fatalf("third layout = %q, want preview", v.LayoutLabel())
	}
}
