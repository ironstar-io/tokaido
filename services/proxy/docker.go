package proxy

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
	"fmt"
	
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

// DockerComposeUp - Lift all containers in the proxy compose file in detached mode
func DockerComposeUp() {
	composeStdout("up", "--remove-orphans", "-d")
}

// dockerComposeDown ...
func dockerComposeDown() {
	composeStdout("down")
}

// dockerComposeRemoveProxy will stop and delete an existing proxy container
// this is necessary to stop conflicts and hard crashes when changing the sync service
func dockerComposeRemoveProxy() {
	composeStdout("kill", "proxy")
	composeStdout("rm", "proxy")
}

// restartContainer ...
func restartContainer(container string) {
	composeStdout("restart", container)
}

// composeStdout - Convenience method for docker-compose shell commands
func composeStdout(args ...string) {
	composeParams := composeArgs(args...)

	// Ensure the proxy docker-compose file exists before trying to run commands against it
	if !fs.CheckExists(filepath.Join(getProxyDir(), "docker-compose.yml")) {
		fmt.Println("A proxy compose file was not found prior to running a docker-compose command")
		fmt.Println("There may be a problem with your proxy service")

		return
	}

	utils.StdoutCmd("docker-compose", composeParams...)
}

func composeArgs(args ...string) []string {
	composeFile := []string{"-f", filepath.Join(getProxyDir(), "docker-compose.yml")}

	return append(composeFile, args...)
}

// localPort - Return the local port of a container
func localPort(containerName string, containerPort string) string {
	// Example return: "unison:32757"
	cs := utils.StdoutCmd("docker-compose", "-f", getComposePath(), "port", containerName, containerPort)

	p := strings.Split(cs, ":")
	if len(p) <= 1 {
		log.Fatal("The required proxy container '" + containerName + "' is not running. Perhaps retry `tok up`?  Exiting...")
		return ""
	}

	return p[1]
}

// getContainerName ...
func getContainerName(name string) (string, error) {
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

// getContainerIP ...
func getContainerIP(name string) (string, error) {
	cn, err := getContainerName(name)
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

// disconnectNetworkEndpoint ...
func disconnectNetworkEndpoint(network, endpoint string) error {
	_, err := utils.BashStringSplitOutput(`docker network disconnect -f ` + network + ` ` + endpoint)
	if err != nil {
		return err
	}

	return nil
}

// unisonPort - Return the local port for unison
func unisonPort() string {
	p := localPort("unison", "5000")
	if p == "" {
		log.Fatal("Unison container doesn't appear to be running!")
	}

	return p
}
