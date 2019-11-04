package version

import (
	"github.com/ironstar-io/tokaido/system/version/goos"
)

// GetInstallPath - Check if tok version is installed or not
func GetInstallPath(version string) string {
	return goos.GetInstallPath(version)
}

// Install - Install a selected tok version
func Install(version string) (string, error) {
	return goos.Install(version)
}
