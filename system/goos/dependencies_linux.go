package goos

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/ironstar-io/tokaido/system/wsl"
)

// commandCheck looks for the presence of command named 'n' on the system
// it returns true or false based on whether or not the command is available
func commandCheck(n string) bool {
	_, err := exec.LookPath(n)
	if err != nil {
		return false
	}

	return true
}

// CheckDependencies - Root executable
func CheckDependencies() {
	if wsl.IsWSL() {
		CheckWSLDependencies()

		return
	}

	d := map[string]bool{
		"docker":           false,
		"docker-compose":   false,
		"unison":           false,
		"unison-fsmonitor": false,
	}

	m := map[string]string{
		"docker":           "😓  docker is missing. You need to install Docker CE or Docker Machine",
		"docker-compose":   "😓  docker-compose is missing. Please visit https://docs.docker.com/compose/install/ for install instructions",
		"unison":           "😓  unison is missing. Please follow the install instructions at https://tokaido.io/docs/getting-started/#installing-on-linux",
		"unison-fsmonitor": "😓  unison-fsmonitor is missing. Please follow the install instructions at https://tokaido.io/docs/getting-started/#installing-on-linux",
	}

	var f bool
	for k := range d {
		if !commandCheck(k) {
			fmt.Println(m[k])
			f = true
		}
	}

	if f {
		log.Fatalf("Unable to proceed. One or more dependency checks failed")
	}

	CheckMaxUserWatches()
}

// CheckWSLDependencies - Root executable
func CheckWSLDependencies() {
	d := map[string]bool{
		"docker":         false,
		"docker-compose": false,
	}

	m := map[string]string{
		"docker":         "😓  docker is missing. You need to install Docker CE or Docker Machine",
		"docker-compose": "😓  docker-compose is missing. Please visit https://docs.docker.com/compose/install/ for install instructions",
	}

	var f bool
	for k := range d {
		if !commandCheck(k) {
			fmt.Println(m[k])
			f = true
		}
	}

	if f {
		log.Fatalf("Unable to proceed. One or more dependency checks failed")
	}

	CheckComposeTimeoutValue()
}
