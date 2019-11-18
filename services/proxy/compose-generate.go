package proxy

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"log"

	yaml "gopkg.in/yaml.v2"
)

// bootstrapProxy generates a docker-compose.yml file and starts the proxy, but only if there isn't
// an instance of proxy already running
func bootstrapProxy() {
	// Identify if the proxy container is running
	proxy := docker.GetContainer("proxy", "proxy")
	if proxy.State == "running" {
		// proxy is already running, so nothing to do
		return
	}

	// generate a new proxy compose file
	generateProxyDockerCompose()

	// start the proxy container
	DockerComposeUp()

}

// generateProxyDockerCompose ...
func generateProxyDockerCompose() {
	utils.DebugString("generating proxy compose file")
	dc := DockerCompose{}

	if conf.GetConfig().Global.Syncservice == "unison" {
		err := yaml.Unmarshal(ComposeDefaultsUnison(), &dc)
		if err != nil {
			fmt.Println(string(ComposeDefaultsUnison()))
			log.Fatalf("Error setting Unison's Compose file defaults for Proxy service: %v", err)
		}
	} else {
		// Read in the default docker-compose.yml file for proxy
		err := yaml.Unmarshal(ComposeDefaults(), &dc)
		if err != nil {
			log.Fatalf("Error setting standard Compose file defaults for Proxy service: %v", err)
		}
	}

	n := []byte(`networks:
  default:
  tokaido_proxy:
    external: true`)

	cy, err := yaml.Marshal(&dc)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	b := append(cy, n...)

	fs.Replace(getComposePath(), b)
}
