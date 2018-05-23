package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// Up - Lift all containers in the compose file
func Up() {
	utils.StdoutCmd("docker-compose", "up", "-d")
}

// Down - Pull down all the containers in the compose file
func Down() {
	utils.StdoutCmd("docker-compose", "down")
}

// Status - Print the container status to the console
func Status() {
	utils.StdoutCmd("docker-compose", "ps")
}
