package fs

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs/goos"
)

// Writable - Check if a file/folder is writable
func Writable(path string) bool {
	return goos.Writable(path)
}
