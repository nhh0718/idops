package cli

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/nhh0718/idops/internal/ports"
	"github.com/spf13/cobra"
)

var dashboardPort string
var dashboardNoOpen bool

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch the web dashboard",
	Long:  `Start the Next.js dashboard server and open it in your browser.`,
	RunE:  runDashboard,
}

func init() {
	dashboardCmd.Flags().StringVarP(&dashboardPort, "port", "p", "3000", "Port to run dashboard on")
	dashboardCmd.Flags().BoolVar(&dashboardNoOpen, "no-open", false, "Don't open browser automatically")
	rootCmd.AddCommand(dashboardCmd)
}

func runDashboard(cmd *cobra.Command, args []string) error {
	// Find dashboard directory
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get executable path: %w", err)
	}
	// Resolve symlinks to get real binary location
	realExecPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		realExecPath = execPath
	}
	execDir := filepath.Dir(execPath)
	realExecDir := filepath.Dir(realExecPath)

	// Also check user home install dir
	homeDir, _ := os.UserHomeDir()
	idopsHome := filepath.Join(homeDir, ".idops")

	// Try multiple locations for dashboard
	possiblePaths := []string{
		filepath.Join(realExecDir, "dashboard"),        // next to real binary
		filepath.Join(execDir, "dashboard"),             // next to binary (or symlink)
		filepath.Join(realExecDir, "..", "dashboard"),   // parent of real binary
		filepath.Join(idopsHome, "dashboard"),           // ~/.idops/dashboard
		filepath.Join(homeDir, ".local", "share", "idops", "dashboard"), // XDG data
		"./dashboard",                                   // CWD
	}

	var dashboardPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
			dashboardPath = path
			break
		}
	}

	if dashboardPath == "" {
		return fmt.Errorf("dashboard not found. Searched:\n  %s\nRun 'idops update' to download dashboard, or clone repo and run from project root", strings.Join(possiblePaths, "\n  "))
	}

	green := "\033[32;1m"
	cyan := "\033[36m"
	yellow := "\033[33m"
	reset := "\033[0m"

	// Kill any previous dashboard instance on target port
	if !isPortFree(dashboardPort) {
		fmt.Printf("%s🔄 Port %s đang bị chiếm, đang dọn dẹp...%s\n", yellow, dashboardPort, reset)
		_ = ports.KillByPort(uint32(mustAtoi(dashboardPort)))
		time.Sleep(1 * time.Second) // Grace period for port release

		// Recheck after kill
		if !isPortFree(dashboardPort) {
			if dashboardPort == "3000" {
				newPort := findFreePort(3001, 3100)
				if newPort == "" {
					return fmt.Errorf("port %s vẫn bị chiếm sau khi kill, không tìm được port trống (3001-3100)", dashboardPort)
				}
				fmt.Printf("%s⚠️  Không thể giải phóng port %s, chuyển sang port %s%s\n", yellow, dashboardPort, newPort, reset)
				dashboardPort = newPort
			} else {
				return fmt.Errorf("port %s is still in use after cleanup. Try --port with a different port", dashboardPort)
			}
		} else {
			fmt.Printf("%s✅ Port %s đã được giải phóng%s\n", green, dashboardPort, reset)
		}
	}

	fmt.Printf("%s🚀 Starting idops dashboard...%s\n", green, reset)
	fmt.Println()

	// Check if npm is available
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("npm not found in PATH. Please install Node.js")
	}

	// Start Next.js dev server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println()
		fmt.Printf("%s⚡ Shutting down dashboard...%s\n", yellow, reset)
		cancel()
	}()

	// Auto-build if .next/BUILD_ID missing
	var npmCmd *exec.Cmd
	nextBuildID := filepath.Join(dashboardPath, ".next", "BUILD_ID")
	if _, err := os.Stat(nextBuildID); err != nil {
		fmt.Printf("%s🔨 Đang build dashboard (lần đầu)...%s\n", yellow, reset)
		buildCmd := exec.CommandContext(ctx, "npm", "run", "build")
		buildCmd.Dir = dashboardPath
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			fmt.Printf("%s⚠️  Build thất bại, chuyển sang dev mode%s\n", yellow, reset)
		}
	}

	// Use production build if .next/BUILD_ID exists, otherwise dev mode
	if _, err := os.Stat(nextBuildID); err == nil {
		npmCmd = exec.CommandContext(ctx, "npm", "start", "--", "-p", dashboardPort)
		fmt.Printf("%s📦 Production mode%s\n", cyan, reset)
	} else {
		npmCmd = exec.CommandContext(ctx, "npm", "run", "dev", "--", "-p", dashboardPort)
		fmt.Printf("%s🔧 Development mode%s\n", yellow, reset)
	}
	npmCmd.Dir = dashboardPath
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr
	npmCmd.Env = append(os.Environ(),
		"PORT="+dashboardPort,
		"IDOPS_CLI_PATH="+execPath,
	)

	serverURL := fmt.Sprintf("http://localhost:%s", dashboardPort)
	ready := make(chan bool)
	npmErr := make(chan error, 1)

	// Start server FIRST
	fmt.Printf("⏳ Đang khởi động dashboard trên port %s...\n", dashboardPort)
	go func() {
		if err := npmCmd.Run(); err != nil && ctx.Err() == nil {
			npmErr <- err
		}
	}()

	// Wait 2s before polling to avoid hitting old server
	go func() {
		time.Sleep(2 * time.Second)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		timeout := time.After(60 * time.Second)

		for {
			select {
			case <-ticker.C:
				resp, err := http.Get(serverURL)
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode == 200 {
						ready <- true
						return
					}
				}
			case <-timeout:
				ready <- false
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	select {
	case success := <-ready:
		if !success {
			return fmt.Errorf("dashboard không khởi động được trong 60 giây")
		}
	case err := <-npmErr:
		return fmt.Errorf("dashboard server lỗi: %w", err)
	case <-ctx.Done():
		return nil
	}

	fmt.Println()
	fmt.Printf("%s✅ Dashboard is ready!%s\n", green, reset)
	fmt.Printf("%s🌐 Dashboard URL: %s%s\n", cyan, serverURL, reset)
	fmt.Println()

	// Open browser
	if !dashboardNoOpen {
		fmt.Println("🖥️  Opening browser...")
		if err := openBrowser(serverURL); err != nil {
			fmt.Printf("%s⚠️  Could not open browser: %v%s\n", yellow, err, reset)
			fmt.Printf("%sPlease open the URL manually%s\n", yellow, reset)
		}
	}

	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the dashboard")
	fmt.Println()

	// Wait for interrupt
	<-ctx.Done()
	return nil
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default: // linux and others
		return exec.Command("xdg-open", url).Start()
	}
}

// isPortFree checks if a port is available by trying to connect to it.
// Uses Dial instead of Listen because on Windows, Listen can succeed
// even when another process holds the port (dual-stack IPv4/IPv6).
func isPortFree(port string) bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 500*time.Millisecond)
	if err != nil {
		// Connection refused = nothing listening = port is free
		return true
	}
	conn.Close()
	// Something is listening = port is NOT free
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
