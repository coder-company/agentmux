package tmux

import (
	"strings"
	"testing"
	"time"
)

// helper to build a session line using the fieldSep delimiter
func sessionLine(name, windows, created, attached, dir string) string {
	return name + fieldSep + windows + fieldSep + created + fieldSep + attached + fieldSep + dir
}

func TestParseSessions(t *testing.T) {
	input := sessionLine("main", "3", "1700000000", "1", "/home/user/code") + "\n" +
		sessionLine("work", "1", "1700001000", "0", "/home/user/work") + "\n"

	sessions, err := ParseSessions(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	s := sessions[0]
	if s.Name != "main" {
		t.Errorf("name: got %q, want %q", s.Name, "main")
	}
	if s.Windows != 3 {
		t.Errorf("windows: got %d, want 3", s.Windows)
	}
	if !s.Attached {
		t.Error("expected attached=true")
	}
	if s.Directory != "/home/user/code" {
		t.Errorf("directory: got %q", s.Directory)
	}
	want := time.Unix(1700000000, 0)
	if !s.Created.Equal(want) {
		t.Errorf("created: got %v, want %v", s.Created, want)
	}

	s2 := sessions[1]
	if s2.Attached {
		t.Error("expected attached=false for second session")
	}
}

func TestParseSessionsEmpty(t *testing.T) {
	sessions, err := ParseSessions("")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions from empty input, got %d", len(sessions))
	}
}

func TestParseSessionsMalformed(t *testing.T) {
	input := "garbage line without delimiter\n" + sessionLine("good", "2", "1700000000", "0", "/tmp") + "\n"
	sessions, err := ParseSessions(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session (skip malformed), got %d", len(sessions))
	}
	if sessions[0].Name != "good" {
		t.Errorf("expected name 'good', got %q", sessions[0].Name)
	}
}

func TestParseSessionNameWithSpaces(t *testing.T) {
	input := sessionLine("my session", "2", "1700000000", "1", "/home/user/project")
	sessions, err := ParseSessions(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Name != "my session" {
		t.Errorf("name: got %q, want %q", sessions[0].Name, "my session")
	}
}

func TestParseSessionNameUnicode(t *testing.T) {
	input := sessionLine("日本語", "1", "1700000000", "0", "/home/user/日本")
	sessions, err := ParseSessions(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Name != "日本語" {
		t.Errorf("name: got %q, want %q", sessions[0].Name, "日本語")
	}
	if sessions[0].Directory != "/home/user/日本" {
		t.Errorf("directory: got %q, want %q", sessions[0].Directory, "/home/user/日本")
	}
}

func TestParsePathWithPipes(t *testing.T) {
	// Paths containing '|' should now work since we use \x1f as delimiter
	input := sessionLine("dev", "1", "1700000000", "0", "/home/user/a|b|c")
	sessions, err := ParseSessions(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Directory != "/home/user/a|b|c" {
		t.Errorf("directory: got %q, want %q", sessions[0].Directory, "/home/user/a|b|c")
	}
}

func TestParseWindowsCountZero(t *testing.T) {
	input := sessionLine("empty", "0", "1700000000", "0", "/tmp")
	sessions, err := ParseSessions(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Windows != 0 {
		t.Errorf("windows: got %d, want 0", sessions[0].Windows)
	}
}

func TestParseVeryLongSessionName(t *testing.T) {
	longName := strings.Repeat("a", 200)
	input := sessionLine(longName, "1", "1700000000", "0", "/tmp")
	sessions, err := ParseSessions(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Name != longName {
		t.Errorf("name length: got %d, want %d", len(sessions[0].Name), len(longName))
	}
}

func TestValidateSessionName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid simple", "main", false, ""},
		{"valid with spaces", "my session", false, ""},
		{"valid unicode", "日本語", false, ""},
		{"valid with hyphens", "my-session", false, ""},
		{"empty", "", true, "must not be empty"},
		{"starts with dash", "-badname", true, "must not start with '-'"},
		{"contains dot", "bad.name", true, "must not contain '.' or ':'"},
		{"contains colon", "bad:name", true, "must not contain '.' or ':'"},
		{"too long", strings.Repeat("x", 129), true, "exceeds maximum length"},
		{"exactly max length", strings.Repeat("x", 128), false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSessionName(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
