package processor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

// FileProcessor handles the file processing logic
type FileProcessor struct {
	filename        string
	wordCount       bool
	characterCount  bool
	pattern         string
	workerCount     int
	processedChunks int32
	totalChunks     int32
	startTime       time.Time
}

// NewFileProcessor creates a new file processor
func NewFileProcessor(filename string, wordCount, characterCount bool, pattern string, workerCount int) *FileProcessor {
	return &FileProcessor{
		filename:       filename,
		wordCount:      wordCount,
		characterCount: characterCount,
		pattern:        pattern,
		workerCount:    workerCount,
		startTime:      time.Now(),
	}
}

func (fp *FileProcessor) initializeFile() (*os.File, error) {
	file, err := os.Open(fp.filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fp.totalChunks = int32(fileInfo.Size()/(100*1000)) + 1

	return file, nil
}

func (fp *FileProcessor) startProgressBar() {
	go func() {
		for {
			processed := atomic.LoadInt32(&fp.processedChunks)
			if processed >= fp.totalChunks {
				break
			}
			fp.updateProgress(processed)
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func (fp *FileProcessor) updateProgress(processed int32) {
	percentage := float64(processed) / float64(fp.totalChunks) * 100
	elapsed := time.Since(fp.startTime)
	width := 40
	completed := int(float64(width) * float64(processed) / float64(fp.totalChunks))
	bar := strings.Repeat("█", completed) + strings.Repeat("░", width-completed)

	fmt.Printf("\r%s %.1f%% (%d/%d chunks) %s",
		bar, percentage, processed, fp.totalChunks,
		elapsed.Round(time.Second))
}

func (fp *FileProcessor) startWorkers(chunks, results chan FileChunk) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < fp.workerCount; i++ {
		wg.Add(1)
		workerID := i
		go fp.worker(workerID, &wg, chunks, results)
	}
	return &wg
}

func (fp *FileProcessor) worker(id int, wg *sync.WaitGroup, chunks, results chan FileChunk) {
	defer wg.Done()
	for chunk := range chunks {
		processedChunk := processChunk(chunk, fp.wordCount, fp.characterCount, fp.pattern)
		results <- processedChunk
		atomic.AddInt32(&fp.processedChunks, 1)
	}
}

func (fp *FileProcessor) readFileInChunks(file *os.File, chunks chan FileChunk) {
	scanner := bufio.NewScanner(file)
	currentChunk := FileChunk{}
	lineCount := 0

	for scanner.Scan() {
		currentChunk.Lines = append(currentChunk.Lines, scanner.Text())
		lineCount++

		if lineCount >= 1000 {
			chunks <- currentChunk
			currentChunk = FileChunk{}
			lineCount = 0
		}
	}

	if len(currentChunk.Lines) > 0 {
		chunks <- currentChunk
	}
	close(chunks)
}

func (fp *FileProcessor) collectResults(results chan FileChunk) *sync.WaitGroup {
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

	return &resultsWg
}

// Process is the main entry point
func (fp *FileProcessor) Process() error {
	file, err := fp.initializeFile()
	if err != nil {
		return err
	}
	defer file.Close()

	chunks := make(chan FileChunk)
	results := make(chan FileChunk)

	fp.startProgressBar()
	workersWg := fp.startWorkers(chunks, results)

	go fp.readFileInChunks(file, chunks)
	resultsWg := fp.collectResults(results)

	workersWg.Wait()
	close(results)

	fmt.Println("Waiting for results collection...") // Debug
	resultsWg.Wait()
	fmt.Println("Results collection complete") // Debug

	fmt.Printf("\r\033[K") // Clear progress bar
	fmt.Println("\nProcessing complete!")
	return nil
}

// Main function becomes much simpler
func readFile(filename string, wordCount, characterCount bool, pattern string, workerCount int) error {
	processor := NewFileProcessor(filename, wordCount, characterCount, pattern, workerCount)
	return processor.Process()
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
