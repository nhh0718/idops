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

	// Auto-find free port if requested port is busy
	if !isPortFree(dashboardPort) {
		if dashboardPort == "3000" {
			// Only auto-switch for default port
			newPort := findFreePort(3001, 3100)
			if newPort == "" {
				return fmt.Errorf("port %s is in use and no free port found (3001-3100)", dashboardPort)
			}
			fmt.Printf("%s⚠️  Port %s đang bị chiếm, chuyển sang port %s%s\n", yellow, dashboardPort, newPort, reset)
			dashboardPort = newPort
		} else {
			return fmt.Errorf("port %s is already in use. Try a different port with --port", dashboardPort)
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

	// Use production build if .next exists, otherwise dev mode
	var npmCmd *exec.Cmd
	nextBuildDir := filepath.Join(dashboardPath, ".next")
	if _, err := os.Stat(nextBuildDir); err == nil {
		// Production: npm start (uses pre-built .next)
		npmCmd = exec.CommandContext(ctx, "npm", "start", "--", "-p", dashboardPort)
		fmt.Printf("%s📦 Using production build%s\n", cyan, reset)
	} else {
		// Development: npm run dev (requires source + node_modules)
		npmCmd = exec.CommandContext(ctx, "npm", "run", "dev", "--", "-p", dashboardPort)
		fmt.Printf("%s🔧 Using development mode%s\n", yellow, reset)
	}
	npmCmd.Dir = dashboardPath
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr
	npmCmd.Env = append(os.Environ(), "IDOPS_CLI_PATH="+execPath)

	// Wait for server to be ready
	serverURL := fmt.Sprintf("http://localhost:%s", dashboardPort)
	ready := make(chan bool)

	go func() {
		// Poll for server readiness
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		timeout := time.After(30 * time.Second)

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

	// Start server in background
	go func() {
		if err := npmCmd.Run(); err != nil && ctx.Err() == nil {
			fmt.Fprintf(os.Stderr, "Dashboard server error: %v\n", err)
		}
	}()

	// Wait for server to be ready
	fmt.Printf("⏳ Waiting for dashboard to start on port %s...\n", dashboardPort)

	select {
	case success := <-ready:
		if !success {
			return fmt.Errorf("dashboard failed to start within 30 seconds")
		}
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

// isPortFree checks if a port is available to listen on.
func isPortFree(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
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
