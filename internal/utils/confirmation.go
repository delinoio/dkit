package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm asks the user for confirmation and returns true if they answer yes
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stdout, "[dkit] %s [y/N]: ", prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// ConfirmOrExit asks for confirmation and exits if the user declines
func ConfirmOrExit(prompt string) {
	if !Confirm(prompt) {
		fmt.Fprintf(os.Stdout, "[dkit] Operation cancelled\n")
		os.Exit(2)
	}
}
