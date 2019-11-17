package utils

import (
	"fmt"
	"log"
	"os/exec"
)

var packageMap = map[string]string{
	"docker-compose": `'docker-compose' was not found on your PATH, it is required in order to be able to run Tokaido.`,
}

// CheckCmdHard - Root executable
func CheckCmdHard(p string) string {
	_, err := exec.LookPath(p)
	if err != nil {
		fmt.Printf("The command '%s' must be available on your system to use Tokaido\n\n%s\n\n", p, packageMap[p])
		log.Fatal(err)
	}

	return ""
}

// CheckCmd - Checks if an executable is available and if not, returns false
func CheckCmd(p string) bool {
	_, err := exec.LookPath(p)
	if err != nil {
		return false
	}

	return true
}
