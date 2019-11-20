package goos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

// CreateSymlink - Check if a file/folder is writable
func CreateSymlink(path string) error {
	tokBinPath := filepath.Join(fs.HomeDir(), "bin", "tok")

	if fs.CheckExists(tokBinPath) == true {
		// Remove any existing soft link
		err := os.Remove(tokBinPath)
		if err != nil {
			return err
		}
	}

	// Create a new symbolic or "soft" link
	// os.Symlink requires elevated permissions in Windows, we can avoid that using `ln -s`. https://github.com/golang/go/issues/22874
	err := os.Symlink(path, tokBinPath)
	if err != nil {
		// As a fallback, this may be flimsy on some systems
		tbp := strings.ReplaceAll(tokBinPath, "C:\\", "/c/")
		t := strings.ReplaceAll(tbp, "\\", "/")
		pa := strings.ReplaceAll(path, "C:\\", "/c/")
		p := strings.ReplaceAll(pa, "\\", "/")
	
		_, err := utils.CommandSubSplitOutput("ln", "-s", p, t)
		if err != nil {
			fmt.Println("Unable to create a symlink when installing Tokaido.")
			fmt.Println("You may need to manually add a symbolic link or directly add the Tokaido binary to your PATH\n")
			
			return err
		}
	}

	return nil
}
