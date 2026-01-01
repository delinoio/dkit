# clipboard Command

## Purpose
Bridge the gap between terminal and system clipboard. Pipe data to/from clipboard, manage clipboard history, and transform clipboard content.

## Command Signature
```bash
dkit clipboard [subcommand] [options]
```

## Subcommands

### copy - Copy to Clipboard

#### Purpose
Copy stdin or file content to system clipboard.

#### Command Signature
```bash
dkit clipboard copy [file] [options]
```

**Arguments:**
- `file` - File to copy (optional, uses stdin if not provided)

**Options:**
- `--trim` - Remove leading/trailing whitespace
- `--append` - Append to existing clipboard content
- `--format <text|html|rtf|image>` - Content format

#### Input Sources

**From stdin (pipe):**
```bash
echo "Hello, World!" | dkit clipboard copy
```

**From file:**
```bash
dkit clipboard copy README.md
```

**From command output:**
```bash
ls -la | dkit clipboard copy
```

**From multiple files (concatenated):**
```bash
cat file1.txt file2.txt | dkit clipboard copy
```

#### Output Format

**Default:**
```
[dkit] ✓ Copied to clipboard (13 bytes)
```

**With --trim:**
```
[dkit] ✓ Copied to clipboard (13 bytes, trimmed)
```

**With --append:**
```
[dkit] ✓ Appended to clipboard (27 bytes total)
```

**Quiet mode:**
```bash
echo "data" | dkit clipboard copy --quiet
# No output, only exit code
```

### paste - Paste from Clipboard

#### Purpose
Output clipboard content to stdout or file.

#### Command Signature
```bash
dkit clipboard paste [options]
```

**Options:**
- `--output <file>` - Write to file instead of stdout
- `--append` - Append to file (if --output specified)
- `--format <text|json|base64>` - Output format

#### Output Examples

**To stdout:**
```bash
dkit clipboard paste

Hello, World!
```

**To file:**
```bash
dkit clipboard paste --output data.txt

[dkit] ✓ Clipboard content saved to data.txt (13 bytes)
```

**As JSON:**
```bash
dkit clipboard paste --format json

{
  "content": "Hello, World!",
  "length": 13,
  "lines": 1,
  "format": "text/plain"
}
```

**As Base64:**
```bash
dkit clipboard paste --format base64

SGVsbG8sIFdvcmxkIQ==
```

**In pipeline:**
```bash
dkit clipboard paste | grep "search term"
```

### clear - Clear Clipboard

#### Purpose
Clear the system clipboard.

#### Command Signature
```bash
dkit clipboard clear [options]
```

**Options:**
- `--force` - Skip confirmation

#### Output

```bash
dkit clipboard clear

[dkit] Clear clipboard? [y/N]: y
[dkit] ✓ Clipboard cleared
```

**Force clear:**
```bash
dkit clipboard clear --force

[dkit] ✓ Clipboard cleared
```

### watch - Watch Clipboard Changes

#### Purpose
Monitor clipboard for changes and trigger actions.

#### Command Signature
```bash
dkit clipboard watch [options]
```

**Options:**
- `--command <cmd>` - Command to run on change
- `--interval <ms>` - Polling interval (default: 500ms)
- `--log <file>` - Log clipboard changes to file

#### Output

```bash
dkit clipboard watch

[dkit] Watching clipboard for changes (Ctrl+C to stop)

[14:30:15] Clipboard changed (13 bytes)
           Hello, World!

[14:32:42] Clipboard changed (28 bytes)
           https://github.com/user/repo

[14:35:10] Clipboard changed (45 bytes)
           {
             "key": "value"
           }
```

**With command:**
```bash
dkit clipboard watch --command "notify-send 'Clipboard changed'"

[dkit] Watching clipboard, running command on change...
```

**Log to file:**
```bash
dkit clipboard watch --log clipboard.log

[dkit] Watching clipboard, logging to clipboard.log...
```

## Common Use Cases

### Development Workflow

**Copy command output:**
```bash
ls -la | dkit clipboard copy
```

**Paste and process:**
```bash
dkit clipboard paste | grep "search" | sort | uniq
```

**Quick code formatting:**
```bash
# Copy code, format, paste back
dkit clipboard transform json-pretty
```

### Data Extraction

**Extract URLs from clipboard:**
```bash
dkit clipboard paste | grep -Eo 'https?://[^ ]+' | dkit clipboard copy
```

**Convert clipboard JSON to CSV:**
```bash
dkit clipboard paste | jq -r '.[] | [.name, .email] | @csv' | dkit clipboard copy
```

### Automation

**Monitor clipboard for URLs:**
```bash
dkit clipboard watch --command 'dkit clipboard paste | grep -E "^https?://" && wget $(dkit clipboard paste)'
```

**Auto-backup clipboard:**
```bash
dkit clipboard watch --log ~/clipboard-backup.log
```

### Testing

**Copy test data:**
```bash
cat test-data.json | dkit clipboard copy
# Manually paste in app to test
```

**Save test output:**
```bash
# Copy app output from UI
dkit clipboard paste > test-output.txt
```

### Note Taking

**Quick note to file:**
```bash
# Copy text, then:
dkit clipboard paste >> notes.txt
```

**Clipboard history as notes:**
```bash
dkit clipboard history --all > clipboard-notes.txt
```

## Integration Examples

### With QR Codes
```bash
# Generate QR code from clipboard
dkit clipboard paste | dkit qr generate

# Copy decoded QR code
dkit qr scan qr.png | dkit clipboard copy
```

### With HTTP
```bash
# Copy API response
curl https://api.example.com/data | dkit clipboard copy

# Send clipboard as POST data
dkit clipboard paste | curl -X POST -d @- https://api.example.com/
```

### With Git
```bash
# Copy current git diff
git diff | dkit clipboard copy

# Copy commit hash
git rev-parse HEAD | dkit clipboard copy
```

### With Time
```bash
# Copy current timestamp
dkit time now --format json | jq -r '.unix' | dkit clipboard copy

# Parse timestamp from clipboard
dkit clipboard paste | xargs dkit time convert
```

### With Encryption
```bash
# Copy encrypted
dkit clipboard paste | openssl enc -aes-256-cbc | base64 | dkit clipboard copy

# Decrypt from clipboard
dkit clipboard paste | base64 -d | openssl enc -d -aes-256-cbc
```

## Platform-Specific Behavior

### macOS
- Uses `pbcopy` and `pbpaste` internally
- Supports rich text and images
- Integrates with Universal Clipboard (iCloud)

### Linux
- Uses `xclip` or `xsel` (auto-detected)
- Requires X11 or Wayland
- Falls back to tmux buffer if available

### Windows
- Uses PowerShell clipboard cmdlets
- Supports Windows clipboard history API
- Compatible with WSL

## Storage & Privacy

### History Storage
- Clipboard history stored in `~/.dkit/clipboard/`
- Encrypted at rest (optional)
- Configurable retention period
- Auto-cleanup of old entries

### Privacy Considerations
- History disabled by default
- Option to exclude sensitive patterns (passwords, tokens)
- Clear history on logout (optional)
- Never log clipboard in verbose mode

## Exit Codes
- `0` - Success
- `1` - Clipboard operation failed
- `2` - Clipboard empty (for paste)
- `3` - History not available
- `127` - Invalid command usage

## Error Handling

### Clipboard Not Available
```
[dkit] ERROR: Clipboard not available
[dkit] No clipboard tool found (pbcopy, xclip, xsel)
[dkit] 
[dkit] Install required tools:
[dkit]   macOS: Built-in (pbcopy)
[dkit]   Linux: sudo apt-get install xclip
[dkit]   Windows: Built-in (PowerShell)
```

### Permission Denied
```
[dkit] ERROR: Permission denied
[dkit] Cannot access clipboard
[dkit] On macOS: System Preferences → Security → Privacy → Accessibility
```

### Clipboard Empty
```
[dkit] ERROR: Clipboard is empty
[dkit] Nothing to paste
```

### Invalid Transformation
```
[dkit] ERROR: Transformation failed
[dkit] Operation: json-pretty
[dkit] Error: Invalid JSON in clipboard
[dkit] 
[dkit] Clipboard content:
[dkit] This is not JSON
```

## Implementation Requirements

### Performance
- Fast clipboard access (< 10ms)
- Efficient history storage
- Minimal memory footprint
- Non-blocking clipboard watch

### Correctness
- Preserve exact clipboard content
- Handle binary data correctly
- Support Unicode/UTF-8
- Proper newline handling

### Reliability
- Graceful degradation if clipboard unavailable
- Atomic clipboard operations
- Safe history management
- Error recovery

### Cross-Platform
- Abstract clipboard access
- Platform-specific implementations
- Consistent behavior across platforms
- Handle platform differences gracefully

## Design Principles
- **Pipe-friendly**: Natural stdin/stdout integration
- **Non-intrusive**: Don't automatically start history
- **Secure**: Privacy-conscious design
- **Fast**: Minimal latency for operations
- **Flexible**: Support multiple data formats
