## Core principles

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
