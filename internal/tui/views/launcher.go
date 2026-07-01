package views

import (
	"fmt"
	"strings"

	"agentmux/internal/core"
	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// LauncherView shows workspaces from config.
type LauncherView struct {
	Workspaces []core.Workspace
	Cursor     int
	CmdCursor  int
	InCommands bool
	Width      int
	Height     int
}

// NewLauncher creates a workspace launcher view.
func NewLauncher(workspaces []core.Workspace) *LauncherView {
	return &LauncherView{Workspaces: workspaces}
}

// SelectedWorkspace returns the highlighted workspace.
func (l *LauncherView) SelectedWorkspace() *core.Workspace {
	if len(l.Workspaces) == 0 {
		return nil
	}
	if l.Cursor >= len(l.Workspaces) {
		l.Cursor = len(l.Workspaces) - 1
	}
	return &l.Workspaces[l.Cursor]
}

// SelectedCommand returns the highlighted command.
func (l *LauncherView) SelectedCommand() *core.Command {
	ws := l.SelectedWorkspace()
	if ws == nil || len(ws.Commands) == 0 {
		return nil
	}
	if l.CmdCursor >= len(ws.Commands) {
		l.CmdCursor = 0
	}
	return &ws.Commands[l.CmdCursor]
}

// MoveUp moves cursor up.
func (l *LauncherView) MoveUp() {
	if l.InCommands {
		if l.CmdCursor > 0 {
			l.CmdCursor--
		}
	} else if l.Cursor > 0 {
		l.Cursor--
		l.CmdCursor = 0
	}
}

// MoveDown moves cursor down.
func (l *LauncherView) MoveDown() {
	if l.InCommands {
		ws := l.SelectedWorkspace()
		if ws != nil && l.CmdCursor < len(ws.Commands)-1 {
			l.CmdCursor++
		}
	} else if l.Cursor < len(l.Workspaces)-1 {
		l.Cursor++
		l.CmdCursor = 0
	}
}

// ToggleCommands enters/exits the commands sub-list.
func (l *LauncherView) ToggleCommands() {
	ws := l.SelectedWorkspace()
	if ws == nil || len(ws.Commands) == 0 {
		l.InCommands = false
		l.CmdCursor = 0
		return
	}
	l.InCommands = !l.InCommands
	l.CmdCursor = 0
}

// Render returns the launcher as a centered overlay.
func (l *LauncherView) Render() string {
	w := l.Width * 62 / 100
	if w < 52 {
		w = 52
	}
	if w > 88 {
		w = 88
	}
	if l.Width > 0 && w > l.Width-4 {
		w = l.Width - 4
	}
	if w < 26 {
		w = 26
	}
	innerW := w - 6
	if innerW < 14 {
		innerW = 14
	}

	title := styles.OverlayTitle.Render("Workspace launcher")
	count := styles.OverlayDim.Render(workspaceCount(len(l.Workspaces)))
	titleRow := title
	gap := innerW - lipgloss.Width(title) - lipgloss.Width(count)
	if gap > 0 {
		titleRow = title + lipgloss.NewStyle().Width(gap).Render("") + count
	}

	if len(l.Workspaces) == 0 {
		content := titleRow + "\n\n" +
			styles.Muted.Render("  No workspaces configured.") + "\n\n" +
			styles.HeaderDim.Render("  Edit ~/.config/agentmux/config.toml") + "\n" +
			styles.HeaderDim.Render("  Run: agentmux init")
		box := styles.Overlay.Width(w).Render(content)
		return lipgloss.Place(l.Width, l.Height, lipgloss.Center, lipgloss.Center, box)
	}

	maxWorkspaces := l.Height - 14
	if maxWorkspaces < 3 {
		maxWorkspaces = 3
	}
	if maxWorkspaces > 8 {
		maxWorkspaces = 8
	}

	var items string
	start, visible := l.visibleWorkspaces(maxWorkspaces)
	for i, ws := range visible {
		index := start + i
		items += renderWorkspaceRow(ws, innerW, index == l.Cursor && !l.InCommands) + "\n"
	}
	if len(l.Workspaces) > maxWorkspaces {
		items += styles.OverlayDim.Render(scrollLabel(start, len(visible), len(l.Workspaces), innerW)) + "\n"
	}

	commandBlock := l.renderCommandBlock(innerW, max(3, l.Height-18))
	footer := styles.HeaderDim.Render(styles.Truncate("↑↓ move · tab commands · enter launch · esc back", innerW))
	content := titleRow + "\n\n" + items + "\n" + commandBlock + "\n" + footer
	box := styles.Overlay.Width(w).Render(content)
	return lipgloss.Place(l.Width, l.Height, lipgloss.Center, lipgloss.Center, box)
}

func (l *LauncherView) visibleWorkspaces(maxItems int) (int, []core.Workspace) {
	if len(l.Workspaces) <= maxItems {
		return 0, l.Workspaces
	}

	start := l.Cursor - maxItems/2
	if start < 0 {
		start = 0
	}
	if start+maxItems > len(l.Workspaces) {
		start = len(l.Workspaces) - maxItems
	}
	return start, l.Workspaces[start : start+maxItems]
}

func (l *LauncherView) renderCommandBlock(width, maxRows int) string {
	ws := l.SelectedWorkspace()
	if ws == nil {
		return ""
	}

	title := "Commands for " + ws.Name
	lines := []string{styles.PanelMeta.Render(styles.Truncate(title, width))}
	if len(ws.Commands) == 0 {
		lines = append(lines, styles.Muted.Render("  enter starts a shell at the workspace root"))
		return strings.Join(lines, "\n")
	}

	limit := maxRows
	if limit > len(ws.Commands) {
		limit = len(ws.Commands)
	}
	if limit < 1 {
		limit = 1
	}

	for i := 0; i < limit; i++ {
		selected := l.InCommands && i == l.CmdCursor
		lines = append(lines, renderCommandRow(ws.Commands[i], width, selected))
	}
	if len(ws.Commands) > limit {
		lines = append(lines, styles.OverlayDim.Render(styles.Truncate(fmt.Sprintf("  %d more commands", len(ws.Commands)-limit), width)))
	}
	return strings.Join(lines, "\n")
}

func renderWorkspaceRow(ws core.Workspace, width int, selected bool) string {
	count := fmt.Sprintf("%dc", len(ws.Commands))
	countW := lipgloss.Width(count)
	nameW := width
	rootW := 0
	if width >= 48 {
		nameW = 20
		rootW = width - nameW - countW - 4
	}
	if rootW < 0 {
		rootW = 0
	}

	name := styles.PadRight(ws.Name, nameW)
	line := name
	if rootW > 0 {
		line += "  " + styles.PadRight(ws.Root, rootW)
	}
	if countW > 0 {
		gap := width - lipgloss.Width(line) - countW
		if gap < 1 {
			gap = 1
		}
		line += strings.Repeat(" ", gap) + count
	}
	line = styles.Truncate(line, width)

	if selected {
		return styles.OverlaySelected.Width(width).Render(line)
	}
	if rootW > 0 {
		return styles.OverlayNormal.Render(styles.PadRight(ws.Name, nameW)) +
			"  " + styles.OverlayDim.Render(styles.PadRight(ws.Root, rootW)) +
			strings.Repeat(" ", max(1, width-nameW-rootW-countW-2)) +
			styles.OverlayDim.Render(count)
	}
	return styles.OverlayNormal.Render(line)
}

func renderCommandRow(cmd core.Command, width int, selected bool) string {
	nameW := 18
	if width < 44 {
		nameW = width / 2
	}
	if nameW < 8 {
		nameW = 8
	}
	cmdW := width - nameW - 4
	if cmdW < 1 {
		cmdW = 1
	}
	line := "  " + styles.PadRight(cmd.Name, nameW) + "  " + styles.Truncate(cmd.Cmd, cmdW)
	line = styles.Truncate(line, width)
	if selected {
		return styles.OverlaySelected.Width(width).Render(line)
	}
	return styles.OverlayNormal.Render("  "+styles.PadRight(cmd.Name, nameW)) +
		"  " + styles.OverlayDim.Render(styles.Truncate(cmd.Cmd, cmdW))
}

func workspaceCount(n int) string {
	if n == 1 {
		return "1 workspace"
	}
	return fmt.Sprintf("%d workspaces", n)
}

func scrollLabel(start, visible, total, width int) string {
	return styles.Truncate(fmt.Sprintf("  showing %d-%d of %d", start+1, start+visible, total), width)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
