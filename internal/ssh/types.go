package ssh

import "time"

// SSHHost represents a single Host block from ~/.ssh/config.
type SSHHost struct {
	Name         string            `json:"name"`
	Hostname     string            `json:"hostname"`
	Port         string            `json:"port"`
	User         string            `json:"user"`
	IdentityFile string            `json:"identityFile"`
	ProxyJump    string            `json:"proxyJump"`
	Options      map[string]string `json:"options,omitempty"`
}

// TestResult holds the outcome of a connectivity test for one host.
type TestResult struct {
	Host    SSHHost       `json:"host"`
	Success bool          `json:"success"`
	Latency time.Duration `json:"latency"`
	Error   error         `json:"error,omitempty"`
}
