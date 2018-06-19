package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"strings"
)

// CheckNetworkUp ...
func CheckNetworkUp() bool {
	project := conf.GetConfig().Project

	_, err := utils.CommandSubSplitOutput("docker", "network", "inspect", project+"_default")
	if err != nil {
		return false
	}

	return true
}

// GetGateway - Get the Gateway IP adress of the docker network
func GetGateway() string {
	project := conf.GetConfig().Project

<<<<<<< HEAD
	gatewayLine := utils.BashStringCmd("docker network inspect " + project + "_default | grep Gateway")
=======
	gatewayLine := utils.SilentBashStringCmd("docker network inspect " + project + "_default | grep Gateway")
>>>>>>> Begin build out for debug configuration

	return strings.Split(gatewayLine, ": ")[1]
}
