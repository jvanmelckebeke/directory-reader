package ignore

import (
	"os"
	"path/filepath"
	"strings"
)

var languageExtensionMap = map[string]string{
	".py":    "python",
	".js":    "javascript,node",
	".ts":    "typescript,node",
	".java":  "java",
	".go":    "go",
	".rb":    "ruby",
	".php":   "php",
	".cs":    "csharp",
	".cpp":   "c++",
	".c":     "c",
	".rs":    "rust",
	".swift": "swift",
	".kt":    "kotlin",
	".scala": "scala",
	".sh":    "bash",
}

// DetectLanguagesInDirectory scans a directory and returns a comma-separated list of detected languages
func DetectLanguagesInDirectory(rootDir string) (string, error) {
	detectedLangs := make(map[string]struct{})

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check file extension
		ext := filepath.Ext(path)
		if lang, exists := languageExtensionMap[ext]; exists {
			// Split in case a single extension maps to multiple languages (like .ts -> typescript,node)
			for _, l := range strings.Split(lang, ",") {
				detectedLangs[l] = struct{}{}
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// Convert detected languages to a comma-separated string
	var langs []string
	for lang := range detectedLangs {
		langs = append(langs, lang)
	}
	return strings.Join(langs, ","), nil
}
