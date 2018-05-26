package utils

import (
	"fmt"
	"log"
	"os/exec"
)

var packageMap = map[string]string{
	"docker-compose": `no docker-compose??`,
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
