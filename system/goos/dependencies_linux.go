package goos

import (
	"log"
	"os/exec"

	"github.com/ironstar-io/tokaido/system/console"
)

// CheckDependencies - Root executable
func CheckDependencies() *string {
	errDocker := CheckDocker()
	errDockerCompose := CheckDockerCompose()
	errUnison := CheckUnison()
	errUnisonFsmonitor := CheckUnisonFsmonitor()

	if (errDocker != nil) ||
		(errDockerCompose != nil) ||
		(errUnison != nil) ||
		(errUnisonFsmonitor != nil) {
		log.Fatalf("Unable to proceed. One or more dependency checks failed")
	}

	return nil
}

// CheckDocker ...
func CheckDocker() error {
	_, err := exec.LookPath("docker")
	if err != nil {
		console.Println(`😓  docker is missing. You need to install Docker CE or Docker Machine`, "×")
		return err
	}

	return nil
}

// CheckDockerCompose ...
func CheckDockerCompose() error {
	_, err := exec.LookPath("docker-compose")
	if err != nil {
		console.Println(`😓  docker-compose is missing. Please visit https://docs.docker.com/compose/install/ for install instructions`, "x")
		return err
	}

	return nil
}

// CheckUnison ...
func CheckUnison() error {
	_, err := exec.LookPath("unison")
	if err != nil {
		console.Println(`😓  unison is missing. Please follow the install instructions at https://tokaido.io/docs/getting-started/#installing-on-linux`, "x")
		return err
	}

	return nil
}

// CheckUnisonFsmonitor ...
func CheckUnisonFsmonitor() error {
	_, err := exec.LookPath("unison-fsmonitor")
	if err != nil {
		console.Println(`😓  unison-fsmonitor is missing. Please follow the install instructions at https://tokaido.io/docs/getting-started/#installing-on-linux`, "x")

		return err
	}

	return nil
}
