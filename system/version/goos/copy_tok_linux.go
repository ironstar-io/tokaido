package goos

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
)

// CopyTokBinary - Copy current binary to the local tok bin directory
func CopyTokBinary(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathLinux, version)
	b := filepath.Join(p, "tok")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory")

		log.Fatal(err)
	}

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	fs.Copy(ex, b)
	// Change file permission bit
	err = os.Chmod(b, 0755)
	if err != nil {
		panic(err)
	}

	// Change file ownership.
	err = os.Chown(b, os.Getuid(), os.Getgid())
	if err != nil {
		panic(err)
	}

	return b, nil
}
