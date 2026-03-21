---
title: "idops - Unified DevOps CLI Tool"
description: "Go CLI gom 5 tools DevOps vao 1 binary: ports, docker, ssh, env, nginx"
status: pending
priority: P1
effort: 2w
branch: main
tags: [go, cli, devops, tui, open-source, idops]
created: 2026-03-21
---

# idops - Unified DevOps CLI Tool

Single Go binary gom 5 DevOps sub-tools: port scanner, docker dashboard, ssh manager, env sync, nginx generator. Full-screen TUI voi Bubble Tea + Cobra CLI framework.

**Project location:** E:/Code-Fun/idops

## Tech Stack
- **Go 1.22+**, Cobra + Viper, Bubble Tea + Bubbles + Lip Gloss v2
- gopsutil/v3, Docker SDK, kevinburke/ssh_config, joho/godotenv, charmbracelet/huh
- GoReleaser (cross-compile + release)

## Architecture
```
cmd/idops/main.go -> internal/cli/ (cobra commands)
                  -> internal/{ports,docker,ssh,env,nginx}/ (business logic + TUI)
                  -> internal/ui/ (shared theme, table, helpers)
                  -> templates/nginx/ (.tmpl files)
```
Config: `~/.config/idops/config.yaml`

## Phases

| # | Phase | Effort | Status | File |
|---|-------|--------|--------|------|
| 1 | Project Setup & Core Architecture | 2d | pending | [phase-01](phase-01-project-setup.md) |
| 2 | Port Scanner (`idops ports`) | 2d | pending | [phase-02](phase-02-port-scanner.md) |
| 3 | Docker Dashboard (`idops docker`) | 3d | pending | [phase-03](phase-03-docker-dashboard.md) |
| 4 | SSH Manager (`idops ssh`) | 2d | pending | [phase-04](phase-04-ssh-manager.md) |
| 5 | Env Sync (`idops env`) | 1.5d | pending | [phase-05](phase-05-env-sync.md) |
| 6 | Nginx Generator (`idops nginx`) | 2.5d | pending | [phase-06](phase-06-nginx-generator.md) |

## Dependencies
- Phase 1 MUST complete first (foundation)
- Phases 2-6 can parallel sau khi Phase 1 xong
- Phase 2 la don gian nhat -> lam truoc de validate architecture
- Phase 3 phuc tap nhat -> can nhieu thoi gian nhat

## Key Decisions
- **internal/ over pkg/**: Private packages, khong expose public API
- **Bubble Tea per subcommand**: Moi tool co TUI model rieng, share theme qua internal/ui
- **Cobra subcommands**: `idops <tool> [flags]` pattern
- **No database**: Config file only (YAML + SSH config file)

## Risk Summary
- Docker SDK requires Docker daemon running -> graceful error handling
- SSH config parsing edge cases (Include directives, wildcards)
- Nginx reload needs sudo -> document privilege requirements
- Cross-platform: Linux primary, macOS secondary, Windows limited
