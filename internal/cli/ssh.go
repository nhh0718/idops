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
	Short: "Print the SSH config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := defaultSSHConfigPath()
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(os.Stdout, f)
		return err
	},
}

var sshImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Append hosts from a file into SSH config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		hosts, err := internalssh.LoadConfig(srcPath)
		if err != nil {
			return fmt.Errorf("reading import file: %w", err)
		}
		if len(hosts) == 0 {
			fmt.Println("No hosts found in import file.")
			return nil
		}
		destPath := defaultSSHConfigPath()
		for _, h := range hosts {
			if err := internalssh.AddHost(destPath, h); err != nil {
				return fmt.Errorf("adding host %q: %w", h.Name, err)
			}
		}
		fmt.Printf("Imported %d host(s) into %s\n", len(hosts), destPath)
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
