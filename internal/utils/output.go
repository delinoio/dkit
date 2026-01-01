package utils

import (
	"fmt"
	"os"
)

// PrintSuccess prints a success message with [dkit] prefix and checkmark
func PrintSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stdout, "[dkit] ✓ %s\n", msg)
}

// PrintError prints an error message with [dkit] prefix to stderr
func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[dkit] ERROR: %s\n", msg)
}

// PrintWarning prints a warning message with [dkit] prefix
func PrintWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[dkit] WARNING: %s\n", msg)
}

// PrintInfo prints an info message with [dkit] prefix
func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stdout, "[dkit] %s\n", msg)
}

// PrintVerbose prints a verbose message if verbose mode is enabled
func PrintVerbose(verbose bool, format string, args ...interface{}) {
	if verbose {
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stdout, "[dkit] %s\n", msg)
	}
}

// Fatal prints an error message and exits with code 1
func Fatal(format string, args ...interface{}) {
	PrintError(format, args...)
	os.Exit(1)
}

// FatalWithCode prints an error message and exits with the specified code
func FatalWithCode(code int, format string, args ...interface{}) {
	PrintError(format, args...)
	os.Exit(code)
}
