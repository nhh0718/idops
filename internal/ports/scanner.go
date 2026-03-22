package ports

import (
	"context"
	"fmt"
	"sort"
	"strings"

	gopsnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ScanOptions holds optional filters applied during scanning.
type ScanOptions struct {
	MinPort    uint32 // 0 = no lower bound
	MaxPort    uint32 // 0 = no upper bound
	Protocol   string // "tcp", "udp", or "" for all
	ShowSystem bool   // false = hide system/unknown processes (default behavior)
}

// isSystemProcess returns true for OS-level or unresolvable processes.
func isSystemProcess(pid int32, name string, port uint32) bool {
	// PID 0 or 4 are kernel/system
	if pid == 0 || pid == 4 {
		return true
	}
	// Unresolvable process name = no permission = system process
	if name == "-" || name == "" {
		return true
	}
	// Well-known system services
	sysNames := map[string]bool{
		"system": true, "svchost.exe": true, "lsass.exe": true,
		"services.exe": true, "wininit.exe": true, "csrss.exe": true,
		"smss.exe": true, "spoolsv.exe": true, "dnsmasq": true,
		"systemd-resolve": true, "avahi-daemon": true,
	}
	if sysNames[strings.ToLower(name)] {
		return true
	}
	return false
}

// Scan returns all LISTEN connections visible to the current user.
// Permission errors on individual processes are silently skipped.
func Scan(ctx context.Context, opts ScanOptions) ([]PortInfo, error) {
	conns, err := gopsnet.ConnectionsWithContext(ctx, "all")
	if err != nil {
		return nil, err
	}

	// Cache process name lookups to avoid repeated syscalls for same PID.
	procCache := make(map[int32]string)

	var results []PortInfo
	for _, c := range conns {
		// Keep only LISTEN / NONE (UDP has no state).
		if c.Status != "LISTEN" && c.Status != "NONE" && c.Status != "" {
			continue
		}

		// gopsutil Type: 1=SOCK_STREAM(tcp), 2=SOCK_DGRAM(udp).
		var proto string
		switch c.Type {
		case 2:
			proto = "udp"
		default:
			proto = "tcp"
		}

		port := uint32(c.Laddr.Port)

		// Apply protocol filter.
		if opts.Protocol != "" && !strings.EqualFold(proto, opts.Protocol) {
			continue
		}
		// Apply port range filter.
		if opts.MinPort > 0 && port < opts.MinPort {
			continue
		}
		if opts.MaxPort > 0 && port > opts.MaxPort {
			continue
		}

		name := processName(c.Pid, procCache)

		// Filter out system processes unless explicitly requested
		if !opts.ShowSystem && isSystemProcess(c.Pid, name, port) {
			continue
		}

		results = append(results, PortInfo{
			Protocol:    proto,
			LocalAddr:   c.Laddr.IP,
			LocalPort:   port,
			RemoteAddr:  c.Raddr.IP,
			PID:         c.Pid,
			ProcessName: name,
			Status:      c.Status,
		})
	}

	// Deduplicate by port+pid+protocol
	seen := make(map[string]bool)
	var deduped []PortInfo
	for _, r := range results {
		key := fmt.Sprintf("%s:%d:%d", r.Protocol, r.LocalPort, r.PID)
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, r)
		}
	}

	sort.Slice(deduped, func(i, j int) bool {
		return deduped[i].LocalPort < deduped[j].LocalPort
	})

	return deduped, nil
}

// processName resolves a PID to its executable name using a local cache.
// Returns "-" on permission error or when process no longer exists.
func processName(pid int32, cache map[int32]string) string {
	if pid == 0 {
		return "-"
	}
	if name, ok := cache[pid]; ok {
		return name
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		cache[pid] = "-"
		return "-"
	}
	name, err := p.Name()
	if err != nil {
		cache[pid] = "-"
		return "-"
	}
	cache[pid] = name
	return name
}

// SortPortInfos sorts a slice in-place by the given field.
func SortPortInfos(infos []PortInfo, by SortField) {
	sort.Slice(infos, func(i, j int) bool {
		a, b := infos[i], infos[j]
		switch by {
		case SortByPID:
			return a.PID < b.PID
		case SortByProcess:
			return strings.ToLower(a.ProcessName) < strings.ToLower(b.ProcessName)
		case SortByProtocol:
			return a.Protocol < b.Protocol
		default: // SortByPort
			return a.LocalPort < b.LocalPort
		}
	})
}
