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
	sessions, err := v.Client.ListSessions()
	if err != nil {
		v.Status = "error: " + err.Error()
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
		v.Preview = components.Preview{Title: "Preview", Content: ""}
		return
	}
	content, err := v.Client.CapturePane(sel.Name, 30)
	if err != nil {
		content = "capture error: " + err.Error()
	}
	v.Preview = components.Preview{
		Title:   sel.Name,
		Content: content,
	}
}

// Render returns the full sessions view.
func (v *SessionsView) Render() string {
	leftWidth := v.Width / 3
	if leftWidth < 24 {
		leftWidth = 24
	}
	rightWidth := v.Width - leftWidth - 4

	// Left panel: session list
	listContent := styles.Title.Render("Sessions") + "\n" + v.List.Render(leftWidth-2)
	left := styles.ActivePanel.Width(leftWidth).Height(v.Height - 4).Render(listContent)

	// Right panel: preview
	previewContent := v.Preview.Render(rightWidth-2, v.Height-4)
	right := styles.Panel.Width(rightWidth).Height(v.Height - 4).Render(previewContent)

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
