package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
	"fmt"
	"os/exec"
)

// CheckBrew - Root executable
func CheckBrew() string {
	_, err := exec.LookPath("brew")
	if err != nil {
		utils.StdoutCmd("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
	}

	fmt.Println("\t✓  brew")
	return ""
}
