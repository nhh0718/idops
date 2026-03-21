---
phase: 4
title: "SSH Manager - idops ssh"
status: pending
effort: 2d
priority: P2
depends_on: [1]
---

# Phase 04: SSH Manager - `idops ssh`

## Context Links
- [Plan Overview](plan.md) | [Phase 01](phase-01-project-setup.md)
- kevinburke/ssh_config: https://pkg.go.dev/github.com/kevinburke/ssh_config
- x/crypto/ssh: https://pkg.go.dev/golang.org/x/crypto/ssh

## Overview
SSH config manager doc/ghi ~/.ssh/config. TUI interactive de list/add/edit/delete hosts. Quick connect, connection health test, import/export.

## Key Insights
- `kevinburke/ssh_config` parse SSH config giu nguyen comments va formatting
- `x/crypto/ssh` cho connection testing (Dial voi timeout)
- SSH config co the co `Include` directives -> can handle recursive
- `ssh_config.Decode` tra ve AST-like structure, co the modify va write lai

## Requirements

### Functional
- List tat ca SSH hosts tu ~/.ssh/config trong TUI table
- Add host moi qua interactive form (hostname, port, user, key, proxy)
- Edit host: select tu list -> form voi pre-filled values
- Delete host voi confirmation
- Quick connect: `idops ssh connect <host>` -> exec ssh
- Test connection: `idops ssh test [host]` -> check reachability voi timeout
- Import/export: `idops ssh export > backup.ssh_config`, `idops ssh import <file>`

### Non-functional
- Preserve comments va formatting khi edit
- Connection test timeout: 5s default, configurable
- Backup config truoc khi modify

## Architecture

```
internal/cli/ssh.go              # Cobra commands (ssh, ssh connect, ssh test)
internal/ssh/config.go           # SSH config parser/writer
internal/ssh/connection.go       # Connection testing
internal/ssh/tui.go              # TUI list/form model
internal/ssh/types.go            # SSHHost struct
```

## Related Code Files

### Files to Create
- `internal/cli/ssh.go`
- `internal/ssh/config.go`
- `internal/ssh/connection.go`
- `internal/ssh/tui.go`
- `internal/ssh/types.go`

### Dependencies
```bash
go get github.com/kevinburke/ssh_config@latest
go get golang.org/x/crypto@latest
```

## Implementation Steps

### 1. Define Types (`internal/ssh/types.go`)
```go
type SSHHost struct {
    Name         string
    Hostname     string
    Port         string
    User         string
    IdentityFile string
    ProxyJump    string
    Options      map[string]string // other options
}
```

### 2. Config Parser (`internal/ssh/config.go`)
```go
func LoadConfig(path string) ([]SSHHost, error) {
    f, _ := os.Open(path)
    cfg, _ := ssh_config.Decode(f)
    var hosts []SSHHost
    for _, host := range cfg.Hosts {
        // Skip wildcards (* patterns)
        // Extract known fields: Hostname, Port, User, IdentityFile, ProxyJump
    }
    return hosts, nil
}

func AddHost(path string, host SSHHost) error {
    // Read file, append new Host block, write back
    // Format: Host <name>\n  Hostname <hostname>\n  Port <port>\n...
}

func UpdateHost(path string, oldName string, host SSHHost) error {
    // Parse, find host by name, replace fields, write back
}

func DeleteHost(path string, name string) error {
    // Parse, remove host block, write back
}

func BackupConfig(path string) error {
    // Copy to path + ".bak." + timestamp
}
```
- **CRITICAL**: Goi `BackupConfig` truoc moi write operation

### 3. Connection Tester (`internal/ssh/connection.go`)
```go
func TestConnection(host SSHHost, timeout time.Duration) error {
    addr := net.JoinHostPort(host.Hostname, host.Port)
    conn, err := net.DialTimeout("tcp", addr, timeout)
    if err != nil {
        return fmt.Errorf("cannot reach %s: %w", addr, err)
    }
    conn.Close()
    return nil
}

func TestAllConnections(hosts []SSHHost, timeout time.Duration) []TestResult {
    // Parallel test voi semaphore (max 10 concurrent)
    // Return TestResult{Host, Status, Latency, Error}
}
```

### 4. TUI Model (`internal/ssh/tui.go`)
```go
type Model struct {
    hosts    []SSHHost
    list     list.Model  // bubbles/list
    mode     viewMode    // list, add, edit, confirm_delete
    form     huh.Form    // charmbracelet/huh for add/edit
    selected int
    configPath string
}

type viewMode int
const (
    modeList viewMode = iota
    modeAdd
    modeEdit
    modeConfirmDelete
)
```
- **List view**: Bubbles list voi host name, hostname, user
- **Add/Edit**: `huh.Form` voi fields: Name, Hostname, Port, User, IdentityFile, ProxyJump
- **Keys**: `a` add, `e` edit, `d` delete, `c` connect, `t` test, `q` quit
- **Status bar**: Hien connection test result (green check / red x)

### 5. Cobra Commands (`internal/cli/ssh.go`)
```go
var sshCmd = &cobra.Command{Use: "ssh", Short: "SSH config manager", RunE: runSSH}
var sshConnectCmd = &cobra.Command{
    Use: "connect <host>", Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // exec.Command("ssh", args[0]).Run() voi Stdin/Stdout/Stderr
    },
}
var sshTestCmd = &cobra.Command{
    Use: "test [host]",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Test 1 host hoac tat ca hosts
    },
}
var sshExportCmd = &cobra.Command{Use: "export", RunE: runExport}
var sshImportCmd = &cobra.Command{Use: "import <file>", Args: cobra.ExactArgs(1), RunE: runImport}
```

### 6. Quick Connect
- `idops ssh connect myserver` -> `syscall.Exec` de replace process voi ssh
- Tren Linux/macOS: `syscall.Exec("/usr/bin/ssh", []string{"ssh", host}, os.Environ())`
- Fallback: `exec.Command("ssh", host)` voi piped I/O

## Todo List
- [ ] Install ssh_config + x/crypto dependencies
- [ ] Define SSHHost type
- [ ] Implement config parser (load/add/update/delete)
- [ ] Implement auto-backup truoc moi write
- [ ] Implement connection tester (single + batch)
- [ ] Implement TUI list view
- [ ] Implement add/edit form (huh library)
- [ ] Implement delete confirmation
- [ ] Implement quick connect (syscall.Exec)
- [ ] Implement export/import commands
- [ ] Implement Cobra commands + register
- [ ] Test voi real SSH config file
- [ ] Handle edge cases: empty config, Include directives

## Success Criteria
- `idops ssh` hien TUI list cua SSH hosts
- Add/edit/delete host thanh cong, config file duoc update dung
- Comments va formatting duoc preserve
- `idops ssh connect <host>` ket noi SSH thanh cong
- `idops ssh test` hien connection status cho tat ca hosts
- Backup tu dong truoc moi modification

## Risk Assessment
- **Config corruption**: Backup truoc moi write, validate sau khi write
- **Include directives**: kevinburke/ssh_config co the khong handle het -> document limitation
- **Wildcard hosts**: Skip `Host *` va pattern hosts trong list view
- **Key permissions**: Warn neu IdentityFile co wrong permissions (khong phai 600)

## Security Considerations
- Khong display private key content
- Backup files co cung permissions voi original
- Quick connect dung syscall.Exec (khong inject commands)
- Validate host name input (khong cho shell metacharacters)

## Next Steps
- Phase 5 (Env Sync) la independent, co the parallel
- Consider adding SSH key management trong future version (YAGNI)
