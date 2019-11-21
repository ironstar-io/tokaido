package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

var binaryName = "tok-windows-amd64.exe"

// SaveTokBinary - Saves the running instance of the Tokaido binary to a persistent path in the user's bin folder
func SaveTokBinary(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathWindows, version)
	b := filepath.Join(p, "tok")

	utils.DebugString("Saving the running Tokaido version [" + version + "] to [" + p + "]")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory: ", err.Error())
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not determine the current wokring directory", err.Error())
		os.Exit(1)
	}
	ib := filepath.Join(pwd, "tokaido", binaryName)

	fs.Copy(ib, b)
	// Change file permission bit
	err = os.Chmod(b, 0755)
	if err != nil {
		fmt.Println("Could not ensure correct ownership on ["+b+"]: ", err.Error())
		os.Exit(1)
	}

	return b, nil
}
