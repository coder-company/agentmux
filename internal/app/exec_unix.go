//go:build !windows

package app

import "syscall"

func execvp(bin string, args []string, env []string) error {
	return syscall.Exec(bin, args, env)
}
