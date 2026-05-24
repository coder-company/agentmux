package views

import (
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

// SelectedCommand returns the highlighted command within the workspace.
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

// MoveUp moves cursor up in the current context.
func (l *LauncherView) MoveUp() {
	if l.InCommands {
		if l.CmdCursor > 0 {
			l.CmdCursor--
		}
	} else {
		if l.Cursor > 0 {
			l.Cursor--
			l.CmdCursor = 0
		}
	}
}

// MoveDown moves cursor down in the current context.
func (l *LauncherView) MoveDown() {
	if l.InCommands {
		ws := l.SelectedWorkspace()
		if ws != nil && l.CmdCursor < len(ws.Commands)-1 {
			l.CmdCursor++
		}
	} else {
		if l.Cursor < len(l.Workspaces)-1 {
			l.Cursor++
			l.CmdCursor = 0
		}
	}
}

// ToggleCommands enters/exits the commands sub-list.
func (l *LauncherView) ToggleCommands() {
	l.InCommands = !l.InCommands
	l.CmdCursor = 0
}

// Render returns the launcher as a centered overlay.
func (l *LauncherView) Render() string {
	w := l.Width * 55 / 100
	if w < 44 {
		w = 44
	}
	if w > 80 {
		w = 80
	}

	title := styles.PanelTitle.Render("Workspaces")

	if len(l.Workspaces) == 0 {
		content := title + "\n\n" +
			styles.Muted.Render("No workspaces configured.") + "\n\n" +
			styles.Subtle.Render("Add [[workspaces]] to ~/.config/agentmux/config.toml") + "\n" +
			styles.Subtle.Render("Run: agentmux init")
		box := styles.Overlay.Width(w).Render(content)
		return lipgloss.Place(l.Width, l.Height, lipgloss.Center, lipgloss.Center, box)
	}

	var items string
	for i, ws := range l.Workspaces {
		name := ws.Name
		root := styles.Muted.Render("  " + ws.Root)

		if i == l.Cursor && !l.InCommands {
			items += styles.Selected.Width(w-6).Render(name+root) + "\n"
		} else if i == l.Cursor {
			items += styles.Bold.Render("▸ "+name) + root + "\n"
		} else {
			items += styles.Normal.Render("  "+name) + root + "\n"
		}

		// Show commands for selected workspace
		if i == l.Cursor && l.InCommands && len(ws.Commands) > 0 {
			for j, cmd := range ws.Commands {
				cmdLine := "    " + cmd.Name + "  " + styles.Subtle.Render(cmd.Cmd)
				if j == l.CmdCursor {
					items += styles.Selected.Width(w-6).Render(cmdLine) + "\n"
				} else {
					items += styles.NormalDim.Render(cmdLine) + "\n"
				}
			}
		}
	}

	content := title + "\n\n" + items
	box := styles.Overlay.Width(w).Render(content)
	return lipgloss.Place(l.Width, l.Height, lipgloss.Center, lipgloss.Center, box)
}
