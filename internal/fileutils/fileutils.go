package fileutils

import (
	"io"
	"os"
	"unicode/utf8"
)

// IsBinaryFile checks if a file is binary by reading a portion of it
func IsBinaryFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	const maxBytes = 8000 // Read up to 8000 bytes to determine if the file is binary
	buf := make([]byte, maxBytes)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Check if the content is valid UTF-8
	if !isText(buf[:n]) {
		return true, nil
	}
	return false, nil
}

func isText(data []byte) bool {
	// Check if data is valid UTF-8
	if !utf8.Valid(data) {
		return false
	}
	// Additional heuristic: check for control characters
	for _, b := range data {
		if b == 0 {
			return false // Null byte detected
		}
		if b < 0x09 {
			return false // Control characters below tab
		}
	}
	return true
}
