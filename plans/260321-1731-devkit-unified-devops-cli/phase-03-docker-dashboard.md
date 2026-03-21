---
phase: 3
title: "Docker Dashboard - idops docker"
status: pending
effort: 3d
priority: P1
depends_on: [1]
---

# Phase 03: Docker Dashboard - `idops docker`

## Context Links
- [Plan Overview](plan.md) | [Phase 01](phase-01-project-setup.md)
- Docker SDK: https://pkg.go.dev/github.com/docker/docker/client

## Overview
Realtime Docker container dashboard trong TUI. Hien thi CPU/RAM/Network stats per container, ho tro sort/filter/search, va interactive actions (start/stop/restart/logs).

## Key Insights
- Docker SDK `ContainerStats` tra ve stream (io.ReadCloser) -> decode JSON stream
- `ContainerList` voi `All: true` lay ca stopped containers
- Stats calculation: CPU% = delta(container CPU) / delta(system CPU) * numCPUs * 100
- Log streaming: `ContainerLogs` voi `Follow: true` + `Tail: "100"`

## Requirements

### Functional
- List tat ca containers (running + stopped) voi status, image, ports, uptime
- Realtime stats: CPU%, Memory (used/limit/%), Network I/O
- Sort theo name, CPU, memory, status
- Search/filter by name hoac image
- Actions: start, stop, restart (voi confirm dialog)
- View logs: `idops docker logs <container>` hoac interactive select
- JSON output: `idops docker --json` (snapshot, khong stream)

### Non-functional
- Stats update moi 1s
- Handle Docker daemon unavailable gracefully
- Memory efficient: khong buffer toan bo log history

## Architecture

```
internal/cli/docker.go            # Cobra command + subcommands
internal/docker/client.go         # Docker SDK wrapper
internal/docker/stats.go          # Stats collection + calculation
internal/docker/tui.go            # TUI dashboard model
internal/docker/tui-logs.go       # Log viewer TUI
internal/docker/types.go          # ContainerInfo, StatsSnapshot structs
```

## Related Code Files

### Files to Create
- `internal/cli/docker.go`
- `internal/docker/client.go`
- `internal/docker/stats.go`
- `internal/docker/tui.go`
- `internal/docker/tui-logs.go`
- `internal/docker/types.go`

### Dependencies
```bash
go get github.com/docker/docker@latest
go get github.com/docker/go-connections@latest
```

## Implementation Steps

### 1. Define Types (`internal/docker/types.go`)
```go
type ContainerInfo struct {
    ID        string
    Name      string
    Image     string
    Status    string // running, exited, paused
    State     string
    Ports     string
    Created   time.Time
    Stats     *StatsSnapshot
}

type StatsSnapshot struct {
    CPUPercent float64
    MemUsage   uint64
    MemLimit   uint64
    MemPercent float64
    NetIn      uint64
    NetOut     uint64
}
```

### 2. Docker Client Wrapper (`internal/docker/client.go`)
```go
func NewClient() (*client.Client, error) {
    // Tu detect: DOCKER_HOST env, default unix socket, TCP
    return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func ListContainers(ctx context.Context, cli *client.Client) ([]ContainerInfo, error) {
    containers, _ := cli.ContainerList(ctx, container.ListOptions{All: true})
    // Map to ContainerInfo
}

func StartContainer(ctx context.Context, cli *client.Client, id string) error
func StopContainer(ctx context.Context, cli *client.Client, id string) error
func RestartContainer(ctx context.Context, cli *client.Client, id string) error
```

### 3. Stats Collector (`internal/docker/stats.go`)
```go
func CollectStats(ctx context.Context, cli *client.Client, containerID string) (*StatsSnapshot, error) {
    resp, _ := cli.ContainerStats(ctx, containerID, false) // one-shot
    defer resp.Body.Close()
    var stats container.StatsResponse
    json.NewDecoder(resp.Body).Decode(&stats)
    return calculateStats(&stats), nil
}

func calculateStats(stats *container.StatsResponse) *StatsSnapshot {
    // CPU%: (delta container usage / delta system usage) * CPU count * 100
    cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
    systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
    cpuPercent := (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
    // Memory: stats.MemoryStats.Usage, stats.MemoryStats.Limit
    // Network: sum all interfaces
}
```

### 4. TUI Dashboard (`internal/docker/tui.go`)
```go
type Model struct {
    client     *client.Client
    containers []ContainerInfo
    table      table.Model
    filter     textinput.Model
    selected   int
    err        error
    width, height int
}
```
- **Layout**: Header (title + filter) | Table (containers) | Footer (keybinds)
- **Keys**: `s` start, `x` stop, `r` restart, `l` logs, `/` filter, `q` quit
- **Confirm dialog**: Truoc khi stop/restart, hien confirm prompt
- **Auto-refresh**: `tea.Tick` moi 1s, collect stats cho tat ca running containers
- **Color coding**: Running=green, Exited=red, Paused=yellow

### 5. Log Viewer (`internal/docker/tui-logs.go`)
```go
type LogModel struct {
    viewport viewport.Model
    containerID string
    client   *client.Client
}
```
- Dung `bubbles/viewport` de scroll logs
- Load last 100 lines, stream new lines
- `q` hoac `Esc` quay lai dashboard

### 6. Cobra Command (`internal/cli/docker.go`)
```go
var dockerCmd = &cobra.Command{
    Use:   "docker",
    Short: "Docker container dashboard",
    RunE:  runDocker,
}
var dockerLogsCmd = &cobra.Command{
    Use:   "logs [container]",
    Short: "View container logs",
    Args:  cobra.ExactArgs(1),
    RunE:  runDockerLogs,
}
```
- `idops docker` -> TUI dashboard
- `idops docker logs <name>` -> log viewer
- `idops docker --json` -> JSON snapshot

## Todo List
- [ ] Install Docker SDK dependencies
- [ ] Define container/stats types
- [ ] Implement Docker client wrapper
- [ ] Implement stats collector with CPU/memory calculation
- [ ] Implement TUI dashboard model
- [ ] Implement container actions (start/stop/restart)
- [ ] Implement confirm dialog cho destructive actions
- [ ] Implement log viewer TUI
- [ ] Implement Cobra commands (docker, docker logs)
- [ ] Handle Docker daemon unavailable error
- [ ] Implement JSON output mode
- [ ] Color coding cho container states
- [ ] Test voi nhieu containers (10+)

## Success Criteria
- `idops docker` hien TUI dashboard voi realtime stats
- CPU/Memory% chinh xac (so sanh voi `docker stats`)
- Start/stop/restart containers thanh cong
- Log viewer hien thi va scroll duoc
- Graceful error khi Docker daemon khong chay
- `idops docker --json` output valid JSON

## Risk Assessment
- **Docker daemon down**: Check connection khi init, hien clear error message
- **Nhieu containers**: Stats collection cho 20+ containers co the cham -> parallel goroutines voi semaphore
- **Docker API version**: `WithAPIVersionNegotiation()` handle backward compat
- **Permission**: Docker socket can group membership -> suggest adding user to docker group

## Security Considerations
- Container actions (stop/restart) can confirm truoc khi execute
- Khong expose Docker socket path trong output
- Log viewer co the chua sensitive data -> warn user

## Next Steps
- TUI patterns tu Phase 2 va 3 duoc reuse cho Phase 4 (SSH)
- Consider adding container exec trong future version (YAGNI for now)
