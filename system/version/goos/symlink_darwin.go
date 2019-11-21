package goos

import (
	"os"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
)

// CreateSymlink - Check if a file/folder is writable
func CreateSymlink(path string) error {
	// Depending on OS, figure out the correct tok path and create symlink

	if fs.CheckExists(constants.ActiveBinaryPathDarwin) == true {
		// Remove any existing soft link
		err := os.Remove(constants.ActiveBinaryPathDarwin)
		if err != nil {
			return err
		}
	}

	// Create a new symbolic or "soft" link
	err := os.Symlink(path, constants.ActiveBinaryPathDarwin)
	if err != nil {
		return err
	}

	return nil
}
