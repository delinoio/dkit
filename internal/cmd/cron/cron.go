package cron

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/delinoio/dkit/internal/utils"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cron",
		Short: "Parse and validate cron expressions",
		Long: `Parse, validate, and explain cron expressions.

Generate cron schedules and calculate execution times. Make cron scheduling 
more accessible and less error-prone.`,
	}

	// Add subcommands
	cmd.AddCommand(newParseCommand())
	cmd.AddCommand(newNextCommand())
	cmd.AddCommand(newValidateCommand())
	cmd.AddCommand(newGenerateCommand())

	return cmd
}

func newParseCommand() *cobra.Command {
	var (
		format  string
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "parse <expression>",
		Short: "Parse cron expression",
		Long:  `Convert a cron expression into human-readable description.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			expression := args[0]

			// Parse the cron expression
			schedule, err := parser.Parse(expression)
			if err != nil {
				utils.PrintError("Invalid cron expression: %v", err)
				return err
			}

			description := describeCronExpression(expression)

			if format == "json" {
				return outputParseJSON(expression, description, schedule, verbose)
			}

			return outputParseText(expression, description, schedule, verbose)
		},
	}

	cmd.Flags().StringVar(&format, "format", "text", "Output format (text|json)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Include detailed field breakdown")

	return cmd
}

func newNextCommand() *cobra.Command {
	var (
		count    int
		from     string
		timezone string
		format   string
	)

	cmd := &cobra.Command{
		Use:   "next <expression>",
		Short: "Calculate next execution times",
		Long:  `Calculate when a cron job will run next.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			expression := args[0]

			// Parse the cron expression
			schedule, err := parser.Parse(expression)
			if err != nil {
				utils.PrintError("Invalid cron expression: %v", err)
				return err
			}

			// Determine starting time
			var startTime time.Time
			if from != "" {
				startTime, err = time.Parse(time.RFC3339, from)
				if err != nil {
					// Try alternative formats
					startTime, err = time.Parse("2006-01-02 15:04:05", from)
					if err != nil {
						utils.PrintError("Invalid time format: %s", from)
						return err
					}
				}
			} else {
				startTime = time.Now()
			}

			// Handle timezone
			if timezone != "" {
				loc, err := time.LoadLocation(timezone)
				if err != nil {
					utils.PrintError("Invalid timezone: %s", timezone)
					return err
				}
				startTime = startTime.In(loc)
			}

			// Calculate next executions
			executions := make([]time.Time, count)
			currentTime := startTime
			for i := 0; i < count; i++ {
				next := schedule.Next(currentTime)
				executions[i] = next
				currentTime = next
			}

			if format == "json" {
				return outputNextJSON(expression, executions)
			} else if format == "csv" {
				return outputNextCSV(executions)
			}

			return outputNextText(expression, executions, startTime)
		},
	}

	cmd.Flags().IntVarP(&count, "count", "n", 5, "Number of future executions to show")
	cmd.Flags().StringVar(&from, "from", "", "Calculate from specific time")
	cmd.Flags().StringVar(&timezone, "timezone", "", "Timezone for calculations")
	cmd.Flags().StringVar(&format, "format", "text", "Output format (text|json|csv)")

	return cmd
}

func newValidateCommand() *cobra.Command {
	var (
		strict bool
		system string
	)

	cmd := &cobra.Command{
		Use:   "validate <expression>",
		Short: "Validate cron expression",
		Long:  `Check if a cron expression is valid and identify issues.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			expression := args[0]

			// Try to parse the expression
			schedule, err := parser.Parse(expression)
			if err != nil {
				utils.PrintError("Invalid cron expression: %s", expression)
				fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
				fmt.Fprintf(os.Stderr, "\nFormat: minute hour day month weekday\n")
				fmt.Fprintf(os.Stderr, "Example: 0 9 * * 1-5 (Every weekday at 9 AM)\n")
				os.Exit(1)
			}

			// Valid expression
			utils.PrintSuccess("Valid cron expression: %s", expression)
			utils.PrintInfo("Type: Standard cron (5 fields)")
			description := describeCronExpression(expression)
			utils.PrintInfo("Description: %s", description)

			// Show next execution
			next := schedule.Next(time.Now())
			utils.PrintInfo("Next execution: %s", next.Format("2006-01-02 15:04:05"))

			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation")
	cmd.Flags().StringVar(&system, "system", "", "Validate for specific cron implementation (linux|macos|freebsd)")

	return cmd
}

func newGenerateCommand() *cobra.Command {
	var (
		every       string
		at          string
		on          string
		interactive bool
		describe    string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate cron expression",
		Long:  `Create cron expressions using natural language or interactive prompts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var expression string

			if interactive {
				utils.PrintError("Interactive mode not yet implemented")
				return fmt.Errorf("interactive mode not yet implemented")
			}

			if every != "" {
				var err error
				expression, err = generateFromInterval(every, at, on)
				if err != nil {
					utils.PrintError("Failed to generate expression: %v", err)
					return err
				}
			} else if at != "" {
				var err error
				expression, err = generateFromTime(at, on)
				if err != nil {
					utils.PrintError("Failed to generate expression: %v", err)
					return err
				}
			} else {
				utils.PrintError("Please specify --every or --at flag")
				return fmt.Errorf("no generation parameters specified")
			}

			// Validate the generated expression
			_, err := parser.Parse(expression)
			if err != nil {
				utils.PrintError("Generated invalid expression: %v", err)
				return err
			}

			description := describeCronExpression(expression)
			utils.PrintInfo("Generated expression: %s", expression)
			utils.PrintInfo("Description: %s", description)

			// Print just the expression to stdout for piping
			fmt.Println(expression)
			return nil
		},
	}

	cmd.Flags().StringVar(&every, "every", "", "Simple interval (1h, 30m, 1d, etc.)")
	cmd.Flags().StringVar(&at, "at", "", "Specific time (09:00, 2:30pm)")
	cmd.Flags().StringVar(&on, "on", "", "Specific days (mon,wed,fri or 1,15)")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive prompt mode")
	cmd.Flags().StringVar(&describe, "describe", "", "Natural language description")

	return cmd
}

// Helper functions

func describeCronExpression(expression string) string {
	parts := strings.Fields(expression)
	if len(parts) != 5 {
		return "Unknown format"
	}

	minute, hour, day, month, weekday := parts[0], parts[1], parts[2], parts[3], parts[4]

	var desc []string

	// Handle special cases first
	if minute == "0" && hour == "0" && day == "*" && month == "*" && weekday == "*" {
		return "Daily at midnight"
	}
	if minute == "0" && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return "Every hour"
	}

	// Describe time
	if strings.Contains(hour, "*/") {
		interval := strings.TrimPrefix(hour, "*/")
		desc = append(desc, fmt.Sprintf("Every %s hours", interval))
	} else if hour == "*" {
		desc = append(desc, "Every hour")
	} else {
		desc = append(desc, fmt.Sprintf("At %s:%s", hour, minute))
	}

	// Describe day/weekday
	if weekday != "*" {
		desc = append(desc, describeWeekday(weekday))
	} else if day != "*" {
		desc = append(desc, fmt.Sprintf("on day %s", day))
	}

	// Describe month
	if month != "*" {
		desc = append(desc, fmt.Sprintf("in month %s", month))
	}

	return strings.Join(desc, ", ")
}

func describeWeekday(weekday string) string {
	days := map[string]string{
		"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
		"4": "Thursday", "5": "Friday", "6": "Saturday", "7": "Sunday",
	}

	if strings.Contains(weekday, "-") {
		parts := strings.Split(weekday, "-")
		if len(parts) == 2 {
			return fmt.Sprintf("%s through %s", days[parts[0]], days[parts[1]])
		}
	}

	if day, ok := days[weekday]; ok {
		return "on " + day
	}

	return "on weekday " + weekday
}

func generateFromInterval(interval, at, on string) (string, error) {
	// Parse interval like "1h", "30m", "1d"
	if len(interval) < 2 {
		return "", fmt.Errorf("invalid interval format")
	}

	unit := interval[len(interval)-1:]
	valueStr := interval[:len(interval)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return "", fmt.Errorf("invalid interval value: %s", valueStr)
	}

	switch unit {
	case "m": // minutes
		if value >= 60 {
			return "", fmt.Errorf("minute interval must be less than 60")
		}
		return fmt.Sprintf("*/%d * * * *", value), nil
	case "h": // hours
		if value >= 24 {
			return "", fmt.Errorf("hour interval must be less than 24")
		}
		return fmt.Sprintf("0 */%d * * *", value), nil
	case "d": // days
		return "0 0 * * *", nil
	default:
		return "", fmt.Errorf("unknown interval unit: %s (use m, h, or d)", unit)
	}
}

func generateFromTime(timeStr, on string) (string, error) {
	// Parse time like "09:00" or "2:30pm"
	var hour, minute int

	// Handle 12-hour format
	isPM := strings.HasSuffix(strings.ToLower(timeStr), "pm")
	isAM := strings.HasSuffix(strings.ToLower(timeStr), "am")

	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "pm"), "PM")
	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "am"), "AM")
	timeStr = strings.TrimSpace(timeStr)

	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid time format (use HH:MM)")
	}

	var err error
	hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid hour: %s", parts[0])
	}
	minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid minute: %s", parts[1])
	}

	// Convert 12-hour to 24-hour
	if isPM && hour != 12 {
		hour += 12
	} else if isAM && hour == 12 {
		hour = 0
	}

	if hour < 0 || hour > 23 {
		return "", fmt.Errorf("hour must be 0-23")
	}
	if minute < 0 || minute > 59 {
		return "", fmt.Errorf("minute must be 0-59")
	}

	weekday := "*"
	day := "*"

	if on != "" {
		// Handle weekday names or ranges
		if strings.Contains(on, "mon") || strings.Contains(on, "tue") || strings.Contains(on, "wed") {
			weekday = parseWeekdayNames(on)
		} else if strings.Contains(on, "-") {
			// Range like "1-5"
			weekday = on
		} else if strings.Contains(on, ",") {
			// List like "1,3,5"
			parts := strings.Split(on, ",")
			// Check if these are day numbers or weekdays
			if isNumeric(parts[0]) && len(parts[0]) <= 2 {
				num, _ := strconv.Atoi(parts[0])
				if num >= 1 && num <= 31 {
					day = on
				} else if num >= 0 && num <= 7 {
					weekday = on
				}
			}
		} else {
			// Single number
			num, err := strconv.Atoi(on)
			if err == nil {
				if num >= 1 && num <= 31 {
					day = on
				} else if num >= 0 && num <= 7 {
					weekday = on
				}
			}
		}
	}

	return fmt.Sprintf("%d %d %s * %s", minute, hour, day, weekday), nil
}

func parseWeekdayNames(on string) string {
	weekdays := map[string]string{
		"mon": "1", "tue": "2", "wed": "3", "thu": "4",
		"fri": "5", "sat": "6", "sun": "0",
	}

	// Handle "mon-fri"
	if strings.Contains(on, "-") {
		parts := strings.Split(on, "-")
		if len(parts) == 2 {
			start := weekdays[strings.ToLower(parts[0])]
			end := weekdays[strings.ToLower(parts[1])]
			if start != "" && end != "" {
				return start + "-" + end
			}
		}
	}

	// Handle "mon,wed,fri"
	if strings.Contains(on, ",") {
		parts := strings.Split(on, ",")
		var nums []string
		for _, part := range parts {
			if num, ok := weekdays[strings.ToLower(strings.TrimSpace(part))]; ok {
				nums = append(nums, num)
			}
		}
		if len(nums) > 0 {
			return strings.Join(nums, ",")
		}
	}

	// Single day
	if num, ok := weekdays[strings.ToLower(on)]; ok {
		return num
	}

	return "*"
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func outputParseText(expression, description string, schedule cron.Schedule, verbose bool) error {
	utils.PrintInfo("Cron expression: %s", expression)
	fmt.Println()
	fmt.Println("Human readable:")
	fmt.Printf("  %s\n", description)
	fmt.Println()

	if verbose {
		parts := strings.Fields(expression)
		if len(parts) == 5 {
			fmt.Println("Schedule breakdown:")
			fmt.Printf("  Minute: %s\n", parts[0])
			fmt.Printf("  Hour: %s\n", parts[1])
			fmt.Printf("  Day: %s\n", parts[2])
			fmt.Printf("  Month: %s\n", parts[3])
			fmt.Printf("  Weekday: %s\n", parts[4])
			fmt.Println()
		}
	}

	// Show next 5 executions
	fmt.Println("Next 5 executions:")
	current := time.Now()
	for i := 0; i < 5; i++ {
		next := schedule.Next(current)
		fmt.Printf("  %s\n", next.Format("2006-01-02 15:04:05"))
		current = next
	}

	return nil
}

func outputParseJSON(expression, description string, schedule cron.Schedule, verbose bool) error {
	result := map[string]interface{}{
		"expression":  expression,
		"description": description,
	}

	if verbose {
		parts := strings.Fields(expression)
		if len(parts) == 5 {
			result["fields"] = map[string]interface{}{
				"minute":  parts[0],
				"hour":    parts[1],
				"day":     parts[2],
				"month":   parts[3],
				"weekday": parts[4],
			}
		}
	}

	// Add next executions
	current := time.Now()
	executions := make([]string, 5)
	for i := 0; i < 5; i++ {
		next := schedule.Next(current)
		executions[i] = next.Format(time.RFC3339)
		current = next
	}
	result["next_executions"] = executions

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func outputNextText(expression string, executions []time.Time, startTime time.Time) error {
	utils.PrintInfo("Next executions for: %s", expression)
	fmt.Println()

	for i, exec := range executions {
		fmt.Printf("%d. %s\n", i+1, exec.Format("Monday, January 02, 2006 at 03:04:05 PM"))

		duration := exec.Sub(time.Now())
		if duration > 0 {
			fmt.Printf("   (in %s)\n", formatDuration(duration))
		}
		fmt.Println()
	}

	return nil
}

func outputNextJSON(expression string, executions []time.Time) error {
	result := map[string]interface{}{
		"expression": expression,
		"executions": make([]map[string]interface{}, len(executions)),
	}

	for i, exec := range executions {
		duration := exec.Sub(time.Now())
		result["executions"].([]map[string]interface{})[i] = map[string]interface{}{
			"datetime":  exec.Format(time.RFC3339),
			"timestamp": exec.Unix(),
			"relative":  formatDuration(duration),
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func outputNextCSV(executions []time.Time) error {
	fmt.Println("datetime,timestamp,relative")
	for _, exec := range executions {
		duration := exec.Sub(time.Now())
		fmt.Printf("%s,%d,%s\n",
			exec.Format("2006-01-02 15:04:05"),
			exec.Unix(),
			formatDuration(duration))
	}
	return nil
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		return "in the past"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var parts []string
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	if minutes > 0 || len(parts) == 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	return strings.Join(parts, ", ")
}
