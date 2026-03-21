---
phase: 1
title: "Project Setup & Core Architecture"
status: pending
effort: 2d
priority: P1
depends_on: []
---

# Phase 01: Project Setup & Core Architecture

## Context Links
- [Plan Overview](plan.md)
- [Research: Go CLI/TUI Libraries](research/)
- [Research: SSH/Nginx/Env Patterns](research/)

## Overview
Khoi tao Go module, setup Cobra root command, Viper config system, shared TUI theme, va build toolchain (Makefile + GoReleaser). Day la foundation cho tat ca phases sau.

## Key Insights
- Cobra + Viper tich hop san: Viper bind duoc voi Cobra flags
- Lip Gloss v2 dung adaptive colors (tu detect terminal capability)
- GoReleaser ho tro cross-compile Linux/macOS/Windows tu 1 config

## Requirements

### Functional
- `idops` command chay duoc, hien help text voi list subcommands
- `idops version` hien version + build info
- `idops --config <path>` override config file path
- Config load tu `~/.config/idops/config.yaml`

### Non-functional
- Build time < 30s
- Binary size < 20MB (stripped)
- Go 1.22+ required

## Architecture

```
cmd/idops/main.go          # os.Args -> cli.Execute()
internal/cli/root.go        # cobra.Command root, Viper init
internal/cli/version.go     # version subcommand
internal/ui/theme.go        # Lip Gloss styles, color palette
internal/ui/table.go        # Reusable Bubble Tea table wrapper
internal/ui/helpers.go      # Common TUI utilities (spinner, error display)
internal/config/config.go   # Config struct + load/save
```

## Related Code Files

### Files to Create
- `cmd/idops/main.go`
- `internal/cli/root.go`
- `internal/cli/version.go`
- `internal/config/config.go`
- `internal/ui/theme.go`
- `internal/ui/table.go`
- `internal/ui/helpers.go`
- `go.mod`, `Makefile`, `.goreleaser.yaml`, `.gitignore`

## Implementation Steps

### 1. Init Go Module
```bash
mkdir -p cmd/idops internal/{cli,config,ui} templates/nginx
go mod init github.com/<user>/idops
```

### 2. Install Core Dependencies
```bash
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss/v2@latest
```

### 3. Create Root Command (`internal/cli/root.go`)
```go
package cli

var (
    cfgFile string
    rootCmd = &cobra.Command{
        Use:   "idops",
        Short: "Unified DevOps CLI tool",
    }
)

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, _ := os.UserHomeDir()
        viper.AddConfigPath(filepath.Join(home, ".config", "idops"))
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
    }
    viper.AutomaticEnv()
    viper.ReadInConfig() // ignore error if no config
}
```

### 4. Create Config Package (`internal/config/config.go`)
```go
type Config struct {
    Docker  DockerConfig  `yaml:"docker"`
    SSH     SSHConfig     `yaml:"ssh"`
    Nginx   NginxConfig   `yaml:"nginx"`
}
type DockerConfig struct {
    Host string `yaml:"host" mapstructure:"host"`
}
type SSHConfig struct {
    ConfigPath string `yaml:"config_path" mapstructure:"config_path"`
}
type NginxConfig struct {
    ConfigDir string `yaml:"config_dir" mapstructure:"config_dir"`
}
```

### 5. Create Shared TUI Theme (`internal/ui/theme.go`)
```go
var (
    Primary   = lipgloss.Color("#7C3AED")
    Success   = lipgloss.Color("#10B981")
    Warning   = lipgloss.Color("#F59E0B")
    Error     = lipgloss.Color("#EF4444")
    Muted     = lipgloss.Color("#6B7280")

    TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(Primary)
    ErrorStyle = lipgloss.NewStyle().Foreground(Error)
)
```

### 6. Create Table Wrapper (`internal/ui/table.go`)
- Wrap `bubbles/table` voi default styles tu theme
- Helper: `NewTable(columns []table.Column, rows []table.Row) table.Model`

### 7. Create Main Entry Point (`cmd/idops/main.go`)
```go
package main

import (
    "os"
    "github.com/<user>/idops/internal/cli"
)

func main() {
    if err := cli.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### 8. Create Build Toolchain
**Makefile:**
- `build`: `go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/idops ./cmd/idops`
- `dev`: `go run ./cmd/idops`
- `test`: `go test ./...`
- `lint`: `golangci-lint run`

**GoReleaser:** Linux amd64/arm64, macOS amd64/arm64, Windows amd64. Homebrew tap, checksum.

### 9. Version Command
- Inject version/commit/date via `-ldflags` at build time
- `idops version` -> `idops v0.1.0 (abc1234) built 2026-03-21`

## Todo List
- [ ] Init Go module + directory structure
- [ ] Install dependencies
- [ ] Implement root command with Viper config
- [ ] Implement config package
- [ ] Implement shared TUI theme
- [ ] Implement table wrapper
- [ ] Implement helpers (spinner, error display)
- [ ] Create main.go entry point
- [ ] Create Makefile
- [ ] Create .goreleaser.yaml
- [ ] Create .gitignore
- [ ] Verify `go build` + `idops --help` works

## Success Criteria
- `idops` binary builds thanh cong
- `idops --help` hien danh sach subcommands
- `idops version` hien version info
- Config loads tu default path hoac --config flag
- `go test ./...` pass (basic tests)

## Risk Assessment
- **Lip Gloss v2 breaking changes**: Pin version trong go.mod
- **Go version compatibility**: Require 1.22+ trong go.mod

## Security Considerations
- Config file permissions: warn neu `config.yaml` la world-readable
- Khong store secrets trong config file

## Next Steps
- Phase 2 (Port Scanner) co the bat dau ngay khi Phase 1 xong
- Shared TUI components se duoc refine khi implement cac phases sau
