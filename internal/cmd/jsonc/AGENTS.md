# jsonc Command

## Purpose
Convert JSONC (JSON with Comments) or JSON5 files to standard JSON format. Designed for use in pipes and automation workflows.

## Command Signature
```bash
dkit jsonc compile [file]
```

## Input/Output Behavior
- **Input Source**:
  - If `file` argument provided: Read from specified file
  - If no argument: Read from stdin (pipe-friendly)
- **Output**: Standard JSON written to stdout
- **Errors**: Written to stderr

## Supported Input Formats

### JSONC (JSON with Comments)
- Single-line comments: `// comment`
- Multi-line comments: `/* comment */`
- Trailing commas allowed
- Follows Visual Studio Code JSONC specification

### JSON5
- Comments (single and multi-line)
- Trailing commas in objects and arrays
- Unquoted object keys
- Single-quoted strings
- Multi-line strings with `\` line continuation
- Hexadecimal numbers (`0x` prefix)
- Leading and trailing decimal points (`.5`, `5.`)
- Positive infinity, negative infinity, and NaN
- Explicit plus sign on numbers

## Processing Behavior

### Comment Removal
- Strip all single-line comments (`//`)
- Strip all multi-line comments (`/* */`)
- Preserve comment-like strings within JSON strings

### Trailing Comma Handling
- Remove trailing commas from objects
- Remove trailing commas from arrays
- Maintain valid JSON syntax

### JSON5-Specific Transformations
- Quote all object keys
- Convert single quotes to double quotes
- Expand multi-line strings into single line
- Convert hexadecimal numbers to decimal
- Normalize number formats (`.5` → `0.5`, `5.` → `5.0`)
- Convert special numeric values to null or error (configurable)

### Output Format
- **Default**: Minified JSON (no whitespace)
- **Option** (future): `--pretty` flag for indented output
- **Encoding**: UTF-8
- **Newline**: Single trailing newline for pipe compatibility

## Exit Codes
- `0` - Successful conversion
- `1` - Invalid input (syntax error in JSONC/JSON5)
- `2` - File not found (when file argument provided)
- `3` - I/O error (read/write failure)
- `127` - Invalid command usage

## Error Handling

### Syntax Errors
```
[dkit] ERROR: Invalid JSONC syntax at line 15, column 8
[dkit] Expected ',' or '}' after object property
```

### File Not Found
```
[dkit] ERROR: File not found: config.jsonc
```

### Invalid UTF-8
```
[dkit] ERROR: Input contains invalid UTF-8 sequences
```

## Usage Examples

### From File
```bash
dkit jsonc compile config.jsonc > config.json
```

### From Stdin (Pipe)
```bash
cat config.jsonc | dkit jsonc compile > config.json
```

### In Shell Pipeline
```bash
curl https://example.com/config.jsonc | dkit jsonc compile | jq '.version'
```

### With Error Handling
```bash
if dkit jsonc compile config.jsonc > output.json 2>error.log; then
  echo "Conversion successful"
else
  cat error.log
fi
```

## Implementation Requirements

### Performance
- Stream-based processing for large files (avoid loading entire file in memory)
- Efficient comment stripping without regex backtracking
- Fast JSON validation

### Correctness
- Must preserve JSON semantics exactly
- Must not alter string contents (even if they look like comments)
- Must handle nested structures correctly
- Must validate output is valid JSON

### Edge Cases
- Empty input → empty output or `{}`/`[]` depending on context
- Comments inside strings must be preserved
- Unicode characters must be preserved
- Large numbers must maintain precision
- Nested comments are invalid and should error

### Compatibility
- Compatible with VS Code JSONC parser
- Compatible with JSON5 specification (v2.2.3)
- Output compatible with all standard JSON parsers

## Integration with Other Tools

### jq Integration
```bash
dkit jsonc compile config.jsonc | jq '.databases[] | select(.primary)'
```

### Configuration Management
```bash
# Convert JSONC config to JSON for deployment
dkit jsonc compile src/config.jsonc > dist/config.json
```

### CI/CD Validation
```bash
# Validate JSONC files in CI pipeline
find . -name "*.jsonc" | xargs -I {} dkit jsonc compile {} > /dev/null
```

## Design Principles
- **Pipe-friendly**: Reads stdin, writes stdout, errors to stderr
- **Fast**: No unnecessary processing or validation beyond conversion
- **Strict**: Fail on invalid input rather than guess
- **Silent success**: Only output converted JSON, no progress messages
- **Verbose errors**: Clear error messages to stderr

