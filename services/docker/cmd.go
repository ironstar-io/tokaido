package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"strings"
)

// Up - Lift all containers in the compose file
func Up() {
	utils.StdoutCmd("docker-compose", "up", "-d")
}

// Stop - Stop all containers in the compose file
func Stop() {
	utils.StdoutCmd("docker-compose", "stop")
}

// Down - Pull down all the containers in the compose file
func Down() {
	utils.StdoutCmd("docker-compose", "down")
}

// Status - Print the container status to the console
func Status() {
	utils.StdoutCmd("docker-compose", "ps")
}

// CheckAllContainers ...
func CheckAllContainers() {

}

// LocalPort - Return the local port of a container
func LocalPort(containerName string, containerPort string) string {
	// Example return: "unison:32757"
	containerStr := utils.SilentStdoutCmd("docker-compose", "port", containerName, containerPort)

	return strings.Split(containerStr, ":")[1]
}
