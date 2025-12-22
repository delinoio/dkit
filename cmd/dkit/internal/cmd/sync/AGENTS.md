# sync Command

## Purpose
Synchronize files and directories with intelligent change detection, real-time watching, and conflict resolution. A modern, user-friendly alternative to rsync with better defaults and error messages.

## Command Signature
```bash
dkit sync [subcommand] [options]
```

## Subcommands

### copy - One-Time Synchronization

#### Purpose
Copy files from source to destination, synchronizing changes intelligently.

#### Command Signature
```bash
dkit sync copy <source> <destination> [options]
```

**Arguments:**
- `source` - Source directory or file
- `destination` - Destination directory or file

**Options:**
- `--delete` - Delete files in destination not present in source
- `--dry-run` - Show what would be copied without doing it
- `--checksum` - Use checksums instead of timestamps
- `--ignore <patterns>` - Comma-separated patterns to ignore
- `--ignore-file <file>` - File containing ignore patterns (like .gitignore)
- `--verbose` - Show detailed progress
- `--quiet` - Only show errors
- `--format <text|json>` - Output format

#### Behavior

1. **Change Detection**: Compare source and destination
   - By default: Compare modification time and size
   - With `--checksum`: Compare file hashes (slower but more accurate)

2. **Sync Strategy**:
   - Copy new files from source
   - Update modified files in destination
   - Optionally delete files not in source (with `--delete`)
   - Preserve permissions and timestamps

3. **Conflict Handling**:
   - Newer source → overwrites destination
   - Same timestamp → skip (no change)
   - Destination only → keep (unless `--delete`)

#### Output Format

**Text (default):**
```
[dkit] Synchronizing: src/ → dist/

Scanning directories...
  Source: 145 files, 12.5 MB
  Destination: 132 files, 11.2 MB

Changes detected:
  ✓ 15 files to copy (new)
  ↑ 8 files to update (modified)
  ✗ 5 files to delete (with --delete)
  → 117 files unchanged

Copying files...
  [████████████████████] 100% (23/23 files)

Summary:
  Copied: 15 files (1.2 MB)
  Updated: 8 files (450 KB)
  Deleted: 0 files (use --delete to remove)
  Duration: 2.3 seconds
  Transfer rate: 735 KB/s
```

**Dry run:**
```
[dkit] DRY RUN: No files will be modified

Would copy (new):
  ✓ src/components/Button.tsx → dist/components/Button.tsx
  ✓ src/utils/helper.js → dist/utils/helper.js

Would update (modified):
  ↑ src/App.tsx → dist/App.tsx (size changed: 1.2 KB → 1.5 KB)
  ↑ src/styles.css → dist/styles.css (modified 2 hours ago)

Would delete (--delete not specified):
  ✗ dist/old-component.js (no longer in source)

Run without --dry-run to apply these changes
```

**JSON:**
```json
{
  "summary": {
    "source": {"files": 145, "size": 13107200},
    "destination": {"files": 132, "size": 11739289},
    "new": 15,
    "updated": 8,
    "deleted": 0,
    "unchanged": 117
  },
  "changes": [
    {
      "type": "new",
      "path": "src/components/Button.tsx",
      "size": 1234
    }
  ],
  "duration_ms": 2300,
  "transfer_rate_bps": 752640
}
```

### watch - Real-Time Synchronization

#### Purpose
Continuously watch source directory and sync changes to destination in real-time.

#### Command Signature
```bash
dkit sync watch <source> <destination> [options]
```

**Arguments:**
- `source` - Source directory to watch
- `destination` - Destination directory to sync to

**Options:**
- `--delete` - Delete files in destination when deleted from source
- `--ignore <patterns>` - Patterns to ignore
- `--ignore-file <file>` - Ignore file (like .gitignore)
- `--debounce <ms>` - Debounce delay for rapid changes (default: 100ms)
- `--initial-sync` - Do full sync before watching (default: true)
- `--verbose` - Show all file events

#### Behavior

1. **Initial Sync**: Full synchronization on startup (unless disabled)
2. **Watch Events**: Monitor file system for changes
3. **Debouncing**: Group rapid changes to same file
4. **Auto Sync**: Immediately sync detected changes
5. **Error Recovery**: Retry on transient errors

#### Output Format

```
[dkit] Starting watch mode: src/ → dist/

Initial synchronization...
  ✓ Synced 145 files (12.5 MB)

Watching for changes (Ctrl+C to stop)...

[10:30:15] File created: src/components/NewButton.tsx
           ✓ Copied to dist/components/NewButton.tsx (1.2 KB)

[10:31:42] File modified: src/App.tsx
           ↑ Updated dist/App.tsx (1.5 KB)

[10:32:08] File deleted: src/old-component.js
           ✗ Deleted dist/old-component.js

[10:33:20] Directory created: src/hooks/
           ✓ Created dist/hooks/

Stats:
  Uptime: 5 minutes
  Events processed: 12
  Files synced: 8
  Errors: 0
```

#### Use Cases
- Development: Sync source to build directory
- Deployment: Sync local changes to remote server
- Backup: Continuous backup to external drive
- Testing: Mirror test data changes

### diff - Compare Directories

#### Purpose
Show differences between two directories without syncing.

#### Command Signature
```bash
dkit sync diff <dir1> <dir2> [options]
```

**Arguments:**
- `dir1` - First directory
- `dir2` - Second directory

**Options:**
- `--checksum` - Use checksums for comparison
- `--ignore <patterns>` - Patterns to ignore
- `--format <text|json|csv>` - Output format
- `--show-content` - Show content diff for text files

#### Output Format

**Text (default):**
```
[dkit] Comparing: dir1/ ↔ dir2/

Only in dir1/ (15 files):
  ✓ components/Button.tsx
  ✓ utils/helper.js
  ✓ assets/logo.png

Only in dir2/ (8 files):
  ✗ old/deprecated.js
  ✗ temp/cache.json

Different content (12 files):
  ≠ App.tsx (modified)
    dir1: 1,234 bytes, modified 2025-12-23 10:30
    dir2: 1,456 bytes, modified 2025-12-23 09:15
  ≠ config.json (modified)
    dir1: 567 bytes
    dir2: 589 bytes

Identical (95 files):
  → index.html
  → package.json
  → ...

Summary:
  15 files only in dir1
  8 files only in dir2
  12 files differ
  95 files identical
```

**With content diff:**
```
[dkit] Content diff for: App.tsx

--- dir1/App.tsx
+++ dir2/App.tsx
@@ -10,7 +10,7 @@
 function App() {
-  const [count, setCount] = useState(0);
+  const [count, setCount] = useState(10);
   
   return (
```

### merge - Merge Directories with Conflict Resolution

#### Purpose
Merge two directories with intelligent conflict resolution strategies.

#### Command Signature
```bash
dkit sync merge <dir1> <dir2> <output> [options]
```

**Arguments:**
- `dir1` - First directory
- `dir2` - Second directory
- `output` - Output directory for merged result

**Options:**
- `--strategy <newer|dir1|dir2|prompt>` - Conflict resolution strategy
- `--dry-run` - Preview merge without creating output
- `--format <text|json>` - Output format

#### Conflict Resolution Strategies

1. **newer** (default): Use newer file based on modification time
2. **dir1**: Always prefer files from dir1
3. **dir2**: Always prefer files from dir2
4. **prompt**: Ask user for each conflict (interactive)

#### Output Format

```
[dkit] Merging: dir1/ + dir2/ → output/

Analyzing directories...
  dir1: 145 files
  dir2: 132 files
  Conflicts: 12 files

Conflict resolution (strategy: newer):
  ≠ App.tsx: Using dir1 (newer: 2025-12-23 10:30)
  ≠ config.json: Using dir2 (newer: 2025-12-23 11:00)
  ≠ styles.css: Using dir1 (newer: 2025-12-23 09:45)

Merging files...
  [████████████████████] 100%

Summary:
  From dir1 only: 15 files
  From dir2 only: 8 files
  From dir1 (conflicts): 7 files
  From dir2 (conflicts): 5 files
  Identical (no conflict): 120 files
  Total output: 155 files (13.8 MB)
```

**Interactive prompt:**
```
[dkit] Conflict: App.tsx exists in both directories

  dir1/App.tsx:
    Size: 1,234 bytes
    Modified: 2025-12-23 10:30:00
    MD5: a1b2c3d4e5f6...

  dir2/App.tsx:
    Size: 1,456 bytes
    Modified: 2025-12-23 09:15:00
    MD5: f6e5d4c3b2a1...

Which version to use?
  1. Use dir1 (newer)
  2. Use dir2
  3. Keep both (rename)
  4. Skip this file
  5. Show diff

Your choice [1]: 
```

## Remote Sync Support

### SSH/SCP Integration
```bash
# Sync to remote server
dkit sync copy local/ user@server:/path/to/remote/

# Sync from remote server
dkit sync copy user@server:/remote/path/ local/

# Watch and sync to remote
dkit sync watch src/ user@server:/var/www/html/
```

### URL Support
```bash
# HTTP/HTTPS download sync
dkit sync copy https://example.com/files/ local/

# S3-compatible storage (future)
dkit sync copy s3://bucket/path/ local/
```

## Ignore Patterns

### Pattern Syntax
```bash
# Exact match
node_modules

# Wildcards
*.log
temp*

# Directory
build/
dist/

# Negation
!important.log

# Comments
# This is a comment
```

### .syncignore File
```
# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
*.min.js

# OS files
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/

# Logs
*.log
logs/

# Don't ignore this specific file
!important.log
```

### Using .gitignore
```bash
# Automatically use .gitignore patterns
dkit sync copy src/ dist/ --ignore-file .gitignore

# Combine multiple ignore files
dkit sync copy src/ dist/ --ignore-file .gitignore --ignore-file .syncignore
```

## Common Use Cases

### Development Workflow
```bash
# Watch source and build to dist
dkit sync watch src/ dist/ --ignore "*.test.js,__tests__"

# Sync to Docker volume for hot reload
dkit sync watch . /var/lib/docker/volumes/myapp/_data/
```

### Deployment
```bash
# Dry run before deploy
dkit sync copy dist/ user@prod:/var/www/app/ --dry-run

# Deploy with confirmation
dkit sync copy dist/ user@prod:/var/www/app/ --delete --verbose
```

### Backup
```bash
# Incremental backup
dkit sync copy ~/projects /mnt/backup/projects --checksum

# Continuous backup
dkit sync watch ~/Documents /mnt/backup/Documents --delete
```

### Monorepo Sync
```bash
# Sync shared packages
dkit sync copy packages/shared/ apps/web/node_modules/@company/shared/

# Watch mode for development
dkit sync watch packages/shared/ apps/web/node_modules/@company/shared/
```

### Testing
```bash
# Mirror test data
dkit sync copy test-data-source/ test-data-local/

# Compare test outputs
dkit sync diff expected-output/ actual-output/ --show-content
```

## Exit Codes
- `0` - Success, all files synced
- `1` - Partial failure, some files failed to sync
- `2` - Complete failure, no files synced
- `3` - Source not found
- `4` - Destination not accessible
- `127` - Invalid command usage

## Error Handling

### Permission Denied
```
[dkit] ERROR: Permission denied
[dkit] Cannot write to: /path/to/destination/file.txt
[dkit] Try running with appropriate permissions or check directory ownership
```

### Disk Full
```
[dkit] ERROR: No space left on device
[dkit] Failed to copy: large-file.zip (500 MB)
[dkit] Available space: 120 MB
[dkit] Required space: 500 MB
```

### Network Error (Remote Sync)
```
[dkit] ERROR: SSH connection failed
[dkit] Host: user@server
[dkit] Error: Connection timeout after 30s
[dkit] Retrying in 5 seconds... (attempt 2/3)
```

### File Conflicts
```
[dkit] WARNING: Conflicting changes detected
[dkit] File: config.json
[dkit] Both source and destination modified since last sync
[dkit] Source: modified 2 hours ago
[dkit] Destination: modified 1 hour ago
[dkit] Use --strategy to resolve automatically or merge manually
```

## Implementation Requirements

### Performance
- Efficient change detection using file metadata
- Parallel file transfers when possible
- Incremental sync (only changed files)
- Progress reporting for large operations
- Optimized for large file counts (100k+ files)

### Correctness
- Atomic operations (don't leave partial files)
- Preserve file permissions and timestamps
- Handle symbolic links correctly
- Detect and prevent infinite loops
- Verify transfers with checksums (optional)

### Reliability
- Retry on transient failures
- Graceful handling of interruptions (Ctrl+C)
- Resume capability for large transfers
- Transaction log for rollback

### Cross-Platform
- Work on macOS, Linux, Windows
- Handle path differences
- Respect platform-specific file attributes
- Deal with case-sensitive vs case-insensitive filesystems

## Design Principles
- **Safe by default**: Don't delete unless explicitly asked
- **Fast**: Efficient change detection and transfer
- **Reliable**: Verify operations, handle errors gracefully
- **User-friendly**: Clear progress and helpful errors
- **Flexible**: Support local and remote sync
- **Pipe-friendly**: JSON output for automation

