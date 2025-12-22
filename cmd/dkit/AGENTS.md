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

### Core Utilities
- `dkit run` - Execute shell commands with AI-optimized output and persistent logging
- `dkit mcp` - MCP (Model Context Protocol) CLI tool with process management

### Configuration & Data Formats
- `dkit env` - Manage environment variables across .env files (parse, merge, validate, convert)
- `dkit jsonc` - Convert JSONC/JSON5 to standard JSON (pipe-friendly)
- `dkit yaml` - Normalize YAML by resolving anchors, aliases, and merge keys (pipe-friendly)

### Development Tools
- `dkit git` - Git utilities and custom merge drivers
- `dkit port` - Manage network ports (check availability, kill processes, list usage)
- `dkit lockfile` - Analyze and manage package manager lockfiles (diff, why, dedupe, stats)
- `dkit cron` - Parse, validate, and generate cron expressions
- `dkit time` - Parse, convert, and manipulate timestamps and dates

### File & System Operations
- `dkit sync` - Synchronize files and directories with intelligent change detection
- `dkit clipboard` - Bridge terminal and system clipboard (copy, paste, history)
- `dkit qr` - Generate and decode QR codes (URLs, WiFi, vCard)

See respective `AGENTS.md` files in each command directory for detailed specifications.

