// +build !unix

package fs

import (
	// "bitbucket.org/ironstar/tokaido-cli/utils"

	"os"
)

// Writable ..
func Writable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// utils.DebugString(`Path doesn't exist: ` + path)
		return false
	}

	err = nil
	if !info.IsDir() {
		// utils.DebugString(`Path isn't a directory: ` + path)
		return false
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		// utils.DebugString(`Write permission is not set on this file: ` + path)
		return false
	}

	return true
}
