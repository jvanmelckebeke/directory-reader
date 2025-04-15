package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// WithContext adds context information to an error
func WithContext(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// WithFile adds file information to an error
func WithFile(err error, filePath string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("file '%s': %w", filePath, err)
}

// SourceLocation returns the file and line where the function was called
func SourceLocation() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "unknown location"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// WrapWithRecovery executes a function and converts any panic to an error
func WrapWithRecovery(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = fmt.Errorf("%s: panic recovered: %v", SourceLocation(), e)
			default:
				err = fmt.Errorf("%s: panic recovered: %v", SourceLocation(), r)
			}
		}
	}()

	return fn()
}
