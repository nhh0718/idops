package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// KeygenOptions holds parameters for SSH key generation.
type KeygenOptions struct {
	Name    string // key file name, e.g. "id_ed25519"
	Type    string // "ed25519" or "rsa"
	Bits    int    // RSA key bits (default 4096), ignored for ed25519
	Comment string // e.g. email address
	Force   bool   // overwrite existing key without asking
}

// KeygenResult contains paths of generated keys.
type KeygenResult struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Output     string `json:"output,omitempty"` // ssh-keygen output (fingerprint, etc)
}

// GenerateKey runs ssh-keygen to create a new SSH key pair.
// It returns the paths to the generated private and public keys.
func GenerateKey(opts KeygenOptions) (KeygenResult, error) {
	if opts.Name == "" {
		opts.Name = "id_ed25519"
	}
	if opts.Type == "" {
		opts.Type = "ed25519"
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return KeygenResult{}, fmt.Errorf("không thể lấy home dir: %w", err)
	}

	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return KeygenResult{}, fmt.Errorf("không thể tạo ~/.ssh: %w", err)
	}

	keyPath := filepath.Join(sshDir, opts.Name)

	// Check if key already exists
	if _, err := os.Stat(keyPath); err == nil {
		if !opts.Force {
			return KeygenResult{}, fmt.Errorf("key '%s' đã tồn tại. Dùng --force để ghi đè", keyPath)
		}
		// Remove existing files so ssh-keygen won't prompt
		os.Remove(keyPath)
		os.Remove(keyPath + ".pub")
	}

	args := []string{"-t", opts.Type}
	if opts.Type == "rsa" && opts.Bits > 0 {
		args = append(args, "-b", fmt.Sprintf("%d", opts.Bits))
	}
	if opts.Comment != "" {
		args = append(args, "-C", opts.Comment)
	}
	args = append(args, "-f", keyPath, "-N", "")

	sshKeygen, err := exec.LookPath("ssh-keygen")
	if err != nil {
		return KeygenResult{}, fmt.Errorf("ssh-keygen không tìm thấy trong PATH: %w", err)
	}

	cmd := exec.Command(sshKeygen, args...)
	// Capture output instead of piping to stdout (prevents mixing with JSON output)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return KeygenResult{}, fmt.Errorf("ssh-keygen thất bại: %s", string(output))
	}

	return KeygenResult{
		PrivateKey: keyPath,
		PublicKey:  keyPath + ".pub",
		Output:     string(output),
	}, nil
}

// SSHKeyInfo describes a key pair found in ~/.ssh/.
type SSHKeyInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`       // ed25519, rsa, ecdsa, dsa, unknown
	PublicKey  string `json:"publicKey"`   // path to .pub file
	PrivateKey string `json:"privateKey"`  // path to private key
	Comment    string `json:"comment"`     // comment from .pub file
	Fingerprint string `json:"fingerprint"` // SHA256 fingerprint
}

// ListKeys scans ~/.ssh/ for key pairs (.pub files) and returns info for each.
func ListKeys() ([]SSHKeyInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get home dir: %w", err)
	}

	sshDir := filepath.Join(home, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read ~/.ssh: %w", err)
	}

	var keys []SSHKeyInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".pub") {
			continue
		}

		pubPath := filepath.Join(sshDir, e.Name())
		name := strings.TrimSuffix(e.Name(), ".pub")
		privPath := filepath.Join(sshDir, name)

		// Check private key exists
		if _, err := os.Stat(privPath); err != nil {
			continue // orphan .pub without private key
		}

		info := SSHKeyInfo{
			Name:       name,
			PublicKey:  pubPath,
			PrivateKey: privPath,
		}

		// Read .pub to get type and comment
		pubData, err := os.ReadFile(pubPath)
		if err == nil {
			parts := strings.Fields(strings.TrimSpace(string(pubData)))
			if len(parts) >= 1 {
				info.Type = keyTypeFromPrefix(parts[0])
			}
			if len(parts) >= 3 {
				info.Comment = parts[2]
			}
		}

		// Get fingerprint via ssh-keygen -lf
		if fp := getFingerprint(pubPath); fp != "" {
			info.Fingerprint = fp
		}

		keys = append(keys, info)
	}

	return keys, nil
}

// DeleteKey removes a key pair from ~/.ssh/.
func DeleteKey(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get home dir: %w", err)
	}

	sshDir := filepath.Join(home, ".ssh")
	privPath := filepath.Join(sshDir, name)
	pubPath := privPath + ".pub"

	if _, err := os.Stat(privPath); err != nil {
		return fmt.Errorf("key '%s' not found", name)
	}

	os.Remove(privPath)
	os.Remove(pubPath)
	return nil
}

// keyTypeFromPrefix maps ssh public key prefix to a readable type.
func keyTypeFromPrefix(prefix string) string {
	switch {
	case strings.Contains(prefix, "ed25519"):
		return "ed25519"
	case strings.Contains(prefix, "ecdsa"):
		return "ecdsa"
	case strings.Contains(prefix, "rsa"):
		return "rsa"
	case strings.Contains(prefix, "dsa"):
		return "dsa"
	default:
		return "unknown"
	}
}

// getFingerprint returns the SHA256 fingerprint of a public key file.
func getFingerprint(pubPath string) string {
	sshKeygen, err := exec.LookPath("ssh-keygen")
	if err != nil {
		return ""
	}
	out, err := exec.Command(sshKeygen, "-lf", pubPath).Output()
	if err != nil {
		return ""
	}
	// Output: "256 SHA256:xxx comment (ED25519)"
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
