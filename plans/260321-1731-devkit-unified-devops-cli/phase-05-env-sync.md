---
phase: 5
title: "Env Sync - idops env"
status: pending
effort: 1.5d
priority: P2
depends_on: [1]
---

# Phase 05: Env Sync - `idops env`

## Context Links
- [Plan Overview](plan.md) | [Phase 01](phase-01-project-setup.md)
- joho/godotenv: https://pkg.go.dev/github.com/joho/godotenv
- charmbracelet/huh: https://pkg.go.dev/github.com/charmbracelet/huh

## Overview
Compare .env.example vs .env, detect missing/extra variables, interactive prompt de fill missing values, validate format. Ho tro nhieu environments (.env.dev, .env.prod, .env.staging).

## Key Insights
- `godotenv` parse .env files thanh `map[string]string`, giu order
- `huh` library cho interactive forms (text input, select, confirm)
- .env files co the co comments (#), empty lines, quotes
- Common issues: trailing spaces, duplicate keys, empty values

## Requirements

### Functional
- Compare: `idops env compare` -> diff giua .env.example va .env
- Sync: `idops env sync` -> interactive prompt cho missing vars
- Validate: `idops env validate` -> check format issues
- Multi-env: `idops env compare --target .env.prod`
- Init: `idops env init` -> generate .env tu .env.example (full interactive)
- Show: `idops env show` -> masked display (.env values voi ****)

### Non-functional
- Khong modify .env.example (read-only source of truth)
- Preserve comments va order trong .env khi sync
- Sensitive detection: auto-mask keys chua PASSWORD, SECRET, KEY, TOKEN

## Architecture

```
internal/cli/env.go             # Cobra commands
internal/env/parser.go          # .env file parser (wrap godotenv + extras)
internal/env/comparator.go      # Diff logic between env files
internal/env/validator.go       # Format validation rules
internal/env/sync.go            # Interactive sync with huh forms
```

## Related Code Files

### Files to Create
- `internal/cli/env.go`
- `internal/env/parser.go`
- `internal/env/comparator.go`
- `internal/env/validator.go`
- `internal/env/sync.go`

### Dependencies
```bash
go get github.com/joho/godotenv@latest
go get github.com/charmbracelet/huh@latest
```

## Implementation Steps

### 1. Parser (`internal/env/parser.go`)
```go
type EnvFile struct {
    Path     string
    Vars     map[string]string
    Order    []string            // preserve insertion order
    Comments map[string]string   // key -> comment above it
}

func Parse(path string) (*EnvFile, error) {
    // godotenv.Read(path) cho vars
    // Custom parse cho comments + order (line by line)
}

func (e *EnvFile) Write(path string) error {
    // Write back voi preserved order va comments
}
```

### 2. Comparator (`internal/env/comparator.go`)
```go
type DiffResult struct {
    Missing []string            // in example but not in target
    Extra   []string            // in target but not in example
    Changed []DiffEntry         // different values (optional compare)
}

type DiffEntry struct {
    Key      string
    Example  string
    Target   string
}

func Compare(example, target *EnvFile) *DiffResult {
    // Set difference operations
}
```

### 3. Validator (`internal/env/validator.go`)
```go
type ValidationIssue struct {
    Line    int
    Key     string
    Type    IssueType // EmptyValue, DuplicateKey, TrailingSpace, InvalidFormat
    Message string
}

func Validate(envFile *EnvFile, rawContent string) []ValidationIssue {
    // Check:
    // - Empty values (KEY= without value)
    // - Duplicate keys (same key defined twice)
    // - Trailing spaces in values
    // - Invalid format (no = sign, spaces in key)
    // - Unquoted values with spaces
}
```

### 4. Interactive Sync (`internal/env/sync.go`)
```go
func SyncInteractive(example, target *EnvFile) error {
    diff := Compare(example, target)
    if len(diff.Missing) == 0 {
        fmt.Println("All variables are in sync!")
        return nil
    }

    // Build huh.Form voi 1 field per missing var
    var fields []huh.Field
    for _, key := range diff.Missing {
        defaultVal := example.Vars[key]
        input := huh.NewInput().
            Title(key).
            Description("Default: " + defaultVal).
            Value(&values[key])
        fields = append(fields, input)
    }

    form := huh.NewForm(huh.NewGroup(fields...))
    form.Run()

    // Add new vars to target, write back
}
```

### 5. Cobra Commands (`internal/cli/env.go`)
```go
var envCmd = &cobra.Command{Use: "env", Short: "Environment file manager"}

var envCompareCmd = &cobra.Command{
    Use:   "compare",
    Short: "Compare .env.example with .env",
    RunE: func(cmd *cobra.Command, args []string) error {
        source, _ := cmd.Flags().GetString("source")  // default .env.example
        target, _ := cmd.Flags().GetString("target")  // default .env
        // Parse both, compare, print diff table
    },
}

var envSyncCmd = &cobra.Command{Use: "sync", Short: "Interactive sync missing vars"}
var envValidateCmd = &cobra.Command{Use: "validate", Short: "Validate .env format"}
var envInitCmd = &cobra.Command{Use: "init", Short: "Generate .env from .env.example"}
var envShowCmd = &cobra.Command{Use: "show", Short: "Display .env with masked secrets"}
```

### 6. Masked Display
```go
var sensitivePatterns = []string{"PASSWORD", "SECRET", "KEY", "TOKEN", "API_KEY", "PRIVATE"}

func IsSensitive(key string) bool {
    upper := strings.ToUpper(key)
    for _, p := range sensitivePatterns {
        if strings.Contains(upper, p) { return true }
    }
    return false
}

func MaskValue(key, value string) string {
    if IsSensitive(key) { return "****" }
    return value
}
```

## Todo List
- [ ] Install godotenv + huh dependencies
- [ ] Implement parser voi order + comments preservation
- [ ] Implement comparator (missing/extra detection)
- [ ] Implement validator (5 rules)
- [ ] Implement interactive sync voi huh forms
- [ ] Implement env init (full generation)
- [ ] Implement masked show
- [ ] Implement Cobra commands + register
- [ ] Implement --source va --target flags
- [ ] Test voi real .env files (comments, quotes, multiline)

## Success Criteria
- `idops env compare` hien chinh xac missing/extra vars
- `idops env sync` prompt cho missing vars va update .env
- `idops env validate` detect format issues
- `idops env init` generate .env tu .env.example
- `idops env show` mask sensitive values
- Comments va order preserved khi sync

## Risk Assessment
- **Multiline values**: godotenv handle, nhung custom parser can chua -> document limitation
- **Encoding**: UTF-8 assumed, other encodings co the fail
- **Large .env files**: Rare nhung possible -> test voi 100+ vars

## Security Considerations
- Auto-mask sensitive values trong `show` command
- Khong log .env values ra stdout ngoai `show` command
- Warn neu .env file la world-readable (permissions check)
- `idops env init` khong overwrite existing .env (require --force)

## Next Steps
- Phase 6 (Nginx) la phase cuoi, co the bat dau parallel
- Consider adding .env encryption feature trong future (YAGNI)
