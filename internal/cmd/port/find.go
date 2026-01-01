package port

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/delinoio/dkit/internal/utils"
)

func runFindCommand(pattern, format string) error {
	// Get all ports
	allPorts, err := ListPorts()
	if err != nil {
		utils.PrintError("Failed to list ports: %v", err)
		os.Exit(127)
	}

	// Parse pattern and filter ports
	var filtered []PortInfo

	// Check if pattern is a range (e.g., "3000-4000")
	if strings.Contains(pattern, "-") {
		start, end, err := utils.ValidatePortRange(pattern)
		if err != nil {
			utils.PrintError("%v", err)
			os.Exit(2)
		}

		for _, p := range allPorts {
			if p.Port >= start && p.Port <= end {
				filtered = append(filtered, p)
			}
		}
	} else if strings.Contains(pattern, ",") {
		// Multiple specific ports (e.g., "3000,8080,9000")
		ports := strings.Split(pattern, ",")
		portSet := make(map[int]bool)
		for _, portStr := range ports {
			port, err := utils.ValidatePort(strings.TrimSpace(portStr))
			if err != nil {
				utils.PrintError("Invalid port in pattern: %s", portStr)
				os.Exit(2)
			}
			portSet[port] = true
		}

		for _, p := range allPorts {
			if portSet[p.Port] {
				filtered = append(filtered, p)
			}
		}
	} else if strings.Contains(pattern, "*") {
		// Wildcard pattern (e.g., "30*" matches 3000-3099, 30000-30999, etc.)
		prefix := strings.TrimSuffix(pattern, "*")
		for _, p := range allPorts {
			portStr := fmt.Sprintf("%d", p.Port)
			if strings.HasPrefix(portStr, prefix) {
				filtered = append(filtered, p)
			}
		}
	} else {
		// Single port
		port, err := utils.ValidatePort(pattern)
		if err != nil {
			utils.PrintError("%v", err)
			os.Exit(2)
		}

		for _, p := range allPorts {
			if p.Port == port {
				filtered = append(filtered, p)
			}
		}
	}

	// Output results
	if len(filtered) == 0 {
		utils.PrintInfo("No ports found matching pattern: %s", pattern)
		return nil
	}

	switch format {
	case "json":
		return outputJSON(filtered)
	case "list":
		for _, p := range filtered {
			fmt.Printf("%d\n", p.Port)
		}
		return nil
	default:
		return outputTable(filtered)
	}
}

// Helper function to check if port matches wildcard pattern
func matchesWildcard(port int, pattern string) bool {
	portStr := strconv.Itoa(port)
	patternParts := strings.Split(pattern, "*")

	if len(patternParts) == 1 {
		return portStr == pattern
	}

	// Check prefix
	if !strings.HasPrefix(portStr, patternParts[0]) {
		return false
	}

	// Check suffix if exists
	if len(patternParts) > 1 && patternParts[1] != "" {
		if !strings.HasSuffix(portStr, patternParts[1]) {
			return false
		}
	}

	return true
}
