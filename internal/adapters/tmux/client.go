package tmux

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"

	"agentmux/internal/core"
)

// maxSessionNameLen is the maximum allowed length for a tmux session name.
const maxSessionNameLen = 128

// ValidateSessionName checks that name is a legal tmux session name.
// It rejects empty names, names containing '.' or ':' (tmux-reserved),
// names starting with '-', and names longer than 128 characters.
func ValidateSessionName(name string) error {
	if name == "" {
		return errors.New("session name must not be empty")
	}
	if utf8.RuneCountInString(name) > maxSessionNameLen {
		return fmt.Errorf("session name exceeds maximum length of %d characters", maxSessionNameLen)
	}
	if strings.HasPrefix(name, "-") {
		return errors.New("session name must not start with '-'")
	}
	if strings.ContainsAny(name, ".:") {
		return fmt.Errorf("session name must not contain '.' or ':' (got %q)", name)
	}
	return nil
}

// Client wraps tmux CLI commands.
type Client struct {
	bin string
}

// New creates a tmux client. Uses "tmux" from PATH by default.
func New() *Client {
	return &Client{bin: "tmux"}
}

// Available returns true if tmux is installed and reachable.
func (c *Client) Available() bool {
	_, err := exec.LookPath(c.bin)
	return err == nil
}

// ListSessions returns all tmux sessions.
func (c *Client) ListSessions() ([]core.Session, error) {
	out, err := c.run("list-sessions", "-F", sessionFormat)
	if err != nil {
		if strings.Contains(err.Error(), "no server running") ||
			strings.Contains(string(out), "no server running") {
			return nil, nil
		}
		return nil, fmt.Errorf("list-sessions: %w", err)
	}
	return ParseSessions(string(out))
}

// NewSession creates a new tmux session with the given name and optional start directory.
func (c *Client) NewSession(name, dir string) error {
	if err := ValidateSessionName(name); err != nil {
		return fmt.Errorf("new-session: %w", err)
	}
	args := []string{"new-session", "-d", "-s", name}
	if dir != "" {
		args = append(args, "-c", dir)
	}
	_, err := c.run(args...)
	return err
}

// KillSession destroys a tmux session.
func (c *Client) KillSession(name string) error {
	_, err := c.run("kill-session", "-t", name)
	return err
}

// RenameSession renames a tmux session.
func (c *Client) RenameSession(old, new string) error {
	if err := ValidateSessionName(new); err != nil {
		return fmt.Errorf("rename-session: %w", err)
	}
	_, err := c.run("rename-session", "-t", old, new)
	return err
}

// AttachSession attaches to a session (replaces the current process).
func (c *Client) AttachSession(name string) error {
	bin, err := exec.LookPath(c.bin)
	if err != nil {
		return err
	}
	cmd := exec.Command(bin, "attach-session", "-t", name)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// CapturePane captures visible pane output from the given session.
// Returns empty string (not error) when capture fails due to no pane content.
func (c *Client) CapturePane(target string, lines int) (string, error) {
	out, err := c.run("capture-pane", "-p", "-t", target,
		"-S", fmt.Sprintf("-%d", lines))
	if err != nil {
		// No pane content is a common case (empty pane, session just created).
		// Return empty rather than propagating the error.
		if strings.Contains(err.Error(), "no pane") ||
			strings.Contains(err.Error(), "can't find pane") ||
			strings.Contains(string(out), "can't find pane") {
			return "", nil
		}
		return "", err
	}
	return string(out), nil
}

func (c *Client) run(args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("tmux: no command specified")
	}
	cmd := exec.Command(c.bin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("%s %s: %w (%s)", c.bin, args[0], err, strings.TrimSpace(string(out)))
	}
	return out, nil
}
