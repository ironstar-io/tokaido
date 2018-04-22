package utils

import (
	"fmt"
	"os/exec"
)

var packageMap = map[string]string{
	"docker-compose": `no docker-compose??`,
}

// CheckPath - Root executable
func CheckPath(p string) string {
	_, err := exec.LookPath(p)
	if err != nil {
		fmt.Printf("The command '%s' must be available on your system to use Tokaido\n\n%s\n\n", p, packageMap[p])
		return FatalError(err)
	}

	return ""
}
