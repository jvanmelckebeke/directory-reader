package ignore

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/denormal/go-gitignore"
)

var defaultIgnorePatterns = []string{
	"*.lock",
	"*.pyc",
	"__pycache__/",
	"node_modules/",
	".DS_Store",
	".git/",
	"directory_content.md",
}

// LoadDefaultIgnorePatterns loads the default ignore patterns and patterns from .gitignore and .readerignore
func LoadDefaultIgnorePatterns(targetDirectory string) (*gitignore.GitIgnore, error) {
	ignorePatterns := defaultIgnorePatterns

	for _, ignoreFile := range []string{".gitignore", ".readerignore"} {
		patterns, err := loadIgnorePatterns(filepath.Join(targetDirectory, ignoreFile))
		if err == nil {
			ignorePatterns = append(ignorePatterns, patterns...)
		}
	}

	// Remove duplicates
	ignorePatterns = removeDuplicates(ignorePatterns)

	// Compile the ignore patterns
	ignorer, err := gitignore.New(strings.NewReader(strings.Join(ignorePatterns, "\n")), targetDirectory, nil)
	if err != nil {
		return nil, err
	}

	return ignorer, nil
}

// AddLanguageIgnorePatterns adds ignore patterns based on specified programming languages
func AddLanguageIgnorePatterns(ignorer *gitignore.GitIgnore, languages []string) error {
	for _, lang := range languages {
		patterns, err := fetchGitignoreForLanguage(lang)
		if err != nil {
			return err
		}
		for _, pattern := range patterns {
			ignorer.AddPattern(pattern)
		}
	}
	return nil
}

func loadIgnorePatterns(ignoreFilePath string) ([]string, error) {
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)


