package app

import (
	"fmt"
	"os"
	"os/exec"

	"agentmux/internal/adapters/tmux"
	"agentmux/internal/config"
	"agentmux/internal/core"
	"agentmux/internal/store"
	"agentmux/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

// App holds all initialized dependencies.
type App struct {
	Client     *tmux.Client
	Store      *store.Store
	Config     *config.Config
	Workspaces []core.Workspace
}

// New initializes the application.
func New() (*App, error) {
	client := tmux.New()

	cfg := config.LoadOrDefault(config.DefaultPath())
	workspaces := cfg.ToWorkspaces()

	st, err := store.Open(store.DefaultPath())
	if err != nil {
		// Non-fatal: run without persistence
		st = nil
	}

	return &App{
		Client:     client,
		Store:      st,
		Config:     cfg,
		Workspaces: workspaces,
	}, nil
}

// Close cleans up resources.
func (a *App) Close() {
	if a.Store != nil {
		a.Store.Close()
	}
}

// RunTUI launches the full-screen terminal UI.
func (a *App) RunTUI() error {
	if !a.Client.Available() {
		return fmt.Errorf("tmux is not installed or not in PATH\n\nInstall tmux:\n  macOS:  brew install tmux\n  Debian: sudo apt install tmux\n  Arch:   sudo pacman -S tmux")
	}

	model := tui.NewModel(a.Client, a.Store, a.Workspaces)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// If user selected a session to attach, do it after TUI exits
	if m, ok := finalModel.(tui.Model); ok {
		if target := m.AttachTarget(); target != "" {
			return execAttach(target)
		}
	}

	return nil
}

// execAttach replaces the process with tmux attach.
func execAttach(name string) error {
	bin, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}
	return Execvp(bin, []string{"tmux", "attach-session", "-t", name}, os.Environ())
}
