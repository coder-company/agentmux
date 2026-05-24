package tui

import (
	"agentmux/internal/adapters/tmux"
	"agentmux/internal/core"
	"agentmux/internal/store"
	"agentmux/internal/tui/views"

	tea "github.com/charmbracelet/bubbletea"
)

// View represents which view is active.
type View int

const (
	ViewSessions View = iota
	ViewPalette
	ViewLauncher
)

// Model is the top-level Bubble Tea model.
type Model struct {
	view         View
	sessions     *views.SessionsView
	palette      *views.PaletteView
	launcher     *views.LauncherView
	client       *tmux.Client
	store        *store.Store
	width        int
	height       int
	quitting     bool
	attachName   string
	confirmKill  bool
	confirmName  string
	renaming     bool
	renameBuffer string
}

// RefreshMsg triggers a session list refresh.
type RefreshMsg struct{}

// NewModel creates the top-level app model.
func NewModel(client *tmux.Client, st *store.Store, workspaces []core.Workspace) Model {
	sessView := views.NewSessionsView(client)
	palette := views.NewPalette(defaultActions())
	launcher := views.NewLauncher(workspaces)

	return Model{
		view:     ViewSessions,
		sessions: sessView,
		palette:  palette,
		launcher: launcher,
		client:   client,
		store:    st,
	}
}

// AttachTarget returns the session to attach to after the TUI exits, if any.
func (m Model) AttachTarget() string {
	return m.attachName
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return RefreshMsg{} }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.sessions.Width = msg.Width
		m.sessions.Height = msg.Height
		m.palette.Width = msg.Width
		m.palette.Height = msg.Height
		m.launcher.Width = msg.Width
		m.launcher.Height = msg.Height

	case RefreshMsg:
		m.sessions.Refresh()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Kill confirmation mode
	if m.confirmKill {
		return m.handleConfirmKey(msg)
	}

	// Rename mode
	if m.renaming {
		return m.handleRenameKey(msg)
	}

	// Palette mode
	if m.view == ViewPalette {
		return m.handlePaletteKey(msg)
	}

	// Launcher mode
	if m.view == ViewLauncher {
		return m.handleLauncherKey(msg)
	}

	// Sessions mode
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "up", "ctrl+p":
		m.sessions.List.MoveUp()
		m.sessions.RefreshPreview()

	case "down", "ctrl+n":
		m.sessions.List.MoveDown()
		m.sessions.RefreshPreview()

	case "enter":
		if sel := m.sessions.List.Selected(); sel != nil {
			m.attachName = sel.Name
			if m.store != nil {
				m.store.RecordSession(sel.Name)
			}
			return m, tea.Quit
		}

	case "n":
		name := generateSessionName(m.sessions.List.Sessions)
		if err := m.client.NewSession(name, ""); err == nil {
			m.sessions.Status = "✓ Created: " + name
			m.sessions.Refresh()
		} else {
			m.sessions.Status = "✗ " + err.Error()
		}

	case "k":
		if sel := m.sessions.List.Selected(); sel != nil {
			m.confirmKill = true
			m.confirmName = sel.Name
			m.sessions.Status = "Kill session \"" + sel.Name + "\"? [y/N]"
		}

	case "r":
		if sel := m.sessions.List.Selected(); sel != nil {
			m.renaming = true
			m.renameBuffer = sel.Name
			m.sessions.Status = "Rename to: " + m.renameBuffer + "█"
		}

	case "/":
		m.view = ViewPalette
		m.palette.SetQuery("")

	case "p":
		m.view = ViewLauncher
	}

	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if err := m.client.KillSession(m.confirmName); err == nil {
			m.sessions.Status = "✓ Killed: " + m.confirmName
			if m.store != nil {
				m.store.RemoveSession(m.confirmName)
			}
			m.sessions.Refresh()
		} else {
			m.sessions.Status = "✗ " + err.Error()
		}
		m.confirmKill = false
		m.confirmName = ""
	default:
		m.confirmKill = false
		m.confirmName = ""
		m.sessions.Status = ""
	}
	return m, nil
}

func (m Model) handleRenameKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		sel := m.sessions.List.Selected()
		if sel != nil && m.renameBuffer != "" && m.renameBuffer != sel.Name {
			if err := tmux.ValidateSessionName(m.renameBuffer); err != nil {
				m.sessions.Status = "✗ " + err.Error()
			} else if err := m.client.RenameSession(sel.Name, m.renameBuffer); err != nil {
				m.sessions.Status = "✗ " + err.Error()
			} else {
				m.sessions.Status = "✓ Renamed: " + sel.Name + " → " + m.renameBuffer
				m.sessions.Refresh()
			}
		}
		m.renaming = false
		m.renameBuffer = ""
	case "esc":
		m.renaming = false
		m.renameBuffer = ""
		m.sessions.Status = ""
	case "backspace":
		if len(m.renameBuffer) > 0 {
			m.renameBuffer = m.renameBuffer[:len(m.renameBuffer)-1]
		}
		m.sessions.Status = "Rename to: " + m.renameBuffer + "█"
	default:
		ch := msg.String()
		if len(ch) == 1 {
			m.renameBuffer += ch
			m.sessions.Status = "Rename to: " + m.renameBuffer + "█"
		}
	}
	return m, nil
}

func (m Model) handlePaletteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.view = ViewSessions
	case "up", "ctrl+p":
		m.palette.MoveUp()
	case "down", "ctrl+n":
		m.palette.MoveDown()
	case "enter":
		if sel := m.palette.Selected(); sel != nil && sel.Do != nil {
			sel.Do()
			m.view = ViewSessions
			m.sessions.Refresh()
		}
	case "backspace":
		m.palette.Backspace()
	default:
		if len(msg.String()) == 1 {
			m.palette.TypeChar(rune(msg.String()[0]))
		}
	}
	return m, nil
}

func (m Model) handleLauncherKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.view = ViewSessions
	case "up", "ctrl+p":
		m.launcher.MoveUp()
	case "down", "ctrl+n":
		m.launcher.MoveDown()
	case "tab":
		m.launcher.ToggleCommands()
	case "enter":
		ws := m.launcher.SelectedWorkspace()
		if ws != nil {
			if err := m.client.NewSession(ws.Name, ws.Root); err == nil {
				m.sessions.Status = "✓ Launched: " + ws.Name
				m.view = ViewSessions
				m.sessions.Refresh()
			} else {
				m.sessions.Status = "✗ " + err.Error()
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	switch m.view {
	case ViewPalette:
		return m.palette.Render()
	case ViewLauncher:
		return m.launcher.Render()
	default:
		return m.sessions.Render()
	}
}

func defaultActions() []views.Action {
	return []views.Action{
		{Name: "New Session", Desc: "Create a new tmux session", Key: "n"},
		{Name: "Kill Session", Desc: "Kill the selected session", Key: "k"},
		{Name: "Rename Session", Desc: "Rename the selected session", Key: "r"},
		{Name: "Workspaces", Desc: "Open workspace launcher", Key: "p"},
		{Name: "Refresh", Desc: "Reload session list", Key: ""},
		{Name: "Quit", Desc: "Exit agentmux", Key: "q"},
	}
}

func generateSessionName(existing []core.Session) string {
	base := "s"
	for i := len(existing); ; i++ {
		name := base + itoa(i)
		found := false
		for _, s := range existing {
			if s.Name == name {
				found = true
				break
			}
		}
		if !found {
			return name
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
