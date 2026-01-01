# retry Command

## Purpose
Execute commands with automatic retry logic, exponential backoff, and failure recovery strategies. Makes flaky commands reliable and reduces manual intervention in CI/CD pipelines.

## Command Signature
```bash
dkit retry [flags] $ARGS
```

## Flags
- `-n, --attempts <number>` - Maximum number of retry attempts (default: 3)
- `-d, --delay <duration>` - Initial delay between retries (default: 1s)
- `--max-delay <duration>` - Maximum delay for exponential backoff (default: 60s)
- `--backoff <linear|exponential|constant>` - Backoff strategy (default: exponential)
- `--backoff-multiplier <float>` - Multiplier for exponential backoff (default: 2.0)
- `--jitter` - Add random jitter to delays to prevent thundering herd
- `--on-exit <codes>` - Comma-separated exit codes to retry (default: all non-zero)
- `--skip-exit <codes>` - Comma-separated exit codes to NOT retry
- `--on-stderr <pattern>` - Retry if stderr matches regex pattern
- `--skip-stderr <pattern>` - Do NOT retry if stderr matches regex pattern
- `--timeout <duration>` - Timeout for each attempt (no timeout by default)
- `--verbose` - Show detailed retry information
- `-w, --workspace` - Execute in project root directory (auto-detected via git)

## Core Behavior

### Execution Flow
1. **Execute Command**: Run the provided command
2. **Check Result**: Evaluate exit code and output
3. **Determine Retry**: Apply retry conditions
4. **Wait**: Apply backoff delay if retrying
5. **Repeat**: Continue until success or max attempts reached

### Retry Conditions

Commands are retried when:
- Exit code is non-zero (configurable with `--on-exit`)
- Exit code is NOT in skip list (configurable with `--skip-exit`)
- stderr matches retry pattern (if `--on-stderr` specified)
- stderr does NOT match skip pattern (if `--skip-stderr` specified)

### Backoff Strategies

#### Constant Backoff
Wait the same amount of time between each retry.
```bash
dkit retry --backoff constant --delay 5s -- flaky-command
# Delays: 5s, 5s, 5s
```

#### Linear Backoff
Increase delay linearly with each retry.
```bash
dkit retry --backoff linear --delay 2s -- flaky-command
# Delays: 2s, 4s, 6s
```

#### Exponential Backoff (default)
Multiply delay by backoff multiplier with each retry.
```bash
dkit retry --backoff exponential --delay 1s --backoff-multiplier 2 -- flaky-command
# Delays: 1s, 2s, 4s, 8s, 16s...
```

#### With Jitter
Add randomness to prevent synchronized retries.
```bash
dkit retry --jitter --delay 1s -- flaky-command
# Delays: 1.2s, 2.3s, 3.8s (random variation ±20%)
```

### Timeout Handling
```bash
# Timeout each attempt at 30 seconds
dkit retry --timeout 30s -- long-running-command

# If any attempt exceeds 30s, it's killed and retried
```

## Exit Codes
- `0` - Command succeeded (on any attempt)
- `N` - Command failed after all retries (original exit code from last attempt)
- `124` - Command timed out on all attempts
- `130` - Interrupted by user (Ctrl+C)

## Output Format

### Default (concise)
```
[dkit retry] Attempt 1/3...
[dkit retry] ✗ Failed with exit code 1
[dkit retry] Waiting 1s before retry...

[dkit retry] Attempt 2/3...
[dkit retry] ✗ Failed with exit code 1
[dkit retry] Waiting 2s before retry...

[dkit retry] Attempt 3/3...
[dkit retry] ✓ Success!
```

### Verbose
```
[dkit retry] Configuration:
  Command: npm install
  Max attempts: 3
  Backoff: exponential (2.0x, max 60s)
  Retry on: all non-zero exits
  Timeout: 30s per attempt

[dkit retry] Attempt 1/3 started at 2025-12-23 10:30:00
[dkit retry] Running: npm install
npm ERR! network timeout
[dkit retry] ✗ Failed after 12.3s with exit code 1
[dkit retry] Error output: network timeout
[dkit retry] Retry condition met: exit code 1 (non-zero)
[dkit retry] Waiting 1s before retry...

[dkit retry] Attempt 2/3 started at 2025-12-23 10:30:13
[dkit retry] Running: npm install
added 234 packages in 5.2s
[dkit retry] ✓ Success after 5.2s
[dkit retry] Total time: 18.5s (2 attempts)
```

### On Final Failure
```
[dkit retry] Attempt 3/3...
[dkit retry] ✗ Failed with exit code 1

[dkit retry] All retry attempts exhausted
[dkit retry] Command failed after 3 attempts
[dkit retry] Total time: 45.2s
[dkit retry] Last exit code: 1
```

## Use Cases

### Flaky Network Commands
```bash
# Retry package installation on network failures
dkit retry --attempts 5 --delay 2s -- npm install

# Retry with exponential backoff for API rate limits
dkit retry --delay 1s --max-delay 30s -- curl https://api.example.com/data
```

### CI/CD Pipelines
```bash
# Retry flaky tests
dkit retry --attempts 3 -- npm test

# Retry docker pulls
dkit retry --timeout 60s --attempts 5 -- docker pull image:tag

# Retry deployments
dkit retry --delay 10s --attempts 3 -- kubectl apply -f deploy.yaml
```

### Conditional Retries
```bash
# Only retry on specific exit codes (network errors)
dkit retry --on-exit 7,28,56 -- curl https://example.com

# Don't retry on authentication failures (exit 401)
dkit retry --skip-exit 401 -- some-api-command

# Retry only if stderr contains "timeout"
dkit retry --on-stderr "timeout|timed out" -- flaky-command

# Don't retry if stderr contains "permission denied"
dkit retry --skip-stderr "permission denied|unauthorized" -- secure-command
```

### Database Operations
```bash
# Retry database migrations with backoff
dkit retry --delay 5s --attempts 10 -- db-migrate up

# Retry connection with linear backoff
dkit retry --backoff linear --delay 2s -- psql -c "SELECT 1"
```

### File Downloads
```bash
# Retry large file download with jitter to avoid server overload
dkit retry --jitter --delay 1s --attempts 10 -- wget https://example.com/large-file.zip

# Retry with timeout per attempt
dkit retry --timeout 120s --attempts 5 -- rsync -av remote:/data ./data
```

## Advanced Examples

### Kubernetes Deployment
```bash
# Wait for deployment with progressive backoff
dkit retry \
  --attempts 20 \
  --delay 5s \
  --max-delay 60s \
  --timeout 30s \
  -- kubectl wait --for=condition=ready pod -l app=myapp
```

### Multi-Region Fallback
```bash
# Try primary region, then fallback
dkit retry --attempts 2 --delay 0s -- deploy-to-us-east-1 || \
dkit retry --attempts 2 --delay 0s -- deploy-to-us-west-2 || \
dkit retry --attempts 2 --delay 0s -- deploy-to-eu-west-1
```

### Smart Test Retries
```bash
# Retry failed tests only
dkit retry \
  --attempts 3 \
  --skip-exit 0 \
  --on-stderr "FAILED|ERROR" \
  -- pytest --last-failed
```

### Rate Limited API
```bash
# Respect rate limits with exponential backoff
dkit retry \
  --attempts 10 \
  --delay 1s \
  --max-delay 300s \
  --on-stderr "rate limit|429" \
  --jitter \
  -- api-client fetch-data
```

## Integration with `dkit run`

Can be combined with `dkit run` for persistent logging:
```bash
dkit run -- dkit retry --attempts 3 -- npm test

# Logs stored in .dkit/ with full retry history
```

## Implementation Requirements

### Core
- Must capture and preserve stdout/stderr separately
- Must maintain exit code semantics
- Must handle signals properly (forward to child, cleanup on SIGINT)
- Must support all shell syntax when command contains pipes/redirects
- Should not buffer output excessively (stream in real-time when possible)

### Timing
- Must accurately track attempt duration
- Must respect timeout per attempt (not total timeout)
- Must implement backoff strategies correctly
- Should add jitter using cryptographically secure random

### Conditionals
- Must support exit code matching (exact values and ranges)
- Must support regex pattern matching on stderr
- Must handle edge cases (empty stderr, no output, etc.)
- Should compile regex patterns once for performance

### Output
- Must clearly indicate which attempt is running
- Must show why retry was triggered (in verbose mode)
- Must display total time and attempt count on completion
- Should provide actionable information on final failure

### Safety
- Must not retry indefinitely (require explicit attempt count)
- Must respect max delay to prevent excessive waiting
- Must allow user interruption (Ctrl+C) at any time
- Should warn if backoff delay exceeds reasonable limits

## Error Handling

### Command Not Found
```
[dkit retry] ERROR: Command not found: nonexistent-command
[dkit retry] No retries will be attempted for command not found errors
```
Exit code: 127

### Invalid Configuration
```
[dkit retry] ERROR: Invalid delay duration: "abc"
[dkit retry] Expected format: 1s, 500ms, 1m30s
```
Exit code: 2

### Timeout on All Attempts
```
[dkit retry] Attempt 3/3...
[dkit retry] ✗ Timeout after 30s
[dkit retry] All attempts timed out
[dkit retry] Total time: 90s (3 × 30s)
```
Exit code: 124

### User Interruption
```
[dkit retry] Attempt 2/3...
[dkit retry] Waiting 5s before retry...
^C
[dkit retry] Interrupted by user
[dkit retry] Command did not complete (1 success, 1 failure)
```
Exit code: 130

### Invalid Regex Pattern
```
[dkit retry] ERROR: Invalid regex pattern in --on-stderr: "[invalid"
[dkit retry] Error: unclosed character class
```
Exit code: 2

## Design Principles

- **Reliable**: Make flaky commands succeed without manual intervention
- **Transparent**: Clear visibility into what's happening and why
- **Configurable**: Flexible retry strategies for different scenarios
- **Safe**: Sensible defaults, prevent infinite retries
- **Efficient**: Minimal overhead, smart backoff strategies
- **Composable**: Works well with other dkit commands and shell tools

## Future Enhancements

- Circuit breaker pattern (stop retrying if too many consecutive failures)
- Success rate tracking and reporting
- Retry budget enforcement (max total time across all attempts)
- Adaptive backoff based on error type
- Integration with monitoring/alerting systems
- Retry statistics export (JSON format)
- Support for retry policies via config file
- Multi-command retry (try alternative commands on failure)
