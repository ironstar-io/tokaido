package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
	"fmt"
	"os/exec"
)

// CheckAndBrewInstall - Root executable
func CheckAndBrewInstall(program string) string {
	_, err := exec.LookPath(program)
	if err != nil {
		utils.StdoutCmd("brew", "install", program)
	}

	fmt.Println("\t✓ ", program)
	return ""
}
