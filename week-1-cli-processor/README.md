# CLI File Processor Project Specification

## Project Overview
Build a command-line tool that processes text files with features similar to `wc`, `grep`, and `awk` but with additional functionality and concurrent processing capabilities.

## Core Requirements

### 1. Basic File Operations
```go
// Example command usage:
go run main.go analyze --file=input.txt
go run main.go analyze --dir=/path/to/files --pattern=*.txt
```

Implement:
- File reading with error handling
- Directory traversal
- Support for both single file and batch processing
- File existence and permission checks

### 2. Text Analysis Features
The tool should calculate and report:
- Total number of lines
- Word count
- Character count (with and without whitespace)
- Number of paragraphs
- Most frequent words (top N, configurable)
- Basic statistics (average words per line, etc.)

### 3. Search Capabilities
```go
// Example command usage:
go run main.go search --file=input.txt --pattern="[0-9]+" --context=2
```

Implement:
- Regular expression pattern matching
- Line number reporting
- Context lines (showing N lines before and after matches)
- Case-sensitive/insensitive search options

### 4. Output Formats
Support multiple output formats:
```json
{
  "filename": "input.txt",
  "statistics": {
    "lines": 100,
    "words": 500,
    "characters": 2500
  },
  "topWords": [
    {"word": "example", "count": 10},
    {"word": "test", "count": 8}
  ]
}
```

### 5. Project Structure
```
cli-processor/
├── cmd/
│   └── root.go       # Main command definitions
├── internal/
│   ├── analyzer/     # Text analysis logic
│   ├── processor/    # File processing logic
│   └── output/       # Output formatting
├── pkg/
│   └── utils/        # Reusable utilities
├── main.go
└── go.mod
```

## Implementation Steps

### Step 1: Project Setup
1. Initialize Go module
```bash
mkdir cli-processor
cd cli-processor
go mod init cli-processor
```

2. Install required packages
```bash
go get github.com/spf13/cobra
go get github.com/spf13/viper
```

### Step 2: Basic CLI Framework
Create the basic command structure using Cobra:
```go
// cmd/root.go
package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "fileproc",
    Short: "A text file processing CLI",
    Long:  `A CLI tool for processing and analyzing text files with various options`,
}
```

### Step 3: Core Processing Logic
Implement the file processor:
```go
// internal/processor/processor.go
type FileProcessor struct {
    FilePath string
    Options  ProcessOptions
}

type ProcessOptions struct {
    IncludeWhitespace bool
    CaseSensitive     bool
    Context           int
}

func (fp *FileProcessor) Process() (Result, error) {
    // Implementation
}
```

### Step 4: Analysis Features
Implement text analysis functions:
```go
// internal/analyzer/analyzer.go
type TextStats struct {
    LineCount     int
    WordCount     int
    CharCount     int
    FrequentWords map[string]int
}

func AnalyzeText(content string) TextStats {
    // Implementation
}
```

## Advanced Features to Add

### 1. Concurrent Processing
- Implement worker pools for processing multiple files
- Use channels for result aggregation
- Add progress reporting for large files/directories

### 2. Memory Efficient Processing
- Process large files in chunks
- Implement streaming for huge files
- Use bufio.Scanner for efficient reading

### 3. Advanced Search Features
- Fuzzy search capabilities
- Support for multiple search patterns
- Pattern exclusion options

## Testing Requirements

1. Unit Tests:
```go
func TestWordCount(t *testing.T) {
    // Test cases
}

func TestPatternMatch(t *testing.T) {
    // Test cases
}
```

2. Integration Tests:
- Test with various file sizes
- Test with different character encodings
- Test concurrent processing

## Bonus Challenges
1. Add support for compressed files (zip, gzip)
2. Implement simple text transformation features
3. Add support for different character encodings
4. Create interactive mode with real-time updates

## Learning Focus Points
1. Go's file handling capabilities
2. Error handling patterns
3. Concurrent processing with goroutines
4. Command-line argument parsing
5. Structured project organization

***

What has been done so far:
- Basic file reading
- Word count
- Character count
- Line count
- Basic statistics

Next steps we were about to explore:
- Character count
- Pattern searching
- Different output formats
- Concurrent processing
