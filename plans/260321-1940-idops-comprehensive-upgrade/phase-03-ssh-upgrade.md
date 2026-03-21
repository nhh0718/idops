---
phase: 3
title: "SSH Manager TUI Upgrade"
status: pending
effort: 3h
priority: P1
---

# Phase 03: SSH Manager TUI Upgrade

## Bugs to Fix

1. **Test results never displayed** — `m.testResults` populated but not rendered in list
2. **Status message clears immediately** — overwritten on next render cycle
3. **No form input validation** — duplicate names, non-numeric ports accepted
4. **Delete confirmation too vague** — only shows name, not user@host
5. **Connect from TUI missing** — `c` key defined but doesn't work properly

## Features to Add

1. **Test result indicators** — green checkmark / red X next to each host in list
2. **Connect from TUI** (`c` key) — exec ssh to selected host (re-exec like menu)
3. **Persistent status messages** — show for 3s with auto-clear
4. **Port validation** — reject non-numeric, range 1-65535
5. **Duplicate name check** — warn if host name already exists when adding

## UX Improvements

1. List item format: `hostname (user@host:port) [OK/FAIL]`
2. Delete confirm: "Delete myserver (root@192.168.1.1:22)? [y/N]"
3. Help footer: show all keys `a add  e edit  d delete  c connect  t test  q quit`
4. Empty state: "No SSH hosts in ~/.ssh/config"
5. After add/edit: show success message

## Files to Modify
- `internal/ssh/tui.go` — test indicators, status, validation, connect
- `internal/ssh/connection.go` — no changes needed
- `internal/ssh/config.go` — add duplicate name check helper

## Implementation Steps
1. Add `testResults map[string]bool` rendering in list view
2. Add status message with auto-clear timer
3. Add port/name validation in form submit
4. Fix delete confirmation to show full host info
5. Implement connect via re-exec (same pattern as menu)
6. Update help footer
