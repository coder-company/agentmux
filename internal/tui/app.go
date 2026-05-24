package tui

import (
	"fmt"

	"agentmux/internal/adapters/tmux"
	"agentmux/internal/core"
	"agentmux/internal/store"
	"agentmux/internal/tui/components"
	"agentmux/internal/tui/styles"
	"agentmux/internal/tui/views"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View represents which view is active.
type View int

const (
	ViewSessions View = iota
	ViewPalette
	ViewLauncher
	ViewHelp
)

// Model is the top-level Bubble Tea model.
type Model struct {
	view         View
	sessions     *views.SessionsView
	palette      *views.PaletteView
	launcher     *views.LauncherView
	help         *views.HelpOverlay
	header       components.Header
	client       *tmux.Client
	store        *store.Store
	workspaces   []core.Workspace
	width        int
	height       int
	quitting     bool
	attachName   string
	confirmKill  bool
	confirmName  string
	renaming     bool
	renameBuffer string
	prevG        bool // for gg binding
}

// RefreshMsg triggers a session list refresh.
type RefreshMsg struct{}

// NewModel creates the top-level app model.
func NewModel(client *tmux.Client, st *store.Store, workspaces []core.Workspace) Model {
	sessView := views.NewSessionsView(client)
	launcher := views.NewLauncher(workspaces)

	return Model{
		view:       ViewSessions,
		sessions:   sessView,
		palette:    views.NewPalette(nil), // populated after first refresh
		launcher:   launcher,
		help:       &views.HelpOverlay{},
		client:     client,
		store:      st,
		workspaces: workspaces,
	}
}

// AttachTarget returns the session to attach to after the TUI exits.
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
		m.help.Width = msg.Width
		m.help.Height = msg.Height
		m.header.Width = msg.Width

	case RefreshMsg:
		m.sessions.Refresh()
		m.header.SessionCount = len(m.sessions.List.Sessions)
		m.rebuildPalette()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Kill confirmation mode — highest priority
	if m.confirmKill {
		return m.handleConfirmKey(msg)
	}

	// Rename mode
	if m.renaming {
		return m.handleRenameKey(msg)
	}

	// Help overlay
	if m.view == ViewHelp {
		switch msg.String() {
		case "esc", "?", "q":
			m.view = ViewSessions
		}
		return m, nil
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
	return m.handleSessionsKey(msg)
}

func (m Model) handleSessionsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// gg binding
	if m.prevG {
		m.prevG = false
		if key == "g" {
			m.sessions.List.MoveTop()
			m.sessions.RefreshPreview()
			return m, nil
		}
	}

	switch key {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		m.sessions.List.MoveUp()
		m.sessions.RefreshPreview()

	case "down", "j":
		m.sessions.List.MoveDown()
		m.sessions.RefreshPreview()

	case "g":
		m.prevG = true

	case "G":
		m.sessions.List.MoveBottom()
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
			m.sessions.Status = "✓ Created " + name
			m.doRefresh()
		} else {
			m.sessions.Status = "✗ " + err.Error()
		}

	case "x":
		if sel := m.sessions.List.Selected(); sel != nil {
			m.confirmKill = true
			m.confirmName = sel.Name
			m.sessions.Status = fmt.Sprintf("Kill \"%s\"? y/n", sel.Name)
		}

	case "r":
		if sel := m.sessions.List.Selected(); sel != nil {
			m.renaming = true
			m.renameBuffer = sel.Name
			m.sessions.Status = "Rename: " + m.renameBuffer + "│"
		}

	case "R":
		m.doRefresh()
		m.sessions.Status = "✓ Refreshed"

	case "/":
		m.rebuildPalette()
		m.view = ViewPalette
		m.palette.SetQuery("")

	case "p":
		m.view = ViewLauncher

	case "?":
		m.view = ViewHelp
	}

	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if err := m.client.KillSession(m.confirmName); err == nil {
			m.sessions.Status = "✓ Killed " + m.confirmName
			if m.store != nil {
				m.store.RemoveSession(m.confirmName)
			}
			m.doRefresh()
		} else {
			m.sessions.Status = "✗ " + err.Error()
		}
	default:
		m.sessions.Status = ""
	}
	m.confirmKill = false
	m.confirmName = ""
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
				m.sessions.Status = "✓ Renamed → " + m.renameBuffer
				m.doRefresh()
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
		m.sessions.Status = "Rename: " + m.renameBuffer + "│"
	default:
		ch := msg.String()
		if len(ch) == 1 {
			m.renameBuffer += ch
			m.sessions.Status = "Rename: " + m.renameBuffer + "│"
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
		if sel := m.palette.Selected(); sel != nil {
			m.view = ViewSessions
			if sel.Do != nil {
				sel.Do()
			}
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
	case "up", "k":
		m.launcher.MoveUp()
	case "down", "j":
		m.launcher.MoveDown()
	case "tab":
		m.launcher.ToggleCommands()
	case "enter":
		ws := m.launcher.SelectedWorkspace()
		if ws == nil {
			break
		}
		// If in commands mode and a command is selected, use workspace-command naming
		sessionName := ws.Name
		if m.launcher.InCommands {
			if cmd := m.launcher.SelectedCommand(); cmd != nil {
				sessionName = ws.Name + "-" + cmd.Name
			}
		}
		if err := m.client.NewSession(sessionName, ws.Root); err == nil {
			m.sessions.Status = "✓ Launched " + sessionName
			m.view = ViewSessions
			m.doRefresh()
		} else {
			m.sessions.Status = "✗ " + err.Error()
			m.view = ViewSessions
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.width == 0 || m.height == 0 {
		return "Loading…"
	}

	// Overlays render on top of everything
	switch m.view {
	case ViewPalette:
		return m.palette.Render()
	case ViewLauncher:
		return m.launcher.Render()
	case ViewHelp:
		return m.help.Render()
	}

	// Main layout: header + body + status + footer
	headerStr := m.header.Render()
	headerH := lipgloss.Height(headerStr)

	footerStr := components.SessionsFooter(m.width)
	footerH := lipgloss.Height(footerStr)

	statusStr := ""
	statusH := 0
	if m.sessions.Status != "" {
		statusStr = styles.StatusBar.Width(m.width).Render(m.sessions.Status)
		statusH = 1
	}

	bodyHeight := m.height - headerH - footerH - statusH
	if bodyHeight < 3 {
		bodyHeight = 3
	}

	bodyStr := m.sessions.Render(bodyHeight)

	parts := []string{headerStr, bodyStr}
	if statusStr != "" {
		parts = append(parts, statusStr)
	}
	parts = append(parts, footerStr)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// doRefresh refreshes sessions and updates header count.
func (m *Model) doRefresh() {
	m.sessions.Refresh()
	m.header.SessionCount = len(m.sessions.List.Sessions)
}

// rebuildPalette builds the command palette actions from current state.
func (m *Model) rebuildPalette() {
	actions := []views.Action{
		{Name: "New Session", Desc: "Create a detached session", Key: "n", Do: func() {
			name := generateSessionName(m.sessions.List.Sessions)
			if err := m.client.NewSession(name, ""); err == nil {
				m.sessions.Status = "✓ Created " + name
				m.doRefresh()
			}
		}},
		{Name: "Kill Session", Desc: "Destroy selected session", Key: "x", Do: func() {
			if sel := m.sessions.List.Selected(); sel != nil {
				m.confirmKill = true
				m.confirmName = sel.Name
				m.sessions.Status = fmt.Sprintf("Kill \"%s\"? y/n", sel.Name)
			}
		}},
		{Name: "Rename Session", Desc: "Rename selected session", Key: "r", Do: func() {
			if sel := m.sessions.List.Selected(); sel != nil {
				m.renaming = true
				m.renameBuffer = sel.Name
				m.sessions.Status = "Rename: " + m.renameBuffer + "│"
			}
		}},
		{Name: "Refresh", Desc: "Reload session list", Key: "R", Do: func() {
			m.doRefresh()
			m.sessions.Status = "✓ Refreshed"
		}},
		{Name: "Workspaces", Desc: "Open workspace launcher", Key: "p", Do: func() {
			m.view = ViewLauncher
		}},
		{Name: "Help", Desc: "Show keybindings", Key: "?", Do: func() {
			m.view = ViewHelp
		}},
		{Name: "Quit", Desc: "Exit agentmux", Key: "q", Do: func() {
			m.quitting = true
		}},
	}

	// Add workspace launchers
	for _, ws := range m.workspaces {
		ws := ws
		actions = append(actions, views.Action{
			Name: "Launch: " + ws.Name,
			Desc: ws.Root,
			Do: func() {
				if err := m.client.NewSession(ws.Name, ws.Root); err == nil {
					m.sessions.Status = "✓ Launched " + ws.Name
					m.doRefresh()
				}
			},
		})
	}

	// Add session-specific actions
	for _, s := range m.sessions.List.Sessions {
		s := s
		actions = append(actions, views.Action{
			Name: "Attach: " + s.Name,
			Desc: s.Directory,
			Do: func() {
				m.attachName = s.Name
				if m.store != nil {
					m.store.RecordSession(s.Name)
				}
				m.quitting = true
			},
		})
	}

	m.palette.UpdateActions(actions)
}

func generateSessionName(existing []core.Session) string {
	for i := len(existing); ; i++ {
		name := fmt.Sprintf("s%d", i)
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
