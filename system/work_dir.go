package system

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"os"
	"path/filepath"
)

// WorkDir - Return the name of the directory the binary is being run in
func WorkDir() string {
	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("Tokaido encountered a fatal error and had to stop")
		return utils.FatalError(err)
	}
	return filepath.Dir(ex)
}
