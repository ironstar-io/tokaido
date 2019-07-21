package proxy

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/version"
	yaml "gopkg.in/yaml.v2"
)

// RemoveProjectFromDockerCompose ...
func RemoveProjectFromDockerCompose() {
	dc := DockerCompose{}
	pn := conf.GetConfig().Tokaido.Project.Name
	nn := docker.GetNetworkName(pn)

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

	ss := conf.GetConfig().Global.Syncservice
	if ss == "unison" {
		uv := version.GetUnisonVersion()
		err = yaml.Unmarshal(SetUnisonVersion(uv), &dc)
		if err != nil {
			log.Fatalf("Error setting Unison version: %v", err)
		}
	}

	// Find and remove the project network if it is duplicated
	for i, v := range dc.Services.Proxy.Networks {
		if v == nn {
			dc.Services.Proxy.Networks = append(dc.Services.Proxy.Networks[:i], dc.Services.Proxy.Networks[i+1:]...)
			break
		}
	}

	n := buildProxyExternalNetworkList()

	cy, err := yaml.Marshal(&dc)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	b := append(cy, n...)

	fs.Replace(getComposePath(), b)
}
