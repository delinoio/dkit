# yaml Command

## Purpose
Normalize YAML files by resolving all anchors, aliases, and merge keys to produce a flat, self-contained YAML output. Designed for use in pipes and automation workflows.

## Command Signature
```bash
dkit yaml [subcommand] [options]
```

## Subcommands

### normalize - Normalize YAML Files

#### Purpose
Normalize YAML files by resolving all anchors, aliases, and merge keys to produce a flat, self-contained YAML output.

#### Command Signature
```bash
dkit yaml normalize [file]
```

#### Input/Output Behavior
- **Input Source**:
  - If `file` argument provided: Read from specified file
  - If no argument: Read from stdin (pipe-friendly)
- **Output**: Normalized YAML written to stdout
- **Errors**: Written to stderr

#### Normalization Process

**Anchor and Alias Resolution:**
- **Anchors** (`&anchor`): Expand all anchored values inline
- **Aliases** (`*alias`): Replace with full referenced content
- Result: No `&` or `*` references in output

**Merge Key Resolution:**
- **Merge keys** (`<<:`): Fully expand merged content
- Resolve nested merges recursively
- Handle multiple merge keys correctly
- Maintain proper override semantics (later values override earlier)

**Other Normalizations:**
- Resolve all YAML tags to their canonical forms
- Expand multi-line strings to consistent format
- Normalize boolean/null representations
- Preserve numeric types and precision
- Maintain UTF-8 encoding

#### YAML Features Handled

**Anchors and Aliases:**
```yaml
# Input
base: &base
  name: example
  version: 1.0

prod:
  <<: *base
  env: production

# Output
base:
  name: example
  version: 1.0

prod:
  name: example
  version: 1.0
  env: production
```

**Multiple Merge Keys:**
```yaml
# Input
defaults: &defaults
  timeout: 30
  retry: 3

advanced: &advanced
  parallel: true
  cache: false

config:
  <<: [*defaults, *advanced]
  timeout: 60  # Override

# Output
defaults:
  timeout: 30
  retry: 3

advanced:
  parallel: true
  cache: false

config:
  timeout: 60
  retry: 3
  parallel: true
  cache: false
```

**Nested Anchors:**
```yaml
# Input
database: &db
  host: localhost
  credentials: &creds
    user: admin
    pass: secret

app:
  db: *db
  auth: *creds

# Output
database:
  host: localhost
  credentials:
    user: admin
    pass: secret

app:
  db:
    host: localhost
    credentials:
      user: admin
      pass: secret
  auth:
    user: admin
    pass: secret
```

#### Output Format
- **Style**: Block style (default YAML formatting)
- **Indentation**: 2 spaces (configurable in future)
- **Encoding**: UTF-8
- **Line Endings**: Unix-style (`\n`)
- **Trailing Newline**: Single newline for pipe compatibility
- **Formatting**: Clean, readable YAML without references

#### Exit Codes
- `0` - Successful normalization
- `1` - Invalid YAML syntax
- `2` - File not found (when file argument provided)
- `3` - Unresolvable anchor/alias reference
- `4` - Circular reference detected
- `5` - I/O error (read/write failure)
- `127` - Invalid command usage

#### Error Handling

**Syntax Errors:**
```
[dkit] ERROR: Invalid YAML syntax at line 23, column 5
[dkit] Unexpected indentation level
```

**Unresolvable References:**
```
[dkit] ERROR: Undefined alias reference: *unknown
[dkit] Referenced at line 45, column 8
```

**Circular References:**
```
[dkit] ERROR: Circular reference detected
[dkit] Anchor 'config' references itself directly or indirectly
```

**File Not Found:**
```
[dkit] ERROR: File not found: config.yaml
```

#### Usage Examples

**From File:**
```bash
dkit yaml normalize config.yaml > normalized.yaml
```

**From Stdin (Pipe):**
```bash
cat deployment.yaml | dkit yaml normalize > flat-deployment.yaml
```

**In Shell Pipeline:**
```bash
curl https://example.com/config.yaml | dkit yaml normalize | yq '.services'
```

**Kubernetes Config Normalization:**
```bash
# Normalize Kubernetes manifests with complex anchors
dkit yaml normalize k8s-template.yaml > k8s-manifest.yaml
kubectl apply -f k8s-manifest.yaml
```

**CI/CD Integration:**
```bash
# Normalize configs before deployment
find ./configs -name "*.yaml" -exec dkit yaml normalize {} > {}.normalized \;
```

## Implementation Requirements

### Performance
- Efficient anchor resolution (single-pass when possible)
- Handle large YAML files (10MB+)
- Avoid exponential expansion with shared structures

### Correctness
- Must preserve YAML semantics exactly
- Must handle all YAML 1.2 core schema types
- Must respect merge key precedence rules
- Must detect and report circular references
- Must maintain data type fidelity

### Edge Cases
- Empty documents → empty output
- Multiple documents in one file → normalize each separately
- Self-referencing anchors → error with clear message
- Anchors defined after use → error (YAML requires definition before use)
- Unused anchors → preserve in output (but expanded)
- Deeply nested merges (10+ levels) → handle correctly

### YAML Compatibility
- Support YAML 1.2 specification
- Handle both YAML 1.1 and 1.2 boolean values
- Preserve custom tags in normalized form
- Support all valid YAML scalar types

## Advanced Features

**Multi-Document Support:**
```bash
# Input: Multiple YAML documents separated by ---
---
config: &cfg
  value: 1
---
app:
  <<: *cfg

# Output: Each document normalized independently
---
config:
  value: 1
---
app:
  value: 1
```

**Preserve Comments (Future):**
```bash
dkit yaml normalize --preserve-comments config.yaml
# Future: Keep comments in normalized output
```

## Integration with Other Tools

**yq Integration:**
```bash
dkit yaml normalize config.yaml | yq '.database.host'
```

**Helm/Kubernetes:**
```bash
# Normalize Helm values for inspection
helm template myapp . | dkit yaml normalize > full-manifest.yaml
```

**Configuration Validation:**
```bash
# Normalize then validate with schema
dkit yaml normalize config.yaml | check-jsonschema --schemafile schema.json -
```

**Git Diff Improvement:**
```bash
# Normalize YAML files before diffing to avoid anchor/alias noise
diff <(dkit yaml normalize old.yaml) <(dkit yaml normalize new.yaml)
```

## Design Principles
- **Pipe-friendly**: Reads stdin, writes stdout, errors to stderr
- **Complete resolution**: No anchors, aliases, or merge keys in output
- **Self-contained**: Output YAML can be used without any context
- **Idempotent**: Running normalize twice produces identical output
- **Silent success**: Only output normalized YAML, no progress messages
- **Verbose errors**: Clear error messages to stderr with line numbers

## Security Considerations
- Prevent billion laughs attack (YAML bomb with exponential expansion)
- Limit maximum output size to prevent DoS
- Validate anchor names for malicious patterns
- Sanitize error messages to avoid information disclosure

