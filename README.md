# agentmux

A modern terminal workspace control plane over tmux.

Keyboard-first TUI for browsing sessions, previewing pane output, and launching
workspaces — without leaving the terminal.

```
 ⣿ agentmux · sessions                                    3 sessions
╭─────────── Sessions ───────────╮╭─────────── Preview ──────────────────╮
│ ● main (3 win)  2h            ││ main  ~/code/agentmux                │
│   api (1 win)  5h             ││                                      │
│   deploy (2 win)  1d          ││ $ go test ./...                       │
│                                ││ ok  agentmux/internal/config  0.01s  │
│                                ││ ok  agentmux/internal/store   0.06s  │
│                                ││ ok  agentmux/internal/adapters 0.00s │
╰────────────────────────────────╯╰──────────────────────────────────────╯
 ↑↓/jk nav │ enter attach │ n new │ x kill │ r rename │ / search │ ? help │ q quit
```

## Install

```bash
# Build and install (creates agentmux + atmux symlink)
git clone https://github.com/youruser/agentmux && cd agentmux
make install

# Or install to a custom prefix
make install PREFIX=~/.local

# Or just build
make build
```

### Shell alias (alternative to symlink)

```bash
# bash/zsh
alias atmux="agentmux"

# fish
alias atmux agentmux
```

### Requirements

- Go 1.21+
- tmux 3.0+

## Usage

```bash
# Launch the TUI (both commands work)
agentmux
atmux

# CLI commands
agentmux list                  # List sessions
agentmux new myproject         # Create a session
agentmux new api --dir ~/code  # Create with working directory
agentmux attach myproject      # Attach (replaces process)
agentmux kill myproject        # Kill a session
agentmux init                  # Generate default config
```

## Keybindings

### Navigation

| Key        | Action           |
|------------|------------------|
| `↑`/`↓`   | Move cursor      |
| `j`/`k`    | Move cursor (vim)|
| `gg`       | Jump to top      |
| `G`        | Jump to bottom   |

### Sessions

| Key     | Action               |
|---------|----------------------|
| `enter` | Attach to session    |
| `n`     | Create new session   |
| `x`     | Kill session (y/n)   |
| `r`     | Rename session       |
| `R`     | Refresh list         |

### Views

| Key     | Action              |
|---------|---------------------|
| `/`     | Command palette     |
| `p`     | Workspace launcher  |
| `?`     | Help overlay        |
| `esc`   | Close overlay       |
| `q`     | Quit                |

## Command Palette

Press `/` to open a unified search over:

- Actions (new, kill, rename, refresh, quit)
- Workspace launchers from config
- Active sessions (quick-attach)

Fuzzy matching — type a few characters to filter instantly.

## Configuration

Config lives at `~/.config/agentmux/config.toml`. Run `agentmux init` to generate a starter file.

```toml
[[workspaces]]
name = "api"
root = "~/code/api"

[[workspaces.commands]]
name = "run"
cmd = "go run ./cmd/server"

[[workspaces.commands]]
name = "test"
cmd = "go test ./..."

[[workspaces]]
name = "frontend"
root = "~/code/web"

[[workspaces.commands]]
name = "dev"
cmd = "npm run dev"
```

When launching from the workspace view, sessions are named `workspace-command`
(e.g., `api-run`, `frontend-dev`).

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
    components/       Header, footer, session list, preview pane
    styles/           Lip Gloss theme (colors, panels, overlays)
    views/            Sessions, command palette, launcher, help
```

## Development

```bash
make check     # fmt + vet + test + build
make smoke     # build + quick tmux smoke test
make build     # build binary
make install   # install agentmux + atmux symlink
```

## Troubleshooting

**tmux not found**

```
agentmux: tmux is not installed or not in PATH

Install tmux:
  macOS:  brew install tmux
  Debian: sudo apt install tmux
  Arch:   sudo pacman -S tmux
```

**Config errors**

Run `agentmux init` to create a valid starter config. Validation errors include
the workspace name and field that failed.

**Corrupt database**

The SQLite state file at `~/.local/share/agentmux/state.db` is auto-recovered
if corrupted. The corrupt file is preserved as `state.db.corrupt.<timestamp>`.

## License

MIT
