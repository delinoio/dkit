# run Command

## Purpose
Execute shell commands in the foreground with real-time output streaming, AI-optimized logging, and persistent process monitoring through the MCP interface.

## Command Signature
```bash
dkit run [flags] $ARGS
```

## Flags
- `-w, --workspace`: Execute command in project root directory (auto-detected via git)
- `--ignore-local-bin`: Skip adding `<project-root>/bin` to PATH

## Core Behavior

### Project Root Detection
- Automatically detects project root directory using git repository detection
- Used when `-w/--workspace` flag is specified or when determining `bin` directory location

### Execution
- **Input**: Accepts arbitrary shell command arguments (`$ARGS`)
- **Working Directory**: 
  - Default: Current directory
  - With `-w/--workspace`: Project root directory (detected via git)
- **PATH Management**:
  - Default: Prepends `<project-root>/bin` to PATH if it exists
  - With `--ignore-local-bin`: Uses system PATH without modification
- **Execution Mode**: 
  - Runs the command in **FOREGROUND** (NOT as a daemon)
  - Blocks until the command completes
  - Real-time output streamed to terminal
  - Returns the command's exit code
- **Process Management**: Creates a watchable process that can be monitored via MCP

### Output Processing
- **Real-time Streaming**:
  - Command output streams directly to terminal in real-time
  - User sees output exactly as if running command directly
  - Ctrl+C interrupts the foreground process
- **AI Optimization** (for logs only):
  - Removes ANSI color codes and escape sequences
  - Strips progress bars and spinner animations
  - Removes redundant whitespace and empty lines
  - Filters out debug/verbose logs unless critical
  - Preserves error messages and warnings
  - Maintains stdout/stderr distinction

### Persistent Logging
- **Storage Location**: `<project-root>/.dkit/`
- **Log Files**:
  - `<process-id>.stdout.log` - Standard output
  - `<process-id>.stderr.log` - Standard error
  - `<process-id>.meta.json` - Process metadata (command, start time, status, exit code)
- **Retention**: Logs persist after process termination for post-mortem analysis

## Process Watching
- Each `dkit run` execution creates a watchable process
- Process can be monitored via `dkit mcp process *` commands
- Monitoring works regardless of process state (running, stopped, crashed)
- Real-time tailing and historical log access

## Storage Structure
```
<project-root>/.dkit/
├── processes/
│   ├── <pid-1>/
│   │   ├── stdout.log
│   │   ├── stderr.log
│   │   └── meta.json
│   ├── <pid-2>/
│   │   ├── stdout.log
│   │   ├── stderr.log
│   │   └── meta.json
│   └── ...
└── index.json  # Process registry
```

## Metadata Format
```json
{
  "id": "unique-process-id",
  "pid": 12345,
  "command": "npm run build",
  "args": ["npm", "run", "build"],
  "cwd": "/path/to/project",
  "started_at": "2025-12-23T10:30:00Z",
  "ended_at": "2025-12-23T10:32:00Z",
  "status": "completed|running|failed",
  "exit_code": 0,
  "stdout_path": ".dkit/processes/<id>/stdout.log",
  "stderr_path": ".dkit/processes/<id>/stderr.log"
}
```

## Use Cases
- Running build commands with persistent logs
- Executing long-running tests with monitoring capability
- Foreground processes with AI-accessible logs and MCP monitoring
- Debugging failed commands through log history
- Automating CI/CD with traceable execution logs

## Implementation Requirements
- **Must run command in FOREGROUND (not as daemon)**
  - Command execution blocks until completion
  - Real-time stdout/stderr streaming to terminal
  - User can interrupt with Ctrl+C
  - Exit code passed through unchanged
- Must create `.dkit` directory if not exists
- Should write logs in real-time (buffered for performance)
- Must handle piped input/output correctly
- Should respect command exit codes
- Must not alter error semantics
- Should handle long-running commands gracefully
- Must support all shell syntax (pipes, redirects, etc.)
- Must clean up resources on process termination
- Should implement log rotation for very large outputs

## Output Format
- **Terminal**: Clean, parseable text optimized for AI processing
- **Logs**: Raw output preserved in `.dkit/` directory
- **Exit codes**: Passed through unchanged

## Error Handling
- Command not found: Exit 127 with clear message
- Permission denied: Exit 126 with explanation
- Command failed: Pass through original exit code
- Invalid syntax: Exit 2 with syntax error details
- .dkit creation failed: Exit 1 with filesystem error details
- Disk full during logging: Warn but continue execution

## Integration with MCP
- Processes registered in MCP-accessible registry
- MCP commands can query and monitor all `dkit run` processes
- See `mcp/AGENTS.md` for process monitoring commands
