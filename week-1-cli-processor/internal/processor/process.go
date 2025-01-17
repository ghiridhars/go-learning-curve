package processor

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var ReadCmd = &cobra.Command{
	Use:   "read [file]",
	Short: "Read/Display a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires a file argument")
		}
		wordCount, _ := cmd.Flags().GetBool("word-count")
		characterCount, _ := cmd.Flags().GetBool("character-count")
		pattern, _ := cmd.Flags().GetString("search")
		return readFile(args[0], wordCount, characterCount, pattern)
	},
}

func init() {
	ReadCmd.Flags().BoolP("word-count", "c", false, "Count words in the file")
	ReadCmd.Flags().BoolP("character-count", "n", false, "Count characters in the file")
	ReadCmd.Flags().StringP("search", "s", "", "Search for a pattern in the file")
}

func readFile(filename string, wordCount bool, characterCount bool, pattern string) error {

	// 1. Open file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// 2. Read file line by line
	scanner := bufio.NewScanner(file)
	lineNum := 1
	var words int
	var chars int

	for scanner.Scan() {
		words += len(strings.Fields(scanner.Text()))
		chars += len(scanner.Text())
		if pattern != "" && strings.Contains(scanner.Text(), pattern) {
			fmt.Printf("Pattern found in line %d: %s\n", lineNum, scanner.Text())
		}
		// fmt.Printf("Line %d: %s\n", lineNum, scanner.Text())
		lineNum++
	}

	if wordCount {
		fmt.Printf("Total words: %d\n", words)
	}

	if characterCount {
		fmt.Printf("Total characters: %d\n", chars)
	}

	if pattern != "" {
		fmt.Printf("Pattern: %s\n", pattern)
	}

	fmt.Printf("File: %s has been processed with a total of %d lines.\n", filename, lineNum)

	// 3. Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}
