# agentmux

A modern terminal workspace control plane over tmux.

agentmux is a keyboard-first TUI that wraps tmux — providing a session browser,
pane preview, command palette, and workspace launcher all from one interface.
It does not replace tmux. It makes tmux faster to navigate.

```
 ⣿ agentmux                                               3 sessions
╭──────────── Sessions ─────────────╮╭──────────── Preview ──────────────────────╮
│ ● main (3 win)  2h               ││ main  ~/code/agentmux                     │
│   api (1 win)  5h                ││                                            │
│   deploy (2 win)  1d             ││ $ go test ./...                             │
│                                   ││ ok  agentmux/internal/config    0.01s      │
│                                   ││ ok  agentmux/internal/store     0.06s      │
│                                   ││ ok  agentmux/internal/adapters  0.00s      │
│                                   ││ ok  agentmux/internal/tui/views 0.00s      │
╰───────────────────────────────────╯╰────────────────────────────────────────────╯
 ↑↓/jk nav │ enter attach │ n new │ x kill │ r rename │ / search │ ? help │ q quit
```

---

## Install

### One-liner (recommended)

```bash
curl -fsSL https://amux.coder.company/install.sh | sh
```

This auto-detects your OS/arch, downloads the latest release, and installs
`agentmux` + `atmux` to `/usr/local/bin`.

Install to a custom directory:

```bash
INSTALL_DIR=~/.local/bin curl -fsSL https://amux.coder.company/install.sh | sh
```

### From source

```bash
git clone https://github.com/codercompany/agentmux
cd agentmux
make install
```

### Build only

```bash
make build        # produces ./agentmux binary
```

### Shell alias (alternative to symlink)

If you prefer not to use the installer, add an alias:

```bash
# bash / zsh (~/.bashrc or ~/.zshrc)
alias atmux="agentmux"

# fish (~/.config/fish/config.fish)
alias atmux agentmux
```

### Requirements

| Dependency | Version |
|-----------|---------|
| Go        | 1.21+ (only for building from source) |
| tmux      | 3.0+    |

---

## Usage

### Launch the TUI

```bash
agentmux    # full-screen interactive mode
atmux       # same thing, shorter name
```

### CLI commands

```bash
agentmux list                    # List all tmux sessions
agentmux new myproject           # Create a detached session
agentmux new api --dir ~/code    # Create with a working directory
agentmux attach myproject        # Attach to a session (replaces process)
agentmux kill myproject          # Kill a session
agentmux init                    # Generate default config file
agentmux --version               # Print version
agentmux --help                  # Print help
```

All CLI commands check for tmux before running and produce clear error messages
with installation instructions if tmux is missing.

---

## TUI Interface

The TUI has four modes:

### Sessions view (default)

Two-panel layout:
- **Left panel**: Session list with name, window count, relative creation time, and attached indicator (`●`)
- **Right panel**: Live pane capture from the selected session (last 40 lines)

A **header** shows the app brand and total session count. A **footer** shows
context-aware keybindings.

### Command palette (`/`)

A centered overlay with fuzzy search over:
- Actions: new session, kill, rename, refresh, workspaces, help, quit
- Workspace launchers from your config
- Active sessions for quick attach

Fuzzy matching is character-order-aware (typing "ns" matches "New Session").

### Workspace launcher (`p`)

A centered overlay listing workspaces from your config. Press `tab` to drill
into workspace commands. Pressing `enter` creates a tmux session named
`workspace-command` (e.g., `api-run`) in the workspace root directory.

### Help overlay (`?`)

A categorized keybinding reference. Press `?` or `esc` to dismiss.

---

## Keybindings

### Navigation

| Key        | Action           |
|------------|------------------|
| `↑` / `k`  | Move cursor up   |
| `↓` / `j`  | Move cursor down |
| `gg`       | Jump to top      |
| `G`        | Jump to bottom   |

### Session actions

| Key     | Action                         |
|---------|--------------------------------|
| `enter` | Attach to selected session     |
| `n`     | Create new session (auto-named)|
| `x`     | Kill session (asks y/n)        |
| `r`     | Rename session inline          |
| `R`     | Refresh session list           |

### Views and overlays

| Key     | Action              |
|---------|---------------------|
| `/`     | Open command palette |
| `p`     | Open workspace launcher |
| `?`     | Toggle help overlay |
| `esc`   | Close any overlay   |
| `q`     | Quit                |
| `Ctrl+C`| Quit                |

### Command palette

| Key     | Action           |
|---------|------------------|
| `↑`/`↓` | Navigate results |
| `enter` | Execute action   |
| `esc`   | Close palette    |
| Type    | Fuzzy filter     |
| `backspace` | Delete character |

### Workspace launcher

| Key     | Action                |
|---------|-----------------------|
| `↑`/`↓`/`j`/`k` | Navigate  |
| `tab`   | Toggle commands list  |
| `enter` | Launch workspace/command |
| `esc`   | Back to sessions      |

---

## Configuration

Configuration lives at `~/.config/agentmux/config.toml`.

Run `agentmux init` to generate a starter file with commented examples.

### Full example

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

[[workspaces.commands]]
name = "build"
cmd = "npm run build"

[[workspaces]]
name = "infra"
root = "~/code/infra"
```

### Config validation

agentmux validates your config on load:
- Each workspace must have a non-empty `name` and `root`
- Workspace names must be unique
- Command names within a workspace must be unique
- Command `cmd` fields must be non-empty
- `~` and `~/path` are expanded to your home directory

Invalid configs are silently skipped (the TUI still works, just without workspaces).

### Session naming

When launching from the workspace view:
- Selecting a workspace creates a session named after the workspace (e.g., `api`)
- Drilling into commands and selecting one creates `workspace-command` (e.g., `api-run`)

---

## Architecture

```
cmd/agentmux/                  CLI entrypoint (main.go)
internal/
  app/                         Application wiring, TUI lifecycle, exec attach
  adapters/tmux/               tmux CLI wrapper, session parser, name validation
  commands/                    Cobra subcommands: list, new, attach, kill, init
  config/                      TOML loader, validation, InitDefault, expandHome
  core/                        Domain models: Session, Workspace, Command
  store/                       SQLite state: recent sessions, schema migrations
  tui/                         Bubble Tea app shell, main Update/View loop
    components/                Header, footer, session list, preview pane
    styles/                    Lip Gloss color palette and style definitions
    views/                     Sessions view, command palette, launcher, help overlay
```

### Key design decisions

| Decision | Rationale |
|----------|-----------|
| `\x1f` as tmux format delimiter | Pipes (`\|`) appear in file paths; Unit Separator never does |
| `syscall.Exec` for attach | Clean process replacement — no orphan parent process |
| SQLite via `modernc.org/sqlite` | Pure Go, no CGO required, works everywhere |
| Schema versioning in store | Forward-compatible migrations without breaking existing DBs |
| Corrupt DB auto-recovery | Renamed to `.corrupt.<timestamp>` and recreated fresh |
| Session name validation | Rejects `.`, `:`, leading `-`, empty, >128 chars — tmux reserved characters |

---

## Data storage

| File | Purpose |
|------|---------|
| `~/.config/agentmux/config.toml` | Workspace definitions |
| `~/.local/share/agentmux/state.db` | SQLite: recent sessions, schema version |

The SQLite database is optional — agentmux works fine without it (store errors are non-fatal).

---

## Development

```bash
make check     # gofmt + go vet + go test -race + go build
make test      # go test -race ./...
make build     # build binary with version from git
make smoke     # build + run CLI smoke test against tmux
make clean     # remove binary
make install   # install to PREFIX/bin (default: /usr/local)
make uninstall # remove from PREFIX/bin
```

### Test coverage

| Package | Tests |
|---------|-------|
| `adapters/tmux` | Parser: normal, empty, malformed, spaces, unicode, pipes, zero windows, long names. Validation: 10 cases |
| `config` | Load, missing file, expand home (bare + path), validation: missing name/root, duplicates, empty cmd, InitDefault |
| `store` | Round-trip, corrupt recovery, ping, schema version, idempotent migration |
| `tui/components` | Session list navigation, empty list safety, cursor clamping, relativeTime |
| `tui/views` | Fuzzy matching (14 cases), palette query/filter, navigation, type/backspace |

All tests are deterministic and do not require tmux.

### CI

GitHub Actions runs on push/PR to `main`:
- `gofmt` check
- `go vet`
- `go test -race`
- `go build`
- CLI smoke test

---

## Troubleshooting

### tmux not found

```
agentmux: tmux is not installed or not in PATH

Install tmux:
  macOS:  brew install tmux
  Debian: sudo apt install tmux
  Arch:   sudo pacman -S tmux
```

### Config parse error

If your config has a syntax error, agentmux logs it and runs without workspaces.
Fix the TOML syntax and restart. Validation errors name the specific workspace
and field (e.g., `workspace "api": missing root directory`).

### Corrupt database

If the SQLite database is corrupted, agentmux automatically:
1. Renames the corrupt file to `state.db.corrupt.<timestamp>`
2. Creates a fresh database
3. Prints a message to stderr

No data loss of tmux sessions — only the "recently used" ordering is reset.

### Terminal too small

If the terminal is narrower than ~50 columns or shorter than ~8 rows, agentmux
shows a "Terminal too small" message instead of crashing. Resize to continue.

### Session names with special characters

agentmux rejects session names containing `.` or `:` (tmux uses these as separators)
and names starting with `-` (tmux flag ambiguity). Spaces, unicode, hyphens, and
underscores are all allowed.

---

## Release

Releases use [GoReleaser](https://goreleaser.com/):

```bash
git tag v0.2.0
git push --tags
goreleaser release
```

Produces binaries for linux/darwin on amd64/arm64 with checksums.

---

## License

MIT — see [LICENSE](./LICENSE).
