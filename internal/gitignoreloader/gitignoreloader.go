package gitignoreloader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// LoadGitIgnoreForLanguages fetches .gitignore files for the specified languages from GitHub
func LoadGitIgnoreForLanguages(languages []string) ([]string, error) {
	var allPatterns []string

	for _, lang := range languages {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			continue
		}

		patterns, err := fetchGitIgnore(lang)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch .gitignore for '%s': %v", lang, err)
		}
		allPatterns = append(allPatterns, patterns...)
	}

	return allPatterns, nil
}

func fetchGitIgnore(language string) ([]string, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/github/gitignore/main/%s.gitignore", normalizeLanguageName(language))
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	patterns := parseGitIgnoreContent(string(content))
	return patterns, nil
}

func normalizeLanguageName(language string) string {
	// Capitalize the first letter and handle special cases
	switch strings.ToLower(language) {
	case "c++":
		return "C++.gitignore"
	case "c":
		return "C.gitignore"
	case "c#":
		return "C%23.gitignore"
	case "f#":
		return "F%23.gitignore"
	case "go", "golang":
		return "Go.gitignore"
	case "objective-c":
		return "Objective-C.gitignore"
	case "vim":
		return "Vim.gitignore"
	default:
		// Capitalize the first letter
		return strings.Title(strings.ToLower(language))
	}
}

func parseGitIgnoreContent(content string) []string {
	var patterns []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns
}
