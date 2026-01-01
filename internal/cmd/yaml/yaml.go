package yaml

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yaml",
		Short: "Normalize YAML files",
		Long: `Normalize YAML files by resolving all anchors, aliases, and merge keys 
to produce a flat, self-contained YAML output.

Designed for use in pipes and automation workflows.`,
	}

	// Add subcommands
	cmd.AddCommand(newNormalizeCommand())

	return cmd
}

func newNormalizeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "normalize [file]",
		Short: "Normalize YAML files",
		Long: `Normalize YAML files by resolving all anchors, aliases, and merge keys.

If no file is specified, reads from stdin.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement yaml normalize
			return nil
		},
	}

	return cmd
}
