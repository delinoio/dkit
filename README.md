# dkit

DevTools by Delino - A collection of developer tools for terminal workflows, automation, and AI-assisted development.

## Installation

```bash
# Build from source
go build -o bin/dkit ./cmd/dkit

# Or install globally
go install github.com/delinoio/dkit/cmd/dkit@latest
```

## Usage

```bash
dkit [command] [subcommand] [flags]
```

## Available Commands

- **clipboard** - Bridge between terminal and system clipboard
- **cron** - Parse, validate, and explain cron expressions
- **env** - Manage environment variables across multiple .env files
- **git** - Git utilities and custom merge drivers
- **jsonc** - Convert JSONC/JSON5 to JSON
- **mcp** - MCP (Model Context Protocol) CLI tool
- **port** - Manage network ports during development
- **retry** - Execute commands with automatic retry logic
- **run** - Execute commands with AI-optimized output and persistent logging
- **yaml** - Normalize YAML files by resolving anchors and aliases

## Quick Examples

### Clipboard Management
```bash
# Copy to clipboard
echo "Hello, World!" | dkit clipboard copy

# Paste from clipboard
dkit clipboard paste
```

### Cron Expression Parsing
```bash
# Parse cron expression
dkit cron parse "0 */6 * * *"

# Calculate next execution times
dkit cron next "0 9 * * 1-5" --count 5
```

### Environment Variable Management
```bash
# List variables from .env files
dkit env list .env

# Merge multiple .env files
dkit env merge .env .env.local .env.production
```

### Port Management
```bash
# Check if port is available
dkit port check 3000

# Kill process on port
dkit port kill 3000
```

### Retry Logic
```bash
# Retry flaky command
dkit retry --attempts 5 -- npm install

# Retry with exponential backoff
dkit retry --delay 1s --max-delay 30s -- curl https://api.example.com
```

### Command Execution with Logging
```bash
# Run command with persistent logging
dkit run -- npm test

# View process logs
dkit mcp process list
dkit mcp process logs <process-id>
```

### YAML Normalization
```bash
# Normalize YAML file
dkit yaml normalize config.yaml

# Use in pipeline
cat deployment.yaml | dkit yaml normalize | kubectl apply -f -
```

## Development

```bash
# Install dependencies
go mod download

# Build
go build -o bin/dkit ./cmd/dkit

# Run tests (when implemented)
go test ./...
```

## Design Principles

### UX First
- Commands must be **short, discoverable, and intuitive**
- Error messages must explain **what failed and how to fix it**
- Every command must return a **meaningful exit code**

### Safety by Default
- Defaults must always be **safe**
- Destructive actions require explicit confirmation (`--force`, `--yes`)
- All user input must be validated

### Extensibility
- Easy to add new subcommands
- Modular architecture for command organization
- Consistent patterns across all commands

### Git Integration
- Custom merge drivers must be automatic and require zero manual intervention
- Git utilities should handle all common package manager lockfiles
- Always regenerate lockfiles using the appropriate package manager
- Prefer "checkout theirs + regenerate" strategy over manual conflict resolution

## License

See [LICENSE](LICENSE) file for details.