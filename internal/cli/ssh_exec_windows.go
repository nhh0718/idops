//go:build windows

package cli

import "errors"

// syscallExec is not available on Windows; callers should use exec.Command instead.
func syscallExec(bin string, args []string) error {
	return errors.New("syscall.Exec not supported on Windows")
}
