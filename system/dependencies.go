package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/goos"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	goos.CheckDependencies()
}
