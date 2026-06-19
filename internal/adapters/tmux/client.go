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

// isNoServer reports whether the error/output indicates that no tmux server
// is currently running. tmux phrases this differently across versions:
//   - "no server running on ..."
//   - "error connecting to /tmp/tmux-<uid>/default (No such file or directory)"
func isNoServer(err error, out []byte) bool {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	outStr := string(out)
	return strings.Contains(msg, "no server running") ||
		strings.Contains(outStr, "no server running") ||
		strings.Contains(msg, "error connecting to") ||
		strings.Contains(outStr, "error connecting to")
}

// ListSessions returns all tmux sessions.
func (c *Client) ListSessions() ([]core.Session, error) {
	out, err := c.run("list-sessions", "-F", sessionFormat)
	if err != nil {
		if isNoServer(err, out) {
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

// ListWindows returns the window names for a session.
func (c *Client) ListWindows(session string) ([]string, error) {
	out, err := c.run("list-windows", "-t", session, "-F", "#{window_name}")
	if err != nil {
		return nil, err
	}
	var windows []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			windows = append(windows, line)
		}
	}
	return windows, nil
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
