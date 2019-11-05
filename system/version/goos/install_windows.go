package goos

import (
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

var baseInstallPath = "/usr/local/ironstar/tokaido/"
var baseBinaryURL = "https://github.com/ironstar-io/tokaido/releases/download/"
var binaryName = "tok-windows.exe"

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
	p := filepath.Join(baseInstallPath, version, "bin", "tok")

	err := utils.DownloadFile(p, baseBinaryURL+version+"/"+binaryName)
	if err != nil {
		return "", err
	}

	return p, nil
}
