package fileutils

import (
	"path/filepath"
	"strings"
)

// FileRank represents the precedence of a file in the output
type FileRank int

// higher rank means more important means lower in the output (more recent in llm context)
const (
	RankTestFile       FileRank = 10 // Test files are usually less important
	RankImplementation FileRank = 20
	RankConfig         FileRank = 30
	RankMarkdown       FileRank = 40
	RankReadme         FileRank = 50 // Higher rank for READMEs
	RankDefault        FileRank = 25
)

// GetFileRank determines the rank of a file based on its path and name
func GetFileRank(path string) FileRank {
	fileName := filepath.Base(path)
	fileNameLower := strings.ToLower(fileName)
	ext := filepath.Ext(path)

	// Check for README files
	if strings.Contains(fileNameLower, "readme") {
		return RankReadme
	}

	// Check for test files
	if strings.Contains(fileNameLower, "test") || strings.Contains(fileNameLower, "_test.go") {
		return RankTestFile
	}

	// Check for configuration files
	switch ext {
	case ".json", ".yaml", ".yml", ".toml", ".xml", ".ini":
		return RankConfig
	case ".md", ".txt":
		return RankMarkdown
	}

	// Check for implementation files
	switch ext {
	case ".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".cs", ".rb", ".php":
		return RankImplementation
	}

	return RankDefault
}
