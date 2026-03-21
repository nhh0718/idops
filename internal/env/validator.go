package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// IssueType categorizes validation problems.
type IssueType int

const (
	EmptyValue IssueType = iota
	DuplicateKey
	TrailingSpace
	InvalidFormat
	UnquotedSpaces
)

// ValidationIssue describes a problem found in an env file.
type ValidationIssue struct {
	Line    int
	Key     string
	Type    IssueType
	Message string
}

// Validate checks a .env file for common issues.
func Validate(path string) ([]ValidationIssue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer f.Close()

	var issues []ValidationIssue
	seen := make(map[string]int) // key -> first line number
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		eqIdx := strings.Index(trimmed, "=")
		if eqIdx < 0 {
			issues = append(issues, ValidationIssue{
				Line: lineNum, Type: InvalidFormat,
				Message: fmt.Sprintf("no '=' found: %s", trimmed),
			})
			continue
		}

		key := trimmed[:eqIdx]
		value := trimmed[eqIdx+1:]

		// Key has spaces
		if strings.ContainsAny(key, " \t") {
			issues = append(issues, ValidationIssue{
				Line: lineNum, Key: key, Type: InvalidFormat,
				Message: "key contains spaces",
			})
		}

		// Empty value
		if strings.TrimSpace(value) == "" {
			issues = append(issues, ValidationIssue{
				Line: lineNum, Key: key, Type: EmptyValue,
				Message: "empty value",
			})
		}

		// Trailing space in value
		if value != strings.TrimRight(value, " \t") {
			issues = append(issues, ValidationIssue{
				Line: lineNum, Key: key, Type: TrailingSpace,
				Message: "trailing whitespace in value",
			})
		}

		// Unquoted value with spaces
		trimVal := strings.TrimSpace(value)
		if strings.Contains(trimVal, " ") && !isQuoted(trimVal) {
			issues = append(issues, ValidationIssue{
				Line: lineNum, Key: key, Type: UnquotedSpaces,
				Message: "value has spaces but is not quoted",
			})
		}

		// Duplicate key
		if firstLine, exists := seen[key]; exists {
			issues = append(issues, ValidationIssue{
				Line: lineNum, Key: key, Type: DuplicateKey,
				Message: fmt.Sprintf("duplicate key (first at line %d)", firstLine),
			})
		}
		seen[key] = lineNum
	}

	return issues, scanner.Err()
}

func isQuoted(s string) bool {
	return (len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"') ||
		(len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'')
}
