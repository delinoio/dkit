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

## Future Expansion
This command will grow significantly with additional subcommands for:
- Protocol inspection and debugging
- Performance monitoring and metrics
- Server discovery and registry
- Multi-server orchestration
- Plugin system for custom tools
