package docker

import "time"

// ContainerInfo holds display-ready container metadata.
type ContainerInfo struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Image   string         `json:"image"`
	Status  string         `json:"status"`
	State   string         `json:"state"`
	Ports   string         `json:"ports"`
	Created time.Time      `json:"created"`
	Stats   *StatsSnapshot `json:"stats,omitempty"`
}

// StatsSnapshot is a single-point-in-time resource usage snapshot.
type StatsSnapshot struct {
	CPUPercent float64 `json:"cpu"`
	MemPercent float64 `json:"memory"`
	MemUsage   uint64  `json:"memUsage"`
	MemLimit   uint64  `json:"memLimit"`
	NetIn      uint64  `json:"netIn"`
	NetOut     uint64  `json:"netOut"`
}
