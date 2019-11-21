package goos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

// ActivateSavedVersion - Copies the specified version (which may be downloaded previously) into /usr/local/bin on macOS
func ActivateSavedVersion(version string) bool {
	// activePath is where the currently active version of Tokaido is found, such as /c/Users/Frank/bin/tok
	activePath := filepath.Join(fs.HomeDir(), "bin", "tok")
	tbp := strings.ReplaceAll(activePath, "C:\\", "/c/")

	// savePath is where to save the Tokaido binary, such as /c/Users/Frank/AppData/Local/Ironstar/Tokaido/{version}/tok
	savePath := filepath.Join(fs.HomeDir(), constants.BaseInstallPathWindows, version)

	// Check if the requested version is not downloaded already
	p := filepath.Join(fs.HomeDir(), savePath, "tok")
	if fs.CheckExists(p) != true {
		utils.DebugString("Tokaido version [" + version + "] was not found at [" + p + "], downloading a new copy...")

		_, err := DownloadTokBinary(version)
		if err != nil {
			fmt.Println("Unexpected error downloading that version: " + err.Error())
			os.Exit(1)
		}
	}

	// Remove any existing global binary
	utils.DebugString("removing any existing Tokaido version at [" + activePath + "]")
	err := os.Remove(activePath)
	if err != nil {
		utils.DebugErrOutput(err)
	}

	// Remove any existing copy of Tokaido at ~/bin/tok
	fs.Remove(activePath)

	// Copy the specified version to ~/bin/tok
	fs.Copy(p, activePath)

	// Make sure the version is executable
	err = os.Chmod(activePath, 0777)
	if err != nil {
		fmt.Println("Unexpected error granting execute permissions to [" + activePath + "]: " + err.Error())
		os.Exit(1)
	}

	return true
}
