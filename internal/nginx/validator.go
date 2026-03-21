package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ValidateNginxConfig runs nginx -t to check configuration syntax.
func ValidateNginxConfig() error {
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx config invalid:\n%s", string(output))
	}
	return nil
}

// Apply enables a config by creating a symlink and reloading nginx.
// Rolls back the symlink if validation fails.
func Apply(configPath, sitesEnabled string) error {
	absConfig, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("cannot resolve path: %w", err)
	}

	linkPath := filepath.Join(sitesEnabled, filepath.Base(absConfig))

	// Create symlink
	if err := os.Symlink(absConfig, linkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// Validate
	if err := ValidateNginxConfig(); err != nil {
		os.Remove(linkPath) // rollback
		return fmt.Errorf("validation failed, rolled back: %w", err)
	}

	// Reload
	if err := exec.Command("nginx", "-s", "reload").Run(); err != nil {
		os.Remove(linkPath) // rollback
		return fmt.Errorf("reload failed, rolled back: %w", err)
	}

	return nil
}

// ListConfigs lists config files in the given directory.
func ListConfigs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", dir, err)
	}

	var configs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			configs = append(configs, entry.Name())
		}
	}
	return configs, nil
}
