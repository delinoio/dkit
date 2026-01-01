package run

import (
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
		Long: `Execute shell commands with AI-optimized output and persistent logging.

Supports watching/monitoring processes through the MCP interface.`,
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement run
			return nil
		},
	}

	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "Execute in project root directory")
	cmd.Flags().BoolVar(&ignoreLocalBin, "ignore-local-bin", false, "Skip adding <project-root>/bin to PATH")

	return cmd
}
