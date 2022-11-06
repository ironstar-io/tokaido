//go:build windows
// +build windows

package goos

import (
	"os"
)

// Writable - Check if a file/folder is writable
func Writable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false
	}

	return true
}
