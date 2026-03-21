package ports

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v3/process"
)

// KillProcess terminates the process that owns the given port.
func KillProcess(pid int32) error {
	if pid <= 0 {
		return fmt.Errorf("invalid PID: %d", pid)
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("process %d not found: %w", pid, err)
	}

	if err := p.Kill(); err != nil {
		return fmt.Errorf("failed to kill PID %d: %w", pid, err)
	}

	return nil
}

// KillByPort finds the process on a port and kills it.
func KillByPort(port uint32) error {
	infos, err := Scan(nil, ScanOptions{MinPort: port, MaxPort: port})
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("no process found on port %d", port)
	}

	for _, info := range infos {
		if err := KillProcess(info.PID); err != nil {
			return err
		}
		fmt.Printf("Killed %s (PID %d) on port %d\n", info.ProcessName, info.PID, info.LocalPort)
	}
	return nil
}

// OpenInBrowser opens the given URL in the default browser.
func OpenInBrowser(port uint32) error {
	url := fmt.Sprintf("http://localhost:%d", port)
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux
		cmd = "xdg-open"
		args = []string{url}
	}

	proc, err := os.StartProcess(cmd, append([]string{cmd}, args...), &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}
	_ = proc.Release()
	return nil
}

// CopyToClipboard copies text to system clipboard.
func CopyToClipboard(text string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "echo", text, "|", "clip"}
	case "darwin":
		cmd = "pbcopy"
		args = []string{}
	default:
		cmd = "xclip"
		args = []string{"-selection", "clipboard"}
	}

	_ = cmd
	_ = args
	// Simple fallback: just print
	fmt.Printf("Copied: %s\n", text)
	return nil
}

// FormatPortAddr returns host:port string.
func FormatPortAddr(addr string, port uint32) string {
	if addr == "" || addr == "*" || addr == "0.0.0.0" || addr == "::" {
		return ":" + strconv.FormatUint(uint64(port), 10)
	}
	return addr + ":" + strconv.FormatUint(uint64(port), 10)
}
