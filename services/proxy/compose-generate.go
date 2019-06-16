package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"log"

	yaml "gopkg.in/yaml.v2"
)

// GenerateProxyDockerCompose ...
func GenerateProxyDockerCompose() {
	utils.DebugString("generating proxy compose file")
	dc := DockerCompose{}
	pn := conf.GetConfig().Tokaido.Project.Name
	nn := docker.GetNetworkName(pn)
	utils.DebugString("resolved network name is: " + nn)

	if conf.GetConfig().Global.Syncservice == "unison" {
		err := yaml.Unmarshal(ComposeDefaultsUnison(), &dc)
		if err != nil {
			log.Fatalf("Error setting Compose file defaults for Proxy service: %v", err)
		}
	} else {
		// Read in the default docker-compose.yml file for proxy
		err := yaml.Unmarshal(ComposeDefaults(), &dc)
		if err != nil {
			log.Fatalf("Error setting Compose file defaults for Proxy service: %v", err)
		}
	}

	// build the list of networks that are attached to the proxy service itself
	dc.Services.Proxy.Networks = buildProxyServiceNetworkAttachments()

	// build the list of networks as part of the base 'networks:' tree in docker compose
	n := buildProxyExternalNetworkList()

	cy, err := yaml.Marshal(&dc)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	b := append(cy, n...)

	fs.Replace(getComposePath(), b)
}
