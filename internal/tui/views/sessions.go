package views

import (
	"agentmux/internal/adapters/tmux"
	"agentmux/internal/tui/components"
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SessionsView is the main sessions browser.
type SessionsView struct {
	List    components.SessionList
	Preview components.Preview
	Help    *components.HelpBar
	Client  *tmux.Client
	Width   int
	Height  int
	Status  string
	Loaded  bool
}

// NewSessionsView creates the sessions view.
func NewSessionsView(client *tmux.Client) *SessionsView {
	return &SessionsView{
		Client: client,
		Help:   components.DefaultHelp(),
	}
}

// Refresh reloads sessions from tmux.
func (v *SessionsView) Refresh() {
	v.Loaded = true
	sessions, err := v.Client.ListSessions()
	if err != nil {
		v.Status = "✗ " + err.Error()
		return
	}
	v.List.Sessions = sessions
	if v.List.Cursor >= len(sessions) && len(sessions) > 0 {
		v.List.Cursor = len(sessions) - 1
	}
	v.RefreshPreview()
}

// RefreshPreview loads the pane capture for the selected session.
func (v *SessionsView) RefreshPreview() {
	sel := v.List.Selected()
	if sel == nil {
		v.Preview = components.Preview{Title: "Preview"}
		return
	}
	content, err := v.Client.CapturePane(sel.Name, 30)
	if err != nil {
		v.Preview = components.Preview{
			Title: sel.Name,
			Error: "Could not capture pane output",
		}
		return
	}
	v.Preview = components.Preview{
		Title:   sel.Name,
		Content: content,
	}
}

// Render returns the full sessions view.
func (v *SessionsView) Render() string {
	if v.Width == 0 || v.Height == 0 {
		return "Loading..."
	}

	leftWidth := v.Width / 3
	if leftWidth < 28 {
		leftWidth = 28
	}
	if leftWidth > v.Width-20 {
		leftWidth = v.Width - 20
	}
	rightWidth := v.Width - leftWidth - 4
	if rightWidth < 10 {
		rightWidth = 10
	}

	panelHeight := v.Height - 4
	if panelHeight < 3 {
		panelHeight = 3
	}

	// Left panel: session list
	listTitle := styles.Title.Render("Sessions")
	if !v.Loaded {
		listTitle += " " + styles.Muted.Render("(loading…)")
	}
	listContent := listTitle + "\n" + v.List.Render(leftWidth-2)
	left := styles.ActivePanel.Width(leftWidth).Height(panelHeight).Render(listContent)

	// Right panel: preview
	previewContent := v.Preview.Render(rightWidth-2, panelHeight)
	right := styles.Panel.Width(rightWidth).Height(panelHeight).Render(previewContent)

	main := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	// Status line
	status := ""
	if v.Status != "" {
		status = styles.Muted.Render(v.Status)
	}

	// Help bar
	help := v.Help.Render(v.Width)

	return lipgloss.JoinVertical(lipgloss.Left, main, status, help)
}
