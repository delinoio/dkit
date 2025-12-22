# compact-run Command

## Purpose
Execute shell commands and format output for AI consumption by removing unnecessary content.

## Command Signature
```bash
dkit compact-run $ARGS
```

## Behavior
- **Input**: Accepts arbitrary shell command arguments (`$ARGS`)
- **Execution**: Runs the command through a shell interpreter
- **Output Processing**: 
  - Removes ANSI color codes and escape sequences
  - Strips progress bars and spinner animations
  - Removes redundant whitespace and empty lines
  - Filters out debug/verbose logs unless critical
  - Preserves error messages and warnings
  - Maintains stdout/stderr distinction
  
## Output Format
- Clean, parseable text optimized for AI processing
- Error context preserved for debugging
- Exit codes passed through unchanged

## Use Cases
- Running build commands without visual clutter
- Executing tests and extracting only results
- Piping command output to AI tools
- Automating CI/CD with cleaner logs

## Implementation Requirements
- Must handle piped input/output correctly
- Should respect command exit codes
- Must not alter error semantics
- Should handle long-running commands gracefully
- Must support all shell syntax (pipes, redirects, etc.)

## Error Handling
- Command not found: Exit 127 with clear message
- Permission denied: Exit 126 with explanation
- Command failed: Pass through original exit code
- Invalid syntax: Exit 2 with syntax error details
