package fs

import (
	"os"
)

// CheckExists - Check if a file exists or not
func CheckExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
