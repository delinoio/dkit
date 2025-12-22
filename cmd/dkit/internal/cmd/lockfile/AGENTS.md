# lockfile Command

## Purpose
Analyze, compare, and manage package manager lockfiles. Understand dependency changes, identify why packages are installed, and optimize lockfile size across all major package managers.

## Command Signature
```bash
dkit lockfile [subcommand] [options]
```

## Subcommands

### diff - Compare Lockfile Changes

#### Purpose
Show a human-readable summary of changes between two lockfile versions. Useful for understanding what changed in a dependency update.

#### Command Signature
```bash
dkit lockfile diff <old-lockfile> [new-lockfile] [options]
```

**Arguments:**
- `old-lockfile` - Previous version of lockfile
- `new-lockfile` - New version of lockfile (optional, defaults to current file)

**Options:**
- `--format <text|json|markdown>` - Output format (default: text)
- `--show-patches` - Include patch version changes
- `--show-dev` - Include devDependencies changes
- `--group-by <package|type|depth>` - Group changes by category
- `--verbose` - Show full dependency trees for changed packages

#### Supported Lockfiles
- `package-lock.json` (npm)
- `yarn.lock` (yarn)
- `pnpm-lock.yaml` (pnpm)
- `bun.lockb` (bun - shows binary diff summary)
- `Cargo.lock` (Rust)
- `poetry.lock` (Python)
- `Pipfile.lock` (Python)
- `composer.lock` (PHP)
- `go.sum` (Go)
- `Gemfile.lock` (Ruby)
- `Package.resolved` (Swift)
- `Podfile.lock` (CocoaPods)

#### Output Format

**Text (default):**
```
[dkit] Lockfile changes: package-lock.json

Added (3):
  ✓ @types/react@18.2.0
  ✓ typescript@5.3.0
  ✓ vite@5.0.0

Updated (5):
  ↑ react: 18.2.0 → 18.3.0
  ↑ react-dom: 18.2.0 → 18.3.0
  ↑ eslint: 8.50.0 → 8.55.0 (minor)
  ↑ @typescript-eslint/parser: 6.0.0 → 7.0.0 (major ⚠️)
  ↑ lodash: 4.17.20 → 4.17.21 (patch - security fix)

Removed (2):
  ✗ webpack@5.88.0
  ✗ webpack-cli@5.1.0

Summary:
  3 added, 5 updated, 2 removed
  1 major update, 1 minor update, 1 patch update
  Total packages: 1,245 → 1,246
```

**Markdown (for PRs):**
```markdown
## Lockfile Changes

### Added Dependencies (3)
- ✓ `@types/react@18.2.0`
- ✓ `typescript@5.3.0`
- ✓ `vite@5.0.0`

### Updated Dependencies (5)
| Package | Old Version | New Version | Change |
|---------|-------------|-------------|--------|
| react | 18.2.0 | 18.3.0 | minor |
| react-dom | 18.2.0 | 18.3.0 | minor |
| eslint | 8.50.0 | 8.55.0 | minor |
| @typescript-eslint/parser | 6.0.0 | 7.0.0 | ⚠️ major |
| lodash | 4.17.20 | 4.17.21 | patch (security) |

### Removed Dependencies (2)
- ✗ `webpack@5.88.0`
- ✗ `webpack-cli@5.1.0`

### Summary
- 3 added, 5 updated, 2 removed
- 1 major update requiring attention
```

**JSON:**
```json
{
  "summary": {
    "added": 3,
    "updated": 5,
    "removed": 2,
    "total_before": 1245,
    "total_after": 1246
  },
  "changes": {
    "added": [
      {"name": "@types/react", "version": "18.2.0"}
    ],
    "updated": [
      {
        "name": "react",
        "old_version": "18.2.0",
        "new_version": "18.3.0",
        "change_type": "minor"
      }
    ],
    "removed": [
      {"name": "webpack", "version": "5.88.0"}
    ]
  }
}
```

#### Git Integration
```bash
# Compare with git HEAD
dkit lockfile diff HEAD:package-lock.json package-lock.json

# Compare between branches
dkit lockfile diff main:package-lock.json feature:package-lock.json

# Show changes in current commit
git show HEAD:package-lock.json | dkit lockfile diff - package-lock.json
```

### why - Explain Why Package is Installed

#### Purpose
Show the dependency tree explaining why a package is installed. Identify which top-level dependencies require a specific package.

#### Command Signature
```bash
dkit lockfile why <package> [lockfile] [options]
```

**Arguments:**
- `package` - Package name to investigate
- `lockfile` - Lockfile to analyze (optional, auto-detected)

**Options:**
- `--all` - Show all dependency paths (not just shortest)
- `--depth <n>` - Limit tree depth
- `--format <tree|list|json>` - Output format (default: tree)
- `--version <version>` - Specific version to analyze

#### Output Format

**Tree (default):**
```
[dkit] Why is "lodash@4.17.21" installed?

Found 3 dependency paths:

Path 1 (direct):
  your-app
  └─ lodash@^4.17.0

Path 2 (via express):
  your-app
  └─ express@4.18.2
     └─ body-parser@1.20.1
        └─ lodash@^4.17.0

Path 3 (via webpack):
  your-app
  └─ webpack@5.88.0
     └─ webpack-cli@5.1.0
        └─ interpret@3.1.1
           └─ lodash@^4.17.0

Summary:
  1 direct dependency, 2 transitive dependencies
  Required by: your-app (direct), express, webpack
```

**List format:**
```
[dkit] Dependency paths for lodash@4.17.21:

1. your-app → lodash@^4.17.0
2. your-app → express@4.18.2 → body-parser@1.20.1 → lodash@^4.17.0
3. your-app → webpack@5.88.0 → webpack-cli@5.1.0 → interpret@3.1.1 → lodash@^4.17.0
```

**JSON:**
```json
{
  "package": "lodash",
  "version": "4.17.21",
  "paths": [
    {
      "path": ["your-app", "lodash@^4.17.0"],
      "type": "direct"
    },
    {
      "path": ["your-app", "express@4.18.2", "body-parser@1.20.1", "lodash@^4.17.0"],
      "type": "transitive"
    }
  ],
  "required_by": ["your-app", "express", "webpack"]
}
```

### dedupe - Find Duplicate Dependencies

#### Purpose
Identify duplicate packages at different versions and suggest deduplication opportunities.

#### Command Signature
```bash
dkit lockfile dedupe [lockfile] [options]
```

**Arguments:**
- `lockfile` - Lockfile to analyze (optional, auto-detected)

**Options:**
- `--fix` - Run package manager deduplication command
- `--dry-run` - Show what would be deduplicated
- `--format <text|json>` - Output format
- `--threshold <bytes>` - Minimum size to report (default: 100KB)

#### Output Format

**Analysis:**
```
[dkit] Analyzing package-lock.json for duplicates...

Found 12 packages with multiple versions:

lodash (3 versions):
  ✓ 4.17.21 - used by 15 packages (preferred)
  ⚠ 4.17.20 - used by 3 packages
  ⚠ 3.10.1 - used by 1 package (outdated)
  Potential savings: ~280 KB

react (2 versions):
  ✓ 18.2.0 - used by 8 packages (preferred)
  ⚠ 17.0.2 - used by 2 packages
  Potential savings: ~340 KB

@babel/core (2 versions):
  ✓ 7.23.0 - used by 10 packages
  ⚠ 7.22.0 - used by 5 packages
  Potential savings: ~180 KB

Summary:
  12 packages with duplicates
  Total potential savings: ~2.1 MB
  Run 'dkit lockfile dedupe --fix' to optimize
```

**With --fix:**
```
[dkit] Running deduplication...
[dkit] Detected package manager: npm
[dkit] Executing: npm dedupe

[dkit] ✓ Deduplication complete
[dkit] Removed 8 duplicate packages
[dkit] Saved ~1.8 MB
```

### stats - Lockfile Statistics

#### Purpose
Show detailed statistics about lockfile contents, sizes, and dependency distribution.

#### Command Signature
```bash
dkit lockfile stats [lockfile] [options]
```

**Arguments:**
- `lockfile` - Lockfile to analyze (optional, auto-detected)

**Options:**
- `--format <text|json>` - Output format
- `--show-top <n>` - Show top N largest packages (default: 10)
- `--group-by <owner|license|size>` - Group statistics

#### Output Format

```
[dkit] Lockfile Statistics: package-lock.json

Overview:
  Total packages: 1,245
  Direct dependencies: 42
  Dev dependencies: 28
  Transitive dependencies: 1,175
  Lockfile size: 1.2 MB
  Total install size: ~450 MB

Version Distribution:
  Major version 0.x: 145 packages (11.6%)
  Major version 1.x: 234 packages (18.8%)
  Major version 2.x+: 866 packages (69.6%)

Top 10 Largest Packages:
  1. typescript - 45.2 MB
  2. webpack - 32.1 MB
  3. @types/node - 28.5 MB
  4. esbuild - 22.8 MB
  5. terser - 18.4 MB
  6. @babel/core - 16.2 MB
  7. react-dom - 14.5 MB
  8. next - 12.8 MB
  9. eslint - 11.2 MB
  10. prettier - 9.8 MB

Dependency Depth:
  Average depth: 4.2
  Maximum depth: 12
  Packages at depth 1: 42
  Packages at depth 2-5: 856
  Packages at depth 6+: 347

License Distribution:
  MIT: 1,105 (88.8%)
  Apache-2.0: 62 (5.0%)
  ISC: 45 (3.6%)
  BSD-3-Clause: 28 (2.2%)
  Other: 5 (0.4%)
```

### validate - Validate Lockfile Integrity

#### Purpose
Check lockfile for common issues, corruption, or inconsistencies with package manifest.

#### Command Signature
```bash
dkit lockfile validate [lockfile] [options]
```

**Arguments:**
- `lockfile` - Lockfile to validate (optional, auto-detected)

**Options:**
- `--strict` - Fail on warnings, not just errors
- `--fix` - Attempt to fix issues automatically
- `--check-manifest` - Verify consistency with package.json/Cargo.toml/etc.

#### Validation Checks

1. **Syntax validation** - Valid JSON/YAML/TOML format
2. **Integrity hashes** - Verify checksums if present
3. **Version consistency** - Same package version resolved consistently
4. **Manifest alignment** - Dependencies match package manifest
5. **Missing packages** - All referenced packages are present
6. **Circular dependencies** - Detect dependency cycles
7. **Deprecated packages** - Flag known deprecated packages
8. **Security vulnerabilities** - Check for known CVEs (optional)

#### Output Format

**Valid lockfile:**
```
[dkit] ✓ package-lock.json is valid
[dkit] No issues found
[dkit] 1,245 packages verified
```

**Issues found:**
```
[dkit] Validating package-lock.json...

ERRORS (2):
  ✗ Missing integrity hash for lodash@4.17.21
  ✗ Version mismatch: react@18.2.0 required by package.json, but 18.1.0 in lockfile

WARNINGS (3):
  ⚠ Package 'request' is deprecated
  ⚠ Circular dependency detected: a → b → c → a
  ⚠ Unused dependency in lockfile: unused-package@1.0.0

Summary:
  2 errors, 3 warnings
  Run 'dkit lockfile validate --fix' to attempt automatic fixes
```

### outdated - Show Outdated Dependencies

#### Purpose
List dependencies that have newer versions available.

#### Command Signature
```bash
dkit lockfile outdated [lockfile] [options]
```

**Arguments:**
- `lockfile` - Lockfile to check (optional, auto-detected)

**Options:**
- `--major` - Only show major version updates
- `--minor` - Only show minor version updates
- `--patch` - Only show patch version updates
- `--security-only` - Only show updates with security fixes
- `--format <table|json>` - Output format

#### Output Format

```
[dkit] Checking for outdated dependencies...

Package                      Current    Wanted     Latest     Type
lodash                       4.17.20    4.17.21    4.17.21    patch
react                        18.2.0     18.2.0     18.3.1     minor
@typescript-eslint/parser    6.0.0      6.21.0     7.1.0      major ⚠️
webpack                      5.88.0     5.89.0     5.89.0     patch

Summary:
  4 packages can be updated
  1 major, 1 minor, 2 patch updates available
  2 security fixes available (lodash, webpack)
```

## Common Use Cases

### Pull Request Review
```bash
# Generate lockfile diff for PR description
dkit lockfile diff main:package-lock.json --format markdown > lockfile-changes.md

# Check for major updates
dkit lockfile diff main:package-lock.json | grep "major ⚠️"
```

### Dependency Investigation
```bash
# Why is this old package still installed?
dkit lockfile why lodash@3.10.1

# Find all paths to a security-vulnerable package
dkit lockfile why minimist --all
```

### Lockfile Optimization
```bash
# Find and remove duplicates
dkit lockfile dedupe --dry-run
dkit lockfile dedupe --fix

# Check for bloat
dkit lockfile stats --show-top 20
```

### CI/CD Validation
```bash
# Validate lockfile in CI
dkit lockfile validate --strict --check-manifest

# Ensure no duplicates in production
if dkit lockfile dedupe --dry-run | grep -q "potential savings"; then
  echo "Duplicates found - run npm dedupe"
  exit 1
fi
```

### Monorepo Management
```bash
# Compare lockfiles across packages
for pkg in packages/*/package-lock.json; do
  echo "Analyzing $pkg"
  dkit lockfile stats "$pkg"
done

# Find common dependencies
dkit lockfile why react packages/app1/package-lock.json
dkit lockfile why react packages/app2/package-lock.json
```

## Exit Codes
- `0` - Success
- `1` - Validation failed or issues found
- `2` - File not found or invalid format
- `3` - Package manager not detected
- `4` - Operation failed
- `127` - Invalid command usage

## Error Handling

### Lockfile Not Found
```
[dkit] ERROR: Lockfile not found
[dkit] Searched for: package-lock.json, yarn.lock, pnpm-lock.yaml, bun.lockb
[dkit] Current directory: /path/to/project
```

### Unsupported Format
```
[dkit] ERROR: Unsupported lockfile format
[dkit] File: custom.lock
[dkit] Supported formats: package-lock.json, yarn.lock, pnpm-lock.yaml, etc.
```

### Invalid Lockfile
```
[dkit] ERROR: Invalid JSON in package-lock.json
[dkit] Line 1245: Unexpected token '}' at column 8
```

## Implementation Requirements

### Performance
- Efficient parsing of large lockfiles (10MB+)
- Incremental diffing algorithm for speed
- Caching of parsed lockfile structures
- Parallel analysis when possible

### Correctness
- Accurate version comparison (semver)
- Proper handling of all lockfile formats
- Correct dependency tree resolution
- Handle edge cases (circular deps, peer deps, optional deps)

### Cross-Platform
- Work with all major package managers
- Handle platform-specific dependencies
- Support different lockfile versions

### Integration
- Work with git for historical comparisons
- Integrate with package registries for outdated checks
- Support monorepo structures
- Compatible with CI/CD pipelines

## Design Principles
- **Package manager agnostic**: Support all major ecosystems
- **Human-readable output**: Clear, actionable information
- **Machine-readable option**: JSON for automation
- **Fast**: Efficient algorithms for large lockfiles
- **Safe**: Don't modify files without explicit --fix flag
- **Git-aware**: Easy integration with version control

## Future Enhancements
- `dkit lockfile visualize` - Interactive dependency graph visualization
- `dkit lockfile audit` - Security audit integration
- `dkit lockfile compress` - Optimize lockfile size
- `dkit lockfile migrate` - Convert between package managers
- `dkit lockfile blame` - Show which commit introduced a dependency
- `dkit lockfile cost` - Calculate bundle size impact
- `--watch` mode for continuous monitoring
- Integration with package vulnerability databases
- Support for custom package registries
- Lockfile merge conflict resolution
- Machine learning for suggesting updates
