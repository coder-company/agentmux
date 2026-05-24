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

	var windowNames []string
	if wins, err := v.Client.ListWindows(sel.Name); err == nil {
		windowNames = wins
	}

	content, err := v.Client.CapturePane(sel.Name, 40)
	if err != nil {
		v.Pane = components.Preview{
			Title:   sel.Name,
			Dir:     sel.Directory,
			Windows: windowNames,
			Error:   "could not capture pane",
		}
		return
	}
	v.Pane = components.Preview{
		Title:   sel.Name,
		Dir:     sel.Directory,
		Windows: windowNames,
		Content: content,
	}
}

// Render returns the full two-panel sessions view.
func (v *SessionsView) Render(bodyHeight int) string {
	if v.Width < 40 || bodyHeight < 5 {
		return styles.Muted.Render(" Terminal too small. Resize to continue.")
	}

	// Panel sizing: 35% left, rest right
	leftOuter := v.Width * 35 / 100
	if leftOuter < 28 {
		leftOuter = 28
	}
	if leftOuter > 50 {
		leftOuter = 50
	}
	rightOuter := v.Width - leftOuter - 1 // 1 char gap between panels

	panelH := bodyHeight
	innerH := panelH - 2 // top + bottom border
	if innerH < 1 {
		innerH = 1
	}

	// Left panel: sessions
	leftInner := leftOuter - 4 // border(2) + padding(2)
	listContent := v.List.Render(leftInner, innerH)
	left := styles.PanelBorderActive.
		Width(leftOuter).
		Height(panelH).
		Render(listContent)

	// Right panel: preview
	rightInner := rightOuter - 4
	previewContent := v.Pane.Render(rightInner, innerH)
	right := styles.PanelBorder.
		Width(rightOuter).
		Height(panelH).
		Render(previewContent)

	gap := lipgloss.NewStyle().Width(1).Render(" ")
	return lipgloss.JoinHorizontal(lipgloss.Top, left, gap, right)
}
