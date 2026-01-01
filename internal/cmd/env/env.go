package env

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/delinoio/dkit/internal/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Manage environment variables",
		Long: `Manage environment variables across multiple .env files.

Parse, merge, validate, and convert environment configurations for 
different deployment scenarios.`,
	}

	// Add subcommands
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newMergeCommand())
	cmd.AddCommand(newValidateCommand())
	cmd.AddCommand(newGetCommand())
	cmd.AddCommand(newSetCommand())

	return cmd
}

func newListCommand() *cobra.Command {
	var (
		format      string
		showSources bool
		noExpand    bool
	)

	cmd := &cobra.Command{
		Use:   "list [file...]",
		Short: "Display environment variables",
		Long:  `Parse and display environment variables from .env files in a readable format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			envVars, sources, err := loadEnvFiles(args, !noExpand)
			if err != nil {
				return err
			}

			return outputEnvVars(envVars, sources, format, showSources)
		},
	}

	cmd.Flags().StringVar(&format, "format", "text", "Output format (text|json|yaml|export)")
	cmd.Flags().BoolVar(&showSources, "show-sources", false, "Show which file each variable comes from")
	cmd.Flags().BoolVar(&noExpand, "no-expand", false, "Don't expand variable references")

	return cmd
}

// loadEnvFiles loads one or more .env files and merges them
func loadEnvFiles(files []string, expand bool) (map[string]string, map[string]string, error) {
	if len(files) == 0 {
		files = []string{".env"}
	}

	envVars := make(map[string]string)
	sources := make(map[string]string)

	for _, file := range files {
		vars, err := godotenv.Read(file)
		if err != nil {
			if os.IsNotExist(err) {
				utils.PrintError("File not found: %s", file)
				return nil, nil, err
			}
			utils.PrintError("Failed to parse %s: %v", file, err)
			return nil, nil, err
		}

		for key, value := range vars {
			envVars[key] = value
			sources[key] = file
		}
	}

	// Expand variable references if requested
	if expand {
		for key, value := range envVars {
			envVars[key] = os.Expand(value, func(k string) string {
				if v, ok := envVars[k]; ok {
					return v
				}
				return os.Getenv(k)
			})
		}
	}

	return envVars, sources, nil
}

// getSortedKeys returns sorted keys from a map
func getSortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// outputEnvVars outputs environment variables in the specified format
func outputEnvVars(envVars map[string]string, sources map[string]string, format string, showSources bool) error {
	keys := getSortedKeys(envVars)

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(envVars)

	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		defer encoder.Close()
		return encoder.Encode(envVars)

	case "export":
		for _, key := range keys {
			value := envVars[key]
			// Quote values that contain spaces or special characters
			if strings.ContainsAny(value, " \t\n\"'$\\") {
				value = fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "'\\''"))
			}
			fmt.Printf("export %s=%s\n", key, value)
		}
		return nil

	case "text":
		fallthrough
	default:
		for _, key := range keys {
			value := envVars[key]
			if showSources && sources != nil {
				source := sources[key]
				fmt.Printf("%s=%s  # %s\n", key, value, source)
			} else {
				fmt.Printf("%s=%s\n", key, value)
			}
		}
		return nil
	}
}

func newMergeCommand() *cobra.Command {
	var (
		output           string
		format           string
		commentConflicts bool
	)

	cmd := &cobra.Command{
		Use:   "merge <file1> <file2> [file...]",
		Short: "Combine multiple .env files",
		Long:  `Merge multiple .env files with proper precedence rules. Later files override earlier ones.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			envVars, sources, err := loadEnvFiles(args, true)
			if err != nil {
				return err
			}

			var writer io.Writer = os.Stdout
			if output != "" {
				file, err := os.Create(output)
				if err != nil {
					utils.PrintError("Failed to create output file: %v", err)
					os.Exit(1)
				}
				defer file.Close()
				writer = file
			}

			switch format {
			case "json":
				encoder := json.NewEncoder(writer)
				encoder.SetIndent("", "  ")
				return encoder.Encode(envVars)

			case "yaml":
				encoder := yaml.NewEncoder(writer)
				encoder.SetIndent(2)
				defer encoder.Close()
				return encoder.Encode(envVars)

			case "export":
				keys := getSortedKeys(envVars)
				for _, key := range keys {
					value := envVars[key]
					if strings.ContainsAny(value, " \t\n\"'$\\") {
						value = fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "'\\''"))
					}
					fmt.Fprintf(writer, "export %s=%s\n", key, value)
				}
				return nil

			case "dotenv":
				fallthrough
			default:
				keys := getSortedKeys(envVars)
				for _, key := range keys {
					value := envVars[key]
					if commentConflicts && sources != nil {
						source := sources[key]
						fmt.Fprintf(writer, "# from: %s\n", source)
					}
					fmt.Fprintf(writer, "%s=%s\n", key, value)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Write to file instead of stdout")
	cmd.Flags().StringVar(&format, "format", "dotenv", "Output format (dotenv|json|yaml|export)")
	cmd.Flags().BoolVar(&commentConflicts, "comment-conflicts", false, "Add comments showing overridden values")

	return cmd
}

func newValidateCommand() *cobra.Command {
	var (
		required     string
		requiredFile string
		schema       string
		allowEmpty   bool
		strict       bool
	)

	cmd := &cobra.Command{
		Use:   "validate [file...]",
		Short: "Check environment configuration",
		Long:  `Validate environment files against a schema or required variables list.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			envVars, _, err := loadEnvFiles(args, true)
			if err != nil {
				return err
			}

			hasErrors := false

			// Get required variables list
			var requiredVars []string
			if required != "" {
				requiredVars = strings.Split(required, ",")
				for i := range requiredVars {
					requiredVars[i] = strings.TrimSpace(requiredVars[i])
				}
			}

			if requiredFile != "" {
				data, err := os.ReadFile(requiredFile)
				if err != nil {
					utils.PrintError("Failed to read required-file: %v", err)
					os.Exit(1)
				}
				lines := strings.Split(string(data), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && !strings.HasPrefix(line, "#") {
						requiredVars = append(requiredVars, line)
					}
				}
			}

			// Check required variables
			var missingVars []string
			var emptyVars []string

			for _, varName := range requiredVars {
				value, exists := envVars[varName]
				if !exists {
					missingVars = append(missingVars, varName)
					hasErrors = true
				} else if value == "" && !allowEmpty {
					emptyVars = append(emptyVars, varName)
					if strict {
						hasErrors = true
					}
				}
			}

			// Print results
			if len(missingVars) > 0 {
				utils.PrintError("Missing required variables:")
				for _, varName := range missingVars {
					fmt.Fprintf(os.Stderr, "  - %s\n", varName)
				}
			}

			if len(emptyVars) > 0 {
				if strict {
					utils.PrintError("Empty values for required variables:")
				} else {
					utils.PrintWarning("Empty values for required variables:")
				}
				for _, varName := range emptyVars {
					fmt.Fprintf(os.Stderr, "  - %s\n", varName)
				}
			}

			if hasErrors {
				utils.PrintError("Validation failed")
				os.Exit(2)
			}

			utils.PrintSuccess("Validation passed")
			utils.PrintInfo("Found %d variables", len(envVars))
			if len(requiredVars) > 0 {
				utils.PrintInfo("All %d required variables present", len(requiredVars))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&required, "required", "", "Comma-separated list of required variables")
	cmd.Flags().StringVar(&requiredFile, "required-file", "", "File containing required variables")
	cmd.Flags().StringVar(&schema, "schema", "", "JSON schema file for validation")
	cmd.Flags().BoolVar(&allowEmpty, "allow-empty", false, "Allow empty values for required variables")
	cmd.Flags().BoolVar(&strict, "strict", false, "Fail on warnings")

	return cmd
}

func newGetCommand() *cobra.Command {
	var (
		defaultValue string
		expand       bool
	)

	cmd := &cobra.Command{
		Use:   "get <variable> [file...]",
		Short: "Retrieve single variable",
		Long:  `Get the value of a specific environment variable from .env files.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			varName := args[0]
			files := args[1:]

			envVars, _, err := loadEnvFiles(files, expand)
			if err != nil {
				return err
			}

			value, exists := envVars[varName]
			if !exists {
				if defaultValue != "" {
					fmt.Println(defaultValue)
					return nil
				}
				utils.PrintError("Variable not found: %s", varName)
				os.Exit(1)
			}

			fmt.Println(value)
			return nil
		},
	}

	cmd.Flags().StringVar(&defaultValue, "default", "", "Default value if variable not found")
	cmd.Flags().BoolVar(&expand, "expand", true, "Expand variable references")

	return cmd
}

func newSetCommand() *cobra.Command {
	var (
		file    string
		create  bool
		quote   string
		comment string
	)

	cmd := &cobra.Command{
		Use:   "set <variable> <value>",
		Short: "Update or add variable",
		Long:  `Set or update a variable in a .env file safely.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			varName := args[0]
			varValue := args[1]

			// Check if file exists
			_, err := os.Stat(file)
			if os.IsNotExist(err) {
				if !create {
					utils.PrintError("File not found: %s", file)
					utils.PrintInfo("Use --create flag to create the file")
					os.Exit(1)
				}
				// Create empty file with secure permissions
				f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					utils.PrintError("Failed to create file: %v", err)
					os.Exit(1)
				}
				f.Close()
			} else if err != nil {
				utils.PrintError("Failed to access file: %v", err)
				os.Exit(1)
			}

			// Load existing variables
			envVars, err := godotenv.Read(file)
			if err != nil {
				envVars = make(map[string]string)
			}

			// Update the variable
			envVars[varName] = varValue

			// Write back to file
			err = godotenv.Write(envVars, file)
			if err != nil {
				utils.PrintError("Failed to write to file: %v", err)
				os.Exit(1)
			}

			utils.PrintSuccess("Set %s in %s", varName, file)
			return nil
		},
	}

	cmd.Flags().StringVar(&file, "file", ".env", "Environment file to modify")
	cmd.Flags().BoolVar(&create, "create", false, "Create file if it doesn't exist")
	cmd.Flags().StringVar(&quote, "quote", "auto", "Quote behavior (always|auto|never)")
	cmd.Flags().StringVar(&comment, "comment", "", "Add inline comment")

	return cmd
}
