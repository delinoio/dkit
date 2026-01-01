package port

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "port",
		Short: "Manage network ports",
		Long: `Manage network ports during development.

Check port availability, identify processes using ports, and terminate 
port-blocking processes with safety features.`,
	}

	// Add subcommands
	cmd.AddCommand(newCheckCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newKillCommand())
	cmd.AddCommand(newFindCommand())
	cmd.AddCommand(newWatchCommand())

	return cmd
}

func newCheckCommand() *cobra.Command {
	var (
		jsonOutput bool
		quiet      bool
	)

	cmd := &cobra.Command{
		Use:   "check <port>",
		Short: "Check port availability",
		Long:  `Determine if a port is available or in use, with details about the process using it.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheckCommand(args[0], jsonOutput, quiet)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output result in JSON format")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Exit code only")

	return cmd
}

func newListCommand() *cobra.Command {
	var (
		portRange string
		listening bool
		all       bool
		tcp       bool
		udp       bool
		format    string
		sortBy    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List used ports",
		Long:  `Display all currently used ports with process information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(portRange, listening, all, tcp, udp, format, sortBy)
		},
	}

	cmd.Flags().StringVar(&portRange, "range", "", "Only show ports in range (e.g., 3000-4000)")
	cmd.Flags().BoolVar(&listening, "listening", true, "Only show listening ports")
	cmd.Flags().BoolVar(&all, "all", false, "Show all network connections")
	cmd.Flags().BoolVar(&tcp, "tcp", false, "TCP ports only")
	cmd.Flags().BoolVar(&udp, "udp", false, "UDP ports only")
	cmd.Flags().StringVar(&format, "format", "table", "Output format (table|json|csv)")
	cmd.Flags().StringVar(&sortBy, "sort", "port", "Sort by column (port|pid|process)")

	return cmd
}

func newKillCommand() *cobra.Command {
	var (
		force   bool
		signal  string
		timeout int
	)

	cmd := &cobra.Command{
		Use:   "kill <port>",
		Short: "Terminate process on port",
		Long:  `Terminate the process using a specific port. Includes safety confirmations for system processes.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runKillCommand(args[0], force, signal, timeout)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	cmd.Flags().StringVar(&signal, "signal", "TERM", "Signal to send (TERM|KILL|HUP)")
	cmd.Flags().IntVar(&timeout, "timeout", 5, "Grace period before SIGKILL")

	return cmd
}

func newFindCommand() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "find <pattern>",
		Short: "Find processes by port pattern",
		Long:  `Search for processes using ports matching a pattern or range.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFindCommand(args[0], format)
		},
	}

	cmd.Flags().StringVar(&format, "format", "table", "Output format (table|json|list)")

	return cmd
}

func newWatchCommand() *cobra.Command {
	var (
		portRange string
		interval  int
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Monitor port activity",
		Long:  `Watch for processes binding to or releasing ports in real-time.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatchCommand(portRange, interval)
		},
	}

	cmd.Flags().StringVar(&portRange, "range", "", "Only watch specific port range")
	cmd.Flags().IntVar(&interval, "interval", 1, "Polling interval in seconds")

	return cmd
}
