package views

import (
	"fmt"
	"sort"
	"strings"

	"agentmux/internal/adapters/tmux"
	"agentmux/internal/core"
	"agentmux/internal/tui/components"
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SessionSort controls how the session list is ordered.
type SessionSort int

const (
	SortCreated SessionSort = iota
	SortName
	SortWindows
	SortAttached
)

// SessionLayout controls the main browser presentation.
type SessionLayout int

const (
	LayoutSplit SessionLayout = iota
	LayoutList
	LayoutPreview
)

// SessionsView is the main sessions browser.
type SessionsView struct {
	List        components.SessionList
	Pane        components.Preview
	Client      *tmux.Client
	AllSessions []core.Session
	Filter      string
	Sort        SessionSort
	Layout      SessionLayout
	Width       int
	Height      int
	Status      string
	Loaded      bool
}

// NewSessionsView creates the sessions view.
func NewSessionsView(client *tmux.Client) *SessionsView {
	return &SessionsView{Client: client}
}

// Refresh reloads sessions from tmux.
func (v *SessionsView) Refresh() {
	v.Loaded = true
	selected := v.selectedName()
	sessions, err := v.Client.ListSessions()
	if err != nil {
		v.Status = "✗ " + err.Error()
		return
	}
	v.AllSessions = sessions
	v.applySessions(selected)
	v.RefreshPreview()
	v.Status = ""
}

// SetFilter updates the visible session filter.
func (v *SessionsView) SetFilter(query string) {
	v.Filter = query
	v.applySessions(v.selectedName())
	v.RefreshPreview()
}

// ClearFilter removes the active filter.
func (v *SessionsView) ClearFilter() {
	v.SetFilter("")
}

// CycleSort advances to the next sort mode.
func (v *SessionsView) CycleSort() {
	v.Sort = (v.Sort + 1) % 4
	v.applySessions(v.selectedName())
	v.RefreshPreview()
}

// CycleLayout advances to the next browser layout.
func (v *SessionsView) CycleLayout() {
	v.Layout = (v.Layout + 1) % 3
}

// SortLabel returns the user-facing sort label.
func (v *SessionsView) SortLabel() string {
	switch v.Sort {
	case SortName:
		return "name"
	case SortWindows:
		return "windows"
	case SortAttached:
		return "attached"
	default:
		return "newest"
	}
}

// LayoutLabel returns the user-facing layout label.
func (v *SessionsView) LayoutLabel() string {
	switch v.Layout {
	case LayoutList:
		return "list"
	case LayoutPreview:
		return "preview"
	default:
		return "split"
	}
}

// VisibleCount returns the number of sessions currently shown.
func (v *SessionsView) VisibleCount() int {
	return len(v.List.Sessions)
}

// TotalCount returns the number of known tmux sessions.
func (v *SessionsView) TotalCount() int {
	return len(v.AllSessions)
}

func (v *SessionsView) applySessions(selected string) {
	filter := strings.ToLower(strings.TrimSpace(v.Filter))
	visible := make([]core.Session, 0, len(v.AllSessions))
	for _, session := range v.AllSessions {
		if filter == "" || sessionMatches(session, filter) {
			visible = append(visible, session)
		}
	}

	sortSessions(visible, v.Sort)
	v.List.Sessions = visible
	v.configureEmptyState()
	v.restoreSelection(selected)
}

func (v *SessionsView) configureEmptyState() {
	if len(v.AllSessions) == 0 {
		v.List.EmptyTitle = ""
		v.List.EmptyBody = ""
		v.List.EmptyHints = nil
		return
	}
	if strings.TrimSpace(v.Filter) == "" {
		v.List.EmptyTitle = ""
		v.List.EmptyBody = ""
		v.List.EmptyHints = nil
		return
	}

	v.List.EmptyTitle = "No matching sessions"
	v.List.EmptyBody = fmt.Sprintf("%q hid %d sessions.", v.Filter, len(v.AllSessions))
	v.List.EmptyHints = []components.HelpBinding{
		{Key: "/", Desc: "edit filter"},
		{Key: "esc", Desc: "clear filter"},
		{Key: ":", Desc: "command palette"},
	}
}

func (v *SessionsView) restoreSelection(selected string) {
	if len(v.List.Sessions) == 0 {
		v.List.Cursor = 0
		return
	}
	if selected != "" {
		for i, session := range v.List.Sessions {
			if session.Name == selected {
				v.List.Cursor = i
				return
			}
		}
	}
	if v.List.Cursor >= len(v.List.Sessions) {
		v.List.Cursor = len(v.List.Sessions) - 1
	}
	if v.List.Cursor < 0 {
		v.List.Cursor = 0
	}
}

func (v *SessionsView) selectedName() string {
	if sel := v.List.Selected(); sel != nil {
		return sel.Name
	}
	return ""
}

// RefreshPreview loads the pane capture for the selected session.
func (v *SessionsView) RefreshPreview() {
	sel := v.List.Selected()
	if sel == nil {
		v.Pane = components.Preview{}
		return
	}

	if v.Client == nil {
		v.Pane = components.Preview{
			Title:       sel.Name,
			Dir:         sel.Directory,
			Created:     sel.Created,
			Attached:    sel.Attached,
			WindowCount: sel.Windows,
		}
		return
	}

	var windowNames []string
	if wins, err := v.Client.ListWindows(sel.Name); err == nil {
		windowNames = wins
	}

	content, err := v.Client.CapturePane(sel.Name, 40)
	if err != nil {
		v.Pane = components.Preview{
			Title:       sel.Name,
			Dir:         sel.Directory,
			Windows:     windowNames,
			Created:     sel.Created,
			Attached:    sel.Attached,
			WindowCount: sel.Windows,
			Error:       "could not capture pane",
		}
		return
	}
	v.Pane = components.Preview{
		Title:       sel.Name,
		Dir:         sel.Directory,
		Windows:     windowNames,
		Created:     sel.Created,
		Attached:    sel.Attached,
		WindowCount: sel.Windows,
		Content:     content,
	}
}

// Render returns the full two-panel sessions view.
func (v *SessionsView) Render(bodyHeight int) string {
	if v.Width < 50 || bodyHeight < 7 {
		return styles.Muted.Render(" Terminal too small. Resize to continue.")
	}

	if v.Layout == LayoutPreview {
		return v.renderPreviewOnly(bodyHeight)
	}
	if v.Layout == LayoutList || v.Width < 88 {
		return v.renderListOnly(bodyHeight)
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
	leftHeader := panelHeader("sessions", v.listMeta(), leftInner)
	listContent := v.List.Render(leftInner, innerH-2)
	leftContent := leftHeader + "\n" + styles.Subtle.Render(rule(leftInner)) + "\n" + listContent
	left := styles.PanelBorderActive.
		Width(leftOuter).
		Height(panelH).
		Render(leftContent)

	rightInner := rightOuter - 2
	selected := "no selection"
	if sel := v.List.Selected(); sel != nil {
		selected = fmt.Sprintf("%dw · tail 40", sel.Windows)
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

func (v *SessionsView) renderListOnly(bodyHeight int) string {
	panelH := bodyHeight
	innerW := v.Width - 2
	innerH := panelH - 2
	header := panelHeader("sessions", v.listMeta(), innerW)
	listContent := v.List.Render(innerW, innerH-2)
	content := header + "\n" + styles.Subtle.Render(rule(innerW)) + "\n" + listContent

	return styles.PanelBorderActive.
		Width(v.Width).
		Height(panelH).
		Render(content)
}

func (v *SessionsView) renderPreviewOnly(bodyHeight int) string {
	panelH := bodyHeight
	innerW := v.Width - 2
	innerH := panelH - 2
	meta := "no selection"
	if sel := v.List.Selected(); sel != nil {
		meta = fmt.Sprintf("%s · %dw · tail 40", styles.Truncate(sel.Name, 24), sel.Windows)
	}
	header := panelHeader("preview", meta, innerW)
	previewContent := v.Pane.Render(innerW, innerH-2)
	content := header + "\n" + styles.Subtle.Render(rule(innerW)) + "\n" + previewContent

	return styles.PanelBorderActive.
		Width(v.Width).
		Height(panelH).
		Render(content)
}

func (v *SessionsView) listMeta() string {
	count := sessionCountLabel(v.VisibleCount(), v.TotalCount())
	parts := []string{count, "sort " + v.SortLabel(), "layout " + v.LayoutLabel()}
	if strings.TrimSpace(v.Filter) != "" {
		parts = append([]string{fmt.Sprintf("filter %q", v.Filter)}, parts...)
	}
	return strings.Join(parts, " · ")
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

func sessionCountLabel(visible, total int) string {
	if total == 0 {
		return "0 sessions"
	}
	if visible != total {
		return fmt.Sprintf("%d/%d sessions", visible, total)
	}
	if visible == 1 {
		return "1 session"
	}
	return fmt.Sprintf("%d sessions", visible)
}

func rule(width int) string {
	if width <= 0 {
		return ""
	}
	return strings.Repeat("─", width)
}

func sessionMatches(session core.Session, filter string) bool {
	return strings.Contains(strings.ToLower(session.Name), filter) ||
		strings.Contains(strings.ToLower(session.Directory), filter) ||
		strings.Contains(fmt.Sprintf("%dw", session.Windows), filter)
}

func sortSessions(sessions []core.Session, mode SessionSort) {
	sort.SliceStable(sessions, func(i, j int) bool {
		a, b := sessions[i], sessions[j]
		switch mode {
		case SortName:
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case SortWindows:
			if a.Windows == b.Windows {
				return strings.ToLower(a.Name) < strings.ToLower(b.Name)
			}
			return a.Windows > b.Windows
		case SortAttached:
			if a.Attached == b.Attached {
				return strings.ToLower(a.Name) < strings.ToLower(b.Name)
			}
			return a.Attached
		default:
			if a.Created.Equal(b.Created) {
				return strings.ToLower(a.Name) < strings.ToLower(b.Name)
			}
			return a.Created.After(b.Created)
		}
	})
}
