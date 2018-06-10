package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"os/exec"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	utils.CheckCmdHard("brew")

	CheckBrew()
	CheckAndBrewInstall("unison")
	CheckDockersync()
}

// CheckBrew ...
func CheckBrew() *string {
	_, err := exec.LookPath("brew")
	if err != nil {
		fmt.Println("     Homebrew isn't installed. Tokaido will install it")
		utils.StdoutCmd("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
		fmt.Println("  ✓  brew")
	}

	return nil
}

// CheckDockersync - Root executable
func CheckDockersync() *string {
	_, err := exec.LookPath("unison-fsmonitor")
	if err != nil {
		fmt.Println("     unison-fsmonitor is missing. Tokaido will install it with Homebrew")
		utils.StdoutCmd("brew", "tap", "eugenmayer/dockersync")
		utils.StdoutCmd("brew", "install", "eugenmayer/dockersync/unox")
		fmt.Println("  ✓  unison-fsmonitor")
	}

	return nil
}

// CheckAndBrewInstall - Root executable
func CheckAndBrewInstall(program string) *string {
	_, err := exec.LookPath(program)
	if err != nil {
		fmt.Println("     " + program + " isn't installed. Tokaido will install it with Homebrew")
		utils.StdoutCmd("brew", "install", program)
		fmt.Println("  ✓ ", program)
	}

	return nil
}
