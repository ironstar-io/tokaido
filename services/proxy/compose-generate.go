package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/version"

	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// GenerateProxyDockerCompose ...
func GenerateProxyDockerCompose() {
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

		if projectInCompose(dc.Services.Proxy.Networks, nn) == true {
			return
		}
	}

	dc.Services.Proxy.Networks = append(dc.Services.Proxy.Networks, nn)

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
