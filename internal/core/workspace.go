package core

// Workspace is a named project directory with predefined commands.
type Workspace struct {
	Name     string
	Root     string
	Commands []Command
}

// Command is a named shell command inside a workspace.
type Command struct {
	Name string
	Cmd  string
}
