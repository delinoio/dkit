package utils

import (
	"fmt"
	"strconv"
)

// ValidatePort validates that a port number is in the valid range (1-65535)
func ValidatePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port number: %s", portStr)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}

	return port, nil
}

// ValidatePortRange validates a port range string (e.g., "3000-4000")
func ValidatePortRange(rangeStr string) (int, int, error) {
	parts := splitPortRange(rangeStr)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid port range format: %s (expected: start-end)", rangeStr)
	}

	start, err := ValidatePort(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start port: %w", err)
	}

	end, err := ValidatePort(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end port: %w", err)
	}

	if start > end {
		return 0, 0, fmt.Errorf("start port (%d) must be less than or equal to end port (%d)", start, end)
	}

	return start, end, nil
}

func splitPortRange(rangeStr string) []string {
	var parts []string
	current := ""
	for _, ch := range rangeStr {
		if ch == '-' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
