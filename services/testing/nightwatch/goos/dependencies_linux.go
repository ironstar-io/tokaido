package goos

import (
	"fmt"
	"log"
	"os/exec"
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
	d := map[string]bool{
		"node": false,
		"yarn": false,
	}

	m := map[string]string{
		"node": "😓  NodeJS is missing and required for Nightwatch testing. Please visit https://nodejs.org/en/download/package-manager/ for install instructions",
		"yarn": "😓  Yarn is missing and required for Nightwatch testing. Please visit https://yarnpkg.com/lang/en/docs/install for install instructions",
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
}
