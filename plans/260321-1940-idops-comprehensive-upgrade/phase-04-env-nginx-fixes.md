---
phase: 4
title: "Env & Nginx CLI Fixes"
status: pending
effort: 2h
priority: P2
---

# Phase 04: Env & Nginx CLI Fixes

## Env Bugs

1. **envShow missing --json flag** — flag never defined on command, always false
2. **Compare output confusing** — "+" for missing is counterintuitive, use clear labels
3. **No error context** — Parse errors don't say which file failed

## Env Fixes
1. Add `--json` flag to `envShowCmd`
2. Change compare output: "MISSING (in .env.example but not .env):" and "EXTRA (in .env but not .env.example):"
3. Wrap Parse errors with file path context

## Nginx Bugs

1. **strconv.Atoi errors silently ignored** — port becomes 0 on invalid input
2. **No input validation** — domain, cert paths not checked
3. **Invalid choice doesn't retry** — exits immediately
4. **Load balancer method not validated** — accepts any string

## Nginx Fixes
1. Wrap all strconv.Atoi with error check + re-prompt
2. Validate domain format (at least non-empty)
3. Validate cert/key file existence when SSL enabled
4. Retry on invalid menu choice (loop)
5. Validate LB method is one of {round-robin, least_conn, ip_hash}
6. Add `promptInt` helper that loops until valid number

## Files to Modify
- `internal/cli/env.go` — json flag, compare labels, error context
- `internal/cli/nginx.go` — validation, retry loops, promptInt helper

## Implementation Steps
1. Add `--json` flag to envShowCmd in init()
2. Rewrite compare output with clear labels
3. Create `promptInt(reader, label, default, min, max)` helper
4. Replace all raw strconv.Atoi with promptInt
5. Add domain/cert validation
6. Add LB method validation with select-style prompt
7. Wrap menu choice in retry loop
