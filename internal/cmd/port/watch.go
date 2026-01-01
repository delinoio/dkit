package port

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/delinoio/dkit/internal/utils"
)

func runWatchCommand(portRange string, interval int) error {
	// Parse port range if specified
	var startPort, endPort int
	if portRange != "" {
		var err error
		startPort, endPort, err = utils.ValidatePortRange(portRange)
		if err != nil {
			utils.PrintError("%v", err)
			os.Exit(2)
		}
	}

	utils.PrintInfo("Watching for port activity (Ctrl+C to stop)")
	fmt.Println()

	// Track current state
	previousPorts := make(map[int]*PortInfo)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Initial scan
	currentPorts, err := scanPorts(startPort, endPort)
	if err != nil {
		utils.PrintError("Failed to scan ports: %v", err)
		os.Exit(127)
	}

	for _, p := range currentPorts {
		previousPorts[p.Port] = &p
	}

	for {
		select {
		case <-sigChan:
			fmt.Println()
			utils.PrintInfo("Stopped watching")
			return nil

		case <-ticker.C:
			currentPorts, err := scanPorts(startPort, endPort)
			if err != nil {
				utils.PrintWarning("Failed to scan ports: %v", err)
				continue
			}

			currentPortMap := make(map[int]*PortInfo)
			for _, p := range currentPorts {
				currentPortMap[p.Port] = &p
			}

			// Check for new ports
			for port, info := range currentPortMap {
				if _, exists := previousPorts[port]; !exists {
					timestamp := time.Now().Format("2006-01-02 15:04:05")
					fmt.Printf("[%s] Port %d opened by %s (PID: %d)\n",
						timestamp, port, info.Process, info.PID)
				}
			}

			// Check for closed ports
			for port, info := range previousPorts {
				if _, exists := currentPortMap[port]; !exists {
					timestamp := time.Now().Format("2006-01-02 15:04:05")
					fmt.Printf("[%s] Port %d closed (PID: %d exited)\n",
						timestamp, port, info.PID)
				}
			}

			previousPorts = currentPortMap
		}
	}
}

func scanPorts(startPort, endPort int) ([]PortInfo, error) {
	allPorts, err := ListPorts()
	if err != nil {
		return nil, err
	}

	// Filter by range if specified
	if startPort > 0 && endPort > 0 {
		var filtered []PortInfo
		for _, p := range allPorts {
			if p.Port >= startPort && p.Port <= endPort {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}

	return allPorts, nil
}
