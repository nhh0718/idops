package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
