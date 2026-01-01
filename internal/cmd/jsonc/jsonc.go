package jsonc

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jsonc",
		Short: "Convert JSONC/JSON5 to JSON",
		Long: `Convert JSONC (JSON with Comments) or JSON5 files to standard JSON format.

Designed for use in pipes and automation workflows.`,
	}

	// Add subcommands
	cmd.AddCommand(newCompileCommand())

	return cmd
}

func newCompileCommand() *cobra.Command {
	var pretty bool

	cmd := &cobra.Command{
		Use:   "compile [file]",
		Short: "Convert JSONC/JSON5 to JSON",
		Long: `Convert JSONC (JSON with Comments) or JSON5 files to standard JSON format.

If no file is specified, reads from stdin.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement jsonc compile
			return nil
		},
	}

	cmd.Flags().BoolVar(&pretty, "pretty", false, "Pretty-print JSON output")

	return cmd
}
