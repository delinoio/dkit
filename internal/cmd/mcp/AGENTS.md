# mcp Command

## Purpose
MCP (Model Context Protocol) server that provides tools for AI coding agents like Claude Code, Cursor, etc.

## Command Signature
```bash
dkit mcp
```

## Overview
When `dkit mcp` is executed, it starts an MCP server that communicates via stdio using JSON-RPC 2.0 protocol. AI coding agents can connect to this server and invoke tools for process management.

## Design Principles
- **Stdio-based**: Communication via standard input/output using JSON-RPC 2.0
- **Tool-oriented**: Exposes capabilities as MCP tools
- **Zero-config**: No setup required, just run and connect
- **Safe**: Operations require explicit parameters

## Protocol
- **Transport**: stdio (standard input/output)
- **Format**: JSON-RPC 2.0
- **Lifecycle**: Server runs until stdin closes or receives shutdown request

## Available Tools
All tools exposed to AI agents for managing `dkit run` processes:

### Data Source
- All process data stored in `<project-root>/.dkit/`
- Process registry: `.dkit/index.json`
- Individual process data: `.dkit/processes/<process-id>/`

### 1. process_list
**Description**: List all processes started by `dkit run`

**Input Schema**:
```json
{
  "status": "running|completed|failed",  // optional
  "limit": 10  // optional
}
```

**Output**: Array of process objects with ID, command, status, start time, exit code

### 2. process_show
**Description**: Show detailed information about a specific process

**Input Schema**:
```json
{
  "process_id": "unique-process-id"
}
```

**Output**: Full metadata including PID, command, working directory, timestamps, status, exit code, log paths

### 3. process_logs
**Description**: View process logs (stdout and stderr)

**Input Schema**:
```json
{
  "process_id": "unique-process-id",
  "stream": "stdout|stderr|both",  // optional, default: both
  "lines": 100  // optional, default: 100
}
```

**Output**: Log content as string

### 4. process_tail
**Description**: Show real-time output from a process

**Input Schema**:
```json
{
  "process_id": "unique-process-id",
  "follow": true,  // optional, default: false
  "stream": "stdout|stderr|both"  // optional, default: both
}
```

**Output**: Log content stream (for follow mode, this is more complex)

### 5. process_kill
**Description**: Send signal to terminate a running process

**Input Schema**:
```json
{
  "process_id": "unique-process-id",
  "signal": "SIGTERM|SIGKILL"  // optional, default: SIGTERM
}
```

**Output**: Confirmation with updated process status

### 6. process_clean
**Description**: Remove process logs and metadata

**Input Schema**:
```json
{
  "all": false,  // optional
  "completed": false,  // optional
  "failed": false,  // optional
  "before": "2025-01-01T00:00:00Z"  // optional, ISO 8601 format
}
```

**Output**: List of removed process IDs

## Usage with AI Agents

### Claude Desktop Configuration
Add to `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "dkit": {
      "command": "dkit",
      "args": ["mcp"]
    }
  }
}
```

### Cursor Configuration
Add to MCP settings to connect the server.

## Future Expansion
Additional tools will be added for other dkit capabilities beyond process management.
