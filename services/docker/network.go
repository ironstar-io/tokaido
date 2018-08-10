package docker

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"

	"strings"
)

// CheckNetworkUp ...
func CheckNetworkUp() bool {
	project := conf.GetConfig().Tokaido.Project.Name

	_, err := utils.CommandSubSplitOutput("docker", "network", "inspect", project+"_default")
	if err != nil {
		return false
	}

	return true
}

// GetGateway - Get the Gateway IP adress of the docker network
func GetGateway(networkName string) string {
	gatewayLine := utils.BashStringCmd("docker network inspect " + networkName + " | grep Gateway")

	return strings.Split(gatewayLine, ": ")[1]
}
