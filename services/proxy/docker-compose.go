package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/version"

	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// GenerateProxyDockerCompose ...
func GenerateProxyDockerCompose() {
	dc := DockerCompose{}
	uv := version.GetUnisonVersion()
	pn := conf.GetConfig().Tokaido.Project.Name

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

		if projectInCompose(dc.Services.Proxy.Networks, pn) == true {
			return
		}
	}

	dc.Services.Proxy.Networks = append(dc.Services.Proxy.Networks, pn+"_default")

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

func getExistingCompose() []byte {
	cp := getComposePath()
	var _, err = os.Stat(cp)

	if os.IsNotExist(err) {
		return nil
	}

	dcf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your proxy docker-compose.yml file\n%v", err)
	}

	return dcf
}

func projectInCompose(networks []string, projectName string) bool {
	for _, v := range networks {
		if v == projectName+"_default" {
			return true
		}
	}

	return false
}

func buildNetworks(networks []string) []byte {
	n := `networks:
  proxy:
`

	for _, v := range networks {
		if v == "proxy" {
			continue
		}

		n = n + `  ` + v + `:
    external: true
`
	}

	return []byte(n)
}

// SetUnisonVersion ...
func SetUnisonVersion(version string) []byte {
	return []byte(`services:
  unison:
    image: tokaido/unison:` + version)
}
