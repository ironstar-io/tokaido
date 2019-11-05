package goos

import (
	"os"
)

var tokBinPath = "/usr/local/bin/tok"

// CreateSymlink - Check if a file/folder is writable
func CreateSymlink(path string) error {
	// Depending on OS, figure out the correct tok path and create symlink

	// Remove any existing soft link
	err := os.Remove(tokBinPath)
	if err != nil {
		return err
	}

	// Create a new symbolic or "soft" link
	err = os.Symlink(tokBinPath, path)
	if err != nil {
		return err
	}

	return nil
}
