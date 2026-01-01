package port

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// GetPortInfo retrieves information about a specific port
func GetPortInfo(port int) (*PortInfo, error) {
	switch runtime.GOOS {
	case "darwin":
		return getPortInfoDarwin(port)
	case "linux":
		return getPortInfoLinux(port)
	case "windows":
		return getPortInfoWindows(port)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// ListPorts lists all ports currently in use
func ListPorts() ([]PortInfo, error) {
	switch runtime.GOOS {
	case "darwin":
		return listPortsDarwin()
	case "linux":
		return listPortsLinux()
	case "windows":
		return listPortsWindows()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// macOS implementation using lsof
func getPortInfoDarwin(port int) (*PortInfo, error) {
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-sTCP:LISTEN", "-n", "-P")
	output, err := cmd.Output()
	if err != nil {
		// Port is not in use
		return nil, nil
	}

	return parseLsofOutput(output, port)
}

func listPortsDarwin() ([]PortInfo, error) {
	cmd := exec.Command("lsof", "-i", "-sTCP:LISTEN", "-n", "-P")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsof: %w", err)
	}

	return parseLsofListOutput(output)
}

func parseLsofOutput(output []byte, port int) (*PortInfo, error) {
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, nil
	}

	// Skip header line
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		return &PortInfo{
			Port:     port,
			Protocol: "tcp",
			PID:      pid,
			Process:  fields[0],
			Command:  strings.Join(fields[8:], " "),
			User:     fields[2],
		}, nil
	}

	return nil, nil
}

func parseLsofListOutput(output []byte) ([]PortInfo, error) {
	var ports []PortInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Skip header
	if !scanner.Scan() {
		return ports, nil
	}

	seenPorts := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		// Extract port from address field (e.g., "*:3000" or "127.0.0.1:3000")
		addrField := fields[8]
		parts := strings.Split(addrField, ":")
		if len(parts) < 2 {
			continue
		}

		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		// Avoid duplicates
		key := fmt.Sprintf("%d-%d", port, pid)
		if seenPorts[key] {
			continue
		}
		seenPorts[key] = true

		ports = append(ports, PortInfo{
			Port:     port,
			Protocol: "tcp",
			PID:      pid,
			Process:  fields[0],
			User:     fields[2],
			Command:  fields[0],
		})
	}

	return ports, nil
}

// Linux implementation using ss/netstat
func getPortInfoLinux(port int) (*PortInfo, error) {
	// Try ss first (modern)
	cmd := exec.Command("ss", "-tlnp", fmt.Sprintf("sport = :%d", port))
	output, err := cmd.Output()
	if err == nil {
		return parseSsOutput(output, port)
	}

	// Fallback to netstat
	cmd = exec.Command("netstat", "-tlnp")
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run netstat: %w", err)
	}

	return parseNetstatOutput(output, port)
}

func listPortsLinux() ([]PortInfo, error) {
	// Try ss first
	cmd := exec.Command("ss", "-tlnp")
	output, err := cmd.Output()
	if err == nil {
		return parseSsListOutput(output)
	}

	// Fallback to netstat
	cmd = exec.Command("netstat", "-tlnp")
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run netstat: %w", err)
	}

	return parseNetstatListOutput(output)
}

func parseSsOutput(output []byte, port int) (*PortInfo, error) {
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, nil
	}

	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// Extract PID from process field
		processField := fields[len(fields)-1]
		pid := extractPIDFromProcess(processField)
		if pid == 0 {
			continue
		}

		return &PortInfo{
			Port:     port,
			Protocol: "tcp",
			PID:      pid,
			Process:  extractProcessName(processField),
			Command:  extractProcessName(processField),
		}, nil
	}

	return nil, nil
}

func parseSsListOutput(output []byte) ([]PortInfo, error) {
	var ports []PortInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Skip header
	if !scanner.Scan() {
		return ports, nil
	}

	seenPorts := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Extract port from local address (e.g., "*:3000" or "127.0.0.1:3000")
		localAddr := fields[3]
		parts := strings.Split(localAddr, ":")
		if len(parts) < 2 {
			continue
		}

		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		// Extract process info from last field
		processField := fields[len(fields)-1]
		pid := extractPIDFromProcess(processField)
		if pid == 0 {
			continue
		}

		key := fmt.Sprintf("%d-%d", port, pid)
		if seenPorts[key] {
			continue
		}
		seenPorts[key] = true

		ports = append(ports, PortInfo{
			Port:     port,
			Protocol: "tcp",
			PID:      pid,
			Process:  extractProcessName(processField),
			Command:  extractProcessName(processField),
		})
	}

	return ports, nil
}

func parseNetstatOutput(output []byte, port int) (*PortInfo, error) {
	lines := strings.Split(string(output), "\n")

	portStr := fmt.Sprintf(":%d", port)
	for _, line := range lines {
		if !strings.Contains(line, portStr) || !strings.Contains(line, "LISTEN") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}

		processField := fields[len(fields)-1]
		pid := extractPIDFromProcess(processField)
		if pid == 0 {
			continue
		}

		return &PortInfo{
			Port:     port,
			Protocol: "tcp",
			PID:      pid,
			Process:  extractProcessName(processField),
			Command:  extractProcessName(processField),
		}, nil
	}

	return nil, nil
}

func parseNetstatListOutput(output []byte) ([]PortInfo, error) {
	var ports []PortInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	seenPorts := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "LISTEN") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}

		// Extract port from local address
		localAddr := fields[3]
		parts := strings.Split(localAddr, ":")
		if len(parts) < 2 {
			continue
		}

		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		processField := fields[len(fields)-1]
		pid := extractPIDFromProcess(processField)
		if pid == 0 {
			continue
		}

		key := fmt.Sprintf("%d-%d", port, pid)
		if seenPorts[key] {
			continue
		}
		seenPorts[key] = true

		ports = append(ports, PortInfo{
			Port:     port,
			Protocol: "tcp",
			PID:      pid,
			Process:  extractProcessName(processField),
			Command:  extractProcessName(processField),
		})
	}

	return ports, nil
}

// Windows implementation using netstat
func getPortInfoWindows(port int) (*PortInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run netstat: %w", err)
	}

	return parseNetstatWindowsOutput(output, port)
}

func listPortsWindows() ([]PortInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run netstat: %w", err)
	}

	return parseNetstatWindowsListOutput(output)
}

func parseNetstatWindowsOutput(output []byte, port int) (*PortInfo, error) {
	lines := strings.Split(string(output), "\n")
	portStr := fmt.Sprintf(":%d", port)

	for _, line := range lines {
		if !strings.Contains(line, portStr) || !strings.Contains(line, "LISTENING") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		pid, err := strconv.Atoi(fields[4])
		if err != nil {
			continue
		}

		processName := getWindowsProcessName(pid)
		return &PortInfo{
			Port:     port,
			Protocol: strings.ToLower(fields[0]),
			PID:      pid,
			Process:  processName,
			Command:  processName,
		}, nil
	}

	return nil, nil
}

func parseNetstatWindowsListOutput(output []byte) ([]PortInfo, error) {
	var ports []PortInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	seenPorts := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "LISTENING") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Extract port from local address
		localAddr := fields[1]
		parts := strings.Split(localAddr, ":")
		if len(parts) < 2 {
			continue
		}

		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		pid, err := strconv.Atoi(fields[4])
		if err != nil {
			continue
		}

		key := fmt.Sprintf("%d-%d", port, pid)
		if seenPorts[key] {
			continue
		}
		seenPorts[key] = true

		processName := getWindowsProcessName(pid)
		ports = append(ports, PortInfo{
			Port:     port,
			Protocol: strings.ToLower(fields[0]),
			PID:      pid,
			Process:  processName,
			Command:  processName,
		})
	}

	return ports, nil
}

func getWindowsProcessName(pid int) string {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		fields := strings.Split(lines[0], ",")
		if len(fields) > 0 {
			// Remove quotes
			return strings.Trim(fields[0], "\"")
		}
	}

	return "unknown"
}

// Helper functions
func extractPIDFromProcess(processField string) int {
	// Process field format: "users:(("process",pid=1234,fd=5))"
	pidStart := strings.Index(processField, "pid=")
	if pidStart == -1 {
		return 0
	}

	pidStart += 4
	pidEnd := pidStart
	for pidEnd < len(processField) && processField[pidEnd] >= '0' && processField[pidEnd] <= '9' {
		pidEnd++
	}

	if pidEnd == pidStart {
		return 0
	}

	pid, err := strconv.Atoi(processField[pidStart:pidEnd])
	if err != nil {
		return 0
	}

	return pid
}

func extractProcessName(processField string) string {
	// Process field format: "users:(("process",pid=1234,fd=5))"
	nameStart := strings.Index(processField, "((\"")
	if nameStart == -1 {
		return "unknown"
	}

	nameStart += 3
	nameEnd := strings.Index(processField[nameStart:], "\"")
	if nameEnd == -1 {
		return "unknown"
	}

	return processField[nameStart : nameStart+nameEnd]
}

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix, FindProcess always succeeds, so we need to send signal 0
	err = process.Signal(os.Signal(nil))
	return err == nil
}
