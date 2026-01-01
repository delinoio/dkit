package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"
)

// Process metadata structure (from run/AGENTS.md)
type ProcessMetadata struct {
	ID         string     `json:"id"`
	PID        int        `json:"pid"`
	Command    string     `json:"command"`
	Args       []string   `json:"args"`
	CWD        string     `json:"cwd"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
	Status     string     `json:"status"` // running, completed, failed
	ExitCode   *int       `json:"exit_code,omitempty"`
	StdoutPath string     `json:"stdout_path"`
	StderrPath string     `json:"stderr_path"`
}

// ProcessIndex represents the process registry
type ProcessIndex struct {
	Processes []ProcessMetadata `json:"processes"`
}

// getDkitDir finds the .dkit directory in the project root
func getDkitDir() (string, error) {
	// Start from current directory and walk up to find .dkit
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	dir := cwd
	for {
		dkitPath := filepath.Join(dir, ".dkit")
		if info, err := os.Stat(dkitPath); err == nil && info.IsDir() {
			return dkitPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf(".dkit directory not found in project")
}

// loadProcessIndex loads the process registry
func loadProcessIndex() (*ProcessIndex, error) {
	dkitDir, err := getDkitDir()
	if err != nil {
		return nil, err
	}

	indexPath := filepath.Join(dkitDir, "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ProcessIndex{Processes: []ProcessMetadata{}}, nil
		}
		return nil, fmt.Errorf("failed to read index.json: %w", err)
	}

	var index ProcessIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index.json: %w", err)
	}

	return &index, nil
}

// saveProcessIndex saves the process registry
func saveProcessIndex(index *ProcessIndex) error {
	dkitDir, err := getDkitDir()
	if err != nil {
		return err
	}

	indexPath := filepath.Join(dkitDir, "index.json")
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write index.json: %w", err)
	}

	return nil
}

// loadProcessMetadata loads metadata for a specific process
func loadProcessMetadata(processID string) (*ProcessMetadata, error) {
	dkitDir, err := getDkitDir()
	if err != nil {
		return nil, err
	}

	metaPath := filepath.Join(dkitDir, "processes", processID, "meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("process not found: %w", err)
	}

	var meta ProcessMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &meta, nil
}

// readLogFile reads the last N lines from a log file
func readLogFile(filePath string, lines int) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	// Split by newlines
	allLines := []string{}
	currentLine := ""
	for _, b := range data {
		if b == '\n' {
			allLines = append(allLines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(b)
		}
	}
	if currentLine != "" {
		allLines = append(allLines, currentLine)
	}

	// Return last N lines
	if lines <= 0 || lines >= len(allLines) {
		return allLines, nil
	}
	return allLines[len(allLines)-lines:], nil
}

// isProcessRunning checks if a process is still running
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// filterProcesses filters processes based on criteria
func filterProcesses(processes []ProcessMetadata, status string, limit int) []ProcessMetadata {
	filtered := []ProcessMetadata{}

	for _, p := range processes {
		// Update status if process is marked as running but actually terminated
		if p.Status == "running" && !isProcessRunning(p.PID) {
			p.Status = "failed"
			if p.ExitCode == nil {
				code := -1
				p.ExitCode = &code
			}
		}

		// Apply status filter
		if status != "" && p.Status != status {
			continue
		}

		filtered = append(filtered, p)
	}

	// Sort by start time (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].StartedAt.After(filtered[j].StartedAt)
	})

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered
}

// killProcess sends a signal to a process
func killProcess(pid int, signal string) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process not found: %w", err)
	}

	var sig syscall.Signal
	switch signal {
	case "SIGTERM":
		sig = syscall.SIGTERM
	case "SIGKILL":
		sig = syscall.SIGKILL
	default:
		return fmt.Errorf("unknown signal: %s", signal)
	}

	if err := process.Signal(sig); err != nil {
		return fmt.Errorf("failed to send signal: %w", err)
	}

	return nil
}

// deleteProcessData removes process directory
func deleteProcessData(processID string) error {
	dkitDir, err := getDkitDir()
	if err != nil {
		return err
	}

	processDir := filepath.Join(dkitDir, "processes", processID)
	if err := os.RemoveAll(processDir); err != nil {
		return fmt.Errorf("failed to remove process directory: %w", err)
	}

	return nil
}
