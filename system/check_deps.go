package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
)

// CheckDeps - Root executable
func CheckDeps() *string {
	var GOOS = utils.CheckOS()

	fmt.Println("\tChecking required dependencies")
	if GOOS == "linux" {
		return linux.CheckDeps()
	}

	if GOOS == "osx" {
		return osx.CheckDeps()
	}

	return nil
}
