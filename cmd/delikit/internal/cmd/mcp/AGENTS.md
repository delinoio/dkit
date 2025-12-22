# mcp Command

## Purpose
MCP (Model Context Protocol) CLI tool for managing and interacting with MCP servers and clients.

## Command Signature
```bash
dkit mcp [subcommand] [options]
```

## Subcommands Overview
This is the main MCP tool with many subcommands to be added. Core functionality includes:

### Server Management
- `dkit mcp server list` - List all configured MCP servers
- `dkit mcp server add <name> <config>` - Add a new MCP server
- `dkit mcp server remove <name>` - Remove an MCP server
- `dkit mcp server start <name>` - Start an MCP server
- `dkit mcp server stop <name>` - Stop a running MCP server
- `dkit mcp server status <name>` - Check server status

### Client Operations
- `dkit mcp connect <server>` - Connect to an MCP server
- `dkit mcp disconnect` - Disconnect from current server
- `dkit mcp call <method> [args]` - Call an MCP method
- `dkit mcp inspect <server>` - Inspect server capabilities

### Tool & Resource Management
- `dkit mcp tools list` - List available tools from connected server
- `dkit mcp tools call <tool-name> [args]` - Execute a specific tool
- `dkit mcp resources list` - List available resources
- `dkit mcp resources get <resource-id>` - Fetch a resource

### Configuration
- `dkit mcp config init` - Initialize MCP configuration
- `dkit mcp config show` - Display current configuration
- `dkit mcp config edit` - Edit configuration file
- `dkit mcp config validate` - Validate configuration syntax

### Development & Debugging
- `dkit mcp dev scaffold <name>` - Generate MCP server template
- `dkit mcp dev test <server>` - Test MCP server locally
- `dkit mcp debug <server>` - Debug MCP server with verbose logging
- `dkit mcp logs <server>` - View server logs

### Process Management (dkit run integration)
- `dkit mcp process list` - List all processes started by `dkit run`
- `dkit mcp process show <id>` - Show process details and metadata
- `dkit mcp process logs <id>` - View process logs (stdout and stderr)
- `dkit mcp process tail <id> [--follow]` - Tail process logs in real-time
- `dkit mcp process kill <id>` - Kill a running process
- `dkit mcp process clean [--all|--completed]` - Clean up old process logs

## Design Principles
- **Extensible**: Easy to add new subcommands
- **Composable**: Commands should work well with pipes and scripts
- **Self-documenting**: Built-in help for all commands
- **Safe**: Confirmation for destructive operations

## Implementation Requirements
- Must support JSON input/output for automation
- Should cache server metadata for performance
- Must handle connection timeouts gracefully
- Should support multiple concurrent server connections
- Must validate all server configurations before use

## Output Format
- Default: Human-readable table format
- `--json`: Machine-readable JSON output
- `--quiet`: Minimal output (IDs/names only)
- `--verbose`: Detailed debugging information

## Error Handling
- Connection failures: Retry with exponential backoff
- Invalid configuration: Show validation errors with line numbers
- Server not found: List available servers
- Method not supported: Show available methods

## Configuration File
Default location: `~/.config/dkit/mcp.yaml`

```yaml
servers:
  example:
    type: stdio
    command: node
    args: ["/path/to/server.js"]
    env:
      API_KEY: "..."
```

## Process Management Details

### Overview
The MCP tool provides comprehensive process management for all processes started via `dkit run`. This allows AI agents and users to monitor, inspect, and control long-running processes.

### Data Source
- All process data stored in `<project-root>/.dkit/`
- Process registry: `.dkit/index.json`
- Individual process data: `.dkit/processes/<process-id>/`

### Process List Command
```bash
dkit mcp process list [--status running|completed|failed] [--limit N]
```
**Output**: Table showing process ID, command, status, start time, exit code

### Process Show Command
```bash
dkit mcp process show <process-id>
```
**Output**: Full metadata including:
- Process ID and PID
- Full command with arguments
- Working directory
- Start/end timestamps
- Current status
- Exit code (if terminated)
- Log file paths

### Process Logs Command
```bash
dkit mcp process logs <process-id> [--stdout|--stderr|--both] [--lines N]
```
**Behavior**:
- Default: Show last 100 lines from both stdout and stderr
- `--stdout`: Only standard output
- `--stderr`: Only standard error
- `--both`: Interleaved output (default)
- `--lines N`: Show last N lines

### Process Tail Command
```bash
dkit mcp process tail <process-id> [--follow] [--stdout|--stderr]
```
**Behavior**:
- Show real-time output from running or completed processes
- `--follow` (`-f`): Continue watching for new output
- Works even if process is not running (shows historical logs)
- Exit with Ctrl+C

### Process Kill Command
```bash
dkit mcp process kill <process-id> [--signal SIGTERM|SIGKILL]
```
**Behavior**:
- Sends signal to running process
- Default: SIGTERM (graceful shutdown)
- Updates process status in registry
- Error if process already terminated

### Process Clean Command
```bash
dkit mcp process clean [--all|--completed|--failed|--before DATE]
```
**Behavior**:
- Remove process logs and metadata
- `--all`: Clean all processes (requires confirmation)
- `--completed`: Only completed processes (exit code 0)
- `--failed`: Only failed processes (exit code != 0)
- `--before DATE`: Processes started before specific date

## Future Expansion
This command will grow significantly with additional subcommands for:
- Protocol inspection and debugging
- Performance monitoring and metrics
- Server discovery and registry
- Multi-server orchestration
- Plugin system for custom tools
- Process analytics and statistics
- Alerting on process failures
- Process restart and scheduling
