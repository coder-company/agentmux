package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the global key bindings.
type KeyMap struct {
	Quit    key.Binding
	Enter   key.Binding
	New     key.Binding
	Kill    key.Binding
	Rename  key.Binding
	Project key.Binding
	Slash   key.Binding
	Escape  key.Binding
	Up      key.Binding
	Down    key.Binding
	Tab     key.Binding
}

// DefaultKeys returns the default key bindings.
func DefaultKeys() KeyMap {
	return KeyMap{
		Quit:    key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
		Enter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "attach")),
		New:     key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new session")),
		Kill:    key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "kill session")),
		Rename:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename")),
		Project: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "projects")),
		Slash:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Escape:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Up:      key.NewBinding(key.WithKeys("up", "ctrl+p"), key.WithHelp("↑", "up")),
		Down:    key.NewBinding(key.WithKeys("down", "ctrl+n"), key.WithHelp("↓", "down")),
		Tab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch pane")),
	}
}
