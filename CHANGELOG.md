# Changelog

## v0.1.0 — 2025-05-24

Initial release.

### Features

- Full-screen TUI with session browser and pane preview
- Command palette with fuzzy search (/)
- Workspace launcher from TOML config (p)
- Inline session rename (r) with validation
- Kill confirmation prompt (y/N)
- CLI commands: `list`, `new`, `attach`, `kill`, `init`
- SQLite local state for recent sessions
- Config at `~/.config/agentmux/config.toml`
- Session name validation (rejects reserved chars)
- Graceful handling of missing tmux, no sessions, corrupt DB
- Schema-versioned database with automatic recovery
