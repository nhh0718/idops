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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return KeygenResult{}, fmt.Errorf("ssh-keygen thất bại: %w", err)
	}

	return KeygenResult{
		PrivateKey: keyPath,
		PublicKey:  keyPath + ".pub",
	}, nil
}
