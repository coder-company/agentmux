//go:build !windows

package app

import "syscall"

// Execvp replaces the current process with the given binary.
func Execvp(bin string, args []string, env []string) error {
	return syscall.Exec(bin, args, env)
}
