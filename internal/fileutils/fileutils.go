package fileutils

import (
	"bytes"
	"github.com/jvanmelckebeke/directory-reader/internal/errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Common binary file signatures/magic numbers
var binarySignatures = [][]byte{
	{0x7F, 'E', 'L', 'F'},                         // ELF binary
	{0x4D, 0x5A},                                  // Windows Executable (MZ)
	{0x50, 0x4B, 0x03, 0x04},                      // ZIP archive including docx, jar, etc.
	{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}, // PNG image
	{0xFF, 0xD8, 0xFF},                            // JPEG image
	{0x47, 0x49, 0x46, 0x38},                      // GIF image
	{0x25, 0x50, 0x44, 0x46},                      // PDF document
}

// CommonBinaryExtensions contains known binary file extensions
var CommonBinaryExtensions = map[string]bool{
	".exe":  true,
	".dll":  true,
	".so":   true,
	".bin":  true,
	".jar":  true,
	".zip":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ttf":  true,
	".woff": true,
}

// IsBinaryFile checks if a file is binary using multiple methods
func IsBinaryFile(filePath string) (bool, error) {
	// First check the extension
	ext := filepath.Ext(filePath)
	if CommonBinaryExtensions[ext] {
		return true, nil
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return false, errors.WithFile(err, filePath)
	}
	defer file.Close()

	// Check file size
	info, err := file.Stat()
	if err != nil {
		return false, errors.WithFile(err, filePath)
	}

	// Empty files are considered text files
	if info.Size() == 0 {
		return false, nil
	}

	// Read the beginning of the file to check for binary signatures
	var maxRead int64 = 8000
	if maxRead > info.Size() {
		maxRead = info.Size()
	}

	buf := make([]byte, maxRead)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, errors.WithFile(err, filePath)
	}

	readBytes := buf[:n]

	// Check for binary signatures (magic numbers)
	for _, sig := range binarySignatures {
		if len(readBytes) >= len(sig) && bytes.Equal(readBytes[:len(sig)], sig) {
			return true, nil
		}
	}

	// Check content type using http.DetectContentType
	contentType := http.DetectContentType(readBytes)
	if !strings.HasPrefix(contentType, "text/") &&
		!strings.Contains(contentType, "json") &&
		!strings.Contains(contentType, "xml") &&
		!strings.Contains(contentType, "javascript") {
		return true, nil
	}

	// Check if the content is valid UTF-8
	if !utf8.Valid(readBytes) {
		return true, nil
	}

	// Count binary characters
	binCount := 0
	totalBytes := len(readBytes)

	for _, b := range readBytes {
		if b == 0 {
			return true, nil // Null byte definitely indicates binary
		}
		// Count control characters excluding common whitespace chars
		if b < 8 || (b > 13 && b < 32) {
			binCount++
		}
	}

	// If more than 10% of bytes are binary characters, consider it binary
	if binCount > totalBytes/10 {
		return true, nil
	}

	return false, nil
}
