package retry

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/delinoio/dkit/internal/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var (
		attempts          int
		delay             string
		maxDelay          string
		backoff           string
		backoffMultiplier float64
		jitter            bool
		onExit            string
		skipExit          string
		onStderr          string
		skipStderr        string
		timeout           string
		verbose           bool
		workspace         bool
	)

	cmd := &cobra.Command{
		Use:   "retry [flags] -- <command>",
		Short: "Execute commands with retry logic",
		Long: `Execute commands with automatic retry logic, exponential backoff, and 
failure recovery strategies.

Makes flaky commands reliable and reduces manual intervention in CI/CD pipelines.`,
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRetry(args, attempts, delay, maxDelay, backoff, backoffMultiplier,
				jitter, onExit, skipExit, onStderr, skipStderr, timeout, verbose, workspace)
		},
	}

	cmd.Flags().IntVarP(&attempts, "attempts", "n", 3, "Maximum number of retry attempts")
	cmd.Flags().StringVarP(&delay, "delay", "d", "1s", "Initial delay between retries")
	cmd.Flags().StringVar(&maxDelay, "max-delay", "60s", "Maximum delay for exponential backoff")
	cmd.Flags().StringVar(&backoff, "backoff", "exponential", "Backoff strategy (linear|exponential|constant)")
	cmd.Flags().Float64Var(&backoffMultiplier, "backoff-multiplier", 2.0, "Multiplier for exponential backoff")
	cmd.Flags().BoolVar(&jitter, "jitter", false, "Add random jitter to delays")
	cmd.Flags().StringVar(&onExit, "on-exit", "", "Comma-separated exit codes to retry")
	cmd.Flags().StringVar(&skipExit, "skip-exit", "", "Comma-separated exit codes to NOT retry")
	cmd.Flags().StringVar(&onStderr, "on-stderr", "", "Retry if stderr matches regex pattern")
	cmd.Flags().StringVar(&skipStderr, "skip-stderr", "", "Do NOT retry if stderr matches regex pattern")
	cmd.Flags().StringVar(&timeout, "timeout", "", "Timeout for each attempt")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed retry information")
	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "Execute in project root directory")

	return cmd
}

func runRetry(args []string, attempts int, delay, maxDelay, backoff string,
	backoffMultiplier float64, jitter bool, onExit, skipExit, onStderr,
	skipStderr, timeout string, verbose, workspace bool) error {

	if len(args) == 0 {
		utils.PrintError("No command specified")
		return fmt.Errorf("command required")
	}

	// Parse durations
	delayDuration, err := parseDuration(delay)
	if err != nil {
		utils.PrintError("Invalid delay duration: %s", delay)
		return err
	}

	maxDelayDuration, err := parseDuration(maxDelay)
	if err != nil {
		utils.PrintError("Invalid max-delay duration: %s", maxDelay)
		return err
	}

	var timeoutDuration time.Duration
	if timeout != "" {
		timeoutDuration, err = parseDuration(timeout)
		if err != nil {
			utils.PrintError("Invalid timeout duration: %s", timeout)
			return err
		}
	}

	// Parse exit codes
	onExitCodes, err := parseExitCodes(onExit)
	if err != nil {
		utils.PrintError("Invalid --on-exit codes: %s", onExit)
		return err
	}

	skipExitCodes, err := parseExitCodes(skipExit)
	if err != nil {
		utils.PrintError("Invalid --skip-exit codes: %s", skipExit)
		return err
	}

	// Compile regex patterns
	var onStderrRegex, skipStderrRegex *regexp.Regexp
	if onStderr != "" {
		onStderrRegex, err = regexp.Compile(onStderr)
		if err != nil {
			utils.PrintError("Invalid regex pattern in --on-stderr: %s", onStderr)
			return err
		}
	}
	if skipStderr != "" {
		skipStderrRegex, err = regexp.Compile(skipStderr)
		if err != nil {
			utils.PrintError("Invalid regex pattern in --skip-stderr: %s", skipStderr)
			return err
		}
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Print configuration if verbose
	if verbose {
		utils.PrintInfo("Configuration:")
		utils.PrintInfo("  Command: %s", strings.Join(args, " "))
		utils.PrintInfo("  Max attempts: %d", attempts)
		utils.PrintInfo("  Backoff: %s (%.1fx, max %s)", backoff, backoffMultiplier, maxDelay)
		if timeout != "" {
			utils.PrintInfo("  Timeout: %s per attempt", timeout)
		}
	}

	startTime := time.Now()
	var lastExitCode int

	for attempt := 1; attempt <= attempts; attempt++ {
		select {
		case <-sigChan:
			utils.PrintInfo("Interrupted by user")
			return fmt.Errorf("interrupted")
		default:
		}

		utils.PrintInfo("Attempt %d/%d...", attempt, attempts)

		attemptStart := time.Now()
		exitCode, stderr := executeCommand(args, timeoutDuration, workspace, verbose)
		attemptDuration := time.Since(attemptStart)

		lastExitCode = exitCode

		// Check if succeeded
		if exitCode == 0 {
			utils.PrintSuccess("Success!")
			if verbose {
				utils.PrintInfo("Total time: %.1fs (%d attempts)", time.Since(startTime).Seconds(), attempt)
			}
			return nil
		}

		// Log failure
		if verbose {
			utils.PrintInfo("✗ Failed after %.1fs with exit code %d", attemptDuration.Seconds(), exitCode)
			if stderr != "" {
				utils.PrintInfo("Error output: %s", truncate(stderr, 200))
			}
		} else {
			utils.PrintInfo("✗ Failed with exit code %d", exitCode)
		}

		// Check if we should retry
		if attempt < attempts {
			shouldRetry := shouldRetryCommand(exitCode, stderr, onExitCodes, skipExitCodes,
				onStderrRegex, skipStderrRegex)

			if !shouldRetry {
				if verbose {
					utils.PrintInfo("Retry condition not met, stopping")
				}
				break
			}

			// Calculate delay
			waitDuration := calculateDelay(attempt, delayDuration, maxDelayDuration,
				backoff, backoffMultiplier, jitter)

			if verbose {
				utils.PrintInfo("Retry condition met: exit code %d", exitCode)
			}
			utils.PrintInfo("Waiting %s before retry...", formatDuration(waitDuration))

			// Wait with interruptibility
			select {
			case <-time.After(waitDuration):
			case <-sigChan:
				utils.PrintInfo("Interrupted by user")
				return fmt.Errorf("interrupted")
			}
		}
	}

	// All attempts exhausted
	utils.PrintInfo("")
	utils.PrintInfo("All retry attempts exhausted")
	utils.PrintInfo("Command failed after %d attempts", attempts)
	utils.PrintInfo("Total time: %.1fs", time.Since(startTime).Seconds())
	utils.PrintInfo("Last exit code: %d", lastExitCode)

	os.Exit(lastExitCode)
	return nil
}

func executeCommand(args []string, timeout time.Duration, workspace bool, verbose bool) (int, string) {
	var cmd *exec.Cmd

	// Create command
	if len(args) == 1 {
		cmd = exec.Command("sh", "-c", args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

	// Set working directory if workspace flag is set
	if workspace {
		if gitRoot, err := utils.FindProjectRoot(""); err == nil {
			cmd.Dir = gitRoot
		}
	}

	// Connect stdout/stderr
	cmd.Stdout = os.Stdout
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return 1, ""
	}

	// Start command
	if err := cmd.Start(); err != nil {
		if verbose {
			utils.PrintError("Failed to start command: %v", err)
		}
		return 127, ""
	}

	// Read stderr
	stderrBytes := make([]byte, 0, 4096)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderrPipe.Read(buf)
			if n > 0 {
				stderrBytes = append(stderrBytes, buf[:n]...)
				os.Stderr.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	// Handle timeout
	var timeoutErr error
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					return exitErr.ExitCode(), string(stderrBytes)
				}
				return 1, string(stderrBytes)
			}
			return 0, string(stderrBytes)
		case <-ctx.Done():
			cmd.Process.Kill()
			timeoutErr = ctx.Err()
		}

		if timeoutErr != nil {
			utils.PrintInfo("✗ Timeout after %s", formatDuration(timeout))
			return 124, string(stderrBytes)
		}
	} else {
		err := cmd.Wait()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return exitErr.ExitCode(), string(stderrBytes)
			}
			return 1, string(stderrBytes)
		}
		return 0, string(stderrBytes)
	}

	return 0, string(stderrBytes)
}

func shouldRetryCommand(exitCode int, stderr string, onExitCodes, skipExitCodes []int,
	onStderrRegex, skipStderrRegex *regexp.Regexp) bool {

	// Check skip exit codes first
	if len(skipExitCodes) > 0 && contains(skipExitCodes, exitCode) {
		return false
	}

	// Check skip stderr pattern
	if skipStderrRegex != nil && skipStderrRegex.MatchString(stderr) {
		return false
	}

	// Check on exit codes
	if len(onExitCodes) > 0 {
		if !contains(onExitCodes, exitCode) {
			return false
		}
	} else {
		// Default: retry on any non-zero exit code
		if exitCode == 0 {
			return false
		}
	}

	// Check on stderr pattern
	if onStderrRegex != nil {
		if !onStderrRegex.MatchString(stderr) {
			return false
		}
	}

	return true
}

func calculateDelay(attempt int, initialDelay, maxDelay time.Duration,
	strategy string, multiplier float64, useJitter bool) time.Duration {

	var delay time.Duration

	switch strategy {
	case "constant":
		delay = initialDelay
	case "linear":
		delay = time.Duration(attempt) * initialDelay
	case "exponential":
		delay = time.Duration(float64(initialDelay) * math.Pow(multiplier, float64(attempt-1)))
	default:
		delay = initialDelay
	}

	// Apply max delay cap
	if delay > maxDelay {
		delay = maxDelay
	}

	// Apply jitter
	if useJitter {
		jitterAmount := float64(delay) * 0.2 // ±20% jitter
		randomJitter, _ := rand.Int(rand.Reader, big.NewInt(int64(jitterAmount*2)))
		jitterDelta := time.Duration(randomJitter.Int64()) - time.Duration(jitterAmount)
		delay += jitterDelta

		// Ensure delay is positive
		if delay < 0 {
			delay = time.Millisecond
		}
	}

	return delay
}

func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func parseExitCodes(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}

	parts := strings.Split(s, ",")
	codes := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		code, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid exit code: %s", part)
		}
		codes = append(codes, code)
	}

	return codes, nil
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return d.String()
}
