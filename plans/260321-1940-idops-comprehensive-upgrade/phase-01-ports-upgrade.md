---
phase: 1
title: "Port Scanner TUI Upgrade"
status: pending
effort: 4h
priority: P1
---

# Phase 01: Port Scanner TUI Upgrade

## Bugs to Fix

1. **Browser open errors silently ignored** — `_ = OpenInBrowser()` must show error
2. **CopyToClipboard non-functional** — function exists but never works, fix or remove
3. **Status message persists forever** — add auto-clear via tick after 3s
4. **Kill confirmation leaks keypresses to table** — block table update during confirm mode
5. **No port range validation in KillByPort** — validate 1-65535

## Features to Add

1. **Copy address** (`c` key) — copy `localhost:PORT` to clipboard (use Go exec for platform clipboard)
2. **Status bar** — bottom bar showing total ports, filtered count, last refresh time
3. **Color-coded protocols** — TCP=blue, UDP=yellow
4. **Permission warning** — detect if non-admin/non-root, show warning banner

## UX Improvements

1. Help footer: show all keys including `k kill`, `o open`, `c copy`
2. Kill success: show green "Killed PID X on port Y" for 3s then clear
3. Kill fail: show red error for 5s
4. Browser open: show "Opened http://localhost:PORT" for 3s
5. Empty state: "No listening ports found" instead of empty table

## Files to Modify
- `internal/ports/tui.go` — status bar, auto-clear, color, confirm fix
- `internal/ports/actions.go` — fix clipboard, add validation
- `internal/ports/types.go` — no changes needed

## Implementation Steps
1. Add `statusTimer` field to TUIModel, clear status after tick
2. Fix confirm mode to fully block table updates
3. Implement real clipboard (exec pbcopy/xclip/clip.exe)
4. Add color to protocol column
5. Add port count + last refresh to footer
6. Add empty state message
7. Fix OpenInBrowser error handling
