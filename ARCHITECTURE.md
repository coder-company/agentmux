# Architecture

agentmux is a Go TUI application that wraps tmux via CLI commands. It does not
use tmux control mode (reserved for a potential v2). The binary is a single
executable with no runtime dependencies beyond tmux itself.

## Dependency graph

```
cmd/agentmux/main.go
  └── internal/commands       (Cobra CLI)
        └── internal/app      (wiring + lifecycle)
              ├── internal/adapters/tmux   (tmux CLI wrapper)
              ├── internal/config          (TOML config)
              ├── internal/store           (SQLite state)
              └── internal/tui            (Bubble Tea app)
                    ├── internal/tui/components  (header, footer, list, preview)
                    ├── internal/tui/styles      (Lip Gloss theme)
                    └── internal/tui/views       (sessions, palette, launcher, help)
```

All packages under `internal/` are unexported. The only public API is the CLI.

## Package responsibilities

### cmd/agentmux

Entrypoint. Calls `commands.Root().Execute()` and exits with code 1 on error.

### internal/commands

Cobra command tree. Each file is one command:
- `root.go` — launches TUI when no subcommand given
- `list.go` — tabwriter output of sessions
- `new.go` — creates a session with `--dir` flag
- `attach.go` — attaches (process replacement via app layer)
- `kill.go` — kills a session
- `init.go` — creates default config file
- `errors.go` — shared error messages (tmux missing)

### internal/app

Application bootstrap:
- `New()` wires all dependencies (client, config, store)
- `RunTUI()` starts Bubble Tea, then does `syscall.Exec` for attach on exit
- `exec_unix.go` / `exec_windows.go` — platform-specific process replacement

### internal/adapters/tmux

tmux CLI wrapper. Never shells out via `sh -c` — all arguments are passed as
exec args (no injection risk).

Key design:
- Format strings use `\x1f` (Unit Separator) as field delimiter instead of `|`
  because pipes can appear in directory paths
- `ValidateSessionName()` rejects characters that tmux treats specially
- `CapturePane()` returns empty string (not error) for missing/empty panes
- `ListSessions()` returns `nil, nil` when no server is running

### internal/config

TOML config loading with validation:
- `Load(path)` — parse + validate
- `LoadOrDefault(path)` — returns empty config on any error (non-fatal)
- `Validate()` — structural checks with joined errors
- `InitDefault(path)` — creates commented example if file doesn't exist
- `expandHome()` — handles both `~` and `~/path`

### internal/core

Pure domain models with no dependencies:
- `Session` — name, windows, created, attached, directory
- `Workspace` — name, root, commands
- `Command` — name, cmd

### internal/store

SQLite persistence via `modernc.org/sqlite` (pure Go, no CGO):
- Schema versioning with incremental migrations
- Corrupt database auto-recovery (rename + recreate)
- `RecordSession()` — upsert with nanosecond timestamp
- `RecentSessions(limit)` — ordered by last use
- `Ping()` — connection health check

### internal/tui

Bubble Tea application shell:
- `Model` — top-level state machine
- Four view modes: Sessions, Palette, Launcher, Help
- Key handling dispatched by current mode
- `doRefresh()` reloads sessions and rebuilds the palette action list
- `rebuildPalette()` dynamically includes workspace and session actions

### internal/tui/components

Reusable rendering primitives:
- `Header` — brand + mode + session count, fills width
- `RenderFooter()` — context-aware keybinding bar with `│` separators
- `SessionList` — cursor-navigated list with attached indicators and timestamps
- `Preview` — pane output with title, directory, error/empty states, line truncation

### internal/tui/styles

Lip Gloss style definitions. A single dark theme with violet accents.
Organized into sections: layout (header, panels, footer), content (selected,
normal, muted), indicators, overlays, and input.

### internal/tui/views

Full-screen view renderers:
- `SessionsView` — two-panel layout (list + preview), delegates to components
- `PaletteView` — centered overlay with fuzzy search, action list, keyboard nav
- `LauncherView` — centered overlay for workspace selection + command drilling
- `HelpOverlay` — categorized keybinding reference

## Data flow

```
User input → tea.KeyMsg → Model.Update() → state mutation → Model.View() → terminal
                              │
                              ├── handleSessionsKey() → tmux.Client (list/new/kill/rename)
                              ├── handlePaletteKey()  → action callbacks
                              ├── handleLauncherKey() → tmux.Client.NewSession()
                              └── handleConfirmKey()  → tmux.Client.KillSession()
```

## Process lifecycle

1. `main()` → `commands.Root().Execute()`
2. If no subcommand: `app.New()` → `app.RunTUI()`
3. Bubble Tea runs until user quits or selects attach
4. On attach: `tea.Program` exits, `execvp("tmux", "attach-session", "-t", name)`
   replaces the process — no parent left behind

## Testing strategy

- Unit tests cover parsers, validators, store operations, and UI logic
- All tests are deterministic — no tmux required
- `make smoke` builds and runs CLI commands against a real tmux server
- CI runs `go test -race` to catch data races
- Race detector is free since the TUI is single-threaded (Bubble Tea model)
