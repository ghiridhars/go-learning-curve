package cmd

import (
	"cli-processor/internal/processor"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fileprocessor",
	Short: "File processing CLI",
	Long:  `A CLI tool for processing and analyzing text files with various options`,
}

// Execute runs the root command and handles any errors that occur during execution.
// It will:
// 1. Execute the root command and all its subcommands
// 2. If an error occurs, print it to stderr instead of stdout for better error handling
// 3. Exit with status code 1 on error
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err) // Write errors to stderr
		os.Exit(1)
	}
}

func init() {
	// This init function adds a persistent flag to the root command
	// Currently it only adds a simple boolean toggle flag that isn't used
	// To optimize this:
	// 1. Remove unused toggle flag
	// 2. Add meaningful flags that match the CLI's purpose like:
	// rootCmd.PersistentFlags().StringP("input", "i", "", "Input file or directory path")
	// rootCmd.PersistentFlags().StringP("output", "o", "", "Output file path for results")
	// rootCmd.PersistentFlags().BoolP("recursive", "r", false, "Process files recursively in directories")
	// rootCmd.PersistentFlags().IntP("workers", "w", 1, "Number of concurrent workers")
	rootCmd.AddCommand(processor.ReadCmd)
}
