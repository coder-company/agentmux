package views

import (
	"testing"

	"agentmux/internal/core"

	"github.com/charmbracelet/lipgloss"
)

func TestLauncherToggleCommandsRequiresCommands(t *testing.T) {
	l := NewLauncher([]core.Workspace{{Name: "api", Root: "/tmp/api"}})

	l.ToggleCommands()
	if l.InCommands {
		t.Fatal("expected command mode to stay disabled without commands")
	}
}

func TestLauncherRenderFitsScreenWidth(t *testing.T) {
	l := NewLauncher([]core.Workspace{{
		Name: "api-service-with-a-long-name",
		Root: "/very/long/workspace/root/that/should/be/truncated",
		Commands: []core.Command{
			{Name: "server", Cmd: "go run ./cmd/server --with-a-long-flag"},
		},
	}})
	l.Width = 70
	l.Height = 22

	out := l.Render()
	if width := lipgloss.Width(out); width > 70 {
		t.Fatalf("rendered width = %d, want <= 70", width)
	}
}
