# agentmux

A modern terminal workspace control plane over tmux. Keyboard-first TUI for managing sessions, previewing pane output, and launching workspaces.

## Install

```bash
go install agentmux/cmd/agentmux@latest
```

Or build from source:

```bash
go build -o agentmux ./cmd/agentmux
```

## Usage

```bash
# Launch the TUI
agentmux

# CLI commands
agentmux list          # List tmux sessions
agentmux new <name>    # Create a new session
agentmux attach <name> # Attach to a session
agentmux kill <name>   # Kill a session
```

## TUI Keybindings

| Key     | Action            |
|---------|-------------------|
| `enter` | Attach to session |
| `n`     | New session       |
| `k`     | Kill session      |
| `r`     | Rename session    |
| `/`     | Command palette   |
| `p`     | Project launcher  |
| `↑/↓`   | Navigate          |
| `tab`   | Switch context    |
| `esc`   | Back              |
| `q`     | Quit              |

## Configuration

Config lives at `~/.config/agentmux/config.toml`:

```toml
[[workspaces]]
name = "myproject"
root = "~/code/myproject"

[[workspaces.commands]]
name = "run"
cmd = "go run ./cmd/server"

[[workspaces.commands]]
name = "test"
cmd = "go test ./..."
```

## Architecture

```
cmd/agentmux/         CLI entrypoint
internal/
  app/                Application initialization and wiring
  adapters/tmux/      tmux CLI wrapper and output parser
  commands/           Cobra CLI commands
  config/             TOML config loader
  core/               Domain models (Session, Workspace)
  store/              SQLite local state
  tui/                Bubble Tea app shell
    components/       Reusable UI components
    styles/           Lip Gloss theme
    views/            View modules (sessions, palette, launcher)
```

## Requirements

- Go 1.21+
- tmux 3.0+

## Development

```bash
go test ./...
go vet ./...
gofmt -l .
```
