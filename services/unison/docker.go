package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"

	"fmt"
	"log"
)

// DockerUp - Lift the unison container
func DockerUp() {
	fmt.Println("Firing up the Unison container!")

	docker.ComposeStdout("up", "-d", "unison")
}

// LocalPort - Return the local port for unison
func LocalPort() string {
	unisonPort := docker.LocalPort("unison", "5000")
	if unisonPort == "" {
		log.Fatal("Unison is not correctly started in the Docker container")
	}

	return unisonPort
}
