package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ProcessStatus represents the status of a process
type ProcessStatus string

const (
	StatusRunning   ProcessStatus = "running"
	StatusCompleted ProcessStatus = "completed"
	StatusFailed    ProcessStatus = "failed"
)

// ProcessMetadata contains metadata about a running or completed process
type ProcessMetadata struct {
	ID         string        `json:"id"`
	PID        int           `json:"pid"`
	Command    string        `json:"command"`
	Args       []string      `json:"args"`
	Cwd        string        `json:"cwd"`
	StartedAt  time.Time     `json:"started_at"`
	EndedAt    *time.Time    `json:"ended_at,omitempty"`
	Status     ProcessStatus `json:"status"`
	ExitCode   *int          `json:"exit_code,omitempty"`
	StdoutPath string        `json:"stdout_path"`
	StderrPath string        `json:"stderr_path"`
}

// ProcessRegistry manages the index of all processes
type ProcessRegistry struct {
	Processes []ProcessMetadata `json:"processes"`
}

// SaveProcessMetadata saves process metadata to the .dkit directory
func SaveProcessMetadata(projectRoot string, meta ProcessMetadata) error {
	dataDir, err := EnsureDkitDataDir(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to create .dkit directory: %w", err)
	}

	// Create processes directory
	processesDir := filepath.Join(dataDir, "processes", meta.ID)
	if err := os.MkdirAll(processesDir, 0755); err != nil {
		return fmt.Errorf("failed to create process directory: %w", err)
	}

	// Save metadata file
	metaPath := filepath.Join(processesDir, "meta.json")
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	// Update index
	if err := updateProcessIndex(dataDir, meta); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return nil
}

// LoadProcessMetadata loads process metadata from the .dkit directory
func LoadProcessMetadata(projectRoot, processID string) (*ProcessMetadata, error) {
	dataDir, err := GetDkitDataDir(projectRoot)
	if err != nil {
		return nil, err
	}

	metaPath := filepath.Join(dataDir, "processes", processID, "meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var meta ProcessMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &meta, nil
}

// updateProcessIndex updates the index.json file with new process metadata
func updateProcessIndex(dataDir string, meta ProcessMetadata) error {
	indexPath := filepath.Join(dataDir, "index.json")

	// Load existing index
	var registry ProcessRegistry
	if data, err := os.ReadFile(indexPath); err == nil {
		json.Unmarshal(data, &registry)
	}

	// Update or add process
	found := false
	for i, p := range registry.Processes {
		if p.ID == meta.ID {
			registry.Processes[i] = meta
			found = true
			break
		}
	}
	if !found {
		registry.Processes = append(registry.Processes, meta)
	}

	// Save index
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0644)
}

// ListProcesses returns all processes from the index
func ListProcesses(projectRoot string) ([]ProcessMetadata, error) {
	dataDir, err := GetDkitDataDir(projectRoot)
	if err != nil {
		return nil, err
	}

	indexPath := filepath.Join(dataDir, "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []ProcessMetadata{}, nil
		}
		return nil, err
	}

	var registry ProcessRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}

	return registry.Processes, nil
}

// GenerateProcessID generates a unique process ID
func GenerateProcessID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
