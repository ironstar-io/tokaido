package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"os/exec"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	utils.CheckPathHard("brew")

	CheckBrew()
	CheckAndBrewInstall("unison")
	CheckDockersync()
}

// CheckBrew ...
func CheckBrew() *string {
	_, err := exec.LookPath("brew")
	if err != nil {
		utils.StdoutCmd("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
	}

	fmt.Println("  ✓  brew")
	return nil
}

// CheckDockersync - Root executable
func CheckDockersync() *string {
	_, err := exec.LookPath("unison-fsmonitor")
	if err != nil {
		utils.StdoutCmd("brew", "tap", "eugenmayer/dockersync")
		utils.StdoutCmd("brew", "install", "eugenmayer/dockersync/unox")
	}

	fmt.Println("  ✓  unison-fsmonitor")
	return nil
}

// CheckAndBrewInstall - Root executable
func CheckAndBrewInstall(program string) *string {
	_, err := exec.LookPath(program)
	if err != nil {
		utils.StdoutCmd("brew", "install", program)
	}

	fmt.Println("  ✓ ", program)
	return nil
}
