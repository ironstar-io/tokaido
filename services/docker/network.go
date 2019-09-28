package docker

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"

	"strings"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// CheckNetworkUp ...
func CheckNetworkUp() bool {
	_, err := utils.CommandSubSplitOutput("docker", "network", "inspect", "tokaido_proxy")
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

// CreateNetwork creates a Docker network with `name`
func CreateNetwork(name string) {
	dcli := GetAPIClient()

	_, err := dcli.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Labels: map[string]string{
			"io.tokaido.managed": "local",
			"io.tokaido.project": conf.GetConfig().Tokaido.Project.Name,
		},
	})

	if err != nil {
		utils.DebugString("silencing error encountered while creating the docker network [" + name + "]:")
		utils.DebugString(err.Error())
	}
}
