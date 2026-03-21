package ports

// PortInfo holds metadata for a single listening port connection.
type PortInfo struct {
	Protocol    string
	LocalAddr   string
	LocalPort   uint32
	RemoteAddr  string
	PID         int32
	ProcessName string
	User        string
	Status      string
}

// SortField defines the column to sort PortInfo slices by.
type SortField int

const (
	SortByPort     SortField = iota // sort by local port number
	SortByPID                       // sort by process ID
	SortByProcess                   // sort by process name
	SortByProtocol                  // sort by protocol (tcp/udp)
)
