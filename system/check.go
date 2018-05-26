package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"os"
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

// CheckAndCreateFolder - Check existence of `./.tok` folder
func CheckAndCreateFolder(path string) {
	_, err := os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
