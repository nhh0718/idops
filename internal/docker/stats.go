package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// CollectStats performs a one-shot stats collection for a container.
func CollectStats(ctx context.Context, cli *client.Client, containerID string) (*StatsSnapshot, error) {
	resp, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, fmt.Errorf("stats for %s: %w", containerID, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read stats body: %w", err)
	}

	var raw container.StatsResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse stats: %w", err)
	}

	return parseSnapshot(&raw), nil
}

// parseSnapshot converts the raw Docker stats response into a StatsSnapshot.
func parseSnapshot(raw *container.StatsResponse) *StatsSnapshot {
	snap := &StatsSnapshot{}

	// CPU percent: (delta_container / delta_system) * num_cpus * 100
	cpuDelta := float64(raw.CPUStats.CPUUsage.TotalUsage) - float64(raw.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(raw.CPUStats.SystemUsage) - float64(raw.PreCPUStats.SystemUsage)
	numCPUs := float64(raw.CPUStats.OnlineCPUs)
	if numCPUs == 0 {
		numCPUs = float64(len(raw.CPUStats.CPUUsage.PercpuUsage))
	}
	if sysDelta > 0 && cpuDelta > 0 {
		snap.CPUPercent = (cpuDelta / sysDelta) * numCPUs * 100.0
	}

	// Memory
	snap.MemUsage = raw.MemoryStats.Usage
	snap.MemLimit = raw.MemoryStats.Limit
	if snap.MemLimit > 0 {
		snap.MemPercent = float64(snap.MemUsage) / float64(snap.MemLimit) * 100.0
	}

	// Network: sum all interfaces
	for _, iface := range raw.Networks {
		snap.NetIn += iface.RxBytes
		snap.NetOut += iface.TxBytes
	}

	return snap
}
