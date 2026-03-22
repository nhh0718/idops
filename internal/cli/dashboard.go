package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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
	// dashboard-dist = flattened standalone (from release)
	// dashboard = source code (for dev)
	searchNames := []string{"dashboard-dist", "dashboard"}
	searchDirs := []string{
		realExecDir,                                       // next to real binary
		execDir,                                            // next to binary (or symlink)
		filepath.Join(realExecDir, ".."),                   // parent of real binary
		idopsHome,                                          // ~/.idops/
		filepath.Join(homeDir, ".local", "share", "idops"), // XDG data
		".",                                                // CWD
	}

	var dashboardPath string
	var isStandalone bool
	var searchedPaths []string

	for _, dir := range searchDirs {
		for _, name := range searchNames {
			path := filepath.Join(dir, name)
			searchedPaths = append(searchedPaths, path)

			// Check for standalone marker or server.js
			if _, err := os.Stat(filepath.Join(path, "server.js")); err == nil {
				dashboardPath = path
				isStandalone = true
				break
			}
			// Check for source code (dev mode)
			if _, err := os.Stat(filepath.Join(path, "app", "page.tsx")); err == nil {
				dashboardPath = path
				break
			}
		}
		if dashboardPath != "" {
			break
		}
	}

	if dashboardPath == "" {
		return fmt.Errorf("dashboard not found. Searched:\n  %s\nRun 'idops update' or reinstall: irm https://raw.githubusercontent.com/nhh0718/idops/main/install.ps1 | iex", strings.Join(searchedPaths, "\n  "))
	}
	_ = isStandalone

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

	// Check if node is available
	nodeBin, err := exec.LookPath("node")
	if err != nil {
		return fmt.Errorf("node not found in PATH. Please install Node.js")
	}

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

	// Determine server mode
	serverJSPath := filepath.Join(dashboardPath, "server.js")
	var serverCmd *exec.Cmd

	if _, err := os.Stat(serverJSPath); err == nil {
		// Standalone mode: server.js at dashboard root (from release)
		fmt.Printf("%s📦 Production mode%s\n", cyan, reset)
		serverCmd = exec.CommandContext(ctx, nodeBin, "server.js")
		serverCmd.Dir = dashboardPath
		serverCmd.Env = append(os.Environ(),
			"PORT="+dashboardPort,
			"HOSTNAME=0.0.0.0",
			"IDOPS_CLI_PATH="+execPath,
		)
	} else {
		// Dev mode: source code, use npm
		npmBin, npmErr := exec.LookPath("npm")
		if npmErr != nil {
			return fmt.Errorf("npm not found. Install Node.js or use release build")
		}
		fmt.Printf("%s🔧 Development mode%s\n", yellow, reset)
		serverCmd = exec.CommandContext(ctx, npmBin, "run", "dev", "--", "-p", dashboardPort)
		serverCmd.Dir = dashboardPath
		serverCmd.Env = append(os.Environ(),
			"PORT="+dashboardPort,
			"IDOPS_CLI_PATH="+execPath,
		)
	}

	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	serverURL := fmt.Sprintf("http://localhost:%s", dashboardPort)
	ready := make(chan bool)
	npmErr := make(chan error, 1)

	// Start server FIRST
	fmt.Printf("⏳ Đang khởi động dashboard trên port %s...\n", dashboardPort)
	go func() {
		if err := serverCmd.Run(); err != nil && ctx.Err() == nil {
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

