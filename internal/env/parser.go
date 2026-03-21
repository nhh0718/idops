package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvFile represents a parsed .env file preserving order and comments.
type EnvFile struct {
	Path     string
	Vars     map[string]string
	Order    []string          // key insertion order
	Comments map[string]string // key -> comment line(s) above it
}

// Parse reads a .env file preserving variable order and comments.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:     path,
		Vars:     make(map[string]string),
		Comments: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	var pendingComment strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Empty line or comment
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			pendingComment.WriteString(line + "\n")
			continue
		}

		// Parse KEY=VALUE
		eqIdx := strings.Index(trimmed, "=")
		if eqIdx < 0 {
			pendingComment.WriteString(line + "\n")
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		value := strings.TrimSpace(trimmed[eqIdx+1:])
		value = unquote(value)

		env.Vars[key] = value
		env.Order = append(env.Order, key)

		if pendingComment.Len() > 0 {
			env.Comments[key] = pendingComment.String()
			pendingComment.Reset()
		}
	}

	return env, scanner.Err()
}

// Write saves the env file preserving order and comments.
func (e *EnvFile) Write(path string) error {
	var buf strings.Builder

	for _, key := range e.Order {
		if comment, ok := e.Comments[key]; ok {
			buf.WriteString(comment)
		}
		buf.WriteString(fmt.Sprintf("%s=%s\n", key, e.Vars[key]))
	}

	return os.WriteFile(path, []byte(buf.String()), 0644)
}

// unquote removes surrounding quotes from a value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
