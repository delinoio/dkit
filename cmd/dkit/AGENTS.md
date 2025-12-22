# dkit - Delino CLI Tool

## Binary Name
`dkit`

## Architecture
This is the main CLI entry point. All subcommands are organized under `cmd/dkit/internal/cmd/`.

## Logging
`dkit` uses Go's standard `log/slog` package for structured logging. All logging throughout the codebase must use `slog` instead of `fmt.Print*` or `log.Print*` functions.

## Project Root Detection
`dkit` automatically detects the project root directory using git repository detection (locating the `.git` directory). This is used for:
- Determining the working directory when `-w/--workspace` flag is used
- Locating the `bin` directory for PATH management
- Storing persistent logs and process metadata in `<project-root>/.dkit/`

## Available Commands
- `dkit run` - Execute shell commands with AI-optimized output and persistent logging
- `dkit mcp` - MCP (Model Context Protocol) CLI tool with process management
- `dkit git` - Git utilities and custom merge drivers
- `dkit jsonc` - Convert JSONC/JSON5 to standard JSON (pipe-friendly)
- `dkit yaml` - Normalize YAML by resolving anchors, aliases, and merge keys (pipe-friendly)
- `dkit env` - Manage environment variables across .env files (parse, merge, validate, convert)

See respective `AGENTS.md` files in each command directory for detailed specifications.

