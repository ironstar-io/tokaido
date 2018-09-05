package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/version"

	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// RemoveProjectFromDockerCompose ...
func RemoveProjectFromDockerCompose() {
	dc := DockerCompose{}
	uv := version.GetUnisonVersion()
	pn := conf.GetConfig().Tokaido.Project.Name
	nn := strings.Replace(pn, ".", "", -1) + "_default"

	var err error
	err = yaml.Unmarshal(ComposeDefaults(), &dc)
	if err != nil {
		log.Fatalf("Error setting Compose file defaults: %v", err)
	}

	dcf := getExistingCompose()
	if dcf != nil {
		err = yaml.Unmarshal(dcf, &dc)
		if err != nil {
			log.Fatalf("Error setting Compose file: %v", err)
		}
	}

	// Find and remove the project network
	for i, v := range dc.Services.Proxy.Networks {
		if v == nn {
			dc.Services.Proxy.Networks = append(dc.Services.Proxy.Networks[:i], dc.Services.Proxy.Networks[i+1:]...)
			break
		}
	}

	n := buildNetworks(dc.Services.Proxy.Networks)

	err = yaml.Unmarshal(SetUnisonVersion(uv), &dc)
	if err != nil {
		log.Fatalf("Error setting Unison version: %v", err)
	}

	cy, err := yaml.Marshal(&dc)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	b := append(cy, n...)

	fs.Replace(getComposePath(), b)
}
