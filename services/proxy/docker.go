package proxy

import (
	"github.com/ironstar-io/tokaido/utils"

	"log"
	"path/filepath"
	"strings"
)

// DockerComposeUp - Lift all containers in the proxy compose file in detached mode
func DockerComposeUp() {
	ComposeStdout("up", "-d")
}

// DockerComposeDown ...
func DockerComposeDown() {
	ComposeStdout("down")
}

// RestartContainer ...
func RestartContainer(container string) {
	ComposeStdout("restart", container)
}

// ComposeStdout - Convenience method for docker-compose shell commands
func ComposeStdout(args ...string) {
	composeParams := composeArgs(args...)

	utils.StdoutCmd("docker-compose", composeParams...)
}

func composeArgs(args ...string) []string {
	composeFile := []string{"-f", filepath.Join(getProxyDir(), "docker-compose.yml")}

	return append(composeFile, args...)
}

// LocalPort - Return the local port of a container
func LocalPort(containerName string, containerPort string) string {
	// Example return: "unison:32757"
	cs := utils.StdoutCmd("docker-compose", "-f", getComposePath(), "port", containerName, containerPort)

	return strings.Split(cs, ":")[1]
}

// UnisonPort - Return the local port for unison
func UnisonPort() string {
	p := LocalPort("unison", "5000")
	if p == "" {
		log.Fatal("Unison container doesn't appear to be running!")
	}

	return p
}
