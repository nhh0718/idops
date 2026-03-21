package docker

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	dockerclient "github.com/docker/docker/client"
)

// runAction executes a container action and returns a Cmd that produces actionResultMsg.
func runAction(cli *dockerclient.Client, action, id, name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		var verb string

		switch action {
		case "s":
			err = StartContainer(ctx, cli, id)
			verb = "Started"
		case "x":
			err = StopContainer(ctx, cli, id)
			verb = "Stopped"
		case "r":
			err = RestartContainer(ctx, cli, id)
			verb = "Restarted"
		case "d":
			err = RemoveContainer(ctx, cli, id)
			verb = "Removed"
		default:
			return actionResultMsg{text: fmt.Sprintf("Unknown action: %s", action), isErr: true}
		}

		if err != nil {
			return actionResultMsg{text: fmt.Sprintf("Error: %s", err.Error()), isErr: true}
		}
		return actionResultMsg{text: fmt.Sprintf("%s container %q", verb, name), isErr: false}
	}
}
