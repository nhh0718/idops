package cli

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

// openBrowser opens a URL in the default browser.
func openBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

// isPortFree checks if a port is available by trying to connect to it.
// Uses Dial instead of Listen because on Windows, Listen can succeed
// even when another process holds the port (dual-stack IPv4/IPv6).
func isPortFree(port string) bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 500*time.Millisecond)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

// findFreePort scans a range and returns the first free port as string.
func findFreePort(start, end int) string {
	for p := start; p <= end; p++ {
		port := fmt.Sprintf("%d", p)
		if isPortFree(port) {
			return port
		}
	}
	return ""
}

// mustAtoi converts string to int, returns 0 on error.
func mustAtoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
