//go:build windows

package app

import (
	"fmt"
	"os"
	"os/exec"
)

// Execvp replaces the current process with the given binary.
func Execvp(bin string, args []string, env []string) error {
	cmd := exec.Command(bin, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("attach: %w", err)
	}
	os.Exit(0)
	return nil
}
