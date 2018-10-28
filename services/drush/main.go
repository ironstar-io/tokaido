package drush

import (
	"log"

	"github.com/ironstar-io/tokaido/services/docker"
)

// DockerUp - Lift the drush container
func DockerUp() {
	docker.ComposeStdout("up", "-d", "drush")
}

// LocalPort - Return the local port for drush
func LocalPort() string {
	p := docker.LocalPort("drush", "5000")
	if p == "" {
		log.Fatal("Drush container doesn't appear to be running!")
	}

	return p
}
