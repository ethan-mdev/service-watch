package core

// Service represents a system service.
type Service struct {
	Name          string  `json:"name"`
	DisplayName   string  `json:"displayName,omitempty"`
	State         string  `json:"state,omitempty"`     // running|stopped|starting|stopping|unknown
	StartType     string  `json:"startType,omitempty"` // auto|manual|disabled|unknown
	CanStop       bool    `json:"canStop,omitempty"`
	PID           int     `json:"pid,omitempty"`
	CPUPercent    float64 `json:"cpuPercent,omitempty"`    // CPU usage percentage
	MemoryMB      float64 `json:"memoryMB,omitempty"`      // Memory usage in MB
	UptimeSeconds int64   `json:"uptimeSeconds,omitempty"` // How long service has been running
}

// WatchlistItem represents an item in the watchlist.
type WatchlistItem struct {
	ServiceName  string   `json:"serviceName"`           // Name of the service being watched
	AutoRestart  bool     `json:"autoRestart"`           // Should we auto-restart if it crashes?
	RestartCount int      `json:"restartCount"`          // How many times have we restarted it?
	FailCount    int      `json:"failCount,omitempty"`   // Consecutive failure count
	LastRestart  string   `json:"lastRestart,omitempty"` // ISO timestamp of last restart
	Service      *Service `json:"service,omitempty"`     // Current service state when fetched
}

// Represents an SSE event.
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
