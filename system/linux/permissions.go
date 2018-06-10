package linux

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// GetPermissionsMask - Return the permissions mask for a file/directory
func GetPermissionsMask(path string) string {
	return utils.NoFatalStdoutCmd("stat", "-c", "%a %n", path)
}
