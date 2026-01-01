package retry

import (
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
			// TODO: Implement retry
			return nil
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
