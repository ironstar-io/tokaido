package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"

	"log"
)

// DockerUp - Lift the unison container
func DockerUp() {
	docker.ComposeStdout("up", "-d", "unison")
}

// LocalPort - Return the local port for unison
func LocalPort() string {
	unisonPort := docker.LocalPort("unison", "5000")
	if unisonPort == "" {
		log.Fatal("Unison container doesn't appear to be running!")
	}

	return unisonPort
}
