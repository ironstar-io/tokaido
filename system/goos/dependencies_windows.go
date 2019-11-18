package goos

import (
	"fmt"
	"os/exec"

	"github.com/ironstar-io/tokaido/utils"
)

// CheckDependencies - Currently Windows noop
func CheckDependencies() {

}

// CheckAndChocoInstall ...
func CheckAndChocoInstall(program string) *string {
	_, err := exec.LookPath(program)
	if err != nil {
		fmt.Println("     " + program + " isn't installed. Tokaido will install it with Chocolatey")
		utils.StdoutCmd("choco", "install", program, "-y")
		fmt.Println(program + " installed successfully")
	}

	return nil
}
