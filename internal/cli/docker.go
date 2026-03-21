package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nhh0718/idops/internal/docker"
	"github.com/spf13/cobra"
)

var dockerJSON bool

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Interactive Docker container dashboard",
	Long:  "Launch an interactive TUI dashboard showing all containers with live stats.",
	RunE:  runDockerDash,
}

var dockerLogsCmd = &cobra.Command{
	Use:   "logs <container>",
	Short: "View logs for a container",
	Args:  cobra.ExactArgs(1),
	RunE:  runDockerLogs,
}

func init() {
	dockerCmd.Flags().BoolVar(&dockerJSON, "json", false, "Output container list as JSON snapshot and exit")
	dockerCmd.AddCommand(dockerLogsCmd)
	rootCmd.AddCommand(dockerCmd)
}

func runDockerDash(cmd *cobra.Command, args []string) error {
	cli, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("cannot connect to Docker daemon: %w", err)
	}
	defer cli.Close()

	if dockerJSON {
		return printDockerJSON(cli)
	}

	m := docker.NewDashboard(cli)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("dashboard error: %w", err)
	}
	return nil
}

func printDockerJSON(cli interface{ Close() error }) error {
	// Re-open a proper typed client for listing
	dockerCli, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	containers, err := docker.ListContainers(ctx, dockerCli)
	if err != nil {
		return fmt.Errorf("list containers: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(containers)
}

func runDockerLogs(cmd *cobra.Command, args []string) error {
	containerID := args[0]

	cli, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("cannot connect to Docker daemon: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()
	lm, err := docker.NewLogViewer(ctx, cli, containerID, containerID)
	if err != nil {
		return fmt.Errorf("load logs: %w", err)
	}

	p := tea.NewProgram(lm, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("log viewer error: %w", err)
	}
	return nil
}
