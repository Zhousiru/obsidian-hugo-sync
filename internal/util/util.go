package util

import (
	"errors"
	"os"
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
