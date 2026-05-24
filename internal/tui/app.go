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
	keys       KeyMap
	view       View
	sessions   *views.SessionsView
	palette    *views.PaletteView
	launcher   *views.LauncherView
	client     *tmux.Client
	store      *store.Store
	width      int
	height     int
	quitting   bool
	attachName string
}

// AttachMsg signals the app should quit and attach to a session.
type AttachMsg struct{ Name string }

// RefreshMsg triggers a session list refresh.
type RefreshMsg struct{}

// NewModel creates the top-level app model.
func NewModel(client *tmux.Client, st *store.Store, workspaces []core.Workspace) Model {
	sessView := views.NewSessionsView(client)
	palette := views.NewPalette(defaultActions())
	launcher := views.NewLauncher(workspaces)

	return Model{
		keys:     DefaultKeys(),
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
		// Create session with auto-generated name
		name := generateSessionName(m.sessions.List.Sessions)
		if err := m.client.NewSession(name, ""); err == nil {
			m.sessions.Status = "created: " + name
			m.sessions.Refresh()
		} else {
			m.sessions.Status = "error: " + err.Error()
		}

	case "k":
		if sel := m.sessions.List.Selected(); sel != nil {
			if err := m.client.KillSession(sel.Name); err == nil {
				m.sessions.Status = "killed: " + sel.Name
				if m.store != nil {
					m.store.RemoveSession(sel.Name)
				}
				m.sessions.Refresh()
			} else {
				m.sessions.Status = "error: " + err.Error()
			}
		}

	case "/":
		m.view = ViewPalette
		m.palette.SetQuery("")

	case "p":
		m.view = ViewLauncher
	}

	return m, nil
}

func (m Model) handlePaletteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
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
				m.sessions.Status = "launched: " + ws.Name
				m.view = ViewSessions
				m.sessions.Refresh()
			} else {
				m.sessions.Status = "error: " + err.Error()
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
