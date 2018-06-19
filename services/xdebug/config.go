package xdebug

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"

	"fmt"
)

// Configure - Amend docker-compose.tok.yml to be compatible with xdebug
func Configure() {
	networkUp := docker.CheckNetworkUp()

	regenerateGateway(networkUp)
}

// regenerateGateway - Regenerate the docker-compose.tok.yml file and change the IP value in
func regenerateGateway(networkUp bool) {
	docker.HardCheckTokCompose()

	if networkUp == false {
		docker.Up()
	}

	networkGateway := docker.GetGateway()
	fmt.Println(networkGateway)

	//TODO Regen compose.tok.yml with gateway IP set

	docker.KillContainer("fpm")
	docker.UpContainer("fpm")
}
