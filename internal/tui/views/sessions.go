package views

import (
	"fmt"
	"strings"

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
	if v.Width < 50 || bodyHeight < 7 {
		return styles.Muted.Render(" Terminal too small. Resize to continue.")
	}

	if v.Width < 88 {
		return v.renderCompact(bodyHeight)
	}

	leftOuter := v.Width * 36 / 100
	if leftOuter < 32 {
		leftOuter = 32
	}
	if leftOuter > 48 {
		leftOuter = 48
	}
	rightOuter := v.Width - leftOuter - 1

	panelH := bodyHeight
	innerH := panelH - 2
	if innerH < 1 {
		innerH = 1
	}

	leftInner := leftOuter - 2
	leftHeader := panelHeader("sessions", sessionCountLabel(len(v.List.Sessions)), leftInner)
	listContent := v.List.Render(leftInner, innerH-2)
	leftContent := leftHeader + "\n" + styles.Subtle.Render(rule(leftInner)) + "\n" + listContent
	left := styles.PanelBorderActive.
		Width(leftOuter).
		Height(panelH).
		Render(leftContent)

	rightInner := rightOuter - 2
	selected := "no selection"
	if sel := v.List.Selected(); sel != nil {
		selected = "tail 40 lines"
	}
	rightHeader := panelHeader("preview", selected, rightInner)
	previewContent := v.Pane.Render(rightInner, innerH-2)
	rightContent := rightHeader + "\n" + styles.Subtle.Render(rule(rightInner)) + "\n" + previewContent
	right := styles.PanelBorder.
		Width(rightOuter).
		Height(panelH).
		Render(rightContent)

	gap := lipgloss.NewStyle().Width(1).Render(" ")
	return lipgloss.JoinHorizontal(lipgloss.Top, left, gap, right)
}

func (v *SessionsView) renderCompact(bodyHeight int) string {
	panelH := bodyHeight
	innerW := v.Width - 2
	innerH := panelH - 2
	header := panelHeader("sessions", "preview hidden on narrow terminals", innerW)
	listContent := v.List.Render(innerW, innerH-2)
	content := header + "\n" + styles.Subtle.Render(rule(innerW)) + "\n" + listContent

	return styles.PanelBorderActive.
		Width(v.Width).
		Height(panelH).
		Render(content)
}

func panelHeader(title, meta string, width int) string {
	title = styles.PanelTitle.Render(title)
	meta = styles.PanelMeta.Render(styles.Truncate(meta, width/2))
	gap := width - lipgloss.Width(title) - lipgloss.Width(meta)
	if gap < 1 {
		return styles.Truncate(title, width)
	}
	return title + lipgloss.NewStyle().Width(gap).Render("") + meta
}

func sessionCountLabel(n int) string {
	if n == 1 {
		return "1 session"
	}
	return fmt.Sprintf("%d sessions", n)
}

func rule(width int) string {
	if width <= 0 {
		return ""
	}
	return strings.Repeat("─", width)
}
