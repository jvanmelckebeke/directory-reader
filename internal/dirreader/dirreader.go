package dirreader

import (
	"fmt"
	"github.com/jvanmelckebeke/directory-reader/internal/fileutils"
	"github.com/sabhiram/go-gitignore"
	"os"
	"path/filepath"
	"sort"
)

type fileInfo struct {
	Path string
	Name string
	Rank fileutils.FileRank
}

// CreateMarkdownFile generates the directory_content.md file
func CreateMarkdownFile(scriptName, targetDirectory string, ignorer *ignore.GitIgnore) (string, error) {
	outputFile := filepath.Join(targetDirectory, "directory_content.md")
	mdFile, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer mdFile.Close()

	mdFile.WriteString("# Directory Structure\n\n")
	mdFile.WriteString("```\n")
	err = WriteDirectoryStructure(mdFile, targetDirectory, ignorer)
	if err != nil {
		return "", err
	}
	mdFile.WriteString("```\n\n")

	mdFile.WriteString("# File Contents\n\n")
	// Collect all files first
	var files []fileInfo

	err = filepath.WalkDir(targetDirectory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(targetDirectory, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			relPath = ""
		}

		if ignorer.MatchesPath(relPath) {
			return nil
		}

		if d.Name() == filepath.Base(outputFile) || d.Name() == scriptName {
			return nil
		}

		// Check if the file is binary
		isBin, err := fileutils.IsBinaryFile(path)
		if err != nil {
			return err
		}
		if isBin {
			return nil // Skip binary files
		}

		rank := fileutils.GetFileRank(relPath)
		files = append(files, fileInfo{
			Path: relPath,
			Name: d.Name(),
			Rank: rank,
		})

		return nil
	})

	if err != nil {
		return "", err
	}

	// Sort files by rank (higher rank files will be displayed later)
	sort.Slice(files, func(i, j int) bool {
		if files[i].Rank != files[j].Rank {
			return files[i].Rank < files[j].Rank
		}
		return files[i].Path < files[j].Path
	})

	// Now write the sorted files to the markdown
	for _, file := range files {
		mdFile.WriteString(fmt.Sprintf("## %s\n\n", file.Path))
		mdFile.WriteString("```\n")

		fileContent, err := os.ReadFile(filepath.Join(targetDirectory, file.Path))
		if err != nil {
			mdFile.WriteString(fmt.Sprintf("Error reading file: %s\n", err.Error()))
		} else {
			mdFile.Write(fileContent)
			mdFile.WriteString("\n")
		}
		mdFile.WriteString("```\n\n")
	}

	return outputFile, nil
}
