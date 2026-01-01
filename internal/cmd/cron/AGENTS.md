# cron Command

## Purpose
Parse, validate, and explain cron expressions. Generate cron schedules and calculate execution times. Make cron scheduling more accessible and less error-prone.

## Command Signature
```bash
dkit cron [subcommand] [options]
```

## Subcommands

### parse - Parse Cron Expression

#### Purpose
Convert a cron expression into human-readable description.

#### Command Signature
```bash
dkit cron parse <expression> [options]
```

**Arguments:**
- `expression` - Cron expression (5 or 6 fields)

**Options:**
- `--format <text|json>` - Output format (default: text)
- `--verbose` - Include detailed field breakdown

#### Supported Formats
- **Standard cron** (5 fields): `* * * * *` (minute hour day month weekday)
- **Extended cron** (6 fields): `* * * * * *` (second minute hour day month weekday)
- **Special strings**: `@yearly`, `@monthly`, `@weekly`, `@daily`, `@hourly`, `@reboot`

#### Output Format

**Text (default):**
```
[dkit] Cron expression: 0 */6 * * *

Human readable:
  Every 6 hours, at minute 0

Schedule breakdown:
  Minute: 0 (at the start of the hour)
  Hour: Every 6 hours (0, 6, 12, 18)
  Day: Every day
  Month: Every month
  Weekday: Any day of the week

Next 5 executions:
  2025-12-23 12:00:00
  2025-12-23 18:00:00
  2025-12-24 00:00:00
  2025-12-24 06:00:00
  2025-12-24 12:00:00
```

**JSON:**
```json
{
  "expression": "0 */6 * * *",
  "description": "Every 6 hours, at minute 0",
  "fields": {
    "minute": {"value": "0", "description": "at minute 0"},
    "hour": {"value": "*/6", "description": "every 6 hours"},
    "day": {"value": "*", "description": "every day"},
    "month": {"value": "*", "description": "every month"},
    "weekday": {"value": "*", "description": "any weekday"}
  },
  "next_executions": [
    "2025-12-23T12:00:00Z",
    "2025-12-23T18:00:00Z"
  ]
}
```

#### Examples
```bash
# Simple hourly
dkit cron parse "0 * * * *"
# Output: Every hour, at minute 0

# Complex expression
dkit cron parse "30 2 * * 1-5"
# Output: At 02:30 AM, Monday through Friday

# Range with step
dkit cron parse "0 9-17/2 * * *"
# Output: Every 2 hours from 9 AM to 5 PM, at minute 0

# Special string
dkit cron parse "@daily"
# Output: Once a day at midnight (0 0 * * *)
```

### next - Calculate Next Execution Times

#### Purpose
Calculate when a cron job will run next.

#### Command Signature
```bash
dkit cron next <expression> [options]
```

**Arguments:**
- `expression` - Cron expression

**Options:**
- `--count <n>` - Number of future executions to show (default: 5)
- `--from <datetime>` - Calculate from specific time (default: now)
- `--timezone <tz>` - Timezone for calculations (default: local)
- `--format <text|json|csv>` - Output format

#### Output Format

**Text (default):**
```
[dkit] Next executions for: 0 9 * * 1-5

1. Monday, December 23, 2025 at 09:00:00 AM
   (in 2 hours, 15 minutes)

2. Tuesday, December 24, 2025 at 09:00:00 AM
   (in 1 day, 2 hours)

3. Wednesday, December 25, 2025 at 09:00:00 AM
   (in 2 days, 2 hours)

4. Thursday, December 26, 2025 at 09:00:00 AM
   (in 3 days, 2 hours)

5. Friday, December 27, 2025 at 09:00:00 AM
   (in 4 days, 2 hours)
```

**JSON:**
```json
{
  "expression": "0 9 * * 1-5",
  "timezone": "America/New_York",
  "executions": [
    {
      "datetime": "2025-12-23T09:00:00-05:00",
      "timestamp": 1735045200,
      "relative": "in 2 hours, 15 minutes"
    }
  ]
}
```

**CSV:**
```csv
datetime,timestamp,relative
2025-12-23 09:00:00,1735045200,in 2 hours
2025-12-24 09:00:00,1735131600,in 1 day
```

### validate - Validate Cron Expression

#### Purpose
Check if a cron expression is valid and identify issues.

#### Command Signature
```bash
dkit cron validate <expression> [options]
```

**Arguments:**
- `expression` - Cron expression to validate

**Options:**
- `--strict` - Enable strict validation (no non-standard extensions)
- `--system <linux|macos|freebsd>` - Validate for specific cron implementation

#### Output Format

**Valid expression:**
```
[dkit] ✓ Valid cron expression: 0 9 * * 1-5
[dkit] Type: Standard cron (5 fields)
[dkit] Description: At 09:00 AM, Monday through Friday
```

**Invalid expression:**
```
[dkit] ✗ Invalid cron expression: 0 25 * * *

Error: Hour value out of range
  Field: hour
  Value: 25
  Valid range: 0-23

Suggestion: Did you mean "0 5 * * *" (05:00 AM)?
```

**Warning (non-standard):**
```
[dkit] ⚠ Cron expression has warnings: @every 1h

Warning: Non-standard syntax
  This syntax is supported by some cron implementations but not all
  Consider using standard format: 0 * * * *
```

#### Validation Checks
1. **Field count** - 5 or 6 fields (or special string)
2. **Value ranges** - Each field within valid range
3. **Step values** - Valid step syntax (`*/n`)
4. **Range syntax** - Valid range format (`1-5`)
5. **List syntax** - Valid list format (`1,3,5`)
6. **Special characters** - Valid use of `*`, `?`, `L`, `W`, `#`
7. **Day/weekday conflict** - Proper use of day and weekday fields

#### Exit Codes
- `0` - Valid expression
- `1` - Invalid expression
- `2` - Valid but has warnings

### generate - Generate Cron Expression

#### Purpose
Create cron expressions using natural language or interactive prompts.

#### Command Signature
```bash
dkit cron generate [options]
```

**Options:**
- `--every <duration>` - Simple interval (1h, 30m, 1d, etc.)
- `--at <time>` - Specific time (09:00, 2:30pm)
- `--on <days>` - Specific days (mon,wed,fri or 1,15)
- `--interactive` - Interactive prompt mode
- `--describe "<text>"` - Natural language description

#### Generation Examples

**Simple intervals:**
```bash
# Every hour
dkit cron generate --every 1h
# Output: 0 * * * *

# Every 30 minutes
dkit cron generate --every 30m
# Output: */30 * * * *

# Every day at midnight
dkit cron generate --every 1d
# Output: 0 0 * * *
```

**Specific times:**
```bash
# Every day at 9 AM
dkit cron generate --at 09:00
# Output: 0 9 * * *

# Weekdays at 2:30 PM
dkit cron generate --at 14:30 --on mon-fri
# Output: 30 14 * * 1-5

# First day of month at noon
dkit cron generate --at 12:00 --on 1
# Output: 0 12 1 * *
```

**Interactive mode:**
```bash
dkit cron generate --interactive

[dkit] Cron Expression Generator

How often should this run?
  1. Every N minutes/hours/days
  2. At specific time(s)
  3. Custom expression

Your choice: 2

What time? (e.g., 09:00, 2:30pm): 09:30

Which days?
  1. Every day
  2. Weekdays (Mon-Fri)
  3. Weekends (Sat-Sun)
  4. Specific days

Your choice: 2

Generated expression: 30 9 * * 1-5
Description: At 09:30 AM, Monday through Friday

Is this correct? [Y/n]: y

30 9 * * 1-5
```

**Natural language (future):**
```bash
dkit cron generate --describe "every weekday at 9:30am"
# Output: 30 9 * * 1-5

dkit cron generate --describe "every 6 hours"
# Output: 0 */6 * * *
```

## Common Use Cases

### Understanding Existing Cron Jobs
```bash
# Parse crontab entries
crontab -l | grep -v '^#' | while read line; do
  expr=$(echo "$line" | awk '{print $1,$2,$3,$4,$5}')
  cmd=$(echo "$line" | cut -d' ' -f6-)
  echo "Command: $cmd"
  dkit cron parse "$expr"
  echo ""
done
```

### Creating New Cron Jobs
```bash
# Generate expression interactively
expr=$(dkit cron generate --interactive)

# Validate before adding
dkit cron validate "$expr" && echo "$expr /path/to/script.sh" | crontab -
```

### Debugging Cron Schedule
```bash
# Why didn't my job run?
dkit cron next "0 9 * * 1-5" --from "2025-12-23 08:00"

# Check if two jobs overlap
dkit cron diff "0 */6 * * *" "0 9-17 * * 1-5"
```

### CI/CD Integration
```bash
# Validate cron expressions in config files
yq '.jobs[].schedule' .github/workflows/scheduled.yml | \
  xargs -I {} dkit cron validate {}
```

### Documentation
```bash
# Generate documentation for cron jobs
echo "# Scheduled Jobs" > CRON.md
crontab -l | grep -v '^#' | while read line; do
  expr=$(echo "$line" | awk '{print $1,$2,$3,$4,$5}')
  cmd=$(echo "$line" | cut -d' ' -f6-)
  echo "## $cmd" >> CRON.md
  dkit cron parse "$expr" --format markdown >> CRON.md
done
```

## Cron Expression Reference

### Field Format
```
* * * * *
│ │ │ │ │
│ │ │ │ └─── Weekday (0-7, 0 and 7 are Sunday)
│ │ │ └───── Month (1-12)
│ │ └─────── Day (1-31)
│ └───────── Hour (0-23)
└─────────── Minute (0-59)
```

### Special Characters
- `*` - Any value
- `,` - List of values (1,3,5)
- `-` - Range of values (1-5)
- `/` - Step values (*/2 = every 2)

### Special Strings
- `@yearly` or `@annually` - `0 0 1 1 *`
- `@monthly` - `0 0 1 * *`
- `@weekly` - `0 0 * * 0`
- `@daily` or `@midnight` - `0 0 * * *`
- `@hourly` - `0 * * * *`
- `@reboot` - Run at startup

### Common Patterns
```bash
# Every minute
* * * * *

# Every hour
0 * * * *

# Every day at midnight
0 0 * * *

# Every Monday at 9 AM
0 9 * * 1

# First day of month
0 0 1 * *

# Every 15 minutes
*/15 * * * *

# Business hours (9-5, Mon-Fri)
0 9-17 * * 1-5

# Twice a day
0 9,21 * * *
```

## Exit Codes
- `0` - Success
- `1` - Invalid expression
- `2` - Valid with warnings
- `127` - Invalid command usage

## Error Handling

### Invalid Field Value
```
[dkit] ERROR: Invalid minute value in cron expression
[dkit] Expression: 75 * * * *
[dkit] Field: minute (position 1)
[dkit] Value: 75
[dkit] Valid range: 0-59
```

### Syntax Error
```
[dkit] ERROR: Invalid cron syntax
[dkit] Expression: 0 9 * *
[dkit] Expected 5 or 6 fields, got 4
[dkit] Format: minute hour day month weekday
```

### Ambiguous Expression
```
[dkit] WARNING: Day and weekday both specified
[dkit] Expression: 0 9 15 * 1
[dkit] This will run on the 15th AND every Monday
[dkit] Use '?' in one field to avoid confusion (non-standard)
```

## Implementation Requirements

### Performance
- Fast expression parsing (< 1ms)
- Efficient next execution calculation
- Caching for repeated queries

### Correctness
- Accurate cron parsing for all standard formats
- Proper timezone handling
- Correct calculation of next execution times
- Handle edge cases (leap years, DST transitions)

### Compatibility
- Support standard POSIX cron
- Recognize common extensions (6-field, special strings)
- Platform-specific validation (Linux/macOS/BSD differences)

### User Experience
- Clear, actionable error messages
- Helpful suggestions for common mistakes
- Human-readable output by default
- Machine-readable JSON for scripting

## Design Principles
- **Educational**: Help users understand cron syntax
- **Safe**: Validate before executing
- **Helpful**: Suggest corrections for errors
- **Cross-platform**: Work consistently across systems
- **Pipe-friendly**: JSON output for automation

