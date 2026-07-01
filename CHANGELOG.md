# Changelog

## v0.2.3 — 2026-07-01

### TUI

- Replaced the old accent treatment with a neutral charcoal palette and
  green/cyan terminal accents.
- Added live session filtering with `/`.
- Added sort cycling with `s`: newest, name, windows, attached.
- Added layout cycling with `tab`: split, list-only, preview-only.
- Expanded preview metadata with attached state, window count, age, and cwd.
- Moved the command palette to `:` and `Ctrl+P`.

## v0.2.0 — 2025-05-24

### Documentation

- Full README with terminal mockup, install methods, keybinding tables,
  config reference, troubleshooting guide, and architecture overview
- ARCHITECTURE.md with dependency graph, package responsibilities,
  data flow, process lifecycle, and testing strategy
- Shell alias documentation for bash, zsh, and fish

### TUI Overhaul

- Header bar with brand, mode indicator, session count
- Context-aware footer with keybinding hints
- Refined dark theme with high-contrast accents
- Session list with relative timestamps and attached indicators
- Preview pane with directory path, safe content truncation
- Centered command palette overlay with fuzzy search
- Centered workspace launcher overlay with command drilling
- Help overlay with categorized keybinding reference (?)
- Empty states with actionable hints throughout

### Navigation

- Vim-style movement: j/k, gg (top), G (bottom)
- Kill moved to x (k reserved for vim up)
- R for manual refresh
- ? for help overlay
- Palette includes attach-by-name and workspace launch actions

### Workspace Launcher

- Auto-names sessions as workspace-command (e.g., api-run)
- Tab into commands sub-list

### Alias Support

- `make install` creates atmux symlink alongside agentmux
- Shell alias docs for bash, zsh, fish

---

## v0.1.0 — 2025-05-24

Initial release.

### Features

- Full-screen TUI with session browser and pane preview
- Command palette with fuzzy search
- Workspace launcher from TOML config (p)
- Inline session rename (r) with validation
- Kill confirmation prompt (y/N)
- CLI commands: `list`, `new`, `attach`, `kill`, `init`
- SQLite local state for recent sessions
- Config at `~/.config/agentmux/config.toml`
- Session name validation (rejects reserved chars)
- Graceful handling of missing tmux, no sessions, corrupt DB
- Schema-versioned database with automatic recovery
