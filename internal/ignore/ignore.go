package ignore

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
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
func LoadDefaultIgnorePatterns(targetDirectory string) ([]string, error) {
	ignorePatterns := defaultIgnorePatterns

	for _, ignoreFile := range []string{".gitignore", ".readerignore"} {
		patterns, err := loadIgnorePatterns(filepath.Join(targetDirectory, ignoreFile))
		if err == nil {
			ignorePatterns = append(ignorePatterns, patterns...)
		}
	}

	// Remove duplicates
	ignorePatterns = removeDuplicates(ignorePatterns)

	return ignorePatterns, nil
}

// CompileIgnorePatterns compiles the ignore patterns into a GitIgnore object
func CompileIgnorePatterns(patterns []string) *ignore.GitIgnore {
	return ignore.CompileIgnoreLines(patterns...)
}

func FetchLanguageAndDefaultIgnorePatterns(languages string) ([]string, error) {
	baseURL := "https://www.toptal.com/developers/gitignore/api/linux,pycharm"
	if languages != "" {
		baseURL += "," + strings.ReplaceAll(languages, ",", ",")
	}

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ignore patterns: %v", err)
	}
	defer resp.Body.Close()

	// check if the status is ok, however, 404 is acceptable as it will return a best-effort response, even if the language is not found
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("failed to fetch ignore patterns: HTTP %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("WARNING: Some languages were not found. gitignore.io returned a best effort response")
	}

	var patterns []string
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
		if err == io.EOF {
			break
		}
	}

	// Load default ignore patterns
	defaultPatterns, err := LoadDefaultIgnorePatterns(".")
	if err != nil {
		return nil, err
	}
	patterns = append(patterns, defaultPatterns...)

	return patterns, nil
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
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}

func removeDuplicates(strings []string) []string {
	uniqueStrings := make(map[string]struct{})
	for _, s := range strings {
		uniqueStrings[s] = struct{}{}
	}
	var result []string
	for s := range uniqueStrings {
		result = append(result, s)
	}
	return result
}
