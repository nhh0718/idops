//go:build !windows

package cli

import "syscall"

// syscallExec replaces the current process with ssh on Unix systems.
func syscallExec(bin string, args []string) error {
	return syscall.Exec(bin, args, syscall.Environ())
}
