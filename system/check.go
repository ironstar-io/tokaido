package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	var GOOS = utils.CheckOS()

	fmt.Println("Checking required dependencies")
	if GOOS == "linux" {
		linux.CheckDependencies()
	}

	if GOOS == "osx" {
		osx.CheckDependencies()
	}
}
