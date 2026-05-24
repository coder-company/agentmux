package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testConfig = `
[[workspaces]]
name = "anode"
root = "~/code/anode"

[[workspaces.commands]]
name = "run"
cmd = "go run ./cmd/anode"

[[workspaces.commands]]
name = "test"
cmd = "go test ./..."

[[workspaces]]
name = "web"
root = "/srv/web"
`

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(testConfig), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Workspaces) != 2 {
		t.Fatalf("expected 2 workspaces, got %d", len(cfg.Workspaces))
	}

	ws := cfg.ToWorkspaces()
	if ws[0].Name != "anode" {
		t.Errorf("expected name 'anode', got %q", ws[0].Name)
	}
	if len(ws[0].Commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(ws[0].Commands))
	}
	if ws[0].Commands[0].Cmd != "go run ./cmd/anode" {
		t.Errorf("unexpected cmd: %q", ws[0].Commands[0].Cmd)
	}
	if ws[1].Root != "/srv/web" {
		t.Errorf("expected root '/srv/web', got %q", ws[1].Root)
	}
}

func TestLoadOrDefaultMissing(t *testing.T) {
	cfg := LoadOrDefault("/nonexistent/path.toml")
	if len(cfg.Workspaces) != 0 {
		t.Errorf("expected empty workspaces for missing config")
	}
}

func TestExpandHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	got := expandHome("~/projects")
	want := filepath.Join(home, "projects")
	if got != want {
		t.Errorf("expandHome: got %q, want %q", got, want)
	}

	abs := "/absolute/path"
	if expandHome(abs) != abs {
		t.Errorf("expandHome should not modify absolute paths")
	}
}

func TestExpandHomeBare(t *testing.T) {
	home, _ := os.UserHomeDir()
	got := expandHome("~")
	if got != home {
		t.Errorf("expandHome(~): got %q, want %q", got, home)
	}
}

func TestValidateMissingName(t *testing.T) {
	cfg := &Config{Workspaces: []workspaceEntry{{Name: "", Root: "/tmp"}}}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "missing name") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateMissingRoot(t *testing.T) {
	cfg := &Config{Workspaces: []workspaceEntry{{Name: "foo", Root: ""}}}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing root")
	}
	if !strings.Contains(err.Error(), `workspace "foo": missing root directory`) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDuplicateWorkspace(t *testing.T) {
	cfg := &Config{Workspaces: []workspaceEntry{
		{Name: "dup", Root: "/a"},
		{Name: "dup", Root: "/b"},
	}}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for duplicate workspace")
	}
	if !strings.Contains(err.Error(), "duplicate name") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDuplicateCommand(t *testing.T) {
	cfg := &Config{Workspaces: []workspaceEntry{
		{Name: "ws", Root: "/a", Commands: []commandEntry{
			{Name: "run", Cmd: "echo 1"},
			{Name: "run", Cmd: "echo 2"},
		}},
	}}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for duplicate command")
	}
	if !strings.Contains(err.Error(), `duplicate command "run"`) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateEmptyCommand(t *testing.T) {
	cfg := &Config{Workspaces: []workspaceEntry{
		{Name: "ws", Root: "/a", Commands: []commandEntry{
			{Name: "bad", Cmd: ""},
		}},
	}}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for empty cmd")
	}
	if !strings.Contains(err.Error(), "empty cmd") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInitDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "config.toml")

	// Should create the file and parent directory.
	if err := InitDefault(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "agentmux") {
		t.Error("expected example content in default config")
	}

	// Write custom content, then call again — should not overwrite.
	custom := []byte("custom")
	if err := os.WriteFile(path, custom, 0644); err != nil {
		t.Fatal(err)
	}
	if err := InitDefault(path); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != "custom" {
		t.Error("InitDefault overwrote existing file")
	}
}
