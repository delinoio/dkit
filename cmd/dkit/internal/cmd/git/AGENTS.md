# git Command

## Purpose
Git utilities and custom merge drivers for automated conflict resolution. Provides intelligent merge strategies for generated files like package manager lockfiles.

## Command Signature
```bash
dkit git [subcommand] [options]
```

## Subcommands

### resolve-conflict - Lockfile Merge Driver

#### Purpose
Custom git merge driver that automatically resolves conflicts in package manager lockfiles by regenerating them using the appropriate package manager.

#### Command Signature
```bash
dkit git resolve-conflict %O %A %B %L %P
```

**Git merge driver arguments (automatically provided by git):**
- `%O` - ancestor's version (base)
- `%A` - current version (ours)
- `%B` - other branches' version (theirs)
- `%L` - conflict marker size
- `%P` - pathname in which the merged result will be stored

#### Supported Package Managers

The driver must automatically detect and support all major package managers:

- **npm** - `package-lock.json`
- **yarn** - `yarn.lock`
- **pnpm** - `pnpm-lock.yaml`
- **bun** - `bun.lockb`
- **cargo** - `Cargo.lock`
- **poetry** - `poetry.lock`
- **pipenv** - `Pipfile.lock`
- **composer** - `composer.lock`
- **go** - `go.sum`
- **gradle** - `gradle.lockfile`
- **maven** - `pom.xml.lock` (if using lockfile plugin)
- **swift** - `Package.resolved`
- **cocoapods** - `Podfile.lock`

#### Resolution Strategy

1. **Detect Package Manager**
   - Identify lockfile type from filename/extension
   - Verify corresponding package manager is installed
   - Locate package manifest file (`package.json`, `Cargo.toml`, etc.)

2. **Checkout Theirs**
   - Accept the incoming lockfile (`%B` - theirs)
   - Replace current lockfile with their version
   - This ensures we start with a consistent state

3. **Regenerate Lockfile**
   - Run appropriate package manager install/update command
   - Use non-interactive mode (no prompts)
   - Preserve existing dependency versions from manifest
   - Commands by package manager:
     - npm: `npm install --package-lock-only`
     - yarn: `yarn install --mode update-lockfile`
     - pnpm: `pnpm install --lockfile-only`
     - bun: `bun install`
     - cargo: `cargo generate-lockfile`
     - poetry: `poetry lock --no-update`
     - pipenv: `pipenv lock`
     - composer: `composer update --lock`
     - go: `go mod tidy`
     - gradle: `gradle dependencies --write-locks`
     - swift: `swift package resolve`
     - cocoapods: `pod install --repo-update`

4. **Mark Resolved**
   - If regeneration succeeds, mark conflict as resolved
   - Exit with code 0 (success)
   - Leave the regenerated lockfile in place

#### Error Handling

- **Package manager not found**: Exit 1, output clear error message
- **Regeneration failed**: Exit 1, preserve conflict markers, show package manager error
- **Unknown lockfile type**: Exit 1, list supported lockfile types
- **Manifest file missing**: Exit 1, explain which manifest file is required
- **Network errors**: Exit 1, suggest offline mode or retry
- **Disk space errors**: Exit 1, clear error message

#### Exit Codes
- `0` - Conflict successfully resolved
- `1` - Resolution failed, manual intervention required
- `127` - Package manager not installed

#### Setup Instructions

The command should be configured as a git merge driver in `.gitattributes`:

```gitattributes
# Package manager lockfiles
package-lock.json merge=dkit-lockfile
yarn.lock merge=dkit-lockfile
pnpm-lock.yaml merge=dkit-lockfile
bun.lockb merge=dkit-lockfile
Cargo.lock merge=dkit-lockfile
poetry.lock merge=dkit-lockfile
Pipfile.lock merge=dkit-lockfile
composer.lock merge=dkit-lockfile
go.sum merge=dkit-lockfile
gradle.lockfile merge=dkit-lockfile
Package.resolved merge=dkit-lockfile
Podfile.lock merge=dkit-lockfile
```

And in `.git/config` or `~/.gitconfig`:

```gitconfig
[merge "dkit-lockfile"]
  name = dkit lockfile merge driver
  driver = dkit git resolve-conflict %O %A %B %L %P
```

#### Output Format

**Success:**
```
[dkit] Detected package-lock.json (npm)
[dkit] Checking out theirs version
[dkit] Regenerating lockfile with: npm install --package-lock-only
[dkit] ✓ Lockfile conflict resolved
```

**Failure:**
```
[dkit] Detected package-lock.json (npm)
[dkit] ERROR: npm not found in PATH
[dkit] Please install npm or resolve conflict manually
```

#### Implementation Requirements

- Must be fast (use `--lockfile-only` or equivalent flags when available)
- Must preserve exact dependency versions from manifest
- Must not modify package manifest files
- Must not install actual packages (only update lockfile)
- Must work in CI/CD environments (non-interactive)
- Must handle monorepos with multiple lockfiles
- Must provide clear progress output for debugging
- Should cache package manager detection for performance
- Should support custom package manager paths via environment variables

#### Design Principles

- **Automatic**: Zero manual intervention for common cases
- **Safe**: Always prefer theirs + regenerate over manual merge
- **Transparent**: Clear logging of what's happening
- **Fast**: Use lockfile-only modes when available
- **Reliable**: Fail gracefully with helpful error messages

#### Future Enhancements

- Support for custom package manager configurations
- Parallel lockfile regeneration in monorepos
- Conflict resolution statistics and reporting
- Integration with `dkit run` for monitored execution
- Support for additional lockfile formats
