package views

import (
	"agentmux/internal/adapters/tmux"
	"agentmux/internal/tui/components"
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SessionsView is the main sessions browser.
type SessionsView struct {
	List   components.SessionList
	Pane   components.Preview
	Client *tmux.Client
	Width  int
	Height int
	Status string
	Loaded bool
}

// NewSessionsView creates the sessions view.
func NewSessionsView(client *tmux.Client) *SessionsView {
	return &SessionsView{Client: client}
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
	v.Status = ""
}

// RefreshPreview loads the pane capture for the selected session.
func (v *SessionsView) RefreshPreview() {
	sel := v.List.Selected()
	if sel == nil {
		v.Pane = components.Preview{}
		return
	}
	content, err := v.Client.CapturePane(sel.Name, 40)
	if err != nil {
		v.Pane = components.Preview{
			Title: sel.Name,
			Dir:   sel.Directory,
			Error: "Could not capture pane",
		}
		return
	}
	v.Pane = components.Preview{
		Title:   sel.Name,
		Dir:     sel.Directory,
		Content: content,
	}
}

// Render returns the full sessions view.
func (v *SessionsView) Render(bodyHeight int) string {
	if v.Width < 20 || bodyHeight < 4 {
		return styles.Muted.Render("Terminal too small")
	}

	// Calculate panel widths
	leftWidth := v.Width * 30 / 100
	if leftWidth < 24 {
		leftWidth = 24
	}
	if leftWidth > 50 {
		leftWidth = 50
	}
	rightWidth := v.Width - leftWidth - 4 // borders eat 4 chars
	if rightWidth < 12 {
		rightWidth = 12
	}

	panelHeight := bodyHeight
	if panelHeight < 3 {
		panelHeight = 3
	}
	innerHeight := panelHeight - 2 // border top + bottom

	// Left panel: session list
	leftTitle := styles.PanelTitle.Render("Sessions")
	if !v.Loaded {
		leftTitle += " " + styles.Muted.Render("…")
	}
	listContent := leftTitle + "\n" + v.List.Render(leftWidth-4, innerHeight-1)
	left := styles.PanelActive.Width(leftWidth).Height(panelHeight).Render(listContent)

	// Right panel: preview
	previewContent := v.Pane.Render(rightWidth-4, innerHeight)
	right := styles.Panel.Width(rightWidth).Height(panelHeight).Render(previewContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
