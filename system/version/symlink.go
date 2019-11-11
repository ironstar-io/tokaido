package version

import (
	"github.com/ironstar-io/tokaido/system/version/goos"
)

// CreateSymlink - Check if a file/folder is writable
func CreateSymlink(version string) error {
	return goos.CreateSymlink(version)
}
