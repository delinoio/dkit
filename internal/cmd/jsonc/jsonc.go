package jsonc

import (
	"fmt"
	"io"
	"os"

	"github.com/delinoio/dkit/internal/utils"
	"github.com/spf13/cobra"
	"github.com/tailscale/hujson"
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
			var input io.Reader

			if len(args) > 0 {
				// Read from file
				file, err := os.Open(args[0])
				if err != nil {
					utils.PrintError("File not found: %s", args[0])
					return err
				}
				defer file.Close()
				input = file
			} else {
				// Read from stdin
				input = os.Stdin
			}

			// Read all input
			data, err := io.ReadAll(input)
			if err != nil {
				utils.PrintError("Failed to read input: %v", err)
				return err
			}

			// Parse JSONC
			ast, err := hujson.Parse(data)
			if err != nil {
				utils.PrintError("Invalid JSONC syntax: %v", err)
				os.Exit(1)
			}

			// Standardize to JSON
			ast.Standardize()

			// Output
			output := ast.Pack()
			if pretty {
				// Pretty print
				formatted, err := hujson.Format(output)
				if err == nil {
					output = formatted
				}
			}

			fmt.Fprintf(os.Stdout, "%s\n", output)
			return nil
		},
	}

	cmd.Flags().BoolVar(&pretty, "pretty", false, "Pretty-print JSON output")

	return cmd
}
