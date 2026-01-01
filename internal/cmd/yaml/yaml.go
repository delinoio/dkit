package yaml

import (
	"io"
	"os"

	"github.com/delinoio/dkit/internal/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
			var input io.Reader

			if len(args) > 0 {
				// Read from file
				file, err := os.Open(args[0])
				if err != nil {
					utils.PrintError("File not found: %s", args[0])
					os.Exit(2)
				}
				defer file.Close()
				input = file
			} else {
				// Read from stdin
				input = os.Stdin
			}

			// Read and parse YAML
			decoder := yaml.NewDecoder(input)
			encoder := yaml.NewEncoder(os.Stdout)
			encoder.SetIndent(2)
			defer encoder.Close()

			// Process all documents
			for {
				var doc interface{}
				err := decoder.Decode(&doc)
				if err == io.EOF {
					break
				}
				if err != nil {
					utils.PrintError("Invalid YAML syntax: %v", err)
					os.Exit(1)
				}

				// Encode back (this resolves all anchors and aliases)
				if err := encoder.Encode(doc); err != nil {
					utils.PrintError("Failed to encode YAML: %v", err)
					os.Exit(3)
				}
			}

			return nil
		},
	}

	return cmd
}
