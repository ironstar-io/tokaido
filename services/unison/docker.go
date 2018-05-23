package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"log"
)

// DockerUp - Lift the unison container
func DockerUp() {
	fmt.Println("Firing up the Unison container!")

	utils.StdoutCmd("docker-compose", "up", "-d", "unison")
}

// LocalPort - Return the local port for unison
func LocalPort() string {
	unisonPort := utils.SilentStdoutCmd("bash", "-c", "docker-compose port unison 5000 | cut -d':' -f2")

	if unisonPort == "" {
		log.Fatal("Unison is not correctly started in the Docker container")
	}

	return unisonPort
}
