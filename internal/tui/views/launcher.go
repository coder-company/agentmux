package views

import (
	"agentmux/internal/core"
	"agentmux/internal/tui/styles"
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
	return &LauncherView{
		Workspaces: workspaces,
	}
}

// SelectedWorkspace returns the highlighted workspace.
func (l *LauncherView) SelectedWorkspace() *core.Workspace {
	if len(l.Workspaces) == 0 {
		return nil
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

// Render returns the launcher view.
func (l *LauncherView) Render() string {
	if len(l.Workspaces) == 0 {
		return styles.Panel.Width(l.Width).Render(
			styles.Title.Render("Workspaces") + "\n" +
				styles.Muted.Render("  No workspaces configured.\n  Add [[workspaces]] to config.toml"),
		)
	}

	var items string
	for i, ws := range l.Workspaces {
		label := ws.Name + "  " + styles.Muted.Render(ws.Root)
		if i == l.Cursor {
			items += styles.Selected.Width(l.Width-6).Render(label) + "\n"
			if l.InCommands {
				for j, cmd := range ws.Commands {
					cmdLabel := "    " + cmd.Name + "  " + styles.Muted.Render(cmd.Cmd)
					if j == l.CmdCursor {
						items += styles.Selected.Width(l.Width-6).Render(cmdLabel) + "\n"
					} else {
						items += styles.Normal.Render(cmdLabel) + "\n"
					}
				}
			}
		} else {
			items += styles.Normal.Width(l.Width-6).Render(label) + "\n"
		}
	}

	help := styles.Muted.Render("enter: open  tab: commands  esc: back")
	content := styles.Title.Render("Workspaces") + "\n" + items + "\n" + help
	return styles.ActivePanel.Width(l.Width).Height(l.Height - 2).Render(content)
}
