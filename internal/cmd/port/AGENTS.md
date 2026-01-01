# port Command

## Purpose
Manage network ports during development. Check port availability, identify processes using ports, and terminate port-blocking processes with safety features.

## Command Signature
```bash
dkit port [subcommand] [options]
```

## Subcommands

### check - Check Port Availability

#### Purpose
Determine if a port is available or in use, with details about the process using it.

#### Command Signature
```bash
dkit port check <port> [options]
```

**Arguments:**
- `port` - Port number to check (1-65535)

**Options:**
- `--json` - Output result in JSON format
- `--quiet` - Exit code only (0 = available, 1 = in use)

#### Output Format

**Available port:**
```
[dkit] Port 3000 is available
```

**Port in use:**
```
[dkit] Port 3000 is in use

Process: node
PID: 12345
Command: node server.js
User: username
Started: 2025-12-23 10:30:15
```

**JSON output:**
```json
{
  "port": 3000,
  "available": false,
  "process": {
    "pid": 12345,
    "name": "node",
    "command": "node server.js",
    "user": "username",
    "started_at": "2025-12-23T10:30:15Z"
  }
}
```

#### Exit Codes
- `0` - Port is available
- `1` - Port is in use
- `2` - Invalid port number
- `127` - Invalid command usage

### list - List Used Ports

#### Purpose
Display all currently used ports with process information.

#### Command Signature
```bash
dkit port list [options]
```

**Options:**
- `--range <start-end>` - Only show ports in range (e.g., `3000-4000`)
- `--listening` - Only show listening ports (default)
- `--all` - Show all network connections
- `--tcp` - TCP ports only (default: both TCP and UDP)
- `--udp` - UDP ports only
- `--format <table|json|csv>` - Output format (default: table)
- `--sort <port|pid|process>` - Sort by column (default: port)

#### Output Format

**Table (default):**
```
PORT   PROTOCOL  PID     PROCESS    COMMAND                USER
3000   TCP       12345   node       node server.js         username
5432   TCP       23456   postgres   /usr/bin/postgres      postgres
8080   TCP       34567   nginx      nginx: master          root
9000   UDP       45678   avahi      avahi-daemon           avahi
```

**JSON:**
```json
{
  "ports": [
    {
      "port": 3000,
      "protocol": "tcp",
      "pid": 12345,
      "process": "node",
      "command": "node server.js",
      "user": "username"
    }
  ],
  "total": 1
}
```

**CSV:**
```csv
port,protocol,pid,process,command,user
3000,tcp,12345,node,node server.js,username
```

#### Filtering Examples
```bash
# Common development ports
dkit port list --range 3000-9000

# Only TCP ports
dkit port list --tcp

# All connections (including established)
dkit port list --all
```

### kill - Terminate Process on Port

#### Purpose
Terminate the process using a specific port. Includes safety confirmations for system processes.

#### Command Signature
```bash
dkit port kill <port> [options]
```

**Arguments:**
- `port` - Port number whose process to terminate

**Options:**
- `--force` - Skip confirmation prompt
- `--signal <TERM|KILL|HUP>` - Signal to send (default: TERM)
- `--timeout <seconds>` - Grace period before SIGKILL (default: 5)

#### Behavior

1. **Port Check**: Verify port is in use
2. **Process Identification**: Find process using the port
3. **Safety Check**: Warn if system process or root-owned
4. **Confirmation**: Require confirmation unless `--force`
5. **Graceful Termination**: Send SIGTERM first
6. **Force Kill**: Send SIGKILL if process doesn't exit within timeout
7. **Verification**: Confirm port is now available

#### Output Format

**Normal flow:**
```
[dkit] Port 3000 is in use by:
[dkit] Process: node (PID: 12345)
[dkit] Command: node server.js
[dkit] 
[dkit] Terminate this process? [y/N]: y
[dkit] Sending SIGTERM to PID 12345...
[dkit] Process terminated successfully
[dkit] Port 3000 is now available
```

**System process warning:**
```
[dkit] WARNING: This is a system process owned by root
[dkit] Process: nginx (PID: 1234)
[dkit] Command: nginx: master process
[dkit] 
[dkit] Are you sure you want to terminate this process? [y/N]: 
```

**With --force:**
```
[dkit] Terminating process 12345 on port 3000...
[dkit] Process terminated successfully
```

**Timeout scenario:**
```
[dkit] Sending SIGTERM to PID 12345...
[dkit] Waiting for process to exit (timeout: 5s)...
[dkit] Process did not exit gracefully
[dkit] Sending SIGKILL to PID 12345...
[dkit] Process terminated successfully
```

#### Exit Codes
- `0` - Process terminated successfully
- `1` - Port not in use (nothing to kill)
- `2` - User cancelled operation
- `3` - Permission denied (cannot kill process)
- `4` - Process termination failed
- `127` - Invalid command usage

#### Safety Features

- **Confirmation required** for all kills (unless `--force`)
- **Extra warning** for system processes
- **Permission check** before attempting kill
- **Graceful shutdown** with SIGTERM before SIGKILL
- **Verification** that port is freed after termination

### find - Find Processes by Port Pattern

#### Purpose
Search for processes using ports matching a pattern or range.

#### Command Signature
```bash
dkit port find <pattern> [options]
```

**Arguments:**
- `pattern` - Port number, range (3000-4000), or wildcard (30*)

**Options:**
- `--format <table|json|list>` - Output format (default: table)

#### Pattern Examples
```bash
# Specific port
dkit port find 3000

# Range
dkit port find 3000-4000

# Wildcard (all ports starting with 3)
dkit port find "3*"

# Multiple ports
dkit port find "3000,8080,9000"
```

#### Output
Uses same format as `dkit port list` but filtered by pattern.

### watch - Monitor Port Activity

#### Purpose
Watch for processes binding to or releasing ports in real-time.

#### Command Signature
```bash
dkit port watch [options]
```

**Options:**
- `--range <start-end>` - Only watch specific port range
- `--interval <seconds>` - Polling interval (default: 1)

#### Output Format
```
[dkit] Watching for port activity (Ctrl+C to stop)
[dkit] 
[2025-12-23 10:30:15] Port 3000 opened by node (PID: 12345)
[2025-12-23 10:31:42] Port 3000 closed (PID: 12345 exited)
[2025-12-23 10:32:01] Port 8080 opened by nginx (PID: 23456)
```

#### Use Cases
- Monitor when dev server actually starts
- Detect port conflicts in real-time
- Track server restart cycles
- Debug port binding issues

## Common Use Cases

### Development Server Port Conflicts
```bash
# Check if dev port is available
dkit port check 3000

# If in use, kill it
dkit port kill 3000 --force

# Start your server
npm run dev
```

### Find All Development Servers
```bash
# List common dev ports
dkit port list --range 3000-9000

# Kill all node processes on dev ports
dkit port list --range 3000-9000 --format json | \
  jq -r '.ports[] | select(.process=="node") | .port' | \
  xargs -I {} dkit port kill {} --force
```

### CI/CD Integration
```bash
# Ensure port is available before test
if ! dkit port check 8080 --quiet; then
  dkit port kill 8080 --force
fi

# Run test server
dkit run npm test
```

### Docker Port Conflicts
```bash
# Find what's using Docker's common ports
dkit port find "5432,3306,6379,27017"

# Kill conflicts before docker-compose up
dkit port kill 5432 --force
docker-compose up
```

## Platform-Specific Behavior

### macOS
- Uses `lsof` for port detection
- Uses `ps` for process details
- Supports both TCP and UDP

### Linux
- Uses `netstat` or `ss` for port detection
- Falls back to `/proc` filesystem
- Full process command line from `/proc/<pid>/cmdline`

### Windows
- Uses `netstat -ano` for port detection
- Uses `tasklist` for process details
- PowerShell integration for enhanced info

## Error Handling

### Invalid Port Number
```
[dkit] ERROR: Invalid port number: 70000
[dkit] Port must be between 1 and 65535
```

### Permission Denied
```
[dkit] ERROR: Permission denied
[dkit] Cannot terminate process 1234 (owned by root)
[dkit] Try running with sudo: sudo dkit port kill 3000
```

### Port Not in Use
```
[dkit] Port 3000 is not in use
```

### Multiple Processes on Same Port
```
[dkit] WARNING: Multiple processes found on port 3000
[dkit] This is unusual and may indicate a configuration issue
[dkit] 
[dkit] Processes:
[dkit] - PID 12345: node server.js
[dkit] - PID 12346: node server.js
[dkit] 
[dkit] Terminate all? [y/N]: 
```

## Implementation Requirements

### Performance
- Fast port scanning using native OS tools
- Efficient process lookup (no full process table scan)
- Cached process information for repeated calls

### Correctness
- Accurate port-to-process mapping
- Handle IPv4 and IPv6 properly
- Detect listening vs. established connections
- Proper signal handling (SIGTERM before SIGKILL)

### Security
- Require confirmation for system process kills
- Check permissions before attempting termination
- Don't expose sensitive process information to unprivileged users
- Validate port numbers to prevent injection attacks

### Cross-Platform
- Detect available tools on each platform
- Graceful degradation if tools unavailable
- Consistent output format across platforms
- Handle platform-specific quirks

## Exit Codes Summary
- `0` - Success
- `1` - Port in use / not in use (context-dependent)
- `2` - Invalid input or user cancelled
- `3` - Permission denied
- `4` - Operation failed
- `127` - Invalid command usage

## Design Principles
- **Safe by default**: Always confirm destructive actions
- **Clear output**: Show what process is using the port
- **Cross-platform**: Work on macOS, Linux, Windows
- **Pipe-friendly**: JSON output for scripting
- **Fast**: Efficient port detection algorithms

