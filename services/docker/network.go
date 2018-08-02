package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"

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
func GetGateway() string {
	project := conf.GetConfig().Tokaido.Project.Name

	gatewayLine := utils.BashStringCmd("docker network inspect " + project + "_default | grep Gateway")

	return strings.Split(gatewayLine, ": ")[1]
}
