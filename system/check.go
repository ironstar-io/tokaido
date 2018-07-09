package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"

	"os"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	var GOOS = CheckOS()

	if GOOS == "linux" {
		linux.CheckDependencies()
	}

	if GOOS == "osx" {
		osx.CheckDependencies()
	}

	if GOOS == "windows" {
		// windows.CheckDependencies()  // TODO Windows
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
