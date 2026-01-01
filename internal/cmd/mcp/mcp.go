package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start MCP server for AI coding agents",
		Long: `Start an MCP (Model Context Protocol) server that communicates via stdio.
AI coding agents like Claude Code can connect to this server to manage dkit processes.`,
		RunE: runMCPServer,
	}

	return cmd
}

// JSON-RPC 2.0 structures
type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP protocol structures
type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type serverCapabilities struct {
	Tools struct{} `json:"tools"`
}

type initializeResult struct {
	ProtocolVersion string              `json:"protocolVersion"`
	Capabilities    serverCapabilities  `json:"capabilities"`
	ServerInfo      serverInfo          `json:"serverInfo"`
}

type tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type toolListResult struct {
	Tools []tool `json:"tools"`
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for scanner.Scan() {
		line := scanner.Bytes()
		
		var req jsonRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			sendError(encoder, nil, -32700, "Parse error", err.Error())
			continue
		}

		// Handle request
		response := handleRequest(&req)
		if err := encoder.Encode(response); err != nil {
			return fmt.Errorf("failed to encode response: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func handleRequest(req *jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return handleInitialize(req)
	case "tools/list":
		return handleToolsList(req)
	case "tools/call":
		return handleToolsCall(req)
	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &rpcError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func handleInitialize(req *jsonRPCRequest) *jsonRPCResponse {
	result := initializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: serverCapabilities{
			Tools: struct{}{},
		},
		ServerInfo: serverInfo{
			Name:    "dkit-mcp",
			Version: "0.1.0",
		},
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func handleToolsList(req *jsonRPCRequest) *jsonRPCResponse {
	tools := []tool{
		{
			Name:        "process_list",
			Description: "List all processes started by dkit run",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"status": map[string]interface{}{
						"type":        "string",
						"description": "Filter by status",
						"enum":        []string{"running", "completed", "failed"},
					},
					"limit": map[string]interface{}{
						"type":        "number",
						"description": "Limit number of results",
					},
				},
			},
		},
		{
			Name:        "process_show",
			Description: "Show detailed information about a specific process",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"process_id": map[string]interface{}{
						"type":        "string",
						"description": "Process ID to show",
					},
				},
				"required": []string{"process_id"},
			},
		},
		{
			Name:        "process_logs",
			Description: "View process logs (stdout and stderr)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"process_id": map[string]interface{}{
						"type":        "string",
						"description": "Process ID",
					},
					"stream": map[string]interface{}{
						"type":        "string",
						"description": "Which stream to show",
						"enum":        []string{"stdout", "stderr", "both"},
						"default":     "both",
					},
					"lines": map[string]interface{}{
						"type":        "number",
						"description": "Number of lines to show",
						"default":     100,
					},
				},
				"required": []string{"process_id"},
			},
		},
		{
			Name:        "process_kill",
			Description: "Send signal to terminate a running process",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"process_id": map[string]interface{}{
						"type":        "string",
						"description": "Process ID to kill",
					},
					"signal": map[string]interface{}{
						"type":        "string",
						"description": "Signal to send",
						"enum":        []string{"SIGTERM", "SIGKILL"},
						"default":     "SIGTERM",
					},
				},
				"required": []string{"process_id"},
			},
		},
		{
			Name:        "process_clean",
			Description: "Remove process logs and metadata",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"all": map[string]interface{}{
						"type":        "boolean",
						"description": "Clean all processes",
					},
					"completed": map[string]interface{}{
						"type":        "boolean",
						"description": "Only completed processes",
					},
					"failed": map[string]interface{}{
						"type":        "boolean",
						"description": "Only failed processes",
					},
					"before": map[string]interface{}{
						"type":        "string",
						"description": "Processes started before date (ISO 8601)",
					},
				},
			},
		},
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: toolListResult{
			Tools: tools,
		},
	}
}

func handleToolsCall(req *jsonRPCRequest) *jsonRPCResponse {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &rpcError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    err.Error(),
			},
		}
	}

	// Route to appropriate tool handler
	var result interface{}
	var err error

	switch params.Name {
	case "process_list":
		result, err = handleProcessList(params.Arguments)
	case "process_show":
		result, err = handleProcessShow(params.Arguments)
	case "process_logs":
		result, err = handleProcessLogs(params.Arguments)
	case "process_kill":
		result, err = handleProcessKill(params.Arguments)
	case "process_clean":
		result, err = handleProcessClean(params.Arguments)
	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &rpcError{
				Code:    -32602,
				Message: "Unknown tool",
				Data:    params.Name,
			},
		}
	}

	if err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &rpcError{
				Code:    -32603,
				Message: "Tool execution failed",
				Data:    err.Error(),
			},
		}
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": formatToolResult(result),
				},
			},
		},
	}
}

// Tool handlers
func handleProcessList(args map[string]interface{}) (interface{}, error) {
	index, err := loadProcessIndex()
	if err != nil {
		return nil, err
	}

	// Extract filter parameters
	status := ""
	if s, ok := args["status"].(string); ok {
		status = s
	}

	limit := 0
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	// Filter processes
	filtered := filterProcesses(index.Processes, status, limit)

	return map[string]interface{}{
		"processes": filtered,
		"total":     len(index.Processes),
		"filtered":  len(filtered),
	}, nil
}

func handleProcessShow(args map[string]interface{}) (interface{}, error) {
	processID, ok := args["process_id"].(string)
	if !ok || processID == "" {
		return nil, fmt.Errorf("process_id is required")
	}

	meta, err := loadProcessMetadata(processID)
	if err != nil {
		return nil, err
	}

	// Update status if process is marked as running but actually terminated
	if meta.Status == "running" && !isProcessRunning(meta.PID) {
		meta.Status = "failed"
		if meta.ExitCode == nil {
			code := -1
			meta.ExitCode = &code
		}
	}

	// Get log file sizes
	dkitDir, _ := getDkitDir()
	stdoutPath := filepath.Join(dkitDir, "processes", processID, "stdout.log")
	stderrPath := filepath.Join(dkitDir, "processes", processID, "stderr.log")

	stdoutSize := int64(0)
	stderrSize := int64(0)

	if info, err := os.Stat(stdoutPath); err == nil {
		stdoutSize = info.Size()
	}
	if info, err := os.Stat(stderrPath); err == nil {
		stderrSize = info.Size()
	}

	return map[string]interface{}{
		"id":          meta.ID,
		"pid":         meta.PID,
		"command":     meta.Command,
		"args":        meta.Args,
		"cwd":         meta.CWD,
		"started_at":  meta.StartedAt,
		"ended_at":    meta.EndedAt,
		"status":      meta.Status,
		"exit_code":   meta.ExitCode,
		"stdout_path": meta.StdoutPath,
		"stderr_path": meta.StderrPath,
		"log_size": map[string]int64{
			"stdout": stdoutSize,
			"stderr": stderrSize,
		},
	}, nil
}

func handleProcessLogs(args map[string]interface{}) (interface{}, error) {
	processID, ok := args["process_id"].(string)
	if !ok || processID == "" {
		return nil, fmt.Errorf("process_id is required")
	}

	stream := "both"
	if s, ok := args["stream"].(string); ok {
		stream = s
	}

	lines := 100
	if l, ok := args["lines"].(float64); ok {
		lines = int(l)
	}

	dkitDir, err := getDkitDir()
	if err != nil {
		return nil, err
	}

	stdoutPath := filepath.Join(dkitDir, "processes", processID, "stdout.log")
	stderrPath := filepath.Join(dkitDir, "processes", processID, "stderr.log")

	result := map[string]interface{}{
		"process_id": processID,
	}

	if stream == "stdout" || stream == "both" {
		stdoutLines, err := readLogFile(stdoutPath, lines)
		if err != nil {
			return nil, fmt.Errorf("failed to read stdout: %w", err)
		}
		result["stdout"] = stdoutLines
	}

	if stream == "stderr" || stream == "both" {
		stderrLines, err := readLogFile(stderrPath, lines)
		if err != nil {
			return nil, fmt.Errorf("failed to read stderr: %w", err)
		}
		result["stderr"] = stderrLines
	}

	return result, nil
}

func handleProcessKill(args map[string]interface{}) (interface{}, error) {
	processID, ok := args["process_id"].(string)
	if !ok || processID == "" {
		return nil, fmt.Errorf("process_id is required")
	}

	signal := "SIGTERM"
	if s, ok := args["signal"].(string); ok {
		signal = s
	}

	meta, err := loadProcessMetadata(processID)
	if err != nil {
		return nil, err
	}

	if meta.Status != "running" {
		return nil, fmt.Errorf("process is not running (status: %s)", meta.Status)
	}

	if !isProcessRunning(meta.PID) {
		return nil, fmt.Errorf("process is no longer running")
	}

	if err := killProcess(meta.PID, signal); err != nil {
		return nil, err
	}

	// Update process status
	index, err := loadProcessIndex()
	if err != nil {
		return nil, err
	}

	for i := range index.Processes {
		if index.Processes[i].ID == processID {
			index.Processes[i].Status = "failed"
			now := time.Now()
			index.Processes[i].EndedAt = &now
			code := -1
			index.Processes[i].ExitCode = &code
			break
		}
	}

	if err := saveProcessIndex(index); err != nil {
		return nil, fmt.Errorf("failed to update process status: %w", err)
	}

	return map[string]interface{}{
		"process_id": processID,
		"signal":     signal,
		"killed_at":  time.Now(),
	}, nil
}

func handleProcessClean(args map[string]interface{}) (interface{}, error) {
	all := false
	if a, ok := args["all"].(bool); ok {
		all = a
	}

	completed := false
	if c, ok := args["completed"].(bool); ok {
		completed = c
	}

	failed := false
	if f, ok := args["failed"].(bool); ok {
		failed = f
	}

	var beforeDate *time.Time
	if b, ok := args["before"].(string); ok {
		t, err := time.Parse(time.RFC3339, b)
		if err != nil {
			return nil, fmt.Errorf("invalid date format (use ISO 8601): %w", err)
		}
		beforeDate = &t
	}

	index, err := loadProcessIndex()
	if err != nil {
		return nil, err
	}

	toDelete := []string{}
	remaining := []ProcessMetadata{}

	for _, p := range index.Processes {
		shouldDelete := false

		if all {
			shouldDelete = true
		} else if completed && p.Status == "completed" {
			shouldDelete = true
		} else if failed && p.Status == "failed" {
			shouldDelete = true
		} else if beforeDate != nil && p.StartedAt.Before(*beforeDate) {
			shouldDelete = true
		}

		if shouldDelete {
			toDelete = append(toDelete, p.ID)
		} else {
			remaining = append(remaining, p)
		}
	}

	// Delete process data
	errors := []string{}
	for _, id := range toDelete {
		if err := deleteProcessData(id); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", id, err))
		}
	}

	// Update index
	index.Processes = remaining
	if err := saveProcessIndex(index); err != nil {
		return nil, fmt.Errorf("failed to update index: %w", err)
	}

	return map[string]interface{}{
		"cleaned": len(toDelete),
		"errors":  errors,
	}, nil
}

func formatToolResult(result interface{}) string {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", result)
	}
	return string(data)
}

func sendError(encoder *json.Encoder, id interface{}, code int, message string, data interface{}) {
	response := &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &rpcError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	encoder.Encode(response)
}
