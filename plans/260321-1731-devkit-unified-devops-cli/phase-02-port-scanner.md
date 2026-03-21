---
phase: 2
title: "Port Scanner - idops ports"
status: pending
effort: 2d
priority: P1
depends_on: [1]
---

# Phase 02: Port Scanner - `idops ports`

## Context Links
- [Plan Overview](plan.md) | [Phase 01](phase-01-project-setup.md)
- gopsutil docs: https://pkg.go.dev/github.com/shirou/gopsutil/v3

## Overview
Scan tat ca listening ports tren may, hien thi process name/PID/user/protocol trong TUI table. Ho tro watch mode (auto-refresh), filter, sort, va JSON output.

## Key Insights
- `gopsutil/v3/net` cung cap `ConnectionsStat` voi PID, local/remote addr, status
- `gopsutil/v3/process` lay process name, user tu PID
- Tren Linux doc `/proc/net/tcp` truc tiep, macOS dung `lsof`
- Bubble Tea `tea.Tick` cho watch mode (refresh interval)

## Requirements

### Functional
- List tat ca listening TCP/UDP ports
- Hien thi: Port, Protocol, PID, Process Name, User, Local Address
- Sort theo bat ky column nao (click header hoac keybind)
- Filter: `--port 8000-9000`, `--protocol tcp`, search box trong TUI
- Watch mode: `idops ports --watch` (refresh moi 2s, configurable)
- JSON output: `idops ports --json`
- Non-TUI mode: `idops ports --plain` (simple table, pipe-friendly)

### Non-functional
- Scan time < 500ms cho ~100 ports
- Watch mode khong leak goroutines

## Architecture

```
internal/cli/ports.go          # Cobra command, flags
internal/ports/scanner.go      # Port scanning logic (gopsutil)
internal/ports/tui.go          # Bubble Tea model for port table
internal/ports/types.go        # PortInfo struct, sort/filter helpers
```

## Related Code Files

### Files to Create
- `internal/cli/ports.go`
- `internal/ports/scanner.go`
- `internal/ports/tui.go`
- `internal/ports/types.go`

### Dependencies
```bash
go get github.com/shirou/gopsutil/v3@latest
```

## Implementation Steps

### 1. Define Types (`internal/ports/types.go`)
```go
type PortInfo struct {
    Protocol    string // tcp, udp
    LocalAddr   string
    LocalPort   uint32
    RemoteAddr  string
    PID         int32
    ProcessName string
    User        string
    Status      string // LISTEN, ESTABLISHED, etc.
}

type SortField int
const (
    SortByPort SortField = iota
    SortByPID
    SortByProcess
    SortByProtocol
)
```

### 2. Implement Scanner (`internal/ports/scanner.go`)
```go
func Scan(ctx context.Context) ([]PortInfo, error) {
    conns, err := net.ConnectionsWithContext(ctx, "all")
    // Filter: chi lay LISTEN status
    // Voi moi connection, lookup process info:
    proc, _ := process.NewProcess(conn.Pid)
    name, _ := proc.Name()
    user, _ := proc.Username()
    // Return sorted by port
}
```
- Cache process lookups (nhieu ports cung 1 PID)
- Handle permission errors gracefully (non-root co the khong thay het)

### 3. Implement Cobra Command (`internal/cli/ports.go`)
```go
var portsCmd = &cobra.Command{
    Use:   "ports",
    Short: "Scan and display listening ports",
    RunE:  runPorts,
}

func init() {
    portsCmd.Flags().Bool("watch", false, "Watch mode with auto-refresh")
    portsCmd.Flags().Duration("interval", 2*time.Second, "Refresh interval")
    portsCmd.Flags().Bool("json", false, "JSON output")
    portsCmd.Flags().Bool("plain", false, "Plain table output")
    portsCmd.Flags().String("port", "", "Filter port range (e.g. 8000-9000)")
    portsCmd.Flags().String("protocol", "", "Filter protocol (tcp/udp)")
    rootCmd.AddCommand(portsCmd)
}
```
- Neu `--json`: scan 1 lan, marshal JSON, print, exit
- Neu `--plain`: scan 1 lan, print tabwriter table, exit
- Neu `--watch` hoac default: khoi dong TUI

### 4. Implement TUI Model (`internal/ports/tui.go`)
```go
type Model struct {
    table     table.Model
    ports     []PortInfo
    filter    textinput.Model
    sortField SortField
    sortAsc   bool
    watching  bool
    interval  time.Duration
    width     int
    height    int
}

func (m Model) Init() tea.Cmd {
    if m.watching {
        return tea.Tick(m.interval, func(t time.Time) tea.Msg { return tickMsg(t) })
    }
    return nil
}
```
- **Keys**: `q` quit, `s` change sort, `/` toggle filter, `r` manual refresh
- **Watch mode**: `tea.Tick` -> rescan -> update table rows
- **Filter**: textinput o tren, filter ports realtime khi gox

### 5. Wire Up & Register Command
- Register `portsCmd` trong `root.go` init
- Test: `go run ./cmd/idops ports`

## Todo List
- [ ] Install gopsutil dependency
- [ ] Define PortInfo type + sort/filter helpers
- [ ] Implement scanner with process lookup + caching
- [ ] Implement Cobra command with all flags
- [ ] Implement TUI model (table, filter, sort)
- [ ] Implement watch mode with tea.Tick
- [ ] Implement JSON output mode
- [ ] Implement plain text output mode
- [ ] Handle permission errors (non-root)
- [ ] Test tren Linux va macOS

## Success Criteria
- `idops ports` hien TUI table voi tat ca listening ports
- Sort/filter hoat dong trong TUI
- `idops ports --watch` auto-refresh
- `idops ports --json` output valid JSON
- `idops ports --plain` pipe-friendly output
- Khong crash khi khong co quyen root (show warning + partial results)

## Risk Assessment
- **Permission**: Non-root khong thay het ports tren Linux -> hien warning, suggest sudo
- **Performance**: Nhieu ports + process lookup cham -> cache PID->name map
- **Platform**: gopsutil API khac nhau giua Linux/macOS -> test ca hai

## Security Considerations
- Khong can elevated privileges de chay, nhung results limited
- Khong expose sensitive info trong JSON output (chi port/process metadata)

## Next Steps
- Sau khi xong, architecture da validate -> cac phase khac follow same pattern
- Shared table component co the duoc reuse cho Phase 3 (Docker) va Phase 4 (SSH)
