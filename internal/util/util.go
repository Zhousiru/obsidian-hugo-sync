package util

import (
	"errors"
	"mime"
	"os"
	"path/filepath"
)

// IsExist checks existence of a path.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		panic(err)
	}

	return true
}

// TrimExt returns filename without ext.
func TrimExt(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

// TrimPrefix trims prefix of specific string.
func TrimPrefix(s, prefix string) string {
	return s[len(prefix):]
}

// GetExt returns ext of specified filename without dot.
func GetExt(filename string) string {
	return filepath.Ext(filename)[1:]
}

// GetContentType returns the content type of specified extension (without dot).
func GetContentType(ext string) string {
	return mime.TypeByExtension("." + ext)
}
