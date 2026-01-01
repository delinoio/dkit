package port

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/delinoio/dkit/internal/utils"
)

func runListCommand(portRange string, listening, all, tcp, udp bool, format, sortBy string) error {
	// Get all ports
	ports, err := ListPorts()
	if err != nil {
		utils.PrintError("Failed to list ports: %v", err)
		os.Exit(127)
	}

	// Filter by port range if specified
	if portRange != "" {
		start, end, err := utils.ValidatePortRange(portRange)
		if err != nil {
			utils.PrintError("%v", err)
			os.Exit(2)
		}

		var filtered []PortInfo
		for _, p := range ports {
			if p.Port >= start && p.Port <= end {
				filtered = append(filtered, p)
			}
		}
		ports = filtered
	}

	// Filter by protocol if specified
	if tcp && !udp {
		var filtered []PortInfo
		for _, p := range ports {
			if strings.ToLower(p.Protocol) == "tcp" {
				filtered = append(filtered, p)
			}
		}
		ports = filtered
	} else if udp && !tcp {
		var filtered []PortInfo
		for _, p := range ports {
			if strings.ToLower(p.Protocol) == "udp" {
				filtered = append(filtered, p)
			}
		}
		ports = filtered
	}

	// Sort ports
	switch sortBy {
	case "port":
		sort.Slice(ports, func(i, j int) bool {
			return ports[i].Port < ports[j].Port
		})
	case "pid":
		sort.Slice(ports, func(i, j int) bool {
			return ports[i].PID < ports[j].PID
		})
	case "process":
		sort.Slice(ports, func(i, j int) bool {
			return ports[i].Process < ports[j].Process
		})
	}

	// Output based on format
	switch format {
	case "json":
		return outputJSON(ports)
	case "csv":
		return outputCSV(ports)
	default:
		return outputTable(ports)
	}
}

func outputJSON(ports []PortInfo) error {
	result := PortListResult{
		Ports: ports,
		Total: len(ports),
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		utils.PrintError("Failed to marshal JSON: %v", err)
		os.Exit(127)
	}

	fmt.Println(string(data))
	return nil
}

func outputCSV(ports []PortInfo) error {
	fmt.Println("port,protocol,pid,process,command,user")
	for _, p := range ports {
		fmt.Printf("%d,%s,%d,%s,%s,%s\n",
			p.Port,
			p.Protocol,
			p.PID,
			p.Process,
			p.Command,
			p.User,
		)
	}
	return nil
}

func outputTable(ports []PortInfo) error {
	if len(ports) == 0 {
		utils.PrintInfo("No ports in use")
		return nil
	}

	// Calculate column widths
	maxPort := 4
	maxProtocol := 8
	maxPID := 3
	maxProcess := 7
	maxCommand := 7
	maxUser := 4

	for _, p := range ports {
		if len(fmt.Sprintf("%d", p.Port)) > maxPort {
			maxPort = len(fmt.Sprintf("%d", p.Port))
		}
		if len(p.Protocol) > maxProtocol {
			maxProtocol = len(p.Protocol)
		}
		if len(fmt.Sprintf("%d", p.PID)) > maxPID {
			maxPID = len(fmt.Sprintf("%d", p.PID))
		}
		if len(p.Process) > maxProcess {
			maxProcess = len(p.Process)
		}
		if len(p.Command) > maxCommand {
			maxCommand = len(p.Command)
		}
		if len(p.User) > maxUser {
			maxUser = len(p.User)
		}
	}

	// Limit command width to 50 characters
	if maxCommand > 50 {
		maxCommand = 50
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s  %-*s\n",
		maxPort, "PORT",
		maxProtocol, "PROTOCOL",
		maxPID, "PID",
		maxProcess, "PROCESS",
		maxCommand, "COMMAND",
		maxUser, "USER",
	)

	// Print rows
	for _, p := range ports {
		command := p.Command
		if len(command) > 50 {
			command = command[:47] + "..."
		}

		fmt.Printf("%-*d  %-*s  %-*d  %-*s  %-*s  %-*s\n",
			maxPort, p.Port,
			maxProtocol, p.Protocol,
			maxPID, p.PID,
			maxProcess, p.Process,
			maxCommand, command,
			maxUser, p.User,
		)
	}

	return nil
}
