package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// GenerateProxyDockerCompose ...
func GenerateProxyDockerCompose() {
	dc := DockerCompose{}
	pn := conf.GetConfig().Tokaido.Project.Name
	nn := strings.Replace(pn, ".", "", -1) + "_default"

	var err error
	err = yaml.Unmarshal(ComposeDefaults(), &dc)
	if err != nil {
		log.Fatalf("Error setting Compose file defaults for Proxy service: %v", err)
	}

	dcf := getExistingCompose()
	if dcf != nil {
		err = yaml.Unmarshal(dcf, &dc)
		if err != nil {
			log.Fatalf("Error setting Compose file for Proxy service: %v", err)
		}

		if projectInCompose(dc.Services.Proxy.Networks, nn) == true {
			utils.DebugString("[proxy] existing network configuration found, not adding it again")

			// Write the updated config to file and then move on
			n := buildNetworks(dc.Services.Proxy.Networks)
			cy, err := yaml.Marshal(&dc)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			b := append(cy, n...)
			fs.Replace(getComposePath(), b)

			return
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
