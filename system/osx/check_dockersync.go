package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
	"fmt"
	"os/exec"
)

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
