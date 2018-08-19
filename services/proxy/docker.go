package proxy

import (
	"github.com/ironstar-io/tokaido/utils"

	"errors"
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

// GetContainerName ...
func GetContainerName(name string) (string, error) {
	cs, err := utils.BashStringSplitOutput("docker-compose -f " + getComposePath() + " ps " + name + " | grep " + name)
	if err != nil {
		return "", err
	}

	cn := strings.Split(cs, " ")[0]
	if cn == "" {
		return "", errors.New("Unable to find the container for " + name + ". Is it currently running?")
	}

	return cn, nil
}

// GetContainerIP ...
func GetContainerIP(name string) (string, error) {
	cn, err := GetContainerName(name)
	if err != nil {
		return "", err
	}

	ia, err := utils.BashStringSplitOutput(`docker inspect ` + cn + ` | grep "IPAddress\": \"1"`)
	if err != nil {
		return "", err
	}

	for _, s := range strings.Split(ia, `"`) {
		if strings.Contains(s, "1") {
			return s, nil
		}
	}

	return "", errors.New("Unable to find the internal IP for container " + name + ". Is it currently running?")
}

// UnisonPort - Return the local port for unison
func UnisonPort() string {
	p := LocalPort("unison", "5000")
	if p == "" {
		log.Fatal("Unison container doesn't appear to be running!")
	}

	return p
}
