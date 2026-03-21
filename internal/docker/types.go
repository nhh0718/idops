package docker

import "time"

// ContainerInfo holds display-ready container metadata.
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Status  string
	State   string
	Ports   string
	Created time.Time
	Stats   *StatsSnapshot
}

// StatsSnapshot is a single-point-in-time resource usage snapshot.
type StatsSnapshot struct {
	CPUPercent float64
	MemPercent float64
	MemUsage   uint64
	MemLimit   uint64
	NetIn      uint64
	NetOut     uint64
}
