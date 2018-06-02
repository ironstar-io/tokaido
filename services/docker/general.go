package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"strings"
)

// CheckAllContainers ...
func CheckAllContainers() {

}

// LocalPort - Return the local port of a container
func LocalPort(containerName string, containerPort string) string {
	CheckForTokComposeFile()

	// Example return: "unison:32757"
	containerStr := utils.SilentStdoutCmd("docker-compose", "port", containerName, containerPort)

	return strings.Split(containerStr, ":")[1]
}

// ComposeStdout - Convenience method for docker-compose shell commands
func ComposeStdout(args ...string) {
	CheckForTokComposeFile()

	composeFile := []string{"-f", fs.WorkDir() + "/docker-compose.tok.yml"}
	composeParams := append(composeFile, args...)

	utils.StdoutCmd("docker-compose", composeParams...)
}
