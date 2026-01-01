package port

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/delinoio/dkit/internal/utils"
)

func runCheckCommand(portStr string, jsonOutput, quiet bool) error {
	// Validate port number
	port, err := utils.ValidatePort(portStr)
	if err != nil {
		if !quiet {
			utils.PrintError("%v", err)
		}
		os.Exit(2)
	}

	// Get port info
	portInfo, err := GetPortInfo(port)
	if err != nil {
		if !quiet {
			utils.PrintError("Failed to check port: %v", err)
		}
		os.Exit(127)
	}

	result := PortCheckResult{
		Port:      port,
		Available: portInfo == nil,
		Process:   portInfo,
	}

	// Handle JSON output
	if jsonOutput {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			utils.PrintError("Failed to marshal JSON: %v", err)
			os.Exit(127)
		}
		fmt.Println(string(data))
		if result.Available {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// Handle quiet mode
	if quiet {
		if result.Available {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// Normal output
	if result.Available {
		utils.PrintSuccess("Port %d is available", port)
		os.Exit(0)
	} else {
		utils.PrintInfo("Port %d is in use", port)
		fmt.Println()
		fmt.Printf("Process: %s\n", portInfo.Process)
		fmt.Printf("PID: %d\n", portInfo.PID)
		fmt.Printf("Command: %s\n", portInfo.Command)
		if portInfo.User != "" {
			fmt.Printf("User: %s\n", portInfo.User)
		}
		os.Exit(1)
	}

	return nil
}
