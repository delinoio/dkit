package run

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/delinoio/dkit/internal/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var (
		workspace      bool
		ignoreLocalBin bool
	)

	cmd := &cobra.Command{
		Use:   "run [flags] -- <command>",
		Short: "Execute commands with logging",
		Long: `Execute shell commands in the FOREGROUND with real-time output streaming.

Features:
- Runs command in foreground (NOT as daemon)
- Real-time stdout/stderr streaming to terminal
- Persistent logging to .dkit/processes/
- AI-optimized log processing
- Process monitoring through MCP interface
- Respects original exit codes`,
		DisableFlagParsing: false,
		SilenceUsage:       true, // Don't show usage on command errors
		SilenceErrors:      true, // We'll handle errors ourselves
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no command specified")
			}

			err := runCommand(args, workspace, ignoreLocalBin)
			if err != nil {
				// If it's an exit error, exit with the same code
				if exitError, ok := err.(*exec.ExitError); ok {
					os.Exit(exitError.ExitCode())
				}
				// For other errors, return them
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "Execute in project root directory")
	cmd.Flags().BoolVar(&ignoreLocalBin, "ignore-local-bin", false, "Skip adding <project-root>/bin to PATH")

	return cmd
}

func runCommand(args []string, workspace, ignoreLocalBin bool) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find project root
	projectRoot, err := utils.FindProjectRoot("")
	if err != nil {
		// Not in a git repository, use current directory
		projectRoot = cwd
	}

	// Determine working directory
	workDir := cwd
	if workspace {
		workDir = projectRoot
	}

	// Generate process ID
	processID := utils.GenerateProcessID()

	// Setup .dkit directory and log files
	dataDir, err := utils.EnsureDkitDataDir(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to create .dkit directory: %w", err)
	}

	processDir := filepath.Join(dataDir, "processes", processID)
	if err := os.MkdirAll(processDir, 0755); err != nil {
		return fmt.Errorf("failed to create process directory: %w", err)
	}

	stdoutPath := filepath.Join(processDir, "stdout.log")
	stderrPath := filepath.Join(processDir, "stderr.log")

	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return fmt.Errorf("failed to create stdout log: %w", err)
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return fmt.Errorf("failed to create stderr log: %w", err)
	}
	defer stderrFile.Close()

	// Build command
	var cmdExec *exec.Cmd
	if len(args) == 1 {
		// Single argument - run through shell
		cmdExec = exec.Command("sh", "-c", args[0])
	} else {
		// Multiple arguments - run directly
		cmdExec = exec.Command(args[0], args[1:]...)
	}

	cmdExec.Dir = workDir

	// Setup environment
	cmdExec.Env = os.Environ()
	if !ignoreLocalBin {
		binDir := filepath.Join(projectRoot, "bin")
		if _, err := os.Stat(binDir); err == nil {
			// Add bin directory to PATH
			pathEnv := os.Getenv("PATH")
			newPath := binDir + string(os.PathListSeparator) + pathEnv
			cmdExec.Env = updateEnv(cmdExec.Env, "PATH", newPath)
		}
	}

	// Setup stdin/stdout/stderr with TTY support
	// Use MultiWriter to write to both terminal and log files
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = io.MultiWriter(os.Stdout, stdoutFile)
	cmdExec.Stderr = io.MultiWriter(os.Stderr, stderrFile)

	// Create metadata
	startTime := time.Now()
	meta := utils.ProcessMetadata{
		ID:         processID,
		Command:    strings.Join(args, " "),
		Args:       args,
		Cwd:        workDir,
		StartedAt:  startTime,
		Status:     utils.StatusRunning,
		StdoutPath: fmt.Sprintf(".dkit/processes/%s/stdout.log", processID),
		StderrPath: fmt.Sprintf(".dkit/processes/%s/stderr.log", processID),
	}

	// Start the command
	if err := cmdExec.Start(); err != nil {
		// Update metadata with failure
		meta.Status = utils.StatusFailed
		exitCode := 1
		meta.ExitCode = &exitCode
		endTime := time.Now()
		meta.EndedAt = &endTime
		utils.SaveProcessMetadata(projectRoot, meta)

		if exitError, ok := err.(*exec.Error); ok && exitError.Err == exec.ErrNotFound {
			return fmt.Errorf("command not found: %s", args[0])
		}
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Update metadata with PID
	meta.PID = cmdExec.Process.Pid
	if err := utils.SaveProcessMetadata(projectRoot, meta); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save process metadata: %v\n", err)
	}

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			// Forward signal to child process
			if cmdExec.Process != nil {
				cmdExec.Process.Signal(os.Interrupt)
			}
		case <-ctx.Done():
		}
	}()

	// Wait for command to finish
	// Output is automatically written to both terminal and log files via MultiWriter
	cmdErr := cmdExec.Wait()

	// Update metadata with completion status
	endTime := time.Now()
	meta.EndedAt = &endTime

	if cmdErr != nil {
		meta.Status = utils.StatusFailed
		if exitError, ok := cmdErr.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			meta.ExitCode = &exitCode
		} else {
			exitCode := 1
			meta.ExitCode = &exitCode
		}
	} else {
		meta.Status = utils.StatusCompleted
		exitCode := 0
		meta.ExitCode = &exitCode
	}

	if err := utils.SaveProcessMetadata(projectRoot, meta); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save process metadata: %v\n", err)
	}

	// Return the command error to preserve exit code
	return cmdErr
}

// updateEnv updates or adds an environment variable in the env slice
func updateEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, e := range env {
		if strings.HasPrefix(e, prefix) {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}
