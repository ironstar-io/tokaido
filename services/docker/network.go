package docker

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"

	"strings"
)

// CheckNetworkUp ...
func CheckNetworkUp() bool {
	p := strings.ToLower(conf.GetConfig().Tokaido.Project.Name)

	// Periods being replaced in recent versions of Docker for network names
	n := GetNetworkName(p)

	_, err := utils.CommandSubSplitOutput("docker", "network", "inspect", n)
	if err != nil {
		return false
	}

	return true
}

// GetGateway - Get the Gateway IP adress of the docker network
func GetGateway(projectName string) string {
	// Periods being replaced in recent versions of Docker for network names
	n := GetNetworkName(projectName)

	gatewayLine := utils.BashStringCmd("docker network inspect " + n + " | grep Gateway")

	return strings.Split(gatewayLine, ": ")[1]
}
