---
phase: 2
title: "Docker Dashboard TUI Upgrade"
status: pending
effort: 4h
priority: P1
---

# Phase 02: Docker Dashboard TUI Upgrade

## Bugs to Fix

1. **Container actions ignore errors** — start/stop/restart use `_ =`, must show result
2. **Color state display unused** — `stateColored` defined but assigned to `_`, wire it up
3. **Filter causes panic** — `m.filtered[:0]` on nil slice can panic
4. **Tick stops on error** — auto-refresh halts permanently after single fetch error
5. **Log viewer hardcoded dimensions** — 120x30 ignores actual terminal size

## Features to Add

1. **Confirm before stop/restart** — "Stop container X? [y/N]" like ports kill
2. **Action feedback** — "Started container X" green message with auto-clear
3. **Docker daemon error** — clear message: "Docker not running. Start Docker Desktop and retry."
4. **Remove container** (`d` key) — with double-confirm since destructive
5. **Inspect container** (`i` key) — show container details in viewport

## UX Improvements

1. Color-coded states: Running=green, Exited=red, Paused=yellow, Created=blue
2. Status bar: total containers, running count, last refresh
3. Log viewer: show container name in title, add scroll hints
4. Help footer: all keys visible
5. Empty state: "No containers found. Is Docker running?"

## Files to Modify
- `internal/docker/tui.go` — confirm mode, action feedback, colors, status bar
- `internal/docker/tui_logs.go` — dynamic sizing, better title, scroll hints
- `internal/docker/client.go` — already has Remove/Inspect (added earlier)
- `internal/cli/docker.go` — better daemon error messages

## Implementation Steps
1. Initialize `m.filtered` properly (not nil)
2. Wire up stateColored function to table rows
3. Add confirm mode (like ports TUI pattern)
4. Show action results with auto-clear status
5. Fix tick to restart even after error
6. Fix log viewer dimensions from WindowSizeMsg
7. Add remove (d) and inspect (i) keybinds
