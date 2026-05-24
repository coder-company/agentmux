package tui

import (
	"agentmux/internal/core"
	"testing"
)

func TestGenerateSessionName(t *testing.T) {
	tests := []struct {
		existing []core.Session
		want     string
	}{
		{nil, "s0"},
		{[]core.Session{{Name: "s0"}}, "s1"},
		{[]core.Session{{Name: "s0"}, {Name: "s1"}}, "s2"},
		{[]core.Session{{Name: "other"}}, "s1"},
		{[]core.Session{{Name: "s0"}, {Name: "s2"}}, "s3"}, // len=2, tries s2 (exists) → s3
	}
	for _, tt := range tests {
		got := generateSessionName(tt.existing)
		if got != tt.want {
			t.Errorf("generateSessionName(%v) = %q, want %q", names(tt.existing), got, tt.want)
		}
	}
}

// Verify the last test case: with s0 and s2 existing, len=2, so it starts at
// "s2" which exists, then tries "s3".
func TestGenerateSessionNameSkipsExisting(t *testing.T) {
	existing := []core.Session{{Name: "s0"}, {Name: "s2"}}
	got := generateSessionName(existing)
	// len=2 → tries s2 (exists) → tries s3
	if got != "s3" {
		t.Errorf("got %q, want s3", got)
	}
}

func names(sessions []core.Session) []string {
	var out []string
	for _, s := range sessions {
		out = append(out, s.Name)
	}
	return out
}
