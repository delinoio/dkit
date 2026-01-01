# Contributing to dkit

Thank you for your interest in contributing to dkit!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/delinoio/dkit.git
cd dkit
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
make build
```

4. Run the binary:
```bash
./bin/dkit --help
```

## Project Structure

```
dkit/
├── cmd/
│   └── dkit/           # Main entry point
│       └── main.go
├── internal/
│   ├── cmd/            # Command implementations
│   │   ├── root/       # Root command
│   │   ├── clipboard/  # Clipboard management
│   │   ├── cron/       # Cron expression handling
│   │   ├── env/        # Environment variable management
│   │   ├── git/        # Git utilities
│   │   ├── jsonc/      # JSONC to JSON conversion
│   │   ├── mcp/        # MCP protocol tools
│   │   ├── port/       # Port management
│   │   ├── retry/      # Retry logic
│   │   ├── run/        # Command execution
│   │   └── yaml/       # YAML normalization
│   └── utils/          # Shared utilities
│       ├── git.go      # Git-related utilities
│       ├── output.go   # Output formatting
│       ├── process.go  # Process management
│       ├── confirmation.go # User confirmation
│       └── validators.go   # Input validation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Adding a New Command

Each command in dkit follows a consistent pattern. Here's how to add a new command:

1. Create a new directory under `internal/cmd/`:
```bash
mkdir internal/cmd/mycommand
```

2. Create the command file `internal/cmd/mycommand/mycommand.go`:
```go
package mycommand

import (
    "github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mycommand",
        Short: "Short description",
        Long:  `Long description of what the command does.`,
    }

    // Add subcommands if needed
    cmd.AddCommand(newSubCommand())

    return cmd
}

func newSubCommand() *cobra.Command {
    var (
        flag1 string
        flag2 bool
    )

    cmd := &cobra.Command{
        Use:   "subcommand [args]",
        Short: "Subcommand description",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation here
            return nil
        },
    }

    cmd.Flags().StringVar(&flag1, "flag1", "default", "Description")
    cmd.Flags().BoolVar(&flag2, "flag2", false, "Description")

    return cmd
}
```

3. Register the command in `internal/cmd/root/root.go`:
```go
import (
    "github.com/delinoio/dkit/internal/cmd/mycommand"
)

func init() {
    // ... existing commands ...
    rootCmd.AddCommand(mycommand.NewCommand())
}
```

4. Create documentation in `internal/cmd/mycommand/AGENTS.md`

5. Build and test:
```bash
make build
./bin/dkit mycommand --help
```

## Command Design Principles

When implementing a new command, follow these principles:

### UX First
- Commands must be **short, discoverable, and intuitive**
- Error messages must explain **what failed and how to fix it**
- Every command must return a **meaningful exit code**

### Safety by Default
- Defaults must always be **safe**
- Destructive actions require explicit confirmation (`--force`, `--yes`)
- All user input must be validated

### Output Format
- Use `utils.PrintSuccess()` for success messages
- Use `utils.PrintError()` for errors
- Use `utils.PrintWarning()` for warnings
- Use `utils.PrintInfo()` for informational messages
- Support JSON output with `--json` flag when appropriate

### Exit Codes
- `0` - Success
- `1` - General error
- `2` - Invalid input / user cancelled
- `3+` - Command-specific errors
- `127` - Invalid command usage

## Utilities

The `internal/utils/` package provides common functionality:

- **Git utilities**: `FindProjectRoot()`, `IsGitRepository()`, `GetDkitDataDir()`
- **Output**: `PrintSuccess()`, `PrintError()`, `PrintWarning()`, `PrintInfo()`
- **Process management**: `SaveProcessMetadata()`, `LoadProcessMetadata()`, `ListProcesses()`
- **Confirmation**: `Confirm()`, `ConfirmOrExit()`
- **Validators**: `ValidatePort()`, `ValidatePortRange()`

## Testing

When tests are implemented, run them with:
```bash
go test ./...
```

## Code Style

- Follow standard Go formatting: `make fmt`
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

## Pull Request Process

1. Create a feature branch from `main`
2. Implement your changes
3. Add/update documentation
4. Ensure code builds: `make build`
5. Run tests (when available): `go test ./...`
6. Format code: `make fmt`
7. Submit a pull request

## Questions?

If you have questions or need help, please open an issue on GitHub.

