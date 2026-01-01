package clipboard

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/delinoio/dkit/internal/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clipboard",
		Short: "Bridge between terminal and system clipboard",
		Long: `Manage clipboard content from the terminal.

Pipe data to/from clipboard, manage clipboard history, and transform clipboard content.`,
	}

	// Add subcommands
	cmd.AddCommand(newCopyCommand())
	cmd.AddCommand(newPasteCommand())
	cmd.AddCommand(newClearCommand())
	cmd.AddCommand(newWatchCommand())

	return cmd
}

func newCopyCommand() *cobra.Command {
	var (
		trim       bool
		appendFlag bool
		format     string
		quiet      bool
	)

	cmd := &cobra.Command{
		Use:   "copy [file]",
		Short: "Copy to clipboard",
		Long:  `Copy stdin or file content to system clipboard.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var content []byte
			var err error

			// Read content from file or stdin
			if len(args) > 0 {
				content, err = os.ReadFile(args[0])
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
			} else {
				content, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read stdin: %w", err)
				}
			}

			// Process content
			text := string(content)
			if trim {
				text = strings.TrimSpace(text)
			}

			// Handle append mode
			if appendFlag {
				existing, err := getClipboard()
				if err == nil && existing != "" {
					text = existing + text
				}
			}

			// Copy to clipboard
			if err := setClipboard(text); err != nil {
				return fmt.Errorf("failed to set clipboard: %w", err)
			}

			// Output
			if !quiet {
				size := len(text)
				msg := fmt.Sprintf("Copied to clipboard (%d bytes", size)
				if trim {
					msg += ", trimmed"
				}
				msg += ")"
				if appendFlag {
					msg = fmt.Sprintf("Appended to clipboard (%d bytes total)", size)
				}
				utils.PrintSuccess(msg)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&trim, "trim", false, "Remove leading/trailing whitespace")
	cmd.Flags().BoolVar(&appendFlag, "append", false, "Append to existing clipboard content")
	cmd.Flags().StringVar(&format, "format", "text", "Content format (text|html|rtf|image)")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "No output, only exit code")

	return cmd
}

func newPasteCommand() *cobra.Command {
	var (
		output     string
		appendFlag bool
		format     string
	)

	cmd := &cobra.Command{
		Use:   "paste",
		Short: "Paste from clipboard",
		Long:  `Output clipboard content to stdout or file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get clipboard content
			content, err := getClipboard()
			if err != nil {
				return fmt.Errorf("failed to get clipboard: %w", err)
			}

			if content == "" {
				return fmt.Errorf("clipboard is empty")
			}

			// Format output
			var outputContent string
			switch format {
			case "text":
				outputContent = content
			case "json":
				data := map[string]interface{}{
					"content": content,
					"length":  len(content),
					"lines":   len(strings.Split(content, "\n")),
					"format":  "text/plain",
				}
				jsonBytes, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to format as JSON: %w", err)
				}
				outputContent = string(jsonBytes)
			case "base64":
				outputContent = base64.StdEncoding.EncodeToString([]byte(content))
			default:
				return fmt.Errorf("invalid format: %s", format)
			}

			// Output to file or stdout
			if output != "" {
				flags := os.O_CREATE | os.O_WRONLY
				if appendFlag {
					flags |= os.O_APPEND
				} else {
					flags |= os.O_TRUNC
				}

				f, err := os.OpenFile(output, flags, 0644)
				if err != nil {
					return fmt.Errorf("failed to open file: %w", err)
				}
				defer f.Close()

				if _, err := f.WriteString(outputContent); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}

				utils.PrintSuccess(fmt.Sprintf("Clipboard content saved to %s (%d bytes)", output, len(outputContent)))
			} else {
				fmt.Print(outputContent)
				if format == "text" && !strings.HasSuffix(outputContent, "\n") {
					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Write to file instead of stdout")
	cmd.Flags().BoolVar(&appendFlag, "append", false, "Append to file (if --output specified)")
	cmd.Flags().StringVar(&format, "format", "text", "Output format (text|json|base64)")

	return cmd
}

func newClearCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear the clipboard",
		Long:  `Clear the system clipboard.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				if !utils.Confirm("Clear clipboard?") {
					return nil
				}
			}

			if err := setClipboard(""); err != nil {
				return fmt.Errorf("failed to clear clipboard: %w", err)
			}

			utils.PrintSuccess("Clipboard cleared")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")

	return cmd
}

func newWatchCommand() *cobra.Command {
	var (
		command  string
		interval int
		logFile  string
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch clipboard changes",
		Long:  `Monitor clipboard for changes and trigger actions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.PrintInfo("Watching clipboard for changes (Ctrl+C to stop)")

			var lastContent string
			ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
			defer ticker.Stop()

			// Open log file if specified
			var logWriter io.Writer
			if logFile != "" {
				f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					return fmt.Errorf("failed to open log file: %w", err)
				}
				defer f.Close()
				logWriter = f
				utils.PrintInfo(fmt.Sprintf("Logging to %s", logFile))
			}

			for range ticker.C {
				content, err := getClipboard()
				if err != nil {
					continue
				}

				if content != lastContent && content != "" {
					timestamp := time.Now().Format("15:04:05")
					fmt.Printf("\n[%s] Clipboard changed (%d bytes)\n", timestamp, len(content))

					// Show preview (first 100 chars)
					preview := content
					if len(preview) > 100 {
						preview = preview[:100] + "..."
					}
					fmt.Printf("%s\n", preview)

					// Log to file
					if logWriter != nil {
						logEntry := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), content)
						logWriter.Write([]byte(logEntry))
					}

					// Run command if specified
					if command != "" {
						execCmd := exec.Command("sh", "-c", command)
						execCmd.Env = append(os.Environ(), fmt.Sprintf("CLIPBOARD_CONTENT=%s", content))
						if err := execCmd.Run(); err != nil {
							fmt.Printf("[dkit] Command failed: %v\n", err)
						}
					}

					lastContent = content
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&command, "command", "", "Command to run on change")
	cmd.Flags().IntVar(&interval, "interval", 500, "Polling interval in milliseconds")
	cmd.Flags().StringVar(&logFile, "log", "", "Log clipboard changes to file")

	return cmd
}

// getClipboard reads from system clipboard
func getClipboard() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--output")
		} else {
			return "", fmt.Errorf("no clipboard tool found (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("powershell.exe", "-command", "Get-Clipboard")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// setClipboard writes to system clipboard
func setClipboard(content string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard tool found (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("powershell.exe", "-command", "$input | Set-Clipboard")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := stdin.Write([]byte(content)); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}
