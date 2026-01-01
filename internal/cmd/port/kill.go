package port

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/delinoio/dkit/internal/utils"
)

func runKillCommand(portStr string, force bool, signalName string, timeout int) error {
	// Validate port number
	port, err := utils.ValidatePort(portStr)
	if err != nil {
		utils.PrintError("%v", err)
		os.Exit(2)
	}

	// Get port info
	portInfo, err := GetPortInfo(port)
	if err != nil {
		utils.PrintError("Failed to check port: %v", err)
		os.Exit(127)
	}

	if portInfo == nil {
		utils.PrintInfo("Port %d is not in use", port)
		os.Exit(1)
	}

	// Show process info
	utils.PrintInfo("Port %d is in use by:", port)
	fmt.Printf("[dkit] Process: %s (PID: %d)\n", portInfo.Process, portInfo.PID)
	fmt.Printf("[dkit] Command: %s\n", portInfo.Command)
	if portInfo.User != "" {
		fmt.Printf("[dkit] User: %s\n", portInfo.User)
	}
	fmt.Println()

	// Check if system process
	isSystemProcess := portInfo.User == "root" || portInfo.PID < 1000
	if isSystemProcess {
		utils.PrintWarning("This is a system process owned by %s", portInfo.User)
	}

	// Ask for confirmation unless force flag is set
	if !force {
		prompt := "Terminate this process?"
		if isSystemProcess {
			prompt = "Are you sure you want to terminate this system process?"
		}

		if !utils.Confirm(prompt) {
			utils.PrintInfo("Operation cancelled")
			os.Exit(2)
		}
	}

	// Terminate process
	if force {
		utils.PrintInfo("Terminating process %d on port %d...", portInfo.PID, port)
	} else {
		utils.PrintInfo("Sending SIG%s to PID %d...", signalName, portInfo.PID)
	}

	err = terminateProcess(portInfo.PID, signalName, timeout)
	if err != nil {
		utils.PrintError("Failed to terminate process: %v", err)
		os.Exit(4)
	}

	utils.PrintSuccess("Process terminated successfully")

	// Verify port is now available
	time.Sleep(100 * time.Millisecond)
	portInfo, err = GetPortInfo(port)
	if err == nil && portInfo == nil {
		utils.PrintSuccess("Port %d is now available", port)
	}

	os.Exit(0)
	return nil
}

func terminateProcess(pid int, signalName string, timeout int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process not found: %w", err)
	}

	// Parse signal
	var sig os.Signal
	switch signalName {
	case "TERM":
		sig = syscall.SIGTERM
	case "KILL":
		sig = syscall.SIGKILL
	case "HUP":
		sig = syscall.SIGHUP
	default:
		return fmt.Errorf("unsupported signal: %s", signalName)
	}

	// Send signal
	if runtime.GOOS == "windows" {
		// On Windows, we can only kill
		return process.Kill()
	}

	err = process.Signal(sig)
	if err != nil {
		return fmt.Errorf("failed to send signal: %w", err)
	}

	// If we sent SIGTERM, wait for graceful shutdown
	if sig == syscall.SIGTERM && timeout > 0 {
		utils.PrintInfo("Waiting for process to exit (timeout: %ds)...", timeout)

		// Wait for timeout
		deadline := time.Now().Add(time.Duration(timeout) * time.Second)
		for time.Now().Before(deadline) {
			if !IsProcessRunning(pid) {
				return nil
			}
			time.Sleep(100 * time.Millisecond)
		}

		// If still running, send SIGKILL
		utils.PrintWarning("Process did not exit gracefully")
		utils.PrintInfo("Sending SIGKILL to PID %d...", pid)
		return process.Signal(syscall.SIGKILL)
	}

	return nil
}
