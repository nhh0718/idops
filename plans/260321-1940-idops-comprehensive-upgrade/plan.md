---
title: "idops Comprehensive TUI/UX Upgrade"
status: pending
priority: P1
effort: 3d
created: 2026-03-21
---

# idops Comprehensive TUI/UX Upgrade

Fix all broken interactions, missing feedback, silent errors, and UX issues across all 5 tools.

## Audit Summary

| Category | Count | Severity |
|---|---|---|
| Silent error ignoring | 8 | HIGH |
| Missing validation | 6 | HIGH |
| Missing user feedback | 7 | MEDIUM |
| UX clarity problems | 5 | MEDIUM |
| Incomplete features | 4 | MEDIUM |
| Logic bugs | 4 | HIGH |

## Phases

| # | Phase | Effort | Status | File |
|---|-------|--------|--------|------|
| 1 | Port Scanner upgrade | 4h | pending | [phase-01](phase-01-ports-upgrade.md) |
| 2 | Docker Dashboard upgrade | 4h | pending | [phase-02](phase-02-docker-upgrade.md) |
| 3 | SSH Manager upgrade | 3h | pending | [phase-03](phase-03-ssh-upgrade.md) |
| 4 | Env & Nginx fixes | 2h | pending | [phase-04](phase-04-env-nginx-fixes.md) |

## Key Principles
- Every action must show feedback (success/error)
- Every destructive action needs confirmation
- Status messages auto-clear after 3s (via tick)
- Errors surfaced to user, never silently ignored
- Input validated before processing
