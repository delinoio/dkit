# time Command

## Purpose
Parse, convert, and manipulate timestamps and dates. A developer-friendly alternative to the `date` command with intuitive syntax and multiple output formats.

## Command Signature
```bash
dkit time [subcommand] [options]
```

## Subcommands

### now - Current Time in Multiple Formats

#### Purpose
Display current time in various formats simultaneously.

#### Command Signature
```bash
dkit time now [options]
```

**Options:**
- `--format <text|json|custom>` - Output format (default: text)
- `--timezone <tz>` - Display in specific timezone
- `--utc` - Display in UTC

#### Output Format

**Text (default):**
```
[dkit] Current time:

Local:     2025-12-23 14:30:45 EST
UTC:       2025-12-23 19:30:45 UTC
ISO 8601:  2025-12-23T14:30:45-05:00
RFC 3339:  2025-12-23T14:30:45-05:00
Unix:      1735064445
Unix ms:   1735064445000
Relative:  now
```

**JSON:**
```json
{
  "local": "2025-12-23T14:30:45-05:00",
  "utc": "2025-12-23T19:30:45Z",
  "unix": 1735064445,
  "unix_ms": 1735064445000,
  "timezone": "America/New_York",
  "offset": "-05:00"
}
```

**Specific timezone:**
```bash
dkit time now --timezone "Asia/Tokyo"

[dkit] Current time in Asia/Tokyo:

Local:     2025-12-24 04:30:45 JST
UTC:       2025-12-23 19:30:45 UTC
Unix:      1735064445
```

### convert - Convert Time Formats

#### Purpose
Convert timestamps between different formats.

#### Command Signature
```bash
dkit time convert <input> [options]
```

**Arguments:**
- `input` - Time to convert (various formats accepted)

**Options:**
- `--from <format>` - Input format (auto-detected if not specified)
- `--to <format>` - Output format (default: iso8601)
- `--timezone <tz>` - Timezone for output
- `--format <text|json>` - Output style

#### Supported Input Formats

**Auto-detected:**
- Unix timestamp: `1735064445`
- Unix milliseconds: `1735064445000`
- ISO 8601: `2025-12-23T14:30:45-05:00`
- RFC 3339: `2025-12-23T14:30:45Z`
- RFC 2822: `Mon, 23 Dec 2025 14:30:45 -0500`
- Common formats: `2025-12-23`, `12/23/2025`, `Dec 23, 2025`
- Relative: `2 hours ago`, `in 3 days`, `yesterday`, `tomorrow`

#### Supported Output Formats
- `iso8601` - ISO 8601 format (default)
- `rfc3339` - RFC 3339 format
- `rfc2822` - RFC 2822 format
- `unix` - Unix timestamp (seconds)
- `unix-ms` - Unix timestamp (milliseconds)
- `date` - Date only (YYYY-MM-DD)
- `time` - Time only (HH:MM:SS)
- `human` - Human-readable format
- `custom` - Custom format string

#### Output Examples

**Unix to ISO:**
```bash
dkit time convert 1735064445

2025-12-23T19:30:45Z
```

**ISO to Unix:**
```bash
dkit time convert "2025-12-23T14:30:45-05:00" --to unix

1735064445
```

**Relative to absolute:**
```bash
dkit time convert "2 hours ago"

2025-12-23T12:30:45-05:00
```

**With timezone conversion:**
```bash
dkit time convert "2025-12-23 14:30:45" --timezone "Asia/Tokyo"

2025-12-24T04:30:45+09:00
```

**Custom format:**
```bash
dkit time convert 1735064445 --to custom --format-string "%Y-%m-%d %H:%M:%S"

2025-12-23 19:30:45
```

**Multiple formats (verbose):**
```bash
dkit time convert 1735064445 --format text

[dkit] Time conversion:

Input:     1735064445 (Unix timestamp)
ISO 8601:  2025-12-23T19:30:45Z
RFC 3339:  2025-12-23T19:30:45Z
RFC 2822:  Mon, 23 Dec 2025 19:30:45 +0000
Human:     Monday, December 23, 2025 at 7:30:45 PM
Relative:  2 hours ago
```

### parse - Parse Natural Language Time

#### Purpose
Parse natural language time expressions into structured time data.

#### Command Signature
```bash
dkit time parse <expression> [options]
```

**Arguments:**
- `expression` - Natural language time expression

**Options:**
- `--base <time>` - Base time for relative expressions (default: now)
- `--timezone <tz>` - Timezone for parsing
- `--format <text|json>` - Output format

#### Supported Expressions

**Relative:**
- `now`, `today`, `yesterday`, `tomorrow`
- `2 hours ago`, `in 3 days`, `5 minutes from now`
- `last Monday`, `next Friday`, `this weekend`
- `beginning of month`, `end of year`
- `3 weeks from now`, `2 months ago`

**Specific:**
- `2025-12-23`
- `Dec 23, 2025`
- `12/23/2025`
- `23 December 2025`
- `Christmas 2025` (common holidays)

**Time of day:**
- `9am`, `2:30pm`, `14:30`, `noon`, `midnight`
- `9am tomorrow`, `2pm next Friday`

#### Output Format

**Text (default):**
```bash
dkit time parse "2 hours ago"

[dkit] Parsed time expression: "2 hours ago"

Result:    2025-12-23T12:30:45-05:00
Unix:      1735057245
Relative:  2 hours ago
Absolute:  Monday, December 23, 2025 at 12:30:45 PM
```

**JSON:**
```json
{
  "input": "2 hours ago",
  "parsed": "2025-12-23T12:30:45-05:00",
  "unix": 1735057245,
  "components": {
    "year": 2025,
    "month": 12,
    "day": 23,
    "hour": 12,
    "minute": 30,
    "second": 45
  }
}
```

### diff - Calculate Time Difference

#### Purpose
Calculate the difference between two timestamps.

#### Command Signature
```bash
dkit time diff <time1> <time2> [options]
```

**Arguments:**
- `time1` - First time
- `time2` - Second time (default: now)

**Options:**
- `--unit <auto|seconds|minutes|hours|days>` - Output unit (default: auto)
- `--absolute` - Show absolute difference (no negative)
- `--format <text|json>` - Output format

#### Output Format

**Auto unit (default):**
```bash
dkit time diff "2025-12-20" "2025-12-23"

3 days
```

**Specific unit:**
```bash
dkit time diff "2025-12-23 10:00" "2025-12-23 14:30" --unit hours

4.5 hours
```

**Detailed breakdown:**
```bash
dkit time diff "2025-12-01" "2025-12-23" --format text

[dkit] Time difference:

From:      2025-12-01T00:00:00Z
To:        2025-12-23T00:00:00Z

Difference:
  22 days
  528 hours
  31,680 minutes
  1,900,800 seconds

Human readable: 3 weeks, 1 day
```

**JSON:**
```json
{
  "from": "2025-12-01T00:00:00Z",
  "to": "2025-12-23T00:00:00Z",
  "difference": {
    "days": 22,
    "hours": 528,
    "minutes": 31680,
    "seconds": 1900800,
    "human": "3 weeks, 1 day"
  }
}
```

**Relative to now:**
```bash
dkit time diff "2025-12-01"

22 days ago
```

### add - Add Duration to Time

#### Purpose
Add a duration to a specific time.

#### Command Signature
```bash
dkit time add <time> <duration> [options]
```

**Arguments:**
- `time` - Base time (default: now)
- `duration` - Duration to add (e.g., 2h, 3d, 1w)

**Options:**
- `--timezone <tz>` - Timezone for calculation
- `--format <text|json>` - Output format

#### Duration Format
- `s` - seconds (e.g., `30s`)
- `m` - minutes (e.g., `45m`)
- `h` - hours (e.g., `2h`)
- `d` - days (e.g., `7d`)
- `w` - weeks (e.g., `2w`)
- `M` - months (e.g., `3M`)
- `y` - years (e.g., `1y`)

Can combine: `1d12h30m` (1 day, 12 hours, 30 minutes)

#### Output Examples

**Simple addition:**
```bash
dkit time add now 2h

2025-12-23T16:30:45-05:00
```

**Complex duration:**
```bash
dkit time add "2025-12-23 10:00" "1d12h30m"

2025-12-24T22:30:00-05:00
```

**Verbose:**
```bash
dkit time add "2025-12-01" "3w" --format text

[dkit] Time calculation:

Original:  2025-12-01T00:00:00Z
Add:       3 weeks (21 days)
Result:    2025-12-22T00:00:00Z

Difference: 21 days
```

### sub - Subtract Duration from Time

#### Purpose
Subtract a duration from a specific time.

#### Command Signature
```bash
dkit time sub <time> <duration> [options]
```

**Arguments:**
- `time` - Base time (default: now)
- `duration` - Duration to subtract

**Options:**
- `--timezone <tz>` - Timezone for calculation
- `--format <text|json>` - Output format

#### Output Examples

```bash
dkit time sub now 2h

2025-12-23T12:30:45-05:00
```

```bash
dkit time sub "2025-12-25" "1w"

2025-12-18T00:00:00Z
```

### zone - Timezone Information and Conversion

#### Purpose
Get information about timezones and convert between them.

#### Command Signature
```bash
dkit time zone [timezone] [options]
```

**Arguments:**
- `timezone` - Timezone name or abbreviation (optional, shows current if not specified)

**Options:**
- `--list` - List all available timezones
- `--search <query>` - Search for timezones
- `--format <text|json>` - Output format

#### Output Format

**Current timezone:**
```bash
dkit time zone

[dkit] Current timezone:

Name:      America/New_York
Abbr:      EST
Offset:    UTC-5
DST:       Not currently observing
Next DST:  2025-03-09 02:00:00 (begins)
```

**Specific timezone:**
```bash
dkit time zone "Asia/Tokyo"

[dkit] Timezone: Asia/Tokyo

Name:      Asia/Tokyo
Abbr:      JST
Offset:    UTC+9
DST:       Does not observe DST
Current:   2025-12-24 04:30:45 JST
```

**List timezones:**
```bash
dkit time zone --list

[dkit] Available timezones:

America/New_York      (UTC-5) EST
America/Los_Angeles   (UTC-8) PST
Europe/London         (UTC+0) GMT
Europe/Paris          (UTC+1) CET
Asia/Tokyo            (UTC+9) JST
Asia/Shanghai         (UTC+8) CST
Australia/Sydney      (UTC+11) AEDT
...
```

**Search timezones:**
```bash
dkit time zone --search "york"

[dkit] Timezones matching "york":

America/New_York      (UTC-5) EST

Current time: 2025-12-23 14:30:45 EST
```

### format - Format Time with Custom Pattern

#### Purpose
Format time using custom format strings.

#### Command Signature
```bash
dkit time format <time> <pattern> [options]
```

**Arguments:**
- `time` - Time to format (default: now)
- `pattern` - Format pattern

**Options:**
- `--timezone <tz>` - Timezone for output

#### Format Tokens

**Date:**
- `%Y` - Year (4 digits): 2025
- `%y` - Year (2 digits): 25
- `%m` - Month (01-12): 12
- `%B` - Month name: December
- `%b` - Month abbr: Dec
- `%d` - Day (01-31): 23
- `%A` - Weekday: Monday
- `%a` - Weekday abbr: Mon

**Time:**
- `%H` - Hour 24h (00-23): 14
- `%I` - Hour 12h (01-12): 02
- `%M` - Minute (00-59): 30
- `%S` - Second (00-59): 45
- `%p` - AM/PM: PM

**Other:**
- `%z` - Timezone offset: -0500
- `%Z` - Timezone name: EST
- `%%` - Literal %

#### Examples

```bash
dkit time format now "%Y-%m-%d %H:%M:%S"

2025-12-23 14:30:45
```

```bash
dkit time format now "%B %d, %Y at %I:%M %p"

December 23, 2025 at 02:30 PM
```

```bash
dkit time format now "%A, %B %d, %Y"

Monday, December 23, 2025
```

## Common Use Cases

### Development & Debugging

**Check API response timestamp:**
```bash
dkit time convert 1735064445

# Or from clipboard
pbpaste | dkit time convert
```

**Calculate request timeout:**
```bash
dkit time add now 30s
```

**Log file timestamp parsing:**
```bash
grep ERROR app.log | awk '{print $1, $2}' | while read ts; do
  dkit time convert "$ts" --to human
done
```

### CI/CD & Deployment

**Calculate deployment window:**
```bash
START=$(dkit time convert "2025-12-23 22:00" --to unix)
END=$(dkit time add "$START" 4h --to unix)
echo "Deployment window: $START to $END"
```

**Check certificate expiry:**
```bash
EXPIRY=$(openssl x509 -enddate -noout -in cert.pem | cut -d= -f2)
dkit time diff "$EXPIRY" now
```

### Scheduling & Planning

**Find next Monday:**
```bash
dkit time parse "next Monday"
```

**Calculate sprint end date:**
```bash
dkit time add "2025-12-01" "2w"
```

**Check timezone for meeting:**
```bash
MEETING="2025-12-24 09:00"
echo "San Francisco: $(dkit time convert "$MEETING" --timezone America/Los_Angeles)"
echo "Tokyo: $(dkit time convert "$MEETING" --timezone Asia/Tokyo)"
echo "London: $(dkit time convert "$MEETING" --timezone Europe/London)"
```

### Data Processing

**Convert log timestamps:**
```bash
cat access.log | while read line; do
  ts=$(echo "$line" | grep -oE '[0-9]{10}')
  human=$(dkit time convert "$ts" --to human)
  echo "$line -> $human"
done
```

**Generate date ranges:**
```bash
for i in {0..6}; do
  dkit time add "2025-12-01" "${i}d" --format "%Y-%m-%d"
done
```

## Integration Examples

### With jq
```bash
echo '{"timestamp": 1735064445}' | \
  jq --arg ts "$(dkit time convert 1735064445)" '.human_time = $ts'
```

### With cron
```bash
# Calculate next cron execution
CRON_EXPR="0 9 * * 1-5"
NEXT=$(dkit cron next "$CRON_EXPR" --count 1 --format json | jq -r '.executions[0].timestamp')
dkit time convert "$NEXT" --to human
```

### With git
```bash
# Convert git commit timestamp
git log -1 --format="%at" | xargs dkit time convert
```

## Exit Codes
- `0` - Success
- `1` - Invalid time format
- `2` - Invalid timezone
- `3` - Invalid duration
- `127` - Invalid command usage

## Error Handling

### Invalid Time Format
```
[dkit] ERROR: Could not parse time input
[dkit] Input: "invalid-time"
[dkit] 
[dkit] Supported formats:
[dkit]   - Unix timestamp: 1735064445
[dkit]   - ISO 8601: 2025-12-23T14:30:45Z
[dkit]   - Common: 2025-12-23, Dec 23 2025
[dkit]   - Relative: 2 hours ago, tomorrow
```

### Invalid Timezone
```
[dkit] ERROR: Unknown timezone
[dkit] Input: "America/Invalid"
[dkit] 
[dkit] Did you mean?
[dkit]   - America/Indiana/Indianapolis
[dkit]   - America/New_York
[dkit] 
[dkit] Use 'dkit time zone --search Invalid' to search
```

### Ambiguous Date
```
[dkit] WARNING: Ambiguous date format
[dkit] Input: "01/02/2025"
[dkit] Interpreted as: January 2, 2025 (MM/DD/YYYY)
[dkit] 
[dkit] Use ISO format to avoid ambiguity: 2025-01-02
```

## Implementation Requirements

### Performance
- Fast parsing and conversion (< 10ms)
- Efficient timezone database lookups
- Caching for repeated timezone queries

### Correctness
- Accurate timezone handling (including DST)
- Proper leap year/second handling
- Correct relative time calculations
- Handle edge cases (midnight, DST transitions)

### Compatibility
- Support standard time formats (ISO 8601, RFC 3339, etc.)
- Work with IANA timezone database
- Handle platform differences (macOS/Linux/Windows)

### User Experience
- Intuitive natural language parsing
- Clear error messages with suggestions
- Helpful output formats
- Smart defaults

## Design Principles
- **Intuitive syntax**: Natural language where possible
- **Flexible input**: Accept many formats
- **Clear output**: Human-readable by default
- **Pipe-friendly**: JSON for scripting
- **Timezone-aware**: Proper handling of timezones and DST
