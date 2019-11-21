package goos

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
)

// SaveTokBinary - Saves the running instance of the Tokaido binary to a persistent path in the user's tok/bin folder
func SaveTokBinary(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathLinux, version)
	b := filepath.Join(p, "tok")

	utils.DebugString("Saving the running Tokaido version ["+version"] to ["+p+"]")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory: ", err.Error)
		os.Exit(1)
	}

	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Unexpected error obtaining this executable's path: ", err.Error())
		os.Exit(1)
	}

	fs.Copy(ex, b)
	// Change file permission bit
	err = os.Chmod(b, 0755)
	if err != nil {
		fmt.Println("Could not ensure correct permissions on ["+b+"]: ", err.Error())
		os.Exit(1)
	}

	// Change file ownership.
	err = os.Chown(b, os.Getuid(), os.Getgid())
	if err != nil {
		fmt.Println("Could not ensure correct ownership on ["+b+"]: ", err.Error())
		os.Exit(1)
	}

	return b, nil
}
