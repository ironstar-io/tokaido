package system

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"log"
)

// UnisonLocalPort - Return the local port for unison
func UnisonLocalPort() string {
	unisonPort := utils.StdoutCmd("bash", "-c", "docker-compose port unison 5000 | cut -d':' -f2")

	if unisonPort == "" {
		log.Fatal("Unison is not correctly started in the Docker container")
	}

	return unisonPort
}
