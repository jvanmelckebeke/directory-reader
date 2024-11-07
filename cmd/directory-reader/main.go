package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jvanmelckebeke/directory-structure/internal/dirreader"
	"github.com/jvanmelckebeke/directory-structure/internal/ignore"
	"github.com/jvanmelckebeke/directory-structure/internal/tokenizer"
)

func main() {
	var (
		countTokens bool
		ignoreLangs string
	)

	// Parse command-line arguments
	flag.BoolVar(&countTokens, "count-tokens", false, "Count tokens in the generated directory_content.md file")
	flag.StringVar(&ignoreLangs, "ignore", "", "Comma-separated list of programming languages to ignore (e.g., 'go,python')")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--count-tokens] [--ignore=lang1,lang2] <target_directory>\n", os.Args[0])
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

	// Get the script name
	scriptName := filepath.Base(os.Args[0])

	// Load ignore patterns
	ignorePatterns, err := ignore.LoadDefaultIgnorePatterns(targetDirectory)
	if err != nil {
		fmt.Println("Error loading ignore patterns:", err)
		os.Exit(1)
	}

	// Process --ignore flag
	if ignoreLangs != "" {
		langs := strings.Split(ignoreLangs, ",")
		err := ignore.AddLanguageIgnorePatterns(ignorePatterns, langs)
		if err != nil {
			fmt.Println("Error adding language ignore patterns:", err)
			os.Exit(1)
		}
	}

	// Create the markdown file
	outputFile, err := dirreader.CreateMarkdownFile(scriptName, targetDirectory, ignorePatterns)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Markdown file '%s' has been created in '%s'.\n", outputFile, targetDirectory)

	// Count tokens if --count-tokens flag is provided
	if countTokens {
		tokenCount, err := tokenizer.CountFileTokens(outputFile)
		if err != nil {
			fmt.Println("Error counting tokens:", err)
			os.Exit(1)
		}
		formattedCount := humanize.Comma(int64(tokenCount))
		fmt.Printf("Number of tokens in '%s': %s\n", outputFile, formattedCount)
	}
}
