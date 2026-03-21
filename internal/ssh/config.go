package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// LoadConfig parses an SSH config file line-by-line and returns all non-wildcard hosts.
func LoadConfig(path string) ([]SSHHost, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	return parseConfig(f)
}

func parseConfig(r io.Reader) ([]SSHHost, error) {
	var hosts []SSHHost
	var current *SSHHost

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := splitKeyValue(line)
		if !ok {
			continue
		}

		switch strings.ToLower(key) {
		case "host":
			if current != nil && current.Name != "*" && !strings.Contains(current.Name, "*") {
				hosts = append(hosts, *current)
			}
			if value == "*" || strings.Contains(value, "*") {
				current = &SSHHost{Name: value}
			} else {
				current = &SSHHost{Name: value, Options: make(map[string]string)}
			}
		case "hostname":
			if current != nil {
				current.Hostname = value
			}
		case "port":
			if current != nil {
				current.Port = value
			}
		case "user":
			if current != nil {
				current.User = value
			}
		case "identityfile":
			if current != nil {
				current.IdentityFile = value
			}
		case "proxyjump":
			if current != nil {
				current.ProxyJump = value
			}
		default:
			if current != nil && current.Options != nil {
				current.Options[key] = value
			}
		}
	}

	if current != nil && current.Name != "*" && !strings.Contains(current.Name, "*") {
		hosts = append(hosts, *current)
	}

	return hosts, scanner.Err()
}

// splitKeyValue splits "Key Value" or "Key=Value" into key and value parts.
func splitKeyValue(line string) (string, string, bool) {
	// Handle key=value format
	if idx := strings.IndexByte(line, '='); idx > 0 {
		return strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
	}
	// Handle "key value" format
	parts := strings.SplitN(line, " ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
	}
	return "", "", false
}

// BackupConfig copies the config file to path.bak.<timestamp>.
func BackupConfig(path string) error {
	src, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // nothing to back up
		}
		return err
	}
	defer src.Close()

	ts := time.Now().Format("20060102-150405")
	dstPath := fmt.Sprintf("%s.bak.%s", path, ts)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// AddHost appends a Host block to the SSH config file.
func AddHost(path string, host SSHHost) error {
	if err := BackupConfig(path); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\n%s", hostBlock(host))
	return err
}

// UpdateHost replaces an existing host block identified by oldName.
func UpdateHost(path, oldName string, host SSHHost) error {
	if err := BackupConfig(path); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	updated, found := replaceHostBlock(string(content), oldName, hostBlock(host))
	if !found {
		return fmt.Errorf("host %q not found", oldName)
	}

	return os.WriteFile(path, []byte(updated), 0600)
}

// DeleteHost removes a host block by name.
func DeleteHost(path, name string) error {
	if err := BackupConfig(path); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	updated, found := replaceHostBlock(string(content), name, "")
	if !found {
		return fmt.Errorf("host %q not found", name)
	}

	return os.WriteFile(path, []byte(updated), 0600)
}

// hostBlock formats a Host struct as SSH config text.
func hostBlock(h SSHHost) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Host %s\n", h.Name)
	if h.Hostname != "" {
		fmt.Fprintf(&sb, "    HostName %s\n", h.Hostname)
	}
	if h.Port != "" {
		fmt.Fprintf(&sb, "    Port %s\n", h.Port)
	}
	if h.User != "" {
		fmt.Fprintf(&sb, "    User %s\n", h.User)
	}
	if h.IdentityFile != "" {
		fmt.Fprintf(&sb, "    IdentityFile %s\n", h.IdentityFile)
	}
	if h.ProxyJump != "" {
		fmt.Fprintf(&sb, "    ProxyJump %s\n", h.ProxyJump)
	}
	for k, v := range h.Options {
		fmt.Fprintf(&sb, "    %s %s\n", k, v)
	}
	return sb.String()
}

// replaceHostBlock finds a Host block by name and replaces it with newBlock.
// Returns updated content and whether the host was found.
func replaceHostBlock(content, name, newBlock string) (string, bool) {
	lines := strings.Split(content, "\n")
	start, end := -1, len(lines)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		key, value, ok := splitKeyValue(trimmed)
		if !ok {
			continue
		}
		if strings.ToLower(key) == "host" {
			if start >= 0 {
				end = i
				break
			}
			if strings.EqualFold(value, name) {
				start = i
			}
		} else if start >= 0 && strings.ToLower(key) == "host" {
			end = i
			break
		}
	}

	// Re-scan to find end of block after start
	if start >= 0 {
		end = len(lines)
		for i := start + 1; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			key, _, ok := splitKeyValue(trimmed)
			if ok && strings.ToLower(key) == "host" {
				end = i
				break
			}
		}
	}

	if start < 0 {
		return content, false
	}

	var result []string
	result = append(result, lines[:start]...)
	if newBlock != "" {
		result = append(result, strings.Split(strings.TrimRight(newBlock, "\n"), "\n")...)
	}
	result = append(result, lines[end:]...)
	return strings.Join(result, "\n"), true
}
