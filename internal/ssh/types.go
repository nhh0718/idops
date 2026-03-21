package ssh

import "time"

// SSHHost represents a single Host block from ~/.ssh/config.
type SSHHost struct {
	Name         string
	Hostname     string
	Port         string
	User         string
	IdentityFile string
	ProxyJump    string
	Options      map[string]string
}

// TestResult holds the outcome of a connectivity test for one host.
type TestResult struct {
	Host    SSHHost
	Success bool
	Latency time.Duration
	Error   error
}
