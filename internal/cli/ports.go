package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"text/tabwriter"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nhh0718/idops/internal/ports"
	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Scan and display listening ports",
	Long:  "List all listening ports with process info. Supports watch mode, JSON output, and filtering.",
	RunE:  runPorts,
}

var portsKillCmd = &cobra.Command{
	Use:   "kill <port>",
	Short: "Kill process listening on a port",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var port uint32
		if _, err := fmt.Sscanf(args[0], "%d", &port); err != nil {
			return fmt.Errorf("invalid port: %s", args[0])
		}
		return ports.KillByPort(port)
	},
}

func init() {
	f := portsCmd.Flags()
	f.Bool("watch", false, "continuously refresh port list")
	f.Duration("interval", 2*time.Second, "refresh interval for watch mode")
	f.Bool("json", false, "output as JSON and exit")
	f.Bool("plain", false, "output as plain text table and exit")
	f.String("port", "", "port range filter, e.g. 1024-65535 or 80")
	f.String("protocol", "", "protocol filter: tcp or udp")

	portsCmd.AddCommand(portsKillCmd)
	rootCmd.AddCommand(portsCmd)
}

func runPorts(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	watchMode, _ := f.GetBool("watch")
	interval, _ := f.GetDuration("interval")
	jsonMode, _ := f.GetBool("json")
	plainMode, _ := f.GetBool("plain")
	portFlag, _ := f.GetString("port")
	protocol, _ := f.GetString("protocol")

	opts, err := buildScanOptions(portFlag, protocol)
	if err != nil {
		return err
	}

	if jsonMode {
		return runJSON(opts)
	}
	if plainMode {
		return runPlain(cmd, opts)
	}

	// Default: launch TUI (watch mode is a TUI option).
	m := ports.NewTUIModel(opts, watchMode, interval)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// runJSON scans once and prints JSON to stdout.
func runJSON(opts ports.ScanOptions) error {
	infos, err := ports.Scan(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	out, err := json.MarshalIndent(infos, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

// runPlain scans once and prints a tab-aligned table to stdout.
func runPlain(cmd *cobra.Command, opts ports.ScanOptions) error {
	infos, err := ports.Scan(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROTO\tADDR\tPORT\tPID\tPROCESS\tSTATUS")
	for _, p := range infos {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\n",
			p.Protocol, p.LocalAddr, p.LocalPort, p.PID, p.ProcessName, p.Status)
	}
	return w.Flush()
}

// buildScanOptions parses CLI flag strings into a ScanOptions struct.
func buildScanOptions(portFlag, protocol string) (ports.ScanOptions, error) {
	opts := ports.ScanOptions{Protocol: protocol}
	if portFlag == "" {
		return opts, nil
	}

	var min, max uint32
	n, _ := fmt.Sscanf(portFlag, "%d-%d", &min, &max)
	if n == 2 {
		opts.MinPort, opts.MaxPort = min, max
		return opts, nil
	}
	// Single port.
	var single uint32
	if _, err := fmt.Sscanf(portFlag, "%d", &single); err != nil {
		return opts, fmt.Errorf("invalid --port value %q: use 80 or 1024-65535", portFlag)
	}
	opts.MinPort, opts.MaxPort = single, single
	return opts, nil
}
