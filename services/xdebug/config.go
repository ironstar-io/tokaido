package xdebug

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var tokComposePath = filepath.Join(fs.WorkDir(), "/docker-compose.tok.yml")

// Configure - Amend docker-compose.tok.yml to be compatible with xdebug
func Configure() {
	xdebugPort := conf.GetConfig().System.Xdebug.Port
	if xdebugPort != "" {
		networkUp := docker.CheckNetworkUp()

		regenerateGateway(networkUp, xdebugPort)
	}
}

// regenerateGateway - Regenerate the docker-compose.tok.yml file and change the IP value in
func regenerateGateway(networkUp bool, xdebugPort string) {
	docker.HardCheckTokCompose()

	if networkUp == false {
		docker.Up()
	}

	networkGateway := docker.GetGateway()

	tokStruct := docker.UnmarshalledDefaults()

	err := yaml.Unmarshal(xdebugFpmEnv(networkGateway, xdebugPort), &tokStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	tokComposeYml, err := yaml.Marshal(&tokStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	docker.CreateOrReplaceTokCompose(tokComposeYml)

	docker.KillContainer("fpm")
	docker.UpContainer("fpm")
}

func xdebugFpmEnv(networkGateway string, port string) []byte {
	return []byte(`
services:
  fpm:
    environment:
      PHP_DISPLAY_ERRORS: "yes"
      XDEBUG_REMOTE_ENABLE: "yes"
      XDEBUG_REMOTE_AUTOSTART: "yes"
      XDEBUG_REMOTE_HOST: ` + networkGateway + `
      XDEBUG_REMOTE_PORT: "` + port + `"`)
}
