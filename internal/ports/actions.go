package ports

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v3/process"
)

// KillProcess terminates the process by PID.
func KillProcess(pid int32) error {
	if pid <= 0 {
		return fmt.Errorf("invalid PID: %d", pid)
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("process %d not found: %w", pid, err)
	}

	if err := p.Kill(); err != nil {
		return fmt.Errorf("failed to kill PID %d (try running as admin): %w", pid, err)
	}

	return nil
}

// KillByPort finds the process on a port and kills it.
func KillByPort(port uint32) error {
	if port == 0 || port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", port)
	}

	infos, err := Scan(context.Background(), ScanOptions{MinPort: port, MaxPort: port})
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("no process found listening on port %d", port)
	}

	for _, info := range infos {
		if err := KillProcess(info.PID); err != nil {
			return err
		}
		fmt.Printf("Killed %s (PID %d) on port %d\n", info.ProcessName, info.PID, info.LocalPort)
	}
	return nil
}

// OpenInBrowser opens localhost:port in the default browser.
func OpenInBrowser(port uint32) error {
	url := fmt.Sprintf("http://localhost:%d", port)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot open browser: %w", err)
	}
	// Don't wait for browser process
	go func() { _ = cmd.Wait() }()
	return nil
}

// FormatPortAddr returns host:port string.
func FormatPortAddr(addr string, port uint32) string {
	if addr == "" || addr == "*" || addr == "0.0.0.0" || addr == "::" {
		return ":" + strconv.FormatUint(uint64(port), 10)
	}
	return addr + ":" + strconv.FormatUint(uint64(port), 10)
}
