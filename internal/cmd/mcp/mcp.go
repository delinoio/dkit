package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

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

// Tool handlers (to be implemented)
func handleProcessList(args map[string]interface{}) (interface{}, error) {
	// TODO: Implement process list
	return "Process list - TODO", nil
}

func handleProcessShow(args map[string]interface{}) (interface{}, error) {
	// TODO: Implement process show
	return "Process show - TODO", nil
}

func handleProcessLogs(args map[string]interface{}) (interface{}, error) {
	// TODO: Implement process logs
	return "Process logs - TODO", nil
}

func handleProcessKill(args map[string]interface{}) (interface{}, error) {
	// TODO: Implement process kill
	return "Process kill - TODO", nil
}

func handleProcessClean(args map[string]interface{}) (interface{}, error) {
	// TODO: Implement process clean
	return "Process clean - TODO", nil
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
