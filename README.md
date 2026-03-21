# idops

All-in-one DevOps CLI tool. Single binary, 5 tools: port scanner, Docker dashboard, SSH manager, env sync, nginx config generator.

Built with Go, Cobra, Bubble Tea.

## Install

### One-line install (recommended)

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/nhh0718/idops/main/install.sh | sh
```

```powershell
# Windows (PowerShell)
irm https://raw.githubusercontent.com/nhh0718/idops/main/install.ps1 | iex
```

### Go install (requires Go 1.22+)

```bash
go install github.com/nhh0718/idops/cmd/idops@latest
```

### Manual download

Grab the latest from [Releases](https://github.com/nhh0718/idops/releases).

### Build from source

```bash
git clone https://github.com/nhh0718/idops.git
cd idops
make build
# Binary at bin/idops
```

## Commands

### `idops ports` - Port Scanner

Scan all listening ports with process info.

```bash
idops ports              # Interactive TUI with sort/filter
idops ports --watch      # Auto-refresh every 2s
idops ports --json       # JSON output
idops ports --plain      # Simple table (pipe-friendly)
idops ports --port 8000-9000 --protocol tcp
```

**TUI keys:** `q` quit, `s` sort, `/` filter, `r` refresh

### `idops docker` - Docker Dashboard

Real-time container dashboard with CPU/memory stats.

```bash
idops docker             # Interactive TUI dashboard
idops docker --json      # JSON snapshot
idops docker logs <name> # View container logs
```

**TUI keys:** `q` quit, `s` start, `x` stop, `r` restart, `l` logs, `/` filter

### `idops ssh` - SSH Config Manager

Manage `~/.ssh/config` hosts with TUI.

```bash
idops ssh                # Interactive TUI manager
idops ssh connect myhost # Quick connect
idops ssh test           # Test all connections
idops ssh test myhost    # Test single host
idops ssh export         # Print config to stdout
idops ssh import file    # Import hosts from file
```

**TUI keys:** `q` quit, `a` add, `e` edit, `d` delete, `c` connect, `t` test

### `idops env` - Environment File Manager

Compare, sync, and validate `.env` files.

```bash
idops env compare                      # .env.example vs .env
idops env compare --target .env.prod   # Compare with specific file
idops env sync                         # Interactive fill missing vars
idops env init                         # Generate .env from .env.example
idops env validate                     # Check format issues
idops env show                         # Display with masked secrets
```

### `idops nginx` - Nginx Config Generator

Interactive nginx config generator with 5 templates.

```bash
idops nginx                        # Interactive generator
idops nginx --output site.conf     # Save directly to file
idops nginx validate               # Run nginx -t
idops nginx apply site.conf        # Enable + validate + reload
idops nginx list                   # List configs
```

**Templates:** Reverse Proxy, Static Site, PHP-FPM, Load Balancer, WebSocket

## Configuration

Config file: `~/.config/idops/config.yaml`

```yaml
docker:
  host: unix:///var/run/docker.sock

ssh:
  config_path: ~/.ssh/config

nginx:
  config_dir: /etc/nginx/sites-available
  templates_dir: /etc/idops/templates

env:
  default_env_file: .env
```

Override with flag: `idops --config /path/to/config.yaml`

## Development

```bash
# Run directly
make dev ARGS="ports --plain"

# Build
make build

# Test
make test

# Lint
make lint

# Cross-compile all platforms
bash scripts/build.sh v0.1.0 all

# Install to /usr/local/bin
sudo make install

# Clean
make clean
```

## Release

### Quick release (recommended)

```bash
bash scripts/release.sh v0.1.0
```

This tags the commit and pushes. GitHub Actions automatically builds binaries for Linux/macOS/Windows (amd64 + arm64) and creates a GitHub Release.

### Manual release

```bash
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions handles the rest
```

### Local release (test)

```bash
goreleaser release --snapshot --clean
```

## CI/CD

| Workflow | Trigger | Action |
|----------|---------|--------|
| CI | Push/PR to `main` | Build + test |
| Release | Push `v*` tag | GoReleaser: cross-compile + GitHub Release |

## Tech Stack

- **Go 1.22+** - Language
- **Cobra + Viper** - CLI framework + config
- **Bubble Tea + Bubbles** - Terminal UI
- **Lip Gloss** - TUI styling
- **gopsutil** - System info (ports)
- **Docker SDK** - Container management
- **GoReleaser** - Cross-platform releases

## License

MIT
