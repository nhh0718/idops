# Nghiên Cứu: Go CLI/TUI Frameworks cho DevKit DevOps Tool
**Ngày:** 2025-03-21 | **Tác giả:** Researcher Agent

## 1. Cobra - Framework CLI Tiêu Chuẩn (spf13/cobra)

### Điểm Mạnh
- De facto standard trong cộng đồng Go (dùng bởi Kubernetes, Docker, Hugo)
- Hỗ trợ nested subcommands với organization logic rõ ràng
- Auto-generate shell completion (bash/zsh/fish/powershell)
- Auto-generate man pages via Mango + GoReleaser
- Command grouping trong help output
- Cobra-cli tool bootstrap application structure nhanh

### Best Practices
- Commands = actions, Args = things, Flags = modifiers (syntax dạng sentence)
- Persistent flags cho global options, local flags cho command-specific
- PersistentPreRun hook để validate args trước execution
- Group commands: `AddGroup()` → group subcommands trong help

### Config & Flags
**Viper Integration** (spf13/viper):
- Universal config reader: files, env vars, remote sources
- Precedence: flags > env vars > config file > defaults
- Binding automatic env vars → flags via cobraflags module
- Hỗ trợ 12-factor app pattern

## 2. TUI Libraries - Charmbracelet Ecosystem

### Bubble Tea (Main Framework)
**Architecture:** Model-Update-View (Elm-inspired)
- **Init()**: return initial command
- **Update(msg)**: handle events, change state
- **View()**: render UI based model

### Best Practices
- I/O operations → Commands (return Msg)
- Use LogToFile vì terminal occupied bởi TUI
- Commands giữ program straightforward & testable
- Real-time refresh: use polling commands or message subscriptions

### Bubbles - Component Library
- **Table**: sortable, selectable rows
- **TextInput**: form input with validation
- **List**: scrollable item selection
- **Viewport**: scrollable content
- **Progress**: progress bars, spinners

### Styling - Lip Gloss + Glamour
**Lip Gloss (v2 - 2025):**
- CSS-like declarative styling
- Table & list sub-packages
- Terminal capability detection

**Glamour:**
- Markdown renderer stylesheet-based
- Multiple built-in styles
- Smart column width, proper wrapping
- Complement Lip Gloss for markdown content

## 3. Port/Process Scanning

### Gopsutil (shirou/gopsutil/v3)
- Pure Go (no CGO) → cross-compile friendly
- Process package: POSIX systems (Linux, macOS, FreeBSD)
- Reads `/proc/net/tcp`, `/proc/net/tcp6` on Linux
- Methods: ConnectionsMax, ProcessConnections
- Simpler than raw syscall calls

### Cross-Platform Strategy
- Linux: read /proc filesystem via gopsutil
- macOS: lsof via subprocess (shirou/gopsutil handles)
- Alternative: exec `ss`/`netstat` with parsing (simpler but subprocess overhead)

## 4. Docker SDK for Go

### Official Docker Client (github.com/docker/docker/client)
```go
ContainerStats(ctx, containerID, stream=true) → StatsResponseReader
- stream=true: continuous stat updates
- stream=false: single snapshot
- Caller responsible for closing io.ReadCloser
```

### Streaming Pattern
- Return container.StatsResponseReader (io.ReadCloser)
- Read JSON from reader: container stats encoded as JSON stream
- Close reader when done

### Container Lifecycle
- Start/Stop/Restart: simple context + containerID
- Logs streaming: similar pattern to stats
- Image pull/push: full image management

## 5. Competitor Analysis

### LazyDocker (29K⭐, gocui-based)
- Simple TUI: containers, images, volumes, services
- Real-time logs, CPU/memory, ASCII charts
- Auto-detect docker-compose.yml
- Full lifecycle management (start/stop/restart)
- Resource-efficient vs GUI tools

### ctop
- Focused: container metrics only
- Top-like interface
- Compressed summary format
- Real-time CPU/memory/network I/O
- Support Docker + runC

### Architecture Patterns
- Event-driven update loops
- Async stat collection (goroutines)
- Config file support (.lazydocker/config.yml)
- Keybinding customization

## 6. DevKit Architecture Recommendation

### Technology Stack
| Component | Choice | Reason |
|-----------|--------|--------|
| CLI Framework | Cobra | nested commands, auto-completion, mature |
| Config | Viper + Cobra | multi-source, precedence rules |
| TUI | Bubble Tea | predictable state model, composable |
| Components | Bubbles | table, list, input reusable |
| Styling | Lip Gloss v2 | CSS-like, built-in layouts |
| Port Scanning | gopsutil/v3 | pure Go, cross-platform, no CGO |
| Docker API | Official SDK | official support, streaming |
| Markdown | Glamour | styled rendering, terminal-friendly |

### Project Structure (YAGNI)
```
devkit/
├── cmd/
│   ├── devkit/main.go
│   ├── commands/
│   │   ├── docker.go (container, stats, logs)
│   │   ├── process.go (port-scan, connections)
│   │   ├── config.go (init, edit)
│   │   └── help.go
├── pkg/
│   ├── docker/client.go
│   ├── process/scanner.go
│   ├── config/manager.go
│   └── ui/
│       ├── tui.go (Bubble Tea app)
│       ├── components/table.go
│       └── styles/theme.go
├── docs/
└── go.mod
```

## 7. Unresolved Questions
1. DevKit có cần interactive TUI full-screen hay simple CLI with tabular output?
2. Real-time refresh interval recommendation? (stats, logs polling)
3. Config file location & format preference? (~/.devkit/config.yaml?)
4. Export format needs? (JSON, CSV logging output?)

## Sources
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Bubble Tea TUI](https://github.com/charmbracelet/bubbletea)
- [Charmbracelet Components](https://github.com/charmbracelet/bubbles)
- [Viper Config Management](https://github.com/spf13/viper)
- [Docker Go SDK](https://pkg.go.dev/github.com/docker/docker/client)
- [gopsutil Library](https://github.com/shirou/gopsutil)
- [LazyDocker Reference](https://github.com/jesseduffield/lazydocker)
- [Lip Gloss Styling](https://github.com/charmbracelet/lipgloss)
- [Glamour Markdown](https://github.com/charmbracelet/glamour)
