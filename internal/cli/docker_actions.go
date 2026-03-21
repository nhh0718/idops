package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/nhh0718/idops/internal/docker"
	"github.com/spf13/cobra"
)

var dockerStartCmd = &cobra.Command{
	Use:   "start <container>",
	Short: "Start a stopped container",
	Args:  cobra.ExactArgs(1),
	RunE:  runDockerAction("start"),
}

var dockerStopCmd = &cobra.Command{
	Use:   "stop <container>",
	Short: "Stop a running container",
	Args:  cobra.ExactArgs(1),
	RunE:  runDockerAction("stop"),
}

var dockerRestartCmd = &cobra.Command{
	Use:   "restart <container>",
	Short: "Restart a container",
	Args:  cobra.ExactArgs(1),
	RunE:  runDockerAction("restart"),
}

var dockerRmCmd = &cobra.Command{
	Use:   "rm <container>",
	Short: "Remove a container (force)",
	Args:  cobra.ExactArgs(1),
	RunE:  runDockerAction("remove"),
}

func init() {
	dockerCmd.AddCommand(dockerStartCmd, dockerStopCmd, dockerRestartCmd, dockerRmCmd)
}

// runDockerAction returns a RunE function for the given Docker action.
func runDockerAction(action string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		containerID := args[0]

		cli, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("cannot connect to Docker: %w", err)
		}
		defer cli.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		switch action {
		case "start":
			err = docker.StartContainer(ctx, cli, containerID)
		case "stop":
			err = docker.StopContainer(ctx, cli, containerID)
		case "restart":
			err = docker.RestartContainer(ctx, cli, containerID)
		case "remove":
			err = docker.RemoveContainer(ctx, cli, containerID)
		default:
			return fmt.Errorf("unknown action: %s", action)
		}

		if err != nil {
			return err
		}
		fmt.Printf("Container %s: %s OK\n", containerID, action)
		return nil
	}
}
