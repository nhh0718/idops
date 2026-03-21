package ports

// PortInfo holds metadata for a single listening port connection.
type PortInfo struct {
	Protocol    string `json:"protocol"`
	LocalAddr   string `json:"address"`
	LocalPort   uint32 `json:"port"`
	RemoteAddr  string `json:"remoteAddr,omitempty"`
	PID         int32  `json:"pid"`
	ProcessName string `json:"process"`
	User        string `json:"user"`
	Status      string `json:"status"`
}

// SortField defines the column to sort PortInfo slices by.
type SortField int

const (
	SortByPort     SortField = iota // sort by local port number
	SortByPID                       // sort by process ID
	SortByProcess                   // sort by process name
	SortByProtocol                  // sort by protocol (tcp/udp)
)
