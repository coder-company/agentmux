package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"agentmux/internal/core"

	"github.com/BurntSushi/toml"
)

// Config is the top-level configuration file structure.
type Config struct {
	Workspaces []workspaceEntry `toml:"workspaces"`
}

type workspaceEntry struct {
	Name     string         `toml:"name"`
	Root     string         `toml:"root"`
	Commands []commandEntry `toml:"commands"`
}

type commandEntry struct {
	Name string `toml:"name"`
	Cmd  string `toml:"cmd"`
}

// DefaultPath returns the default config file path.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "agentmux", "config.toml")
}

// Validate checks the config for structural correctness.
func (c *Config) Validate() error {
	var errs []error
	seen := make(map[string]bool)
	for _, w := range c.Workspaces {
		if w.Name == "" {
			errs = append(errs, fmt.Errorf("workspace: missing name"))
			continue
		}
		if w.Root == "" {
			errs = append(errs, fmt.Errorf("workspace %q: missing root directory", w.Name))
		}
		if seen[w.Name] {
			errs = append(errs, fmt.Errorf("workspace %q: duplicate name", w.Name))
		}
		seen[w.Name] = true

		cmdSeen := make(map[string]bool)
		for _, cmd := range w.Commands {
			if cmd.Cmd == "" {
				errs = append(errs, fmt.Errorf("workspace %q: command %q: empty cmd", w.Name, cmd.Name))
			}
			if cmd.Name != "" {
				if cmdSeen[cmd.Name] {
					errs = append(errs, fmt.Errorf("workspace %q: duplicate command %q", w.Name, cmd.Name))
				}
				cmdSeen[cmd.Name] = true
			}
		}
	}
	return errors.Join(errs...)
}

// Load reads and parses the config file at path.
func Load(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return &cfg, nil
}

// LoadOrDefault tries to load the config; returns an empty config if the file doesn't exist.
func LoadOrDefault(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		return &Config{}
	}
	return cfg
}

// Workspaces converts the parsed config into domain models.
func (c *Config) ToWorkspaces() []core.Workspace {
	out := make([]core.Workspace, len(c.Workspaces))
	for i, w := range c.Workspaces {
		cmds := make([]core.Command, len(w.Commands))
		for j, cmd := range w.Commands {
			cmds[j] = core.Command{Name: cmd.Name, Cmd: cmd.Cmd}
		}
		out[i] = core.Workspace{
			Name:     w.Name,
			Root:     expandHome(w.Root),
			Commands: cmds,
		}
	}
	return out
}

// InitDefault creates a default config file at path if it doesn't already exist.
func InitDefault(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil // file exists, do nothing
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("config: create directory: %w", err)
	}
	const example = `# agentmux configuration
# Uncomment and customize the example workspace below.

# [[workspaces]]
# name = "myproject"
# root = "~/code/myproject"
#
# [[workspaces.commands]]
# name = "run"
# cmd = "go run ./cmd/server"
#
# [[workspaces.commands]]
# name = "test"
# cmd = "go test ./..."
`
	return os.WriteFile(path, []byte(example), 0644)
}

func expandHome(path string) string {
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	if len(path) > 1 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
