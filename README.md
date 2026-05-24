# agentmux

A modern terminal workspace control plane over tmux.

Keyboard-first TUI for browsing sessions, previewing pane output, and launching workspaces — without leaving the terminal.

```
╭─ Sessions ──────────╮╭─ Preview ───────────────────────╮
│ ● main (3 win)      ││ ~/code/agentmux                 │
│   api (1 win)       ││ $ go test ./...                  │
│   deploy (2 win)    ││ ok  agentmux/internal/config     │
│                     ││ ok  agentmux/internal/store      │
│                     ││ ok  agentmux/internal/adapters   │
╰─────────────────────╯╰─────────────────────────────────╯
 enter attach  n new  k kill  r rename  / search  p projects  q quit
```

## Install

```bash
# From source
go install github.com/youruser/agentmux/cmd/agentmux@latest

# Or build locally
git clone https://github.com/youruser/agentmux && cd agentmux
make install
```

### Requirements

- Go 1.21+
- tmux 3.0+

## Usage

```bash
# Launch the TUI
agentmux

# CLI commands
agentmux list              # List sessions
agentmux new myproject     # Create a session
agentmux new api --dir ~/code/api  # Create with working directory
agentmux attach myproject  # Attach (replaces process)
agentmux kill myproject    # Kill a session
agentmux init              # Generate default config
```

## Keybindings

| Key       | Action              |
|-----------|---------------------|
| `enter`   | Attach to session   |
| `n`       | New session         |
| `k`       | Kill session (y/N)  |
| `r`       | Rename session      |
| `/`       | Command palette     |
| `p`       | Workspace launcher  |
| `↑`/`↓`  | Navigate            |
| `tab`     | Switch context      |
| `esc`     | Back                |
| `q`       | Quit                |

## Configuration

Config lives at `~/.config/agentmux/config.toml`. Run `agentmux init` to generate it.

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
  app/                Application wiring
  adapters/tmux/      tmux CLI wrapper and output parser
  commands/           Cobra CLI commands
  config/             TOML config loader with validation
  core/               Domain models (Session, Workspace)
  store/              SQLite local state with schema versioning
  tui/                Bubble Tea app shell
    components/       Reusable UI components
    styles/           Lip Gloss theme
    views/            View modules (sessions, palette, launcher)
```

## Development

```bash
make check     # fmt + vet + test
make build     # build binary
make smoke     # build + quick smoke test
```

## License

MIT
