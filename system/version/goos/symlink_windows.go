package goos

import (
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/fs"
)

// CreateSymlink - Check if a file/folder is writable
func CreateSymlink(path string) error {
	// Depending on OS, figure out the correct tok path and create symlink
	var tokBinPath = filepath.Join("C:/", "Program Files", "Ironstar", "Tokaido", "tok.exe")

	if fs.CheckExists(tokBinPath) == true {
		// Remove any existing soft link
		err := os.Remove(tokBinPath)
		if err != nil {
			return err
		}
	}

	// Create a new symbolic or "soft" link
	err := os.Symlink(tokBinPath, path)
	if err != nil {
		return err
	}

	return nil
}
