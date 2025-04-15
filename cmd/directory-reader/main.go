package main

import (
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jvanmelckebeke/directory-reader/internal/dirreader"
	"github.com/jvanmelckebeke/directory-reader/internal/ignore"
	"github.com/jvanmelckebeke/directory-reader/internal/tokenizer"
	"os"
	"path/filepath"
)

func main() {
	var (
		ignoreLangs string
	)

	// Parse command-line arguments - removed the count-tokens flag
	flag.StringVar(&ignoreLangs, "ignore", "", "Comma-separated list of programming languages to ignore (e.g., 'go,python')")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--ignore=lang1,lang2] <target_directory>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	targetDirectory := flag.Arg(0)
	if info, err := os.Stat(targetDirectory); err != nil || !info.IsDir() {
		fmt.Printf("Error: '%s' is not a valid directory.\n", targetDirectory)
		os.Exit(1)
	}

	// Auto-detect languages in the directory
	detectedLangs, err := ignore.DetectLanguagesInDirectory(targetDirectory)
	if err != nil {
		fmt.Println("Warning: Failed to auto-detect languages:", err)
	} else if detectedLangs != "" {
		fmt.Printf("Auto-detected languages: %s\n", detectedLangs)
		if ignoreLangs != "" {
			ignoreLangs = ignoreLangs + "," + detectedLangs
		} else {
			ignoreLangs = detectedLangs
		}
	}

	// Get the script name
	scriptName := filepath.Base(os.Args[0])

	// Default ignore patterns (Linux and PyCharm)
	ignorePatterns, err := ignore.FetchLanguageAndDefaultIgnorePatterns(ignoreLangs)

	if err != nil {
		fmt.Println("Error fetching ignore patterns:", err)
		os.Exit(1)
	}

	// Compile the ignore patterns
	ignorer := ignore.CompileIgnorePatterns(ignorePatterns)

	// Create the markdown file
	outputFile, err := dirreader.CreateMarkdownFile(scriptName, targetDirectory, ignorer)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Markdown file '%s' has been created in '%s'.\n", outputFile, targetDirectory)

	// Always count tokens now
	tokenCount, err := tokenizer.CountFileTokens(outputFile)
	if err != nil {
		fmt.Println("Error counting tokens:", err)
		os.Exit(1)
	}

	formattedCount := humanize.Comma(int64(tokenCount))

	// Show warning if token count exceeds 16k
	fmt.Printf("Token count: %s\n", formattedCount)
}
