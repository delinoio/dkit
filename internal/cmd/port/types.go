package port

import "time"

// PortInfo contains information about a port and the process using it
type PortInfo struct {
	Port     int       `json:"port"`
	Protocol string    `json:"protocol"`
	PID      int       `json:"pid"`
	Process  string    `json:"process"`
	Command  string    `json:"command"`
	User     string    `json:"user"`
	Started  time.Time `json:"started,omitempty"`
}

// PortCheckResult contains the result of a port availability check
type PortCheckResult struct {
	Port      int       `json:"port"`
	Available bool      `json:"available"`
	Process   *PortInfo `json:"process,omitempty"`
}

// PortListResult contains a list of ports
type PortListResult struct {
	Ports []PortInfo `json:"ports"`
	Total int        `json:"total"`
}
