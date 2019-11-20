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
	tbp := strings.ReplaceAll(tokBinPath, "C:\\", "/c/")
	symlinkPath := strings.ReplaceAll(tbp, "\\", "/")
	pa := strings.ReplaceAll(path, "C:\\", "/c/")
	executablePath := strings.ReplaceAll(pa, "\\", "/")

	// Remove any existing soft link
	err := os.Remove(tokBinPath)
	if err != nil {
		utils.DebugErrOutput(err)
	}
	_, err = utils.CommandSubSplitOutput("rm", symlinkPath)
	if err != nil {
		utils.DebugErrOutput(err)
	}

	// Create a new symbolic or "soft" link
	// os.Symlink requires elevated permissions in Windows, we can avoid that using `ln -s`. https://github.com/golang/go/issues/22874
	err = os.Symlink(path, tokBinPath)
	if err != nil {
		// As a fallback, this may be flimsy on some systems
		utils.DebugString("The golang native method failed to create a symbolic link. Usually this is due to symlink creation requiring elevation. Falling back to attempt `mklink /d`")
		utils.DebugErrOutput(err)
		
		_, err := utils.CommandSubSplitOutput("mklink", "/d", symlinkPath, executablePath)
		if err != nil {
			utils.DebugString("The `mklink` command failed to create a symbolic link. Falling back to attempt `ln -s`")
			utils.DebugErrOutput(err)

			_, err := utils.CommandSubSplitOutput("ln", "-s", executablePath, symlinkPath)
			if err != nil {
				fmt.Println("Unable to create a symlink when installing Tokaido.")
				fmt.Println("You may need to manually add a symbolic link or directly add the Tokaido binary to your PATH\n")
				
				utils.DebugString("The `ln -s` command to create a symbolic link has also failed.")

				return err
			}
		}
	}

	return nil
}
