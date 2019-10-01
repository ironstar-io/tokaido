package proxy

import (
	"fmt"

	"github.com/ironstar-io/tokaido/services/docker"
)

// CreateProxyNetwork creates the `proxy` network if it does not already exist
func CreateProxyNetwork() {
	if docker.CheckNetworkUp() {
		return
	}

	docker.CreateNetwork("tokaido_proxy")
}

// getContainerProxyIP retrieves the IP address for this container connection to the 'tokaido_proxy' network
func getContainerProxyIP(containerName, projectName string) (string, error) {
	con := docker.GetContainer(containerName, projectName)

	if len(con.Names) == 0 {
		return "", fmt.Errorf("Data for the container [%s] could not be found", containerName)
	}

	if len(con.NetworkSettings.Networks) != 2 {
		return "", fmt.Errorf("The container [%s] did not have an accurate network setup", containerName)
	}

	if con.NetworkSettings.Networks["tokaido_proxy"].IPAddress == "" {
		return "", fmt.Errorf("The container [%s] was missing a 'tokaido_proxy' IP address", containerName)
	}

	return con.NetworkSettings.Networks["tokaido_proxy"].IPAddress, nil
}
