package processor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

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
		workerCount, _ := cmd.Flags().GetInt("worker-count")
		return readFile(args[0], wordCount, characterCount, pattern, workerCount)
	},
}

func init() {
	ReadCmd.Flags().BoolP("word-count", "c", false, "Count words in the file")
	ReadCmd.Flags().BoolP("character-count", "n", false, "Count characters in the file")
	ReadCmd.Flags().StringP("search", "s", "", "Search for a pattern in the file")
	ReadCmd.Flags().IntP("worker-count", "w", 4, "Number of concurrent workers") // Add worker flag
}

type FileChunk struct {
	Lines        []string
	WordCount    int
	CharCount    int
	MatchedLines []string // For pattern matching
}

func readFile(filename string, wordCount, characterCount bool, pattern string, workerCount int) error {
	fmt.Println("Starting file processing...") // Debug

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fmt.Printf("File opened successfully: %s\n", filename) // Debug

	chunks := make(chan FileChunk)
	results := make(chan FileChunk)

	var wg sync.WaitGroup
	fmt.Printf("Starting %d workers...\n", workerCount) // Debug

	// Workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		workerID := i
		go func() {
			defer wg.Done()
			fmt.Printf("Worker %d started\n", workerID) // Debug
			for chunk := range chunks {
				fmt.Printf("Worker %d processing chunk of size %d\n", workerID, len(chunk.Lines))
				processedChunk := processChunk(chunk, wordCount, characterCount, pattern)
				fmt.Printf("Worker %d processed chunk: words=%d, chars=%d\n",
					workerID, processedChunk.WordCount, processedChunk.CharCount)
				results <- processedChunk
			}
		}()
	}

	// File reader
	go func() {
		scanner := bufio.NewScanner(file)
		currentChunk := FileChunk{}
		lineCount := 0

		fmt.Println("Starting to read file...") // Debug
		for scanner.Scan() {
			currentChunk.Lines = append(currentChunk.Lines, scanner.Text())
			lineCount++

			if lineCount >= 1000 {
				fmt.Printf("Sending chunk with %d lines\n", lineCount) // Debug
				chunks <- currentChunk
				currentChunk = FileChunk{}
				lineCount = 0
			}
		}

		if len(currentChunk.Lines) > 0 {
			fmt.Printf("Sending final chunk with %d lines\n", len(currentChunk.Lines)) // Debug
			chunks <- currentChunk
		}

		fmt.Println("Closing chunks channel...") // Debug
		close(chunks)
	}()

	var totalStats struct {
		words int
		chars int
		lines int
		sync.Mutex
	}

	var resultsWg sync.WaitGroup
	resultsWg.Add(1)

	// Results collector
	go func() {
		defer resultsWg.Done()
		fmt.Println("Starting to collect results...") // Debug
		for result := range results {
			totalStats.Lock()
			totalStats.words += result.WordCount
			totalStats.chars += result.CharCount
			totalStats.lines += len(result.Lines)
			fmt.Printf("Updated totals: words=%d, chars=%d, lines=%d\n", // Debug
				totalStats.words, totalStats.chars, totalStats.lines)
			totalStats.Unlock()
		}
		fmt.Println("Finished collecting results") // Debug
	}()

	fmt.Println("Waiting for workers to finish...") // Debug
	wg.Wait()
	fmt.Println("Workers finished, closing results channel...") // Debug
	close(results)

	fmt.Println("Waiting for results collection...") // Debug
	resultsWg.Wait()
	fmt.Println("Results collection complete") // Debug

	fmt.Printf("\nFinal Statistics:\n")
	fmt.Printf("Total Words: %d\n", totalStats.words)
	fmt.Printf("Total Characters: %d\n", totalStats.chars)
	fmt.Printf("Total Lines: %d\n", totalStats.lines)

	return nil
}

func processChunk(chunk FileChunk, wordCount, characterCount bool, pattern string) FileChunk {
	results := FileChunk{
		Lines: chunk.Lines,
	}
	for _, line := range chunk.Lines {
		if wordCount {
			results.WordCount += len(strings.Fields(line))
		}
		if characterCount {
			results.CharCount += len(line)
		}
		if pattern != "" {
			if strings.Contains(line, pattern) {
				results.MatchedLines = append(chunk.MatchedLines, line)
			}
		}
	}
	return results
}
