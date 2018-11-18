package goos

import (
	"fmt"
	"os/exec"

	"github.com/ironstar-io/tokaido/utils"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	CheckBrew()
	CheckAndBrewInstall("node")
	CheckAndBrewInstall("yarn")
}

// CheckBrew ...
func CheckBrew() *string {
	_, err := exec.LookPath("brew")
	if err != nil {
		fmt.Println("    Homebrew isn't installed. Tokaido will install it")
		utils.StdoutCmd("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
		fmt.Println("  √  brew")
	}

	return nil
}

// CheckAndBrewInstall - Root executable
func CheckAndBrewInstall(program string) *string {
	_, err := exec.LookPath(program)
	if err != nil {
		fmt.Println("     " + program + " isn't installed. Tokaido will install it with Homebrew")
		utils.StdoutCmd("brew", "install", program)
		fmt.Println("  √ ", program)
	}

	return nil
}
