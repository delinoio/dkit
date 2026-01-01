package root

import (
	"fmt"
	"os"

	"github.com/delinoio/dkit/internal/cmd/clipboard"
	"github.com/delinoio/dkit/internal/cmd/cron"
	"github.com/delinoio/dkit/internal/cmd/env"
	"github.com/delinoio/dkit/internal/cmd/git"
	"github.com/delinoio/dkit/internal/cmd/jsonc"
	"github.com/delinoio/dkit/internal/cmd/mcp"
	"github.com/delinoio/dkit/internal/cmd/port"
	"github.com/delinoio/dkit/internal/cmd/retry"
	"github.com/delinoio/dkit/internal/cmd/run"
	"github.com/delinoio/dkit/internal/cmd/yaml"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dkit",
	Short: "DevTools by Delino",
	Long: `dkit - DevTools by Delino

A collection of developer tools for terminal workflows, automation, 
and AI-assisted development.`,
	Version: "0.1.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(clipboard.NewCommand())
	rootCmd.AddCommand(cron.NewCommand())
	rootCmd.AddCommand(env.NewCommand())
	rootCmd.AddCommand(git.NewCommand())
	rootCmd.AddCommand(jsonc.NewCommand())
	rootCmd.AddCommand(mcp.NewCommand())
	rootCmd.AddCommand(port.NewCommand())
	rootCmd.AddCommand(retry.NewCommand())
	rootCmd.AddCommand(run.NewCommand())
	rootCmd.AddCommand(yaml.NewCommand())

	// Disable completion command for cleaner output
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Set custom usage template
	rootCmd.SetUsageTemplate(usageTemplate)
}

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[dkit] ERROR: %v\n", err)
		os.Exit(1)
	}
}
