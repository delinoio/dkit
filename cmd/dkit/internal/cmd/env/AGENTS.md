# env Command

## Purpose
Manage environment variables across multiple `.env` files. Parse, merge, validate, and convert environment configurations for different deployment scenarios.

## Command Signature
```bash
dkit env [subcommand] [options]
```

## Subcommands

### list - Display Environment Variables

#### Purpose
Parse and display environment variables from `.env` files in a readable format.

#### Command Signature
```bash
dkit env list [file...] [options]
```

**Arguments:**
- `file...` - One or more `.env` files to parse (default: `.env` in current directory)

**Options:**
- `--format <text|json|yaml|export>` - Output format (default: text)
- `--show-sources` - Show which file each variable comes from
- `--no-expand` - Don't expand variable references (e.g., `${VAR}`)

#### Output Formats

**Text (default):**
```
DATABASE_URL=postgresql://localhost:5432/mydb
API_KEY=abc123xyz
DEBUG=true
```

**JSON:**
```json
{
  "DATABASE_URL": "postgresql://localhost:5432/mydb",
  "API_KEY": "abc123xyz",
  "DEBUG": "true"
}
```

**YAML:**
```yaml
DATABASE_URL: postgresql://localhost:5432/mydb
API_KEY: abc123xyz
DEBUG: true
```

**Export (shell-compatible):**
```bash
export DATABASE_URL='postgresql://localhost:5432/mydb'
export API_KEY='abc123xyz'
export DEBUG='true'
```

**With sources:**
```
DATABASE_URL=postgresql://localhost:5432/mydb  # .env
API_KEY=abc123xyz                               # .env.local
DEBUG=true                                      # .env.development
```

### merge - Combine Multiple .env Files

#### Purpose
Merge multiple `.env` files with proper precedence rules. Later files override earlier ones.

#### Command Signature
```bash
dkit env merge <file1> <file2> [file...] [options]
```

**Arguments:**
- `file1 file2 ...` - Environment files in precedence order (later overrides earlier)

**Options:**
- `--output <file>` - Write to file instead of stdout
- `--format <dotenv|json|yaml|export>` - Output format (default: dotenv)
- `--comment-conflicts` - Add comments showing overridden values

#### Merge Behavior
- Later files override earlier files
- Preserves comments from the last file that defined each variable
- Empty values are considered valid and will override previous values
- Variable references (`${VAR}`) are expanded after merging

#### Example
```bash
# .env
DATABASE_URL=postgresql://localhost:5432/dev
API_KEY=dev-key-123

# .env.production
DATABASE_URL=postgresql://prod-db:5432/prod
DEBUG=false

# Result of: dkit env merge .env .env.production
DATABASE_URL=postgresql://prod-db:5432/prod
API_KEY=dev-key-123
DEBUG=false
```

**With conflict comments:**
```bash
# .env: postgresql://localhost:5432/dev
DATABASE_URL=postgresql://prod-db:5432/prod
API_KEY=dev-key-123
DEBUG=false
```

### validate - Check Environment Configuration

#### Purpose
Validate environment files against a schema or required variables list.

#### Command Signature
```bash
dkit env validate [file...] [options]
```

**Arguments:**
- `file...` - Environment files to validate (default: `.env`)

**Options:**
- `--required <vars>` - Comma-separated list of required variables
- `--required-file <file>` - File containing required variables (one per line)
- `--schema <file>` - JSON schema file for validation
- `--allow-empty` - Allow empty values for required variables
- `--strict` - Fail on warnings (not just errors)

#### Validation Checks
1. **Syntax validation** - Proper `.env` format
2. **Required variables** - All specified variables are present
3. **Empty values** - Flag variables with empty values
4. **Duplicate keys** - Warn about duplicate variable definitions
5. **Invalid references** - Detect unresolvable `${VAR}` references
6. **Schema validation** - If schema provided, validate against JSON schema

#### Output Format

**Success:**
```
[dkit] ✓ .env validated successfully
[dkit] Found 15 variables
[dkit] All required variables present
```

**Failure:**
```
[dkit] ERROR: Validation failed for .env

Missing required variables:
  - DATABASE_URL
  - API_KEY

Invalid variable references:
  - REDIS_URL references undefined ${REDIS_HOST}

Duplicate variables:
  - DEBUG defined at lines 12 and 45

Exit code: 1
```

#### Schema Format
Use JSON Schema to define expected environment variables:

```json
{
  "type": "object",
  "required": ["DATABASE_URL", "API_KEY"],
  "properties": {
    "DATABASE_URL": {
      "type": "string",
      "pattern": "^postgresql://.+"
    },
    "API_KEY": {
      "type": "string",
      "minLength": 10
    },
    "DEBUG": {
      "type": "string",
      "enum": ["true", "false"]
    },
    "PORT": {
      "type": "string",
      "pattern": "^[0-9]+$"
    }
  }
}
```

### get - Retrieve Single Variable

#### Purpose
Get the value of a specific environment variable from `.env` files.

#### Command Signature
```bash
dkit env get <variable> [file...] [options]
```

**Arguments:**
- `variable` - Variable name to retrieve
- `file...` - Environment files to search (default: `.env`)

**Options:**
- `--default <value>` - Default value if variable not found
- `--expand` - Expand variable references (default: true)
- `--no-expand` - Don't expand variable references

#### Output
Prints only the variable value (no formatting) for easy use in scripts.

#### Example
```bash
# .env contains: DATABASE_URL=postgresql://localhost:5432/mydb
$ dkit env get DATABASE_URL
postgresql://localhost:5432/mydb

# With default value
$ dkit env get MISSING_VAR --default "fallback"
fallback
```

#### Exit Codes
- `0` - Variable found
- `1` - Variable not found and no default provided

### set - Update or Add Variable

#### Purpose
Set or update a variable in a `.env` file safely.

#### Command Signature
```bash
dkit env set <variable> <value> [options]
```

**Arguments:**
- `variable` - Variable name to set
- `value` - Value to assign

**Options:**
- `--file <file>` - Environment file to modify (default: `.env`)
- `--create` - Create file if it doesn't exist
- `--quote <always|auto|never>` - Quote behavior (default: auto)
- `--comment <text>` - Add inline comment

#### Behavior
- If variable exists, update in place
- If variable doesn't exist, append to end of file
- Preserves file formatting and comments
- Auto-quotes values with spaces or special characters

#### Example
```bash
# Set simple value
dkit env set API_KEY abc123xyz

# Set value with spaces (auto-quoted)
dkit env set APP_NAME "My Application"

# Add with comment
dkit env set DEBUG true --comment "Enable debug mode in development"

# Result in .env:
API_KEY=abc123xyz
APP_NAME="My Application"
DEBUG=true  # Enable debug mode in development
```

## Supported .env Format

### Basic Syntax
```bash
# Comments start with #
VARIABLE_NAME=value

# Values with spaces need quotes
APP_NAME="My App"

# Single quotes preserve literal values
PATH='/usr/local/bin'

# Multi-line values
PRIVATE_KEY="-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC...
-----END PRIVATE KEY-----"
```

### Variable Expansion
```bash
# Reference other variables
HOME=/home/user
CONFIG_DIR=${HOME}/config

# With default values
REDIS_URL=${REDIS_HOST:-localhost:6379}

# Nested expansion
API_URL=${PROTOCOL:-http}://${API_HOST}:${API_PORT}
```

### Special Cases
- Empty values: `VAR=` (valid, sets empty string)
- No value: `VAR` (invalid, will error)
- Inline comments: `VAR=value  # comment`
- Escaped quotes: `VAR="He said \"Hello\""`
- Export prefix: `export VAR=value` (preserved)

## Common .env File Patterns

### Precedence Order (from lowest to highest)
1. `.env` - Committed to git, contains safe defaults
2. `.env.local` - Local overrides (gitignored)
3. `.env.development` - Development-specific
4. `.env.production` - Production-specific
5. `.env.test` - Test environment

### Typical Merge Commands
```bash
# Development
dkit env merge .env .env.local .env.development

# Production
dkit env merge .env .env.production

# Test
dkit env merge .env .env.test
```

## Exit Codes
- `0` - Success
- `1` - General error (invalid syntax, missing file, etc.)
- `2` - Validation failed
- `3` - Required variable missing
- `4` - Invalid variable reference
- `127` - Invalid command usage

## Error Handling

### Syntax Errors
```
[dkit] ERROR: Invalid .env syntax in .env.local at line 23
[dkit] Expected format: VARIABLE_NAME=value
[dkit] Got: INVALID LINE HERE
```

### Missing File
```
[dkit] ERROR: File not found: .env.production
[dkit] Use --create flag to create the file
```

### Circular References
```
[dkit] ERROR: Circular variable reference detected
[dkit] A=${B} -> B=${C} -> C=${A}
```

### Invalid Schema
```
[dkit] ERROR: Schema validation failed for DATABASE_URL
[dkit] Expected pattern: ^postgresql://.+
[dkit] Got: mysql://localhost:3306/db
```

## Usage Examples

### Development Workflow
```bash
# Check current environment
dkit env list

# Merge development configs
dkit env merge .env .env.local .env.development --output .env.merged

# Validate before deployment
dkit env validate .env.production --required DATABASE_URL,API_KEY

# Get specific value for scripts
DB_URL=$(dkit env get DATABASE_URL)
```

### CI/CD Integration
```bash
# Validate required variables in CI
dkit env validate .env.production \
  --required-file required-vars.txt \
  --strict

# Generate export script
dkit env list .env.production --format export > set-env.sh
source set-env.sh
```

### Docker Integration
```bash
# Convert .env to docker-compose compatible format
dkit env merge .env .env.production --format dotenv > .env.docker

# Use in docker-compose.yml
docker-compose --env-file .env.docker up
```

### JSON/YAML Config Generation
```bash
# Generate JSON config from .env
dkit env list .env --format json > config.json

# Generate YAML config
dkit env list .env --format yaml > config.yaml
```

## Implementation Requirements

### Performance
- Stream-based parsing for large files
- Lazy variable expansion (only when needed)
- Efficient merge without loading all files into memory

### Correctness
- Must handle all edge cases in .env format
- Proper quote handling (single, double, escaped)
- Correct variable expansion with shell-like semantics
- Preserve file integrity when updating

### Security
- Don't log sensitive values in error messages
- Secure file permissions when creating files (0600)
- Prevent command injection through variable expansion
- Sanitize output in different formats

### Compatibility
- Compatible with dotenv libraries (Node.js, Python, Ruby, etc.)
- Support docker-compose `.env` format
- Handle common variations (export prefix, etc.)

## Design Principles
- **Pipe-friendly**: Commands work well in Unix pipelines
- **Safe by default**: Don't overwrite without confirmation
- **Clear errors**: Helpful messages with line numbers
- **Format agnostic**: Easy conversion between formats
- **Shell compatible**: Output can be sourced directly

