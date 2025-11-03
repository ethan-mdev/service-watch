package core

// Service represents a system service.
type Service struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	State       string `json:"state,omitempty"`     // running|stopped|starting|stopping|unknown
	StartType   string `json:"startType,omitempty"` // auto|manual|disabled|unknown
	CanStop     bool   `json:"canStop,omitempty"`
	PID         int    `json:"pid,omitempty"`
}
