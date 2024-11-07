package dirreader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jvanmelckebeke/directory-reader/internal/fileutils"
	"github.com/sabhiram/go-gitignore"
)

// WriteDirectoryStructure writes the directory structure to the provided writer
func WriteDirectoryStructure(w io.Writer, startPath string, ignorer *ignore.GitIgnore) error {
	var walk func(string, int) error
	walk = func(path string, level int) error {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		// Exclude directories and files based on ignore patterns
		var dirs []os.DirEntry
		var files []os.DirEntry

		relPath, err := filepath.Rel(startPath, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			relPath = ""
		}

		for _, entry := range entries {
			entryPath := filepath.Join(relPath, entry.Name())
			matchPath := entryPath
			if entry.IsDir() {
				matchPath += string(os.PathSeparator)
			}

			if ignorer.MatchesPath(matchPath) {
				continue
			}

			if entry.IsDir() {
				dirs = append(dirs, entry)
			} else {
				isBin, err := fileutils.IsBinaryFile(filepath.Join(path, entry.Name()))
				if err != nil {
					return err
				}
				if isBin {
					continue // Skip binary files
				}
				files = append(files, entry)
			}
		}

		indent := strings.Repeat("  ", level)
		dirName := filepath.Base(path)
		if relPath == "" {
			dirName = filepath.Base(startPath)
		}
		fmt.Fprintf(w, "%s- %s/\n", indent, dirName)

		for _, file := range files {
			fmt.Fprintf(w, "%s  - %s\n", indent, file.Name())
		}

		for _, dir := range dirs {
			err := walk(filepath.Join(path, dir.Name()), level+1)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return walk(startPath, 0)
}
