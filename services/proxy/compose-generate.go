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

	dc.Services.Proxy.Networks = append(dc.Services.Proxy.Networks, nn)

	n := buildNetworks(dc.Services.Proxy.Networks)

	cy, err := yaml.Marshal(&dc)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	b := append(cy, n...)

	fs.Replace(getComposePath(), b)
}
