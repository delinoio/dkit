# mcp Command

## Purpose
MCP (Model Context Protocol) CLI tool for managing and interacting with MCP servers and clients.

## Command Signature
```bash
dkit mcp [subcommand] [options]
```

## Subcommands Overview
This is the main MCP tool with many subcommands to be added.

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
- Must handle connection timeouts gracefully
- Must validate all configurations before use

## Output Format
- Default: Human-readable table format
- `--json`: Machine-readable JSON output
- `--quiet`: Minimal output (IDs/names only)
- `--verbose`: Detailed debugging information

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
This command will grow significantly with additional subcommands.
