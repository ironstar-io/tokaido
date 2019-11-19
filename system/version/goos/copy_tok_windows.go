package goos

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/fs"
)

// CopyTokBinary - Copy current binary to the local tok bin directory
func CopyTokBinary(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), "AppData", "Local", "Ironstar", "Tokaido", version)
	b := filepath.Join(p, "tok")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory")

		log.Fatal(err)
	}

	return b, nil
}
