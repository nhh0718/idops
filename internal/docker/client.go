package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// NewClient creates a Docker client from environment variables with API version negotiation.
func NewClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("docker client init: %w", err)
	}
	return cli, nil
}

// ListContainers returns all containers (running and stopped).
func ListContainers(ctx context.Context, cli *client.Client) ([]ContainerInfo, error) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		ports := make([]string, 0, len(c.Ports))
		for _, p := range c.Ports {
			if p.PublicPort != 0 {
				ports = append(ports, fmt.Sprintf("%s:%d->%d/%s", p.IP, p.PublicPort, p.PrivatePort, p.Type))
			} else {
				ports = append(ports, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
			}
		}

		result = append(result, ContainerInfo{
			ID:      c.ID[:12],
			Name:    name,
			Image:   c.Image,
			Status:  c.Status,
			State:   c.State,
			Ports:   strings.Join(ports, ", "),
			Created: time.Unix(c.Created, 0),
		})
	}
	return result, nil
}

// StartContainer starts a stopped container by ID.
func StartContainer(ctx context.Context, cli *client.Client, id string) error {
	if err := cli.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		return fmt.Errorf("start container %s: %w", id, err)
	}
	return nil
}

// StopContainer stops a running container by ID.
func StopContainer(ctx context.Context, cli *client.Client, id string) error {
	if err := cli.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
		return fmt.Errorf("stop container %s: %w", id, err)
	}
	return nil
}

// RestartContainer restarts a container by ID.
func RestartContainer(ctx context.Context, cli *client.Client, id string) error {
	if err := cli.ContainerRestart(ctx, id, container.StopOptions{}); err != nil {
		return fmt.Errorf("restart container %s: %w", id, err)
	}
	return nil
}

// RemoveContainer force-removes a container by ID.
func RemoveContainer(ctx context.Context, cli *client.Client, id string) error {
	if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("remove container %s: %w", id, err)
	}
	return nil
}

// InspectContainer returns detailed container info as formatted string.
func InspectContainer(ctx context.Context, cli *client.Client, id string) (string, error) {
	info, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		return "", fmt.Errorf("inspect container %s: %w", id, err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name:    %s\n", info.Name))
	sb.WriteString(fmt.Sprintf("ID:      %s\n", info.ID[:12]))
	sb.WriteString(fmt.Sprintf("Image:   %s\n", info.Config.Image))
	sb.WriteString(fmt.Sprintf("State:   %s\n", info.State.Status))
	sb.WriteString(fmt.Sprintf("Started: %s\n", info.State.StartedAt))
	if info.State.Health != nil {
		sb.WriteString(fmt.Sprintf("Health:  %s\n", info.State.Health.Status))
	}
	sb.WriteString(fmt.Sprintf("IP:      %s\n", info.NetworkSettings.IPAddress))

	if len(info.Config.Env) > 0 {
		sb.WriteString("\nEnvironment:\n")
		for _, e := range info.Config.Env {
			sb.WriteString(fmt.Sprintf("  %s\n", e))
		}
	}
	return sb.String(), nil
}
