package system

import (
	"errors"
	"os"
)

// GetPermissionsMask - Return the permissions mask for a file/directory
func GetPermissionsMask(path string) (os.FileMode, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, errors.New("Tokaido doesn't have the required permissions on the file '" + path + "'")
	}

	return file.Mode().Perm(), nil
}
