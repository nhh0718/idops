package ssh

import "time"

// DefaultTimeout is the TCP dial timeout used for connection tests.
const DefaultTimeout = 5 * time.Second

// defaultTimeout is an internal alias for DefaultTimeout.
const defaultTimeout = DefaultTimeout
