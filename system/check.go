package system

import (
	"bitbucket.org/ironstar/tokaido-cli/system/goos"

	"os"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	goos.CheckDependencies()
}

// CheckAndCreateFolder - Check existence of `./.tok` folder
func CheckAndCreateFolder(path string) {
	_, err := os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
