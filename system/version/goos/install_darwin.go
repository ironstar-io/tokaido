package goos

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

var baseInstallPath = "/.tok/bin/"
var baseBinaryURL = "https://github.com/ironstar-io/tokaido/releases/download/"
var binaryName = "tok-macos"

// GetInstallPath - Check if tok version is installed or not
func GetInstallPath(version string) string {
	p := filepath.Join(baseInstallPath, version, "bin", "tok")
	if fs.CheckExists(p) == true {
		return p
	}

	return ""
}

// Install - Install a selected tok version and returns install path
func Install(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), baseInstallPath, version)
	b := filepath.Join(p, "tok")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory")

		log.Fatal(err)
	}

	err = utils.DownloadFile(b, baseBinaryURL+version+"/"+binaryName)
	if err != nil {
		return "", err
	}

	return b, nil
}
