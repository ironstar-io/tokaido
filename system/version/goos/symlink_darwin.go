package goos

import (
	"os"
)

// CreateSymlink - Check if a file/folder is writable
func CreateSymlink(version string) error {
	// Depending on OS, figure out the correct tok path and create symlink

	// create a new symbolic or "soft" link
	err := os.Symlink("file.txt", "file-symlink.txt")
	if err != nil {
		return err
	}

	return nil
}
