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
	Latency time.Duration `json:"-"`
	Error   error         `json:"-"`
	// JSON-friendly fields populated by MarshalJSON-aware callers
	LatencyStr string `json:"latency"`
	ErrorStr   string `json:"error,omitempty"`
}

// FillJSON populates JSON-friendly string fields from native types.
func (r *TestResult) FillJSON() {
	r.LatencyStr = r.Latency.Round(time.Millisecond).String()
	if r.Error != nil {
		r.ErrorStr = r.Error.Error()
	}
}
