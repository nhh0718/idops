package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	internalssh "github.com/nhh0718/idops/internal/ssh"
)

// exportJSON writes hosts as a JSON array to w.
func exportJSON(hosts []internalssh.SSHHost, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(hosts)
}

// importFromJSON tries to parse data as a JSON array of SSHHost.
func importFromJSON(data []byte) ([]internalssh.SSHHost, error) {
	var hosts []internalssh.SSHHost
	if err := json.Unmarshal(data, &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

var sshJSON bool

func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.AddCommand(sshConnectCmd)
	sshCmd.AddCommand(sshTestCmd)
	sshCmd.AddCommand(sshExportCmd)
	sshCmd.AddCommand(sshImportCmd)
	sshCmd.AddCommand(sshListCmd)

	sshCmd.Flags().BoolVar(&sshJSON, "json", false, "Output SSH hosts as JSON")
	sshTestCmd.Flags().Bool("json", false, "Output test results as JSON")

	sshExportCmd.Flags().Bool("raw", false, "Xuất raw ssh config text thay vì JSON")
	sshExportCmd.Flags().String("file", "", "Lưu vào file thay vì stdout")

	sshImportCmd.Flags().Bool("dry-run", false, "Xem trước mà không ghi vào config")
}

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Manage SSH hosts via TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := defaultSSHConfigPath()
		hosts, err := internalssh.LoadConfig(path)
		if err != nil {
			return fmt.Errorf("loading ssh config: %w", err)
		}

		if sshJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(hosts)
		}

		model, err := internalssh.NewTUIModel(path)
		if err != nil {
			return fmt.Errorf("loading ssh config: %w", err)
		}
		p := tea.NewProgram(model, tea.WithAltScreen())
		_, err = p.Run()
		return err
	},
}

var sshConnectCmd = &cobra.Command{
	Use:   "connect <host>",
	Short: "Connect to an SSH host by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostName := args[0]
		hosts, err := internalssh.LoadConfig(defaultSSHConfigPath())
		if err != nil {
			return err
		}
		for _, h := range hosts {
			if h.Name == hostName {
				return execSSH(h)
			}
		}
		return fmt.Errorf("host %q not found in ssh config", hostName)
	},
}

var sshTestCmd = &cobra.Command{
	Use:   "test [host]",
	Short: "Test SSH host connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts, err := internalssh.LoadConfig(defaultSSHConfigPath())
		if err != nil {
			return err
		}
		if len(args) > 0 {
			name := args[0]
			filtered := hosts[:0]
			for _, h := range hosts {
				if h.Name == name {
					filtered = append(filtered, h)
				}
			}
			if len(filtered) == 0 {
				return fmt.Errorf("host %q not found", name)
			}
			hosts = filtered
		}
		results := internalssh.TestAllConnections(hosts, internalssh.DefaultTimeout)

		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			for i := range results {
				results[i].FillJSON()
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(results)
		}

		for _, r := range results {
			if r.Success {
				fmt.Printf("  [OK]  %-20s  %v\n", r.Host.Name, r.Latency.Round(1))
			} else {
				fmt.Printf(" [FAIL] %-20s  %v\n", r.Host.Name, r.Error)
			}
		}
		return nil
	},
}

var sshExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Xuất danh sách SSH host ra JSON (mặc định) hoặc raw config",
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, _ := cmd.Flags().GetBool("raw")
		filePath, _ := cmd.Flags().GetString("file")

		if raw {
			// Legacy: print raw ssh config text
			path := defaultSSHConfigPath()
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(os.Stdout, f)
			return err
		}

		hosts, err := internalssh.LoadConfig(defaultSSHConfigPath())
		if err != nil {
			return fmt.Errorf("đọc ssh config thất bại: %w", err)
		}

		var out io.Writer = os.Stdout
		var f *os.File
		if filePath != "" {
			f, err = os.Create(filePath)
			if err != nil {
				return fmt.Errorf("không thể tạo file %s: %w", filePath, err)
			}
			defer f.Close()
			out = f
		}

		if err := exportJSON(hosts, out); err != nil {
			return err
		}
		if filePath != "" {
			fmt.Printf("Đã export %d host(s) ra %s\n", len(hosts), filePath)
		}
		return nil
	},
}

var sshImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Nhập host từ file JSON hoặc ssh_config vào SSH config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		data, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("không thể đọc file %s: %w", srcPath, err)
		}

		// Auto-detect: try JSON first, then fall back to ssh_config format
		hosts, jsonErr := importFromJSON(data)
		if jsonErr != nil {
			hosts, err = internalssh.LoadConfig(srcPath)
			if err != nil {
				return fmt.Errorf("đọc file import thất bại (không phải JSON hoặc ssh_config hợp lệ): %w", err)
			}
		}

		if len(hosts) == 0 {
			fmt.Println("Không tìm thấy host nào trong file import.")
			return nil
		}

		if dryRun {
			fmt.Printf("[Dry-run] Sẽ import %d host(s):\n", len(hosts))
			for _, h := range hosts {
				fmt.Printf("  - %s (%s@%s)\n", h.Name, h.User, h.Hostname)
			}
			return nil
		}

		destPath := defaultSSHConfigPath()
		for _, h := range hosts {
			if err := internalssh.AddHost(destPath, h); err != nil {
				return fmt.Errorf("thêm host %q thất bại: %w", h.Name, err)
			}
			fmt.Printf("Đã import host: %s\n", h.Name)
		}
		fmt.Printf("Đã import %d host(s) từ %s\n", len(hosts), srcPath)
		return nil
	},
}

var sshListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all SSH hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts, err := internalssh.LoadConfig(defaultSSHConfigPath())
		if err != nil {
			return err
		}
		if len(hosts) == 0 {
			fmt.Println("No hosts found")
			return nil
		}
		for _, h := range hosts {
			fmt.Printf("  %s -> %s@%s", h.Name, h.User, h.Hostname)
			if h.Port != "" && h.Port != "22" {
				fmt.Printf(":%s", h.Port)
			}
			fmt.Println()
		}
		return nil
	},
}

// execSSH replaces the current process with an ssh invocation.
func execSSH(host internalssh.SSHHost) error {
	sshArgs := buildSSHArgs(host)
	sshBin, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh not found in PATH: %w", err)
	}
	// On Windows exec.Command is used since syscall.Exec is not available.
	if runtime.GOOS == "windows" {
		c := exec.Command(sshBin, sshArgs...)
		c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
		return c.Run()
	}
	return syscallExec(sshBin, append([]string{"ssh"}, sshArgs...))
}

func buildSSHArgs(h internalssh.SSHHost) []string {
	var args []string
	if h.Port != "" {
		args = append(args, "-p", h.Port)
	}
	if h.User != "" {
		args = append(args, "-l", h.User)
	}
	if h.IdentityFile != "" {
		args = append(args, "-i", h.IdentityFile)
	}
	if h.ProxyJump != "" {
		args = append(args, "-J", h.ProxyJump)
	}
	hostname := h.Hostname
	if hostname == "" {
		hostname = h.Name
	}
	args = append(args, hostname)
	return args
}

// defaultSSHConfigPath returns ~/.ssh/config.
func defaultSSHConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssh", "config")
}
