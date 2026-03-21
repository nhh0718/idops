package ssh

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const maxConcurrentTests = 10

// TestConnection performs a TCP dial to the host's address and returns latency.
func TestConnection(host SSHHost, timeout time.Duration) (time.Duration, error) {
	addr := resolveAddr(host)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return 0, err
	}
	latency := time.Since(start)
	conn.Close()
	return latency, nil
}

// TestAllConnections runs TCP tests for all hosts in parallel (max 10 goroutines).
func TestAllConnections(hosts []SSHHost, timeout time.Duration) []TestResult {
	results := make([]TestResult, len(hosts))
	sem := make(chan struct{}, maxConcurrentTests)
	var wg sync.WaitGroup

	for i, h := range hosts {
		wg.Add(1)
		go func(idx int, host SSHHost) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			latency, err := TestConnection(host, timeout)
			results[idx] = TestResult{
				Host:    host,
				Success: err == nil,
				Latency: latency,
				Error:   err,
			}
		}(i, h)
	}

	wg.Wait()
	return results
}

// resolveAddr returns "hostname:port" for dialing, defaulting port to 22.
func resolveAddr(host SSHHost) string {
	hostname := host.Hostname
	if hostname == "" {
		hostname = host.Name
	}
	port := host.Port
	if port == "" {
		port = "22"
	}
	return fmt.Sprintf("%s:%s", hostname, port)
}
